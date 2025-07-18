# Post Scorer

A Go package that uses OpenAI's GPT to score Reddit posts based on some arbitrary criteria specified in a custom prompt.

## Overview

The scorer evaluates Reddit posts and returns a slice of `ScoredPost` structs containing:

- The original post
- A relevance score (0-100)
- A reason for the score

## Installation

**Note**: This library currently depends on `github.com/JohnPlummer/reddit-client` which must be available as a published module. Currently, this dependency uses a local replace directive in `go.mod`.

To use this library externally, you will need to either:
1. Publish the `reddit-client` module to a public repository
2. Use a local clone of both repositories
3. Replace the dependency with an alternative Reddit client

```bash
go get github.com/JohnPlummer/post-scorer
```

## Usage

```go
package main

import (
    "context"
    "github.com/JohnPlummer/post-scorer/scorer"
    "github.com/JohnPlummer/reddit-client/reddit"
)

func main() {
    // Initialize the scorer
    s, err := scorer.New(scorer.Config{
        OpenAIKey: "your-api-key",
    })
    if err != nil {
        panic(err)
    }

    // Score some posts
    posts := []*reddit.Post{
        {
            ID:    "post1",
            Title: "Best restaurants in town?",
        },
    }

    scored, err := s.ScorePosts(context.Background(), posts)
    if err != nil {
        panic(err)
    }

    // Use the scored posts
    for _, post := range scored {
        fmt.Printf("Post: %s\nScore: %d\nReason: %s\n\n", 
            post.Post.Title, 
            post.Score, 
            post.Reason)
    }
}
```

## Configuration

The `Config` struct accepts:

- `OpenAIKey` (required): Your OpenAI API key
- `PromptText` (optional): Custom prompt template
- `MaxConcurrent` (optional): For rate limiting

## Custom Prompts

Your prompt must instruct the LLM to return JSON in this exact format:

```json
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
```

Critical requirements:

1. Output must be ONLY valid JSON (no markdown or other formatting)
2. All fields are required
3. Score must be between 0-100
4. Every post must receive a score and reason
5. Include `%s` as placeholder where posts will be injected

See `examples/basic/custom_prompt.txt` for a complete example prompt.

## Documentation

For comprehensive documentation, see the `docs/` directory:

- **[Project Overview](docs/project-overview.md)** - Architecture, features, and use cases
- **[Development Setup](docs/development-setup.md)** - Installation, dependencies, and coding standards  
- **[Package Usage](docs/package-usage.md)** - Complete API reference and examples
- **[Key Components](docs/key-components.md)** - Core interfaces and implementation details
- **[Deployment Guide](docs/deployment-guide.md)** - Production deployment and configuration
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions
- **[Recent Changes](docs/recent-changes.md)** - Latest updates and improvements

## Examples

Check the `examples` directory for complete usage examples, including CSV data loading and custom prompt configuration.
