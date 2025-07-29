package scorer

import (
	"context"
	"errors"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
)

// Scorer provides methods to score Reddit posts using ChatGPT
type Scorer interface {
	ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error)
	ScorePostsWithOptions(ctx context.Context, posts []*reddit.Post, opts ...ScoringOption) ([]*ScoredPost, error)
}

// ScoredPost represents a Reddit post with its AI-generated score
type ScoredPost struct {
	Post   *reddit.Post
	Score  int
	Reason string
}

// Config holds the configuration for the scorer
type Config struct {
	OpenAIKey     string
	Model         string
	PromptText    string
	MaxConcurrent int
}

// OpenAIClient defines the interface for interacting with OpenAI API
type OpenAIClient interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

type scorer struct {
	client OpenAIClient
	config Config
	prompt string
}

// ErrMissingAPIKey is returned when no OpenAI API key is provided
var ErrMissingAPIKey = errors.New("OpenAI API key is required")

type scoreResponse struct {
	Version string      `json:"version"`
	Scores  []scoreItem `json:"scores"`
}

type scoreItem struct {
	PostID string `json:"post_id"`
	Title  string `json:"title"`
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

// ScoringOption is a functional option for configuring scoring behavior
type ScoringOption func(*scoringOptions)

// scoringOptions holds the options for a scoring request
type scoringOptions struct {
	model string
}

// WithModel sets the model for this scoring request
func WithModel(model string) ScoringOption {
	return func(opts *scoringOptions) {
		opts.model = model
	}
}
