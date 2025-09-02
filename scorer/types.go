package scorer

import (
	"context"
	"errors"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sony/gobreaker/v2"
)

// TextItem represents a generic text item to be scored
type TextItem struct {
	ID       string                 // Unique identifier for the text item
	Content  string                 // The text content to be scored
	Metadata map[string]interface{} // Optional metadata for context
}

// ScoredItem represents a text item with its AI-generated score
type ScoredItem struct {
	Item   TextItem // Original text item
	Score  int      // Score between 0-100
	Reason string   // AI explanation for the score
}

// Scorer provides methods to score generic text items
type Scorer interface {
	// ScoreTexts scores a slice of text items
	ScoreTexts(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error)
	
	// ScoreTextsWithOptions scores text items with runtime options
	ScoreTextsWithOptions(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error)
	
	// GetHealth returns the current health status of the scorer
	GetHealth(ctx context.Context) HealthStatus
}

// HealthStatus represents the health state of the scorer
type HealthStatus struct {
	Healthy bool                   // Overall health status
	Status  string                 // Human-readable status message
	Details map[string]interface{} // Additional health details
}

// Config holds the configuration for the scorer
type Config struct {
	APIKey               string                // OpenAI API key (required)
	Model                string                // OpenAI model to use
	PromptText           string                // Custom prompt template
	MaxConcurrent        int                   // Maximum concurrent API calls
	MaxContentLength     int                   // Maximum content length per text item (0 = use default)
	EnableCircuitBreaker bool                  // Enable circuit breaker pattern
	EnableRetry          bool                  // Enable retry with backoff
	Timeout              time.Duration         // Request timeout
	CircuitBreakerConfig *CircuitBreakerConfig // Circuit breaker configuration
	RetryConfig          *RetryConfig          // Retry configuration
}

// CircuitBreakerConfig holds circuit breaker settings
type CircuitBreakerConfig struct {
	MaxRequests   uint32                                          // Max requests in half-open state
	Interval      time.Duration                                   // Interval for closed state
	Timeout       time.Duration                                   // Timeout for open state
	ReadyToTrip   func(counts gobreaker.Counts) bool            // Custom trip condition
	OnStateChange func(name string, from, to gobreaker.State)    // State change callback
}

// RetryConfig holds retry settings
type RetryConfig struct {
	MaxAttempts  int           // Maximum number of retry attempts
	Strategy     RetryStrategy // Backoff strategy to use
	InitialDelay time.Duration // Initial delay between retries
	MaxDelay     time.Duration // Maximum delay between retries
}

// RetryStrategy defines the backoff strategy for retries
type RetryStrategy string

const (
	RetryStrategyExponential RetryStrategy = "exponential"
	RetryStrategyConstant    RetryStrategy = "constant"
	RetryStrategyFibonacci   RetryStrategy = "fibonacci"
	
	// Content length limits
	DefaultMaxContentLength = 10000 // Default maximum content length in characters
	MinContentLength        = 1     // Minimum content length to be valid
)

// OpenAIClient defines the interface for interacting with OpenAI API
type OpenAIClient interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// Internal scorer implementation
type scorer struct {
	client OpenAIClient
	config Config
	prompt string
}

// Error definitions
var (
	ErrMissingAPIKey     = errors.New("OpenAI API key is required")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrEmptyInput        = errors.New("input items cannot be empty")
	ErrContentTooLong    = errors.New("content exceeds maximum length")
	ErrContentTooShort   = errors.New("content is too short")
	ErrContentWhitespace = errors.New("content contains only whitespace")
)

// Internal response types for JSON parsing
type scoreResponse struct {
	Version string      `json:"version"`
	Scores  []scoreItem `json:"scores"`
}

type scoreItem struct {
	ItemID string `json:"item_id"`
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

// ScoringOption is a functional option for configuring scoring behavior
type ScoringOption func(*scoringOptions)

// scoringOptions holds the options for a scoring request (internal)
type scoringOptions struct {
	model        string                 // Model to use for this request
	promptText   string                 // Custom prompt for this request
	extraContext map[string]interface{} // Additional context data
}

// ScoringOptions is the exported version for testing (uppercase)
type ScoringOptions struct {
	Model        string                 // Model to use for this request
	PromptText   string                 // Custom prompt for this request
	ExtraContext map[string]interface{} // Additional context data
}

// WithModel sets the model for this scoring request
func WithModel(model string) ScoringOption {
	return func(opts *scoringOptions) {
		opts.model = model
	}
}

// WithPromptTemplate sets a custom prompt template for this scoring request
func WithPromptTemplate(prompt string) ScoringOption {
	return func(opts *scoringOptions) {
		opts.promptText = prompt
	}
}

// WithExtraContext adds extra context data for template substitution
func WithExtraContext(context map[string]interface{}) ScoringOption {
	return func(opts *scoringOptions) {
		opts.extraContext = context
	}
}

