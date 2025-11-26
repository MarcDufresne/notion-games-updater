.PHONY: help install build-frontend build-backend build run dev-setup dev-backend dev-frontend dev clean docker-build docker-run lint

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
	@echo "Copying frontend for embed..."
	@rm -rf cmd/server/frontend
	@cp -r frontend cmd/server/frontend
	@cd cmd/server && go build -o ../../server .
	@rm -rf cmd/server/frontend
	@echo "Build complete!"

build: build-backend ## Build everything

run: build ## Build and run the server
	./server

dev-setup: build-frontend ## Setup development environment (run once)
	@echo "Setting up development environment..."
	@rm -rf cmd/server/frontend
	@cp -r frontend cmd/server/frontend
	@echo "Development environment ready!"
	@echo "Run 'make dev-backend' to start the backend server"
	@echo "Run 'make dev-frontend' in another terminal for frontend hot reload"

dev-backend: ## Run backend in development mode (run dev-setup first)
	@if [ ! -d "cmd/server/frontend/dist" ]; then \
		echo "Error: Frontend not found. Run 'make dev-setup' first."; \
		exit 1; \
	fi
	@echo "Starting backend server..."
	@go run cmd/server/main.go

dev-frontend: ## Run frontend dev server
	cd frontend && npm run dev

dev: dev-setup ## Setup and show development instructions
	@echo ""
	@echo "=========================================="
	@echo "Development Environment Ready!"
	@echo "=========================================="
	@echo ""
	@echo "To start developing:"
	@echo "  1. Run backend:  make dev-backend"
	@echo "  2. Run frontend: make dev-frontend (in another terminal)"
	@echo ""
	@echo "The frontend dev server (Vite) provides hot reload."
	@echo "The backend serves the API at http://localhost:8080"
	@echo ""
	@echo "Note: Backend uses embedded production build."
	@echo "For frontend changes, use the Vite dev server at http://localhost:5173"
	@echo "=========================================="
	@echo ""

clean: ## Clean build artifacts
	rm -f server migrate
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf cmd/server/frontend

docker-build: ## Build Docker image
	docker build -t game-tracker .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env.docker -v ./firebase_key.json:/firebase_key.json game-tracker
