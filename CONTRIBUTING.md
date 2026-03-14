# Contributing to Brainless WAF

Thank you for your interest in contributing! Brainless WAF is a community-driven project and every contribution — whether it's a bug report, a new rule, a documentation fix, or a major feature — makes the project better for everyone.

---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Ways to Contribute](#ways-to-contribute)
- [Before You Start](#before-you-start)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Commit Messages](#commit-messages)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Writing Tests](#writing-tests)
- [Contributing Rules](#contributing-rules)
- [Contributing Documentation](#contributing-documentation)
- [Reporting Bugs](#reporting-bugs)
- [Requesting Features](#requesting-features)
- [Security Vulnerabilities](#security-vulnerabilities)
- [Getting Help](#getting-help)

---

## Code of Conduct

This project follows our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold it. Please report unacceptable behavior to `conduct@brainless-security.io`.

---

## Ways to Contribute

You don't need to write code to contribute. Here are all the ways you can help:

- **Report bugs** — file an issue using the bug report template
- **Write rules** — contribute detection rules for new attack patterns
- **Improve docs** — fix typos, add examples, clarify confusing sections
- **Answer questions** — help others in GitHub Discussions or Discord
- **Write tests** — improve test coverage for existing features
- **Build features** — pick up an issue labeled `help-wanted`
- **Review PRs** — provide thoughtful feedback on open pull requests
- **Share** — write a blog post, give a talk, star the repo

---

## Before You Start

For anything beyond a small fix (typo, broken link), **please open or comment on an issue first**. This avoids wasted effort if your change doesn't align with the project's direction.

- Check [existing issues](https://github.com/brainless-security/brainless-waf/issues) before opening a new one
- For major changes, open a Discussion to get early feedback
- For security vulnerabilities, see [Security Vulnerabilities](#security-vulnerabilities) — do NOT open a public issue

---

## Development Setup

### Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.22+ | Core engine |
| Node.js | 20+ | Dashboard |
| Python | 3.11+ | Management API |
| Docker | 24.0+ | Local stack |
| Docker Compose | v2+ | Local stack |
| Make | any | Build automation |

### Clone and Start

```bash
# 1. Fork the repo on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/brainless-waf.git
cd brainless-waf

# 2. Add upstream remote
git remote add upstream https://github.com/brainless-security/brainless-waf.git

# 3. Install all dependencies
make deps

# 4. Start the full dev stack (WAF + mock backend + PostgreSQL + Redis)
make dev

# 5. Verify everything is running
make health-check
```

The dev stack starts:
- WAF core on `http://localhost:9090`
- Management API on `http://localhost:8000`
- Dashboard on `http://localhost:5173` (Vite dev server with HMR)
- Mock backend on `http://localhost:3000`
- PostgreSQL on `localhost:5432`

### Pre-commit Hooks

Install the pre-commit hooks to catch issues before they reach CI:

```bash
pip install pre-commit
pre-commit install
```

Hooks run: `golangci-lint`, `ruff`, `mypy`, `prettier`, `trufflehog` (secret scanning).

---

## Making Changes

### 1. Sync with upstream

```bash
git checkout develop
git fetch upstream
git merge upstream/develop
```

### 2. Create a branch

```bash
# Branch naming: type/short-description
git checkout -b feature/add-graphql-detection
git checkout -b fix/rate-limit-header-overflow
git checkout -b docs/improve-rule-examples
```

Branch types: `feature/`, `fix/`, `security/`, `docs/`, `refactor/`, `test/`, `chore/`

### 3. Make your changes

Write your code, add tests, update docs. See [Coding Standards](#coding-standards).

### 4. Run checks locally

```bash
make test          # Run all unit tests
make lint          # Run all linters
make fmt           # Auto-format code
make build         # Ensure everything compiles
```

All checks must pass before opening a PR.

---

## Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/). This enables automatic changelog generation and clear history.

**Format:**
```
<type>(<scope>): <short description>

[optional body]

[optional footer]
```

**Types:**

| Type | When to use |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `test` | Adding or updating tests |
| `refactor` | Code change without feature or fix |
| `perf` | Performance improvement |
| `chore` | Tooling, dependencies, CI |
| `security` | Security patch or hardening |
| `ci` | CI/CD configuration changes |

**Scopes:** `core`, `api`, `dashboard`, `rules`, `deploy`, `docs`, `deps`

**Examples:**

```
feat(rules): add YARA rule parser support

fix(core): prevent null pointer dereference in URI normalizer

security(core): patch path traversal bypass via double URL encoding
Closes #142

docs(api): add authentication examples to rate-limit endpoints

test(rules): add SQLi bypass regression tests for CVE-2024-1234
```

**Rules:**
- Use the imperative mood: "add feature" not "added feature"
- Keep the subject line under 72 characters
- Reference issues in the footer: `Closes #123`, `Fixes #456`, `Relates to #789`
- Breaking changes: add `BREAKING CHANGE:` in the footer

---

## Pull Request Process

### Before Opening

- [ ] All tests pass: `make test`
- [ ] Linter is clean: `make lint`
- [ ] New code has tests (≥80% coverage on changed files)
- [ ] Documentation is updated if needed
- [ ] `CHANGELOG.md` has an entry under `## [Unreleased]` for user-visible changes
- [ ] Branch is up to date with `develop`

### Opening the PR

1. Open the PR against the `develop` branch (never `main` directly)
2. Fill in the PR template completely — incomplete PRs will be asked to update
3. Link the related issue: `Closes #123`
4. Add appropriate labels
5. Request a reviewer if you know who owns that area (see [CODEOWNERS](.github/CODEOWNERS))

### Review Process

- A maintainer will review within **5 business days** (usually faster)
- Address all review comments — use "Resolve conversation" only after the reviewer is satisfied
- Keep the branch up to date if `develop` moves ahead
- Once approved + CI green → a maintainer will merge (squash merge preferred)

### After Merge

Your contribution will appear in the next release. The maintainer will handle the CHANGELOG and release notes.

---

## Coding Standards

### Go (Core Engine)

- Follow [Effective Go](https://go.dev/doc/effective_go) and the [Google Go Style Guide](https://google.github.io/styleguide/go/)
- All exported functions and types **must** have godoc comments
- Error handling: never ignore errors with `_`; wrap with `fmt.Errorf("context: %w", err)`
- No global mutable state outside of explicitly documented singletons
- Benchmarks required for any change to hot-path code (request parsing, rule matching)

```bash
# Linting
golangci-lint run ./...

# Formatting (enforced by CI)
gofmt -w .
goimports -w .

# Vulnerability scan
govulncheck ./...
```

### Python (Management API)

- Follow [PEP 8](https://peps.python.org/pep-0008/) with `ruff` for enforcement
- Type hints on all function signatures (enforced by `mypy --strict`)
- Use `async/await` throughout — no synchronous blocking calls in async context
- Pydantic models for all request/response schemas

```bash
ruff check .
ruff format .
mypy management/ --strict
```

### TypeScript (Dashboard)

- Strict TypeScript (`"strict": true` in `tsconfig.json`)
- No `any` types — use `unknown` and narrow properly
- Components: functional with hooks, no class components
- State management: React context for global state, local `useState` for UI state

```bash
cd dashboard
npm run lint
npm run type-check
```

### General Rules

- No secrets or credentials in code — use environment variables
- No `TODO` comments in merged code — convert to issues before merging
- Keep functions small and focused — if a function is >50 lines, consider splitting
- Prefer explicit over implicit

---

## Writing Tests

### Go Tests

```bash
# Run unit tests
go test ./...

# Run with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run a specific test
go test -run TestRuleEngine_SQLiDetection ./core/rules/...

# Benchmarks
go test -bench=BenchmarkRequestParser -benchmem ./core/...

# Fuzz testing (for parsers)
go test -fuzz=FuzzURIParser -fuzztime=60s ./core/...
```

### Test File Naming

```
core/
  rules/
    engine.go
    engine_test.go       ← unit tests
    engine_bench_test.go ← benchmarks
    engine_fuzz_test.go  ← fuzz tests
```

### Test Data

- Place test payloads in `tests/fixtures/`
- Attack payloads go in `tests/fixtures/attacks/` organized by attack type
- Safe payloads (should NOT trigger rules) go in `tests/fixtures/safe/`

### Integration Tests

```bash
make test-integration    # Requires Docker (starts test stack automatically)
```

Integration tests live in `tests/integration/` and test the full request pipeline.

---

## Contributing Rules

Detection rules are one of the most impactful contributions. A good rule protects every deployment.

### Rule File Location

```
rules/
  custom/
    YOUR_RULE_FILE.rules   ← your new rules go here
  tests/
    YOUR_RULE_FILE_test.yaml  ← required test file
```

### Rule Requirements

Every new rule must include:

1. **A unique ID** in the range `50000–59999` (community rules range)
2. **A test file** with at least one positive (should-block) and one negative (should-pass) payload
3. **A comment** explaining what the rule detects and why
4. **Tags** for OWASP category, attack type, and CVE if applicable

### Rule Test Format

```yaml
# rules/tests/my_rule_test.yaml
rule_id: 50001
description: "Test for My Attack Pattern"

should_block:
  - description: "Basic attack payload"
    request:
      method: GET
      uri: "/search?q=<script>alert(1)</script>"
      headers:
        Host: example.com

should_pass:
  - description: "Normal search query"
    request:
      method: GET
      uri: "/search?q=hello+world"
      headers:
        Host: example.com
```

Run rule tests:
```bash
make test-rules RULE=rules/custom/my_rule.rules
```

---

## Contributing Documentation

Documentation lives in the `docs/` directory and is written in Markdown.

- Keep language clear and direct — assume technical readers, not security experts
- Include working code examples for every feature you describe
- Update the table of contents when adding new sections
- Screenshots go in `docs/assets/images/` — use PNG, max 1200px wide

For major documentation additions, open an issue first so we can agree on structure.

---

## Reporting Bugs

Use the [bug report template](https://github.com/brainless-security/brainless-waf/issues/new?template=bug_report.md).

A good bug report includes:
- **Version:** output of `brainless-waf --version`
- **Environment:** OS, Docker version, Kubernetes version if applicable
- **Steps to reproduce:** minimal, numbered steps
- **Expected behavior:** what should happen
- **Actual behavior:** what actually happened
- **Logs:** relevant log output (sanitize any sensitive data)

---

## Requesting Features

Use the [feature request template](https://github.com/brainless-security/brainless-waf/issues/new?template=feature_request.md).

Before requesting, check:
- [Existing issues](https://github.com/brainless-security/brainless-waf/issues) (including closed ones)
- [The roadmap](ROADMAP.md) — it may already be planned
- [GitHub Discussions](https://github.com/brainless-security/brainless-waf/discussions) — it may already be discussed

---

## Security Vulnerabilities

**Do not open a public issue for security vulnerabilities.**

Please report security issues privately to `security@brainless-security.io`. We follow responsible disclosure:

1. You report privately
2. We acknowledge within 48 hours
3. We investigate and develop a fix (typically 7–14 days)
4. We release a patch and notify you
5. You may publish details 30 days after the patch is released

See [SECURITY.md](SECURITY.md) for full details.

---

## Getting Help

- **Discord `#dev` channel** — fastest way to get help from maintainers and contributors
- **GitHub Discussions** — for longer questions and design discussions
- **Code review** — leave questions directly on PR lines

We're a friendly community. No question is too basic — ask away.

---

*Thank you for making Brainless WAF better for everyone.* 🛡️
