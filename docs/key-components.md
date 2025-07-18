# Key Components

## Core Interfaces

### Scorer Interface
**Location**: `scorer/types.go:12-14`

The main interface that defines the scoring functionality:

```go
type Scorer interface {
    ScorePosts(ctx context.Context, posts []*reddit.Post) ([]*ScoredPost, error)
}
```

**Purpose**: Provides a clean abstraction for scoring Reddit posts, enabling easy testing and alternative implementations.

### OpenAIClient Interface  
**Location**: `scorer/types.go:30-33`

Defines the contract for OpenAI API interactions:

```go
type OpenAIClient interface {
    CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}
```

**Purpose**: Enables dependency injection and mocking for testing without making actual API calls.

## Core Types

### ScoredPost
**Location**: `scorer/types.go:16-21`

The primary output type containing scored results:

```go
type ScoredPost struct {
    Post   *reddit.Post
    Score  int
    Reason string
}
```

**Fields**:
- `Post`: Original Reddit post data
- `Score`: AI-generated relevance score (0-100)
- `Reason`: Explanation for the assigned score

### Config
**Location**: `scorer/types.go:23-28`

Configuration structure for scorer initialization:

```go
type Config struct {
    OpenAIKey     string
    PromptText    string
    MaxConcurrent int
}
```

**Fields**:
- `OpenAIKey`: Required OpenAI API key
- `PromptText`: Optional custom prompt (uses default if empty)
- `MaxConcurrent`: Optional rate limiting (not yet implemented)

## Implementation Details

### scorer struct
**Location**: `scorer/types.go:35-39`

The main implementation of the Scorer interface:

```go
type scorer struct {
    client OpenAIClient
    config Config
    prompt string
}
```

**Responsibilities**:
- Manages OpenAI client instance
- Handles batch processing logic
- Applies configured prompts
- Processes and validates responses

## Processing Pipeline

### 1. Batch Creation
**Location**: `scorer/scorer.go:68-87`

The `ScorePosts` method splits input posts into batches of maximum 10 posts each:

```go
for i := 0; i < len(posts); i += maxBatchSize {
    results, err := s.processBatch(ctx, posts[i:min(i+maxBatchSize, len(posts))])
    // ...
}
```

**Key Features**:
- Automatic batching with `maxBatchSize = 10`
- Sequential batch processing
- Comprehensive error handling with batch context

### 2. API Request Formation
**Location**: `scorer/batch.go:37-70`

The `createChatCompletion` method constructs structured API requests:

- **System Prompt**: Loaded from embedded `prompts/system_prompt.txt`
- **User Prompt**: Formatted JSON containing posts to score
- **Response Format**: Enforced JSON schema validation
- **Model**: Uses `openai.GPT4oMini` for cost efficiency

### 3. Response Processing
**Location**: `scorer/batch.go:72-94`

The `parseResponse` method handles API responses with robust validation:

```go
func (s *scorer) parseResponse(resp *openai.ChatCompletionResponse, expectedPostCount int) (map[string]scoreItem, error)
```

**Validation Steps**:
- JSON unmarshaling with error handling
- Score range validation (0-100)
- Missing score detection and logging
- Response completeness checking

### 4. Result Assembly
**Location**: `scorer/batch.go:96-127`

The `createScoredPosts` method creates final results with graceful fallbacks:

**Fallback Behavior**:
- Missing scores default to 0
- Automatic reason generation for missing scores
- Comprehensive logging for debugging

## Prompt Management

### System Prompt
**Location**: `scorer/prompts/system_prompt.txt`

Embedded prompt that defines the AI's role and scoring criteria:
- Focuses on location-based recommendations and events
- Requires scoring of ALL input posts
- Emphasizes use of comments for additional context

### Custom Prompts
**Location**: `scorer/prompts.go`

Default batch processing prompt with JSON formatting requirements:
- Must include `%s` placeholder for post injection
- Must produce valid JSON responses
- Must include version and scores array structure

## Data Flow

```
Input Posts
    ↓
Batch Splitting (max 10 posts)
    ↓
JSON Formatting
    ↓
OpenAI API Request
    ↓
JSON Schema Validation
    ↓
Response Parsing
    ↓
Score Validation (0-100)
    ↓
ScoredPost Creation
    ↓
Result Assembly
```

## Error Handling Strategy

### API Errors
- Network failures: Propagated with context
- Invalid responses: Detailed error messages
- Schema violations: Immediate failure

### Data Errors
- Missing scores: Default to 0 with logging
- Invalid score ranges: Validation failure
- Malformed JSON: Parsing failure

### Batch Errors
- Individual batch failures stop processing
- Comprehensive error context provided
- Batch number included in error messages

## Testing Architecture

### Mock Support
The `OpenAIClient` interface enables comprehensive testing:
- API call mocking without external dependencies
- Response simulation for edge cases
- Performance testing with controlled responses

### Test Coverage
**Location**: `scorer/scorer_test.go`

Uses Ginkgo BDD framework for:
- Interface compliance testing
- Error condition validation
- Response parsing verification
- Batch processing validation