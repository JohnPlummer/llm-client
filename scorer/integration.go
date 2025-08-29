package scorer

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sony/gobreaker/v2"
)

// IntegratedScorer combines all resilience patterns and features
type IntegratedScorer struct {
	baseScorer TextScorer
	metrics    *MetricsRecorder
	config     Config
}

// NewIntegratedScorer creates a fully integrated scorer with all features
func NewIntegratedScorer(cfg Config) (TextScorer, error) {
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Create base scorer
	baseScorer, err := NewTextScorer(cfg)
	if err != nil {
		return nil, err
	}

	// Apply resilience patterns based on configuration
	var scorer TextScorer = baseScorer

	// Layer 1: Add retry logic (innermost)
	if cfg.EnableRetry {
		slog.Info("Enabling retry logic",
			"max_attempts", cfg.RetryConfig.MaxAttempts,
			"strategy", cfg.RetryConfig.Strategy)
		scorer = NewRetryScorer(scorer, cfg.RetryConfig)
	}

	// Layer 2: Add circuit breaker (wraps retry)
	if cfg.EnableCircuitBreaker {
		slog.Info("Enabling circuit breaker",
			"max_requests", cfg.CircuitBreakerConfig.MaxRequests,
			"timeout", cfg.CircuitBreakerConfig.Timeout)
		
		// Add metrics callback to circuit breaker
		if cfg.CircuitBreakerConfig.OnStateChange == nil {
			cfg.CircuitBreakerConfig.OnStateChange = func(name string, from, to gobreaker.State) {
				metrics := NewMetricsRecorder(true)
				metrics.RecordCircuitBreakerState(name, stateToInt(to))
				if to == gobreaker.StateOpen {
					metrics.RecordCircuitBreakerTrip(name)
				}
			}
		}
		
		scorer = NewCircuitBreakerScorer(scorer, cfg.CircuitBreakerConfig)
	}

	// Create integrated scorer with metrics
	integrated := &IntegratedScorer{
		baseScorer: scorer,
		metrics:    NewMetricsRecorder(true),
		config:     cfg,
	}

	slog.Info("Integrated scorer created",
		"model", cfg.Model,
		"max_concurrent", cfg.MaxConcurrent,
		"circuit_breaker", cfg.EnableCircuitBreaker,
		"retry", cfg.EnableRetry)

	return integrated, nil
}

// ScoreTexts implements TextScorer with full integration
func (s *IntegratedScorer) ScoreTexts(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	return s.ScoreTextsWithOptions(ctx, items, opts...)
}

// ScoreTextsWithOptions implements TextScorer with metrics and monitoring
func (s *IntegratedScorer) ScoreTextsWithOptions(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	start := time.Now()
	
	// Record batch size
	s.metrics.RecordBatchSize(len(items))
	
	// Track concurrent requests
	s.metrics.RecordConcurrentRequests(1)
	defer s.metrics.RecordConcurrentRequests(-1)

	// Apply options to get model
	options := &scoringOptions{
		model: s.config.Model,
	}
	for _, opt := range opts {
		opt(options)
	}
	model := options.model
	if model == "" {
		model = openai.GPT4oMini
	}

	// Call underlying scorer
	results, err := s.baseScorer.ScoreTextsWithOptions(ctx, items, opts...)
	
	// Record metrics
	duration := time.Since(start).Seconds()
	s.metrics.RecordRequestDuration(duration, model)
	
	if err != nil {
		s.metrics.RecordRequest("error", model)
		s.metrics.RecordError(classifyError(err))
		return nil, err
	}
	
	s.metrics.RecordRequest("success", model)
	s.metrics.RecordItemsScored(len(results))
	
	// Record score distribution
	for _, result := range results {
		s.metrics.RecordScore(result.Score)
	}
	
	return results, nil
}

// GetHealth returns comprehensive health status
func (s *IntegratedScorer) GetHealth(ctx context.Context) HealthStatus {
	baseHealth := s.baseScorer.GetHealth(ctx)
	
	// Add integration-specific health checks
	baseHealth.Details["integration"] = map[string]interface{}{
		"circuit_breaker_enabled": s.config.EnableCircuitBreaker,
		"retry_enabled":           s.config.EnableRetry,
		"metrics_enabled":         true,
		"model":                   s.config.Model,
		"max_concurrent":          s.config.MaxConcurrent,
	}
	
	return baseHealth
}

// BuildProductionScorer creates a production-ready scorer with all features
func BuildProductionScorer(apiKey string) (TextScorer, error) {
	cfg := NewProductionConfig(apiKey)
	return NewIntegratedScorer(cfg)
}

// BuildCustomScorer creates a scorer with custom configuration
func BuildCustomScorer(cfg Config) (TextScorer, error) {
	return NewIntegratedScorer(cfg)
}

// stateToInt converts circuit breaker state to int for metrics
func stateToInt(state gobreaker.State) int {
	switch state {
	case gobreaker.StateClosed:
		return 0
	case gobreaker.StateHalfOpen:
		return 1
	case gobreaker.StateOpen:
		return 2
	default:
		return -1
	}
}

// classifyError returns error type for metrics
func classifyError(err error) string {
	if err == nil {
		return "none"
	}
	
	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		switch {
		case apiErr.HTTPStatusCode == 429:
			return "rate_limit"
		case apiErr.HTTPStatusCode >= 500:
			return "server_error"
		case apiErr.HTTPStatusCode >= 400:
			return "client_error"
		default:
			return "api_error"
		}
	}
	
	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout"
	}
	
	if errors.Is(err, context.Canceled) {
		return "cancelled"
	}
	
	if errors.Is(err, gobreaker.ErrOpenState) {
		return "circuit_open"
	}
	
	if errors.Is(err, gobreaker.ErrTooManyRequests) {
		return "circuit_half_open"
	}
	
	return "unknown"
}

// WithMetrics wraps any TextScorer with metrics recording
func WithMetrics(scorer TextScorer, metrics *MetricsRecorder) TextScorer {
	return &metricsScorer{
		scorer:  scorer,
		metrics: metrics,
	}
}

type metricsScorer struct {
	scorer  TextScorer
	metrics *MetricsRecorder
}

func (m *metricsScorer) ScoreTexts(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	return m.ScoreTextsWithOptions(ctx, items, opts...)
}

func (m *metricsScorer) ScoreTextsWithOptions(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error) {
	start := time.Now()
	m.metrics.RecordBatchSize(len(items))
	
	results, err := m.scorer.ScoreTextsWithOptions(ctx, items, opts...)
	
	duration := time.Since(start).Seconds()
	m.metrics.RecordRequestDuration(duration, "unknown")
	
	if err != nil {
		m.metrics.RecordError(classifyError(err))
	} else {
		m.metrics.RecordItemsScored(len(results))
		for _, result := range results {
			m.metrics.RecordScore(result.Score)
		}
	}
	
	return results, err
}

func (m *metricsScorer) GetHealth(ctx context.Context) HealthStatus {
	return m.scorer.GetHealth(ctx)
}