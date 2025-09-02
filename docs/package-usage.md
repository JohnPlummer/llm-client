# Package Usage

## Installation

```bash
go get github.com/JohnPlummer/llm-client@v0.9.0
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/JohnPlummer/llm-client/scorer"
)

func main() {
    // Initialize the scorer
    s, err := scorer.NewScorer(scorer.Config{
        OpenAIKey: "your-openai-api-key",
    })
    if err != nil {
        panic(err)
    }

    // Prepare text items to score
    items := []scorer.TextItem{
        {
            ID:      "item1",
            Content: "Best restaurants in downtown? Looking for good dinner spots",
            Metadata: map[string]interface{}{"title": "Best restaurants in downtown?"},
        },
        {
            ID:      "item2",
            Content: "Weekend events happening? What's going on this weekend?",
            Metadata: map[string]interface{}{"title": "Weekend events happening?"},
        },
    }

    // Score the text items
    results, err := s.ScoreTexts(context.Background(), items)
    if err != nil {
        panic(err)
    }

    // Use the results
    for _, result := range results {
        title := result.Item.Metadata["title"]
        fmt.Printf("Title: %s\nScore: %d\nReason: %s\n\n", 
            title, 
            result.Score, 
            result.Reason)
    }
}
```

## Configuration Options

### Config Structure

```go
type Config struct {
    OpenAIKey     string // Required: Your OpenAI API key
    PromptText    string // Optional: Custom prompt template
    MaxConcurrent int    // Optional: Rate limiting (future feature)
}
```

### Configuration Examples

#### Default Configuration
```go
config := scorer.Config{
    OpenAIKey: os.Getenv("OPENAI_API_KEY"),
}
```

#### Custom Prompt Configuration
```go
customPrompt := `Score posts for tech events (0-100)...
Posts to score: %s`

config := scorer.Config{
    OpenAIKey:  os.Getenv("OPENAI_API_KEY"),
    PromptText: customPrompt,
}
```

## Advanced Usage

### Dependency Injection

For testing or custom OpenAI client configuration:

```go
// Create custom client
client := openai.NewClient("your-api-key")

// Use dependency injection
scorer := scorer.NewWithClient(client, scorer.WithPrompt("custom prompt"))
```

### With Custom Logging

```go
import "log/slog"

// Set up structured logging
opts := &slog.HandlerOptions{Level: slog.LevelDebug}
logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
slog.SetDefault(logger)

// Scorer will use the configured logger
s, err := scorer.NewScorer(config)
```

## Input Data Structure

### TextItem Structure

The library expects text items conforming to the `TextItem` type:

```go
type TextItem struct {
    ID       string                 // Required: Unique identifier
    Content  string                 // Required: Text content to score
    Metadata map[string]interface{} // Optional: Additional context
}
```

### Metadata Integration

Metadata provides additional context for scoring:

```go
item := scorer.TextItem{
    ID:      "item1",
    Content: "Restaurant recommendations? Looking for good food",
    Metadata: map[string]interface{}{
        "title": "Restaurant recommendations?",
        "comments": []string{
            "Try Luigi's on Main Street!",
            "The corner café has great coffee",
        },
        },
        "location": "downtown",
    },
}
```

## Output Structure

### ScoredItem Type

```go
type ScoredItem struct {
    Item   TextItem // Original text item data
    Score  int      // AI-generated score (0-100)
    Reason string   // Explanation for the score
}
```

### Score Interpretation

- **90-100**: Highly relevant content with specific venue/event details
- **70-89**: Likely relevant with activity or location discussion
- **40-69**: Potentially relevant but requires investigation
- **1-39**: Low relevance probability
- **0**: No relevant information or processing failed

## Custom Prompts

### Prompt Requirements

Custom prompts **MUST** include:

1. **Text placeholder**: Include `%s` where text items will be injected
2. **JSON output requirement**: Specify exact JSON structure needed
3. **Scoring criteria**: Define clear 0-100 scoring guidelines
4. **Complete coverage**: Require scoring of ALL input text items

### JSON Response Format

Your prompt must instruct the AI to return JSON in this exact structure:

```json
{
  "version": "1.0",
  "scores": [
    {
      "item_id": "<item_id>",
      "score": <0-100>,
      "reason": "<explanation>"
    }
  ]
}
```

### Example Custom Prompt

```text
Score text items for technology events and meetups (0-100 scale).

Categories to consider:
- Tech conferences and workshops
- Developer meetups and hackathons
- Product launches and demos
- Startup events and networking

Scoring guidelines:
90-100: Specific tech event details with venue/time
70-89: Discussion of tech activities or gatherings
40-69: Tech-related but unclear about events
1-39: Minimal tech relevance
0: No tech event content

CRITICAL: Score every text item, return only valid JSON.

Text items to score:
%s
```

### Validation Requirements

- **No markdown formatting** in JSON output
- **All fields required**: item_id, score, reason
- **Score range**: Must be 0-100 integer
- **Complete coverage**: Every input text item must receive a score
- **Valid JSON**: Must parse without errors

## Error Handling

### Common Error Scenarios

```go
// Missing API key
s, err := scorer.NewScorer(scorer.Config{})
if err == scorer.ErrMissingAPIKey {
    // Handle missing API key
}

// API failures
results, err := s.ScoreTexts(ctx, items)
if err != nil {
    // Check if it's a context cancellation
    if errors.Is(err, context.Canceled) {
        // Handle timeout
    }
    // Handle other API errors
}
```

### Graceful Degradation

The library provides automatic fallbacks:

- **Missing scores**: Automatically assigned score of 0
- **Invalid score ranges**: Validation error returned
- **Partial responses**: Warning logged, processing continues
- **API timeouts**: Error propagated with context

## Performance Considerations

### Batch Processing

- **Automatic batching**: Text items are processed in batches of 10
- **Sequential processing**: Batches are processed one at a time
- **Error isolation**: Batch failures don't affect other batches

### Rate Limiting

```go
// Future feature - MaxConcurrent configuration
config := scorer.Config{
    OpenAIKey:     "your-key",
    MaxConcurrent: 3, // Limit concurrent API calls
}
```

### Cost Optimization

- Uses `openai.GPT4oMini` for cost efficiency
- Batch processing reduces API call overhead
- JSON schema validation prevents invalid responses

## Example Applications

### Content Curation

```go
func filterRelevantItems(items []scorer.TextItem, minScore int) []scorer.TextItem {
    scored, err := scorer.ScoreTexts(ctx, items)
    if err != nil {
        return nil
    }
    
    var filtered []scorer.TextItem
    for _, si := range scored {
        if si.Score >= minScore {
            filtered = append(filtered, si.Item)
        }
    }
    return filtered
}
```

### Analytics Dashboard

```go
func generateScoreStats(scored []scorer.ScoredItem) map[string]int {
    stats := make(map[string]int)
    for _, si := range scored {
        switch {
        case si.Score >= 90:
            stats["high"]++
        case si.Score >= 70:
            stats["medium"]++
        case si.Score >= 40:
            stats["low"]++
        default:
            stats["irrelevant"]++
        }
    }
    return stats
}
```

### CSV Integration

See `examples/basic/main.go` for a complete example of:
- Loading text items from CSV files
- Including metadata and context in text items
- Environment variable configuration
- Structured logging integration
- Error handling and reporting