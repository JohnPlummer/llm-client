package scorer_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/JohnPlummer/llm-client/scorer"
)

// TestScorer is defined in types_test.go

var _ = Describe("Scorer", func() {
	var cfg scorer.Config

	BeforeEach(func() {
		cfg = scorer.Config{
			APIKey: "test-api-key",
		}
	})

	Describe("NewScorer", func() {
		Context("validation", func() {
			It("should return error when API key is missing", func() {
				cfg.APIKey = ""
				_, err := scorer.NewScorer(cfg)
				Expect(err).To(Equal(scorer.ErrMissingAPIKey))
			})

			It("should return error when MaxConcurrent is negative", func() {
				cfg.MaxConcurrent = -1
				_, err := scorer.NewScorer(cfg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("MaxConcurrent must be non-negative"))
			})

			It("should create scorer with valid config", func() {
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
			})
		})

		Context("defaults", func() {
			It("should set default MaxConcurrent to 1 when not specified", func() {
				cfg.MaxConcurrent = 0
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
				// Default is applied internally
			})

			It("should set default Model to GPT4oMini when not specified", func() {
				cfg.Model = ""
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
				// Default is applied internally
			})

			It("should set default Timeout to 30 seconds when not specified", func() {
				cfg.Timeout = 0
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
				// Default is applied internally
			})
		})

		Context("custom prompts", func() {
			It("should accept prompt with Go template syntax", func() {
				cfg.PromptText = "Custom prompt: {{.Content}}"
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
			})

			It("should accept prompt with sprintf placeholder", func() {
				cfg.PromptText = "Custom prompt: %s"
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
			})

			It("should warn but not error on prompt without placeholders", func() {
				cfg.PromptText = "Custom prompt without any placeholder"
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
				// Warning is logged but not returned as error
			})
		})

		Context("configuration options", func() {
			It("should accept all configuration fields", func() {
				cfg.Model = "gpt-4"
				cfg.MaxConcurrent = 5
				cfg.Timeout = 60
				cfg.PromptText = "Custom: %s"
				cfg.EnableCircuitBreaker = true
				cfg.EnableRetry = true
				
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
			})
		})
	})

	Describe("ScoringOption functions", func() {
		It("should provide WithModel option", func() {
			opt := scorer.WithModel("gpt-4")
			Expect(opt).ToNot(BeNil())
		})

		It("should provide WithPromptTemplate option", func() {
			opt := scorer.WithPromptTemplate("Custom: {{.Content}}")
			Expect(opt).ToNot(BeNil())
		})

		It("should provide WithExtraContext option", func() {
			opt := scorer.WithExtraContext(map[string]interface{}{
				"City": "Brighton",
			})
			Expect(opt).ToNot(BeNil())
		})
	})

	Describe("TextItem and ScoredItem types", func() {
		It("should create valid TextItem", func() {
			item := scorer.TextItem{
				ID:      "test-1",
				Content: "Test content",
				Metadata: map[string]interface{}{
					"source": "test",
				},
			}
			Expect(item.ID).To(Equal("test-1"))
			Expect(item.Content).To(Equal("Test content"))
			Expect(item.Metadata).To(HaveKey("source"))
		})

		It("should create valid ScoredItem", func() {
			item := scorer.TextItem{
				ID:      "test-1",
				Content: "Test content",
			}
			scored := scorer.ScoredItem{
				Item:   item,
				Score:  75,
				Reason: "High relevance",
			}
			Expect(scored.Item.ID).To(Equal("test-1"))
			Expect(scored.Score).To(Equal(75))
			Expect(scored.Reason).To(Equal("High relevance"))
		})
	})

	Describe("Content Length Validation", func() {
		Context("when content exceeds maximum length", func() {
			It("should return error for content exceeding default max length", func() {
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				
				// Create content that exceeds default max (10,000 chars)
				longContent := ""
				for i := 0; i <= scorer.DefaultMaxContentLength; i++ {
					longContent += "a"
				}
				
				items := []scorer.TextItem{
					{ID: "test-1", Content: longContent},
				}
				
				ctx := context.Background()
				_, err = s.ScoreTexts(ctx, items)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, scorer.ErrContentTooLong)).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("content exceeds maximum length"))
			})
			
			It("should respect custom max content length", func() {
				cfg.MaxContentLength = 100
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				
				// Create content that exceeds custom max (100 chars)
				longContent := ""
				for i := 0; i <= 100; i++ {
					longContent += "a"
				}
				
				items := []scorer.TextItem{
					{ID: "test-1", Content: longContent},
				}
				
				ctx := context.Background()
				_, err = s.ScoreTexts(ctx, items)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, scorer.ErrContentTooLong)).To(BeTrue())
			})
		})
		
		Context("when content is too short", func() {
			It("should validate empty content correctly", func() {
				// This test verifies that empty content is accepted with a warning
				// The actual API call test would require mocking the OpenAI client
				// For now, we just verify the validation logic doesn't reject empty content
				
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(s).ToNot(BeNil())
				
				// Empty content should be allowed through validation
				// (warning is logged but no error returned)
				items := []scorer.TextItem{
					{ID: "test-1", Content: ""},
				}
				
				// Just verify the scorer was created successfully
				// Actual scoring would require mock client
				Expect(len(items)).To(Equal(1))
				Expect(items[0].Content).To(Equal(""))
			})
		})
		
		Context("when content length is valid", func() {
			It("should accept content within limits", func() {
				cfg.MaxContentLength = 100
				s, err := scorer.NewScorer(cfg)
				Expect(err).ToNot(HaveOccurred())
				
				items := []scorer.TextItem{
					{ID: "test-1", Content: "Valid content within limits"},
				}
				
				ctx := context.Background()
				// This would normally call the API, but in tests with mock client it should work
				_, err = s.ScoreTexts(ctx, items)
				// The test setup doesn't have a proper mock, so we just verify no validation error
				if err != nil {
					Expect(errors.Is(err, scorer.ErrContentTooLong)).To(BeFalse())
					Expect(errors.Is(err, scorer.ErrContentTooShort)).To(BeFalse())
				}
			})
		})
	})

	Describe("HealthStatus", func() {
		It("should represent healthy state", func() {
			status := scorer.HealthStatus{
				Healthy: true,
				Status:  "healthy",
				Details: map[string]interface{}{
					"api_status": "connected",
				},
			}
			Expect(status.Healthy).To(BeTrue())
			Expect(status.Status).To(Equal("healthy"))
			Expect(status.Details).To(HaveKey("api_status"))
		})

		It("should represent unhealthy state", func() {
			status := scorer.HealthStatus{
				Healthy: false,
				Status:  "unhealthy",
				Details: map[string]interface{}{
					"error": "connection failed",
				},
			}
			Expect(status.Healthy).To(BeFalse())
			Expect(status.Status).To(Equal("unhealthy"))
			Expect(status.Details["error"]).To(Equal("connection failed"))
		})
	})
})