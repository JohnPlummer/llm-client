package scorer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	Post  reddit.Post
	Score float64
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

const batchScorePrompt = `Score each of the following Reddit posts on a scale of 0-100 based on how likely they contain information about:
- Events or activities
- Restaurant recommendations
- Bar recommendations
- Cafe/coffee shop recommendations
- Other actionable location-based activities

A score of 100 means the post definitely contains specific recommendations or event details.
A score of 0 means the post has no relevant recommendations or event information.

IMPORTANT: For each post, respond with ONLY the post ID, title, and a SINGLE overall score between 0 and 100.
Do not break down the score by category.

IMPORTANT: You MUST provide a score for EVERY post in the input list. Do not skip any posts.
If a post is not relevant, give it a score of 0, but still include it in the response.

Example format (exactly like this):
abc123 "Example Title": 85
def456 "Another Title": 30
ghi789 "Third Title": 95

Posts to score:
%s`

func formatPostsForBatch(posts []reddit.Post) string {
	var sb strings.Builder
	for _, post := range posts {
		fmt.Fprintf(&sb, "%s %q: %s\n\n", post.ID, post.Title, post.SelfText)
	}
	return sb.String()
}

func parseBatchResponse(response string, posts []reddit.Post) ([]ScoredPost, error) {
	scores := make(map[string]float64)

	lines := strings.Split(strings.TrimSpace(response), "\n")
	for _, line := range lines {
		// Find the position of the last colon (after the title)
		lastColon := strings.LastIndex(line, ":")
		if lastColon == -1 {
			continue
		}

		// Extract the score part
		scoreStr := strings.TrimSpace(line[lastColon+1:])
		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil || score < 0 || score > 100 {
			return nil, fmt.Errorf("invalid score in line %q: %s", line, scoreStr)
		}

		// Extract the ID (it's before the first quote)
		firstQuote := strings.Index(line, "\"")
		if firstQuote == -1 {
			continue
		}
		postID := strings.TrimSpace(line[:firstQuote])

		scores[postID] = score
	}

	// Create results and verify all posts were scored
	results := make([]ScoredPost, len(posts))
	for i, post := range posts {
		score, exists := scores[post.ID]
		if !exists {
			return nil, fmt.Errorf("missing score for post %s: %q", post.ID, post.Title)
		}
		results[i] = ScoredPost{
			Post:  post,
			Score: score,
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

		batchResults, err := parseBatchResponse(resp.Choices[0].Message.Content, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to parse batch %d-%d: %w", i, end-1, err)
		}

		allResults = append(allResults, batchResults...)

		// Optional: add delay between batches to respect rate limits
		// time.Sleep(100 * time.Millisecond)
	}

	return allResults, nil
}
