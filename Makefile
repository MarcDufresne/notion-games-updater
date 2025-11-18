.PHONY: build run test clean docker-build docker-run install help

# Variables
BINARY_NAME=updater
DOCKER_IMAGE=notion-games-updater
DOCKER_TAG=go

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

install: ## Download Go dependencies
	go mod download
	go mod tidy

build: ## Build the binary
	go build -o $(BINARY_NAME) cmd/updater/main.go

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux cmd/updater/main.go

run: ## Run the updater once
	go run cmd/updater/main.go

run-loop: ## Run the updater in loop mode
	go run cmd/updater/main.go -run-forever -interval 15

test: ## Run tests (when implemented)
	go test -v ./...

clean: ## Remove built binaries
	rm -f $(BINARY_NAME) $(BINARY_NAME)-linux

docker-build: ## Build Docker image
	docker build -f Dockerfile -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-run-once: ## Run Docker container once
	docker run --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG) ./updater

fmt: ## Format Go code
	go fmt ./cmd/... ./internal/...

vet: ## Run go vet
	go vet ./cmd/... ./internal/...

lint: fmt vet ## Run formatters and linters
