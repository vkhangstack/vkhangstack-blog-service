.PHONY: help build run test clean docker-build docker-run docker-stop docker-clean compose-up compose-down compose-logs migrate migrate-create lint fmt vet

# Variables
APP_NAME := golang-hexagonal
IMAGE_NAME := $(APP_NAME):latest
DOCKER_COMPOSE := docker-compose
GO := go
MAIN_PATH := ./cmd

# Default target
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Development
build: ## Build the Go application
	@echo "Building application..."
	$(GO) build -o bin/$(APP_NAME) $(MAIN_PATH)

run: ## Run the application locally
	@echo "Running application..."
	$(GO) run $(MAIN_PATH)

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test -v ./...

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	$(GO) test -v ./internal/tests/unit/...

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GO) test -v ./internal/tests/integration/...

test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GO) test -bench=. -benchmem ./internal/tests/benchmark/...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality
lint: ## Run linter
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GO) fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	$(GO) mod tidy

# Docker - Single Container
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME) .

docker-run: ## Run application in Docker container
	@echo "Running Docker container..."
	docker run -d --name $(APP_NAME) -p 8000:8000 --env-file .env $(IMAGE_NAME)

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

docker-logs: ## Show Docker container logs
	docker logs -f $(APP_NAME)

docker-clean: ## Remove Docker image and container
	@echo "Cleaning Docker resources..."
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true
	docker rmi $(IMAGE_NAME) || true

# Docker Compose
compose-up: ## Start all services with docker-compose
	@echo "Starting services with docker-compose..."
	$(DOCKER_COMPOSE) up -d

compose-down: ## Stop all services with docker-compose
	@echo "Stopping services with docker-compose..."
	$(DOCKER_COMPOSE) down

compose-logs: ## Show docker-compose logs
	$(DOCKER_COMPOSE) logs -f

compose-ps: ## Show docker-compose services status
	$(DOCKER_COMPOSE) ps

compose-build: ## Build services with docker-compose
	@echo "Building services..."
	$(DOCKER_COMPOSE) build

compose-restart: ## Restart all services
	@echo "Restarting services..."
	$(DOCKER_COMPOSE) restart

compose-clean: ## Stop and remove all containers, volumes, and images
	@echo "Cleaning docker-compose resources..."
	$(DOCKER_COMPOSE) down -v --rmi all

# Database
migrate-up: ## Run database migrations
	@echo "Running database migrations..."
	@if [ -f "./internal/migrations/hex_arch_db_schema.sql" ]; then \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -f ./internal/migrations/hex_arch_db_schema.sql; \
	else \
		echo "Migration file not found"; \
	fi

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=<migration_name>"; exit 1; fi
	@mkdir -p ./internal/migrations
	$(eval TIMESTAMP := $(shell date +%Y%m%d%H%M%S))
	@touch ./internal/migrations/$(TIMESTAMP)_$(NAME).up.sql
	@touch ./internal/migrations/$(TIMESTAMP)_$(NAME).down.sql
	@echo "Created migration: ./internal/migrations/$(TIMESTAMP)_$(NAME).up.sql and ./internal/migrations/$(TIMESTAMP)_$(NAME).down.sql"
	@echo "Please edit the migration files to add your SQL statements."

migrate-down: ## Rollback database migrations
	@echo "Rolling back database migrations..."
	@echo "Manual rollback required"

db-connect: ## Connect to database
	@echo "Connecting to database..."
	PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME)

# Clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f dump.rdb

clean-all: clean docker-clean compose-clean ## Clean everything including Docker resources

# Full deployment workflow
deploy-local: docker-build compose-up ## Build and deploy locally with docker-compose
	@echo "Application deployed locally"
	@echo "Access at: http://localhost:8000"

# Development workflow
dev: deps fmt vet test ## Run development checks (format, vet, test)

# Production build
prod-build: deps test docker-build ## Production build (test + docker build)
	@echo "Production build complete"

# Quick start
quick-start: ## Quick start with docker-compose
	@echo "Quick starting application..."
	@if [ ! -f ".env" ]; then \
		echo "Creating .env from .env.docker..."; \
		cp .env.docker .env; \
	fi
	make docker-build
	make compose-up
	@echo ""
	@echo "Application started!"
	@echo "Access at: http://localhost:8000"
	@echo "View logs: make compose-logs"
