# TLS & Rate Limiting Design (v0.5)

## TLS Termination
The Go Core Engine will natively support TLS termination to inspect encrypted traffic.

### Configuration
`config.yaml` will be extended:
```yaml
server:
  listen: ":80"
  tls:
    enabled: true
    cert_file: "/etc/bwaf/certs/tls.crt"
    key_file: "/etc/bwaf/certs/tls.key"
    listen_tls: ":443"
```

### Implementation
- `main.go` will spawn two listeners (HTTP and HTTPS) if TLS is enabled.
- For v0.5 local dev, we will provide a self-signed certificate generator script or instructions.

## Rate Limiting
Rate limiting protects the backend from flood attacks and brute-force attempts.

### Strategy: Token Bucket (In-Memory)
- **Bucket Key**: Source IP.
- **Rate**: Configurable requests per second (RPS).
- **Burst**: Configurable maximum burst size.

### Configuration
```yaml
rate_limiting:
  enabled: true
  rps: 10
  burst: 20
```

### Implementation
- New `limiter` package in `pkg/`.
- Integration into `WAFProxy.ServeHTTP` as the first security check.

## Implementation Plan
1.  Extend `common.Config` and `loader.go`.
2.  Implement `lib/limiter` using `golang.org/x/time/rate`.
3.  Update `WAFProxy` to use the limiter.
4.  Update `main.go` for dual listening.
5.  Verification: Curl tests and unit tests for the limiter.
