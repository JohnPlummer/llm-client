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
	MaxConcurrent int // for rate limiting
}

// OpenAIClient interface allows us to mock the OpenAI API
type OpenAIClient interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// scorer implements the Scorer interface
type scorer struct {
	client OpenAIClient
	config Config
}

// New creates a new instance of the Scorer
func New(cfg Config) (Scorer, error) {
	client := openai.NewClient(cfg.OpenAIKey)

	return &scorer{
		client: client,
		config: cfg,
	}, nil
}

// ErrMissingAPIKey is returned when no OpenAI API key is provided
var ErrMissingAPIKey = errors.New("OpenAI API key is required")

const scorePrompt = `Evaluate this Reddit post on a scale of 0-100 based on how likely it contains information about:
- Events or activities
- Restaurant recommendations
- Bar recommendations
- Cafe/coffee shop recommendations
- Other actionable location-based activities

A score of 100 means the post definitely contains specific recommendations or event details.
A score of 0 means the post has no relevant recommendations or event information.
Respond with only a number between 0 and 100.

Post Title: %s
Post Content: %s`

// ScorePosts evaluates and scores a slice of Reddit posts
func (s *scorer) ScorePosts(ctx context.Context, posts []reddit.Post) ([]ScoredPost, error) {
	results := make([]ScoredPost, len(posts))
	for i, post := range posts {
		score, err := s.scorePost(ctx, post)
		if err != nil {
			return nil, err
		}
		results[i] = ScoredPost{Post: post, Score: score}
	}
	return results, nil
}

func (s *scorer) scorePost(ctx context.Context, post reddit.Post) (float64, error) {
	prompt := fmt.Sprintf(scorePrompt, post.Title, post.SelfText)

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a content analyzer focused on identifying posts containing location-based recommendations and events.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return 0, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return 0, errors.New("no response from OpenAI")
	}

	scoreStr := strings.TrimSpace(resp.Choices[0].Message.Content)
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse score '%s': %w", scoreStr, err)
	}

	if score < 0 || score > 100 {
		return 0, fmt.Errorf("score %f out of range (0-100)", score)
	}

	return score, nil
}
