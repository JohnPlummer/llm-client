package scorer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sony/gobreaker/v2"
)

// CircuitBreakerWrapper wraps an OpenAI client with circuit breaker functionality
type CircuitBreakerWrapper struct {
	client OpenAIClient
	cb     *gobreaker.CircuitBreaker[openai.ChatCompletionResponse]
}

// NewCircuitBreakerWrapper creates a new circuit breaker wrapper around an OpenAI client
func NewCircuitBreakerWrapper(client OpenAIClient, config *CircuitBreakerConfig) *CircuitBreakerWrapper {
	if config == nil {
		// Default configuration
		config = &CircuitBreakerConfig{
			MaxRequests: 10,
			Interval:    60 * time.Second,
			Timeout:     30 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.ConsecutiveFailures >= 5 || 
					   (counts.Requests >= 10 && failureRatio > 0.6)
			},
		}
	}

	settings := gobreaker.Settings{
		Name:        "openai-api",
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: config.ReadyToTrip,
		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Warn("Circuit breaker state changed",
				"name", name,
				"from", from.String(),
				"to", to.String())
			
			if config.OnStateChange != nil {
				config.OnStateChange(name, from, to)
			}
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			
			// Don't count rate limits and timeouts as circuit breaker failures
			// These are temporary and should be retried
			return !ShouldTripCircuit(err)
		},
	}

	cb := gobreaker.NewCircuitBreaker[openai.ChatCompletionResponse](settings)

	return &CircuitBreakerWrapper{
		client: client,
		cb:     cb,
	}
}

// CreateChatCompletion executes the API call through the circuit breaker
func (w *CircuitBreakerWrapper) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	resp, err := w.cb.Execute(func() (openai.ChatCompletionResponse, error) {
		return w.client.CreateChatCompletion(ctx, req)
	})

	if err != nil {
		// Log the error with context
		if errors.Is(err, gobreaker.ErrOpenState) {
			slog.Debug("Circuit breaker is open, request rejected",
				"error", err)
		} else if errors.Is(err, gobreaker.ErrTooManyRequests) {
			slog.Debug("Circuit breaker in half-open state, too many requests",
				"error", err)
		} else {
			slog.Debug("Request failed through circuit breaker",
				"error", err,
				"should_trip", ShouldTripCircuit(err))
		}
	}

	return resp, err
}

// State returns the current state of the circuit breaker
func (w *CircuitBreakerWrapper) State() gobreaker.State {
	return w.cb.State()
}

// Counts returns the current counts of the circuit breaker
func (w *CircuitBreakerWrapper) Counts() gobreaker.Counts {
	return w.cb.Counts()
}

// GetHealth returns the health status of the circuit breaker
func (w *CircuitBreakerWrapper) GetHealth() HealthStatus {
	state := w.cb.State()
	counts := w.cb.Counts()
	
	var healthy bool
	var status string
	
	switch state {
	case gobreaker.StateClosed:
		healthy = true
		status = "closed"
	case gobreaker.StateHalfOpen:
		healthy = true // Degraded but operational
		status = "half-open"
	case gobreaker.StateOpen:
		healthy = false
		status = "open"
	default:
		status = "unknown"
	}

	details := map[string]interface{}{
		"state":                state.String(),
		"requests":             counts.Requests,
		"total_successes":      counts.TotalSuccesses,
		"total_failures":       counts.TotalFailures,
		"consecutive_failures": counts.ConsecutiveFailures,
		"consecutive_successes": counts.ConsecutiveSuccesses,
	}

	return HealthStatus{
		Healthy: healthy,
		Status:  status,
		Details: details,
	}
}

// ShouldTripCircuit determines if an error should cause the circuit to trip
func ShouldTripCircuit(err error) bool {
	if err == nil {
		return false
	}

	// Check for OpenAI API errors
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.HTTPStatusCode {
		case 429: // Rate limit - don't trip, this is expected
			return false
		case 401, 403: // Auth errors - trip immediately
			return true
		case 500, 502, 503, 504: // Server errors - trip
			return true
		default:
			// Other 4xx errors might be client issues, trip to be safe
			if apiErr.HTTPStatusCode >= 400 && apiErr.HTTPStatusCode < 500 {
				return true
			}
			// 5xx errors should trip
			if apiErr.HTTPStatusCode >= 500 {
				return true
			}
		}
	}

	// Check for timeout errors - don't trip on timeouts
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}

	// Unknown errors should trip the circuit
	return true
}

// WrapWithCircuitBreaker wraps an existing TextScorer with circuit breaker functionality
func WrapWithCircuitBreaker(scorer TextScorer, config *CircuitBreakerConfig) TextScorer {
	// This would require extracting the client from the scorer
	// For now, this is a placeholder for future enhancement
	slog.Info("Circuit breaker wrapper for TextScorer not yet implemented")
	return scorer
}

// circuitBreakerScorer wraps a TextScorer with circuit breaker functionality
type circuitBreakerScorer struct {
	scorer TextScorer
	cb     *gobreaker.CircuitBreaker[[]ScoredItem]
	config *CircuitBreakerConfig
}

// NewCircuitBreakerScorer creates a new circuit breaker wrapper for a TextScorer
func NewCircuitBreakerScorer(scorer TextScorer, config *CircuitBreakerConfig) TextScorer {
	if config == nil {
		config = &CircuitBreakerConfig{
			MaxRequests: 10,
			Interval:    60 * time.Second,
			Timeout:     30 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
				return counts.ConsecutiveFailures >= 5 || 
					   (counts.Requests >= 10 && failureRatio > 0.6)
			},
		}
	}

	settings := gobreaker.Settings{
		Name:        "text-scorer",
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: config.ReadyToTrip,
		OnStateChange: func(name string, from, to gobreaker.State) {
			slog.Warn("Text scorer circuit breaker state changed",
				"name", name,
				"from", from.String(),
				"to", to.String())
			
			if config.OnStateChange != nil {
				config.OnStateChange(name, from, to)
			}
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			return !ShouldTripCircuit(err)
		},
	}

	cb := gobreaker.NewCircuitBreaker[[]ScoredItem](settings)

	return &circuitBreakerScorer{
		scorer: scorer,
		cb:     cb,
		config: config,
	}
}

// ScoreTexts implements TextScorer interface with circuit breaker
func (s *circuitBreakerScorer) ScoreTexts(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	return s.cb.Execute(func() ([]ScoredItem, error) {
		return s.scorer.ScoreTexts(ctx, items, opts...)
	})
}

// ScoreTextsWithOptions implements TextScorer interface with circuit breaker
func (s *circuitBreakerScorer) ScoreTextsWithOptions(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	return s.cb.Execute(func() ([]ScoredItem, error) {
		return s.scorer.ScoreTextsWithOptions(ctx, items, opts...)
	})
}

// GetHealth implements TextScorer interface
func (s *circuitBreakerScorer) GetHealth(ctx context.Context) HealthStatus {
	state := s.cb.State()
	counts := s.cb.Counts()
	
	baseHealth := s.scorer.GetHealth(ctx)
	
	// Merge circuit breaker status with base health
	baseHealth.Details["circuit_breaker_state"] = state.String()
	baseHealth.Details["circuit_breaker_requests"] = counts.Requests
	baseHealth.Details["circuit_breaker_failures"] = counts.TotalFailures
	
	// Override health if circuit is open
	if state == gobreaker.StateOpen {
		baseHealth.Healthy = false
		baseHealth.Status = fmt.Sprintf("circuit open (%s)", baseHealth.Status)
	} else if state == gobreaker.StateHalfOpen {
		baseHealth.Status = fmt.Sprintf("degraded (%s)", baseHealth.Status)
	}
	
	return baseHealth
}