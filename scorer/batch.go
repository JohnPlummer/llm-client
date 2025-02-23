package scorer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func (s *scorer) processBatch(ctx context.Context, batch []*reddit.Post) ([]*ScoredPost, error) {
	prompt := fmt.Sprintf(s.prompt, formatPostsForBatch(batch))

	slog.Info("Processing batch of posts", "batch_size", len(batch))

	schema, err := jsonschema.GenerateSchemaForType(scoreResponse{})
	if err != nil {
		return nil, fmt.Errorf("generating schema: %w", err)
	}

	resp, err := s.createChatCompletion(ctx, prompt, schema)
	if err != nil {
		return nil, err
	}

	scores, err := s.parseResponse(resp)
	if err != nil {
		return nil, err
	}

	return s.createScoredPosts(batch, scores)
}

func (s *scorer) createChatCompletion(ctx context.Context, prompt string, schema *jsonschema.Definition) (*openai.ChatCompletionResponse, error) {
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
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return &resp, nil
}

func (s *scorer) parseResponse(resp *openai.ChatCompletionResponse) (map[string]scoreItem, error) {
	var result scoreResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	scores := make(map[string]scoreItem)
	for _, score := range result.Scores {
		if score.Score < 0 || score.Score > 100 {
			return nil, fmt.Errorf("invalid score %d for post %s: score must be between 0 and 100", score.Score, score.PostID)
		}
		scores[score.PostID] = score
	}

	return scores, nil
}

func (s *scorer) createScoredPosts(batch []*reddit.Post, scores map[string]scoreItem) ([]*ScoredPost, error) {
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

		slog.Debug("Post scored",
			"id", post.ID,
			"title", post.Title,
			"score", score.Score,
			"reason", score.Reason)
	}

	return results, nil
}

func formatPostsForBatch(posts []*reddit.Post) string {
	input := struct {
		Posts []*reddit.Post `json:"posts"`
	}{
		Posts: posts,
	}

	jsonData, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		slog.Error("failed to marshal posts", "error", err)
		return ""
	}

	return string(jsonData)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
