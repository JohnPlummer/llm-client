package scorer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog" // Using standard library slog
	"strings"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
)

// Scorer provides methods to score Reddit posts using ChatGPT
type Scorer interface {
	// ScorePosts evaluates and scores a slice of Reddit posts
	ScorePosts(ctx context.Context, posts []reddit.Post) ([]ScoredPost, error)
}

// ScoredPost represents a Reddit post with its AI-generated score
type ScoredPost struct {
	Post   reddit.Post
	Score  float64
	Reason string
	// We could add additional fields like:
	// Reasoning string    // explanation for the score
	// Confidence float64 // how confident the AI is in its scoring
}

// Config holds the configuration for the scorer
type Config struct {
	OpenAIKey     string
	PromptText    string // If empty, will use default prompt
	MaxConcurrent int    // for rate limiting
}

// OpenAIClient interface allows us to mock the OpenAI API
type OpenAIClient interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// scorer implements the Scorer interface
type scorer struct {
	client OpenAIClient
	config Config
	prompt string
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

// ErrMissingAPIKey is returned when no OpenAI API key is provided
var ErrMissingAPIKey = errors.New("OpenAI API key is required")

const (
	maxBatchSize = 10                                            // Maximum number of posts to process in one API call
	scorePrompt  = `Score each of the following Reddit posts...` // existing prompt
)

const batchScorePrompt = `Score each of the following Reddit post titles and output as JSON. Consider these categories:
- Regular venues (restaurants, bars, cafes, museums, galleries, etc.)
- Local attractions and points of interest
- Entertainment events (music, theatre, comedy, sports, etc.)
- Cultural events and festivals
- Markets and shopping areas
- Parks and outdoor spaces
- Family-friendly activities
- Seasonal or special events
- Hidden gems and local recommendations

Scoring guidelines:
90-100: Title directly references specific venues, events, or activities
70-89: Title suggests discussion of activities or places
40-69: Title might contain some relevant information
1-39: Title has low probability of relevant information
0: Title clearly indicates no relevant activity information

CRITICAL RULES:
1. Output must be ONLY valid JSON with no markdown or other formatting
2. Response must follow this exact format:
{
  "version": "1.0",
  "scores": [
    {
      "post_id": "<id>",
      "title": "<title>",
      "score": <0-100>,
      "reason": "<explanation>"
    }
  ]
}
3. Every post must receive a score and reason
4. Empty/invalid posts must get score 0
5. Never skip posts - score everything
6. Score must be between 0-100
7. Include clear reasoning for each score

Posts to score:
%s`

// Add new types for JSON parsing
type scoreResponse struct {
	Version string      `json:"version"`
	Scores  []scoreItem `json:"scores"`
}

type scoreItem struct {
	PostID string  `json:"post_id"`
	Title  string  `json:"title"`
	Score  float64 `json:"score"`
	Reason string  `json:"reason"`
}

func formatPostsForBatch(posts []reddit.Post) string {
	var sb strings.Builder
	for _, post := range posts {
		fmt.Fprintf(&sb, "%s %q: %s\n\n", post.ID, post.Title, post.SelfText)
	}
	return sb.String()
}

// Update parseBatchResponse to handle JSON
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

		// Optional: add delay between batches to respect rate limits
		// time.Sleep(100 * time.Millisecond)
	}

	return allResults, nil
}
