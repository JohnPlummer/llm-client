# Project Overview

## Purpose

LLM Client is a Go library that leverages OpenAI's GPT API to intelligently score text content based on customizable criteria. The library is designed to help developers automatically evaluate and rank any text content for relevance, quality, or other specified metrics.

## Core Functionality

The scorer evaluates text items and returns structured results containing:
- The original text data
- A relevance score (0-100 scale)
- A detailed reason explaining the score

## Architecture

This is a Go library built around the `scorer` package, which provides a clean interface for batch processing text content through OpenAI's API.

### Core Components

- **Scorer Interface** (`scorer/types.go`): Main interface defining the `ScoreTexts()` method
- **scorer struct** (`scorer/scorer.go`): Primary implementation with OpenAI client integration
- **Batch Processing** (`scorer/batch.go`): Handles efficient processing of text items in batches (maximum 10 per batch)
- **Types** (`scorer/types.go`): Core data structures including `ScoredItem`, `TextItem`, `Config`, and `scoreResponse`

### Key Features

- **Batch Processing**: Efficiently processes up to 10 text items per API call
- **JSON Schema Validation**: Ensures OpenAI responses conform to expected structure
- **Graceful Error Handling**: Assigns default scores (0) for missing or invalid responses
- **Embedded Prompt System**: Includes default scoring criteria with location-based recommendations focus
- **Custom Prompt Support**: Allows complete customization via `Config.PromptText`
- **Rate Limiting**: Built-in support for concurrent request limiting

### Technology Stack

- **Language**: Go 1.23.1+
- **AI Integration**: OpenAI GPT API via `github.com/sashabaranov/go-openai`
- **Data Source**: Generic text items via `TextItem` struct with ID, Content, and Metadata fields
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

- **Content Curation**: Automatically identify high-quality text content for various applications
- **Location-Based Filtering**: Score text items based on geographic relevance (configurable via metadata)
- **Custom Criteria Scoring**: Evaluate text content against any custom criteria using tailored prompts
- **Batch Content Analysis**: Process large volumes of text content efficiently
- **Quality Assessment**: Rank text content by relevance, quality, or informational value