# Deployment Guide

## Library Distribution

### Go Module Publishing

This project is distributed as a Go module through Git-based versioning.

#### Version Tagging

```bash
# Create a new version tag
git tag v1.2.3
git push origin v1.2.3

# Verify tag creation
git tag -l
```

#### Semantic Versioning

Follow semantic versioning (SemVer) for releases:
- **Major** (v2.0.0): Breaking API changes
- **Minor** (v1.1.0): New features, backward compatible
- **Patch** (v1.0.1): Bug fixes, backward compatible

### Release Checklist

1. **Run full test suite**: `make check`
2. **Update documentation**: Ensure all docs reflect changes
3. **Verify examples**: Test `examples/basic/main.go`
4. **Generate coverage**: `make coverage` and review
5. **Update version references**: In README and docs
6. **Create git tag**: Follow semantic versioning
7. **Push to repository**: Include tags

## Integration Deployment

### Application Integration

#### Environment Variables

```bash
# Required
export OPENAI_API_KEY="sk-your-openai-api-key"

# Optional
export LOG_LEVEL="info"  # debug, info, warn, error
```

#### Docker Integration

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o scorer-app ./cmd/your-app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/scorer-app .

# Set environment variables
ENV OPENAI_API_KEY=""
ENV LOG_LEVEL="info"

CMD ["./scorer-app"]
```

#### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: llm-client-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: llm-client
  template:
    metadata:
      labels:
        app: llm-client
    spec:
      containers:
      - name: scorer
        image: your-registry/llm-client:v1.0.0
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: openai-secret
              key: api-key
        - name: LOG_LEVEL
          value: "info"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Secret
metadata:
  name: openai-secret
type: Opaque
data:
  api-key: <base64-encoded-api-key>
```

## Production Considerations

### API Key Management

#### Security Best Practices

- **Never commit API keys** to version control
- **Use environment variables** or secret management systems
- **Rotate keys regularly** as per security policy
- **Monitor API usage** for unusual patterns

#### Secret Management Options

```go
// Using HashiCorp Vault
import "github.com/hashicorp/vault/api"

func getAPIKeyFromVault() (string, error) {
    config := vault.DefaultConfig()
    client, err := vault.NewClient(config)
    if err != nil {
        return "", err
    }
    
    secret, err := client.Logical().Read("secret/openai")
    if err != nil {
        return "", err
    }
    
    return secret.Data["api_key"].(string), nil
}
```

### Rate Limiting and Quotas

#### OpenAI API Limits

- **Requests per minute**: Varies by subscription tier
- **Tokens per minute**: Monitor usage patterns
- **Concurrent requests**: Plan for `MaxConcurrent` configuration

#### Application-Level Rate Limiting

```go
import "golang.org/x/time/rate"

type RateLimitedScorer struct {
    scorer.Scorer
    limiter *rate.Limiter
}

func (r *RateLimitedScorer) ScoreTexts(ctx context.Context, items []scorer.TextItem) ([]scorer.ScoredItem, error) {
    if err := r.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    return r.Scorer.ScoreTexts(ctx, items)
}
```

### Monitoring and Observability

#### Structured Logging

```go
import "log/slog"

// Production logging configuration
opts := &slog.HandlerOptions{
    Level: slog.LevelInfo,
    AddSource: true,
}
logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
slog.SetDefault(logger)
```

#### Metrics Collection

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    itemsScored = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "items_scored_total",
            Help: "Total number of text items scored",
        },
        []string{"status"},
    )
    
    scoringDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "scoring_duration_seconds",
            Help: "Time spent scoring text items",
        },
        []string{"batch_size"},
    )
)

func init() {
    prometheus.MustRegister(itemsScored, scoringDuration)
}
```

#### Health Checks

```go
func (s *scorer) HealthCheck(ctx context.Context) error {
    // Test with minimal text item
    testItem := []scorer.TextItem{{
        ID:      "health-check",
        Content: "Health check text item",
        Metadata: map[string]interface{}{"title": "Health check"},
    }}
    
    _, err := s.ScoreTexts(ctx, testItem)
    return err
}
```

### Error Recovery and Resilience

#### Retry Logic

```go
import "github.com/cenkalti/backoff/v4"

func ScoreWithRetry(ctx context.Context, s scorer.Scorer, items []scorer.TextItem) ([]scorer.ScoredItem, error) {
    var result []scorer.ScoredItem
    
    operation := func() error {
        var err error
        result, err = s.ScoreTexts(ctx, items)
        return err
    }
    
    backoffStrategy := backoff.WithContext(
        backoff.NewExponentialBackOff(),
        ctx,
    )
    
    err := backoff.Retry(operation, backoffStrategy)
    return result, err
}
```

#### Circuit Breaker Pattern

```go
import "github.com/sony/gobreaker"

func NewCircuitBreakerScorer(s scorer.Scorer) scorer.Scorer {
    cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "openai-scorer",
        MaxRequests: 3,
        Interval:    time.Minute,
        Timeout:     30 * time.Second,
    })
    
    return &CircuitBreakerScorer{
        scorer: s,
        cb:     cb,
    }
}
```

## Performance Optimization

### Batch Size Tuning

The library uses a fixed batch size of 10 text items. Monitor performance metrics to determine if this needs adjustment:

```go
// Future configuration option
config := scorer.Config{
    OpenAIKey: key,
    BatchSize: 15, // Adjust based on performance testing
}
```

### Concurrent Processing

Plan for future concurrent batch processing:

```go
// Proposed enhancement
config := scorer.Config{
    OpenAIKey:     key,
    MaxConcurrent: 3, // Process 3 batches simultaneously
}
```

### Caching Strategy

For repeated scoring of similar content:

```go
import "github.com/allegro/bigcache"

type CachedScorer struct {
    scorer.Scorer
    cache *bigcache.BigCache
}

func (c *CachedScorer) ScoreTexts(ctx context.Context, items []scorer.TextItem) ([]scorer.ScoredItem, error) {
    // Check cache for previously scored text items
    // Fall back to API for cache misses
}
```

## Security Considerations

### API Key Protection

- **Environment isolation**: Separate keys for dev/staging/prod
- **Key rotation**: Regular rotation schedule
- **Access logging**: Monitor API key usage
- **Minimal permissions**: Use least-privilege principle

### Input Validation

```go
func ValidateTextItems(items []scorer.TextItem) error {
    for _, item := range items {
        if item.ID == "" {
            return errors.New("text item ID cannot be empty")
        }
        if len(item.Content) > 10000 {
            return errors.New("text item content too long")
        }
    }
    return nil
}
```

### Output Sanitization

Ensure scored content doesn't contain sensitive information:

```go
func SanitizeScoredItems(items []scorer.ScoredItem) {
    for _, item := range items {
        // Remove or mask sensitive data
        item.Reason = sanitizeReason(item.Reason)
    }
}
```

## Troubleshooting Deployment Issues

### Common Problems

1. **API Key Issues**
   - Verify key format and validity
   - Check environment variable loading
   - Confirm API quota availability

2. **Network Connectivity**
   - Test OpenAI API accessibility
   - Check firewall and proxy settings
   - Verify TLS certificate validation

3. **Memory Usage**
   - Monitor batch processing memory
   - Consider batch size adjustment
   - Implement garbage collection tuning

4. **Response Timeouts**
   - Adjust context timeout values
   - Monitor API response times
   - Implement retry strategies

### Diagnostic Commands

```bash
# Test API connectivity
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
     https://api.openai.com/v1/models

# Check application health
kubectl exec -it pod-name -- ./scorer-app --health-check

# Monitor resource usage
kubectl top pods -l app=llm-client
```