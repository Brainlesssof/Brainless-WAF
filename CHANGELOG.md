# Changelog

All notable changes to Brainless WAF are documented here.

This project follows [Semantic Versioning](https://semver.org/) and [Conventional Commits](https://www.conventionalcommits.org/).

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

---

## [Unreleased]

### Added
- Initial project structure and repository setup
- Architecture Decision Records (ADR-001 through ADR-006)
- Development environment configuration (`docker-compose.dev.yml`, `Makefile`)
- CI/CD pipeline (GitHub Actions: lint → test → build → docker push)
- Issue templates: bug report, feature request, security vulnerability
- Pull request template
- Community documents: README, CONTRIBUTING, CODE_OF_CONDUCT, SECURITY

### Changed
- Nothing yet

### Fixed
- Nothing yet

---

## [0.1.0] — TBD

> First functional release — basic HTTP reverse proxy with health check.

### Added
- HTTP/HTTPS reverse proxy core (Go, `net/http`)
- Structured logging with `zerolog` (JSON format, request ID tracking)
- YAML configuration parser (`server.listen`, `server.upstream`, `server.tls`)
- Health check endpoint (`GET /health`)
- Docker image (multi-stage build, non-root user, <50MB)
- Docker Compose for local deployment
- Basic unit test suite (≥80% coverage on core package)

---

## [0.2.0] — TBD

> Request parsing and normalization.

### Added
- URL decoding (single and double encoded)
- Unicode normalization (NFC)
- Header normalization (canonical names, hop-by-hop header stripping)
- Multipart form body parsing
- JSON body parsing
- URL-encoded form data parsing
- Internal `Request` struct passed through detection pipeline
- Parser fuzz tests (`go test -fuzz`)

### Performance
- Request parser benchmarks added: target <0.5ms per request

---

## [0.3.0] — TBD

> Rule engine MVP — Brainless Rule Format (BRF) + OWASP CRS compatibility.

### Added
- Brainless Rule Format (BRF) specification (`docs/rules.md`)
- Rule file parser (`.rules` file format)
- Operators: `@rx`, `@streq`, `@contains`, `@beginsWith`, `@endsWith`, `@gt`, `@lt`
- Variables: `REQUEST_URI`, `REQUEST_HEADERS`, `ARGS`, `REMOTE_ADDR`, `REQUEST_BODY`
- Actions: `deny`, `allow`, `log`, `pass`, `redirect`
- Phase support: `phase:1` (request headers), `phase:2` (request body)
- Rule hot-reload without WAF restart
- OWASP CRS 4.x compatibility layer (import existing ModSecurity rules)
- Rule evaluation benchmarks

### Performance
- 10,000 rules evaluated in <2ms (all regex pre-compiled at startup)

---

## [0.4.0] — TBD

> Anomaly scoring engine and paranoia levels.

### Added
- Transaction-scoped anomaly score accumulation
- Severity weights: CRITICAL=10, ERROR=5, WARNING=3, NOTICE=1
- Configurable `anomaly_threshold` (default: 10)
- `paranoia_level` setting (1–4) to control rule strictness
- `setvar` and `expirevar` actions for stateful rule counters
- Per-IP counters with optional Redis backend
- Integration tests for scoring thresholds

---

## [0.5.0] — TBD

> TLS termination, HTTP/2, and rate limiting.

### Added
- TLS 1.2/1.3 termination with configurable certificate paths
- ACME protocol (Let's Encrypt) with automatic renewal
- HTTP/2 server push support
- Token bucket rate limiter (configurable RPS and burst per IP)
- Rate limit response headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After`
- 429 Too Many Requests responses with proper `Retry-After` header
- Static IP allowlist and blocklist (config file + runtime API)

---

## [0.6.0] — TBD

> Management REST API.

### Added
- FastAPI management service (`/management`)
- PostgreSQL schema with Alembic migrations
- JWT authentication (`POST /api/v1/auth/token`)
- API key management (create, revoke, list, scope)
- Rule CRUD endpoints (`GET/POST/PUT/DELETE /api/v1/rules`)
- IP list management (`/api/v1/blocklist`, `/api/v1/allowlist`)
- Configuration management endpoints
- Security event query with filtering and pagination
- Live config push from API to WAF core via gRPC
- OpenAPI 3.1 schema at `/api/v1/docs`
- API rate limiting (brute-force protection)

---

## [0.7.0] — TBD

> Web dashboard.

### Added
- React + TypeScript + Vite dashboard (`/dashboard`)
- Login page with JWT authentication and auto-refresh
- RBAC enforcement in UI (Admin, Analyst, Read-only)
- Overview page: live traffic metrics (req/s, block rate, p50/p95/p99 latency)
- Events page: paginated security event log with filter, search, CSV export
- Rules page: list, create, edit, delete with CodeMirror syntax highlighting
- Rule tester: paste HTTP request, see matching rules
- IP Management page: blocklist/allowlist with CIDR range support
- Settings page: WAF mode, paranoia level, rate limit configuration
- User Management page (Admin only)
- WebSocket real-time traffic feed on Overview page
- Dark mode support
- E2E tests with Playwright for critical user paths

---

## [0.8.0] — TBD

> Advanced detection: bots, SSRF, response scanning, virtual patching.

### Added
- Bot detection: user-agent fingerprinting, behavioral analysis
- JS challenge for suspicious clients (configurable)
- TLS fingerprinting (JA3/JA4 hash support)
- Tor exit node and known proxy/VPN blocklist (daily updates)
- Phase 4 (response headers) and Phase 5 (response body) support
- PII leak detection in responses (credit card, SSN, email patterns)
- SSRF protection: block outbound requests to RFC1918 and metadata endpoints
- Enhanced path traversal detection
- File upload scanning: MIME type and magic byte verification
- Virtual patching CLI: `brainless-ctl vpatch apply CVE-XXXX-XXXX`

---

## [0.9.0] — TBD

> Kubernetes-native deployment, Prometheus metrics, SIEM export.

### Added
- Helm chart (`/deploy/helm/brainless-waf`) with HPA, PDB, and ServiceMonitor
- Kubernetes Ingress Controller mode (annotation-driven policies)
- Prometheus metrics endpoint (`/metrics`):
  - `brainless_requests_total`
  - `brainless_blocked_total`
  - `brainless_latency_seconds` (histogram)
  - `brainless_rule_evaluations_total`
- Grafana dashboard JSON (importable, covers all key metrics)
- Elasticsearch, Splunk HEC, and Syslog (RFC 5424) SIEM export
- OpenTelemetry tracing integration (trace ID propagation)

---

## [1.0.0] — TBD

> General Availability — stable, production-ready release.

### Added
- Binary releases for: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
- SHA256 checksums for all release artifacts
- Signed releases with cosign
- `brainless-ctl import-crs` tool for migrating from ModSecurity
- Migration guide from ModSecurity to Brainless WAF
- OWASP CRS 4.x bundled as default ruleset

### Security
- Internal security audit completed
- All known bypass techniques added to regression test suite
- Docker image scan clean (Trivy)
- Supply chain: all dependencies verified, SBOM published

### Performance
- p99 latency <5ms at 10,000 req/s on 8-core reference hardware
- Memory usage ≤380MB at full load
- Rule evaluation: <2ms for 10,000 rules

### Documentation
- Full documentation site live at `https://docs.brainless-security.io`
- All public APIs documented with request/response examples
- Troubleshooting guide with 30+ common issues
- Video quickstart guide

---

## Version History Summary

| Version | Highlights | Date |
|---------|-----------|------|
| 0.1.0 | HTTP proxy skeleton | TBD |
| 0.2.0 | Request parser | TBD |
| 0.3.0 | Rule engine MVP | TBD |
| 0.4.0 | Anomaly scoring | TBD |
| 0.5.0 | TLS + rate limiting | TBD |
| 0.6.0 | Management API | TBD |
| 0.7.0 | Dashboard | TBD |
| 0.8.0 | Advanced detection | TBD |
| 0.9.0 | Kubernetes + observability | TBD |
| 1.0.0 | GA release | TBD |

---

[Unreleased]: https://github.com/brainless-security/brainless-waf/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v1.0.0
[0.9.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.9.0
[0.8.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.8.0
[0.7.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.7.0
[0.6.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.6.0
[0.5.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.5.0
[0.4.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.4.0
[0.3.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.3.0
[0.2.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.2.0
[0.1.0]: https://github.com/brainless-security/brainless-waf/releases/tag/v0.1.0
