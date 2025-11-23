.PHONY: all version fmt lint test coverage benchmark air deps release clean docker-build docker-up docker-down docker-logs help

BINARY_NAME := $(shell basename $(PWD))
GIT_VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0")
VERSION ?= $(GIT_VERSION)
DOCKER_CR ?= $(shell basename $$(dirname $(PWD)))
DOCKER_IMAGE := ${DOCKER_CR}/$(BINARY_NAME):$(VERSION)

all: fmt lint coverage ## Run all tests and checks

version: ## Display current version
	@echo "Current version: $(VERSION)"

fmt: ## Format the code
	golangci-lint fmt

lint: ## Lint the code
	golangci-lint run --timeout=5m

test: ## Run tests
	go test -race -shuffle=on -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...

coverage: test ## Generate coverage
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

benchmark: ## Run benchmarks
	go test -run=^$$ -bench=. -benchmem ./... | tee benchmark.txt

air: ## Run development server
	@command -v air >/dev/null 2>&1 || { \
      echo "Please install air: go install github.com/cosmtrek/air@latest"; \
      exit 1; \
    }
	@echo "Starting development server with air..."
	@air

deps: ## Install dependencies
	go mod download

release: ## Create release
	goreleaser release --snapshot --clean

clean: ## Remove build artifacts
	rm -f coverage.* benchmark.txt
	rm -rf dist

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-up: ## Start Docker services
	@echo "Starting Docker services..."
	@docker compose up --build -d

docker-down: ## Stop Docker services
	@echo "Stopping Docker services..."
	@docker compose down -v

docker-logs: ## Show Docker logs
	@echo "Showing Docker logs..."
	@docker compose logs -f

help: ## Show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
