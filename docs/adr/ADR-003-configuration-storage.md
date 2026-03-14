# ADR-003: Choice of PostgreSQL for Configuration Storage

## Status
Accepted

## Context
The Management API needs a persistent store for user accounts, rules, configuration, and security event logs.

## Decision
We will use **PostgreSQL** as the primary relational database.

## Rationale
- **Reliability:** PostgreSQL is a mature, ACID-compliant database suitable for critical security configuration.
- **JSON Support:** Excellent support for JSONB allows us to store flexible rule data and event details alongside structured data.
- **Scalability:** Handles large volumes of event logs efficiently with partitioning and indexing.
- **Ecosystem:** Strong integration with SQLAlchemy (Python) and Go-based migration tools.

## Consequences
- Adds PostgreSQL as a required infrastructure dependency.
- Requires careful schema management for telemetry data growth.
