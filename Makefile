# Makefile for Geopolitical Domain Service

.PHONY: help build run test clean docker-build docker-run docker-compose-up docker-compose-down

# Variables
APP_NAME=geopolitical-service
DOCKER_IMAGE=geopolitical-service:latest
GO_VERSION=1.21

# Help target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) cmd/server/main.go

build-linux: ## Build for Linux
	@echo "Building $(APP_NAME) for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux cmd/server/main.go

# Run targets
run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run cmd/server/main.go

run-dev: ## Run with development environment
	@echo "Running $(APP_NAME) in development mode..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	go run cmd/server/main.go

# Test targets
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint and format
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Docker Compose targets
docker-compose-up: ## Start all services with docker-compose
	@echo "Starting services with docker-compose..."
	docker-compose up -d

docker-compose-down: ## Stop all services
	@echo "Stopping services..."
	docker-compose down

docker-compose-logs: ## Show logs from all services
	docker-compose logs -f

# Database targets
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	# Add migration command here

db-seed: ## Seed database with initial data
	@echo "Seeding database..."
	# Add seed command here

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

clean-docker: ## Clean Docker images and containers
	@echo "Cleaning Docker artifacts..."
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE) 2>/dev/null || true

# Development targets
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file from template"; fi
	go mod download
	@echo "Development environment ready!"

# Production targets
prod-build: ## Build for production
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o bin/$(APP_NAME) cmd/server/main.go

# Kubernetes targets (for future cloud deployment)
k8s-deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	# kubectl apply -f k8s/

k8s-delete: ## Delete from Kubernetes
	@echo "Deleting from Kubernetes..."
	# kubectl delete -f k8s/

# Default target
.DEFAULT_GOAL := help