package scorer

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
)

//go:embed prompts/*.txt
var promptFS embed.FS

var systemPrompt string
var initError error

func init() {
	// Load system prompt during package initialization
	promptBytes, err := promptFS.ReadFile("prompts/system_prompt.txt")
	if err != nil {
		initError = fmt.Errorf("failed to load system prompt: %w", err)
		return
	}
	systemPrompt = string(promptBytes)
}

// New creates a new instance of the Scorer
func New(cfg Config) (Scorer, error) {
	if initError != nil {
		return nil, initError
	}
	
	if batchPromptError != nil {
		return nil, batchPromptError
	}
	
	if cfg.OpenAIKey == "" {
		return nil, ErrMissingAPIKey
	}
	
	if cfg.PromptText != "" && !strings.Contains(cfg.PromptText, "%s") {
		return nil, errors.New("custom prompt must contain %s placeholder for posts")
	}
	
	if cfg.MaxConcurrent < 0 {
		return nil, errors.New("MaxConcurrent must be non-negative")
	}
	
	// Set default MaxConcurrent if not specified
	if cfg.MaxConcurrent == 0 {
		cfg.MaxConcurrent = 1
	}

	prompt := batchScorePrompt
	if cfg.PromptText != "" {
		prompt = cfg.PromptText
	}

	client := openai.NewClient(cfg.OpenAIKey)
	return &scorer{
		client: client,
		config: cfg,
		prompt: prompt,
	}, nil
}


// WithPrompt returns a function that sets a custom prompt for the scorer
func WithPrompt(prompt string) func(*scorer) {
	return func(s *scorer) {
		s.prompt = prompt
	}
}

// WithMaxConcurrent returns a function that sets MaxConcurrent for the scorer
func WithMaxConcurrent(maxConcurrent int) func(*scorer) {
	return func(s *scorer) {
		s.config.MaxConcurrent = maxConcurrent
	}
}

// NewWithClient creates a new scorer with a custom OpenAI client and options
func NewWithClient(client OpenAIClient, opts ...func(*scorer)) Scorer {
	s := &scorer{
		client: client,
		prompt: batchScorePrompt,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// NewWithClientAndOptions creates a new scorer with a custom OpenAI client and ScoringOptions
func NewWithClientAndOptions(client OpenAIClient, prompt string, maxConcurrent int) Scorer {
	s := &scorer{
		client: client,
		prompt: prompt,
		config: Config{
			MaxConcurrent: maxConcurrent,
		},
	}
	if s.prompt == "" {
		s.prompt = batchScorePrompt
	}
	if s.config.MaxConcurrent == 0 {
		s.config.MaxConcurrent = 1
	}
	return s
}

// ScorePosts evaluates and scores a slice of Reddit posts
func (s *scorer) ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error) {
	return s.ScorePostsWithOptions(ctx, posts)
}

// ScorePostsWithOptions evaluates and scores a slice of Reddit posts with options
func (s *scorer) ScorePostsWithOptions(ctx context.Context, posts []*reddit.Post, opts ...ScoringOption) ([]*ScoredPost, error) {
	if posts == nil {
		return nil, errors.New("posts cannot be nil")
	}
	
	if len(posts) == 0 {
		return []*ScoredPost{}, nil
	}
	
	for i, post := range posts {
		if post == nil {
			return nil, fmt.Errorf("post at index %d is nil", i)
		}
		if post.ID == "" {
			return nil, fmt.Errorf("post at index %d has empty ID", i)
		}
	}

	// Apply default options
	options := &scoringOptions{
		model: s.config.Model, // Default from config
	}
	
	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	// Create batches
	var batches [][]*reddit.Post
	for i := 0; i < len(posts); i += maxBatchSize {
		batch := posts[i:min(i+maxBatchSize, len(posts))]
		batches = append(batches, batch)
	}

	// Process batches based on MaxConcurrent setting
	if s.config.MaxConcurrent <= 1 {
		return s.processSequentially(ctx, batches, options)
	}
	return s.processConcurrently(ctx, batches, options)
}

func (s *scorer) processSequentially(ctx context.Context, batches [][]*reddit.Post, options *scoringOptions) ([]*ScoredPost, error) {
	var allResults []*ScoredPost
	for i, batch := range batches {
		results, err := s.processBatch(ctx, batch, options)
		if err != nil {
			return nil, fmt.Errorf("processing batch %d: %w", i, err)
		}
		allResults = append(allResults, results...)
	}

	slog.Info("All posts scored successfully",
		"total_posts", len(allResults),
		"total_batches", len(batches),
		"mode", "sequential")

	return allResults, nil
}

func (s *scorer) processConcurrently(ctx context.Context, batches [][]*reddit.Post, options *scoringOptions) ([]*ScoredPost, error) {
	type batchResult struct {
		index   int
		results []*ScoredPost
		err     error
	}

	// Semaphore to limit concurrent processing
	sem := make(chan struct{}, s.config.MaxConcurrent)
	results := make(chan batchResult, len(batches))

	// Process batches concurrently
	for i, batch := range batches {
		go func(index int, batch []*reddit.Post) {
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			batchResults, err := s.processBatch(ctx, batch, options)
			results <- batchResult{
				index:   index,
				results: batchResults,
				err:     err,
			}
		}(i, batch)
	}

	// Collect results in order
	allResults := make([][]*ScoredPost, len(batches))
	for i := 0; i < len(batches); i++ {
		result := <-results
		if result.err != nil {
			return nil, fmt.Errorf("processing batch %d: %w", result.index, result.err)
		}
		allResults[result.index] = result.results
	}

	// Flatten results
	var flatResults []*ScoredPost
	for _, batchResults := range allResults {
		flatResults = append(flatResults, batchResults...)
	}

	slog.Info("All posts scored successfully",
		"total_posts", len(flatResults),
		"total_batches", len(batches),
		"mode", "concurrent",
		"max_concurrent", s.config.MaxConcurrent)

	return flatResults, nil
}

// ScorePostsWithContext evaluates and scores posts with additional context
func (s *scorer) ScorePostsWithContext(ctx context.Context, contexts []ScoringContext, opts ...ScoringOption) ([]*ScoredPost, error) {
	if contexts == nil {
		return nil, errors.New("contexts cannot be nil")
	}
	
	if len(contexts) == 0 {
		return []*ScoredPost{}, nil
	}
	
	// Extract posts from contexts for validation
	posts := make([]*reddit.Post, len(contexts))
	for i, context := range contexts {
		if context.Post == nil {
			return nil, fmt.Errorf("context at index %d has nil post", i)
		}
		if context.Post.ID == "" {
			return nil, fmt.Errorf("post at index %d has empty ID", i)
		}
		posts[i] = context.Post
	}

	// Apply default options
	options := &scoringOptions{
		model: s.config.Model, // Default from config
	}
	
	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	// Create batches of contexts
	var batches [][]ScoringContext
	for i := 0; i < len(contexts); i += maxBatchSize {
		batch := contexts[i:min(i+maxBatchSize, len(contexts))]
		batches = append(batches, batch)
	}

	// Process batches based on MaxConcurrent setting
	if s.config.MaxConcurrent <= 1 {
		return s.processContextsSequentially(ctx, batches, options)
	}
	return s.processContextsConcurrently(ctx, batches, options)
}

func (s *scorer) processContextsSequentially(ctx context.Context, batches [][]ScoringContext, options *scoringOptions) ([]*ScoredPost, error) {
	var allResults []*ScoredPost
	for i, batch := range batches {
		// Extract posts from contexts
		posts := make([]*reddit.Post, len(batch))
		for j, context := range batch {
			posts[j] = context.Post
		}
		
		results, err := s.processBatchWithContext(ctx, posts, batch, options)
		if err != nil {
			return nil, fmt.Errorf("processing batch %d: %w", i, err)
		}
		allResults = append(allResults, results...)
	}

	slog.Info("All posts scored successfully with context",
		"total_posts", len(allResults),
		"total_batches", len(batches),
		"mode", "sequential")

	return allResults, nil
}

func (s *scorer) processContextsConcurrently(ctx context.Context, batches [][]ScoringContext, options *scoringOptions) ([]*ScoredPost, error) {
	type batchResult struct {
		index   int
		results []*ScoredPost
		err     error
	}

	// Semaphore to limit concurrent processing
	sem := make(chan struct{}, s.config.MaxConcurrent)
	results := make(chan batchResult, len(batches))

	// Process batches concurrently
	for i, batch := range batches {
		go func(index int, batch []ScoringContext) {
			sem <- struct{}{} // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			// Extract posts from contexts
			posts := make([]*reddit.Post, len(batch))
			for j, context := range batch {
				posts[j] = context.Post
			}
			
			batchResults, err := s.processBatchWithContext(ctx, posts, batch, options)
			results <- batchResult{
				index:   index,
				results: batchResults,
				err:     err,
			}
		}(i, batch)
	}

	// Collect results in order
	allResults := make([][]*ScoredPost, len(batches))
	for i := 0; i < len(batches); i++ {
		result := <-results
		if result.err != nil {
			return nil, fmt.Errorf("processing batch %d: %w", result.index, result.err)
		}
		allResults[result.index] = result.results
	}

	// Flatten results
	var flatResults []*ScoredPost
	for _, batchResults := range allResults {
		flatResults = append(flatResults, batchResults...)
	}

	slog.Info("All posts scored successfully with context",
		"total_posts", len(flatResults),
		"total_batches", len(batches),
		"mode", "concurrent",
		"max_concurrent", s.config.MaxConcurrent)

	return flatResults, nil
}
