package scorer_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
	"github.com/sony/gobreaker/v2"

	"github.com/JohnPlummer/post-scorer/scorer"
)

var _ = Describe("CircuitBreaker", func() {
	var (
		cb      *scorer.CircuitBreakerWrapper
		mockAPI *mockAPIClient
		ctx     context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockAPI = &mockAPIClient{}
		
		config := scorer.CircuitBreakerConfig{
			MaxRequests: 3,
			Interval:    10 * time.Second,
			Timeout:     5 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				// Trip after 3 consecutive failures
				return counts.ConsecutiveFailures >= 3
			},
		}
		
		cb = scorer.NewCircuitBreakerWrapper(mockAPI, &config)
	})

	Describe("Normal Operation", func() {
		It("should pass through successful requests", func() {
			mockAPI.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{Message: openai.ChatCompletionMessage{Content: "test response"}},
				},
			}

			resp, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Choices[0].Message.Content).To(Equal("test response"))
		})

		It("should count successful requests", func() {
			mockAPI.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{Message: openai.ChatCompletionMessage{Content: "test"}},
				},
			}

			for i := 0; i < 5; i++ {
				_, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
				Expect(err).ToNot(HaveOccurred())
			}

			state := cb.State()
			Expect(state).To(Equal(gobreaker.StateClosed))
		})
	})

	Describe("Error Handling", func() {
		Context("with retryable errors", func() {
			It("should not trip on rate limit errors (429)", func() {
				mockAPI.err = &openai.APIError{
					Code:           "rate_limit_exceeded",
					Message:        "Rate limit exceeded",
					HTTPStatusCode: 429,
				}

				// Send multiple 429 errors
				for i := 0; i < 5; i++ {
					_, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
					Expect(err).To(HaveOccurred())
				}

				// Circuit should still be closed
				state := cb.State()
				Expect(state).To(Equal(gobreaker.StateClosed))
			})

			It("should not trip on timeout errors", func() {
				mockAPI.err = context.DeadlineExceeded

				// Send multiple timeout errors
				for i := 0; i < 5; i++ {
					_, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
					Expect(err).To(HaveOccurred())
				}

				// Circuit should still be closed
				state := cb.State()
				Expect(state).To(Equal(gobreaker.StateClosed))
			})
		})

		Context("with circuit-breaking errors", func() {
			It("should trip on server errors (5xx)", func() {
				mockAPI.err = &openai.APIError{
					Code:           "internal_server_error",
					Message:        "Internal server error",
					HTTPStatusCode: 500,
				}

				// Send 3 server errors to trip the circuit
				for i := 0; i < 3; i++ {
					_, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
					Expect(err).To(HaveOccurred())
				}

				// Circuit should be open
				state := cb.State()
				Expect(state).To(Equal(gobreaker.StateOpen))

				// Next request should fail immediately
				_, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, gobreaker.ErrOpenState)).To(BeTrue())
			})

			It("should trip on authentication errors", func() {
				mockAPI.err = &openai.APIError{
					Code:           "invalid_api_key",
					Message:        "Invalid API key",
					HTTPStatusCode: 401,
				}

				// Send 3 auth errors to trip the circuit
				for i := 0; i < 3; i++ {
					_, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
					Expect(err).To(HaveOccurred())
				}

				// Circuit should be open
				state := cb.State()
				Expect(state).To(Equal(gobreaker.StateOpen))
			})

			It("should trip on unknown errors", func() {
				mockAPI.err = errors.New("unknown error")

				// Send 3 unknown errors to trip the circuit
				for i := 0; i < 3; i++ {
					_, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
					Expect(err).To(HaveOccurred())
				}

				// Circuit should be open
				state := cb.State()
				Expect(state).To(Equal(gobreaker.StateOpen))
			})
		})
	})

	Describe("Recovery", func() {
		It("should transition to half-open after timeout", func() {
			// Trip the circuit
			mockAPI.err = errors.New("server error")
			for i := 0; i < 3; i++ {
				cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			}

			// Verify circuit is open
			Expect(cb.State()).To(Equal(gobreaker.StateOpen))

			// Wait for timeout (using shorter timeout for test)
			config := scorer.CircuitBreakerConfig{
				MaxRequests: 1,
				Interval:    1 * time.Second,
				Timeout:     100 * time.Millisecond,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.ConsecutiveFailures >= 3
				},
			}
			cb = scorer.NewCircuitBreakerWrapper(mockAPI, &config)
			
			// Trip it
			for i := 0; i < 3; i++ {
				cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			}

			// Wait for timeout
			time.Sleep(150 * time.Millisecond)

			// Should be half-open, ready to try again
			mockAPI.err = nil
			mockAPI.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{Message: openai.ChatCompletionMessage{Content: "recovered"}},
				},
			}

			// This request should succeed and close the circuit
			resp, err := cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Choices[0].Message.Content).To(Equal("recovered"))
			
			// Circuit should be closed again
			Eventually(func() gobreaker.State {
				return cb.State()
			}).Should(Equal(gobreaker.StateClosed))
		})
	})

	Describe("Health Check", func() {
		It("should report healthy when circuit is closed", func() {
			health := cb.GetHealth()
			Expect(health.Healthy).To(BeTrue())
			Expect(health.Status).To(Equal("closed"))
			Expect(health.Details["state"]).To(Equal("closed"))
		})

		It("should report unhealthy when circuit is open", func() {
			// Trip the circuit
			mockAPI.err = errors.New("server error")
			for i := 0; i < 3; i++ {
				cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			}

			health := cb.GetHealth()
			Expect(health.Healthy).To(BeFalse())
			Expect(health.Status).To(Equal("open"))
			Expect(health.Details["state"]).To(Equal("open"))
		})

		It("should report degraded when circuit is half-open", func() {
			// Create CB with very short timeout for testing
			config := scorer.CircuitBreakerConfig{
				MaxRequests: 1,
				Interval:    1 * time.Second,
				Timeout:     50 * time.Millisecond,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.ConsecutiveFailures >= 3
				},
			}
			cb = scorer.NewCircuitBreakerWrapper(mockAPI, &config)

			// Trip the circuit
			mockAPI.err = errors.New("server error")
			for i := 0; i < 3; i++ {
				cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			}

			// Wait for half-open
			time.Sleep(60 * time.Millisecond)

			health := cb.GetHealth()
			Expect(health.Status).To(Equal("half-open"))
			Expect(health.Details["state"]).To(Equal("half-open"))
		})
	})

	Describe("State Change Callbacks", func() {
		It("should call state change callback", func() {
			var stateChanges []string
			
			config := scorer.CircuitBreakerConfig{
				MaxRequests: 3,
				Interval:    10 * time.Second,
				Timeout:     5 * time.Second,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.ConsecutiveFailures >= 3
				},
				OnStateChange: func(name string, from, to gobreaker.State) {
					stateChanges = append(stateChanges, 
						from.String()+"->"+to.String())
				},
			}
			
			cb = scorer.NewCircuitBreakerWrapper(mockAPI, &config)

			// Trip the circuit
			mockAPI.err = errors.New("server error")
			for i := 0; i < 3; i++ {
				cb.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			}

			Expect(stateChanges).To(ContainElement("closed->open"))
		})
	})

	Describe("Smart Error Classification", func() {
		It("should classify errors correctly", func() {
			// Rate limit - should not trip
			Expect(scorer.ShouldTripCircuit(&openai.APIError{HTTPStatusCode: 429})).To(BeFalse())
			
			// Server error - should trip
			Expect(scorer.ShouldTripCircuit(&openai.APIError{HTTPStatusCode: 500})).To(BeTrue())
			
			// Auth error - should trip
			Expect(scorer.ShouldTripCircuit(&openai.APIError{HTTPStatusCode: 401})).To(BeTrue())
			
			// Timeout - should not trip
			Expect(scorer.ShouldTripCircuit(context.DeadlineExceeded)).To(BeFalse())
			
			// Unknown error - should trip
			Expect(scorer.ShouldTripCircuit(errors.New("unknown"))).To(BeTrue())
		})
	})
})

// Mock API client for testing
type mockAPIClient struct {
	response openai.ChatCompletionResponse
	err      error
	calls    int
}

func (m *mockAPIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	m.calls++
	return m.response, m.err
}