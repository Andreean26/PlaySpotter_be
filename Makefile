.PHONY: up down logs migrate test build clean help

help: ## Show this help
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

up: ## Start all services
	docker compose up -d

down: ## Stop and remove all services
	docker compose down

logs: ## Show API logs
	docker compose logs -f api

migrate: ## Run database migrations
	@echo "Running migrations..."
	docker compose exec -T db psql -U postgres -d playspotter < migrations/0001_init.sql
	@echo "Migrations completed successfully!"

test: ## Run tests
	go test ./... -v

build: ## Build the application locally
	go build -o bin/playspotter ./cmd/api

clean: ## Clean up containers and volumes
	docker compose down -v

restart: down up ## Restart all services

swagger: ## Generate Swagger documentation
	swag init -g cmd/api/main.go -o docs

dev: ## Run in development mode
	go run cmd/api/main.go
