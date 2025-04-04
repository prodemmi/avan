# Define variables
GO=go
APP_NAME=messaging-service
BIN_DIR=bin
BIN_NAME=$(BIN_DIR)/$(APP_NAME)

# Default target (help)
.PHONY: help
help:
	@echo "Makefile commands for $(APP_NAME)"
	@echo "Available commands:"
	@echo "  build       Build the Go executable"
	@echo "  run         Run the server"
	@echo "  test        Run tests"
	@echo "  fmt         Format Go code"
	@echo "  lint        Lint the code"
	@echo "  clean       Clean the bin directory"
	@echo "  docker      Build Docker image"
	@echo "  metrics     Run the Prometheus metrics"

# Build the Go executable
.PHONY: build
build: clean
	@echo "Building the project..."
	$(GO) build -o $(BIN_NAME) ./cmd/server/main.go  # Point to the main.go file

# Run the server
.PHONY: run
run:
	@echo "Starting server on port $(SERVER_PORT)..."
	$(BIN_NAME)

# Run sqlc
.PHONY: sqlc
sqlc:
	@echo "🚀 Generating sqlc data..."
	@rm -rf ./internal/db/sqlc/
	@sqlc generate
	@sleep 2

# Run sqlc
.PHONY: templ
templ:
	@echo "🚀 Generating templ data..."
	@templ generate

# Run the dev server
.PHONY: dev
dev:
	@echo "🚀 Starting dev server with hot reload..."
	$(MAKE) sqlc
	@templ generate --watch &
	@npm run dev &
	@air

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./...

# Format Go code
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	$(GO) fmt ./...

# Lint the code (using golangci-lint)
.PHONY: lint
lint:
	@echo "Linting Go code..."
	golangci-lint run

# Clean the bin directory
.PHONY: clean
clean:
	@echo "Cleaning bin directory..."
	rm -rf $(BIN_DIR)

# Docker image build
.PHONY: docker
docker:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

# Prometheus metrics (Test the /metrics endpoint locally)
.PHONY: metrics
metrics:
	@echo "Testing metrics endpoint..."
	curl http://localhost:8080/metrics
