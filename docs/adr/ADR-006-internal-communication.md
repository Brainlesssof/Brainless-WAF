# ADR-006: Internal Communication via gRPC

## Status
Accepted

## Context
We need a fast and reliable way for the Management API (Python) to communicate configuration and rule updates to the Core Engine (Go).

## Decision
We will use **gRPC** (over Protocol Buffers) for internal component communication.

## Rationale
- **Performance:** Binary serialization is more efficient than JSON over REST, reducing latency for rule hot-reloads.
- **Type Safety:** Protobuf definitions act as a cross-language contract, preventing schema mismatches between Python and Go.
- **Streaming:** Supports bidirectional streaming if we later need real-time event feeds from the Core to the API.
- **Efficiency:** Low overhead compared to traditional HTTP/1.1 REST calls.

## Consequences
- Requires Protobuf compilation steps in the build process.
- Requires gRPC libraries in both Python and Go environments.
