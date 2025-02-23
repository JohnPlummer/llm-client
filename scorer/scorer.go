package scorer

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

//go:embed prompts/*.txt
var promptFS embed.FS

var systemPrompt string

func init() {
	// Load system prompt during package initialization
	promptBytes, err := promptFS.ReadFile("prompts/system_prompt.txt")
	if err != nil {
		slog.Error("failed to load system prompt", "error", err)
		os.Exit(1)
	}
	systemPrompt = string(promptBytes)
}

// New creates a new instance of the Scorer
func New(cfg Config) (Scorer, error) {
	if cfg.OpenAIKey == "" {
		return nil, ErrMissingAPIKey
	}

	prompt := batchScorePrompt
	if cfg.PromptText != "" {
		prompt = cfg.PromptText
	}

	client := openai.NewClient(cfg.OpenAIKey)
	return &scorer{
		client: client,
		config: cfg,
		prompt: prompt,
	}, nil
}

// WithPrompt returns a function that sets a custom prompt for the scorer
func WithPrompt(prompt string) func(*scorer) {
	return func(s *scorer) {
		s.prompt = prompt
	}
}

// NewWithClient creates a new scorer with a custom OpenAI client and options
func NewWithClient(client OpenAIClient, opts ...func(*scorer)) Scorer {
	s := &scorer{
		client: client,
		prompt: batchScorePrompt,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func formatPostsForBatch(posts []*reddit.Post) string {
	input := struct {
		Posts []*reddit.Post `json:"posts"`
	}{
		Posts: posts,
	}

	// Using json.MarshalIndent for better readability in the prompt
	jsonData, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		// This should never happen when marshaling reddit.Post
		slog.Error("failed to marshal posts", "error", err)
		return ""
	}

	return string(jsonData)
}

// ScorePosts evaluates and scores a slice of Reddit posts
func (s *scorer) ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error) {
	if len(posts) == 0 {
		return nil, nil
	}

	var allResults []*ScoredPost

	// Process posts in batches of maxBatchSize
	for i := 0; i < len(posts); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(posts) {
			end = len(posts)
		}

		batch := posts[i:end]
		prompt := fmt.Sprintf(s.prompt, formatPostsForBatch(batch))

		// Log start of batch processing
		slog.Info("Processing batch of posts",
			"batch_start", i,
			"batch_end", end-1,
			"batch_size", len(batch))

		// Generate schema from our response type
		schema, err := jsonschema.GenerateSchemaForType(scoreResponse{})
		if err != nil {
			return nil, fmt.Errorf("failed to generate schema: %w", err)
		}

		resp, err := s.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT4oMini,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: systemPrompt,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
				ResponseFormat: &openai.ChatCompletionResponseFormat{
					Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
					JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
						Schema: schema,
						Name:   "post_scoring",
					},
				},
			},
		)
		if err != nil {
			return nil, fmt.Errorf("OpenAI API error in batch %d-%d: %w", i, end-1, err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no response from OpenAI for batch %d-%d", i, end-1)
		}

		var result scoreResponse
		if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		scores := make(map[string]scoreItem)
		for _, score := range result.Scores {
			// Validate that scores are within the required range
			if score.Score < 0 || score.Score > 100 {
				return nil, fmt.Errorf("invalid score %d for post %s: score must be between 0 and 100", score.Score, score.PostID)
			}
			scores[score.PostID] = score
		}

		results := make([]*ScoredPost, len(batch))
		for i, post := range batch {
			score, exists := scores[post.ID]
			if !exists {
				return nil, fmt.Errorf("missing score for post %s: %q", post.ID, post.Title)
			}
			results[i] = &ScoredPost{
				Post:   post,
				Score:  score.Score,
				Reason: score.Reason,
			}
		}

		// Log batch scoring results
		slog.Info("Batch scoring completed",
			"batch_start", i,
			"batch_end", end-1,
			"posts_scored", len(results))

		// Log individual post scores at info level
		for _, post := range results {
			slog.Debug("Post scored",
				"id", post.Post.ID,
				"title", post.Post.Title,
				"score", post.Score,
				"reason", post.Reason)
		}

		allResults = append(allResults, results...)
	}

	slog.Info("All posts scored successfully",
		"total_posts", len(posts),
		"total_batches", (len(posts)+maxBatchSize-1)/maxBatchSize)

	return allResults, nil
}
