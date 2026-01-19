.PHONY: generate build test lint clean

# Generate swagger documentation
generate:
	swag init -g cmd/shorty-server/main.go -o api/swagger

# Build the server
build:
	go build -o bin/shorty-server ./cmd/shorty-server

# Run all tests
test:
	go test -v ./pkg/...
	go test -v ./tests/integration/...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/

# Install development tools
tools:
	go install github.com/swaggo/swag/cmd/swag@latest

# Check that generated files are up to date
check-generate: generate
	@if [ -n "$$(git status --porcelain api/swagger/)" ]; then \
		echo "Error: Generated files are out of date. Run 'make generate' and commit the changes."; \
		git status --porcelain api/swagger/; \
		git diff api/swagger/; \
		exit 1; \
	fi
	@echo "Generated files are up to date."
