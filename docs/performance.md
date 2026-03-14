# Performance Tuning Guide

This guide covers how to measure, understand, and optimize Brainless WAF performance for your specific workload.

---

## Table of Contents

- [Performance Baselines](#performance-baselines)
- [Key Metrics](#key-metrics)
- [Profiling](#profiling)
- [CPU Tuning](#cpu-tuning)
- [Memory Tuning](#memory-tuning)
- [Latency Tuning](#latency-tuning)
- [Rule Performance](#rule-performance)
- [TLS Performance](#tls-performance)
- [Rate Limiter Performance](#rate-limiter-performance)
- [High-Traffic Deployments](#high-traffic-deployments)
- [Benchmarking Your Setup](#benchmarking-your-setup)

---

## Performance Baselines

Reference benchmarks on `c5.2xlarge` (8 vCPU, 16GB RAM, Linux 6.1):

| Scenario | Throughput | p50 Latency | p99 Latency | Memory |
|----------|-----------|-------------|-------------|--------|
| Idle | — | — | — | 45 MB |
| 1,000 req/s (CRS paranoia 2) | 1,000 req/s | 0.3ms | 1.1ms | 120 MB |
| 5,000 req/s (CRS paranoia 2) | 5,000 req/s | 0.6ms | 2.4ms | 210 MB |
| 10,000 req/s (CRS paranoia 2) | 10,000 req/s | 0.8ms | 3.2ms | 380 MB |
| 10,000 req/s (paranoia 4) | 8,200 req/s | 1.1ms | 5.8ms | 520 MB |
| 10,000 req/s + response body scan | 6,400 req/s | 2.3ms | 12ms | 850 MB |

These numbers are for HTTP (no TLS). Add ~0.3ms per request for TLS 1.3 termination on modern hardware.

---

## Key Metrics

Monitor these Prometheus metrics to understand WAF performance:

```promql
# Request throughput
rate(brainless_requests_total[1m])

# Block rate
rate(brainless_blocked_total[1m]) / rate(brainless_requests_total[1m])

# Latency percentiles
histogram_quantile(0.50, brainless_latency_seconds_bucket)
histogram_quantile(0.95, brainless_latency_seconds_bucket)
histogram_quantile(0.99, brainless_latency_seconds_bucket)

# Rule evaluation time
histogram_quantile(0.99, brainless_rule_evaluation_seconds_bucket)

# Active connections
brainless_connections_active

# Memory usage
process_resident_memory_bytes{job="brainless-waf"}

# CPU usage
rate(process_cpu_seconds_total{job="brainless-waf"}[1m])
```

Set up alerts for:
- `p99 > 10ms` — investigate bottleneck
- `p99 > 50ms` — critical, likely a blocking issue
- `memory > 2GB` — possible memory leak
- `CPU > 80%` — consider scaling out

---

## Profiling

### CPU Profiling

```bash
# Enable pprof (set in config or environment)
BRAINLESS_PPROF=true brainless-waf --config config/local.yaml

# Sample CPU for 30 seconds
curl http://localhost:6060/debug/pprof/profile?seconds=30 -o cpu.prof

# Open interactive flamegraph
go tool pprof -http=:8888 cpu.prof
```

Look for: rule evaluation functions, regex matching, JSON/body parsing.

### Memory Profiling

```bash
# Heap snapshot
curl http://localhost:6060/debug/pprof/heap -o heap.prof
go tool pprof -http=:8888 heap.prof
```

Look for: large allocations in transaction variables, cached regex patterns growing unexpectedly.

### Goroutine Profiling

```bash
# If goroutines are piling up (connection leaks, blocking operations)
curl http://localhost:6060/debug/pprof/goroutine?debug=2
```

---

## CPU Tuning

### Worker Threads

By default, Brainless WAF uses all available CPU cores. In containerized environments, set this explicitly:

```yaml
advanced:
  worker_threads: 4  # Match your CPU limit
```

If running in a container with a CPU limit (e.g., Kubernetes `resources.limits.cpu: 2`), set `worker_threads: 2` to avoid Go's scheduler fighting over CPU time with the container runtime.

### GOMAXPROCS

For containers, use the `automaxprocs` library (already built into Brainless WAF) which reads `cgroups` CPU quota and sets `GOMAXPROCS` automatically. No manual tuning needed.

### Paranoia Level Impact on CPU

| Paranoia Level | Rules Active | CPU Impact |
|---------------|-------------|------------|
| 1 | ~800 | Baseline |
| 2 | ~2,100 | +15% |
| 3 | ~3,900 | +35% |
| 4 | ~5,200 | +55% |

If you need paranoia 3 or 4 but can't afford the CPU, consider applying higher paranoia levels only to specific paths:

```yaml
rules:
  per_path_paranoia:
    - path: /admin/
      paranoia_level: 4
    - path: /api/
      paranoia_level: 3
    - path: /
      paranoia_level: 2
```

---

## Memory Tuning

### Memory Breakdown (at 5,000 req/s)

| Component | Memory |
|-----------|--------|
| Rule set (compiled regex) | ~80 MB |
| Active connection buffers | ~60 MB |
| Rate limit counters (in-memory) | ~20 MB |
| Transaction variable pools | ~30 MB |
| **Total** | ~190 MB |

### Request Body Buffer Limit

Brainless WAF buffers request bodies for phase 2 inspection. Large uploads consume significant memory:

```yaml
advanced:
  request_body_limit: 10485760   # 10MB (reduce from 128MB default if uploads aren't needed)
```

For APIs with large JSON payloads, consider phase 2 inspection only on endpoints that need it:

```yaml
rules:
  phase2_paths:
    - /api/auth/
    - /api/admin/
  # Other paths skip phase 2 body inspection
```

### Response Body Inspection

Response body inspection (phase 4/5) is the most memory-intensive feature — it buffers the full response:

```yaml
detection:
  response_body_inspection: false    # Disabled by default — only enable if needed
  response_body_limit: 1048576       # Cap at 1MB per response when enabled
```

---

## Latency Tuning

### Sources of Latency

| Source | Typical contribution | How to reduce |
|--------|---------------------|---------------|
| Rule evaluation | 0.3–2ms | Lower paranoia level, optimize custom rules |
| TLS handshake | 0.2–0.5ms | TLS session resumption (enabled by default) |
| Rate limit lookup | 0.05ms (memory) / 1ms (Redis) | Use in-memory for single instance |
| Body parsing | 0.1–0.5ms | Reduce `request_body_limit` |
| Response body scan | 1–10ms | Disable unless necessary |
| Backend response | Varies | This is your app's problem, not the WAF's |

### Connection Keep-Alive

Ensure keep-alive is enabled on both sides:

```yaml
server:
  idle_timeout: 120s     # Keep client connections alive
  upstream_keepalive: true
  upstream_keepalive_idle: 60s
  upstream_max_idle_conns: 1000
```

### TLS Session Resumption

Enabled by default. Session tickets reduce TLS handshake time for returning clients from ~0.5ms to ~0.1ms. No configuration needed.

---

## Rule Performance

### Regex Performance

The #1 cause of slow rule evaluation is poorly written regular expressions. Watch for:

**Catastrophic backtracking:** Regex like `(a+)+b` on long input can take exponential time.

```bash
# The WAF detects and warns about risky patterns at startup
grep "regex performance warning" /var/log/brainless/brainless.log
```

**Use `@pm` instead of `@rx` for keyword lists:**

```
# Slow: many @rx rules for individual keywords
SecRule ARGS "@rx (?i)union" "..."
SecRule ARGS "@rx (?i)select" "..."
SecRule ARGS "@rx (?i)insert" "..."

# Fast: single @pm rule with all keywords
SecRule ARGS "@pmFromFile /etc/brainless/rules/sql_keywords.txt" "..."
```

`@pm` uses the Aho-Corasick algorithm — it matches thousands of patterns in a single pass through the input, as fast as matching one pattern.

### Rule Ordering

Rules are evaluated in order. Put fast, common rules first:

1. IP-based rules (`@ipMatch`) — fastest, evaluated first
2. URI-based rules (`@rx` on `REQUEST_URI`) — very fast, small target
3. Header rules — medium speed
4. Body rules (`ARGS`, `REQUEST_BODY`) — slowest, put last

### Phase Optimization

Put blocking rules in phase 1 (request headers) where possible. If you can block a request before reading the body, you save body parsing time:

```
# Good: block by IP in phase 1, before body is read
SecRule REMOTE_ADDR "@ipMatchFromFile blocklist.txt" "id:...,phase:1,deny"

# Less optimal: same rule in phase 2 (body already parsed unnecessarily)
SecRule REMOTE_ADDR "@ipMatchFromFile blocklist.txt" "id:...,phase:2,deny"
```

---

## TLS Performance

### Session Tickets vs Session IDs

Both are enabled by default. Session tickets are preferred for multi-instance deployments (stateless).

For multi-instance, ensure all WAF instances share the same session ticket key:

```yaml
server:
  tls:
    session_ticket_key: /etc/brainless/tls-session-key.bin  # Shared across instances
```

Generate a shared key:
```bash
openssl rand 48 > /etc/brainless/tls-session-key.bin
# Copy to all WAF instances
```

### TLS 1.3 Only

For maximum performance and security, restrict to TLS 1.3 (at the cost of compatibility with old clients):

```yaml
server:
  tls:
    min_version: "1.3"
```

TLS 1.3 is ~30% faster handshake than 1.2 due to fewer round trips.

---

## Rate Limiter Performance

### In-Memory vs Redis

| Storage | Latency | Suitable for |
|---------|---------|-------------|
| In-memory | ~0.05ms | Single WAF instance |
| Redis (local) | ~0.5ms | Multi-instance, same datacenter |
| Redis (remote) | ~1–5ms | Multi-instance, different regions |

For high-throughput deployments, in-memory rate limiting with Redis as a secondary sync is the best trade-off:

```yaml
rate_limiting:
  storage: memory
  redis_sync:
    enabled: true
    redis_url: redis://localhost:6379
    sync_interval: 1s   # Sync counters to Redis every second
```

This gives in-memory speed for the hot path while sharing state across instances with 1-second eventual consistency.

---

## High-Traffic Deployments

### Horizontal Scaling

When a single WAF instance reaches CPU saturation (>80% at p99):

1. Deploy 2+ WAF core instances behind a load balancer
2. Switch rate limiting to Redis storage (shared counters)
3. Ensure Management API can push config to all instances simultaneously

```yaml
# In Management API config
core_instances:
  - grpc://waf-core-1:9091
  - grpc://waf-core-2:9091
  - grpc://waf-core-3:9091
```

### Load Balancer Configuration

The load balancer in front of WAF instances should:
- Use **least-connections** algorithm (not round-robin) — requests vary in processing time
- Enable **health check** on `GET /health`
- Preserve client IP via `X-Forwarded-For` (configure WAF's `trusted_proxies` accordingly)

```nginx
# nginx upstream example
upstream brainless_waf {
    least_conn;
    server waf-1:80;
    server waf-2:80;
    server waf-3:80;
    keepalive 100;
}
```

### Kernel Tuning for High Connection Counts

```bash
# /etc/sysctl.d/99-brainless-waf.conf
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 65535
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.ip_local_port_range = 1024 65535
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_fin_timeout = 15
fs.file-max = 1000000

# Apply
sysctl -p /etc/sysctl.d/99-brainless-waf.conf

# Set file descriptor limit for the WAF process
ulimit -n 1000000
# Or in systemd service: LimitNOFILE=1000000
```

---

## Benchmarking Your Setup

### Quick benchmark with `wrk`

```bash
# Install wrk
apt install wrk

# Benchmark a safe endpoint (should show WAF overhead)
wrk -t8 -c400 -d30s http://your-waf-host/

# Compare with direct backend (no WAF)
wrk -t8 -c400 -d30s http://your-backend-host/

# The difference is WAF overhead
```

### Benchmark with attack traffic mixed in

```bash
# Use the included benchmark script
scripts/benchmark.sh \
  --target http://your-waf-host \
  --duration 60s \
  --concurrency 200 \
  --attack-mix 0.05   # 5% of requests will be attack payloads
```

This gives you a realistic picture: throughput, latency, and block rate under mixed traffic.

### Load test profile

Run at least these scenarios before going to production:

1. **Ramp test** — gradually increase from 0 to 150% of expected peak over 10 minutes
2. **Soak test** — run at 80% of expected peak for 2 hours (catch memory leaks)
3. **Spike test** — jump from 10% to 200% of expected peak instantly (test burst handling)
4. **Attack flood** — send 100% attack traffic at expected peak rate (verify blocking doesn't crash the WAF)
