package scorer

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

// NewScorer creates a new instance of the Scorer
func NewScorer(cfg Config) (Scorer, error) {
	if initError != nil {
		return nil, initError
	}
	
	if batchPromptError != nil {
		return nil, batchPromptError
	}
	
	if cfg.APIKey == "" {
		return nil, ErrMissingAPIKey
	}
	
	// Validate prompt template if provided
	if cfg.PromptText != "" {
		// Check for either Go template syntax or sprintf placeholder
		hasTemplate := strings.Contains(cfg.PromptText, "{{") && strings.Contains(cfg.PromptText, "}}")
		hasSprintf := strings.Contains(cfg.PromptText, "%s")
		if !hasTemplate && !hasSprintf {
			slog.Warn("Custom prompt has no placeholders, items will be appended",
				"prompt_preview", cfg.PromptText[:min(50, len(cfg.PromptText))])
		}
	}
	
	if cfg.MaxConcurrent < 0 {
		return nil, errors.New("MaxConcurrent must be non-negative")
	}
	
	// Set default MaxConcurrent if not specified
	if cfg.MaxConcurrent == 0 {
		cfg.MaxConcurrent = 1
	}

	// Set default model if not specified
	if cfg.Model == "" {
		cfg.Model = openai.GPT4oMini
	}

	// Set default timeout if not specified
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 // 30 seconds default
	}

	prompt := batchScorePrompt
	if cfg.PromptText != "" {
		prompt = cfg.PromptText
	}

	client := openai.NewClient(cfg.APIKey)
	return &scorer{
		client: client,
		config: cfg,
		prompt: prompt,
	}, nil
}

// ScoreTexts scores a slice of text items
func (s *scorer) ScoreTexts(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	return s.ScoreTextsWithOptions(ctx, items, opts...)
}

// ScoreTextsWithOptions scores text items with runtime options
func (s *scorer) ScoreTextsWithOptions(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	if items == nil {
		return nil, errors.New("items cannot be nil")
	}
	
	if len(items) == 0 {
		return []ScoredItem{}, nil
	}
	
	// Validate items
	for i, item := range items {
		if item.ID == "" {
			return nil, fmt.Errorf("item at index %d has empty ID", i)
		}
		if item.Content == "" {
			slog.Warn("Item has empty content", "item_id", item.ID, "index", i)
		}
	}

	// Apply default options
	options := &scoringOptions{
		model: s.config.Model,
	}
	
	// Apply provided options
	for _, opt := range opts {
		opt(options)
	}

	// Create batches
	var batches [][]TextItem
	for i := 0; i < len(items); i += maxBatchSize {
		batch := items[i:min(i+maxBatchSize, len(items))]
		batches = append(batches, batch)
	}

	// Process batches based on MaxConcurrent setting
	if s.config.MaxConcurrent <= 1 {
		return s.processSequentially(ctx, batches, options)
	}
	return s.processConcurrently(ctx, batches, options)
}

// GetHealth returns the current health status of the scorer
func (s *scorer) GetHealth(ctx context.Context) HealthStatus {
	// Basic health check - attempt a simple API call
	testItem := []TextItem{
		{ID: "health-check", Content: "test"},
	}
	
	_, err := s.ScoreTexts(ctx, testItem)
	if err != nil {
		return HealthStatus{
			Healthy: false,
			Status:  "unhealthy",
			Details: map[string]interface{}{
				"error":    err.Error(),
				"api_key":  s.config.APIKey != "",
				"model":    s.config.Model,
			},
		}
	}
	
	return HealthStatus{
		Healthy: true,
		Status:  "healthy",
		Details: map[string]interface{}{
			"api_status":      "connected",
			"model":           s.config.Model,
			"max_concurrent":  s.config.MaxConcurrent,
			"circuit_breaker": s.config.EnableCircuitBreaker,
			"retry_enabled":   s.config.EnableRetry,
		},
	}
}

func (s *scorer) processSequentially(ctx context.Context, batches [][]TextItem, options *scoringOptions) ([]ScoredItem, error) {
	var allResults []ScoredItem
	for i, batch := range batches {
		results, err := s.processBatch(ctx, batch, options)
		if err != nil {
			return nil, fmt.Errorf("processing batch %d: %w", i, err)
		}
		allResults = append(allResults, results...)
	}

	slog.Info("All items scored successfully",
		"total_items", len(allResults),
		"total_batches", len(batches),
		"mode", "sequential")

	return allResults, nil
}

func (s *scorer) processConcurrently(ctx context.Context, batches [][]TextItem, options *scoringOptions) ([]ScoredItem, error) {
	type batchResult struct {
		index   int
		results []ScoredItem
		err     error
	}

	// Semaphore to limit concurrent processing
	sem := make(chan struct{}, s.config.MaxConcurrent)
	results := make(chan batchResult, len(batches))

	// Process batches concurrently
	for i, batch := range batches {
		go func(index int, batch []TextItem) {
			sem <- struct{}{}        // Acquire semaphore
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
	allResults := make([][]ScoredItem, len(batches))
	for i := 0; i < len(batches); i++ {
		result := <-results
		if result.err != nil {
			return nil, fmt.Errorf("processing batch %d: %w", result.index, result.err)
		}
		allResults[result.index] = result.results
	}

	// Flatten results
	var flatResults []ScoredItem
	for _, batchResults := range allResults {
		flatResults = append(flatResults, batchResults...)
	}

	slog.Info("All items scored successfully",
		"total_items", len(flatResults),
		"total_batches", len(batches),
		"mode", "concurrent",
		"max_concurrent", s.config.MaxConcurrent)

	return flatResults, nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

