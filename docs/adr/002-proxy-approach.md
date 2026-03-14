# ADR-002: Use Go Standard Library for HTTP Proxy (Not NGINX or Envoy)

**Status:** Accepted  
**Date:** 2025  
**Deciders:** Founding team  

---

## Context

The WAF core needs an HTTP proxy layer to receive client requests, inspect them, and forward them to the upstream backend. We need to decide whether to build on top of an existing proxy (NGINX, Envoy, Caddy) or implement our own using Go's standard library.

## Options Considered

### Option A: Build on NGINX (ModSecurity approach)

Pros:
- NGINX is a proven, extremely high-performance proxy used everywhere
- ModSecurity already works this way — familiar approach for the WAF community
- HTTP/2, HTTP/3, TLS are handled by NGINX

Cons:
- NGINX uses C modules — contributors need to know C and the NGINX module API
- Tight coupling to NGINX version — upgrades are risky
- NGINX licensing (NGINX Plus features are commercial) creates confusion for open-source
- Deployment complexity — two components to configure (NGINX + WAF module)
- Debugging is harder — NGINX and WAF module logs are separate

### Option B: Build on Envoy as a filter

Pros:
- Envoy is the standard for cloud-native proxying
- Excellent Kubernetes and service mesh integration
- Filter API is well-defined and stable

Cons:
- Envoy filters must be in C++ (or Wasm, which has performance overhead)
- Very complex to build and test
- Heavy — Envoy binary is ~150MB, requires significant memory
- Overkill for standalone WAF deployments not using a service mesh

### Option C: Build on Caddy

Pros:
- Caddy is written in Go — contributors can use same language as WAF core
- Caddy handles TLS, ACME, HTTP/2, HTTP/3 out of the box
- Caddy plugin API is well-designed

Cons:
- Caddy's extension model couples our release cycle to Caddy's
- Less control over the request lifecycle
- Caddy adds ~30MB binary overhead
- Caddy is opinionated about configuration format (Caddyfile vs YAML)

### Option D: Go standard library (`net/http` + `httputil.ReverseProxy`)

Pros:
- Full control over the entire request lifecycle
- No third-party proxy in the dependency chain
- Single binary with no external process
- Go standard library is extremely well-tested and maintained
- `httputil.ReverseProxy` provides a solid foundation, we extend it
- We can expose exactly the hooks we need for WAF inspection
- HTTP/2 supported natively, HTTP/3 via `golang.org/x/net/http3`

Cons:
- We own TLS termination code (more responsibility)
- ACME integration requires `golang.org/x/crypto/acme` (small external dep)
- More code to write vs. using an existing proxy

## Decision

**Go standard library (`net/http` + custom reverse proxy).**

The control and simplicity advantages outweigh the extra implementation work. We specifically want to avoid a two-component architecture (separate proxy + WAF) that complicates deployment and debugging. With Go's standard library, everything from TLS handshake to response proxying is a single Go binary that contributors can read, understand, and modify.

## Consequences

- We own the HTTP/1.1, HTTP/2 stack via `net/http` (Go stdlib)
- HTTP/3 support via `golang.org/x/net/http3` in v0.5+
- TLS via Go's `crypto/tls` with ACME via `golang.org/x/crypto/acme`
- We build on `httputil.ReverseProxy` and extend it with inspection hooks
- Performance is our responsibility to maintain — we will benchmark every release
