.PHONY: help build build-dev build-prod run-dev run-prod frontend-install frontend-build migrate docker-build docker-run clean test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: build-prod ## Build production binary (alias for build-prod)

build-dev: ## Build development binary (serves frontend from filesystem)
	@echo "Building development binary..."
	go build -tags dev -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate
	@echo "Done! Run with: ./bin/server"

build-prod: frontend-build ## Build production binary (with embedded frontend)
	@echo "Building production binary..."
	mkdir -p cmd/server/frontend
	cp -r frontend/dist cmd/server/frontend/
	go build -o bin/server ./cmd/server
	go build -o bin/migrate ./cmd/migrate
	rm -rf cmd/server/frontend
	@echo "Done! Run with: ./bin/server"

run-dev: ## Run development server (build + run)
	@echo "Starting development server..."
	go run -tags dev ./cmd/server/main.go

run-prod: build-prod ## Run production server
	@echo "Starting production server..."
	./bin/server

frontend-install: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

frontend-build: frontend-install ## Build frontend
	@echo "Building frontend..."
	cd frontend && npm run build

frontend-dev: ## Run frontend development server
	@echo "Starting frontend development server..."
	cd frontend && npm run dev

migrate: ## Build migration tool
	@echo "Building migration tool..."
	go build -o bin/migrate ./cmd/migrate
	@echo "Done! Run with: ./bin/migrate --help"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t game-tracker:latest .

docker-run: ## Run Docker container
	@echo "Starting Docker container..."
	docker-compose up -d

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	docker-compose down

docker-logs: ## View Docker container logs
	docker-compose logs -f

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/
	rm -rf cmd/server/frontend/
	@echo "Done!"

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

deps: ## Download Go dependencies
	@echo "Downloading Go dependencies..."
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
