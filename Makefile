# Chat Agent Makefile

# Variables
APP_NAME = chat-agent
BINARY_NAME = $(APP_NAME)
DOCKER_IMAGE = $(APP_NAME):latest
MAIN_PATH = ./cmd/server
BUILD_DIR = ./bin

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOFMT = $(GOCMD) fmt

# Build flags
BUILD_FLAGS = -ldflags="-w -s"

.PHONY: all build clean test deps fmt vet run docker-build docker-run docker-stop help

# Default target
all: clean deps fmt vet test build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Test with coverage
test-coverage: test
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

# Run with hot reload (requires air)
dev:
	@echo "Running $(APP_NAME) with hot reload..."
	air

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running Docker container..."
	docker-compose up -d

docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f chat-agent

# Database commands
db-migrate:
	@echo "Running database migrations..."
	$(GOCMD) run $(MAIN_PATH) migrate

db-seed:
	@echo "Seeding database..."
	$(GOCMD) run $(MAIN_PATH) seed

# Linting
lint:
	@echo "Running linter..."
	golangci-lint run

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/cosmtrek/air@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Generate swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g $(MAIN_PATH)/main.go

# Security scan
security:
	@echo "Running security scan..."
	gosec ./...

# Performance test
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Create release
release: clean deps fmt vet test build-all
	@echo "Creating release..."
	@mkdir -p release
	@cp -r $(BUILD_DIR)/* release/
	@tar -czf release/$(APP_NAME)-$(shell date +%Y%m%d-%H%M%S).tar.gz -C release .

# Help
help:
	@echo "Available commands:"
	@echo "  build          Build the application"
	@echo "  clean          Clean build artifacts"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage"
	@echo "  deps           Download dependencies"
	@echo "  fmt            Format code"
	@echo "  vet            Vet code"
	@echo "  run            Run the application"
	@echo "  dev            Run with hot reload"
	@echo "  build-all      Build for multiple platforms"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run with Docker Compose"
	@echo "  docker-stop    Stop Docker containers"
	@echo "  docker-logs    Show Docker logs"
	@echo "  lint           Run linter"
	@echo "  security       Run security scan"
	@echo "  bench          Run benchmarks"
	@echo "  release        Create release package"
	@echo "  help           Show this help message"

