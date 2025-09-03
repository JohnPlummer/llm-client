// Package scorer provides a production-ready Go library for scoring text content
// using OpenAI's GPT API to determine relevance for various use cases.
//
// The library provides batch processing, concurrent execution, and flexible prompt
// customization with comprehensive error handling and resilience patterns.
//
// Features:
//   - Batch processing of text items (10 items per API call for efficiency)
//   - Concurrent processing with configurable parallelism
//   - Interface-first design for testing and extensibility
//   - Custom prompt templates with Go template support
//   - Circuit breaker pattern for resilience
//   - Retry logic with exponential backoff
//   - Prometheus metrics integration
//   - Content validation and sanitization utilities
//
// Basic usage:
//
//	cfg := scorer.Config{
//	    OpenAIKey:     os.Getenv("OPENAI_API_KEY"),
//	    MaxConcurrent: 5,
//	}
//	s, err := scorer.NewScorer(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	results, err := s.ScoreTexts(ctx, items)
package scorer
