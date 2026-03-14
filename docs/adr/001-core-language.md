# ADR-001: Use Go for the WAF Core Engine

**Status:** Accepted  
**Date:** 2025  
**Deciders:** Founding team  

---

## Context

We need to choose a programming language for the WAF core engine — the component that intercepts, inspects, and proxies every HTTP request. This is the most performance-critical component in the system. The wrong choice here has long-term consequences for throughput, latency, memory usage, and contributor accessibility.

## Decision Drivers

- **Performance:** Needs to handle 10,000+ req/s with <5ms p99 latency overhead
- **Memory safety:** Buffer overflows and use-after-free bugs in a security tool are catastrophic
- **Concurrency:** HTTP proxying is inherently concurrent; the language must handle this well
- **Developer accessibility:** A broader contributor pool means a healthier open-source community
- **Operational simplicity:** Self-contained binary, easy to deploy, no runtime dependencies
- **Ecosystem:** Libraries for HTTP, TLS, regex, gRPC, Prometheus must be available and maintained

## Options Considered

### Option A: Go

Pros:
- Excellent concurrency model (goroutines, channels)
- Fast enough for our throughput targets (10k+ req/s on commodity hardware)
- Strong standard library — `net/http`, `crypto/tls`, `regexp` are battle-tested
- Single static binary — no runtime dependencies, easy to containerize
- Large contributor pool (top 5 most-used languages for backend services)
- Memory-managed (GC) — eliminates most memory safety bugs without manual effort
- Fast compile times (contributor experience)
- GoReleaser for cross-platform binary releases

Cons:
- Garbage collector can cause latency spikes (mitigated by tuning `GOGC` and using sync.Pool)
- Not as fast as C/Rust at the raw level (~2–3x slower for pure compute)
- Less fine-grained memory control than Rust

### Option B: Rust

Pros:
- Maximum performance (comparable to C)
- Zero-cost abstractions
- Memory safety guaranteed at compile time (no GC, no runtime)
- Growing ecosystem for web/network programming

Cons:
- Steep learning curve — significantly smaller contributor pool
- Slower development velocity (borrow checker, longer compile times)
- Standard library is more minimal (more third-party dependencies needed)
- Less mature HTTP server ecosystem vs Go (Hyper is good, but Axum/Actix ecosystem is newer)

### Option C: C with libev/libuv

Pros:
- Maximum performance
- Used by battle-hardened tools (nginx, HAProxy, ModSecurity)

Cons:
- Manual memory management — very high risk of security vulnerabilities in a security tool
- Very small contributor pool for open-source projects
- Complex build system, cross-compilation is hard
- Not acceptable for an open-source tool that relies on community contributions

### Option D: Java / JVM (Go + Kotlin/Scala)

Pros:
- Large developer ecosystem
- Excellent libraries

Cons:
- JVM startup time (~1–5 seconds) unacceptable for our deployment model
- JVM memory overhead (~200MB minimum) vs Go (~45MB idle)
- GC pauses more problematic than Go's GC for latency-sensitive proxying
- Docker image size much larger

## Decision

**Go.**

The primary reasoning:
1. Go's performance is sufficient for our targets (proven by projects like Traefik, Caddy, and Envoy's Go control plane handling much higher throughput)
2. The contributor accessibility advantage over Rust is significant for an open-source project
3. Go's concurrency model is a natural fit for an HTTP proxy
4. The deployment simplicity (single binary) is a meaningful user experience advantage
5. The standard library covers ~70% of what we need without third-party dependencies

The GC latency concern is real but manageable. We will tune `GOGC=200` and use `sync.Pool` for request/response buffer allocation to minimize GC pressure.

## Consequences

- Core engine in `core/` will be Go 1.22+
- Plugins will have a Go API (via `plugin` package or gRPC sidecar for other languages)
- We accept the ~2–3x performance gap vs C/Rust as a trade-off for contributor accessibility
- We will benchmark regularly and optimize hot paths if needed
- Latency spikes from GC must be monitored and kept below 1ms p99.9

## Revisit Condition

If we cannot meet latency targets (p99 >5ms) under target load after optimization, we will evaluate rewriting the hottest paths in C via CGo, or evaluate Rust for v2.0.
