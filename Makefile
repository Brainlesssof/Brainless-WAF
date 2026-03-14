.PHONY: build test lint docker dev clean

# Variables
MANAGEMENT_DIR=management
CORE_DIR=core
DOCKER_COMPOSE=docker-compose.dev.yml

# Default target
all: build

# Build all components
build:
	@echo "Building all components..."
	cd $(MANAGEMENT_DIR) && pip install .
	# cd $(CORE_DIR) && go build -o brainless-waf

# Run tests
test:
	@echo "Running tests..."
	cd $(MANAGEMENT_DIR) && pytest
	# cd $(CORE_DIR) && go test ./...

# Linting
lint:
	@echo "Running linters..."
	cd $(MANAGEMENT_DIR) && flake8 .
	# cd $(CORE_DIR) && golangci-lint run

# Docker
docker:
	@echo "Building Docker images..."
	docker compose -f $(DOCKER_COMPOSE) build

# Local Development
dev:
	@echo "Starting local development environment..."
	docker compose -f $(DOCKER_COMPOSE) up -d

# Cleanup
clean:
	@echo "Cleaning up..."
	docker compose -f $(DOCKER_COMPOSE) down -v
	rm -rf $(MANAGEMENT_DIR)/build/
	rm -rf $(MANAGEMENT_DIR)/dist/
	rm -rf $(CORE_DIR)/brainless-waf
