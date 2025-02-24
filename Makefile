.PHONY: run test tidy tidy-examples

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
