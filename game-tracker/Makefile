.PHONY: help install build-frontend build-backend build run dev clean docker-build docker-run migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	go mod download
	cd frontend && npm install

build-frontend: ## Build the frontend
	cd frontend && npm run build

build-backend: build-frontend ## Build the backend (embeds frontend)
	go build -o server cmd/server/main.go

build: build-backend ## Build everything

run: build ## Build and run the server
	./server

dev-backend: ## Run backend in development mode
	go run cmd/server/main.go

dev-frontend: ## Run frontend dev server
	cd frontend && npm run dev

clean: ## Clean build artifacts
	rm -f server migrate
	rm -rf frontend/dist
	rm -rf frontend/node_modules

docker-build: ## Build Docker image
	docker build -t game-tracker .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env game-tracker

docker-compose-up: ## Start with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop docker-compose
	docker-compose down

migrate: ## Build and run migration tool
	go build -o migrate cmd/migrate/main.go
	./migrate --user-id=$(USER_ID)

test: ## Run tests
	go test -v ./...

lint: ## Run linters
	golangci-lint run
	cd frontend && npm run lint
