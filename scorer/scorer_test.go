package scorer_test

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

	"github.com/JohnPlummer/post-scorer/scorer"
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
		s          scorer.Scorer
		posts      []*reddit.Post
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockClient = &mockOpenAIClient{}
		s = scorer.NewWithClient(mockClient)
		posts = []*reddit.Post{
			{
				ID:       "post1",
				Title:    "Test Post 1",
				SelfText: "Test Content 1",
			},
		}
	})

	Context("New", func() {
		It("should return error when API key is missing", func() {
			_, err := scorer.New(scorer.Config{})
			Expect(err).To(Equal(scorer.ErrMissingAPIKey))
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

			s := scorer.NewWithClient(mockClient)
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
			s = scorer.NewWithClient(mockClient)
		})

		It("should return empty slice for empty posts", func() {
			scored, err := s.ScorePosts(ctx, []*reddit.Post{})
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(Equal([]*scorer.ScoredPost{}))
		})
		
		It("should return error for nil posts", func() {
			_, err := s.ScorePosts(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("posts cannot be nil"))
		})
		
		It("should return error for custom prompt without placeholder", func() {
			cfg := scorer.Config{
				OpenAIKey:  "test-key",
				PromptText: "Rate these posts without placeholder",
			}
			_, err := scorer.New(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must contain %s placeholder"))
		})
		
		It("should return error for negative MaxConcurrent", func() {
			cfg := scorer.Config{
				OpenAIKey:     "test-key",
				MaxConcurrent: -1,
			}
			_, err := scorer.New(cfg)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be non-negative"))
		})
		
		It("should return error for post with empty ID", func() {
			posts := []*reddit.Post{
				{ID: "", Title: "Test Post"},
			}
			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty ID"))
		})
		
		It("should return error for post that is nil", func() {
			posts := []*reddit.Post{nil}
			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("is nil"))
		})
		
		It("should work with MaxConcurrent set", func() {
			cfg := scorer.Config{
				OpenAIKey:     "test-key",
				MaxConcurrent: 2,
			}
			s, err := scorer.New(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(s).NotTo(BeNil())
		})
		
		It("should handle concurrent processing", func() {
			// Create more posts to trigger multiple batches
			manyPosts := make([]*reddit.Post, 25)
			for i := 0; i < 25; i++ {
				manyPosts[i] = &reddit.Post{
					ID:    fmt.Sprintf("post%d", i+1),
					Title: fmt.Sprintf("Test Post %d", i+1),
				}
			}

			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							{
								Message: openai.ChatCompletionMessage{
									Content: `{"version": "1.0", "scores": [{"post_id": "post1", "title": "Test Post 1", "score": 50, "reason": "Test reason"}]}`,
								},
							},
						},
					}, nil
				},
			}
			s := scorer.NewWithClient(mockClient)

			scored, err := s.ScorePosts(ctx, manyPosts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(scored)).To(Equal(25))
		})
		
		It("should actually use concurrent processing with MaxConcurrent > 1", func() {
			// Create enough posts to trigger multiple batches (15 posts = 2 batches)
			manyPosts := make([]*reddit.Post, 15)
			for i := 0; i < 15; i++ {
				manyPosts[i] = &reddit.Post{
					ID:    fmt.Sprintf("post%d", i+1),
					Title: fmt.Sprintf("Test Post %d", i+1),
				}
			}

			callCount := 0
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					callCount++
					// Generate responses for all posts in the batch
					var scores []string
					for i := 1; i <= 10; i++ { // Each batch has up to 10 posts
						if callCount == 1 && i <= 10 {
							scores = append(scores, fmt.Sprintf(`{"post_id": "post%d", "title": "Test Post %d", "score": 50, "reason": "Test reason"}`, i, i))
						} else if callCount == 2 && i <= 5 {
							scores = append(scores, fmt.Sprintf(`{"post_id": "post%d", "title": "Test Post %d", "score": 60, "reason": "Test reason"}`, i+10, i+10))
						}
					}
					
					content := fmt.Sprintf(`{"version": "1.0", "scores": [%s]}`, strings.Join(scores, ","))
					
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{
							{
								Message: openai.ChatCompletionMessage{
									Content: content,
								},
							},
						},
					}, nil
				},
			}
			
			// Create scorer with concurrent processing enabled
			s := scorer.NewWithClient(mockClient, 
				scorer.WithPrompt("Test prompt with %s"),
				scorer.WithMaxConcurrent(3))
			
			scored, err := s.ScorePosts(ctx, manyPosts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(scored)).To(Equal(15))
			Expect(callCount).To(Equal(2)) // Should make 2 batch calls
		})
		
		It("should respect context cancellation", func() {
			ctx, cancel := context.WithCancel(context.Background())
			
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					// Check if context is cancelled
					select {
					case <-ctx.Done():
						return openai.ChatCompletionResponse{}, ctx.Err()
					default:
						return openai.ChatCompletionResponse{
							Choices: []openai.ChatCompletionChoice{
								{
									Message: openai.ChatCompletionMessage{
										Content: `{"version": "1.0", "scores": [{"post_id": "post1", "title": "Test Post 1", "score": 50, "reason": "Test reason"}]}`,
									},
								},
							},
						}, nil
					}
				},
			}
			
			s := scorer.NewWithClient(mockClient)
			
			// Cancel the context before making the call
			cancel()
			
			posts := []*reddit.Post{
				{ID: "post1", Title: "Test Post 1"},
			}
			
			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("context canceled"))
		})

		It("should handle API errors", func() {
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					return openai.ChatCompletionResponse{}, errors.New("API error")
				},
			}
			s = scorer.NewWithClient(mockClient)

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
			s = scorer.NewWithClient(mockClient)

			_, err := s.ScorePosts(ctx, posts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("empty response with no choices"))
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
				s = scorer.NewWithClient(mockClient)

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
			s = scorer.NewWithClient(mockClient)

			scored, err := s.ScorePosts(ctx, posts)
			// Now we expect no error, but a post with score 0
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(1))
			Expect(scored[0].Score).To(Equal(0))
			Expect(scored[0].Reason).To(ContainSubstring("No score provided by model"))
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
			s = scorer.NewWithClient(mockClient, scorer.WithPrompt(customPrompt))

			scored, err := s.ScorePosts(ctx, posts)
			Expect(err).NotTo(HaveOccurred())
			Expect(scored).To(HaveLen(1))
			Expect(scored[0].Score).To(Equal(85))
			Expect(scored[0].Reason).To(Equal("Custom prompt test"))
		})

		It("should handle more than maxBatchSize posts", func() {
			// Create 15 test posts
			largePosts := make([]*reddit.Post, 15)
			for i := range largePosts {
				largePosts[i] = &reddit.Post{
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
			s = scorer.NewWithClient(mockClient)

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

	Describe("ScorePostsWithOptions", func() {
		It("should use custom model when WithModel option is provided", func() {
			var capturedModel string
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					capturedModel = req.Model
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version":"1.0","scores":[{"post_id":"post1","title":"Test Post 1","score":75,"reason":"Test"}]}`,
							},
						}},
					}, nil
				},
			}
			
			cfg := scorer.Config{
				OpenAIKey: "test-key",
				Model: "gpt-4", // Default model
			}
			s, err := scorer.New(cfg)
			Expect(err).NotTo(HaveOccurred())
			
			// Use internal test method to inject mock client
			type scorerWithClient interface {
				scorer.Scorer
				SetClient(client interface{})
			}
			if sc, ok := s.(scorerWithClient); ok {
				sc.SetClient(mockClient)
			} else {
				// If we can't inject, create with NewWithClient
				s = scorer.NewWithClient(mockClient)
			}
			
			// Test with custom model
			_, err = s.ScorePostsWithOptions(ctx, posts, scorer.WithModel("gpt-3.5-turbo"))
			Expect(err).NotTo(HaveOccurred())
			Expect(capturedModel).To(Equal("gpt-3.5-turbo"))
		})

		It("should use custom prompt template with WithPromptTemplate option", func() {
			var capturedPrompt string
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					// Capture the user message content
					for _, msg := range req.Messages {
						if msg.Role == openai.ChatMessageRoleUser {
							capturedPrompt = msg.Content
						}
					}
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version":"1.0","scores":[{"post_id":"post1","title":"Test Post 1","score":75,"reason":"Test"}]}`,
							},
						}},
					}, nil
				},
			}
			
			cfg := scorer.Config{
				OpenAIKey: "test-key",
			}
			s, err := scorer.New(cfg)
			Expect(err).NotTo(HaveOccurred())
			
			// Use NewWithClient for testing
			s = scorer.NewWithClient(mockClient)
			
			// Test with custom template
			customTemplate := "Custom prompt: {{.Posts}}"
			_, err = s.ScorePostsWithOptions(ctx, posts, scorer.WithPromptTemplate(customTemplate))
			Expect(err).NotTo(HaveOccurred())
			Expect(capturedPrompt).To(ContainSubstring("Custom prompt:"))
		})

		It("should pass extra context with WithExtraContext option", func() {
			var capturedPrompt string
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					for _, msg := range req.Messages {
						if msg.Role == openai.ChatMessageRoleUser {
							capturedPrompt = msg.Content
						}
					}
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version":"1.0","scores":[{"post_id":"post1","title":"Test Post 1","score":75,"reason":"Test"}]}`,
							},
						}},
					}, nil
				},
			}
			
			s := scorer.NewWithClient(mockClient)
			
			// Test with extra context
			extraContext := map[string]string{
				"City": "Brighton",
				"Type": "Gatekeeper",
			}
			customTemplate := "Scoring for {{.City}} - Type: {{.Type}}\nPosts: {{.Posts}}"
			
			_, err := s.ScorePostsWithOptions(ctx, posts, 
				scorer.WithPromptTemplate(customTemplate),
				scorer.WithExtraContext(extraContext))
			Expect(err).NotTo(HaveOccurred())
			Expect(capturedPrompt).To(ContainSubstring("Scoring for Brighton"))
			Expect(capturedPrompt).To(ContainSubstring("Type: Gatekeeper"))
		})
	})

	Describe("ScorePostsWithContext", func() {
		It("should score posts with additional context data", func() {
			mockClient := &mockOpenAIClient{
				response: openai.ChatCompletionResponse{
					Choices: []openai.ChatCompletionChoice{{
						Message: openai.ChatCompletionMessage{
							Content: `{"version":"1.0","scores":[{"post_id":"post1","title":"Test Post 1","score":85,"reason":"Has relevant comments"}]}`,
						},
					}},
				},
			}
			
			s := scorer.NewWithClient(mockClient)
			
			contexts := []scorer.ScoringContext{
				{
					Post: posts[0],
					ExtraData: map[string]string{
						"Comments": "Great place for coffee!",
					},
				},
			}
			
			scored, err := s.ScorePostsWithContext(ctx, contexts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(scored)).To(Equal(1))
			Expect(scored[0].Score).To(Equal(85))
			Expect(scored[0].Reason).To(Equal("Has relevant comments"))
		})

		It("should use template variables from context", func() {
			var capturedPrompt string
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					for _, msg := range req.Messages {
						if msg.Role == openai.ChatMessageRoleUser {
							capturedPrompt = msg.Content
						}
					}
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: `{"version":"1.0","scores":[{"post_id":"post1","title":"Test Post 1","score":85,"reason":"Test"}]}`,
							},
						}},
					}, nil
				},
			}
			
			s := scorer.NewWithClient(mockClient)
			
			contexts := []scorer.ScoringContext{
				{
					Post: posts[0],
					ExtraData: map[string]string{
						"Comments": "Great coffee shop recommendation",
					},
				},
			}
			
			// Template that uses context data
			template := "Post: {{range .Contexts}}{{.PostTitle}} - Comments: {{.Comments}}{{end}}"
			
			_, err := s.ScorePostsWithContext(ctx, contexts, scorer.WithPromptTemplate(template))
			Expect(err).NotTo(HaveOccurred())
			Expect(capturedPrompt).To(ContainSubstring("Test Post 1 - Comments: Great coffee shop recommendation"))
		})

		It("should handle empty contexts", func() {
			s := scorer.NewWithClient(mockClient)
			
			scored, err := s.ScorePostsWithContext(ctx, []scorer.ScoringContext{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(scored)).To(Equal(0))
		})

		It("should return error for nil contexts", func() {
			s := scorer.NewWithClient(mockClient)
			
			_, err := s.ScorePostsWithContext(ctx, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("contexts cannot be nil"))
		})

		It("should return error for context with nil post", func() {
			s := scorer.NewWithClient(mockClient)
			
			contexts := []scorer.ScoringContext{
				{
					Post:      nil,
					ExtraData: map[string]string{},
				},
			}
			
			_, err := s.ScorePostsWithContext(ctx, contexts)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("context at index 0 has nil post"))
		})

		It("should work with concurrent processing", func() {
			var callCount int
			mockClient := &mockOpenAIClient{
				createChatCompletionFunc: func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
					callCount++
					// Return appropriate number of scores based on request
					scores := `[{"post_id":"post1","title":"Test","score":70,"reason":"Test"}]`
					if callCount == 1 {
						scores = `[{"post_id":"post1","title":"Test","score":70,"reason":"Test"},{"post_id":"post2","title":"Test","score":70,"reason":"Test"},{"post_id":"post3","title":"Test","score":70,"reason":"Test"},{"post_id":"post4","title":"Test","score":70,"reason":"Test"},{"post_id":"post5","title":"Test","score":70,"reason":"Test"},{"post_id":"post6","title":"Test","score":70,"reason":"Test"},{"post_id":"post7","title":"Test","score":70,"reason":"Test"},{"post_id":"post8","title":"Test","score":70,"reason":"Test"},{"post_id":"post9","title":"Test","score":70,"reason":"Test"},{"post_id":"post10","title":"Test","score":70,"reason":"Test"}]`
					}
					return openai.ChatCompletionResponse{
						Choices: []openai.ChatCompletionChoice{{
							Message: openai.ChatCompletionMessage{
								Content: fmt.Sprintf(`{"version":"1.0","scores":%s}`, scores),
							},
						}},
					}, nil
				},
			}
			
			// Create scorer with concurrent processing
			s := scorer.NewWithClient(mockClient, scorer.WithMaxConcurrent(2))
			
			// Create many contexts
			var contexts []scorer.ScoringContext
			for i := 1; i <= 15; i++ {
				contexts = append(contexts, scorer.ScoringContext{
					Post: &reddit.Post{
						ID:    fmt.Sprintf("post%d", i),
						Title: fmt.Sprintf("Test Post %d", i),
					},
					ExtraData: map[string]string{
						"Comments": fmt.Sprintf("Comments for post %d", i),
					},
				})
			}
			
			scored, err := s.ScorePostsWithContext(ctx, contexts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(scored)).To(Equal(15))
			Expect(callCount).To(Equal(2)) // Should make 2 batch calls
		})
	})
})
