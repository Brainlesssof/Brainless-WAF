# ADR-003: Use PostgreSQL for Management Plane Storage

**Status:** Accepted  
**Date:** 2025  

---

## Context

The management plane needs persistent storage for rules, configuration, users, events, and IP lists.

## Options Considered

### SQLite
- Pros: No separate process, simple deployment, file-based backup
- Cons: Poor write concurrency (WAL mode helps but has limits), not suitable for multi-instance management API, limited data types

### PostgreSQL
- Pros: Excellent concurrency, JSONB for flexible event data, full-text search on logs, mature replication, widely understood by ops teams, managed services available on every cloud
- Cons: Separate process to manage, more complex local dev setup (mitigated by Docker)

### MySQL / MariaDB
- Pros: Widely used, good performance
- Cons: Less feature-rich than PostgreSQL (JSONB, full-text search, array types), PostgreSQL has better community reputation in the open-source/DevOps world

## Decision

**PostgreSQL 15+.**

PostgreSQL's JSONB type is particularly valuable for event storage — security events have variable structure and JSONB lets us query event details without a rigid schema. The managed service availability (RDS, Cloud SQL, Supabase) means operators don't have to self-manage PostgreSQL if they don't want to.

## Consequences

- Management API requires PostgreSQL 15+
- Docker Compose bundles PostgreSQL for convenience
- Alembic for schema migrations
- Events table uses JSONB for the `details` field
- We provide a Docker Compose with PostgreSQL for local dev — no need to install it locally

---

# ADR-004: Use React + TypeScript for the Dashboard

**Status:** Accepted  
**Date:** 2025  

---

## Context

We need a web UI for the WAF dashboard. We need to choose a frontend framework.

## Options Considered

### React + TypeScript
- Pros: Largest ecosystem, most contributors familiar with it, excellent TypeScript support, Vite for fast dev experience, component libraries available
- Cons: Not opinionated (more boilerplate than Vue/Svelte for some patterns)

### Vue 3 + TypeScript
- Pros: More opinionated (less boilerplate), excellent TypeScript support, growing ecosystem
- Cons: Smaller community than React, fewer contributors likely to be familiar

### Svelte / SvelteKit
- Pros: Excellent performance, very little boilerplate, small bundle sizes
- Cons: Smallest contributor pool of the three, less mature ecosystem, fewer component libraries

### HTMX + server-rendered HTML
- Pros: Minimal JavaScript, simpler mental model
- Cons: Poor fit for real-time dashboards with WebSocket updates and complex interactive charts, cannot build a rich rule editor with this approach

## Decision

**React 18 + TypeScript (strict mode) + Vite.**

The contributor pool argument strongly favors React. The dashboard is the component most likely to receive contributions from developers who are not security experts — frontend developers who want to improve the UI. React gives us the highest chance those contributors already know the framework.

## Consequences

- Dashboard in `dashboard/` is React 18 + TypeScript
- State management: React Context for global state (no Redux — too much boilerplate for our needs)
- Charts: Recharts (React-native, good TypeScript types)
- Code editor: CodeMirror 6 (for rule editor with syntax highlighting)
- API client: Auto-generated from OpenAPI schema using openapi-typescript-codegen

---

# ADR-005: Target ModSecurity CRS 4.x Compatibility

**Status:** Accepted  
**Date:** 2025  

---

## Context

We need to decide on our rule format. We could create a completely new format or be compatible with existing formats (primarily ModSecurity's SecLang / CRS).

## Options Considered

### Completely new rule format
- Pros: Can be designed optimally for our architecture, no legacy baggage
- Cons: Zero ecosystem — no existing rules, requires everyone to learn a new format, no migration path from ModSecurity

### ModSecurity 2.x compatible
- Pros: Massive existing ecosystem (CRS, OWASP rules, commercial rulesets), low migration friction
- Cons: ModSecurity 2.x has known design flaws and deprecated features, compatibility is complex

### ModSecurity 3.x / CRS 4.x compatible (our choice)
- Pros: Modern CRS 4.x is the current standard, large ecosystem, migration path for existing ModSecurity users, OWASP backing ensures ongoing development
- Cons: CRS 4.x is still more complex than a clean-room design, some ModSecurity directives are hard to implement

## Decision

**Target CRS 4.x compatibility with BRF (Brainless Rule Format) as the canonical format.**

BRF is a clean superset of the SecLang syntax used by CRS 4.x, with extensions for features not possible in ModSecurity (Lua scripting, ML-based operators, per-path config). This gives us:
1. Import existing CRS rulesets without changes
2. Migration path from ModSecurity deployments
3. A clean foundation for future innovation

## Consequences

- BRF spec documented in `docs/rules.md`
- `brainless-ctl import-crs` tool for importing CRS rules
- Compatibility matrix maintained in `docs/crs-compatibility.md`
- ModSecurity directives that cannot be supported are documented and a BRF equivalent provided

---

# ADR-006: Use gRPC for Internal Core↔Management Communication

**Status:** Accepted  
**Date:** 2025  

---

## Context

The Management API needs to push configuration changes (rule updates, IP list changes, mode changes) to the WAF core in real time. We need a communication protocol.

## Options Considered

### REST / HTTP API on the WAF core
- Pros: Simple, well-understood
- Cons: Polling-based (management API must poll for changes) or we need webhooks, less efficient for streaming config updates

### Message queue (Redis pub/sub, NATS, RabbitMQ)
- Pros: Decoupled, supports fan-out to multiple WAF instances
- Cons: Additional dependency, more complex failure modes, eventual consistency

### gRPC
- Pros: Bi-directional streaming (WAF core can push metrics to management), strongly typed with Protobuf, efficient binary protocol, generated client/server code, built-in TLS (mTLS)
- Cons: More complex than REST, requires Protobuf toolchain, binary protocol is harder to debug

### Unix domain socket with custom binary protocol
- Pros: Very fast for single-host deployments
- Cons: Doesn't work for multi-instance (network required), custom protocol is maintenance burden

## Decision

**gRPC with Protobuf.**

The bi-directional streaming capability is the deciding factor. We want the WAF core to push real-time traffic metrics to the Management API (which then pushes to the dashboard via WebSocket). gRPC's server-streaming RPC handles this cleanly. The mTLS authentication built into gRPC also removes a class of security concerns about internal API access.

## Consequences

- Protobuf schemas in `core/proto/`
- Generated Go code in `core/internal/grpc/`
- Generated Python code in `management/app/grpc/`
- `make gen-grpc` regenerates both from `.proto` files
- gRPC server binds to `127.0.0.1:9091` (not externally accessible)
- mTLS certificates generated at startup for internal use (self-signed, rotated on restart)
