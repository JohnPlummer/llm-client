.PHONY: run test tidy tidy-examples check coverage

# Run the basic example
run:
	@echo "Running basic example..."
	cd examples/basic && go run main.go

# Run tests using Ginkgo
test:
	@echo "Running tests..."
	ginkgo -v ./...

# Run go mod tidy in root project
tidy:
	@echo "Running go mod tidy in root project..."
	go mod tidy

# Run go mod tidy in examples
tidy-examples:
	@echo "Running go mod tidy in examples..."
	cd examples/basic && go mod tidy

# Run go mod tidy everywhere
tidy-all: tidy tidy-examples

# Run all checks: tidy, lint, test, and run example
check:
	@echo "Running all checks..."
	@echo "Step 1: Running tidy-all..."
	@make tidy-all
	@echo "Step 2: Running tests..."
	@make test
	@echo "Step 3: Running example..."
	@make run
	@echo "All checks completed successfully!"

# Generate coverage report in markdown format
coverage:
	@echo "Generating coverage report..."
	ginkgo -v --coverprofile=coverage.out ./...
	@echo "# Coverage Report\n" > coverage.md
	@echo "| Function | Coverage |" >> coverage.md
	@echo "|----------|----------|" >> coverage.md
	@go tool cover -func=coverage.out | grep -v "total:" | awk '{printf "| %s | %s |\n", $$1, $$3}' >> coverage.md
	@echo "\n## Total Coverage" >> coverage.md
	@go tool cover -func=coverage.out | grep "total:" | awk '{printf "**%s**\n", $$3}' >> coverage.md
	@echo "Coverage report generated: coverage.md"
