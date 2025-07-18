# Troubleshooting

## Common Issues

### API Key Problems

#### Missing API Key Error
```
Error: OpenAI API key is required
```

**Solutions:**
1. **Set environment variable:**
   ```bash
   export OPENAI_API_KEY="sk-your-actual-api-key"
   ```

2. **Check .env file loading:**
   ```go
   if err := godotenv.Load(); err != nil {
       fmt.Println("Error loading .env file:", err)
   }
   ```

3. **Verify key format:**
   - Should start with `sk-`
   - Should be exactly 51 characters long
   - No trailing spaces or newlines

#### Invalid API Key Error
```
Error: OpenAI API error: Incorrect API key provided
```

**Solutions:**
1. **Verify key validity:**
   ```bash
   curl -H "Authorization: Bearer $OPENAI_API_KEY" \
        https://api.openai.com/v1/models
   ```

2. **Check key permissions:**
   - Ensure key has access to chat completions
   - Verify organization access if applicable

3. **Regenerate key:**
   - Create new key in OpenAI dashboard
   - Update environment variables

### Scoring Issues

#### No Scores Returned
```
Warning: incomplete scores from OpenAI, expected_count=5, received_count=0
```

**Causes and Solutions:**

1. **Prompt formatting issues:**
   ```go
   // Ensure custom prompt includes %s placeholder
   prompt := "Score these posts: %s"
   ```

2. **Invalid JSON response:**
   - Check OpenAI response format requirements
   - Verify JSON schema compliance
   - Use default prompt to test

3. **Model limitations:**
   - Posts may be too long for context window
   - Try reducing batch size manually

#### Invalid Score Ranges
```
Error: invalid score 150 for post abc123: score must be between 0 and 100
```

**Solutions:**
1. **Update custom prompt:**
   ```text
   IMPORTANT: Scores must be integers between 0 and 100.
   ```

2. **Add validation example:**
   ```text
   Example valid response:
   {"version": "1.0", "scores": [{"post_id": "123", "title": "Title", "score": 85, "reason": "Explanation"}]}
   ```

#### Missing Scores for Some Posts
```
Warning: missing score from OpenAI response, assigning default score, post_id=xyz789
```

**Expected Behavior:** This is normal graceful degradation. The library automatically assigns score 0 with explanation.

**To reduce frequency:**
1. **Improve prompt clarity:**
   ```text
   CRITICAL: You MUST score every single post. Never skip any posts.
   ```

2. **Use smaller batches:**
   ```go
   // Process posts in smaller groups
   const customBatchSize = 5
   ```

### Network and Connectivity Issues

#### Request Timeout Errors
```
Error: context deadline exceeded
```

**Solutions:**
1. **Increase context timeout:**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
   defer cancel()
   ```

2. **Check network connectivity:**
   ```bash
   ping api.openai.com
   ```

3. **Verify proxy settings:**
   ```bash
   export HTTPS_PROXY="http://proxy.company.com:8080"
   ```

#### Rate Limiting Errors
```
Error: OpenAI API error: Rate limit reached
```

**Solutions:**
1. **Implement backoff strategy:**
   ```go
   import "time"
   
   func retryWithBackoff(operation func() error) error {
       for i := 0; i < 3; i++ {
           if err := operation(); err != nil {
               if strings.Contains(err.Error(), "rate limit") {
                   time.Sleep(time.Duration(i+1) * time.Second)
                   continue
               }
               return err
           }
           return nil
       }
       return errors.New("max retries exceeded")
   }
   ```

2. **Reduce concurrent requests:**
   ```go
   // Future: Use MaxConcurrent configuration
   config.MaxConcurrent = 1
   ```

3. **Check OpenAI usage dashboard** for quota limits

### Data Processing Issues

#### JSON Parsing Errors
```
Error: unmarshaling response: invalid character 'T' looking for beginning of value
```

**Causes:**
- OpenAI returned non-JSON response
- Response wrapped in markdown code blocks
- Partial or truncated response

**Solutions:**
1. **Check response format in logs:**
   ```go
   slog.Debug("Raw OpenAI response", "content", resp.Choices[0].Message.Content)
   ```

2. **Update prompt to enforce JSON:**
   ```text
   Return ONLY valid JSON. No markdown, no explanations, no code blocks.
   ```

3. **Implement response cleaning:**
   ```go
   func cleanJSONResponse(content string) string {
       // Remove markdown code blocks
       content = strings.TrimPrefix(content, "```json")
       content = strings.TrimSuffix(content, "```")
       return strings.TrimSpace(content)
   }
   ```

#### CSV Loading Errors
```
Error: reading CSV records: record on line 3: wrong number of fields
```

**Solutions:**
1. **Check CSV format:**
   ```csv
   id,title,text
   post1,"Restaurant recommendation","Looking for good food"
   post2,"Events this weekend","What's happening?"
   ```

2. **Handle malformed CSV:**
   ```go
   reader := csv.NewReader(file)
   reader.LazyQuotes = true
   reader.TrimLeadingSpace = true
   ```

3. **Validate CSV headers:**
   ```go
   header, err := reader.Read()
   if err != nil || len(header) < 3 {
       return fmt.Errorf("invalid CSV format")
   }
   ```

### Testing Issues

#### Test Failures
```
Error: failed to create scorer: OpenAI API key is required
```

**Solutions:**
1. **Use test configuration:**
   ```go
   // In tests, use mock client
   scorer := scorer.NewWithClient(mockClient)
   ```

2. **Set test environment:**
   ```bash
   export OPENAI_API_KEY="test-key-for-mocking"
   ```

3. **Check Ginkgo setup:**
   ```bash
   ginkgo version
   go mod tidy
   ```

#### Coverage Issues
```
Error: no test files found
```

**Solutions:**
1. **Run from correct directory:**
   ```bash
   cd scorer/
   ginkgo -v ./...
   ```

2. **Check test file naming:**
   - Must end with `_test.go`
   - Must be in same package

### Performance Issues

#### Slow Processing
**Symptoms:** Long delays between batches

**Diagnosis:**
1. **Enable debug logging:**
   ```go
   opts := &slog.HandlerOptions{Level: slog.LevelDebug}
   logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
   slog.SetDefault(logger)
   ```

2. **Monitor batch timing:**
   ```bash
   # Look for batch processing logs
   grep "Processing batch" application.log
   ```

**Solutions:**
1. **Reduce batch size** if posts are very long
2. **Check network latency** to OpenAI API
3. **Monitor API response times**

#### Memory Usage
**Symptoms:** High memory consumption with large datasets

**Solutions:**
1. **Process in smaller chunks:**
   ```go
   const maxPosts = 100
   for i := 0; i < len(allPosts); i += maxPosts {
       chunk := allPosts[i:min(i+maxPosts, len(allPosts))]
       results, err := scorer.ScorePosts(ctx, chunk)
       // Process results immediately
   }
   ```

2. **Monitor garbage collection:**
   ```bash
   GODEBUG=gctrace=1 ./your-app
   ```

## Debugging Strategies

### Enable Detailed Logging

```go
// Set debug level for maximum visibility
opts := &slog.HandlerOptions{
    Level:     slog.LevelDebug,
    AddSource: true,
}
logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
slog.SetDefault(logger)
```

### Test with Minimal Example

```go
func debugScoring() {
    s, err := scorer.New(scorer.Config{
        OpenAIKey: os.Getenv("OPENAI_API_KEY"),
    })
    if err != nil {
        panic(err)
    }

    // Test with single simple post
    posts := []*reddit.Post{{
        ID:    "debug1",
        Title: "Test restaurant post",
    }}

    results, err := s.ScorePosts(context.Background(), posts)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    for _, result := range results {
        fmt.Printf("Score: %d, Reason: %s\n", result.Score, result.Reason)
    }
}
```

### Validate Configuration

```go
func validateSetup() {
    // Check environment
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        fmt.Println("❌ OPENAI_API_KEY not set")
        return
    }
    fmt.Printf("✅ API key present (length: %d)\n", len(apiKey))

    // Test API connectivity
    client := openai.NewClient(apiKey)
    _, err := client.ListModels(context.Background())
    if err != nil {
        fmt.Printf("❌ API test failed: %v\n", err)
        return
    }
    fmt.Println("✅ OpenAI API connectivity confirmed")

    // Test scorer creation
    _, err = scorer.New(scorer.Config{OpenAIKey: apiKey})
    if err != nil {
        fmt.Printf("❌ Scorer creation failed: %v\n", err)
        return
    }
    fmt.Println("✅ Scorer created successfully")
}
```

## Getting Help

### Log Collection

When reporting issues, include:

1. **Full error message** with stack trace
2. **Configuration** (without API key)
3. **Debug logs** from the problematic operation
4. **Go version** and dependency versions
5. **Sample input data** that causes the issue

### Diagnostic Commands

```bash
# Check Go version
go version

# Check dependencies
go mod graph | grep -E "(openai|reddit-client)"

# Run with verbose output
go run -v ./examples/basic/

# Test API independently
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"Hello"}]}' \
     https://api.openai.com/v1/chat/completions
```

### Environment Information

```bash
# System information
echo "OS: $(uname -a)"
echo "Go: $(go version)"
echo "Git: $(git --version)"

# Project information
echo "Current directory: $(pwd)"
echo "Git branch: $(git branch --show-current)"
echo "Git commit: $(git rev-parse HEAD)"
```