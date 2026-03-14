# ADR-002: Custom TCP/HTTP Stack vs. NGINX

## Status
Accepted

## Context
We need to decide whether to build a custom traffic interceptor or use an existing solution like NGINX/OpenResty for the WAF Core.

## Decision
We will build a **custom HTTP reverse proxy** in Go, rather than relying on NGINX or OpenResty.

## Rationale
- **Control:** A custom stack allows for deeper integration of the Brainless Rule Format (BRF) without the limitations of NGINX module development.
- **Performance:** We can optimize the request parsing path specifically for security scanning, avoiding the overhead of general-purpose proxy features we don't need.
- **Portability:** A Go-based proxy is easier to distribute and run in diverse environments (K8s, Edge, Bare-metal) without managing NGINX dependencies.
- **Observability:** Custom implementation allows for native integration of Prometheus metrics and OpenTelemetry tracing directly into the request lifecycle.

## Consequences
- Increased development effort for standard proxy features (retries, timeouts, connection pooling).
- Need to implement and maintain a robust TLS stack.
