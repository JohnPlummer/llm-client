# Post-Scorer Basic Example

This example demonstrates the new v0.10.0 API features including resilience patterns and generic text scoring.

## Features Demonstrated

1. **Production Configuration** - Ready-to-use setup with all resilience features
2. **Custom Configuration** - Fine-tuned resilience settings
3. **Generic Text Scoring** - New TextItem/ScoredItem API
4. **Health Monitoring** - Service health checks
5. **Prometheus Metrics** - Production monitoring endpoints

## Running the Example

### Prerequisites

1. Go 1.23.1 or higher
2. OpenAI API key

### Setup

```bash
# Install dependencies
go mod download

# Set your OpenAI API key
export OPENAI_API_KEY="your-api-key"

# Or create a .env file
echo "OPENAI_API_KEY=your-api-key" > .env
```

### Run

```bash
go run main.go
```

The example will:
1. Score sample texts using production configuration
2. Demonstrate custom configuration with specific settings
3. Show health check status
4. Start a metrics server on http://localhost:8080

## What's New in v0.10.0

### Production-Ready Features

```go
// One-line production setup with all resilience patterns
scorer, err := scorer.BuildProductionScorer(apiKey)
```

This includes:
- **Circuit Breaker**: Prevents cascade failures
- **Retry Logic**: Handles transient errors with exponential backoff
- **Prometheus Metrics**: Full observability
- **Health Checks**: Service status monitoring

### Generic Text Scoring

```go
// New generic API - not limited to Reddit posts
items := []scorer.TextItem{
    {ID: "1", Title: "Title", Body: "Content"},
}
results, err := scorer.ScoreTexts(ctx, items)
```

### Custom Configuration

```go
cfg := scorer.NewDefaultConfig(apiKey)
cfg = cfg.WithCircuitBreaker()
cfg = cfg.WithRetry()
cfg = cfg.WithMaxConcurrent(5)

// Fine-tune retry strategy
cfg.RetryConfig.Strategy = scorer.RetryStrategyExponential
cfg.RetryConfig.MaxAttempts = 3
```

## Monitoring

While the example is running, you can access:

- **Metrics**: http://localhost:8080/metrics (Prometheus format)
- **Health**: http://localhost:8080/health (JSON status)

### Key Metrics

- `text_scorer_requests_total` - Total requests by status
- `text_scorer_request_duration_seconds` - Request latency
- `text_scorer_circuit_breaker_state` - Circuit breaker status (0=closed, 1=half-open, 2=open)
- `text_scorer_retry_attempts` - Retry attempt distribution
- `text_scorer_score_distribution` - Score value histogram

## Files

- `main.go` - Example implementation
- `example_posts.csv` - Sample data
- `example_comments.csv` - Additional context data
- `custom_prompt.txt` - Custom scoring prompt

## Migration from v0.9.x

See [MIGRATION.md](../../MIGRATION.md) for upgrading existing code.