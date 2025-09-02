# LLM Client

A production-ready Go library for scoring text content using OpenAI's GPT models, with built-in resilience patterns, metrics, and generic text support.

## Overview

LLM-Client provides intelligent text scoring with enterprise-grade features:

- **Generic Text Scoring**: Score any text content, not just Reddit posts
- **Resilience Patterns**: Circuit breaker and retry with backoff
- **Production Monitoring**: Prometheus metrics integration
- **Template Support**: Dynamic prompt generation with Go templates
- **Batch Processing**: Efficient handling of multiple items

## Installation

```bash
go get github.com/JohnPlummer/llm-client@v0.11.0
```

## Quick Start

### Basic Usage (Generic Text)

```go
package main

import (
    "context"
    "github.com/JohnPlummer/llm-client/scorer"
)

func main() {
    // Create scorer with production configuration
    s, err := scorer.BuildProductionScorer("your-api-key")
    if err != nil {
        panic(err)
    }

    // Score text items
    items := []scorer.TextItem{
        {ID: "1", Content: "Best coffee shops in Seattle - Looking for recommendations...", Metadata: map[string]interface{}{"title": "Best coffee shops in Seattle"}},
        {ID: "2", Content: "Moving to Portland - What neighborhoods are family-friendly?", Metadata: map[string]interface{}{"title": "Moving to Portland"}},
    }

    results, err := s.ScoreTexts(context.Background(), items)
    if err != nil {
        panic(err)
    }

    // Use the results
    for _, result := range results {
        title := result.Item.Metadata["title"]
        fmt.Printf("Content: %s\nScore: %d\nReason: %s\n\n",
            title,
            result.Score,
            result.Reason)
    }
}
```

## Resilience Features

### Production Configuration

The library includes production-ready resilience patterns:

```go
// Create a production-ready scorer with all resilience features
scorer, err := scorer.BuildProductionScorer("api-key")

// Or customize configuration
cfg := scorer.NewDefaultConfig("api-key")
cfg = cfg.WithCircuitBreaker()  // Add circuit breaker
cfg = cfg.WithRetry()            // Add retry logic
cfg = cfg.WithMaxConcurrent(10)  // Set concurrency limit

scorer, err := scorer.NewIntegratedScorer(cfg)
```

### Circuit Breaker

Prevents cascade failures by stopping requests when error threshold is reached:

```go
cfg := scorer.Config{
    APIKey:               "api-key",
    EnableCircuitBreaker: true,
    CircuitBreakerConfig: &scorer.CircuitBreakerConfig{
        MaxRequests: 3,                    // Requests allowed in half-open state
        Interval:    10 * time.Second,      // Reset interval
        Timeout:     60 * time.Second,      // Time before trying half-open
        OnStateChange: func(name string, from, to gobreaker.State) {
            log.Printf("Circuit breaker %s: %v -> %v", name, from, to)
        },
    },
}
```

### Retry with Backoff

Automatically retries transient failures with configurable strategies:

```go
cfg := scorer.Config{
    APIKey:      "api-key",
    EnableRetry: true,
    RetryConfig: &scorer.RetryConfig{
        MaxAttempts:  3,
        Strategy:     scorer.RetryStrategyExponential,
        InitialDelay: 100 * time.Millisecond,
        MaxDelay:     5 * time.Second,
        Jitter:       0.1,  // 10% jitter to prevent thundering herd
    },
}
```

Available retry strategies:

- `RetryStrategyConstant`: Fixed delay between attempts
- `RetryStrategyExponential`: Exponentially increasing delay
- `RetryStrategyFibonacci`: Fibonacci sequence delays

### Prometheus Metrics

Built-in metrics for production monitoring:

```go
// Metrics are automatically recorded when using IntegratedScorer
// Expose metrics endpoint
http.Handle("/metrics", scorer.GetMetricsHandler())

// Available metrics:
// - text_scorer_requests_total
// - text_scorer_request_duration_seconds
// - text_scorer_errors_total
// - text_scorer_circuit_breaker_state
// - text_scorer_retry_attempts
// - text_scorer_score_distribution
```

## Configuration

### Core Configuration

- `APIKey` (required): Your OpenAI API key
- `Model` (optional): OpenAI model to use (defaults to GPT-4o-mini)
- `PromptText` (optional): Custom prompt template
- `MaxConcurrent` (optional): Concurrent batch processing limit
- `Timeout` (optional): Request timeout (default: 30s)

### Resilience Configuration

- `EnableCircuitBreaker`: Enable circuit breaker pattern
- `CircuitBreakerConfig`: Circuit breaker settings
- `EnableRetry`: Enable retry with backoff
- `RetryConfig`: Retry behavior settings

## Advanced Usage

### Per-Request Model Selection

Override the model for specific scoring requests:

```go
// Use GPT-4 for more accurate scoring
results, err := scorer.ScoreTextsWithOptions(ctx, items,
    scorer.WithModel("gpt-4"))

// Use GPT-3.5-turbo for faster, cheaper scoring
results, err := scorer.ScoreTextsWithOptions(ctx, items,
    scorer.WithModel("gpt-3.5-turbo"))
```

### Custom Prompt Templates

Use Go template syntax for dynamic prompts:

```go
// Template with extra context
template := "Score text items for {{.City}}: {{.Items}}"
results, err := scorer.ScoreTextsWithOptions(ctx, items,
    scorer.WithPromptTemplate(template),
    scorer.WithExtraContext(map[string]interface{}{"City": "Brighton"}))
```

### Scoring with Additional Context

Score text items with extra metadata:

```go
items := []scorer.TextItem{
    {
        ID: "1",
        Content: "Best coffee shops in Seattle - Looking for recommendations...",
        Metadata: map[string]interface{}{
            "title": "Best coffee shops in Seattle",
            "comments": "Great coffee! Been there many times.",
            "location": "Brighton",
        },
    },
}

// Use a template that includes the metadata
template := `Score this text item:
Title: {{range .Items}}{{.Metadata.title}}{{end}}
Content: {{range .Items}}{{.Content}}{{end}}
Comments: {{range .Items}}{{.Metadata.comments}}{{end}}`

results, err := scorer.ScoreTextsWithOptions(ctx, items,
    scorer.WithPromptTemplate(template),
    scorer.WithModel("gpt-4o"))
```

## Custom Prompts

Your prompt must instruct the LLM to return JSON in this exact format:

```json
{
  "version": "1.0",
  "scores": [
    {
      "item_id": "<id>",
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
4. Every text item must receive a score and reason
5. Include `%s` as placeholder for simple prompts, or use Go template syntax for advanced prompts with `.Items` field access

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

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
