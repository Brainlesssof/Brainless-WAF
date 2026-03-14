# Configuration Reference

Complete reference for all Brainless WAF configuration options.

Configuration is loaded from a YAML file (default: `config/local.yaml`). Environment variables override config file values using the pattern `BRAINLESS_SECTION_KEY` (e.g., `BRAINLESS_SERVER_UPSTREAM`).

---

## Full Configuration Example

```yaml
# ─── SERVER ───────────────────────────────────────────────────────────────────
server:
  listen: 0.0.0.0:80            # HTTP listen address
  listen_tls: 0.0.0.0:443       # HTTPS listen address (omit to disable TLS)
  upstream: http://backend:8080  # Backend server to proxy to
  upstream_timeout: 30s          # Timeout waiting for upstream response
  max_connections: 10000         # Maximum concurrent connections
  read_timeout: 10s              # Client read timeout
  write_timeout: 30s             # Client write timeout
  idle_timeout: 120s             # Keep-alive idle timeout
  real_ip_header: X-Forwarded-For
  trusted_proxies:
    - 127.0.0.0/8
    - 10.0.0.0/8
    - 172.16.0.0/12

  tls:
    cert: ""                    # Path to TLS certificate (PEM)
    key: ""                     # Path to TLS private key (PEM)
    min_version: "1.2"          # Minimum TLS version: 1.2 or 1.3
    ciphers: []                 # Custom cipher list (empty = secure defaults)
    acme: false                 # Enable automatic Let's Encrypt certificates
    acme_email: ""              # Required if acme: true
    acme_domains: []            # Required if acme: true
    acme_cache_dir: /var/lib/brainless/acme

# ─── DETECTION ────────────────────────────────────────────────────────────────
detection:
  mode: learning                # learning | detect | block
  paranoia_level: 2             # 1 (permissive) to 4 (strict)
  anomaly_threshold: 10         # Request blocked when score >= this value
  inbound_anomaly_threshold: 10 # Alias for anomaly_threshold
  outbound_anomaly_threshold: 4 # Threshold for response rule scoring
  ml_bot_detection: true        # Enable ML-based bot detection
  response_body_inspection: false  # Enable phase 4/5 (adds latency)
  response_body_limit: 1048576  # Max response body size to inspect (bytes)

  # Per-severity score values
  severity_scores:
    critical: 10
    error: 5
    warning: 3
    notice: 1

# ─── RULES ────────────────────────────────────────────────────────────────────
rules:
  crs_enabled: true             # Load bundled OWASP CRS 4.x rules
  crs_paranoia_level: 2         # CRS paranoia level (overrides detection.paranoia_level for CRS)
  custom_rules_dir: /etc/brainless/rules/custom/
  scripts_dir: /etc/brainless/rules/scripts/
  auto_update: true             # Auto-pull rule updates
  update_interval: 6h           # How often to check for updates
  update_url: https://rules.brainless-security.io/v1/rules.tar.gz
  update_verify_signature: true # Verify rule update signatures

# ─── RATE LIMITING ────────────────────────────────────────────────────────────
rate_limiting:
  enabled: true
  default_rps: 100              # Requests per second per IP
  burst: 200                    # Burst capacity above default_rps
  per_endpoint:                 # Per-endpoint overrides
    - path: /api/auth/login
      methods: [POST]
      rps: 5
      burst: 10
    - path: /api/
      rps: 50
      burst: 100
  storage: memory               # memory | redis (use redis for multi-instance)
  redis_url: ""                 # Required if storage: redis

# ─── IP LISTS ─────────────────────────────────────────────────────────────────
ip_lists:
  blocklist:
    - 198.51.100.0/24           # CIDR ranges are supported
    - 203.0.113.42
  allowlist:
    - 127.0.0.1
    - 10.0.0.0/8
  block_tor_exits: false        # Block known Tor exit nodes (updated daily)
  block_known_scanners: false   # Block known vulnerability scanner IPs

# ─── LOGGING ──────────────────────────────────────────────────────────────────
logging:
  level: info                   # debug | info | warn | error
  format: json                  # json | text
  output: stdout                # stdout | /path/to/file.log
  audit_log: /var/log/brainless/audit.log  # Detailed event log (set empty to disable)
  rotate:
    enabled: true
    max_size_mb: 100
    max_backups: 10
    max_age_days: 30

  # SIEM / log forwarding
  siem:
    enabled: false
    type: elasticsearch          # elasticsearch | splunk | syslog | kafka
    endpoint: http://elk:9200
    index: brainless-waf-logs
    tls_verify: true
    batch_size: 100
    flush_interval: 5s

# ─── METRICS ──────────────────────────────────────────────────────────────────
metrics:
  enabled: true
  listen: 0.0.0.0:9113         # Prometheus scrape endpoint
  path: /metrics

# ─── MANAGEMENT API ───────────────────────────────────────────────────────────
management:
  grpc_listen: 127.0.0.1:9091  # Internal gRPC server (management API connects here)
  grpc_tls: true                # Use mTLS for management gRPC channel

# ─── ADVANCED ─────────────────────────────────────────────────────────────────
advanced:
  worker_threads: 0             # 0 = use all available CPU cores
  connection_buffer_size: 65536 # TCP buffer size per connection (bytes)
  request_body_limit: 134217728 # Max request body to buffer (128MB default)
  uri_limit: 8192               # Max URI length (bytes)
  header_limit: 32768           # Max total header size (bytes)
  fail_open: false              # If true, pass traffic on internal errors (NOT recommended)
```

---

## Key Options Explained

### `detection.mode`

| Mode | Behavior |
|------|---------|
| `learning` | Passively observes traffic. Logs all rule matches but never blocks. Use for initial deployment to understand your traffic baseline. |
| `detect` | Logs matching requests and adds detection headers but does NOT block. Use to tune rules and identify false positives before going live. |
| `block` | Actively blocks requests that exceed the anomaly threshold. Production mode. |

**Recommended workflow:** `learning` (72h) → `detect` (48h, review FPs) → `block`

### `detection.paranoia_level`

Controls how strict the OWASP CRS rules are. Higher levels add more rules but more false positives.

| Level | Description | Recommended for |
|-------|-------------|----------------|
| 1 | Only the most reliable rules. Very few false positives. | New deployments, legacy applications |
| 2 | Good balance of protection vs. false positives. | Most production deployments |
| 3 | Adds rules that may trigger on unusual-but-legitimate traffic. | High-security environments |
| 4 | Maximum detection. Expect false positives to tune. | WAF testing, highly targeted applications |

### `rate_limiting.storage`

- `memory` — rate limit counters stored in the WAF process memory. Simple, no dependencies, but counters reset on restart and are not shared between multiple WAF instances.
- `redis` — counters stored in Redis. Required for multi-instance deployments to avoid attackers bypassing limits by hitting different WAF nodes.

### `server.trusted_proxies`

If your WAF is behind a load balancer or CDN, set `trusted_proxies` to the IP ranges of your upstream proxy infrastructure. The WAF will then trust `X-Forwarded-For` (or `real_ip_header`) from these sources when determining the real client IP for rate limiting and logging.

**Security warning:** Do not add broad ranges like `0.0.0.0/0`. Only add the IPs of your own proxy infrastructure.

### `advanced.fail_open`

**Default: `false`** — on internal errors, return 502 rather than passing traffic through.

Setting `fail_open: true` means the WAF will pass traffic to the backend even when it cannot inspect it. This prioritizes availability over security. Do not use in security-critical deployments.

---

## Environment Variable Overrides

All config values can be overridden via environment variables. The format is:
`BRAINLESS_<SECTION>_<KEY>` (uppercase, underscores replace dots and hyphens).

```bash
BRAINLESS_SERVER_UPSTREAM=http://backend:9000
BRAINLESS_DETECTION_MODE=block
BRAINLESS_DETECTION_ANOMALY_THRESHOLD=15
BRAINLESS_LOGGING_LEVEL=debug
DATABASE_URL=postgresql://user:pass@host:5432/db
JWT_SECRET=your-secret-key-minimum-32-characters
```

Environment variables take precedence over the config file. Useful for container deployments where secrets should not be in config files.

---

## Secrets Management

Never put secrets directly in the config file if it will be committed to version control.

**Recommended approaches:**

```bash
# Docker secrets
echo "your-jwt-secret" | docker secret create jwt_secret -

# Kubernetes secrets
kubectl create secret generic brainless-secrets \
  --from-literal=jwt-secret=your-secret \
  --from-literal=db-password=your-password

# HashiCorp Vault (v1.2+)
vault write secret/brainless-waf jwt_secret=your-secret
# Set in config: jwt_secret: "vault:secret/brainless-waf#jwt_secret"
```
