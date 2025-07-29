package scorer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/JohnPlummer/reddit-client/reddit"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func (s *scorer) processBatch(ctx context.Context, batch []*reddit.Post, options *scoringOptions) ([]*ScoredPost, error) {
	return s.processBatchWithContext(ctx, batch, nil, options)
}

func (s *scorer) processBatchWithContext(ctx context.Context, batch []*reddit.Post, contexts []ScoringContext, options *scoringOptions) ([]*ScoredPost, error) {
	// Determine which prompt to use
	promptText := s.prompt
	if options != nil && options.promptText != "" {
		promptText = options.promptText
	}
	
	// Format the prompt with appropriate data
	var prompt string
	var err error
	if contexts != nil {
		prompt, err = s.formatPromptWithContext(promptText, contexts, options)
	} else {
		prompt, err = s.formatPrompt(promptText, batch, options)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to format prompt: %w", err)
	}

	slog.Info("Processing batch of posts", "batch_size", len(batch))

	schema, err := jsonschema.GenerateSchemaForType(scoreResponse{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate JSON schema for batch of %d posts: %w", len(batch), err)
	}

	resp, err := s.createChatCompletion(ctx, prompt, schema, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion for batch of %d posts: %w", len(batch), err)
	}

	scores, err := s.parseResponse(resp, len(batch))
	if err != nil {
		return nil, fmt.Errorf("failed to parse response for batch of %d posts: %w", len(batch), err)
	}

	return s.createScoredPosts(batch, scores)
}

func (s *scorer) buildChatRequest(prompt string, schema *jsonschema.Definition, options *scoringOptions) openai.ChatCompletionRequest {
	// Determine which model to use
	model := s.config.Model
	if model == "" {
		model = openai.GPT4oMini // Default model
	}
	if options != nil && options.model != "" {
		model = options.model // Override with options
	}
	
	return openai.ChatCompletionRequest{
		Model: model,
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
	}
}

func (s *scorer) createChatCompletion(ctx context.Context, prompt string, schema *jsonschema.Definition, options *scoringOptions) (*openai.ChatCompletionResponse, error) {
	req := s.buildChatRequest(prompt, schema, options)
	
	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI returned empty response with no choices")
	}

	return &resp, nil
}

func (s *scorer) validateResponse(result scoreResponse, expectedPostCount int) error {
	// Check if we received scores for all expected posts
	if len(result.Scores) < expectedPostCount {
		slog.WarnContext(context.Background(), "incomplete scores from OpenAI",
			"expected_count", expectedPostCount,
			"received_count", len(result.Scores))
	}

	for _, score := range result.Scores {
		if score.Score < 0 || score.Score > 100 {
			return fmt.Errorf("invalid score %d for post %s: score must be between 0 and 100", score.Score, score.PostID)
		}
	}
	
	return nil
}

func (s *scorer) parseResponse(resp *openai.ChatCompletionResponse, expectedPostCount int) (map[string]scoreItem, error) {
	var result scoreResponse
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OpenAI response (expected %d posts): %w", expectedPostCount, err)
	}

	if err := s.validateResponse(result, expectedPostCount); err != nil {
		return nil, err
	}

	scores := make(map[string]scoreItem)
	for _, score := range result.Scores {
		scores[score.PostID] = score
	}

	return scores, nil
}

func (s *scorer) createScoredPosts(batch []*reddit.Post, scores map[string]scoreItem) ([]*ScoredPost, error) {
	results := make([]*ScoredPost, len(batch))
	for i, post := range batch {
		score, exists := scores[post.ID]
		if !exists {
			slog.WarnContext(context.Background(), "missing score from OpenAI response, assigning default score",
				"post_id", post.ID,
				"title", post.Title)

			results[i] = &ScoredPost{
				Post:   post,
				Score:  0,
				Reason: "No score provided by model - automatically assigned lowest relevance score",
			}
			continue
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

// formatPrompt formats a prompt with template support
func (s *scorer) formatPrompt(promptText string, posts []*reddit.Post, options *scoringOptions) (string, error) {
	// Check if prompt contains template syntax
	if strings.Contains(promptText, "{{") && strings.Contains(promptText, "}}") {
		// Use template processing
		tmpl, err := template.New("prompt").Parse(promptText)
		if err != nil {
			return "", fmt.Errorf("failed to parse prompt template: %w", err)
		}
		
		// Create template data
		data := map[string]interface{}{
			"Posts": formatPostsForBatch(posts),
		}
		
		// Add extra context from options if available
		if options != nil && options.extraContext != nil {
			for k, v := range options.extraContext {
				data[k] = v
			}
		}
		
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to execute prompt template: %w", err)
		}
		
		return buf.String(), nil
	}
	
	// Fall back to simple sprintf for backward compatibility
	if strings.Contains(promptText, "%s") {
		return fmt.Sprintf(promptText, formatPostsForBatch(posts)), nil
	}
	
	// No placeholders, return as-is
	return promptText, nil
}

// formatPromptWithContext formats a prompt with scoring contexts
func (s *scorer) formatPromptWithContext(promptText string, contexts []ScoringContext, options *scoringOptions) (string, error) {
	// Convert contexts to posts for JSON formatting
	posts := make([]*reddit.Post, len(contexts))
	for i, ctx := range contexts {
		posts[i] = ctx.Post
	}
	
	// Check if prompt contains template syntax
	if strings.Contains(promptText, "{{") && strings.Contains(promptText, "}}") {
		// Use template processing
		tmpl, err := template.New("prompt").Parse(promptText)
		if err != nil {
			return "", fmt.Errorf("failed to parse prompt template: %w", err)
		}
		
		// Create template data with contexts
		contextData := make([]map[string]interface{}, len(contexts))
		for i, ctx := range contexts {
			postData := map[string]interface{}{
				"ID":       ctx.Post.ID,
				"Title":    ctx.Post.Title,
				"Body":     ctx.Post.SelfText,
				"PostTitle": ctx.Post.Title,
				"PostBody":  ctx.Post.SelfText,
			}
			
			// Add extra data from context
			for k, v := range ctx.ExtraData {
				postData[k] = v
			}
			
			contextData[i] = postData
		}
		
		data := map[string]interface{}{
			"Posts":    formatPostsForBatch(posts),
			"Contexts": contextData,
		}
		
		// Add extra context from options if available
		if options != nil && options.extraContext != nil {
			for k, v := range options.extraContext {
				data[k] = v
			}
		}
		
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", fmt.Errorf("failed to execute prompt template: %w", err)
		}
		
		return buf.String(), nil
	}
	
	// Fall back to simple formatting
	return s.formatPrompt(promptText, posts, options)
}

