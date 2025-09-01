# Migration Guide: v0.10.0 to v0.11.0

This guide helps you migrate from v0.10.0 to v0.11.0, which renames the library from `post-scorer` to `llm-client`.

## Breaking Change: Module Rename

The module has been renamed to better reflect its broader LLM capabilities beyond just scoring posts.

### Import Path Change

**Before (v0.10.0):**

```go
import "github.com/JohnPlummer/post-scorer/scorer"
```

**After (v0.11.0):**

```go
import "github.com/JohnPlummer/llm-client/scorer"
```

## Migration Steps

### Step 1: Update Dependencies

Update your `go.mod` file to use the new module:

```bash
# Remove old dependency
go get github.com/JohnPlummer/post-scorer@none

# Add new dependency
go get github.com/JohnPlummer/llm-client@v0.11.0
```

### Step 2: Update Import Statements

Find and replace all import statements in your codebase:

```bash
# Find all Go files with the old import
grep -r "github.com/JohnPlummer/post-scorer" --include="*.go" .

# Update them to use the new import path
find . -name "*.go" -type f -exec sed -i '' 's|github.com/JohnPlummer/post-scorer|github.com/JohnPlummer/llm-client|g' {} \;
```

### Step 3: Verify and Test

After updating imports:

```bash
# Tidy dependencies
go mod tidy

# Run tests to ensure everything works
go test ./...

# Build your application
go build
```

## No API Changes

**Important:** This release only changes the module name. All APIs, functions, and features remain exactly the same as v0.10.0. You only need to update the import paths.

### What Stays the Same

- ✅ All function signatures
- ✅ Configuration options
- ✅ Resilience patterns (circuit breaker, retry)
- ✅ Generic text scoring APIs
- ✅ Prometheus metrics
- ✅ Health checks

## Example Migration

### Before (v0.10.0)

```go
package main

import (
    "context"
    "github.com/JohnPlummer/post-scorer/scorer"
)

func main() {
    s, err := scorer.BuildProductionScorer(apiKey)
    if err != nil {
        panic(err)
    }

    items := []scorer.TextItem{
        {ID: "1", Title: "Title", Body: "Content"},
    }

    results, err := s.ScoreTexts(context.Background(), items)
    // ... rest of code
}
```

### After (v0.11.0)

```go
package main

import (
    "context"
    "github.com/JohnPlummer/llm-client/scorer"  // Only this line changes!
)

func main() {
    s, err := scorer.BuildProductionScorer(apiKey)
    if err != nil {
        panic(err)
    }

    items := []scorer.TextItem{
        {ID: "1", Title: "Title", Body: "Content"},
    }

    results, err := s.ScoreTexts(context.Background(), items)
    // ... rest of code remains exactly the same
}
```

## GitHub Redirect

The GitHub repository has been renamed from `post-scorer` to `llm-client`. GitHub automatically redirects from the old URL to the new one, so existing links and clones will continue to work. However, you may want to update your git remotes:

```bash
# Check current remote
git remote -v

# If it shows the old URL, update it
git remote set-url origin git@github.com:JohnPlummer/llm-client.git
```

## Rollback Plan

If you need to rollback to v0.10.0:

```bash
# Rollback to previous version
go get github.com/JohnPlummer/post-scorer@v0.10.0

# Update imports back to old path
find . -name "*.go" -type f -exec sed -i '' 's|github.com/JohnPlummer/llm-client|github.com/JohnPlummer/post-scorer|g' {} \;
```

Note: The old module path will continue to work for existing versions due to GitHub's redirect.

## Support

For issues or questions about migration:

- Open an issue: <https://github.com/JohnPlummer/llm-client/issues>
- Check examples: `examples/` directory
- Review the README for updated documentation
