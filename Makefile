.PHONY: help build run test clean docker-build docker-up docker-down lint format

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  test-cover   - Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up    - Start Docker services"
	@echo "  docker-down  - Stop Docker services"
	@echo "  lint         - Run linter"
	@echo "  format       - Format code"
	@echo "  deps         - Install dependencies"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/api-test-framework cmd/api-server/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/api-server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f *.exe
	rm -f api-test-framework
	rm -f debug.log

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t api-test-framework .

# Start Docker services
docker-up:
	@echo "Starting Docker services..."
	docker-compose -f docker-compose.dev.yml up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	docker-compose -f docker-compose.dev.yml down

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
format:
	@echo "Formatting code..."
	go fmt ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Development setup
dev-setup: deps docker-up
	@echo "Development environment setup complete!"

# Production build
prod-build:
	@echo "Building production binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o bin/api-test-framework cmd/api-server/main.go
