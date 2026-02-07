.PHONY: generate generate-clients build test lint clean

# Generate swagger documentation
generate: generate-swagger generate-clients

# Generate swagger docs only
generate-swagger:
	swag init -g cmd/shorty-server/main.go -o api/swagger

# Generate client libraries from OpenAPI spec
generate-clients:
	./scripts/generate-clients.sh

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
	go install github.com/swaggo/swag/cmd/swag@v1.16.4
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1

# Check that generated files are up to date
check-generate: generate
	@if [ -n "$$(git status --porcelain api/swagger/docs.go api/swagger/swagger.json api/swagger/swagger.yaml)" ]; then \
		echo "Error: Generated swagger files are out of date. Run 'make generate' and commit the changes."; \
		git status --porcelain api/swagger/; \
		git diff api/swagger/; \
		exit 1; \
	fi
	@echo "Generated files are up to date."
