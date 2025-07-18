# CLAUDE.md

This file provides AI assistant guidance for working with this Go library project. For comprehensive documentation, see the `docs/` directory.

## Documentation Reference

For detailed information, refer to:
- **[Project Overview](docs/project-overview.md)** - Architecture, features, and use cases
- **[Development Setup](docs/development-setup.md)** - Installation, dependencies, and coding standards
- **[Key Components](docs/key-components.md)** - Core interfaces, types, and implementation details
- **[Package Usage](docs/package-usage.md)** - API reference, examples, and custom prompts
- **[Deployment Guide](docs/deployment-guide.md)** - Production deployment and configuration
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and debugging strategies
- **[Recent Changes](docs/recent-changes.md)** - Latest updates and improvements

## Development Commands

### Primary Commands
- `make test` - Run tests using Ginkgo BDD framework
- `make run` - Run the basic example (examples/basic/main.go)
- `make check` - Run all checks: tidy, tests, and example validation
- `make coverage` - Generate coverage report in markdown format

### Module Management
- `make tidy` - Run go mod tidy in root project
- `make tidy-examples` - Run go mod tidy in examples directory
- `make tidy-all` - Run go mod tidy everywhere

### Testing Framework
- **Framework**: Ginkgo BDD with Gomega matchers
- **Test Files**: `scorer/scorer_test.go`
- **Direct Commands**:
  - `ginkgo -v ./...` - Run tests with verbose output
  - `ginkgo -v --coverprofile=coverage.out ./...` - Generate coverage

## AI Assistant Context

### Project Type
This is a **Go library** that scores Reddit posts using OpenAI's GPT API. The primary package is `scorer` with batch processing capabilities.

### Key Development Patterns
- **Interface-first design**: `Scorer` interface enables testing and mocking
- **Dependency injection**: `OpenAIClient` interface for API abstraction
- **Graceful degradation**: Missing scores default to 0 with logging
- **Batch processing**: Fixed 10-post batches for API efficiency
- **Embedded resources**: Prompts stored in `scorer/prompts/` directory

### Critical Implementation Details
- **Batch size limit**: Maximum 10 posts per OpenAI API call
- **Score validation**: Must be integers 0-100, validated in `parseResponse()`
- **Error handling**: Missing scores get default value with warning logs
- **JSON schema**: OpenAI responses validated against `scoreResponse` struct
- **Model selection**: Hard-coded to `openai.GPT4oMini` in batch processing

### Testing Requirements
- Must use Ginkgo BDD framework for new tests
- Mock `OpenAIClient` interface for unit tests
- Validate error conditions and edge cases
- Run `make check` before code changes