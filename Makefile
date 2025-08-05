# Makefile for CodeBase Golang

.PHONY: build run test clean deps docker-build docker-run docker-compose-up docker-compose-down deploy-prep

# Application name
APP_NAME=codebase-golang
DOCKER_IMAGE=codebase-golang:latest

# Build the application
build:
	@echo "Building application..."
	mkdir -p bin
	go build -o bin/$(APP_NAME) main.go

# Run the application locally
run:
	@echo "Running application..."
	go run main.go

# Run with air for hot reload (install with: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Running in development mode with hot reload..."
	air

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f $(APP_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Security check
security:
	@echo "Running security check..."
	gosec ./...

# Setup environment
setup:
	@echo "Setting up environment..."
	cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

# Run with Docker
docker-run:
	@echo "Running with Docker..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Run with docker-compose
docker-compose-up:
	@echo "Starting with docker-compose..."
	docker-compose up -d

# Stop docker-compose
docker-compose-down:
	@echo "Stopping docker-compose..."
	docker-compose down

# Run with nginx proxy
docker-compose-nginx:
	@echo "Starting with nginx proxy..."
	docker-compose --profile nginx up -d

# View logs
logs:
	docker-compose logs -f app

# Database migrations (if using migrate tool)
migrate-install:
	@echo "Installing migrate tool..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-create:
	@echo "Creating new migration..."
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	@echo "Running migrations up..."
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	@echo "Running migrations down..."
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Production deployment preparation
deploy-prep:
	@echo "Preparing for deployment..."
	@echo "1. Building optimized binary..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o $(APP_NAME) main.go
	@echo "2. Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "3. Running security scan..."
	docker run --rm -v $(PWD):/app securecodewarrior/docker-security-scan $(DOCKER_IMAGE)

# Push to registry (update with your registry)
docker-push:
	@echo "Pushing to Docker registry..."
	docker tag $(DOCKER_IMAGE) your-registry/$(APP_NAME):latest
	docker push your-registry/$(APP_NAME):latest

# Benchmark
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Profile
profile:
	@echo "Running with profiling..."
	go build -o bin/$(APP_NAME) main.go
	./bin/$(APP_NAME) &
	sleep 5
	go tool pprof http://localhost:8080/debug/pprof/profile

# Help
help:
	@echo "CodeBase Golang - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  setup         - Setup environment (.env file)"
	@echo "  run           - Run application locally"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  benchmark     - Run benchmarks"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt           - Format code"
	@echo "  lint          - Run linter"
	@echo "  security      - Run security check"
	@echo ""
	@echo "Build & Deploy:"
	@echo "  build         - Build application binary"
	@echo "  docker-build  - Build Docker image"
	@echo "  deploy-prep   - Prepare for production deployment"
	@echo "  docker-push   - Push image to registry"
	@echo ""
	@echo "Docker Compose:"
	@echo "  docker-compose-up    - Start all services"
	@echo "  docker-compose-down  - Stop all services"
	@echo "  docker-compose-nginx - Start with nginx proxy"
	@echo "  logs                 - View application logs"
	@echo ""
	@echo "Database:"
	@echo "  migrate-install - Install migrate tool"
	@echo "  migrate-create  - Create migration (use: make migrate-create name=migration_name)"
	@echo "  migrate-up      - Run migrations up"
	@echo "  migrate-down    - Run migrations down"
	@echo ""
	@echo "Utilities:"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  profile       - Run with profiling"
	@echo "  help          - Show this help"