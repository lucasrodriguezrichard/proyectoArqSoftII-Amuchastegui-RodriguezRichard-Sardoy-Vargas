.PHONY: help build run stop clean test migrate seed logs docker-up docker-down

# Variables
APP_NAME=restaurant-system
DOCKER_COMPOSE=docker-compose
GO=go

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	$(GO) build -o bin/$(APP_NAME) ./cmd/api

run: ## Run the application locally
	@echo "Running application..."
	$(GO) run ./cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build files
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-up: ## Start all services with Docker Compose
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d
	@echo "Services started!"
	@echo "API: http://localhost:8080"
	@echo "Adminer: http://localhost:8081"

docker-down: ## Stop all Docker services
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-rebuild: ## Rebuild and restart Docker services
	@echo "Rebuilding Docker services..."
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f api

logs-db: ## Show database logs
	$(DOCKER_COMPOSE) logs -f db

stop: ## Stop all services
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) stop

restart: ## Restart all services
	@echo "Restarting services..."
	$(DOCKER_COMPOSE) restart

ps: ## Show running containers
	$(DOCKER_COMPOSE) ps

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

migrate: ## Run database migrations
	@echo "Running migrations..."
	$(GO) run ./cmd/migrate

fmt: ## Format Go code
	@echo "Formatting code..."
	$(GO) fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

dev: ## Start development environment
	@echo "Starting development environment..."
	$(MAKE) docker-up
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Development environment ready!"

prod: ## Start production environment
	@echo "Starting production environment..."
	$(DOCKER_COMPOSE) -f docker-compose.prod.yml up -d

health: ## Check API health
	@echo "Checking API health..."
	@curl -s http://localhost:8080/health | jq .

db-shell: ## Open PostgreSQL shell
	$(DOCKER_COMPOSE) exec db psql -U restaurant_user -d restaurant_db

db-backup: ## Backup database
	@echo "Backing up database..."
	$(DOCKER_COMPOSE) exec -T db pg_dump -U restaurant_user restaurant_db > backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Backup completed!"

db-restore: ## Restore database from backup (usage: make db-restore FILE=backup.sql)
	@echo "Restoring database..."
	$(DOCKER_COMPOSE) exec -T db psql -U restaurant_user restaurant_db < $(FILE)
	@echo "Restore completed!"