package scorer

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
)

//go:embed prompts/*.txt
var promptFS embed.FS

var systemPrompt string

func init() {
	// Load system prompt during package initialization
	promptBytes, err := promptFS.ReadFile("prompts/system_prompt.txt")
	if err != nil {
		slog.Error("failed to load system prompt", "error", err)
		os.Exit(1)
	}
	systemPrompt = string(promptBytes)
}

// New creates a new instance of the Scorer
func New(cfg Config) (Scorer, error) {
	if cfg.OpenAIKey == "" {
		return nil, ErrMissingAPIKey
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

// ScorePosts evaluates and scores a slice of Reddit posts
func (s *scorer) ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error) {
	if len(posts) == 0 {
		return nil, nil
	}

	var allResults []*ScoredPost
	for i := 0; i < len(posts); i += maxBatchSize {
		results, err := s.processBatch(ctx, posts[i:min(i+maxBatchSize, len(posts))])
		if err != nil {
			return nil, fmt.Errorf("processing batch %d: %w", i/maxBatchSize, err)
		}
		allResults = append(allResults, results...)
	}

	slog.Info("All posts scored successfully",
		"total_posts", len(posts),
		"total_batches", (len(posts)+maxBatchSize-1)/maxBatchSize)

	return allResults, nil
}
