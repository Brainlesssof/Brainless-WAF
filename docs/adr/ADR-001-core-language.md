# ADR-001: Choice of Go for the Core Engine

## Status
Accepted

## Context
The Core Engine of Brainless WAF must handle high volumes of concurrent HTTP traffic with minimal latency overhead. It needs to perform complex request parsing, normalization, and rule evaluation in real-time.

## Decision
We will use **Go** (Golang) as the primary language for the Core Engine.

## Rationale
- **Concurrency:** Go's Goroutines and Channels provide a robust and lightweight model for handling thousands of simultaneous connections.
- **Performance:** As a compiled language, Go offers near-C performance while maintaining high developer productivity.
- **Networking Stack:** Go's standard library (`net/http`) and high-performance alternatives (like `fasthttp`) are industry-leading for building proxies.
- **Security:** Go's memory safety prevents many common vulnerabilities like buffer overflows.
- **Static Binaries:** Go produces single static binaries, simplifying deployment in Docker and bare-metal environments.

## Consequences
- Developers must be proficient in Go.
- Integration with Python-based Management API will require a cross-language protocol (gRPC).
