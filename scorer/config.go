package scorer

import (
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/sony/gobreaker/v2"
)

// NewDefaultConfig creates a config with sensible defaults
func NewDefaultConfig(apiKey string) Config {
	if apiKey == "" {
		panic("API key is required")
	}
	
	return Config{
		APIKey:        apiKey,
		Model:         openai.GPT4oMini,
		MaxConcurrent: 1,
		Timeout:       30 * time.Second,
	}
}

// NewProductionConfig creates a production-ready config with all resilience features
func NewProductionConfig(apiKey string) Config {
	cfg := NewDefaultConfig(apiKey)
	cfg.MaxConcurrent = 5
	cfg.Timeout = 60 * time.Second
	
	// Enable circuit breaker with production settings
	cfg = cfg.WithCircuitBreaker()
	
	// Enable retry with production settings
	cfg = cfg.WithRetry()
	
	return cfg
}

// WithCircuitBreaker enables circuit breaker with default settings
func (c Config) WithCircuitBreaker() Config {
	c.EnableCircuitBreaker = true
	c.CircuitBreakerConfig = &CircuitBreakerConfig{
		MaxRequests: 10,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if 5 consecutive failures OR failure rate > 60%
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.ConsecutiveFailures >= 5 || 
				   (counts.Requests >= 10 && failureRatio > 0.6)
		},
	}
	return c
}

// WithCircuitBreakerConfig enables circuit breaker with custom settings
func (c Config) WithCircuitBreakerConfig(config *CircuitBreakerConfig) Config {
	c.EnableCircuitBreaker = true
	c.CircuitBreakerConfig = config
	return c
}

// WithRetry enables retry with default exponential backoff
func (c Config) WithRetry() Config {
	c.EnableRetry = true
	c.RetryConfig = &RetryConfig{
		MaxAttempts:  3,
		Strategy:     RetryStrategyExponential,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
	}
	return c
}

// WithRetryStrategy enables retry with specified strategy
func (c Config) WithRetryStrategy(strategy RetryStrategy, maxAttempts int) Config {
	c.EnableRetry = true
	c.RetryConfig = &RetryConfig{
		MaxAttempts:  maxAttempts,
		Strategy:     strategy,
		InitialDelay: 1 * time.Second,
		MaxDelay:     30 * time.Second,
	}
	return c
}

// WithRetryConfig enables retry with custom settings
func (c Config) WithRetryConfig(config *RetryConfig) Config {
	c.EnableRetry = true
	c.RetryConfig = config
	return c
}

// WithModel sets the OpenAI model
func (c Config) WithModel(model string) Config {
	c.Model = model
	return c
}

// WithTimeout sets the request timeout
func (c Config) WithTimeout(timeout time.Duration) Config {
	if timeout < 0 {
		panic("timeout must be positive")
	}
	c.Timeout = timeout
	return c
}

// WithMaxConcurrent sets the maximum concurrent requests
func (c Config) WithMaxConcurrent(max int) Config {
	if max < 0 {
		panic("MaxConcurrent must be non-negative")
	}
	c.MaxConcurrent = max
	return c
}

// WithPromptTemplate sets a custom prompt template
func (c Config) WithPromptTemplate(templateText string) Config {
	// Validate template syntax
	_, err := template.New("prompt").Parse(templateText)
	if err != nil {
		panic(fmt.Sprintf("invalid template syntax: %v", err))
	}
	c.PromptText = templateText
	return c
}

// Validate checks if the config is valid
func (c Config) Validate() error {
	// Required fields
	if c.APIKey == "" {
		return errors.New("API key is required")
	}
	
	// Model validation
	if c.Model != "" && !isValidModel(c.Model) {
		return fmt.Errorf("unsupported model: %s", c.Model)
	}
	
	// Timeout validation
	if c.Timeout < 0 {
		return errors.New("timeout must be positive")
	}
	
	// Concurrency validation
	if c.MaxConcurrent < 0 {
		return errors.New("MaxConcurrent must be non-negative")
	}
	
	// Circuit breaker validation
	if c.EnableCircuitBreaker && c.CircuitBreakerConfig == nil {
		return errors.New("circuit breaker enabled but config is nil")
	}
	
	// Retry validation
	if c.EnableRetry {
		if c.RetryConfig == nil {
			return errors.New("retry enabled but config is nil")
		}
		
		if !isValidRetryStrategy(c.RetryConfig.Strategy) {
			return fmt.Errorf("invalid retry strategy: %s", c.RetryConfig.Strategy)
		}
		
		if c.RetryConfig.MaxAttempts <= 0 {
			return errors.New("retry MaxAttempts must be positive")
		}
		
		if c.RetryConfig.InitialDelay <= 0 {
			return errors.New("retry InitialDelay must be positive")
		}
		
		if c.RetryConfig.MaxDelay <= 0 {
			return errors.New("retry MaxDelay must be positive")
		}
	}
	
	// Template validation
	if c.PromptText != "" {
		if strings.Contains(c.PromptText, "{{") && strings.Contains(c.PromptText, "}}") {
			_, err := template.New("prompt").Parse(c.PromptText)
			if err != nil {
				return fmt.Errorf("invalid prompt template: %w", err)
			}
		}
	}
	
	return nil
}

// isValidModel checks if the model is supported
func isValidModel(model string) bool {
	validModels := []string{
		openai.GPT4,
		openai.GPT4o,
		openai.GPT4oMini,
		openai.GPT4Turbo,
		openai.GPT432K,
		openai.GPT3Dot5Turbo,
		openai.GPT3Dot5Turbo16K,
	}
	
	for _, valid := range validModels {
		if model == valid {
			return true
		}
	}
	return false
}

// isValidRetryStrategy checks if the retry strategy is valid
func isValidRetryStrategy(strategy RetryStrategy) bool {
	switch strategy {
	case RetryStrategyExponential, RetryStrategyConstant, RetryStrategyFibonacci:
		return true
	default:
		return false
	}
}