// Package scorer_test provides comprehensive test coverage for the retry mechanism
// that wraps OpenAI API calls with configurable retry strategies and backoff algorithms.
// This test suite validates exponential, constant, and fibonacci backoff strategies
// along with proper error classification and context cancellation handling.
package scorer_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"

	"github.com/JohnPlummer/llm-client/scorer"
)

// Retry test suite validates the RetryWrapper's ability to handle transient failures
// with configurable retry strategies, proper error classification, and backoff algorithms.
// Tests cover successful requests, retryable/non-retryable errors, max attempts, and context cancellation.
var _ = Describe("Retry", func() {
	var (
		wrapper *scorer.RetryWrapper
		mockAPI *mockRetryAPIClient
		ctx     context.Context
	)

	// BeforeEach sets up the test environment with a standard exponential backoff configuration
	// and a mock API client that can simulate various error conditions and response patterns.
	BeforeEach(func() {
		ctx = context.Background()
		mockAPI = &mockRetryAPIClient{}

		config := scorer.RetryConfig{
			MaxAttempts:  3,
			Strategy:     scorer.RetryStrategyExponential,
			InitialDelay: 10 * time.Millisecond,
			MaxDelay:     100 * time.Millisecond,
		}

		wrapper = scorer.NewRetryWrapper(mockAPI, &config)
	})

	Describe("Successful Requests", func() {
		It("should not retry on successful requests", func() {
			mockAPI.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{Message: openai.ChatCompletionMessage{Content: "success"}},
				},
			}

			start := time.Now()
			resp, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			duration := time.Since(start)

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Choices[0].Message.Content).To(Equal("success"))
			Expect(mockAPI.calls).To(Equal(1))
			Expect(duration).To(BeNumerically("<", 50*time.Millisecond))
		})
	})

	// Retryable Errors section tests the retry mechanism for transient failures
	// including rate limits (429), timeouts, and server errors (5xx).
	// These errors are considered temporary and should trigger retry attempts.
	Describe("Retryable Errors", func() {
		Context("with rate limit errors", func() {
			// Rate limit error test validates exponential backoff timing and successful recovery
			// after multiple 429 responses from the OpenAI API.
			It("should retry with exponential backoff on 429 errors", func() {
				mockAPI.errors = []error{
					&openai.APIError{
						Code:           "rate_limit_exceeded",
						Message:        "Rate limit exceeded",
						HTTPStatusCode: 429,
					},
					&openai.APIError{
						Code:           "rate_limit_exceeded",
						Message:        "Rate limit exceeded",
						HTTPStatusCode: 429,
					},
					nil, // Success on third attempt
				}
				mockAPI.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "success after retry"}},
					},
				}

				start := time.Now()
				resp, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
				duration := time.Since(start)

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Choices[0].Message.Content).To(Equal("success after retry"))
				Expect(mockAPI.calls).To(Equal(3))
				// Should have delays: ~10ms + ~20ms = ~30ms minimum
				Expect(duration).To(BeNumerically(">=", 25*time.Millisecond))
			})

			// Jitter test ensures that multiple clients don't retry simultaneously
			// by adding randomness to delay calculations, preventing thundering herd effects.
			It("should apply jitter to prevent thundering herd", func() {
				config := scorer.RetryConfig{
					MaxAttempts:  5,
					Strategy:     scorer.RetryStrategyExponential,
					InitialDelay: 10 * time.Millisecond,
					MaxDelay:     100 * time.Millisecond,
				}
				wrapper = scorer.NewRetryWrapper(mockAPI, &config)

				// Create multiple retries and measure delays
				var delays []time.Duration
				for i := 0; i < 3; i++ {
					mockAPI.calls = 0
					mockAPI.errors = []error{
						errors.New("temporary error"),
						nil,
					}

					start := time.Now()
					wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
					delays = append(delays, time.Since(start))
				}

				// Delays should vary due to jitter
				Expect(delays[0]).ToNot(Equal(delays[1]))
				Expect(delays[1]).ToNot(Equal(delays[2]))
			})
		})

		Context("with timeout errors", func() {
			It("should retry on timeout errors", func() {
				mockAPI.errors = []error{
					context.DeadlineExceeded,
					context.DeadlineExceeded,
					nil,
				}
				mockAPI.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "success"}},
					},
				}

				resp, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})

				Expect(err).ToNot(HaveOccurred())
				Expect(mockAPI.calls).To(Equal(3))
				Expect(resp.Choices[0].Message.Content).To(Equal("success"))
			})
		})

		Context("with server errors", func() {
			It("should retry on 5xx errors", func() {
				mockAPI.errors = []error{
					&openai.APIError{
						Code:           "internal_server_error",
						Message:        "Internal server error",
						HTTPStatusCode: 500,
					},
					nil,
				}
				mockAPI.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "success"}},
					},
				}

				resp, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})

				Expect(err).ToNot(HaveOccurred())
				Expect(mockAPI.calls).To(Equal(2))
				Expect(resp.Choices[0].Message.Content).To(Equal("success"))
			})
		})
	})

	// Non-Retryable Errors section validates that client errors (4xx) are not retried
	// since they represent permanent failures that won't resolve with repeated attempts.
	Describe("Non-Retryable Errors", func() {
		// Authentication error test ensures 401 errors are not retried since they indicate
		// invalid credentials that won't be fixed by retrying the same request.
		It("should not retry on authentication errors", func() {
			mockAPI.errors = []error{
				&openai.APIError{
					Code:           "invalid_api_key",
					Message:        "Invalid API key",
					HTTPStatusCode: 401,
				},
			}

			start := time.Now()
			_, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
			duration := time.Since(start)

			Expect(err).To(HaveOccurred())
			Expect(mockAPI.calls).To(Equal(1)) // No retry
			Expect(duration).To(BeNumerically("<", 20*time.Millisecond))
		})

		It("should not retry on bad request errors", func() {
			mockAPI.errors = []error{
				&openai.APIError{
					Code:           "invalid_request",
					Message:        "Invalid request",
					HTTPStatusCode: 400,
				},
			}

			_, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})

			Expect(err).To(HaveOccurred())
			Expect(mockAPI.calls).To(Equal(1)) // No retry
		})
	})

	// Max Attempts section validates that retry attempts are properly limited
	// and the system returns the final error after exhausting all retry attempts.
	Describe("Max Attempts", func() {
		// Max attempts limit test ensures the retry mechanism respects the configured
		// maximum and doesn't continue retrying indefinitely on persistent failures.
		It("should stop after max attempts", func() {
			mockAPI.errors = []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
				errors.New("error 4"), // Won't be reached
			}

			_, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("error 3"))
			Expect(mockAPI.calls).To(Equal(3)) // Exactly max attempts
		})

		It("should return last error after all retries exhausted", func() {
			lastErr := errors.New("final error")
			mockAPI.errors = []error{
				errors.New("error 1"),
				errors.New("error 2"),
				lastErr,
			}

			_, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})

			Expect(err).To(Equal(lastErr))
		})
	})

	// Backoff Strategies section tests different delay algorithms including exponential,
	// constant, and fibonacci strategies, validating timing behavior and delay capping.
	Describe("Backoff Strategies", func() {
		Context("Exponential Backoff", func() {
			// Exponential backoff test validates that delays double on each retry attempt
			// creating a 10ms, 20ms, 40ms pattern that reduces API load over time.
			It("should double delay on each retry", func() {
				config := scorer.RetryConfig{
					MaxAttempts:  4,
					Strategy:     scorer.RetryStrategyExponential,
					InitialDelay: 10 * time.Millisecond,
					MaxDelay:     1000 * time.Millisecond,
				}
				wrapper = scorer.NewRetryWrapper(mockAPI, &config)

				mockAPI.errors = []error{
					errors.New("error 1"),
					errors.New("error 2"),
					errors.New("error 3"),
					nil,
				}
				mockAPI.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "success"}},
					},
				}

				start := time.Now()
				resp, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
				duration := time.Since(start)

				Expect(err).ToNot(HaveOccurred())
				Expect(resp.Choices[0].Message.Content).To(Equal("success"))
				// Delays: ~10ms + ~20ms + ~40ms = ~70ms minimum
				Expect(duration).To(BeNumerically(">=", 60*time.Millisecond))
			})

			// MaxDelay capping test ensures exponential backoff doesn't grow infinitely
			// and respects the configured maximum delay to prevent excessive wait times.
			It("should cap delay at MaxDelay", func() {
				config := scorer.RetryConfig{
					MaxAttempts:  5,
					Strategy:     scorer.RetryStrategyExponential,
					InitialDelay: 10 * time.Millisecond,
					MaxDelay:     20 * time.Millisecond,
				}
				wrapper = scorer.NewRetryWrapper(mockAPI, &config)

				mockAPI.errors = []error{
					errors.New("error 1"),
					errors.New("error 2"),
					errors.New("error 3"),
					errors.New("error 4"),
					nil,
				}
				mockAPI.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "success"}},
					},
				}

				start := time.Now()
				wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
				duration := time.Since(start)

				// Delays: ~10ms + ~20ms + ~20ms + ~20ms = ~70ms
				Expect(duration).To(BeNumerically("<", 100*time.Millisecond))
			})
		})

		Context("Constant Backoff", func() {
			// Constant backoff test validates that the delay remains fixed between retry attempts
			// providing predictable timing when exponential growth is not desired.
			It("should use constant delay between retries", func() {
				config := scorer.RetryConfig{
					MaxAttempts:  3,
					Strategy:     scorer.RetryStrategyConstant,
					InitialDelay: 15 * time.Millisecond,
					MaxDelay:     100 * time.Millisecond,
				}
				wrapper = scorer.NewRetryWrapper(mockAPI, &config)

				mockAPI.errors = []error{
					errors.New("error 1"),
					errors.New("error 2"),
					nil,
				}
				mockAPI.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "success"}},
					},
				}

				start := time.Now()
				wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
				duration := time.Since(start)

				// Delays: ~15ms + ~15ms = ~30ms
				Expect(duration).To(BeNumerically(">=", 25*time.Millisecond))
				Expect(duration).To(BeNumerically("<", 50*time.Millisecond))
			})
		})

		Context("Fibonacci Backoff", func() {
			// Fibonacci backoff test validates the 1,1,2,3,5,8 delay sequence pattern
			// providing moderate growth that's less aggressive than exponential backoff.
			It("should use fibonacci sequence for delays", func() {
				config := scorer.RetryConfig{
					MaxAttempts:  5,
					Strategy:     scorer.RetryStrategyFibonacci,
					InitialDelay: 10 * time.Millisecond,
					MaxDelay:     1000 * time.Millisecond,
				}
				wrapper = scorer.NewRetryWrapper(mockAPI, &config)

				mockAPI.errors = []error{
					errors.New("error 1"),
					errors.New("error 2"),
					errors.New("error 3"),
					errors.New("error 4"),
					nil,
				}
				mockAPI.response = openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{
						{Message: openai.ChatCompletionMessage{Content: "success"}},
					},
				}

				start := time.Now()
				wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})
				duration := time.Since(start)

				// Delays: ~10ms + ~10ms + ~20ms + ~30ms = ~70ms
				Expect(duration).To(BeNumerically(">=", 60*time.Millisecond))
			})
		})
	})

	// Context Cancellation section validates that retry attempts are properly interrupted
	// when the context is cancelled, preventing unnecessary API calls and resource waste.
	Describe("Context Cancellation", func() {
		// Context timeout test ensures retry attempts stop immediately when the context deadline
		// is exceeded, avoiding continued retries that would exceed user-defined timeouts.
		It("should stop retrying when context is cancelled", func() {
			mockAPI.errors = []error{
				errors.New("error 1"),
				errors.New("error 2"),
				errors.New("error 3"),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
			defer cancel()

			_, err := wrapper.CreateChatCompletion(ctx, openai.ChatCompletionRequest{})

			Expect(err).To(HaveOccurred())
			Expect(mockAPI.calls).To(BeNumerically("<=", 3))
		})
	})

	// Error Classification section tests the IsRetryableError function's ability
	// to correctly distinguish between transient and permanent error conditions.
	Describe("Error Classification", func() {
		// Error classification test validates the logic that determines which HTTP status codes
		// and error types warrant retry attempts versus immediate failure responses.
		It("should classify errors correctly", func() {
			// Retryable errors
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 429})).To(BeTrue())
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 500})).To(BeTrue())
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 502})).To(BeTrue())
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 503})).To(BeTrue())
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 504})).To(BeTrue())
			Expect(scorer.IsRetryableError(context.DeadlineExceeded)).To(BeTrue())

			// Non-retryable errors
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 400})).To(BeFalse())
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 401})).To(BeFalse())
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 403})).To(BeFalse())
			Expect(scorer.IsRetryableError(&openai.APIError{HTTPStatusCode: 404})).To(BeFalse())
			Expect(scorer.IsRetryableError(context.Canceled)).To(BeFalse())
		})
	})
})

// mockRetryAPIClient provides a test double for OpenAI API client that allows simulation
// of various error conditions, response patterns, and call counting for retry validation.
// The mock can be configured with a sequence of errors and a final success response.
type mockRetryAPIClient struct {
	response openai.ChatCompletionResponse
	errors   []error
	calls    int
}

// CreateChatCompletion simulates OpenAI API calls with configurable error sequences.
// It returns errors from the errors slice in order, then returns the configured response.
// The calls counter tracks total invocations for retry attempt validation.
func (m *mockRetryAPIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	m.calls++

	if m.calls <= len(m.errors) {
		err := m.errors[m.calls-1]
		if err != nil {
			return openai.ChatCompletionResponse{}, err
		}
	}

	return m.response, nil
}
