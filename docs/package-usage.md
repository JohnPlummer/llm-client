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
    "github.com/JohnPlummer/reddit-client/reddit"
)

func main() {
    // Initialize the scorer
    s, err := scorer.New(scorer.Config{
        OpenAIKey: "your-openai-api-key",
    })
    if err != nil {
        panic(err)
    }

    // Prepare posts to score
    posts := []*reddit.Post{
        {
            ID:       "post1",
            Title:    "Best restaurants in downtown?",
            SelfText: "Looking for good dinner spots",
        },
        {
            ID:       "post2", 
            Title:    "Weekend events happening?",
            SelfText: "What's going on this weekend?",
        },
    }

    // Score the posts
    scoredPosts, err := s.ScorePosts(context.Background(), posts)
    if err != nil {
        panic(err)
    }

    // Use the results
    for _, post := range scoredPosts {
        fmt.Printf("Title: %s\nScore: %d\nReason: %s\n\n", 
            post.Post.Title, 
            post.Score, 
            post.Reason)
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
s, err := scorer.New(config)
```

## Input Data Structure

### Reddit Post Structure

The library expects posts conforming to the `reddit.Post` type:

```go
type Post struct {
    ID       string     // Required: Unique identifier
    Title    string     // Required: Post title
    SelfText string     // Optional: Post body text
    Comments []Comment  // Optional: Associated comments
    // ... other fields
}
```

### Comments Integration

Comments provide additional context for scoring:

```go
post := &reddit.Post{
    ID:       "post1",
    Title:    "Restaurant recommendations?",
    SelfText: "Looking for good food",
    Comments: []reddit.Comment{
        {Body: "Try Luigi's on Main Street!"},
        {Body: "The corner café has great coffee"},
    },
}
```

## Output Structure

### ScoredPost Type

```go
type ScoredPost struct {
    Post   *reddit.Post // Original post data
    Score  int          // AI-generated score (0-100)
    Reason string       // Explanation for the score
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

1. **Post placeholder**: Include `%s` where posts will be injected
2. **JSON output requirement**: Specify exact JSON structure needed
3. **Scoring criteria**: Define clear 0-100 scoring guidelines
4. **Complete coverage**: Require scoring of ALL input posts

### JSON Response Format

Your prompt must instruct the AI to return JSON in this exact structure:

```json
{
  "version": "1.0",
  "scores": [
    {
      "post_id": "<post_id>",
      "title": "<post_title>",
      "score": <0-100>,
      "reason": "<explanation>"
    }
  ]
}
```

### Example Custom Prompt

```text
Score Reddit posts for technology events and meetups (0-100 scale).

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

CRITICAL: Score every post, return only valid JSON.

Posts to score:
%s
```

### Validation Requirements

- **No markdown formatting** in JSON output
- **All fields required**: post_id, title, score, reason
- **Score range**: Must be 0-100 integer
- **Complete coverage**: Every input post must receive a score
- **Valid JSON**: Must parse without errors

## Error Handling

### Common Error Scenarios

```go
// Missing API key
s, err := scorer.New(scorer.Config{})
if err == scorer.ErrMissingAPIKey {
    // Handle missing API key
}

// API failures
scoredPosts, err := s.ScorePosts(ctx, posts)
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

- **Automatic batching**: Posts are processed in batches of 10
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
func filterRelevantPosts(posts []*reddit.Post, minScore int) []*reddit.Post {
    scored, err := scorer.ScorePosts(ctx, posts)
    if err != nil {
        return nil
    }
    
    var filtered []*reddit.Post
    for _, sp := range scored {
        if sp.Score >= minScore {
            filtered = append(filtered, sp.Post)
        }
    }
    return filtered
}
```

### Analytics Dashboard

```go
func generateScoreStats(scored []*ScoredPost) map[string]int {
    stats := make(map[string]int)
    for _, sp := range scored {
        switch {
        case sp.Score >= 90:
            stats["high"]++
        case sp.Score >= 70:
            stats["medium"]++
        case sp.Score >= 40:
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
- Loading posts from CSV files
- Associating comments with posts
- Environment variable configuration
- Structured logging integration
- Error handling and reporting