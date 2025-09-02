# Key Components

## Core Interfaces

### Scorer Interface
**Location**: `scorer/types.go:12-14`

The main interface that defines the scoring functionality:

```go
type Scorer interface {
    // ScoreTexts scores a slice of text items
    ScoreTexts(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error)
    
    // ScoreTextsWithOptions scores text items with runtime options
    ScoreTextsWithOptions(ctx context.Context, items []TextItem, opts ...ScoringOption) ([]ScoredItem, error)
    
    // GetHealth returns the current health status of the scorer
    GetHealth(ctx context.Context) HealthStatus
}
```

**Purpose**: Provides a clean abstraction for scoring generic text content, enabling easy testing and alternative implementations.

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

### ScoredItem
**Location**: `scorer/types.go:19-24`

The primary output type containing scored results:

```go
type ScoredItem struct {
    Item   TextItem // Original text item
    Score  int      // Score between 0-100
    Reason string   // AI explanation for the score
}
```

**Fields**:
- `Item`: Original text item data (TextItem with ID, Content, Metadata)
- `Score`: AI-generated relevance score (0-100)
- `Reason`: Explanation for the assigned score

### Config
**Location**: `scorer/types.go:23-28`

Configuration structure for scorer initialization:

```go
type Config struct {
    APIKey               string                // OpenAI API key (required)
    Model                string                // OpenAI model to use
    PromptText           string                // Custom prompt template
    MaxConcurrent        int                   // Maximum concurrent API calls
    EnableCircuitBreaker bool                  // Enable circuit breaker pattern
    EnableRetry          bool                  // Enable retry with backoff
    Timeout              time.Duration         // Request timeout
    CircuitBreakerConfig *CircuitBreakerConfig // Circuit breaker configuration
    RetryConfig          *RetryConfig          // Retry configuration
}
```

**Fields**:
- `APIKey`: Required OpenAI API key
- `Model`: Optional model selection (defaults to GPT-4o-mini)
- `PromptText`: Optional custom prompt template
- `MaxConcurrent`: Optional concurrent request limiting
- `EnableCircuitBreaker`: Enable circuit breaker for resilience
- `EnableRetry`: Enable retry with backoff
- `Timeout`: Request timeout duration

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

The `ScoreTexts` method splits input text items into batches of maximum 10 items each:

```go
for i := 0; i < len(items); i += maxBatchSize {
    results, err := s.processBatch(ctx, items[i:min(i+maxBatchSize, len(items))])
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
- **User Prompt**: Formatted JSON containing text items to score
- **Response Format**: Enforced JSON schema validation
- **Model**: Uses `openai.GPT4oMini` for cost efficiency

### 3. Response Processing
**Location**: `scorer/batch.go:72-94`

The `parseResponse` method handles API responses with robust validation:

```go
func (s *scorer) parseResponse(resp *openai.ChatCompletionResponse, expectedItemCount int) (map[string]scoreItem, error)
```

**Validation Steps**:
- JSON unmarshaling with error handling
- Score range validation (0-100)
- Missing score detection and logging
- Response completeness checking

### 4. Result Assembly
**Location**: `scorer/batch.go:96-127`

The `createScoredItems` method creates final results with graceful fallbacks:

**Fallback Behavior**:
- Missing scores default to 0
- Automatic reason generation for missing scores
- Comprehensive logging for debugging

## Prompt Management

### System Prompt
**Location**: `scorer/prompts/system_prompt.txt`

Embedded prompt that defines the AI's role and scoring criteria:
- Focuses on generic text content scoring
- Requires scoring of ALL input text items
- Utilizes metadata for additional context

### Custom Prompts
**Location**: `scorer/prompts.go`

Default batch processing prompt with JSON formatting requirements:
- Must include `%s` placeholder for text item injection
- Must produce valid JSON responses
- Must include version and scores array structure

## Data Flow

```
Input Text Items
    ↓
Batch Splitting (max 10 items)
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
ScoredItem Creation
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