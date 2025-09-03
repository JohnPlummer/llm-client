package scorer

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sethvargo/go-retry"
)

// RetryWrapper wraps an OpenAI client with retry logic
type RetryWrapper struct {
	client OpenAIClient
	config *RetryConfig
}

// NewRetryWrapper creates a new retry wrapper around an OpenAI client
func NewRetryWrapper(client OpenAIClient, config *RetryConfig) *RetryWrapper {
	if config == nil {
		config = &RetryConfig{
			MaxAttempts:  3,
			Strategy:     RetryStrategyExponential,
			InitialDelay: 1 * time.Second,
			MaxDelay:     30 * time.Second,
		}
	}

	return &RetryWrapper{
		client: client,
		config: config,
	}
}

// CreateChatCompletion executes the API call with retry logic
func (w *RetryWrapper) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	var lastErr error
	var attempts int

	backoff := w.getBackoffStrategy()

	for {
		attempts++

		// Try the request
		resp, err := w.client.CreateChatCompletion(ctx, req)
		if err == nil {
			if attempts > 1 {
				slog.Info("Request succeeded after retry",
					"attempts", attempts)
			}
			return resp, nil
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryableError(err) {
			slog.Debug("Non-retryable error, giving up",
				"error", err,
				"attempts", attempts)
			return openai.ChatCompletionResponse{}, err
		}

		// Check if we've exceeded max attempts
		if attempts >= w.config.MaxAttempts {
			slog.Warn("Max retry attempts reached",
				"attempts", attempts,
				"error", lastErr)
			return openai.ChatCompletionResponse{}, lastErr
		}

		// Calculate next delay
		delay, stop := backoff.Next()
		if stop {
			slog.Warn("Backoff strategy stopped",
				"attempts", attempts,
				"error", lastErr)
			return openai.ChatCompletionResponse{}, lastErr
		}

		slog.Debug("Retrying request after delay",
			"attempt", attempts,
			"delay", delay,
			"error", err)

		// Wait with context awareness
		select {
		case <-ctx.Done():
			return openai.ChatCompletionResponse{}, ctx.Err()
		case <-time.After(delay):
			// Continue to next retry
		}
	}
}

// getBackoffStrategy returns the appropriate backoff strategy
func (w *RetryWrapper) getBackoffStrategy() retry.Backoff {
	switch w.config.Strategy {
	case RetryStrategyConstant:
		return retry.WithMaxRetries(
			uint64(w.config.MaxAttempts),
			retry.BackoffFunc(func() (time.Duration, bool) {
				// Add jitter to prevent thundering herd
				jitter := time.Duration(rand.Int63n(int64(w.config.InitialDelay / 10)))
				return w.config.InitialDelay + jitter, false
			}),
		)

	case RetryStrategyFibonacci:
		return retry.WithMaxRetries(
			uint64(w.config.MaxAttempts),
			retry.WithCappedDuration(
				w.config.MaxDelay,
				retry.WithJitter(
					w.config.InitialDelay/10,
					retry.NewFibonacci(w.config.InitialDelay),
				),
			),
		)

	case RetryStrategyExponential:
		fallthrough
	default:
		return retry.WithMaxRetries(
			uint64(w.config.MaxAttempts),
			retry.WithCappedDuration(
				w.config.MaxDelay,
				retry.WithJitter(
					w.config.InitialDelay/10,
					retry.NewExponential(w.config.InitialDelay),
				),
			),
		)
	}
}

// IsRetryableError determines if an error should trigger a retry
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for OpenAI API errors
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.HTTPStatusCode {
		case 429: // Rate limit - definitely retry
			return true
		case 500, 502, 503, 504: // Server errors - retry
			return true
		case 400, 401, 403, 404: // Client errors - don't retry
			return false
		default:
			// Unknown 5xx errors should be retried
			if apiErr.HTTPStatusCode >= 500 {
				return true
			}
			// Other errors shouldn't be retried
			return false
		}
	}

	// Timeout errors are retryable
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// Cancelled context is not retryable
	if errors.Is(err, context.Canceled) {
		return false
	}

	// Network errors might be retryable
	// For now, consider unknown errors as retryable
	return true
}

// retryScorer wraps a Scorer with retry functionality
type retryScorer struct {
	scorer Scorer
	config *RetryConfig
}

// NewRetryScorer creates a new retry wrapper for a Scorer
func NewRetryScorer(scorer Scorer, config *RetryConfig) Scorer {
	if config == nil {
		config = &RetryConfig{
			MaxAttempts:  3,
			Strategy:     RetryStrategyExponential,
			InitialDelay: 1 * time.Second,
			MaxDelay:     30 * time.Second,
		}
	}

	return &retryScorer{
		scorer: scorer,
		config: config,
	}
}

// ScoreTexts implements Scorer interface with retry logic
func (s *retryScorer) ScoreTexts(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	return s.retryOperation(ctx, func() ([]ScoredItem, error) {
		return s.scorer.ScoreTexts(ctx, items, opts...)
	})
}

// ScoreTextsWithOptions implements Scorer interface with retry logic
func (s *retryScorer) ScoreTextsWithOptions(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	return s.retryOperation(ctx, func() ([]ScoredItem, error) {
		return s.scorer.ScoreTextsWithOptions(ctx, items, opts...)
	})
}

// GetHealth implements Scorer interface
func (s *retryScorer) GetHealth(ctx context.Context) HealthStatus {
	// Health checks shouldn't be retried
	return s.scorer.GetHealth(ctx)
}

// retryOperation performs an operation with retry logic
func (s *retryScorer) retryOperation(ctx context.Context, operation func() ([]ScoredItem, error)) ([]ScoredItem, error) {
	var lastErr error
	var attempts int

	wrapper := &RetryWrapper{config: s.config}
	backoff := wrapper.getBackoffStrategy()

	for {
		attempts++

		// Try the operation
		result, err := operation()
		if err == nil {
			if attempts > 1 {
				slog.Info("Text scoring succeeded after retry",
					"attempts", attempts)
			}
			return result, nil
		}

		lastErr = err

		// Check if error is retryable
		if !IsRetryableError(err) {
			slog.Debug("Non-retryable error in text scoring",
				"error", err,
				"attempts", attempts)
			return nil, err
		}

		// Check if we've exceeded max attempts
		if attempts >= s.config.MaxAttempts {
			slog.Warn("Max retry attempts reached for text scoring",
				"attempts", attempts,
				"error", lastErr)
			return nil, lastErr
		}

		// Calculate next delay
		delay, stop := backoff.Next()
		if stop {
			return nil, lastErr
		}

		slog.Debug("Retrying text scoring after delay",
			"attempt", attempts,
			"delay", delay,
			"error", err)

		// Wait with context awareness
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			// Continue to next retry
		}
	}
}

// CombineWithCircuitBreaker creates a scorer with both retry and circuit breaker
func CombineWithCircuitBreaker(scorer Scorer, retryConfig *RetryConfig, cbConfig *CircuitBreakerConfig) Scorer {
	// First wrap with retry (inner layer)
	withRetry := NewRetryScorer(scorer, retryConfig)

	// Then wrap with circuit breaker (outer layer)
	withCB := NewCircuitBreakerScorer(withRetry, cbConfig)

	return withCB
}

// CalculateRetryDelay calculates the delay for a given retry attempt
func CalculateRetryDelay(attempt int, config *RetryConfig) time.Duration {
	if config == nil {
		return 0
	}

	var delay time.Duration

	switch config.Strategy {
	case RetryStrategyConstant:
		delay = config.InitialDelay

	case RetryStrategyFibonacci:
		// Calculate fibonacci number
		a, b := config.InitialDelay, config.InitialDelay
		for i := 2; i <= attempt; i++ {
			a, b = b, a+b
		}
		delay = b

	case RetryStrategyExponential:
		fallthrough
	default:
		// 2^(attempt-1) * InitialDelay
		multiplier := 1 << (attempt - 1)
		delay = time.Duration(multiplier) * config.InitialDelay
	}

	// Cap at MaxDelay
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	// Add jitter (±10%)
	jitter := time.Duration(rand.Int63n(int64(delay / 10)))
	if rand.Intn(2) == 0 {
		delay += jitter
	} else {
		delay -= jitter
	}

	return delay
}

// GetRetryStats returns statistics about retry operations
func GetRetryStats(err error) (attempts int, finalError error) {
	// This would be enhanced with actual retry tracking
	// For now, return basic info
	return 1, err
}
