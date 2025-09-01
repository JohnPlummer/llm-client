package scorer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/JohnPlummer/llm-client/scorer"
	"github.com/sashabaranov/go-openai"
)

func TestScorer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scorer Suite")
}

var _ = Describe("Types", func() {
	Describe("TextItem", func() {
		It("should create a valid text item with required fields", func() {
			item := scorer.TextItem{
				ID:      "test-1",
				Content: "This is test content",
			}
			Expect(item.ID).To(Equal("test-1"))
			Expect(item.Content).To(Equal("This is test content"))
		})

		It("should support optional metadata", func() {
			item := scorer.TextItem{
				ID:      "test-2",
				Content: "Content with metadata",
				Metadata: map[string]interface{}{
					"source":    "reddit",
					"timestamp": time.Now().Unix(),
					"score":     42,
				},
			}
			Expect(item.Metadata).To(HaveLen(3))
			Expect(item.Metadata["source"]).To(Equal("reddit"))
			Expect(item.Metadata["score"]).To(Equal(42))
		})

		It("should work with nil metadata", func() {
			item := scorer.TextItem{
				ID:      "test-3",
				Content: "Content without metadata",
			}
			Expect(item.Metadata).To(BeNil())
		})
	})

	Describe("ScoredItem", func() {
		It("should create a scored item with all fields", func() {
			originalItem := scorer.TextItem{
				ID:      "test-1",
				Content: "Original content",
			}
			scored := scorer.ScoredItem{
				Item:   originalItem,
				Score:  85,
				Reason: "High relevance to topic",
			}
			Expect(scored.Item.ID).To(Equal("test-1"))
			Expect(scored.Score).To(Equal(85))
			Expect(scored.Reason).To(Equal("High relevance to topic"))
		})

		It("should allow zero score", func() {
			scored := scorer.ScoredItem{
				Item: scorer.TextItem{
					ID:      "test-2",
					Content: "Irrelevant content",
				},
				Score:  0,
				Reason: "Not relevant",
			}
			Expect(scored.Score).To(Equal(0))
		})

		It("should allow maximum score", func() {
			scored := scorer.ScoredItem{
				Item: scorer.TextItem{
					ID:      "test-3",
					Content: "Perfect match",
				},
				Score:  100,
				Reason: "Perfect relevance",
			}
			Expect(scored.Score).To(Equal(100))
		})
	})

	Describe("Scorer Interface", func() {
		Context("ScoreTexts method", func() {
			It("should be implemented by scorer struct", func() {
				var _ scorer.Scorer = (*mockTextScorer)(nil)
			})

			It("should accept context and text items", func() {
				mock := &mockTextScorer{
					scoreFunc: func(ctx context.Context, items []scorer.TextItem, opts ...scorer.ScoringOption) ([]scorer.ScoredItem, error) {
						return []scorer.ScoredItem{
							{Item: items[0], Score: 75, Reason: "Test reason"},
						}, nil
					},
				}

				items := []scorer.TextItem{
					{ID: "1", Content: "Test content"},
				}

				results, err := mock.ScoreTexts(context.Background(), items)
				Expect(err).ToNot(HaveOccurred())
				Expect(results).To(HaveLen(1))
				Expect(results[0].Score).To(Equal(75))
			})

			It("should handle errors appropriately", func() {
				mock := &mockTextScorer{
					scoreFunc: func(ctx context.Context, items []scorer.TextItem, opts ...scorer.ScoringOption) ([]scorer.ScoredItem, error) {
						return nil, errors.New("API error")
					},
				}

				items := []scorer.TextItem{
					{ID: "1", Content: "Test content"},
				}

				results, err := mock.ScoreTexts(context.Background(), items)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("API error"))
				Expect(results).To(BeNil())
			})
		})

		Context("ScoreTextsWithOptions method", func() {
			It("should accept scoring options", func() {
				var capturedOpts []scorer.ScoringOption
				mock := &mockTextScorer{
					scoreFunc: func(ctx context.Context, items []scorer.TextItem, opts ...scorer.ScoringOption) ([]scorer.ScoredItem, error) {
						capturedOpts = opts
						return []scorer.ScoredItem{}, nil
					},
				}

				items := []scorer.TextItem{
					{ID: "1", Content: "Test"},
				}

				opt1 := scorer.WithModel(openai.GPT4)
				opt2 := scorer.WithPromptTemplate("Custom {{.Content}}")

				_, err := mock.ScoreTextsWithOptions(context.Background(), items, opt1, opt2)
				Expect(err).ToNot(HaveOccurred())
				Expect(capturedOpts).To(HaveLen(2))
			})
		})

		Context("GetHealth method", func() {
			It("should return health status", func() {
				mock := &mockTextScorer{
					healthFunc: func(ctx context.Context) scorer.HealthStatus {
						return scorer.HealthStatus{
							Healthy: true,
							Status:  "operational",
							Details: map[string]interface{}{
								"api_status":       "connected",
								"circuit_breaker":  "closed",
								"requests_pending": 0,
							},
						}
					},
				}

				health := mock.GetHealth(context.Background())
				Expect(health.Healthy).To(BeTrue())
				Expect(health.Status).To(Equal("operational"))
				Expect(health.Details).To(HaveKey("api_status"))
				Expect(health.Details["api_status"]).To(Equal("connected"))
			})

			It("should report unhealthy status", func() {
				mock := &mockTextScorer{
					healthFunc: func(ctx context.Context) scorer.HealthStatus {
						return scorer.HealthStatus{
							Healthy: false,
							Status:  "degraded",
							Details: map[string]interface{}{
								"error":           "Circuit breaker open",
								"circuit_breaker": "open",
								"last_error_time": time.Now().Unix(),
							},
						}
					},
				}

				health := mock.GetHealth(context.Background())
				Expect(health.Healthy).To(BeFalse())
				Expect(health.Status).To(Equal("degraded"))
				Expect(health.Details["error"]).To(Equal("Circuit breaker open"))
			})
		})
	})

	Describe("HealthStatus", func() {
		It("should represent healthy state", func() {
			status := scorer.HealthStatus{
				Healthy: true,
				Status:  "healthy",
				Details: map[string]interface{}{
					"uptime_seconds": 3600,
					"requests_total": 1000,
				},
			}
			Expect(status.Healthy).To(BeTrue())
			Expect(status.Status).To(Equal("healthy"))
			Expect(status.Details).To(HaveLen(2))
		})

		It("should represent unhealthy state with error details", func() {
			status := scorer.HealthStatus{
				Healthy: false,
				Status:  "error",
				Details: map[string]interface{}{
					"error_message": "Connection timeout",
					"retry_after":   60,
				},
			}
			Expect(status.Healthy).To(BeFalse())
			Expect(status.Status).To(Equal("error"))
			Expect(status.Details["error_message"]).To(Equal("Connection timeout"))
		})

		It("should work with nil details", func() {
			status := scorer.HealthStatus{
				Healthy: true,
				Status:  "ok",
			}
			Expect(status.Details).To(BeNil())
		})
	})

	Describe("ScoringOption", func() {
		It("should be a function type that modifies scoring options", func() {
			// Note: ScoringOption functions work with the internal scoringOptions type
			// We test the functions' behavior through their effects on the scorer
			modelOption := scorer.WithModel(openai.GPT4)
			Expect(modelOption).ToNot(BeNil())

			templateOption := scorer.WithPromptTemplate("Custom template")
			Expect(templateOption).ToNot(BeNil())
			
			contextOption := scorer.WithExtraContext(map[string]interface{}{"key": "value"})
			Expect(contextOption).ToNot(BeNil())
		})
	})

	Describe("Generic Config", func() {
		It("should support all necessary configuration", func() {
			cfg := scorer.Config{
				APIKey:               "test-key",
				Model:                openai.GPT4oMini,
				EnableCircuitBreaker: true,
				EnableRetry:          true,
				Timeout:              30 * time.Second,
				MaxConcurrent:        5,
				CircuitBreakerConfig: &scorer.CircuitBreakerConfig{
					MaxRequests:     10,
					Interval:        time.Minute,
					Timeout:         30 * time.Second,
					ReadyToTrip:     nil,
					OnStateChange:   nil,
				},
				RetryConfig: &scorer.RetryConfig{
					MaxAttempts: 3,
					Strategy:    scorer.RetryStrategyExponential,
					InitialDelay: time.Second,
					MaxDelay:     30 * time.Second,
				},
			}

			Expect(cfg.APIKey).To(Equal("test-key"))
			Expect(cfg.EnableCircuitBreaker).To(BeTrue())
			Expect(cfg.EnableRetry).To(BeTrue())
			Expect(cfg.Timeout).To(Equal(30 * time.Second))
			Expect(cfg.CircuitBreakerConfig).ToNot(BeNil())
			Expect(cfg.RetryConfig).ToNot(BeNil())
			Expect(cfg.RetryConfig.Strategy).To(Equal(scorer.RetryStrategyExponential))
		})
	})
})

// Mock implementation for testing
type mockTextScorer struct {
	scoreFunc  func(context.Context, []scorer.TextItem, ...scorer.ScoringOption) ([]scorer.ScoredItem, error)
	healthFunc func(context.Context) scorer.HealthStatus
}

func (m *mockTextScorer) ScoreTexts(ctx context.Context, items []scorer.TextItem, opts ...scorer.ScoringOption) ([]scorer.ScoredItem, error) {
	if m.scoreFunc != nil {
		return m.scoreFunc(ctx, items, opts...)
	}
	return nil, nil
}

func (m *mockTextScorer) ScoreTextsWithOptions(ctx context.Context, items []scorer.TextItem, opts ...scorer.ScoringOption) ([]scorer.ScoredItem, error) {
	return m.ScoreTexts(ctx, items, opts...)
}

func (m *mockTextScorer) GetHealth(ctx context.Context) scorer.HealthStatus {
	if m.healthFunc != nil {
		return m.healthFunc(ctx)
	}
	return scorer.HealthStatus{Healthy: true, Status: "ok"}
}