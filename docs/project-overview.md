# Project Overview

## Purpose

Post Scorer is a Go library that leverages OpenAI's GPT API to intelligently score Reddit posts based on customizable criteria. The library is designed to help developers automatically evaluate and rank social media content for relevance, quality, or other specified metrics.

## Core Functionality

The scorer evaluates Reddit posts and returns structured results containing:
- The original post data
- A relevance score (0-100 scale)
- A detailed reason explaining the score

## Architecture

This is a Go library built around the `scorer` package, which provides a clean interface for batch processing Reddit posts through OpenAI's API.

### Core Components

- **Scorer Interface** (`scorer/types.go`): Main interface defining the `ScorePosts()` method
- **scorer struct** (`scorer/scorer.go`): Primary implementation with OpenAI client integration
- **Batch Processing** (`scorer/batch.go`): Handles efficient processing of posts in batches (maximum 10 per batch)
- **Types** (`scorer/types.go`): Core data structures including `ScoredPost`, `Config`, and `scoreResponse`

### Key Features

- **Batch Processing**: Efficiently processes up to 10 Reddit posts per API call
- **JSON Schema Validation**: Ensures OpenAI responses conform to expected structure
- **Graceful Error Handling**: Assigns default scores (0) for missing or invalid responses
- **Embedded Prompt System**: Includes default scoring criteria with location-based recommendations focus
- **Custom Prompt Support**: Allows complete customization via `Config.PromptText`
- **Rate Limiting**: Built-in support for concurrent request limiting

### Technology Stack

- **Language**: Go 1.23.1+
- **AI Integration**: OpenAI GPT API via `github.com/sashabaranov/go-openai`
- **Data Source**: Reddit posts via `github.com/JohnPlummer/reddit-client`
- **Testing**: Ginkgo BDD framework with Gomega matchers
- **Build System**: Make-based with comprehensive targets

## Project Structure

```
llm-client/
├── scorer/               # Main library package
│   ├── scorer.go        # Core implementation
│   ├── types.go         # Interface and type definitions
│   ├── batch.go         # Batch processing logic
│   ├── prompts.go       # Prompt management
│   ├── scorer_test.go   # Test suite
│   └── prompts/         # Embedded prompt templates
│       └── system_prompt.txt
├── examples/
│   └── basic/           # Usage examples (separate go.mod)
│       ├── main.go
│       ├── custom_prompt.txt
│       └── *.csv        # Sample data files
├── docs/                # Comprehensive documentation
├── CLAUDE.md           # AI assistant development context
├── README.md           # Quick start guide
├── Makefile            # Build and development commands
└── go.mod              # Go module definition
```

## Use Cases

- **Content Curation**: Automatically identify high-quality posts for community highlights
- **Location-Based Filtering**: Score posts based on geographic relevance (default behavior)
- **Custom Criteria Scoring**: Evaluate posts against any custom criteria using tailored prompts
- **Batch Content Analysis**: Process large volumes of social media content efficiently
- **Quality Assessment**: Rank posts by engagement potential or informational value