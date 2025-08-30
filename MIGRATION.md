# Migration Guide: v0.9.x to v0.10.0

This guide helps you migrate from post-scorer v0.9.x to v0.10.0, which introduces resilience patterns and generic text scoring support.

## Breaking Changes

### 1. New Generic Text Scoring Interface

The library now supports generic text scoring, not just Reddit posts:

**Before (v0.9.x):**
```go
// Only Reddit posts supported
scorer.ScorePosts(ctx, []*reddit.Post{...})
```

**After (v0.10.0):**
```go
// Generic text scoring (recommended)
scorer.ScoreTexts(ctx, []scorer.TextItem{
    {ID: "1", Title: "Title", Body: "Content"},
})

// Reddit posts still supported for backward compatibility
scorer.ScorePosts(ctx, []*reddit.Post{...})
```

### 2. Configuration Changes

**Before (v0.9.x):**
```go
cfg := scorer.Config{
    OpenAIKey: "key",
    Model:     "gpt-4",
}
s, _ := scorer.New(cfg)
```

**After (v0.10.0):**
```go
// Option 1: Production configuration (recommended)
s, _ := scorer.BuildProductionScorer("key")

// Option 2: Custom configuration with resilience
cfg := scorer.NewDefaultConfig("key")
cfg = cfg.WithCircuitBreaker()
cfg = cfg.WithRetry()
s, _ := scorer.NewIntegratedScorer(cfg)

// Option 3: Legacy configuration (still works)
s, _ := scorer.New(scorer.Config{OpenAIKey: "key"})
```

## New Features

### Circuit Breaker Pattern

Protect against cascading failures:

```go
cfg := scorer.Config{
    APIKey:               "key",
    EnableCircuitBreaker: true,
    CircuitBreakerConfig: &scorer.CircuitBreakerConfig{
        MaxRequests: 3,
        Interval:    10 * time.Second,
        Timeout:     60 * time.Second,
    },
}
```

### Retry with Backoff

Automatic retry for transient failures:

```go
cfg := scorer.Config{
    APIKey:      "key",
    EnableRetry: true,
    RetryConfig: &scorer.RetryConfig{
        MaxAttempts:  3,
        Strategy:     scorer.RetryStrategyExponential,
        InitialDelay: 100 * time.Millisecond,
        MaxDelay:     5 * time.Second,
    },
}
```

### Prometheus Metrics

Built-in production monitoring:

```go
// Metrics automatically recorded with IntegratedScorer
http.Handle("/metrics", scorer.GetMetricsHandler())
```

### Health Checks

Monitor scorer health:

```go
health := scorer.GetHealth(ctx)
if health.Status != scorer.HealthStatusHealthy {
    log.Printf("Scorer unhealthy: %v", health.Details)
}
```

## Migration Steps

### Step 1: Update Import

```bash
go get github.com/JohnPlummer/post-scorer@v0.10.0
```

### Step 2: Choose Migration Path

#### Option A: Minimal Changes (Keep existing code working)

No code changes required. Your existing code will continue to work with the legacy interfaces.

#### Option B: Adopt Resilience Features (Recommended)

Replace your scorer initialization:

```go
// Old
s, err := scorer.New(scorer.Config{
    OpenAIKey: apiKey,
})

// New - with production resilience
s, err := scorer.BuildProductionScorer(apiKey)
```

#### Option C: Full Migration to Generic Text

Convert Reddit posts to generic text items:

```go
// Convert Reddit posts to TextItems
items := make([]scorer.TextItem, len(posts))
for i, post := range posts {
    items[i] = scorer.TextItem{
        ID:    post.ID,
        Title: post.Title,
        Body:  post.Selftext,
    }
}

// Use generic scoring
results, err := scorer.ScoreTexts(ctx, items)
```

### Step 3: Add Monitoring (Optional)

```go
// Expose metrics endpoint
http.Handle("/metrics", scorer.GetMetricsHandler())

// Add health check endpoint
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    health := scorer.GetHealth(r.Context())
    if health.Status == scorer.HealthStatusHealthy {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
    json.NewEncoder(w).Encode(health)
})
```

## Example: Complete Migration

### Before (v0.9.x)

```go
package main

import (
    "context"
    "github.com/JohnPlummer/post-scorer/scorer"
    "github.com/JohnPlummer/reddit-client/reddit"
)

func main() {
    cfg := scorer.Config{
        OpenAIKey:     getAPIKey(),
        Model:         "gpt-4",
        MaxConcurrent: 5,
    }
    
    s, err := scorer.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    posts := loadRedditPosts()
    scored, err := s.ScorePosts(context.Background(), posts)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, post := range scored {
        fmt.Printf("Score: %d\n", post.Score)
    }
}
```

### After (v0.10.0)

```go
package main

import (
    "context"
    "net/http"
    "github.com/JohnPlummer/post-scorer/scorer"
)

func main() {
    // Create production-ready scorer with all resilience features
    s, err := scorer.BuildProductionScorer(getAPIKey())
    if err != nil {
        log.Fatal(err)
    }
    
    // Convert to generic text items
    posts := loadRedditPosts()
    items := make([]scorer.TextItem, len(posts))
    for i, post := range posts {
        items[i] = scorer.TextItem{
            ID:    post.ID,
            Title: post.Title,
            Body:  post.Selftext,
        }
    }
    
    // Score with resilience patterns active
    results, err := s.ScoreTexts(context.Background(), items)
    if err != nil {
        log.Printf("Scoring failed: %v", err)
        // Check if circuit breaker is open
        health := s.GetHealth(context.Background())
        if health.Status == scorer.HealthStatusUnhealthy {
            log.Printf("Circuit breaker open: %v", health.Details)
        }
    }
    
    for _, result := range results {
        fmt.Printf("Score: %d\n", result.Score)
    }
    
    // Expose metrics for monitoring
    http.Handle("/metrics", scorer.GetMetricsHandler())
    go http.ListenAndServe(":8080", nil)
}
```

## Rollback Plan

If you encounter issues with v0.10.0, you can safely rollback:

```bash
go get github.com/JohnPlummer/post-scorer@v0.9.1
```

The v0.9.x API remains fully supported in v0.10.0, so gradual migration is possible.

## Support

For issues or questions about migration:
- Open an issue: https://github.com/JohnPlummer/post-scorer/issues
- Check examples: `examples/` directory
- Review tests: `scorer/*_test.go` for usage patterns