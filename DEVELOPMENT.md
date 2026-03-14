# Development Guide

Everything you need to set up a local development environment for Brainless WAF, understand the codebase, run tests, and ship changes confidently.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Repository Structure](#repository-structure)
- [Local Setup](#local-setup)
- [Running the Stack](#running-the-stack)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Debugging](#debugging)
- [Environment Variables](#environment-variables)
- [Database](#database)
- [Working on the Core (Go)](#working-on-the-core-go)
- [Working on the API (Python)](#working-on-the-api-python)
- [Working on the Dashboard (TypeScript)](#working-on-the-dashboard-typescript)
- [Working on Rules](#working-on-rules)
- [CI/CD Pipeline](#cicd-pipeline)
- [Makefile Reference](#makefile-reference)

---

## Prerequisites

Install these tools before starting:

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.22+ | https://go.dev/dl/ |
| Node.js | 20+ | https://nodejs.org or `nvm install 20` |
| Python | 3.11+ | https://python.org or `pyenv install 3.11` |
| Docker | 24.0+ | https://docs.docker.com/get-docker/ |
| Docker Compose | v2+ | Bundled with Docker Desktop |
| Make | any | Pre-installed on Linux/macOS |
| git | 2.40+ | https://git-scm.com |

Optional but recommended:
- `golangci-lint` 1.57+ — Go linter
- `pre-commit` — git hook manager
- VS Code with the Go, Python, and ESLint extensions

---

## Repository Structure

```
brainless-waf/
│
├── core/                    # Go — WAF engine
│   ├── cmd/                 # Entry points (brainless-waf, brainless-ctl)
│   ├── internal/
│   │   ├── proxy/           # HTTP reverse proxy
│   │   ├── parser/          # Request/response parsing + normalization
│   │   ├── rules/           # Rule engine (BRF parser, operators, actions)
│   │   ├── detection/       # Detection pipeline + anomaly scoring
│   │   ├── ratelimit/       # Token bucket rate limiter
│   │   ├── tls/             # TLS termination + ACME
│   │   ├── metrics/         # Prometheus metrics
│   │   └── config/          # Config loading and validation
│   ├── pkg/                 # Public Go packages (re-usable by plugins)
│   └── tests/               # Integration + fuzz tests
│
├── management/              # Python — REST API + config management
│   ├── app/
│   │   ├── api/             # FastAPI routers
│   │   ├── models/          # SQLAlchemy models
│   │   ├── schemas/         # Pydantic request/response schemas
│   │   ├── services/        # Business logic
│   │   └── grpc/            # gRPC client → core engine
│   ├── migrations/          # Alembic database migrations
│   └── tests/               # pytest test suite
│
├── dashboard/               # TypeScript — React web UI
│   ├── src/
│   │   ├── pages/           # Page components
│   │   ├── components/      # Shared UI components
│   │   ├── hooks/           # Custom React hooks
│   │   ├── api/             # API client (auto-generated from OpenAPI)
│   │   └── store/           # Global state (React context)
│   └── e2e/                 # Playwright E2E tests
│
├── rules/                   # Detection rules
│   ├── crs/                 # OWASP Core Rule Set (submodule)
│   ├── brainless/           # Built-in Brainless rules
│   ├── community/           # Community-contributed rules
│   └── tests/               # Rule test cases
│
├── deploy/                  # Deployment configuration
│   ├── helm/                # Helm chart
│   ├── docker/              # Dockerfiles
│   └── terraform/           # Terraform modules
│
├── docs/                    # Documentation
│   ├── adr/                 # Architecture Decision Records
│   └── assets/              # Images, diagrams
│
├── tests/                   # Cross-component integration tests
│   ├── integration/
│   ├── e2e/
│   └── fixtures/            # Test payloads (attacks + safe)
│
├── scripts/                 # Build, release, utility scripts
├── config/                  # Example configurations
│   ├── default.yaml
│   └── production.yaml.example
│
├── Makefile
├── docker-compose.yml       # Production-style compose
├── docker-compose.dev.yml   # Development compose (hot reload, debug ports)
└── .github/                 # CI workflows, issue templates, CODEOWNERS
```

---

## Local Setup

### 1. Clone

```bash
git clone https://github.com/YOUR_USERNAME/brainless-waf.git
cd brainless-waf
git remote add upstream https://github.com/brainless-security/brainless-waf.git
```

### 2. Install all dependencies

```bash
make deps
```

This installs Go modules, Python packages (in a virtualenv), and npm packages.

### 3. Install pre-commit hooks

```bash
pip install pre-commit
pre-commit install
```

### 4. Copy and configure dev environment

```bash
cp config/default.yaml config/local.yaml
cp .env.example .env
```

Edit `.env` — the defaults work for local development without changes.

### 5. Start the full stack

```bash
make dev
```

This starts all services via `docker-compose.dev.yml`:

| Service | URL | Notes |
|---------|-----|-------|
| WAF core | `http://localhost:9090` | Proxies to mock backend |
| Management API | `http://localhost:8000` | API docs at `/api/v1/docs` |
| Dashboard | `http://localhost:5173` | Vite HMR — changes reload instantly |
| Mock backend | `http://localhost:3000` | Simple HTTP echo server |
| PostgreSQL | `localhost:5432` | DB: `brainless`, user: `brainless`, pass: `devpassword` |
| Redis | `localhost:6379` | Rate limit counters (optional in dev) |
| Adminer | `http://localhost:8888` | DB browser UI |

### 6. Verify everything works

```bash
make health-check
# Should output: All services healthy ✓
```

---

## Running the Stack

### Full stack (recommended for most development)

```bash
make dev           # Start with hot-reload
make dev-logs      # Tail logs from all services
make dev-stop      # Stop all services
make dev-clean     # Stop + delete volumes (fresh start)
```

### Individual components

If you're only working on one component, you can run it directly:

```bash
# Core engine only (requires Docker for postgres/redis)
cd core && go run ./cmd/brainless-waf --config ../config/local.yaml

# Management API only
cd management && uvicorn app.main:app --reload --port 8000

# Dashboard only
cd dashboard && npm run dev
```

---

## Development Workflow

```
1. Sync with upstream
   git fetch upstream && git merge upstream/develop

2. Create a branch
   git checkout -b feature/my-change

3. Make changes
   # Edit code, write tests

4. Run checks
   make test && make lint

5. Commit
   git commit -m "feat(rules): add WebSocket injection detection"

6. Push and open PR
   git push origin feature/my-change
   # Open PR on GitHub against develop
```

---

## Testing

### Run all tests

```bash
make test                  # Unit tests (all components)
make test-integration      # Integration tests (requires Docker)
make test-e2e              # End-to-end tests (requires full stack running)
make test-coverage         # Coverage report → coverage.html
```

### Go unit tests

```bash
cd core

# All tests
go test ./...

# Specific package
go test ./internal/rules/...

# Specific test
go test -run TestRuleEngine_SQLiDetection ./internal/rules/...

# With verbose output
go test -v ./...

# With race detector (always use for concurrent code)
go test -race ./...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

### Benchmarks

```bash
cd core
go test -bench=. -benchmem ./internal/...
go test -bench=BenchmarkRequestParser -benchmem ./internal/parser/...
```

### Fuzz tests

```bash
cd core
go test -fuzz=FuzzURIParser -fuzztime=60s ./internal/parser/...
go test -fuzz=FuzzRuleParser -fuzztime=60s ./internal/rules/...
```

### Python API tests

```bash
cd management
pytest                              # All tests
pytest tests/api/                   # API endpoint tests only
pytest -v -k "test_auth"            # Tests matching a keyword
pytest --cov=app --cov-report=html  # Coverage report
```

### Dashboard tests

```bash
cd dashboard
npm run test                # Vitest unit tests
npm run test:coverage       # With coverage
npm run test:e2e            # Playwright E2E (requires full stack)
```

### Rule tests

```bash
# Test all rules
make test-rules

# Test a specific rule file
make test-rules RULE=rules/brainless/sqli.rules

# Test with a specific payload
brainless-ctl rule-test --rule rules/brainless/sqli.rules \
  --request tests/fixtures/attacks/sqli/union_select.http
```

---

## Debugging

### WAF Core (Go)

```bash
# Enable debug logging
export BRAINLESS_LOG_LEVEL=debug
cd core && go run ./cmd/brainless-waf --config ../config/local.yaml

# Attach Delve debugger
dlv debug ./cmd/brainless-waf -- --config ../config/local.yaml

# CPU profiling
go test -cpuprofile=cpu.prof -bench=BenchmarkDetectionPipeline ./...
go tool pprof -http=:6060 cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=. ./...
go tool pprof -http=:6060 mem.prof
```

In `docker-compose.dev.yml`, the WAF core exposes `localhost:6060` for pprof automatically.

### Management API (Python)

```bash
# Debug mode with auto-reload
cd management && uvicorn app.main:app --reload --log-level debug

# Open interactive Python shell with app context
cd management && python -c "from app.main import app; import IPython; IPython.embed()"
```

### Dashboard (TypeScript)

- React DevTools browser extension is recommended
- Redux/context state is inspectable via browser DevTools
- Vite dev server proxies `/api` to `localhost:8000` automatically

### Useful Docker commands

```bash
# Enter a running container
docker compose -f docker-compose.dev.yml exec waf-core sh

# View logs for one service
docker compose -f docker-compose.dev.yml logs -f waf-core

# Restart one service without stopping others
docker compose -f docker-compose.dev.yml restart management-api
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `BRAINLESS_CONFIG` | `config/local.yaml` | Path to config file |
| `BRAINLESS_LOG_LEVEL` | `info` | Log level: debug, info, warn, error |
| `BRAINLESS_LOG_FORMAT` | `json` | Log format: json, text |
| `DATABASE_URL` | `postgresql://brainless:devpassword@localhost:5432/brainless` | PostgreSQL connection string |
| `REDIS_URL` | `redis://localhost:6379/0` | Redis connection string (optional) |
| `JWT_SECRET` | *(required)* | Secret key for JWT signing — min 32 chars |
| `API_KEY_ENCRYPTION_KEY` | *(required)* | AES-256 key for encrypting API keys at rest |
| `GRPC_CORE_ADDR` | `localhost:9091` | Address of WAF core gRPC server |
| `CORS_ORIGINS` | `http://localhost:5173` | Allowed CORS origins for dashboard |

In development, these are set via `.env` file (see `.env.example`).

---

## Database

### Migrations

```bash
cd management

# Apply all pending migrations
alembic upgrade head

# Create a new migration
alembic revision --autogenerate -m "add user preferences table"

# Downgrade one step
alembic downgrade -1

# Show migration history
alembic history --verbose
```

### Reset dev database

```bash
make dev-db-reset    # Drop, recreate, and seed with test data
```

### Seed data

The dev database seeds with:
- Admin user: `admin` / `changeme`
- Analyst user: `analyst` / `changeme`
- 5 sample security events
- OWASP CRS rules loaded

---

## Working on the Core (Go)

### Key packages to know

| Package | Purpose |
|---------|---------|
| `internal/proxy` | HTTP reverse proxy, request/response flow |
| `internal/parser` | Request normalization, body parsing |
| `internal/rules` | BRF parser, rule evaluation, variable extraction |
| `internal/detection` | Pipeline orchestration, anomaly scoring |
| `internal/ratelimit` | Token bucket implementation |
| `internal/config` | Config struct, YAML loading, validation |

### Adding a new operator

1. Add the operator function to `internal/rules/operators.go`
2. Register it in `internal/rules/operators_registry.go`
3. Add unit tests in `internal/rules/operators_test.go`
4. Document it in `docs/rules.md`

### Adding a new variable

1. Add extraction logic to `internal/rules/variables.go`
2. Register it in `internal/rules/variables_registry.go`
3. Add unit tests
4. Document in `docs/rules.md`

---

## Working on the API (Python)

### Adding a new endpoint

1. Create a new router file in `management/app/api/v1/`
2. Add the router to `management/app/api/v1/__init__.py`
3. Define request/response Pydantic schemas in `management/app/schemas/`
4. Add business logic in `management/app/services/`
5. Write tests in `management/tests/api/`
6. The OpenAPI docs update automatically

### Database models

All models live in `management/app/models/`. After adding a model:
```bash
cd management
alembic revision --autogenerate -m "add your_model table"
alembic upgrade head
```

---

## Working on the Dashboard (TypeScript)

### API client

The API client is auto-generated from the OpenAPI schema:
```bash
cd dashboard
npm run gen-api    # Regenerates src/api/client.ts from management API schema
```

Always regenerate after changing API schemas.

### Adding a new page

1. Create `src/pages/MyPage.tsx`
2. Add the route to `src/App.tsx`
3. Add navigation link to `src/components/Sidebar.tsx`
4. Add E2E test in `e2e/my-page.spec.ts`

---

## Working on Rules

### Rule ID ranges

| Range | Owner |
|-------|-------|
| 1000–9999 | OWASP CRS (imported, do not modify) |
| 10000–19999 | Brainless built-in rules |
| 20000–29999 | Reserved for future use |
| 50000–59999 | Community rules |

### Testing a rule before PR

```bash
# Test your rule against known payloads
make test-rules RULE=rules/community/my_rule.rules

# Test against all attack fixtures (no rule should block safe fixtures)
make test-rules-safe RULE=rules/community/my_rule.rules
```

---

## CI/CD Pipeline

GitHub Actions runs on every PR and push to `develop` and `main`:

```
PR opened / push
    │
    ├── lint.yml
    │     ├── golangci-lint (Go)
    │     ├── ruff + mypy (Python)
    │     └── eslint + tsc (TypeScript)
    │
    ├── test.yml
    │     ├── go test -race ./... (Go unit tests)
    │     ├── pytest (Python unit tests)
    │     ├── vitest (TypeScript unit tests)
    │     └── integration tests (Docker)
    │
    ├── security.yml
    │     ├── CodeQL analysis
    │     ├── govulncheck (Go)
    │     ├── safety (Python)
    │     └── npm audit (Node.js)
    │
    └── build.yml (on develop/main only)
          ├── docker build + push to registry
          └── helm chart lint + package
```

All checks must pass before a PR can be merged. Maintainers cannot bypass CI.

---

## Makefile Reference

```bash
# Setup
make deps              # Install all dependencies
make hooks             # Install pre-commit hooks

# Development
make dev               # Start dev stack (docker-compose.dev.yml)
make dev-stop          # Stop dev stack
make dev-clean         # Stop + remove volumes
make dev-logs          # Tail all service logs
make dev-db-reset      # Reset and reseed database
make health-check      # Verify all services are running

# Testing
make test              # All unit tests
make test-integration  # Integration tests
make test-e2e          # End-to-end tests
make test-rules        # Rule tests
make test-coverage     # Coverage report

# Code quality
make lint              # All linters
make fmt               # Auto-format all code
make security-scan     # govulncheck + trivy + trufflehog

# Building
make build             # Build all binaries (output: bin/)
make build-core        # Build WAF core only
make build-docker      # Build Docker image
make build-helm        # Package Helm chart

# Release
make changelog         # Generate changelog since last tag
make release-notes     # Draft release notes

# Utilities
make clean             # Remove build artifacts
make gen-api           # Regenerate dashboard API client from OpenAPI schema
make gen-grpc          # Regenerate gRPC code from .proto files
make deps-update       # Update all dependencies to latest compatible versions
```

---

## Getting Help

- **Discord `#dev`** — fastest, maintainers are usually online
- **GitHub Discussions** — for design questions and longer threads
- Check the [Troubleshooting Guide](docs/troubleshooting.md) first for common issues

When asking for help, include:
1. What you were trying to do
2. What you expected to happen
3. What actually happened (include error messages and logs)
4. Your environment: OS, Go/Python/Node version, Docker version
