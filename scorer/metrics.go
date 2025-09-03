// Package scorer provides comprehensive Prometheus metrics for monitoring text scoring operations.
// This metrics system tracks performance, errors, and usage patterns across the scoring pipeline.
package scorer

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics for monitoring all aspects of text scoring operations.
// These metrics enable comprehensive observability into performance, reliability, and usage patterns.
// Metrics are organized into functional groups: requests, batching, errors, resilience, and API interactions.
var (
	// Core request metrics track overall system usage and performance
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "text_scorer_requests_total",
			Help: "Total number of text scoring requests",
		},
		[]string{"status", "model"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "text_scorer_request_duration_seconds",
			Help:    "Duration of text scoring requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"model"},
	)

	// Batch processing metrics monitor efficiency of batch operations
	batchSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "text_scorer_batch_size",
			Help:    "Size of text scoring batches",
			Buckets: []float64{1, 2, 5, 10, 20, 50, 100},
		},
	)

	itemsScored = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "text_scorer_items_scored_total",
			Help: "Total number of text items scored",
		},
	)

	// Error tracking provides insights into failure patterns and reliability
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "text_scorer_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"error_type"},
	)

	// Circuit breaker metrics monitor system resilience and fault tolerance
	circuitBreakerState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "text_scorer_circuit_breaker_state",
			Help: "Current state of circuit breaker (0=closed, 1=half-open, 2=open)",
		},
		[]string{"name"},
	)

	circuitBreakerTrips = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "text_scorer_circuit_breaker_trips_total",
			Help: "Total number of circuit breaker trips",
		},
		[]string{"name"},
	)

	// Retry mechanism metrics track system robustness under transient failures
	retryAttempts = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "text_scorer_retry_attempts",
			Help:    "Number of retry attempts per request",
			Buckets: []float64{1, 2, 3, 4, 5},
		},
	)

	retryTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "text_scorer_retry_total",
			Help: "Total number of retries by reason",
		},
		[]string{"reason"},
	)

	// OpenAI API interaction metrics monitor external service performance and costs
	apiCallDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "text_scorer_api_call_duration_seconds",
			Help:    "Duration of API calls to OpenAI",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint", "status"},
	)

	// Token usage tracking enables cost monitoring and optimization
	apiTokensUsed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "text_scorer_api_tokens_used_total",
			Help: "Total number of tokens used in API calls",
		},
		[]string{"type"}, // prompt, completion, total
	)

	// Score distribution analysis provides insights into scoring patterns and quality
	scoreDistribution = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "text_scorer_score_distribution",
			Help:    "Distribution of scores (0-100)",
			Buckets: []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		},
	)

	// Concurrency metrics monitor system load and resource utilization
	concurrentRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "text_scorer_concurrent_requests",
			Help: "Number of concurrent requests being processed",
		},
	)

	queuedRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "text_scorer_queued_requests",
			Help: "Number of requests waiting in queue",
		},
	)
)

// MetricsRecorder provides methods to record metrics with optional enablement control.
// When disabled, all recording operations become no-ops for zero performance impact.
// This design enables metrics collection to be toggled without code changes.
type MetricsRecorder struct {
	enabled bool
}

// NewMetricsRecorder creates a new metrics recorder with the specified enablement state.
// Production systems typically enable metrics, while test environments may disable them.
func NewMetricsRecorder(enabled bool) *MetricsRecorder {
	return &MetricsRecorder{enabled: enabled}
}

// RecordRequest records a request metric
func (m *MetricsRecorder) RecordRequest(status string, model string) {
	if !m.enabled {
		return
	}
	requestsTotal.WithLabelValues(status, model).Inc()
}

// RecordRequestDuration records request duration
func (m *MetricsRecorder) RecordRequestDuration(seconds float64, model string) {
	if !m.enabled {
		return
	}
	requestDuration.WithLabelValues(model).Observe(seconds)
}

// RecordBatchSize records the size of a batch
func (m *MetricsRecorder) RecordBatchSize(size int) {
	if !m.enabled {
		return
	}
	batchSize.Observe(float64(size))
}

// RecordItemsScored records the number of items scored
func (m *MetricsRecorder) RecordItemsScored(count int) {
	if !m.enabled {
		return
	}
	itemsScored.Add(float64(count))
}

// RecordError records an error
func (m *MetricsRecorder) RecordError(errorType string) {
	if !m.enabled {
		return
	}
	errorsTotal.WithLabelValues(errorType).Inc()
}

// RecordCircuitBreakerState records circuit breaker state transitions.
// State values: 0=closed (normal), 1=half-open (testing), 2=open (failing).
func (m *MetricsRecorder) RecordCircuitBreakerState(name string, state int) {
	if !m.enabled {
		return
	}
	circuitBreakerState.WithLabelValues(name).Set(float64(state))
}

// RecordCircuitBreakerTrip records when a circuit breaker transitions to open state.
// Frequent trips indicate persistent downstream service issues requiring investigation.
func (m *MetricsRecorder) RecordCircuitBreakerTrip(name string) {
	if !m.enabled {
		return
	}
	circuitBreakerTrips.WithLabelValues(name).Inc()
}

// RecordRetryAttempt records retry attempts
func (m *MetricsRecorder) RecordRetryAttempt(attempts int) {
	if !m.enabled {
		return
	}
	retryAttempts.Observe(float64(attempts))
}

// RecordRetry records a retry
func (m *MetricsRecorder) RecordRetry(reason string) {
	if !m.enabled {
		return
	}
	retryTotal.WithLabelValues(reason).Inc()
}

// RecordAPICall records an API call duration
func (m *MetricsRecorder) RecordAPICall(endpoint string, status string, seconds float64) {
	if !m.enabled {
		return
	}
	apiCallDuration.WithLabelValues(endpoint, status).Observe(seconds)
}

// RecordTokensUsed records API token consumption for cost tracking and optimization.
// TokenType should be "prompt", "completion", or "total" to categorize usage patterns.
func (m *MetricsRecorder) RecordTokensUsed(tokenType string, count int) {
	if !m.enabled {
		return
	}
	apiTokensUsed.WithLabelValues(tokenType).Add(float64(count))
}

// RecordScore records individual score values to analyze scoring distribution patterns.
// Score distribution analysis helps identify bias, quality issues, and prompt effectiveness.
func (m *MetricsRecorder) RecordScore(score int) {
	if !m.enabled {
		return
	}
	scoreDistribution.Observe(float64(score))
}

// RecordConcurrentRequests updates concurrent request count
func (m *MetricsRecorder) RecordConcurrentRequests(delta float64) {
	if !m.enabled {
		return
	}
	concurrentRequests.Add(delta)
}

// RecordQueuedRequests updates queued request count
func (m *MetricsRecorder) RecordQueuedRequests(delta float64) {
	if !m.enabled {
		return
	}
	queuedRequests.Add(delta)
}

// GetMetricsHandler returns an HTTP handler for exposing Prometheus metrics.
// Mount this handler at /metrics to enable scraping by Prometheus servers.
// The handler serves metrics in the standard Prometheus text format.
func GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}

// RegisterCustomMetrics allows registration of application-specific metrics beyond the standard set.
// Use this for domain-specific measurements that complement the built-in metrics.
// Returns error if metric name conflicts with existing registrations.
func RegisterCustomMetrics(collector prometheus.Collector) error {
	return prometheus.Register(collector)
}
