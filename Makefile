APP_NAME := go-scaffold
VERSION := 1.0.0
BUILD_DIR := build

.PHONY: all build clean test test-cover run run-dev api-docs help

all: clean build

build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o $(BUILD_DIR)/$(APP_NAME) -v ./cmd/scaffold

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

test: ## Run all tests
	@echo "Running all tests..."
	go test -v ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	go test -v ./test/unit/...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	go test -v ./test/integration/...

test-functional: ## Run functional tests only
	@echo "Running functional tests..."
	go test -v ./test/functional/...

test-cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -cover -v ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	go run ./cmd/scaffold

run-dev: ## Run the application in development mode (requires air)
	@echo "Running $(APP_NAME) in development mode..."
	air

api-docs: ## Generate API documentation
	@echo "API documentation is available at http://localhost:8081/api-docs"
	@echo "OpenAPI spec is available at http://localhost:8081/api/docs"

help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
