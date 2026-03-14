# Architecture

This document describes the internal architecture of Brainless WAF — how components are organized, how they communicate, and the design decisions behind each one.

---

## Overview

Brainless WAF is a **multi-process, pipeline-based** system. Three independent processes handle traffic inspection, configuration management, and the user interface. They communicate over well-defined interfaces, which means each component can be scaled, replaced, or disabled independently.

```
                    ┌─────────────────────────────────────────┐
                    │           Brainless WAF System           │
                    │                                          │
  Internet          │  ┌──────────────┐    ┌────────────────┐ │
  ─────────────────►│  │  WAF Core    │───►│ Your Backend   │ │
  HTTP/HTTPS        │  │  (Go)        │    │ (any server)   │ │
                    │  └──────┬───────┘    └────────────────┘ │
                    │         │ gRPC                           │
                    │  ┌──────▼───────┐                        │
                    │  │  Mgmt API    │◄── REST API ─── Admin  │
                    │  │  (Python)    │                         │
                    │  └──────┬───────┘                        │
                    │         │ SQL                             │
                    │  ┌──────▼───────┐    ┌────────────────┐ │
                    │  │  PostgreSQL  │    │   Dashboard    │ │
                    │  │              │    │   (React)      │ │
                    │  └──────────────┘    └────────────────┘ │
                    └─────────────────────────────────────────┘
```

---

## Component: WAF Core (Go)

The WAF core is the only component in the data path. Every HTTP request passes through it. It is written in Go for performance and safety.

### Request Pipeline

Every request flows through stages in order. A stage can pass the request forward, modify it, or terminate it (block, redirect):

```
Incoming Request
      │
      ▼
 ┌─────────────┐
 │ 1. TLS      │  Terminate HTTPS, extract client certificate if present
 └──────┬──────┘
        │
        ▼
 ┌─────────────┐
 │ 2. Parser   │  Normalize URI, headers, body (URL decode, Unicode normalize)
 └──────┬──────┘
        │
        ▼
 ┌─────────────┐
 │ 3. Phase 1  │  Evaluate rules on request headers and URI (before body read)
 │   Rules     │
 └──────┬──────┘
        │
        ▼
 ┌─────────────┐
 │ 4. Phase 2  │  Evaluate rules on request body (POST/PUT data)
 │   Rules     │
 └──────┬──────┘
        │
        ▼
 ┌─────────────┐
 │ 5. Score    │  Sum anomaly scores; block if threshold exceeded
 │   Check     │
 └──────┬──────┘
        │
        ▼
 ┌─────────────┐
 │ 6. Proxy    │  Forward to upstream backend
 └──────┬──────┘
        │
        ▼
 ┌─────────────┐
 │ 7. Phase 4  │  Evaluate rules on response headers
 │   Rules     │
 └──────┬──────┘
        │
        ▼
 ┌─────────────┐
 │ 8. Phase 5  │  Evaluate rules on response body (optional, adds latency)
 │   Rules     │
 └──────┬──────┘
        │
        ▼
 Response to Client
```

### Rule Engine

The rule engine is the most performance-critical part of the system. Key design decisions:

- All regular expressions are **pre-compiled** at startup — zero compilation cost per request
- Rules are stored in an ordered list per phase, evaluated sequentially
- Short-circuit evaluation: a `deny` action immediately terminates the pipeline
- Anomaly scoring accumulates across all phases; the final block decision is at stage 5

### Rate Limiter

The rate limiter uses a **token bucket** algorithm per IP address:

- Each IP gets a bucket with capacity `burst`
- Tokens refill at `requests_per_second` rate
- A request costs 1 token; if the bucket is empty, return 429
- Buckets are stored in memory (single instance) or Redis (clustered)
- Stale buckets (no activity for 5 minutes) are evicted from memory

### Configuration Hot-Reload

The WAF core exposes a **gRPC server** on an internal port (not publicly accessible). The Management API sends config updates via gRPC. The core applies changes atomically:

1. Receive new config/rules via gRPC
2. Validate and compile (regex pre-compilation)
3. Acquire write lock
4. Swap the active config pointer
5. Release lock

In-flight requests complete with the old config. New requests use the new config. Zero downtime, no restart needed.

---

## Component: Management API (Python/FastAPI)

The Management API is the **control plane**. It handles:

- Authentication and authorization (JWT + API keys)
- CRUD operations on rules, IP lists, site configuration
- Persisting everything to PostgreSQL
- Pushing config changes to WAF core via gRPC
- Serving logs and event data to the dashboard

It is intentionally **not** in the request data path. Even if the Management API is down, the WAF core continues protecting traffic using its last-known configuration.

### Why Python?

The Management API is I/O-bound (database queries, gRPC calls). Python with async/await (FastAPI + asyncpg) provides excellent performance for this workload with a much faster development cycle than Go for CRUD-heavy services.

### Database Schema (simplified)

```
users          → id, username, hashed_password, role, created_at
api_keys       → id, user_id, key_hash, name, scopes, expires_at
rules          → id, rule_text, enabled, created_by, created_at, updated_at
ip_blocklist   → id, cidr, reason, expires_at, created_by
ip_allowlist   → id, cidr, comment, created_by
events         → id, timestamp, severity, rule_id, src_ip, uri, action, details
config         → id, key, value, updated_at, updated_by
```

---

## Component: Dashboard (React/TypeScript)

The dashboard is a **single-page application** served as static files. It communicates exclusively with the Management API over HTTP. It has no direct connection to the WAF core.

The dashboard is optional — all operations available in the UI are also available via the REST API. Organizations that want to integrate Brainless WAF into their own tooling can disable the dashboard entirely.

### Real-time Data

The overview page uses a **WebSocket connection** to the Management API for live traffic metrics. The Management API subscribes to a Redis pub/sub channel that the WAF core publishes metrics to every second.

```
WAF Core → Redis pub/sub → Management API → WebSocket → Dashboard
```

---

## Data Flow: Blocked Request Example

```
1. Attacker sends: GET /search?q=1' OR '1'='1 HTTP/1.1

2. TLS terminates the connection (if HTTPS)

3. Parser normalizes:
   - URL-decodes: /search?q=1' OR '1'='1
   - Extracts ARGS: { "q": "1' OR '1'='1" }

4. Phase 1 rules evaluate REQUEST_URI:
   - Rule 942100 (@rx for SQL tautology patterns) matches
   - Adds anomaly score: +5 (ERROR severity)

5. Phase 2 rules evaluate ARGS:
   - Rule 942200 (@rx for SQL injection in parameters) matches
   - Adds anomaly score: +10 (CRITICAL severity)

6. Score check: 5 + 10 = 15 >= threshold (10)
   → BLOCK

7. WAF returns 403 Forbidden to attacker
   Body: {"error": "Request blocked by security policy", "id": "evt_abc123"}

8. Event written to PostgreSQL:
   { severity: CRITICAL, rule: 942200, src_ip: 1.2.3.4, uri: "/search", action: "block" }

9. Management API broadcasts event to dashboard via WebSocket

10. Dashboard updates block counter in real-time
```

---

## Security Architecture

### No Trust Between Components

Even though the WAF core and Management API run on the same host (or cluster), they authenticate each other:

- Management API → WAF Core: mTLS on the gRPC channel
- Dashboard → Management API: JWT Bearer token

### WAF Core Isolation

The WAF core process:
- Runs as non-root (`uid: 1000` in Docker)
- Has read-only access to the filesystem (except `/tmp` and log dir)
- Cannot make outbound network connections (except to the upstream backend and the gRPC server)
- Has no access to the PostgreSQL database — config is pushed to it, not pulled

This means a compromise of the WAF core cannot directly access the credential database.

### Fail Safe

If any internal error occurs in the detection pipeline (panic recovery, timeout, resource exhaustion):
- The WAF returns **502 Bad Gateway** to the client
- The request is **not passed through** to the backend
- The error is logged with full context for debugging

We prefer a brief outage over accidentally allowing malicious traffic.

---

## Scalability

### Single Instance

A single WAF core instance handles 12,000+ req/s on 8 CPU cores. For most deployments, one instance is sufficient.

### Horizontal Scaling

For high-traffic deployments:

```
                    ┌──────────────────────────────────┐
     Internet       │     Load Balancer (L4/L7)         │
    ──────────────► │   (nginx, HAProxy, cloud LB)      │
                    └────────────┬────────────┬─────────┘
                                 │            │
                          ┌──────▼───┐  ┌─────▼────┐
                          │  WAF     │  │  WAF     │  ← Stateless
                          │  Core 1  │  │  Core 2  │
                          └──────┬───┘  └─────┬────┘
                                 │            │
                                 └─────┬──────┘
                                       │ gRPC
                                ┌──────▼───────┐
                                │  Mgmt API    │  ← Single instance (or HA)
                                └──────┬───────┘
                                       │
                                ┌──────▼───────┐
                                │  PostgreSQL  │
                                │  + Redis     │
                                └──────────────┘
```

WAF core instances are **stateless** with respect to request processing — all persistent state (rate limit counters, block/allowlists) is in Redis. Multiple instances share the same Redis and receive the same configuration from the Management API.

---

## Architecture Decision Records

Major architectural choices are documented as ADRs in `docs/adr/`:

| ADR | Decision |
|-----|---------|
| [ADR-001](adr/001-core-language.md) | Use Go for the WAF core |
| [ADR-002](adr/002-proxy-approach.md) | Use Go standard library `net/http` (not NGINX or Envoy) |
| [ADR-003](adr/003-database.md) | Use PostgreSQL for management plane |
| [ADR-004](adr/004-dashboard-framework.md) | Use React + TypeScript for dashboard |
| [ADR-005](adr/005-rule-compatibility.md) | Target ModSecurity CRS 4.x compatibility |
| [ADR-006](adr/006-internal-comms.md) | Use gRPC for core↔management communication |
