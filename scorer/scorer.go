package scorer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog" // Using standard library slog
	"strings"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
)

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

func formatPostsForBatch(posts []reddit.Post) string {
	var sb strings.Builder
	for _, post := range posts {
		fmt.Fprintf(&sb, "%s %q: %s\n\n", post.ID, post.Title, post.SelfText)
	}
	return sb.String()
}

func parseBatchResponse(response string, posts []reddit.Post) ([]ScoredPost, error) {
	var resp scoreResponse
	if err := json.Unmarshal([]byte(response), &resp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	if resp.Version != "1.0" {
		slog.Warn("Unexpected version in response", "version", resp.Version)
	}

	scores := make(map[string]scoreItem)
	for _, score := range resp.Scores {
		if score.Score < 0 || score.Score > 100 {
			return nil, fmt.Errorf("invalid score %f for post %s", score.Score, score.PostID)
		}
		scores[score.PostID] = score
	}

	results := make([]ScoredPost, len(posts))
	for i, post := range posts {
		score, exists := scores[post.ID]
		if !exists {
			return nil, fmt.Errorf("missing score for post %s: %q", post.ID, post.Title)
		}
		results[i] = ScoredPost{
			Post:   post,
			Score:  score.Score,
			Reason: score.Reason,
		}
	}

	return results, nil
}

// ScorePosts evaluates and scores a slice of Reddit posts
func (s *scorer) ScorePosts(ctx context.Context, posts []reddit.Post) ([]ScoredPost, error) {
	if len(posts) == 0 {
		return nil, nil
	}

	var allResults []ScoredPost

	// Process posts in batches of maxBatchSize
	for i := 0; i < len(posts); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(posts) {
			end = len(posts)
		}

		batch := posts[i:end]
		prompt := fmt.Sprintf(s.prompt, formatPostsForBatch(batch))

		slog.Debug("Sending prompt", "prompt", prompt)

		resp, err := s.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleSystem,
						Content: "You are a content analyzer focused on identifying posts containing location-based recommendations and events. Respond only with post IDs and scores.",
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
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

		slog.Debug("Raw response from OpenAI", "response", resp.Choices[0].Message.Content)

		batchResults, err := parseBatchResponse(resp.Choices[0].Message.Content, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch %d-%d: %w", i, end-1, err)
		}

		// Add debug logging for each scored post
		for _, post := range batchResults {
			slog.Debug("Post scored",
				"id", post.Post.ID,
				"title", post.Post.Title,
				"score", post.Score,
				"reason", post.Reason)
		}

		allResults = append(allResults, batchResults...)

	}

	return allResults, nil
}
