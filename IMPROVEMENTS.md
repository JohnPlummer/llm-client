# Post Scorer Improvement Plan

This document outlines a sequence of improvements for the post-scorer library. Each step should be implemented independently, tested, and committed before moving to the next.

## Instructions

1. Implement each improvement step by step
2. After each step, run `make check` to ensure all tests pass and fix any errors or warnings
3. Mark the step as complete by adding ✅ at the beginning of the step title
4. Commit the changes with a descriptive message before moving to the next step

---

## ✅ Step 1: Fix Init Error Handling

**Problem**: The library calls `os.Exit(1)` in init() which is inappropriate for a library.

**Files to modify**:

- `scorer/scorer.go`

**Implementation**:

1. Add a package-level error variable to capture init errors
2. Modify init() to set the error instead of exiting
3. Check for init error in New() function
4. Update tests to handle the new error case

**Expected changes**:

```go
// At package level
var initError error

// In init()
if err != nil {
    initError = fmt.Errorf("failed to load system prompt: %w", err)
    return
}

// In New()
if initError != nil {
    return nil, initError
}
```

After implementation:
- [ ] Run `make check`
- [ ] Fix any errors or warnings
- [ ] Mark as complete: ✅ Step 1: Fix Init Error Handling
- [ ] Commit: `git commit -m "fix: remove os.Exit from library init function"`

---

## ✅ Step 2: Fix Documentation Inconsistencies

**Problem**: README examples have type errors and incorrect slice types.

**Files to modify**:
- `README.md`

**Implementation**:
1. Change line 40 from `[]reddit.Post` to `[]*reddit.Post`
2. Change line 54 from `%.2f` to `%d` for Score formatting
3. Ensure all code examples compile correctly

After implementation:
- [ ] Run `make check`
- [ ] Verify examples still work with `make run`
- [ ] Mark as complete: ✅ Step 2: Fix Documentation Inconsistencies
- [ ] Commit: `git commit -m "docs: fix README code examples and type errors"`

---

## ✅ Step 3: Add Config Validation

**Problem**: No validation for required fields and prompt format.

**Files to modify**:
- `scorer/scorer.go`

**Implementation**:
1. Add validation in New() for empty OpenAIKey when client is nil
2. Add validation for custom prompt to ensure it contains `%s` placeholder
3. Add validation for MaxConcurrent to be >= 0
4. Return appropriate errors with context

**Expected changes**:
```go
// In New()
if cfg.OpenAIKey == "" && client == nil {
    return nil, errors.New("OpenAI API key is required")
}

if cfg.PromptText != "" && !strings.Contains(cfg.PromptText, "%s") {
    return nil, errors.New("custom prompt must contain %s placeholder for posts")
}

if cfg.MaxConcurrent < 0 {
    return nil, errors.New("MaxConcurrent must be non-negative")
}
```

After implementation:
- [ ] Run `make check`
- [ ] Add tests for validation errors
- [ ] Mark as complete: ✅ Step 3: Add Config Validation
- [ ] Commit: `git commit -m "feat: add comprehensive config validation"`

---

## ✅ Step 4: Improve Error Context

**Problem**: Errors lack sufficient context for debugging.

**Files to modify**:
- `scorer/batch.go`
- `scorer/scorer.go`

**Implementation**:
1. Wrap all errors with context using fmt.Errorf
2. Include relevant information (batch number, post count, etc.)
3. Use consistent error message format

**Example**:
```go
// In processBatch
if err != nil {
    return nil, fmt.Errorf("failed to create chat completion for batch of %d posts: %w", len(posts), err)
}
```

After implementation:
- [ ] Run `make check`
- [ ] Fix any errors or warnings
- [ ] Mark as complete: ✅ Step 4: Improve Error Context
- [ ] Commit: `git commit -m "feat: add comprehensive error context throughout"`

---

## ✅ Step 5: Add Input Validation

**Problem**: No validation for nil posts or empty slices.

**Files to modify**:
- `scorer/scorer.go`

**Implementation**:
1. Add validation at the start of ScorePosts
2. Check for nil posts in the slice
3. Handle empty slice gracefully (return empty result)

**Expected changes**:
```go
// In ScorePosts
if posts == nil {
    return nil, errors.New("posts cannot be nil")
}

if len(posts) == 0 {
    return []*ScoredPost{}, nil
}

for i, post := range posts {
    if post == nil {
        return nil, fmt.Errorf("post at index %d is nil", i)
    }
    if post.ID == "" {
        return nil, fmt.Errorf("post at index %d has empty ID", i)
    }
}
```

After implementation:
- [ ] Run `make check`
- [ ] Add tests for validation cases
- [ ] Mark as complete: ✅ Step 5: Add Input Validation
- [ ] Commit: `git commit -m "feat: add input validation for posts"`

---

## ✅ Step 6: Extract Helper Functions

**Problem**: Complex functions that do too many things.

**Files to modify**:
- `scorer/batch.go`
- `examples/basic/main.go`

**Implementation**:
1. Extract JSON building logic from processBatch into buildChatRequest
2. Extract response validation logic into validateResponse
3. In example, extract CSV loading, scoring, and output into separate functions

After implementation:
- [ ] Run `make check`
- [ ] Ensure refactoring doesn't break functionality
- [ ] Mark as complete: ✅ Step 6: Extract Helper Functions
- [ ] Commit: `git commit -m "refactor: extract helper functions for better readability"`

---

## ✅ Step 7: Remove Min Function Duplication

**Problem**: Custom min function when Go 1.21+ has built-in min.

**Files to modify**:
- `scorer/batch.go`
- `go.mod` (ensure Go version is 1.21+)

**Implementation**:
1. Remove the custom min function
2. Use built-in min function
3. Update go.mod if needed to require Go 1.21+

After implementation:
- [ ] Run `make check`
- [ ] Fix any errors or warnings
- [ ] Mark as complete: ✅ Step 7: Remove Min Function Duplication
- [ ] Commit: `git commit -m "refactor: use built-in min function from Go 1.21+"`

---

## ✅ Step 8: Add Missing Test Cases

**Problem**: Low test coverage in critical areas.

**Files to modify**:
- `scorer/scorer_test.go`

**Implementation**:
1. Add test for prompt validation error
2. Add test for JSON marshaling error path
3. Add test for nil posts validation
4. Add test for empty ID validation
5. Increase coverage for New() function

After implementation:
- [ ] Run `make check`
- [ ] Run `make coverage` and verify improvement
- [ ] Mark as complete: ✅ Step 8: Add Missing Test Cases
- [ ] Commit: `git commit -m "test: add comprehensive test coverage for edge cases"`

---

## ✅ Step 9: Create Custom Prompt Example

**Problem**: Custom prompt example file is referenced but missing.

**Files to create**:
- `examples/basic/custom_prompt.txt`

**Implementation**:
1. Create the custom prompt file with proper format
2. Ensure it includes the required %s placeholder
3. Add example showing different scoring criteria

After implementation:
- [ ] Run `make check`
- [ ] Test with `make run` to ensure it works
- [ ] Mark as complete: ✅ Step 9: Create Custom Prompt Example
- [ ] Commit: `git commit -m "docs: add custom prompt example file"`

---

## ✅ Step 10: Implement MaxConcurrent Feature

**Problem**: MaxConcurrent is documented but not implemented.

**Files to modify**:
- `scorer/scorer.go`
- `scorer/batch.go`
- `scorer/scorer_test.go`

**Implementation**:
1. Add concurrent batch processing logic
2. Use semaphore pattern with buffered channel
3. Handle context cancellation properly
4. Add tests for concurrent processing
5. Set default MaxConcurrent to 1 for backward compatibility

**High-level approach**:
```go
func (s *scorer) processConcurrently(ctx context.Context, batches [][]*reddit.Post) ([]*ScoredPost, error) {
    if s.config.MaxConcurrent <= 1 {
        // Fall back to sequential processing
        return s.processSequentially(ctx, batches)
    }
    
    sem := make(chan struct{}, s.config.MaxConcurrent)
    // ... concurrent implementation
}
```

After implementation:
- [ ] Run `make check`
- [ ] Add tests for concurrent processing
- [ ] Verify no race conditions with `go test -race`
- [ ] Mark as complete: ✅ Step 10: Implement MaxConcurrent Feature
- [ ] Commit: `git commit -m "feat: implement MaxConcurrent for parallel batch processing"`

---

## ✅ Step 11: Improve Prompt Management

**Problem**: Inconsistent prompt storage (init vs constant).

**Files to modify**:
- `scorer/prompts/batch_prompt.txt` (create)
- `scorer/batch.go`
- `scorer/scorer.go`

**Implementation**:
1. Create batch_prompt.txt in prompts directory
2. Load it the same way as system_prompt.txt
3. Remove the hardcoded constant
4. Ensure consistent error handling

After implementation:
- [ ] Run `make check`
- [ ] Fix any errors or warnings
- [ ] Mark as complete: ✅ Step 11: Improve Prompt Management
- [ ] Commit: `git commit -m "refactor: unify prompt management with embedded files"`

---

## Step 12: Fix Local Dependency Issue

**Problem**: Local replace directive makes library unusable for external users.

**Files to modify**:
- Document the issue in README.md
- Add comment in go.mod explaining the local replacement

**Implementation**:
1. Add a section in README explaining the reddit-client dependency
2. Document how to handle this for external users
3. Consider options for properly publishing reddit-client

After implementation:
- [ ] Run `make check`
- [ ] Mark as complete: ✅ Step 12: Fix Local Dependency Issue
- [ ] Commit: `git commit -m "docs: document local dependency handling for external users"`

---

## Completion Summary

After completing all steps:
1. Review the complete set of changes
2. Update version in documentation if appropriate
3. Create a final commit summarizing all improvements
4. Consider creating a new release tag

Total improvements implemented:
- [ ] Init error handling
- [ ] Documentation fixes
- [ ] Config validation
- [ ] Error context
- [ ] Input validation
- [ ] Code refactoring
- [ ] Duplicate code removal
- [ ] Test coverage
- [ ] Missing examples
- [ ] MaxConcurrent feature
- [ ] Prompt management
- [ ] Dependency documentation

Remember: Run `make check` after EVERY step and commit your changes!