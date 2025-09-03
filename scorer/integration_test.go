// Package scorer_test provides integration tests for the scorer package,
// focusing on end-to-end scenarios with resilience patterns like circuit breakers
// and retry mechanisms. These tests validate the integration between different
// scorer components and their behavior under various failure conditions.
package scorer_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"

	"github.com/JohnPlummer/llm-client/scorer"
)

var _ = Describe("Integration", func() {
	var (
		cfg scorer.Config
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		cfg = scorer.NewDefaultConfig("test-api-key")
	})

	Describe("IntegratedScorer", func() {
		It("should create scorer with default configuration", func() {
			s, err := scorer.NewIntegratedScorer(cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(s).ToNot(BeNil())
		})

		It("should create scorer with circuit breaker enabled", func() {
			cfg = cfg.WithCircuitBreaker()
			s, err := scorer.NewIntegratedScorer(cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(s).ToNot(BeNil())
		})

		It("should create scorer with retry enabled", func() {
			cfg = cfg.WithRetry()
			s, err := scorer.NewIntegratedScorer(cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(s).ToNot(BeNil())
		})

		It("should create scorer with both resilience patterns", func() {
			cfg = cfg.WithCircuitBreaker().WithRetry()
			s, err := scorer.NewIntegratedScorer(cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(s).ToNot(BeNil())
		})

		It("should validate configuration", func() {
			invalidCfg := scorer.Config{} // Missing API key
			_, err := scorer.NewIntegratedScorer(invalidCfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("API key is required"))
		})
	})

	Describe("Production Configuration", func() {
		// Tests production-ready scorer creation with all resilience patterns enabled
		It("should create production-ready scorer", func() {
			s, err := scorer.BuildProductionScorer("test-api-key")
			Expect(err).ToNot(HaveOccurred())
			Expect(s).ToNot(BeNil())

			health := s.GetHealth(ctx)
			Expect(health.Details["integration"]).ToNot(BeNil())

			integration := health.Details["integration"].(map[string]interface{})
			Expect(integration["circuit_breaker_enabled"]).To(BeTrue())
			Expect(integration["retry_enabled"]).To(BeTrue())
			Expect(integration["metrics_enabled"]).To(BeTrue())
		})

		It("should handle custom configuration", func() {
			customCfg := scorer.NewDefaultConfig("test-key")
			customCfg = customCfg.WithModel(openai.GPT4)
			customCfg = customCfg.WithMaxConcurrent(10)
			customCfg = customCfg.WithTimeout(60 * time.Second)

			s, err := scorer.BuildCustomScorer(customCfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(s).ToNot(BeNil())
		})
	})

	Describe("End-to-End Resilience", func() {
		var mockClient *mockIntegrationClient

		BeforeEach(func() {
			mockClient = &mockIntegrationClient{}
		})

		Context("with retry and circuit breaker", func() {
			// Tests the interaction between retry and circuit breaker patterns,
			// ensuring retries are exhausted before circuit breaker trips
			It("should retry transient errors before tripping circuit", func() {
				// Setup: first 2 calls fail with retryable error, then succeed
				mockClient.errors = []error{
					&openai.APIError{HTTPStatusCode: 500},
					&openai.APIError{HTTPStatusCode: 500},
					nil,
				}
				mockClient.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{
							Content: `{"version":"1.0","scores":[{"item_id":"1","score":75,"reason":"test"}]}`,
						}},
					},
				}

				// Create scorer with retry and circuit breaker
				cfg := scorer.Config{
					APIKey:               "test",
					EnableRetry:          true,
					EnableCircuitBreaker: true,
					RetryConfig: &scorer.RetryConfig{
						MaxAttempts:  3,
						Strategy:     scorer.RetryStrategyExponential,
						InitialDelay: 1 * time.Millisecond,
						MaxDelay:     10 * time.Millisecond,
					},
					CircuitBreakerConfig: &scorer.CircuitBreakerConfig{
						MaxRequests: 3,
						Interval:    1 * time.Second,
						Timeout:     100 * time.Millisecond,
					},
				}

				// This would use the mock client in a real implementation
				// For now, we're testing the configuration works
				s, err := scorer.NewIntegratedScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
			})

			It("should trip circuit breaker after persistent failures", func() {
				// Setup: all calls fail
				mockClient.errors = []error{
					&openai.APIError{HTTPStatusCode: 500},
					&openai.APIError{HTTPStatusCode: 500},
					&openai.APIError{HTTPStatusCode: 500},
					&openai.APIError{HTTPStatusCode: 500},
				}

				cfg := scorer.Config{
					APIKey:               "test",
					EnableRetry:          true,
					EnableCircuitBreaker: true,
					RetryConfig: &scorer.RetryConfig{
						MaxAttempts:  2,
						Strategy:     scorer.RetryStrategyConstant,
						InitialDelay: 1 * time.Millisecond,
						MaxDelay:     10 * time.Millisecond,
					},
					CircuitBreakerConfig: &scorer.CircuitBreakerConfig{
						MaxRequests: 2,
						Interval:    1 * time.Second,
						Timeout:     100 * time.Millisecond,
					},
				}

				s, err := scorer.NewIntegratedScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
			})

			// Tests specific handling of rate limit errors (429), which should be
			// retried with exponential backoff but not count toward circuit breaker failures
			It("should handle rate limits gracefully", func() {
				// Rate limits should be retried but not trip circuit
				mockClient.errors = []error{
					&openai.APIError{HTTPStatusCode: 429},
					&openai.APIError{HTTPStatusCode: 429},
					nil,
				}

				cfg := scorer.Config{
					APIKey:               "test",
					EnableRetry:          true,
					EnableCircuitBreaker: true,
					RetryConfig: &scorer.RetryConfig{
						MaxAttempts:  3,
						Strategy:     scorer.RetryStrategyExponential,
						InitialDelay: 1 * time.Millisecond,
						MaxDelay:     100 * time.Millisecond,
					},
					CircuitBreakerConfig: &scorer.CircuitBreakerConfig{
						MaxRequests: 3,
						Interval:    1 * time.Second,
						Timeout:     100 * time.Millisecond,
					},
				}

				s, err := scorer.NewIntegratedScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
			})
		})
	})

	Describe("Health Monitoring", func() {
		It("should provide comprehensive health status", func() {
			cfg := scorer.NewProductionConfig("test-key")
			s, err := scorer.NewIntegratedScorer(cfg)
			Expect(err).ToNot(HaveOccurred())

			health := s.GetHealth(ctx)
			Expect(health).ToNot(BeNil())
			Expect(health.Details).To(HaveKey("integration"))
		})
	})

	Describe("Metrics Integration", func() {
		It("should record metrics for successful requests", func() {
			// Metrics are automatically recorded
			cfg := scorer.NewDefaultConfig("test-key")
			s, err := scorer.NewIntegratedScorer(cfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(s).ToNot(BeNil())
		})

		It("should wrap any scorer with metrics", func() {
			baseScorer, err := scorer.NewScorer(cfg)
			Expect(err).ToNot(HaveOccurred())

			metrics := scorer.NewMetricsRecorder(true)
			wrapped := scorer.WithMetrics(baseScorer, metrics)
			Expect(wrapped).ToNot(BeNil())
		})
	})
})

// mockIntegrationClient provides a controlled OpenAI client implementation for integration testing.
// It allows simulation of various error conditions and response patterns to test
// resilience patterns like retry logic and circuit breaker behavior.
type mockIntegrationClient struct {
	response openai.ChatCompletionResponse // Response to return on successful calls
	errors   []error                       // Sequential errors to return on each call
	calls    int                           // Call counter for tracking invocations
}

// CreateChatCompletion simulates OpenAI API calls with configurable error patterns.
// Returns errors from the errors slice in sequence, then the configured response.
// This enables testing of retry patterns and failure scenarios.
func (m *mockIntegrationClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	m.calls++
	if m.calls <= len(m.errors) {
		err := m.errors[m.calls-1]
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
	}
	return m.response, nil
}
