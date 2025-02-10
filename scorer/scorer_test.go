package scorer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/JohnPlummer/reddit-client/reddit"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sashabaranov/go-openai"
)

func TestScorer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scorer Suite")
}

// mockOpenAIClient implements the OpenAIClient interface for testing
type mockOpenAIClient struct {
	response                 openai.ChatCompletionResponse
	err                      error
	createChatCompletionFunc func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

func (m *mockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	if m.createChatCompletionFunc != nil {
		return m.createChatCompletionFunc(ctx, req)
	}
	return m.response, m.err
}

var _ = Describe("Scorer", func() {
	var (
		mockClient *mockOpenAIClient
		s          Scorer
		posts      []reddit.Post
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockClient = &mockOpenAIClient{}
		s = newWithClient(mockClient)
		posts = []reddit.Post{
			{
				ID:       "post1",
				Title:    "Test Post 1",
				SelfText: "Test Content 1",
			},
		}
	})

	Context("New", func() {
		It("should return error when API key is missing", func() {
			_, err := New(Config{})
			Expect(err).To(Equal(ErrMissingAPIKey))
		})

		It("should create a working scorer with valid API key", func() {
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version": "1.0","scores": [{"post_id": "post1", "title": "Test Post", "score": 85, "reason": "Test reason"}]}`,
							},
						}},
					}, nil
				},
			}

			s := &scorer{
				client: mockClient,
				config: Config{OpenAIKey: "test-key"},
				prompt: batchScorePrompt,
			}

			scored, err := s.ScorePosts(ctx, posts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(1))

			result := scored[0]
			Expect(result.Post).To(Equal(posts[0]))
			Expect(result.Score).To(BeNumerically(">=", 0))
			Expect(result.Score).To(BeNumerically("<=", 100))
			Expect(result.Reason).NotTo(BeEmpty())
		})
	})

	Context("ScorePosts", func() {
		BeforeEach(func() {
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version": "1.0","scores": [{"post_id": "post1", "title": "Test Post", "score": 85, "reason": "Test reason"}]}`,
							},
						}},
					}, nil
				},
			}
			s = &scorer{
				client: mockClient,
				config: Config{},
				prompt: batchScorePrompt,
			}
		})

		It("should return nil for empty posts", func() {
			scored, err := s.ScorePosts(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(BeNil())
		})

		It("should handle API errors", func() {
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{}, errors.New("API error")
				},
			}
			s = &scorer{client: mockClient, config: Config{}, prompt: batchScorePrompt}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("API error"))
		})

		It("should handle empty API response", func() {
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{},
					}, nil
				},
			}
			s = &scorer{client: mockClient, config: Config{}, prompt: batchScorePrompt}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no response from OpenAI"))
		})

		It("should handle out of range scores", func() {
			// Test both too high and too low scores
			testCases := []struct {
				name  string
				score float64
			}{
				{"too high", 101},
				{"too low", -1},
			}

			for _, tc := range testCases {
				mockClient := &mockOpenAIClient{
					createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
						return openai.ChatCompletionResponse{
							Choices: []openai.ChatCompletionChoice{{
								Message: openai.ChatCompletionMessage{
									Content: fmt.Sprintf(`{"version": "1.0","scores": [{"post_id": "post1", "title": "Test Post", "score": %f, "reason": "Test reason"}]}`, tc.score),
								},
							}},
						}, nil
					},
				}
				s = &scorer{client: mockClient, config: Config{}, prompt: batchScorePrompt}

				// Verify through public interface that invalid scores are handled
				scored, err := s.ScorePosts(ctx, posts)
				Expect(err).To(HaveOccurred(), "case: %s", tc.name)
				Expect(scored).To(BeNil(), "case: %s", tc.name)
			}
		})

		It("should handle missing scores", func() {
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version": "1.0","scores": []}`,
							},
						}},
					}, nil
				},
			}
			s = &scorer{client: mockClient, config: Config{}, prompt: batchScorePrompt}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("missing score"))
		})

		It("should use custom prompt when provided", func() {
			customPrompt := "Custom prompt text %s"
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					// Verify custom prompt is used
					Expect(req.Messages[1].Content).To(ContainSubstring("Custom prompt text"))

					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version": "1.0","scores": [{"post_id": "post1", "title": "Test Post", "score": 85, "reason": "Custom prompt test"}]}`,
							},
						}},
					}, nil
				},
			}
			s = &scorer{
				client: mockClient,
				config: Config{},
				prompt: customPrompt,
			}

			scored, err := s.ScorePosts(ctx, posts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(1))
			Expect(scored[0].Score).To(Equal(85.0))
			Expect(scored[0].Reason).To(Equal("Custom prompt test"))
		})

		It("should handle more than maxBatchSize posts", func() {
			// Create 15 test posts
			largePosts := make([]reddit.Post, 15)
			for i := range largePosts {
				largePosts[i] = reddit.Post{
					ID:       fmt.Sprintf("post%d", i+1),
					Title:    fmt.Sprintf("Test Post %d", i+1),
					SelfText: "Content",
				}
			}

			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					// Determine which batch we're processing based on the content
					batchStart := 0
					if len(req.Messages) > 1 && req.Messages[1].Content != "" {
						if strings.Contains(req.Messages[1].Content, "post11") {
							batchStart = 10
						}
					}

					// Build response for current batch
					var scores strings.Builder
					scores.WriteString(`{"version": "1.0","scores": [`)

					batchEnd := batchStart + 10
					if batchEnd > 15 {
						batchEnd = 15
					}

					for i := batchStart; i < batchEnd; i++ {
						if i > batchStart {
							scores.WriteString(",")
						}
						fmt.Fprintf(&scores,
							`{"post_id": "post%d", "title": "Test Post %d", "score": %d, "reason": "Test reason %d"}`,
							i+1, i+1, 85-i*5, i+1)
					}
					scores.WriteString("]}")

					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: scores.String(),
							},
						}},
					}, nil
				},
			}
			s = &scorer{
				client: mockClient,
				config: Config{},
				prompt: batchScorePrompt,
			}

			scored, err := s.ScorePosts(ctx, largePosts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(15))

			for i, result := range scored {
				Expect(result.Post).To(Equal(largePosts[i]))
				Expect(result.Score).To(BeNumerically(">=", 0))
				Expect(result.Score).To(BeNumerically("<=", 100))
				Expect(result.Reason).NotTo(BeEmpty())
			}
		})
	})
})

// Helper functions
func newWithClient(client *mockOpenAIClient) Scorer {
	return &scorer{
		client: client,
		config: Config{},
		prompt: batchScorePrompt,
	}
}
