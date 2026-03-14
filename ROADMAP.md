# 🛡️ Brainless WAF — Development Roadmap & Workflow

> **Project:** Brainless Security — Open-Source Web Application Firewall  
> **Maintainer:** Brainless Security Community  
> **License:** Apache 2.0  
> **Repo:** `https://github.com/brainless-security/brainless-waf`

---

## 📋 Table of Contents

- [Project Phases Overview](#project-phases-overview)
- [Development Workflow](#development-workflow)
- [Phase 0 — Foundation (Now)](#phase-0--foundation-now)
- [Phase 1 — Core Engine (v0.1 – v0.5)](#phase-1--core-engine-v01--v05)
- [Phase 2 — Feature Complete (v0.6 – v0.9)](#phase-2--feature-complete-v06--v09)
- [Phase 3 — Stable Release (v1.0)](#phase-3--stable-release-v10)
- [Phase 4 — Post-Launch (v1.1 – v1.2)](#phase-4--post-launch-v11--v12)
- [Phase 5 — Next Generation (v2.0)](#phase-5--next-generation-v20)
- [Task Board](#task-board)
- [Definition of Done](#definition-of-done)
- [Branch & Release Strategy](#branch--release-strategy)
- [Team Roles](#team-roles)

---

## Project Phases Overview

```
Phase 0       Phase 1          Phase 2           Phase 3     Phase 4       Phase 5
Foundation → Core Engine → Feature Complete → v1.0 GA → Post-Launch → Next Gen
  (Now)      Q2 2025          Q3 2025          Q4 2025     Q1–Q2 2026     2026+
```

| Phase | Milestone | Target Date | Status |
|-------|-----------|-------------|--------|
| Phase 0 | Project bootstrap, repo setup, architecture decision | Now | 🟡 In Progress |
| Phase 1 | Working traffic interceptor + basic rule engine | Q2 2025 | ⬜ Planned |
| Phase 2 | Full OWASP coverage, dashboard, API | Q3 2025 | ⬜ Planned |
| Phase 3 | v1.0 GA — stable, tested, documented | Q4 2025 | ⬜ Planned |
| Phase 4 | Plugin system, GraphQL, HA clustering | Q1–Q2 2026 | ⬜ Planned |
| Phase 5 | ML detection, eBPF, service mesh | 2026+ | ⬜ Future |

---

## Development Workflow

### How We Work

All development follows a **GitHub Flow** model with short-lived feature branches.

```
main (always deployable)
 └── develop (integration branch)
      ├── feature/xxx   ← your work happens here
      ├── fix/xxx
      └── docs/xxx
```

### Step-by-Step Contribution Flow

```
1. Pick a task from the board (see Task Board below)
       ↓
2. Create a branch from develop
   git checkout develop && git pull
   git checkout -b feature/your-feature-name
       ↓
3. Write code + tests
   make test        # unit tests
   make lint        # linter checks
   make e2e         # end-to-end (if applicable)
       ↓
4. Open a Pull Request → develop
   - Fill in the PR template
   - Link the issue: "Closes #123"
   - Request at least 1 reviewer
       ↓
5. Review & CI checks pass
   - All tests green
   - golangci-lint / ruff / mypy pass
   - Coverage ≥ 80% on changed files
       ↓
6. Merge to develop (squash merge preferred)
       ↓
7. Periodic release cut: develop → main
   - Tagged with semantic version (v1.0.0)
   - Changelog auto-generated from commit messages
```

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat:     New feature
fix:      Bug fix
docs:     Documentation only
test:     Adding or updating tests
refactor: Code change, no feature or fix
perf:     Performance improvement
chore:    Tooling, deps, CI changes
ci:       Changes to CI/CD configuration
security: Security-related fix (use this for vulnerability patches)
```

**Examples:**
```
feat(rule-engine): add YARA rule parser support
fix(tls): handle expired certificate renewal race condition
security(core): patch path traversal bypass in URI normalizer
docs(api): add rate limiting endpoint examples
```

---

## Phase 0 — Foundation (Now)

> **Goal:** Lay the technical and organizational groundwork so every future contributor can hit the ground running.

### 0.1 Repository & Project Setup

- [ ] Initialize GitHub repository with branch protection rules on `main` and `develop`
- [ ] Set up GitHub Actions CI pipeline (lint → test → build → docker push)
- [ ] Create `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`
- [ ] Set up issue templates: bug report, feature request, security vulnerability
- [ ] Set up PR template with checklist
- [ ] Configure Dependabot for dependency updates (Go modules, npm, pip)
- [ ] Set up CodeQL scanning for automated security analysis
- [ ] Create Discord server with channels: #announcements, #dev, #help, #security
- [ ] Register domain `brainless-security.io`
- [ ] Set up project board in GitHub Projects (Kanban: Backlog → In Progress → Review → Done)

### 0.2 Architecture Decision Records (ADRs)

- [ ] ADR-001: Choose Go as the core engine language (vs. C/Rust/Java)
- [ ] ADR-002: Choose NGINX vs. custom TCP stack for traffic interception
- [ ] ADR-003: Choose PostgreSQL vs. SQLite for configuration storage
- [ ] ADR-004: Choose React vs. Vue for the dashboard
- [ ] ADR-005: ModSecurity rule compatibility strategy
- [ ] ADR-006: Choose gRPC vs. REST for internal component communication
- [ ] Store all ADRs in `/docs/adr/` directory

### 0.3 Development Environment

- [ ] Create `docker-compose.dev.yml` for local development (WAF + mock backend + DB)
- [ ] Write `Makefile` with targets: `build`, `test`, `lint`, `e2e`, `docker`, `clean`
- [ ] Create `.devcontainer/` configuration for VS Code Dev Containers
- [ ] Document local setup in `DEVELOPMENT.md` (should work in under 10 minutes)
- [ ] Set up `pre-commit` hooks: lint, format, secret scanning (truffleHog)

### 0.4 Monorepo Structure Decision

- [ ] Decide on repo structure and document it:

```
brainless-waf/
├── core/               # Go — Traffic interceptor + detection engine
├── management/         # Python/FastAPI — REST API + config management
├── dashboard/          # React/TypeScript — Web UI
├── rules/              # Built-in BRF rules + OWASP CRS
├── plugins/            # Official plugin directory
├── deploy/             # Helm chart, Docker Compose, Terraform
├── docs/               # Documentation source
│   └── adr/            # Architecture Decision Records
├── tests/              # Integration + E2E tests
└── scripts/            # Build, release, utility scripts
```

---

## Phase 1 — Core Engine (v0.1 – v0.5)

> **Goal:** A working WAF that can intercept HTTP traffic and apply rules. No dashboard. No fancy features. Just rock-solid fundamentals.

### v0.1 — Skeleton (2 weeks)

- [ ] Initialize Go module (`github.com/brainless-security/brainless-waf`)
- [ ] Implement basic HTTP reverse proxy (net/http or fasthttp benchmark first)
- [ ] Add structured logging (zerolog) with request ID tracking
- [ ] Write health check endpoint `/health`
- [ ] Docker image builds and runs successfully
- [ ] Basic YAML config parsing (`server.listen`, `server.upstream`)
- [ ] Unit tests for config parsing (≥80% coverage)
- [ ] CI pipeline passes on PR

### v0.2 — Request Parsing (2 weeks)

- [ ] Parse and normalize incoming HTTP requests:
  - [ ] URL decoding (single and double encoded)
  - [ ] Unicode normalization
  - [ ] Header normalization (canonicalize names, strip hop-by-hop)
  - [ ] Multipart body parsing
  - [ ] JSON body parsing
  - [ ] Form data parsing
- [ ] Build internal `Request` struct passed to all subsequent stages
- [ ] Add request parsing benchmarks (target: <0.5ms per request)
- [ ] Fuzz test the parser with `go-fuzz`

### v0.3 — Rule Engine MVP (3 weeks)

- [ ] Design and document Brainless Rule Format (BRF) spec
- [ ] Implement rule file parser (`.rules` files)
- [ ] Implement operators: `@rx` (regex), `@streq`, `@contains`, `@beginsWith`, `@endsWith`, `@gt`, `@lt`
- [ ] Implement variables: `REQUEST_URI`, `REQUEST_HEADERS`, `ARGS`, `REMOTE_ADDR`, `REQUEST_BODY`
- [ ] Implement actions: `deny`, `allow`, `log`, `pass`, `redirect`
- [ ] Implement phases: `phase:1` (request headers), `phase:2` (request body)
- [ ] Rule hot-reload without WAF restart
- [ ] Load and parse OWASP CRS 4.x rules (compatibility layer)
- [ ] Unit tests for each operator and action
- [ ] Benchmark: 10,000 rules evaluated in <2ms

### v0.4 — Anomaly Scoring (2 weeks)

- [ ] Implement transaction-scoped anomaly score accumulation
- [ ] Assign severity weights: CRITICAL=10, ERROR=5, WARNING=3, NOTICE=1
- [ ] Implement configurable `anomaly_threshold` (default: 10)
- [ ] Add `paranoia_level` (1–4) to control rule strictness
- [ ] Implement `setvar` and `expirevar` for stateful rule logic
- [ ] Per-IP and per-session counters using in-memory store (Redis optional)
- [ ] Integration tests: ensure scoring correctly blocks at threshold

### v0.5 — TLS + Rate Limiting (2 weeks)

- [ ] TLS termination with configurable cert/key paths
- [ ] ACME protocol integration (Let's Encrypt auto-renewal via `golang.org/x/crypto/acme`)
- [ ] HTTP/2 support (built-in with Go's `net/http`)
- [ ] Basic rate limiting: requests per second per IP using token bucket algorithm
- [ ] Rate limit headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After`
- [ ] IP allowlist and blocklist (static config + API-managed)
- [ ] Return proper 429 Too Many Requests with `Retry-After` header

**✅ Phase 1 Exit Criteria:**
- WAF proxies real traffic without errors under 1000 req/s load
- OWASP CRS rules load and block basic SQLi/XSS payloads
- Docker image is <50MB and starts in <2 seconds
- All tests pass, coverage ≥ 80%

---

## Phase 2 — Feature Complete (v0.6 – v0.9)

> **Goal:** Build everything needed for real-world production use: management API, dashboard, bot protection, and full OWASP coverage.

### v0.6 — Management API (3 weeks)

- [ ] Initialize FastAPI project in `/management`
- [ ] PostgreSQL schema design and Alembic migrations
- [ ] JWT authentication endpoint (`POST /api/v1/auth/token`)
- [ ] API key management (create, revoke, list)
- [ ] CRUD endpoints for rules (`GET/POST/PUT/DELETE /api/v1/rules`)
- [ ] CRUD endpoints for IP lists (`/api/v1/blocklist`, `/api/v1/allowlist`)
- [ ] Configuration management endpoints
- [ ] Event/log query endpoint with filtering and pagination
- [ ] gRPC bridge: Management API → Core engine (push config changes live)
- [ ] OpenAPI 3.1 schema auto-generated, served at `/api/v1/docs`
- [ ] Rate limiting on the API itself (prevent brute force)
- [ ] API integration tests (pytest + httpx)

### v0.7 — Dashboard (4 weeks)

- [ ] Initialize React + TypeScript + Vite project in `/dashboard`
- [ ] Authentication: Login page, JWT token storage, auto-refresh
- [ ] RBAC: Admin, Analyst, Read-only roles enforced in UI
- [ ] **Overview page:** requests/sec, block rate, top threats, latency p50/p95/p99
- [ ] **Events page:** paginated security event log with severity filter, search, export to CSV
- [ ] **Rules page:** list, create, edit, delete rules with syntax highlighting (CodeMirror)
- [ ] **Rule tester:** paste a raw HTTP request and see which rules trigger
- [ ] **IP Management page:** blocklist/allowlist management with CIDR support
- [ ] **Settings page:** WAF mode, paranoia level, rate limit config
- [ ] **User Management page:** create users, assign roles (Admin only)
- [ ] WebSocket for real-time traffic feed on Overview page
- [ ] Responsive design (works on tablet, not required on mobile)
- [ ] Dark mode support
- [ ] E2E tests with Playwright for critical paths (login, create rule, block event)

### v0.8 — Advanced Detection (3 weeks)

- [ ] **Bot detection:**
  - [ ] User-agent fingerprinting database (known scanners: Nikto, sqlmap, Nuclei, etc.)
  - [ ] Headless browser detection via JS challenge (configurable)
  - [ ] Behavioral analysis: request rate, path patterns, header anomalies
  - [ ] TLS fingerprinting (JA3/JA4 hash support)
- [ ] **Response scanning:**
  - [ ] Phase 4 (response headers) + Phase 5 (response body) support
  - [ ] PII leak detection in responses (credit card patterns, SSN, email)
  - [ ] Error message suppression (stack traces, DB errors)
- [ ] **SSRF protection:** block requests to RFC1918, loopback, metadata endpoints
- [ ] **Path traversal:** enhanced `../` and null byte detection
- [ ] **File upload scanning:** MIME type validation, magic byte verification
- [ ] Virtual patching command: `brainless-ctl vpatch apply CVE-XXXX-XXXX`

### v0.9 — Kubernetes & Observability (3 weeks)

- [ ] **Helm chart** (`/deploy/helm/brainless-waf`):
  - [ ] `values.yaml` with all configurable parameters documented
  - [ ] HorizontalPodAutoscaler (HPA) configuration
  - [ ] PodDisruptionBudget for HA deployments
  - [ ] ServiceMonitor for Prometheus scraping
  - [ ] Ingress resource support
- [ ] **Kubernetes Ingress Controller mode** (annotation-driven WAF policies)
- [ ] **Metrics export (Prometheus):**
  - [ ] `brainless_requests_total` (labels: status, method, upstream)
  - [ ] `brainless_blocked_total` (labels: rule_id, severity, attack_type)
  - [ ] `brainless_latency_seconds` (histogram: p50, p95, p99)
  - [ ] `brainless_rule_evaluations_total`
- [ ] **Grafana dashboard JSON** (importable, covers all key metrics)
- [ ] **SIEM export:** Elasticsearch, Splunk HEC, Syslog (RFC 5424)
- [ ] **Distributed tracing:** OpenTelemetry integration (trace ID propagation)

**✅ Phase 2 Exit Criteria:**
- Full OWASP Top 10 blocked in automated test suite
- Dashboard fully functional end-to-end
- Helm chart deploys successfully on k3s and EKS
- Performance: <5ms added latency at p99 under 5000 req/s

---

## Phase 3 — Stable Release (v1.0)

> **Goal:** Production-ready. Every feature is tested, documented, and secure.

### v1.0-rc1 — Release Candidate (3 weeks)

- [ ] **Security audit:**
  - [ ] Internal security review of WAF core (focus: bypass techniques)
  - [ ] Dependency vulnerability scan (govulncheck, safety, npm audit)
  - [ ] Docker image scan (Trivy)
  - [ ] Invite community to responsible disclosure testing period
- [ ] **Performance hardening:**
  - [ ] Profile with `pprof` and eliminate top 3 bottlenecks
  - [ ] Rule engine: pre-compile all regex patterns at startup
  - [ ] Connection pool tuning for upstream proxy
  - [ ] Load test report: 10,000 req/s sustained, p99 <5ms
- [ ] **Documentation complete:**
  - [ ] `README.md` — project overview, quick start, badges
  - [ ] `docs/installation.md` — all deployment methods
  - [ ] `docs/configuration.md` — every config key documented
  - [ ] `docs/rules.md` — full BRF spec with examples
  - [ ] `docs/api.md` — all endpoints with request/response examples
  - [ ] `docs/troubleshooting.md` — common issues and fixes
  - [ ] `CHANGELOG.md` — full history from v0.1
  - [ ] All ADRs written and reviewed
- [ ] Migration guide from ModSecurity to Brainless WAF
- [ ] OWASP CRS import tool: `brainless-ctl import-crs /path/to/crs/`

### v1.0 GA — General Availability

- [ ] Tag `v1.0.0` on `main`
- [ ] Publish Docker image to Docker Hub (`brainlesssecurity/brainless-waf:1.0.0`, `:latest`)
- [ ] Publish Helm chart to `charts.brainless-security.io`
- [ ] GitHub Release with:
  - [ ] Binary downloads for: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64
  - [ ] SHA256 checksums file
  - [ ] Signed release (cosign)
- [ ] Announcement: GitHub Discussions, Discord, Reddit (r/netsec, r/selfhosted), Hacker News
- [ ] Submit to Awesome WAF and Awesome Security lists

---

## Phase 4 — Post-Launch (v1.1 – v1.2)

### v1.1 Target: Q1 2026

- [ ] **Plugin system:**
  - [ ] Define plugin API spec (`/plugins/PLUGIN_API.md`)
  - [ ] Plugin loader: discover and load `.so` files from plugin directory
  - [ ] Official plugins: GeoIP blocking, Cloudflare IP forwarding, reCAPTCHA
  - [ ] Plugin registry page on website
- [ ] **GraphQL attack detection:**
  - [ ] Parse GraphQL queries from request body
  - [ ] Detect introspection abuse, deeply nested queries, batch attacks
  - [ ] Configurable query depth and complexity limits
- [ ] **WebSocket inspection:**
  - [ ] Intercept WS upgrade handshake
  - [ ] Apply rules to WS message payloads
- [ ] **Terraform provider** (`terraform-provider-brainless`)
  - [ ] Resources: `brainless_rule`, `brainless_blocklist`, `brainless_site`
  - [ ] Data sources: `brainless_events`, `brainless_stats`
  - [ ] Publish to Terraform Registry
- [ ] **API key scoping:** restrict API keys by IP, endpoint prefix, or HTTP method

### v1.2 Target: Q2 2026

- [ ] **Active/Active HA clustering:**
  - [ ] Raft-based leader election for config distribution
  - [ ] Shared Redis for rate limit counters across nodes
  - [ ] Gossip protocol for rule synchronization
  - [ ] Cluster health dashboard widget
- [ ] **HashiCorp Vault integration:**
  - [ ] Fetch TLS certs from Vault PKI secrets engine
  - [ ] Fetch database credentials from Vault dynamic secrets
  - [ ] Auto-rotate credentials on lease expiry
- [ ] **Multi-tenancy:**
  - [ ] Site/tenant isolation (separate rule sets per virtual host)
  - [ ] Per-tenant rate limits and quotas
  - [ ] Tenant admin role (manage own rules, cannot see other tenants)
- [ ] **Automated rule update channel:**
  - [ ] Signed rule feed from `rules.brainless-security.io`
  - [ ] Community-submitted rules with review process
  - [ ] One-click update with rollback support

---

## Phase 5 — Next Generation (v2.0)

> Long-horizon items. Community input will shape prioritization.

- [ ] **eBPF-based traffic capture:** kernel-level packet inspection for near-zero overhead
- [ ] **ML anomaly detection model:**
  - [ ] Train on labeled attack traffic dataset (open dataset: CICIDS2017)
  - [ ] ONNX model format for cross-platform inference
  - [ ] Target: <5ms inference, <100MB model size
  - [ ] Online learning mode: adapt to traffic patterns
- [ ] **Service mesh integration:** Istio/Envoy filter support
- [ ] **Distributed attack correlation:** share threat intel between independent deployments (opt-in)
- [ ] **AI-assisted rule generation:** describe an attack in plain English, get a BRF rule
- [ ] **FIPS 140-2 compliance mode:** use BoringCrypto build tags
- [ ] **IPv6 full support** (currently partial)

---

## Task Board

### How to Use This Board

Copy tasks into GitHub Issues with the appropriate labels. Each task maps to one issue.

**Labels:**
- `phase/0`, `phase/1`, `phase/2`, `phase/3`, `phase/4` — which phase
- `component/core`, `component/api`, `component/dashboard`, `component/rules`, `component/deploy`
- `priority/critical`, `priority/high`, `priority/medium`, `priority/low`
- `type/feature`, `type/fix`, `type/docs`, `type/test`, `type/security`, `type/perf`
- `good-first-issue` — suitable for new contributors
- `help-wanted` — extra eyes needed

### Current Sprint (Phase 0 Focus)

| # | Task | Owner | Status | Priority |
|---|------|-------|--------|----------|
| 1 | Initialize GitHub repo with branch protection | — | ⬜ Todo | 🔴 Critical |
| 2 | Set up GitHub Actions CI (lint + test + build) | — | ⬜ Todo | 🔴 Critical |
| 3 | Write CONTRIBUTING.md | — | ⬜ Todo | 🔴 Critical |
| 4 | Create issue and PR templates | — | ⬜ Todo | 🟠 High |
| 5 | Write ADR-001: Core language choice | — | ⬜ Todo | 🟠 High |
| 6 | Set up docker-compose.dev.yml | — | ⬜ Todo | 🟠 High |
| 7 | Write Makefile with standard targets | — | ⬜ Todo | 🟠 High |
| 8 | Set up CodeQL + Dependabot | — | ⬜ Todo | 🟡 Medium |
| 9 | Register Discord server | — | ⬜ Todo | 🟡 Medium |
| 10 | Document monorepo folder structure | — | ⬜ Todo | 🟡 Medium |

### Backlog — Phase 1 (Ready to pick up)

| # | Task | Component | Estimated Effort |
|---|------|-----------|-----------------|
| 11 | Basic HTTP reverse proxy | core | 3 days |
| 12 | Structured logging with zerolog | core | 1 day |
| 13 | YAML config parser | core | 2 days |
| 14 | URL/Unicode normalization | core | 3 days |
| 15 | Multipart + JSON body parser | core | 3 days |
| 16 | BRF rule file parser | rules | 4 days |
| 17 | Regex operator (@rx) | rules | 2 days |
| 18 | String operators (@streq, @contains, etc.) | rules | 1 day |
| 19 | Deny / allow / log actions | rules | 2 days |
| 20 | OWASP CRS compatibility layer | rules | 5 days |
| 21 | Anomaly scoring engine | core | 3 days |
| 22 | TLS termination | core | 2 days |
| 23 | ACME / Let's Encrypt integration | core | 3 days |
| 24 | Token bucket rate limiter | core | 2 days |
| 25 | IP blocklist / allowlist | core | 1 day |

---

## Definition of Done

A task is **Done** only when ALL of the following are true:

```
✅ Code is written and self-reviewed
✅ Unit tests written, coverage ≥ 80% on changed code
✅ Linter passes (golangci-lint / ruff / mypy)
✅ PR is reviewed and approved by at least 1 maintainer
✅ CI pipeline is fully green (no skipped steps)
✅ Relevant documentation is updated (inline comments + /docs if applicable)
✅ CHANGELOG.md entry added (for user-visible changes)
✅ No new security vulnerabilities introduced (CodeQL + govulncheck clean)
✅ Merged to develop (or main for hotfixes)
```

For **security-related tasks**, add:
```
✅ Reviewed by at least 2 maintainers
✅ Tested against known bypass techniques
✅ Added to regression test suite
```

---

## Branch & Release Strategy

### Branch Names

```
feature/short-description       # New features
fix/short-description           # Bug fixes
security/cve-or-description     # Security patches
docs/short-description          # Docs-only changes
refactor/short-description      # Refactoring
release/v1.0.0                  # Release preparation
hotfix/v1.0.1                   # Urgent production fix
```

### Versioning

We use **Semantic Versioning** (`MAJOR.MINOR.PATCH`):

| Change type | Version bump | Example |
|-------------|-------------|---------|
| Breaking change | MAJOR | `1.0.0 → 2.0.0` |
| New feature (backward compatible) | MINOR | `1.0.0 → 1.1.0` |
| Bug fix or patch | PATCH | `1.0.0 → 1.0.1` |
| Security fix | PATCH (expedited) | `1.0.0 → 1.0.1` |

### Release Cadence

- **Patch releases** — as needed (security fixes get same-day release)
- **Minor releases** — every 6–8 weeks
- **Major releases** — annually or on major breaking changes

### Release Checklist

```
[ ] All planned issues for this milestone are closed or deferred
[ ] CHANGELOG.md updated with full list of changes
[ ] Version bumped in: version.go, package.json, Chart.yaml, pyproject.toml
[ ] docker-compose.yml uses the new tag (not :latest)
[ ] Release branch created: release/vX.Y.Z
[ ] Final CI run passes on release branch
[ ] Docker images built and pushed (amd64 + arm64)
[ ] GitHub Release created with changelog + binary assets + checksums
[ ] Helm chart version bumped and published
[ ] Announcement posted to Discord #announcements
```

---

## Team Roles

| Role | Responsibilities | How to become one |
|------|-----------------|-------------------|
| **Core Maintainer** | Merge PRs, triage issues, cut releases, set direction | 6+ months of consistent contributions + vote by existing maintainers |
| **Component Owner** | Deep expertise in one component (core/api/dashboard/rules), reviews PRs for that area | 3+ months + nominated by a Core Maintainer |
| **Security Researcher** | Review security-sensitive PRs, manage CVE disclosures | Invitation only, after verifying background |
| **Contributor** | Submit PRs, report bugs, answer community questions | Open to all — just submit a PR! |
| **Community Moderator** | Manage Discord, enforce Code of Conduct | Apply in Discord #meta channel |

---

## Quick Reference

### Key Commands

```bash
# Local dev setup
make dev                  # Start all services in watch mode

# Testing
make test                 # Unit tests
make test-integration     # Integration tests (needs Docker)
make test-e2e             # End-to-end tests (needs full stack)
make test-coverage        # Coverage report → coverage.html

# Code quality
make lint                 # All linters
make fmt                  # Auto-format all code
make security-scan        # govulncheck + trivy + trufflehog

# Building
make build                # Build all binaries
make docker               # Build Docker image
make docker-push          # Push to registry (CI only)

# Release
make changelog            # Generate changelog since last tag
make release-notes        # Draft release notes
```

### Useful Links

| Resource | URL |
|----------|-----|
| GitHub Repository | `https://github.com/brainless-security/brainless-waf` |
| Issue Tracker | `https://github.com/brainless-security/brainless-waf/issues` |
| Project Board | `https://github.com/orgs/brainless-security/projects/1` |
| Documentation | `https://docs.brainless-security.io` |
| Discord Community | `https://discord.gg/brainless-waf` |
| Security Issues | `security@brainless-security.io` |
| OWASP CRS | `https://github.com/coreruleset/coreruleset` |
| Helm Charts | `https://charts.brainless-security.io` |

---

*Last updated: 2025 — maintained by the Brainless Security core team.*  
*To propose changes to this roadmap, open a Discussion in GitHub.*
