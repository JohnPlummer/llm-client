package scorer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func (s *scorer) processBatch(ctx context.Context, batch []TextItem, options *scoringOptions) ([]ScoredItem, error) {
	// Determine which prompt to use
	promptText := s.prompt
	if options != nil && options.promptText != "" {
		promptText = options.promptText
	}
	
	// Format the prompt with appropriate data
	prompt, err := s.formatPrompt(promptText, batch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to format prompt: %w", err)
	}

	slog.Info("Processing batch of text items", "batch_size", len(batch))

	schema, err := jsonschema.GenerateSchemaForType(scoreResponse{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate JSON schema for batch of %d items: %w", len(batch), err)
	}

	resp, err := s.createChatCompletion(ctx, prompt, schema, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion for batch of %d items: %w", len(batch), err)
	}

	// Parse response
	content := resp.Choices[0].Message.Content
	
	slog.Debug("Received response from OpenAI", "content_length", len(content))
	
	var scores scoreResponse
	if err := json.Unmarshal([]byte(content), &scores); err != nil {
		slog.Error("Failed to parse response JSON", "error", err, "content", content)
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	slog.Info("Received scores from OpenAI", "scores_count", len(scores.Scores))

	// Map scores back to items
	return s.mapScoresToItems(batch, scores.Scores), nil
}

func (s *scorer) createChatCompletion(ctx context.Context, prompt string, schema *jsonschema.Definition, options *scoringOptions) (openai.ChatCompletionResponse, error) {
	// Determine model to use
	model := s.config.Model
	if model == "" {
		model = openai.GPT4oMini
	}
	if options != nil && options.model != "" {
		model = options.model
	}

	request := openai.ChatCompletionRequest{
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
				Name:   "score_response",
				Strict: true,
				Schema: schema,
			},
		},
	}

	slog.Debug("Sending request to OpenAI", "model", model, "prompt_length", len(prompt))

	return s.client.CreateChatCompletion(ctx, request)
}

func (s *scorer) mapScoresToItems(items []TextItem, scores []scoreItem) []ScoredItem {
	scoreMap := make(map[string]scoreItem)
	for _, score := range scores {
		scoreMap[score.ItemID] = score
	}

	results := make([]ScoredItem, len(items))
	for i, item := range items {
		if score, found := scoreMap[item.ID]; found {
			// Validate score range
			if score.Score < 0 || score.Score > 100 {
				slog.Warn("Score out of range, clamping to valid range",
					"item_id", item.ID,
					"original_score", score.Score)
				if score.Score < 0 {
					score.Score = 0
				} else if score.Score > 100 {
					score.Score = 100
				}
			}
			
			results[i] = ScoredItem{
				Item:   item,
				Score:  score.Score,
				Reason: score.Reason,
			}
			slog.Debug("Mapped score to item",
				"item_id", item.ID,
				"score", score.Score)
		} else {
			slog.Warn("Score not found for item, using default",
				"item_id", item.ID)
			results[i] = ScoredItem{
				Item:   item,
				Score:  0,
				Reason: "Score not found in response",
			}
		}
	}
	
	return results
}

func (s *scorer) formatPrompt(promptText string, items []TextItem, options *scoringOptions) (string, error) {
	// Check if prompt uses Go template syntax
	if strings.Contains(promptText, "{{") && strings.Contains(promptText, "}}") {
		return s.formatPromptWithTemplate(promptText, items, options)
	}
	
	// Legacy sprintf-style formatting
	if strings.Contains(promptText, "%s") {
		itemsText := s.formatItemsAsText(items)
		return fmt.Sprintf(promptText, itemsText), nil
	}
	
	// If no placeholders, append items to the prompt
	itemsText := s.formatItemsAsText(items)
	return fmt.Sprintf("%s\n\nItems to score:\n%s", promptText, itemsText), nil
}

func (s *scorer) formatPromptWithTemplate(promptText string, items []TextItem, options *scoringOptions) (string, error) {
	tmpl, err := template.New("prompt").Parse(promptText)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template: %w", err)
	}
	
	// Prepare template data
	data := map[string]interface{}{
		"Items": items,
		"Count": len(items),
	}
	
	// Add extra context if provided
	if options != nil && options.extraContext != nil {
		for k, v := range options.extraContext {
			data[k] = v
		}
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// Provide helpful error message with template preview
		preview := promptText
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		return "", fmt.Errorf("failed to execute template '%s': %w", preview, err)
	}
	
	return buf.String(), nil
}

func (s *scorer) formatItemsAsText(items []TextItem) string {
	var sb strings.Builder
	for i, item := range items {
		sb.WriteString(fmt.Sprintf("Item %d (ID: %s):\n", i+1, item.ID))
		sb.WriteString(item.Content)
		if item.Metadata != nil && len(item.Metadata) > 0 {
			sb.WriteString("\nMetadata: ")
			for k, v := range item.Metadata {
				sb.WriteString(fmt.Sprintf("%s=%v ", k, v))
			}
		}
		sb.WriteString("\n\n")
	}
	return sb.String()
}