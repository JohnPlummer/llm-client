package scorer_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
	"github.com/sony/gobreaker/v2"

	"github.com/JohnPlummer/post-scorer/scorer"
)

var _ = Describe("Config", func() {
	Describe("NewDefaultConfig", func() {
		It("should create config with sensible defaults", func() {
			cfg := scorer.NewDefaultConfig("test-api-key")
			
			Expect(cfg.APIKey).To(Equal("test-api-key"))
			Expect(cfg.Model).To(Equal(openai.GPT4oMini))
			Expect(cfg.MaxConcurrent).To(Equal(1))
			Expect(cfg.Timeout).To(Equal(30 * time.Second))
			Expect(cfg.EnableCircuitBreaker).To(BeFalse())
			Expect(cfg.EnableRetry).To(BeFalse())
			Expect(cfg.CircuitBreakerConfig).To(BeNil())
			Expect(cfg.RetryConfig).To(BeNil())
		})
		
		It("should panic with empty API key", func() {
			Expect(func() {
				scorer.NewDefaultConfig("")
			}).To(Panic())
		})
	})
	
	Describe("WithCircuitBreaker", func() {
		It("should enable circuit breaker with default settings", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithCircuitBreaker()
			
			Expect(cfg.EnableCircuitBreaker).To(BeTrue())
			Expect(cfg.CircuitBreakerConfig).ToNot(BeNil())
			Expect(cfg.CircuitBreakerConfig.MaxRequests).To(Equal(uint32(10)))
			Expect(cfg.CircuitBreakerConfig.Interval).To(Equal(60 * time.Second))
			Expect(cfg.CircuitBreakerConfig.Timeout).To(Equal(30 * time.Second))
			Expect(cfg.CircuitBreakerConfig.ReadyToTrip).ToNot(BeNil())
		})
		
		It("should use custom settings when provided", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			customCB := &scorer.CircuitBreakerConfig{
				MaxRequests: 5,
				Interval:    30 * time.Second,
				Timeout:     15 * time.Second,
			}
			cfg = cfg.WithCircuitBreakerConfig(customCB)
			
			Expect(cfg.EnableCircuitBreaker).To(BeTrue())
			Expect(cfg.CircuitBreakerConfig).To(Equal(customCB))
		})
		
		It("should provide ready to trip function that trips after 5 consecutive failures", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithCircuitBreaker()
			
			tripFunc := cfg.CircuitBreakerConfig.ReadyToTrip
			Expect(tripFunc).ToNot(BeNil())
			
			// Should not trip with 4 failures
			counts := gobreaker.Counts{
				Requests:             10,
				TotalFailures:        4,
				ConsecutiveFailures:  4,
			}
			Expect(tripFunc(counts)).To(BeFalse())
			
			// Should trip with 5 consecutive failures
			counts.ConsecutiveFailures = 5
			counts.TotalFailures = 5
			Expect(tripFunc(counts)).To(BeTrue())
			
			// Should trip when failure rate > 60%
			counts = gobreaker.Counts{
				Requests:             100,
				TotalFailures:        61,
				ConsecutiveFailures:  3,
			}
			Expect(tripFunc(counts)).To(BeTrue())
		})
	})
	
	Describe("WithRetry", func() {
		It("should enable retry with default exponential backoff", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithRetry()
			
			Expect(cfg.EnableRetry).To(BeTrue())
			Expect(cfg.RetryConfig).ToNot(BeNil())
			Expect(cfg.RetryConfig.MaxAttempts).To(Equal(3))
			Expect(cfg.RetryConfig.Strategy).To(Equal(scorer.RetryStrategyExponential))
			Expect(cfg.RetryConfig.InitialDelay).To(Equal(1 * time.Second))
			Expect(cfg.RetryConfig.MaxDelay).To(Equal(30 * time.Second))
		})
		
		It("should support constant backoff strategy", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithRetryStrategy(scorer.RetryStrategyConstant, 5)
			
			Expect(cfg.EnableRetry).To(BeTrue())
			Expect(cfg.RetryConfig.Strategy).To(Equal(scorer.RetryStrategyConstant))
			Expect(cfg.RetryConfig.MaxAttempts).To(Equal(5))
		})
		
		It("should support fibonacci backoff strategy", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithRetryStrategy(scorer.RetryStrategyFibonacci, 4)
			
			Expect(cfg.EnableRetry).To(BeTrue())
			Expect(cfg.RetryConfig.Strategy).To(Equal(scorer.RetryStrategyFibonacci))
			Expect(cfg.RetryConfig.MaxAttempts).To(Equal(4))
		})
		
		It("should use custom retry config when provided", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			customRetry := &scorer.RetryConfig{
				MaxAttempts:  5,
				Strategy:     scorer.RetryStrategyConstant,
				InitialDelay: 2 * time.Second,
				MaxDelay:     60 * time.Second,
			}
			cfg = cfg.WithRetryConfig(customRetry)
			
			Expect(cfg.EnableRetry).To(BeTrue())
			Expect(cfg.RetryConfig).To(Equal(customRetry))
		})
	})
	
	Describe("WithModel", func() {
		It("should set the OpenAI model", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithModel(openai.GPT4)
			
			Expect(cfg.Model).To(Equal(openai.GPT4))
		})
	})
	
	Describe("WithTimeout", func() {
		It("should set the timeout", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithTimeout(60 * time.Second)
			
			Expect(cfg.Timeout).To(Equal(60 * time.Second))
		})
		
		It("should not allow negative timeout", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			Expect(func() {
				cfg.WithTimeout(-1 * time.Second)
			}).To(Panic())
		})
	})
	
	Describe("WithMaxConcurrent", func() {
		It("should set max concurrent requests", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithMaxConcurrent(5)
			
			Expect(cfg.MaxConcurrent).To(Equal(5))
		})
		
		It("should not allow negative concurrency", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			Expect(func() {
				cfg.WithMaxConcurrent(-1)
			}).To(Panic())
		})
		
		It("should allow zero to mean sequential processing", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg = cfg.WithMaxConcurrent(0)
			
			Expect(cfg.MaxConcurrent).To(Equal(0))
		})
	})
	
	Describe("WithPromptTemplate", func() {
		It("should set custom prompt template", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			template := "Custom template: {{.Items}}"
			cfg = cfg.WithPromptTemplate(template)
			
			Expect(cfg.PromptText).To(Equal(template))
		})
		
		It("should validate Go template syntax", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			
			// Valid template
			Expect(func() {
				cfg.WithPromptTemplate("{{.Items}} {{.Count}}")
			}).ToNot(Panic())
			
			// Invalid template
			Expect(func() {
				cfg.WithPromptTemplate("{{.Items")
			}).To(Panic())
		})
	})
	
	Describe("Validate", func() {
		It("should validate a complete config", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			Expect(cfg.Validate()).To(Succeed())
		})
		
		It("should error on missing API key", func() {
			cfg := scorer.Config{}
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("API key is required"))
		})
		
		It("should error on invalid model", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg.Model = "invalid-model"
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unsupported model"))
		})
		
		It("should error on invalid retry strategy", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg.EnableRetry = true
			cfg.RetryConfig = &scorer.RetryConfig{
				Strategy: "invalid",
			}
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid retry strategy"))
		})
		
		It("should error on negative timeout", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg.Timeout = -1 * time.Second
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("timeout must be positive"))
		})
		
		It("should error on negative max concurrent", func() {
			cfg := scorer.NewDefaultConfig("test-key")
			cfg.MaxConcurrent = -1
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("MaxConcurrent must be non-negative"))
		})
	})
	
	Describe("Production Config Builder", func() {
		It("should build a production-ready config with all resilience features", func() {
			cfg := scorer.NewProductionConfig("test-key")
			
			// Should have circuit breaker enabled
			Expect(cfg.EnableCircuitBreaker).To(BeTrue())
			Expect(cfg.CircuitBreakerConfig).ToNot(BeNil())
			
			// Should have retry enabled
			Expect(cfg.EnableRetry).To(BeTrue())
			Expect(cfg.RetryConfig).ToNot(BeNil())
			
			// Should have reasonable concurrency
			Expect(cfg.MaxConcurrent).To(Equal(5))
			
			// Should have longer timeout for production
			Expect(cfg.Timeout).To(Equal(60 * time.Second))
			
			// Should use cost-effective model
			Expect(cfg.Model).To(Equal(openai.GPT4oMini))
			
			// Should validate successfully
			Expect(cfg.Validate()).To(Succeed())
		})
	})
})