# Development Setup

## Prerequisites

- **Go**: Version 1.23.1 or later
- **OpenAI API Key**: Required for scoring functionality
- **Git**: For version control and dependency management

## Installation

### As a Library Dependency

```bash
go get github.com/JohnPlummer/llm-client
```

### For Development

```bash
git clone https://github.com/JohnPlummer/llm-client.git
cd llm-client
go mod tidy
```

## Dependencies

### Core Dependencies

- `github.com/sashabaranov/go-openai v1.37.0` - OpenAI API client
- `github.com/JohnPlummer/reddit-client v0.0.0` - Reddit post types (local dependency)

### Development Dependencies

- `github.com/onsi/ginkgo/v2 v2.23.4` - BDD testing framework
- `github.com/onsi/gomega v1.36.3` - Matcher library for tests

### Local Dependencies

The project uses a local `reddit-client` dependency via Go module replacement:
```go
replace github.com/JohnPlummer/reddit-client => ../reddit-client
```

## Development Commands

### Primary Commands

```bash
make test      # Run tests using Ginkgo framework
make run       # Run the basic example (examples/basic/main.go)
make check     # Run all checks: tidy, tests, and example
make coverage  # Generate coverage report in markdown format
```

### Module Management

```bash
make tidy          # Run go mod tidy in root project
make tidy-examples # Run go mod tidy in examples directory
make tidy-all      # Run go mod tidy everywhere
```

## Testing

### Framework

- **Testing Framework**: Ginkgo BDD (Behavior-Driven Development)
- **Assertion Library**: Gomega
- **Test Location**: `scorer/scorer_test.go`

### Running Tests

```bash
# Using Make (recommended)
make test

# Direct Ginkgo commands
ginkgo -v ./...
ginkgo -v --coverprofile=coverage.out ./...
```

### Coverage Reports

```bash
# Generate markdown coverage report
make coverage

# View coverage in browser
go tool cover -html=coverage.out
```

## Environment Configuration

### Required Environment Variables

```bash
export OPENAI_API_KEY="your-openai-api-key"
```

### Optional Configuration

The library accepts configuration through the `Config` struct:

```go
config := scorer.Config{
    OpenAIKey:     "your-api-key",     // Required
    PromptText:    "custom prompt",    // Optional
    MaxConcurrent: 5,                  // Optional, for rate limiting
}
```

## Code Style Guidelines

### Go Best Practices

- **Naming**: Use MixedCaps naming convention (no snake_case)
- **Standard Library**: Prioritize standard library solutions over external dependencies
- **State Management**: Avoid package-level state and global variables
- **Dependency Injection**: Use dependency injection patterns for testability
- **Function Design**: Keep functions modular and testable
- **Error Handling**: Provide meaningful error messages with context
- **Logging**: Accept `slog.Logger` as parameter rather than initializing internally

### Code Organization

- **Interface First**: Define interfaces before implementations
- **Separation of Concerns**: Keep business logic separate from I/O operations
- **Testability**: Design for easy unit testing and mocking
- **Documentation**: Include clear package and function documentation

### Project-Specific Guidelines

- **Batch Size**: Maintain maximum batch size of 10 posts per OpenAI API call
- **Error Resilience**: Always provide fallback behavior (default score of 0)
- **JSON Validation**: Strictly validate OpenAI response structure
- **Prompt Management**: Keep prompts in embedded files for version control

## IDE Configuration

### VS Code

For optimal development experience, ensure these extensions are installed:
- Go extension for syntax highlighting and IntelliSense
- Test Explorer for Ginkgo test integration

### Development Rules

The project includes development style guides in `.cursor/rules/` directory for AI-assisted development tools.

## Building and Testing

### Local Development Workflow

1. **Make changes** to source code
2. **Run tests** with `make test`
3. **Check coverage** with `make coverage`
4. **Test examples** with `make run`
5. **Validate everything** with `make check`

### Continuous Integration

The project is designed to support CI/CD pipelines with the `make check` command that runs:
- Module tidying
- Full test suite
- Example execution
- Coverage validation