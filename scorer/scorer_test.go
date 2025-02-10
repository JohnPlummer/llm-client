package scorer

import (
	"context"
	"errors"
	"fmt"
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

// mockOpenAIClient implements the openAIClient interface for testing
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
				ID:       "123",
				Title:    "Best restaurants in town",
				SelfText: "Check out these amazing places...",
			},
		}
	})

	Context("ScorePosts", func() {
		It("should successfully score posts", func() {
			mockClient.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: `123 "Best restaurants in town": 85 [Restaurant recommendations]`,
						},
					},
				},
			}

			scored, err := s.ScorePosts(ctx, posts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(1))
			Expect(scored[0].Score).To(Equal(85.0))
			Expect(scored[0].Reason).To(Equal("Restaurant recommendations"))
		})

		It("should handle API errors", func() {
			mockClient.err = errors.New("API error")

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(MatchError(ContainSubstring("API error")))
		})

		It("should handle invalid scores", func() {
			mockClient.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: `123 "Best restaurants in town": invalid [Test reason]`,
						},
					},
				},
			}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(MatchError(ContainSubstring("invalid score")))
		})

		It("should handle out of range scores", func() {
			mockClient.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: `123 "Best restaurants in town": 150 [Test reason]`,
						},
					},
				},
			}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(MatchError(ContainSubstring("invalid score")))
		})

		It("should successfully score posts with reasons", func() {
			mockClient.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: `123 "Best restaurants in town": 85 [Contains specific restaurant recommendations]`,
						},
					},
				},
			}

			scored, err := s.ScorePosts(ctx, posts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(1))
			Expect(scored[0].Score).To(Equal(85.0))
			Expect(scored[0].Reason).To(Equal("Contains specific restaurant recommendations"))
		})

		It("should handle missing reason brackets", func() {
			mockClient.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: `123 "Best restaurants in town": 85`,
						},
					},
				},
			}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(MatchError(ContainSubstring("without reason")))
		})

		It("should handle more than maxBatchSize posts", func() {
			// 1. Setup: Create a slice of 15 posts
			largePosts := make([]reddit.Post, 15)
			for i := range largePosts {
				largePosts[i] = reddit.Post{
					ID:       fmt.Sprintf("post%d", i+1),
					Title:    fmt.Sprintf("Post %d", i+1),
					SelfText: "Content",
				}
			}

			// 2. Setup responses for both batches
			responses := []openai.ChatCompletionResponse{
				{
					// First batch response (posts 1-10)
					Choices: []openai.ChatCompletionChoice{{
						Message: openai.ChatCompletionMessage{
							Content: `post1 "Post 1": 85 [Event listing]
post2 "Post 2": 70 [Restaurant review]
post3 "Post 3": 60 [General discussion]
post4 "Post 4": 55 [Partial information]
post5 "Post 5": 50 [Some relevance]
post6 "Post 6": 45 [Limited details]
post7 "Post 7": 40 [Vague content]
post8 "Post 8": 35 [Minimal relevance]
post9 "Post 9": 30 [Few details]
post10 "Post 10": 25 [Not very relevant]`,
						},
					}},
				},
				{
					// Second batch response (posts 11-15)
					Choices: []openai.ChatCompletionChoice{{
						Message: openai.ChatCompletionMessage{
							Content: `post11 "Post 11": 20 [Low relevance]
post12 "Post 12": 15 [Very low relevance]
post13 "Post 13": 10 [Almost irrelevant]
post14 "Post 14": 5 [Barely relevant]
post15 "Post 15": 0 [Not relevant]`,
						},
					}},
				},
			}

			// 3. Create a mock client that returns different responses for each batch
			responseIndex := 0
			mockClient.createChatCompletionFunc = func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
				resp := responses[responseIndex]
				responseIndex++
				return resp, nil
			}

			// 4. Score all posts in one call
			scored, err := s.ScorePosts(ctx, largePosts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(15))

			// 5. Verify scores are as expected
			Expect(scored[0].Score).To(Equal(85.0))
			Expect(scored[14].Score).To(Equal(0.0))

			// Add reason verification
			Expect(scored[0].Reason).To(Equal("Event listing"))
			Expect(scored[14].Reason).To(Equal("Not relevant"))
		})

		It("should use custom prompt when provided", func() {
			customPrompt := "Custom prompt text %s"
			var capturedMessages []openai.ChatCompletionMessage

			mockClient.createChatCompletionFunc = func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
				capturedMessages = req.Messages
				return openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{{
						Message: openai.ChatCompletionMessage{
							Content: `123 "Best restaurants in town": 85 [Custom prompt test]`,
						},
					}},
				}, nil
			}

			s = newWithConfig(mockClient, Config{
				PromptText: customPrompt,
			})

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).NotTo(HaveOccurred())

			// Updated to check the user message (index 1) instead of system message (index 0)
			Expect(capturedMessages[1].Content).To(ContainSubstring("Custom prompt text"))
			Expect(capturedMessages[1].Content).NotTo(ContainSubstring("Score each of the following Reddit posts"))
		})
	})
})

// newWithClient creates a new scorer with a custom OpenAI client for testing
func newWithClient(client *mockOpenAIClient) Scorer {
	return &scorer{
		client: client,
		config: Config{},
	}
}

// Update helper function to properly initialize the scorer with custom prompt
func newWithConfig(client *mockOpenAIClient, cfg Config) Scorer {
	if cfg.PromptText == "" {
		cfg.PromptText = batchScorePrompt
	}
	return &scorer{
		client: client,
		config: cfg,
		prompt: cfg.PromptText,
	}
}
