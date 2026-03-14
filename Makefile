.PHONY: build test lint docker dev clean

# Variables
MANAGEMENT_DIR=management
CORE_DIR=core
DASHBOARD_DIR=dashboard
DOCKER_COMPOSE=docker-compose.dev.yml

# Default target
all: build

# Build all components
build:
	@echo "Building all components..."
	cd $(MANAGEMENT_DIR) && pip install .
	cd $(DASHBOARD_DIR) && npm ci && npm run build
	# cd $(CORE_DIR) && go build -o brainless-waf

# Run tests
test:
	@echo "Running tests..."
	cd $(MANAGEMENT_DIR) && pytest
	cd $(DASHBOARD_DIR) && npm ci && npm run test -- --run
	# cd $(CORE_DIR) && go test ./...

# Linting
lint:
	@echo "Running linters..."
	cd $(MANAGEMENT_DIR) && flake8 .
	cd $(DASHBOARD_DIR) && npm ci && npm run lint && npm run type-check
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
