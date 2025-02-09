package scorer

import (
	"context"
	"errors"
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
	response openai.ChatCompletionResponse
	err      error
}

func (m *mockOpenAIClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
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
							Content: "85",
						},
					},
				},
			}

			scored, err := s.ScorePosts(ctx, posts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(1))
			Expect(scored[0].Score).To(Equal(85.0))
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
							Content: "invalid",
						},
					},
				},
			}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(MatchError(ContainSubstring("failed to parse score")))
		})

		It("should handle out of range scores", func() {
			mockClient.response = openai.ChatCompletionResponse{
				Choices: []openai.ChatCompletionChoice{
					{
						Message: openai.ChatCompletionMessage{
							Content: "150",
						},
					},
				},
			}

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(MatchError(ContainSubstring("out of range")))
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
