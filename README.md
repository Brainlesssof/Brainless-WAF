<div align="center">

# 🛡️ Brainless WAF

**Enterprise-grade Web Application Firewall. Open-Source. Zero Cost.**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8.svg)](https://golang.org)
[![Build Status](https://img.shields.io/github/actions/workflow/status/brainless-security/brainless-waf/ci.yml?branch=main)](https://github.com/brainless-security/brainless-waf/actions)
[![Coverage](https://img.shields.io/codecov/c/github/brainless-security/brainless-waf)](https://codecov.io/gh/brainless-security/brainless-waf)
[![Discord](https://img.shields.io/discord/placeholder?label=Discord&logo=discord)](https://discord.gg/brainless-waf)
[![GitHub Stars](https://img.shields.io/github/stars/brainless-security/brainless-waf?style=social)](https://github.com/brainless-security/brainless-waf/stargazers)

[Documentation](https://docs.brainless-security.io) · [Quick Start](#quick-start) · [Report Bug](https://github.com/brainless-security/brainless-waf/issues/new?template=bug_report.md) · [Request Feature](https://github.com/brainless-security/brainless-waf/issues/new?template=feature_request.md) · [Discord](https://discord.gg/brainless-waf)

</div>

---

## What is Brainless WAF?

Brainless WAF is a high-performance, open-source Web Application Firewall built to deliver enterprise-grade protection without enterprise-grade licensing costs. Inspired by F5 Advanced WAF and Cloudflare WAF, Brainless WAF gives you a full-featured, extensible security platform that any organization can deploy, customize, and contribute to.

```
Client → [Brainless WAF] → Your Backend
           │
           ├── OWASP Top 10 Protection
           ├── Custom Rule Engine (BRF + ModSecurity Compatible)
           ├── Bot Detection & DDoS Mitigation
           ├── Rate Limiting per IP / Endpoint
           ├── TLS Termination (ACME / Let's Encrypt)
           └── Real-time Dashboard + REST API
```

---

## Features

| Category | Capabilities |
|----------|-------------|
| **Attack Detection** | OWASP Top 10, SQLi, XSS, CSRF, RFI, LFI, SSRF, Path Traversal |
| **Rule Engine** | BRF (native), ModSecurity CRS 4.x compatible, YARA, Lua scripting |
| **Bot Protection** | User-agent fingerprinting, JS challenge, TLS fingerprinting (JA3/JA4) |
| **Rate Limiting** | Per-IP, per-endpoint, token bucket with burst support |
| **TLS** | TLS 1.2/1.3, ACME/Let's Encrypt auto-renewal, HTTP/2 + HTTP/3 |
| **Deployment** | Docker, Kubernetes (Helm), bare metal, API Gateway mode |
| **Observability** | Prometheus metrics, Grafana dashboard, OpenTelemetry tracing, SIEM export |
| **Management** | REST API, React dashboard, RBAC (Admin / Analyst / Read-only) |

---

## Quick Start

### Docker (Recommended)

```bash
# Clone and start
git clone https://github.com/brainless-security/brainless-waf.git
cd brainless-waf
cp config/default.yaml config/local.yaml

# Edit config/local.yaml and set your upstream server
docker compose up -d

# Dashboard → http://localhost:8080
# Default: admin / changeme  ← change this immediately!
```

### Kubernetes (Helm)

```bash
helm repo add brainless https://charts.brainless-security.io
helm repo update
helm install brainless-waf brainless/brainless-waf \
  --namespace brainless-waf \
  --create-namespace
```

### Binary

```bash
# Linux / macOS
curl -sSL https://get.brainless-security.io | sh

# Configure and run
brainless-waf --config config/local.yaml
```

> **First deployment?** Read the [Installation Guide](docs/installation.md) and start in `learning` mode for 72 hours before switching to `block`.

---

## Minimum Configuration

```yaml
# config/local.yaml
server:
  listen: 0.0.0.0:80
  listen_tls: 0.0.0.0:443
  upstream: http://your-backend:8080

detection:
  mode: learning        # learning → detect → block
  paranoia_level: 2

rules:
  crs_enabled: true
  auto_update: true
```

---

## Documentation

| Document | Description |
|----------|-------------|
| [Installation Guide](docs/installation.md) | All deployment methods (Docker, K8s, bare metal) |
| [Configuration Reference](docs/configuration.md) | Every config option explained |
| [Rule Engine Guide](docs/rules.md) | BRF syntax, operators, actions, examples |
| [API Reference](docs/api.md) | REST API endpoints with examples |
| [Dashboard Guide](docs/dashboard.md) | Using the web UI |
| [Troubleshooting](docs/troubleshooting.md) | Common issues and fixes |
| [Development Guide](DEVELOPMENT.md) | Set up a local dev environment |
| [Roadmap](ROADMAP.md) | What's planned and when |
| [Changelog](CHANGELOG.md) | Full version history |

---

## Performance

Benchmarked on a single `c5.2xlarge` instance (8 vCPU, 16GB RAM):

| Metric | Value |
|--------|-------|
| Throughput | 12,000 req/s (sustained) |
| Latency added (p50) | 0.8ms |
| Latency added (p99) | 3.2ms |
| Memory (idle) | ~45MB |
| Memory (12k req/s) | ~380MB |
| Rule evaluation (10k rules) | <2ms |

---

## Comparison

| Feature | Brainless WAF | F5 Advanced WAF | Cloudflare WAF |
|---------|--------------|----------------|----------------|
| License | Apache 2.0 (Free) | Commercial | Commercial SaaS |
| Self-hosted | ✅ | ✅ | ❌ |
| Source code | ✅ Full access | ❌ | ❌ |
| OWASP Top 10 | ✅ | ✅ | ✅ |
| Custom rules | ✅ Unlimited | ✅ | Limited by plan |
| Lua scripting | ✅ | ✅ | ❌ |
| Kubernetes native | ✅ | Partial | ❌ |
| Annual cost | $0 | $15k–$80k+ | $200–$5k+/mo |

---

## Contributing

We welcome all contributions — code, docs, rules, bug reports, ideas.

```bash
# Fork → branch → change → test → PR
git checkout -b feature/my-improvement
make test && make lint
git commit -m "feat: my improvement"
# Open PR against develop
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for the full guide.  
New to the project? Look for [`good-first-issue`](https://github.com/brainless-security/brainless-waf/labels/good-first-issue) labels.

---

## Community

- **Discord:** [discord.gg/brainless-waf](https://discord.gg/brainless-waf) — real-time chat, help, and announcements
- **GitHub Discussions:** questions, ideas, RFCs
- **Security issues:** `security@brainless-security.io` — please do **not** use public issues for vulnerabilities

---

## License

Brainless WAF is licensed under the [Apache License 2.0](LICENSE).

Third-party components retain their own licenses. See [THIRD-PARTY-LICENSES.md](THIRD-PARTY-LICENSES.md).

---

<div align="center">
Made with ❤️ by the Brainless Security community
</div>
