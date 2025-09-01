# CLAUDE.md

This file provides comprehensive guidance to Claude Code (claude.ai/code) when working with the llm-client Go library.

## Repository Overview

**Post-Scorer** is a production-ready Go library (v0.9.1) that scores Reddit posts using OpenAI's GPT API to determine their relevance for location-based recommendations and events. The library provides batch processing, concurrent execution, and flexible prompt customization.

### Project Statistics
- **Primary language**: Go 1.23.1+
- **Lines of code**: ~2,500
- **Active since**: February 9, 2025
- **Test coverage**: 86.5%
- **Total commits**: 52
- **Latest version**: 0.9.1

### Core Value Proposition
- Batch processing of Reddit posts (10 posts per API call for efficiency)
- Concurrent processing with configurable parallelism
- Interface-first design for testing and extensibility
- Custom prompt templates with Go template support
- Graceful error handling with detailed logging

## Quick Start

### Prerequisites
- Go 1.23.1 or higher
- OpenAI API key
- Git for version control

### Initial Setup
```bash
# Install the library
go get github.com/JohnPlummer/llm-client

# Set environment variable
export OPENAI_API_KEY="your-api-key"

# For development, clone the repository
git clone https://github.com/JohnPlummer/llm-client.git
cd llm-client

# Install dependencies
go mod download

# Verify installation
make test
```

### Basic Usage
```go
import "github.com/JohnPlummer/llm-client/scorer"

// Create scorer with config
cfg := scorer.Config{
    OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
    Model:         openai.GPT4oMini, // optional, defaults to GPT-4o-mini
    MaxConcurrent: 5,                 // optional, defaults to 1
}
s, err := scorer.New(cfg)

// Score posts
scoredPosts, err := s.ScorePosts(ctx, posts)
```

## Essential Commands

### Development
```bash
# Run tests with Ginkgo BDD framework
make test

# Run the basic example
make run

# Comprehensive validation (tidy + test + run)
make check

# Generate coverage report in markdown
make coverage
```

### Module Management
```bash
# Clean up root dependencies
make tidy

# Clean up example dependencies
make tidy-examples

# Clean up all dependencies
make tidy-all
```

### Testing Commands
```bash
# Run tests with verbose output
ginkgo -v ./...

# Run tests with race detection
go test -race ./...

# Generate coverage with profile
ginkgo -v --coverprofile=coverage.out ./...
```

## Architecture and Key Concepts

### System Architecture
The library follows clean architecture principles with interface-first design:

```
Application Layer (examples/)
    ↓
Domain Layer (scorer.Scorer interface)
    ↓
Implementation Layer (scorer.scorer struct)
    ↓
Infrastructure Layer (OpenAI API client)
```

### 1. **Scorer Interface**
The core abstraction that enables dependency injection and testing:
```go
type Scorer interface {
    ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error)
    ScorePostsWithOptions(ctx context.Context, posts []*reddit.Post, opts ...ScoringOption) ([]*ScoredPost, error)
    ScorePostsWithContext(ctx context.Context, contexts []ScoringContext, opts ...ScoringOption) ([]*ScoredPost, error)
}
```
- **Location**: `scorer/types.go`
- **Key files**: `scorer/scorer.go` (implementation)
- **Testing**: Mock implementation in `scorer/scorer_test.go`

### 2. **Batch Processing Architecture**
Fixed batch size of 10 posts per OpenAI API call for efficiency:
```go
const maxBatchSize = 10 // Hard limit for optimal API performance
```
- **Location**: `scorer/batch.go`
- **Rationale**: Token limits and cost optimization
- **Processing**: Sequential by default, concurrent with `MaxConcurrent > 1`

### 3. **Functional Options Pattern**
Per-request customization without breaking changes:
```go
scoredPosts, err := scorer.ScorePostsWithOptions(ctx, posts,
    scorer.WithModel(openai.GPT4),
    scorer.WithPromptTemplate("Custom: {{.PostTitle}}"),
    scorer.WithExtraContext(map[string]string{"location": "NYC"}),
)
```
- **Location**: `scorer/types.go`
- **Added**: v0.9.1 (July 29, 2025)

### Data Flow
1. **Input Validation** → Check for nil posts, empty IDs
2. **Batching** → Split into 10-post batches
3. **Prompt Generation** → Apply template with context
4. **API Call** → Send to OpenAI with JSON schema
5. **Response Parsing** → Validate scores (0-100 range)
6. **Result Assembly** → Map scores back to original posts

### Concurrency Model
- **Sequential Mode** (`MaxConcurrent <= 1`): Default for backward compatibility
- **Concurrent Mode** (`MaxConcurrent > 1`): Semaphore-based goroutine limiting
- **Order Preservation**: Indexed channels maintain input sequence

## Project Structure

```
llm-client/
├── scorer/                      # Core library package
│   ├── scorer.go                # Main implementation and constructors
│   ├── batch.go                 # Batch processing logic
│   ├── types.go                 # Interfaces and data structures
│   ├── prompts.go               # Prompt loading utilities
│   ├── prompts/                 # Embedded prompt templates
│   │   ├── system_prompt.txt   # Default system instructions
│   │   └── batch_prompt.txt    # Batch scoring template
│   └── scorer_test.go           # Ginkgo BDD test suite (697 lines)
├── examples/                    # Usage examples
│   └── basic/                   # Self-contained example module
│       ├── main.go              # Example implementation
│       ├── custom_prompt.txt   # Custom prompt template
│       ├── example_posts.csv   # Sample post data
│       ├── example_comments.csv # Sample comment data
│       ├── go.mod               # Separate module definition
│       └── go.sum               # Dependency checksums
├── docs/                        # Comprehensive documentation
│   ├── project-overview.md     # Architecture and features
│   ├── development-setup.md    # Installation guide
│   ├── key-components.md       # Core components reference
│   ├── package-usage.md        # API documentation
│   ├── deployment-guide.md     # Production deployment
│   ├── troubleshooting.md      # Common issues
│   └── recent-changes.md       # Version history
├── Makefile                     # Build automation (8 targets)
├── go.mod                       # Module dependencies
├── version.go                   # Version information (0.9.1)
├── README.md                    # Primary documentation
├── LICENSE                      # MIT license
└── IMPROVEMENTS.md              # Completed improvements (12 items)
```

## Important Patterns

### Interface-First Design
Always define interfaces before implementations:
```go
// Define interface
type OpenAIClient interface {
    CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

// Accept interface in constructors
func NewWithClient(client OpenAIClient) Scorer
```

### Error Handling Philosophy
**Context-rich error messages with wrapping:**
```go
// Always wrap errors with context
fmt.Errorf("processing batch %d of %d: %w", batchNum, totalBatches, err)

// Validate at boundaries
if post.ID == "" {
    return nil, fmt.Errorf("post at index %d has empty ID", i)
}
```

**Graceful degradation:**
- Missing scores default to 0 with warning logs
- Continue processing other posts on individual failures
- Log at appropriate levels (Debug, Info, Warn, Error)

### Adding New Features
1. Define interface changes in `types.go`
2. Implement in appropriate file (`scorer.go`, `batch.go`)
3. Add comprehensive tests in `scorer_test.go` using Ginkgo
4. Update examples in `examples/basic/`
5. Document in `docs/` if significant
6. Run `make check` before committing

### Dependency Injection Pattern
```go
// Production: use real client
cfg := scorer.Config{OpenAIKey: apiKey}
s, _ := scorer.New(cfg)

// Testing: use mock client
mockClient := &mockOpenAIClient{response: testResponse}
s := scorer.NewWithClient(mockClient)
```

### Embedded Resources
Prompts are embedded at compile time for zero runtime dependencies:
```go
//go:embed prompts/*.txt
var promptsFS embed.FS
```

## Code Style

### Naming Conventions
- **Files**: snake_case (`scorer_test.go`, `batch.go`)
- **Directories**: lowercase or kebab-case (`scorer/`, `llm-client/`)
- **Exported types/functions**: PascalCase (`ScorePosts`, `ScoredPost`)
- **Private types/functions**: camelCase (`processBatch`, `scoreResponse`)
- **Interfaces**: Descriptive nouns (`Scorer`, `OpenAIClient`)
- **Error variables**: `Err` prefix (`ErrMissingAPIKey`)
- **Constants**: camelCase or PascalCase (`maxBatchSize`)

### Import Organization
```go
import (
    // Standard library
    "context"
    "errors"
    "fmt"
    
    // Third-party packages
    "github.com/sashabaranov/go-openai"
    "github.com/JohnPlummer/reddit-client/reddit"
    
    // Internal packages (if any)
)
```

### File Organization
- `types.go`: All interfaces and data structures
- `scorer.go`: Constructors and main implementation
- `batch.go`: Batch processing logic
- `prompts.go`: Prompt management utilities
- Test files use `_test` suffix with external test package

### Documentation Standards
```go
// ScorePosts scores a slice of Reddit posts and returns scored posts with explanations.
// It processes posts in batches of 10 for API efficiency.
func (s *scorer) ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error) {
```

## Testing Approach

### Framework
**Ginkgo BDD with Gomega matchers** - Structured, readable tests:
```go
var _ = Describe("Scorer", func() {
    Context("when scoring posts", func() {
        It("should handle empty input gracefully", func() {
            result, err := scorer.ScorePosts(ctx, []*reddit.Post{})
            Expect(err).ToNot(HaveOccurred())
            Expect(result).To(BeEmpty())
        })
    })
})
```

### Test Coverage Requirements
- Current: 86.5% (must maintain or improve)
- All new features must include tests
- Edge cases and error conditions required
- Race condition testing with `-race` flag

### Mock Strategy
```go
type mockOpenAIClient struct {
    response openai.ChatCompletionResponse
    err      error
    createChatCompletionFunc func(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}
```

### Running Tests
```bash
# Quick test
make test

# With coverage
make coverage

# Specific package
ginkgo -v ./scorer

# With race detection
go test -race ./...
```

## API Reference

### Core Types

#### Config
```go
type Config struct {
    OpenAIKey     string // Required: OpenAI API key
    Model         string // Optional: defaults to GPT-4o-mini
    PromptText    string // Optional: custom prompt text
    MaxConcurrent int    // Optional: concurrent batch limit (default: 1)
}
```

#### ScoredPost
```go
type ScoredPost struct {
    Post   *reddit.Post // Original Reddit post
    Score  int          // Score 0-100
    Reason string       // AI explanation for score
}
```

#### ScoringContext
```go
type ScoringContext struct {
    Post      *reddit.Post       // Reddit post to score
    ExtraData map[string]string  // Additional context (e.g., comments)
}
```

### Constructor Functions

```go
// Standard constructor with config
func New(cfg Config) (Scorer, error)

// Constructor with custom client (for testing)
func NewWithClient(client OpenAIClient) Scorer

// Constructor with client and options
func NewWithClientAndOptions(client OpenAIClient, opts ...func(*scorer)) Scorer
```

### Scoring Methods

```go
// Basic scoring
ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error)

// Scoring with runtime options
ScorePostsWithOptions(ctx context.Context, posts []*reddit.Post, opts ...ScoringOption) ([]*ScoredPost, error)

// Scoring with extra context
ScorePostsWithContext(ctx context.Context, contexts []ScoringContext, opts ...ScoringOption) ([]*ScoredPost, error)
```

### Scoring Options

```go
// Use different model
WithModel(model string) ScoringOption

// Custom prompt template (Go template syntax)
WithPromptTemplate(template string) ScoringOption

// Additional context for all posts
WithExtraContext(context map[string]string) ScoringOption
```

## Hidden Context and Gotchas

### Fixed Batch Size Limitation
The batch size is hard-coded to 10 posts (`maxBatchSize = 10`). This cannot be configured and is based on OpenAI token limits. Attempting to score more than 10 posts will automatically split them into multiple batches.

### Score Range Validation
Scores MUST be integers between 0-100. Out-of-range scores from the API will cause validation errors. The system enforces this through JSON schema validation.

### Missing Scores Behavior
If OpenAI returns fewer scores than posts sent, missing scores default to 0 with the reason "Score not found in response". This is logged as a warning but doesn't fail the operation.

### Template Syntax Requirements
Custom prompts MUST either:
1. Include `%s` placeholder for sprintf-style formatting (legacy)
2. Use Go template syntax with valid field references: `{{.Posts}}`, `{{.PostTitle}}`, etc.

Failing to include proper placeholders will cause a validation error.

### Concurrent Processing Order
When `MaxConcurrent > 1`, batches are processed in parallel but results are always returned in the original input order. This is guaranteed through indexed channel collection.

### Environment Variable Loading
The examples use `godotenv` to load `.env` files, but the library itself doesn't require this. In production, set `OPENAI_API_KEY` directly in the environment.

### Model Selection Caveats
While you can specify any OpenAI model, the default prompts are optimized for GPT-4o-mini. Using different models may require prompt adjustments for optimal results.

### Import Path Changes
The library moved from local development to published module. Ensure imports use:
```go
import "github.com/JohnPlummer/llm-client/scorer"
```
Not the old local path references.

### Embedded Prompts Location
Prompts are embedded from `scorer/prompts/` at compile time. Modifying these files requires recompilation. For runtime customization, use `WithPromptTemplate()`.

### Test Framework Version
The project uses Ginkgo v2. When running tests, ensure you have the v2 CLI:
```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

## Version History and Evolution

### v0.9.1 (July 29, 2025) - Current
- Added functional options pattern for runtime configuration
- Introduced custom prompt templates with Go template support
- Added per-request model selection
- Enhanced error messages with template preview
- Improved test coverage (+257 lines)

### v0.9.0 (July 18, 2025) - Major Release
- Implemented 12 autonomous improvements
- Added concurrent batch processing
- Migrated to published module
- Enhanced error handling and logging
- Achieved 89.3% test coverage

### Breaking Changes
- **February 10, 2025**: Migrated from `log` to `slog` for structured logging
- **February 10, 2025**: Replaced `Debug` config flag with `LOG_LEVEL` environment variable

### Migration Guide
If upgrading from pre-v0.9.0:
1. Update import paths to `github.com/JohnPlummer/llm-client/scorer`
2. Replace `Debug: true` with `LOG_LEVEL=debug` environment variable
3. Update to Go 1.23.1+ for built-in `min` function support

## Performance Considerations

### Batch Processing Efficiency
- **Optimal**: 10 posts per batch (current implementation)
- **API calls**: N/10 calls for N posts
- **Cost**: Optimized for GPT-4o-mini pricing

### Concurrency Settings
- **Sequential** (`MaxConcurrent = 1`): ~1-2 seconds per batch
- **Concurrent** (`MaxConcurrent = 5`): Can process 50 posts in ~2-3 seconds
- **Maximum recommended**: 10 (to avoid rate limiting)

### Memory Usage
- Prompts embedded at compile time (no runtime file I/O)
- Batch processing limits memory usage
- CSV streaming in examples for large datasets

## Debugging Guide

### Common Issues

1. **"OpenAI API key is required" error**
   - Ensure `OPENAI_API_KEY` is set in environment
   - Check for typos in environment variable name
   - Verify key is valid with OpenAI

2. **"context deadline exceeded" errors**
   - Increase context timeout
   - Reduce `MaxConcurrent` to lower load
   - Check network connectivity

3. **Scores returning 0 for all posts**
   - Check API response in debug logs (`LOG_LEVEL=debug`)
   - Verify prompt template includes scoring instructions
   - Ensure posts have meaningful content

4. **"score out of range" validation errors**
   - Custom prompts must specify 0-100 score range
   - Check JSON schema in prompt matches expected format

### Debugging Tools
```bash
# Enable debug logging
export LOG_LEVEL=debug

# Run with verbose test output
ginkgo -v ./...

# Check for race conditions
go test -race ./...

# Inspect embedded prompts
go run -tags debug ./scorer/prompts.go
```

### Log Levels
- `debug`: Individual post scores, API requests/responses
- `info`: Batch completion, configuration
- `warn`: Missing scores, fallback behavior
- `error`: API failures, validation errors

## Resources

### Internal Documentation
- `README.md`: Primary documentation
- `docs/`: Comprehensive guides (7 documents)
- `examples/basic/`: Working example with CSV data
- `IMPROVEMENTS.md`: Completed enhancements tracking

### External Resources
- [OpenAI API Documentation](https://platform.openai.com/docs)
- [Ginkgo Testing Framework](https://onsi.github.io/ginkgo/)
- [Go OpenAI Client](https://github.com/sashabaranov/go-openai)

### Code Ownership
- **Maintainer**: John Plummer (52 commits, 100% contribution)
- **Core components**: `scorer/` package
- **Testing**: Comprehensive Ginkgo test suite
- **Documentation**: Extensive docs in `docs/` directory

## Contributing Guidelines

### Before Contributing
1. Read existing documentation in `docs/`
2. Check `IMPROVEMENTS.md` for completed work
3. Review test patterns in `scorer/scorer_test.go`

### Code Submission Process
1. Fork the repository
2. Create a feature branch
3. Write tests first (TDD approach)
4. Implement feature
5. Run `make check` to validate
6. Submit PR with clear description

### Code Review Checklist
- [ ] Tests pass (`make test`)
- [ ] Coverage maintained or improved (currently 86.5%)
- [ ] Documentation updated if needed
- [ ] Examples work (`make run`)
- [ ] No linting issues
- [ ] Error messages are descriptive
- [ ] Logging at appropriate levels

### Definition of Done
- Feature fully implemented with tests
- Documentation updated
- Examples demonstrate usage
- All quality gates pass
- Code follows established patterns