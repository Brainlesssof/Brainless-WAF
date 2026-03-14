# Troubleshooting Guide

This guide covers common issues and how to resolve them. If your issue isn't listed, check [GitHub Issues](https://github.com/brainless-security/brainless-waf/issues) or ask in [Discord](https://discord.gg/brainless-waf).

---

## Table of Contents

- [Diagnosing Issues](#diagnosing-issues)
- [Installation Issues](#installation-issues)
- [Traffic & Blocking Issues](#traffic--blocking-issues)
- [False Positives](#false-positives)
- [Performance Issues](#performance-issues)
- [TLS / Certificate Issues](#tls--certificate-issues)
- [Dashboard Issues](#dashboard-issues)
- [Rule Issues](#rule-issues)
- [Kubernetes Issues](#kubernetes-issues)
- [Log Reference](#log-reference)

---

## Diagnosing Issues

### Start here — check the health endpoint

```bash
curl http://localhost/health
# Expected: {"status":"ok","version":"1.0.0","rules_loaded":4523}

# If the WAF isn't responding:
docker compose ps          # Are containers running?
docker compose logs waf-core --tail=50
```

### Enable debug logging temporarily

```yaml
# config/local.yaml
logging:
  level: debug
```

```bash
docker compose restart waf-core
docker compose logs waf-core -f | grep -E "ERROR|WARN|blocked"
```

Remember to set `level: info` again after debugging — debug logging is very verbose.

### Identify why a specific request was blocked

Every blocked request gets an event ID in the response body:
```json
{"error": "Request blocked by security policy", "event_id": "evt_abc123xyz"}
```

Look up the event:
```bash
# Via CLI
brainless-ctl event show evt_abc123xyz

# Via API
curl http://localhost:8000/api/v1/events/evt_abc123xyz \
  -H "Authorization: Bearer <token>"
```

This shows which rule triggered, the matched value, and the full request details.

---

## Installation Issues

### Docker containers exit immediately

```bash
docker compose logs waf-core
```

Common causes:

**Config file not found:**
```
FATAL config file not found: /etc/brainless/config.yaml
```
→ Ensure `config/local.yaml` exists and is mounted correctly.

**Upstream not reachable at startup:**
```
WARN upstream health check failed: connection refused
```
→ This is a warning, not fatal. The WAF starts and retries. Check that your backend is running.

**Port already in use:**
```
FATAL listen tcp :80: bind: address already in use
```
→ Something else is using port 80. Find and stop it: `sudo lsof -i :80`

### Permission denied on port 80/443

Non-root processes cannot bind to ports below 1024 on Linux by default.

```bash
# Option 1: Use Docker (handles this automatically)
docker compose up -d

# Option 2: Grant the binary permission
sudo setcap 'cap_net_bind_service=+ep' /usr/local/bin/brainless-waf

# Option 3: Use a higher port + iptables redirect
# In config: listen: 0.0.0.0:8080
sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8080
```

---

## Traffic & Blocking Issues

### Legitimate requests are being blocked

See [False Positives](#false-positives) below.

### All requests return 502 Bad Gateway

The WAF cannot reach your backend.

```bash
# Test connectivity from inside the container
docker exec brainless-waf-core wget -qO- http://your-backend:8080/health

# Check the upstream setting in config
grep upstream config/local.yaml

# If using Docker Compose, ensure services are on the same network
docker network ls
docker inspect brainless-waf-core | grep NetworkMode
```

### Real client IPs are showing as the load balancer IP

Configure trusted proxy headers:

```yaml
# config/local.yaml
server:
  trusted_proxies:
    - 10.0.0.0/8
    - 172.16.0.0/12
  real_ip_header: X-Forwarded-For  # or X-Real-IP
```

### Requests are slow (high latency)

See [Performance Issues](#performance-issues).

### WAF blocks requests in `learning` or `detect` mode

In `detect` mode, the WAF logs but does NOT block. If you're seeing blocks, check:
1. You may still be in `block` mode — verify: `brainless-ctl config get detection.mode`
2. Some rules have `allow` or `deny` regardless of mode (IP blocklist always blocks)

---

## False Positives

False positives are legitimate requests incorrectly blocked by WAF rules. They are normal when starting out, especially at higher paranoia levels.

### Step 1 — Identify the triggering rule

Look up the blocked event (see [Diagnosing Issues](#diagnosing-issues)).

The event will show:
```json
{
  "rule_id": 942200,
  "rule_msg": "Detects MySQL comments and classic SQL injection attempts",
  "matched_var": "ARGS:description",
  "matched_value": "SELECT option from dropdown",
  "anomaly_score": 10
}
```

### Step 2 — Decide: tune the rule or add an exception

**Option A: Create an exception rule** (recommended for specific cases)

```
# Exclude a specific parameter from a rule
SecRuleUpdateTargetById 942200 "!ARGS:description"

# Exclude an entire URL path from a rule
SecRule REQUEST_URI "@beginsWith /editor" \
    "id:900100,phase:1,allow,nolog"

# Exclude a rule for a specific user agent (e.g., internal tool)
SecRule REQUEST_HEADERS:User-Agent "@contains MyInternalTool/1.0" \
    "id:900101,phase:1,ctl:ruleRemoveById=942200"
```

**Option B: Lower the paranoia level**

```yaml
detection:
  paranoia_level: 1  # Down from 2
```

**Option C: Raise the anomaly threshold**

```yaml
detection:
  anomaly_threshold: 15  # Up from 10
```

### Step 3 — Use learning mode

Switch back to `learning` mode if false positive rate is high:

```bash
brainless-ctl config set detection.mode learning
# Wait 72 hours
brainless-ctl learning-report --output exceptions.rules
# Review and apply generated exceptions
```

---

## Performance Issues

### High CPU usage

```bash
# Profile the WAF (requires debug build or pprof endpoint)
curl http://localhost:6060/debug/pprof/profile?seconds=30 -o cpu.prof
go tool pprof -http=:8888 cpu.prof
```

Common causes:
- **Too many regex rules:** Check if you have a large number of complex `@rx` rules. Use `@pm` (multi-pattern) for simple keyword lists instead.
- **Response body scanning enabled:** Phase 4/5 rules buffer the entire response. Disable for endpoints with large responses.
- **High connection count:** Increase `server.max_connections` and check for connection leaks.

### High memory usage

```bash
docker stats brainless-waf-core
```

If memory grows without bound (memory leak):
1. Check for large `IP:` variable sets with long TTLs
2. Enable the memory profiler: `BRAINLESS_PPROF=true` and sample at `localhost:6060/debug/pprof/heap`
3. Report as a bug with the profile attached

### Latency spikes

```bash
# Check latency histogram
curl http://localhost/metrics | grep brainless_latency
```

If p99 is high but p50 is normal, look for:
- Slow upstream backend responses (WAF can't reduce backend latency)
- Phase 4/5 (response body scanning) on a slow endpoint
- Redis connectivity issues (rate limit counter lookups timing out)

---

## TLS / Certificate Issues

### Certificate not renewing (Let's Encrypt)

```bash
# Check ACME status
brainless-ctl tls status

# Force renewal
brainless-ctl tls renew

# Common causes:
# - Port 80 not accessible from the internet (ACME HTTP challenge requires it)
# - DNS not pointing to this server
# - Rate limited by Let's Encrypt (max 5 failures per account per hour)
```

### Mixed content warnings in browser

Ensure your backend is also using HTTPS, or that the WAF rewrites `http://` links in responses:

```yaml
server:
  rewrite_redirects: true
  force_https: true
```

### TLS handshake failures

```bash
# Test TLS configuration
openssl s_client -connect your-domain:443 -servername your-domain

# Check minimum TLS version (should be 1.2)
grep tls_min_version config/local.yaml
```

---

## Dashboard Issues

### Cannot log in

1. Check Management API is running: `curl http://localhost:8000/health`
2. Reset admin password: `brainless-ctl passwd admin`
3. Clear browser cookies and try again

### Dashboard shows stale data

The dashboard polls the API every 30 seconds by default. For live updates, ensure WebSocket connection is working:

- Open browser DevTools → Network → WS tab
- You should see a WebSocket connection to `/api/v1/ws/metrics`
- If blocked, check your reverse proxy is configured to pass WebSocket upgrades:

```nginx
# nginx example
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "upgrade";
```

### 403 Forbidden when accessing dashboard

The WAF may be blocking its own dashboard traffic. Add an allowlist rule:

```
SecRule REQUEST_URI "@beginsWith /dashboard" \
    "id:900200,phase:1,allow,nolog"
```

---

## Rule Issues

### Rule not loading

```bash
brainless-ctl rules validate /path/to/my.rules

# Common errors:
# - Missing required action (id, phase)
# - Rule ID already in use
# - Invalid regex syntax
# - Missing closing quote
```

### Rule not triggering when expected

```bash
# Test a specific request against your rule
brainless-ctl rule-test \
  --rules /path/to/my.rules \
  --request tests/my_test_request.http \
  --verbose
```

`--verbose` shows the full variable extraction and every operator evaluation.

### Rule causes too many false positives after deployment

```bash
# Disable a specific rule without removing it
brainless-ctl rule disable 942200

# Re-enable
brainless-ctl rule enable 942200

# Or add a targeted exception (better than disabling entirely)
SecRuleUpdateTargetById 942200 "!ARGS:safe_field"
```

---

## Kubernetes Issues

### WAF pods crash-looping

```bash
kubectl describe pod -n brainless-waf <pod-name>
kubectl logs -n brainless-waf <pod-name> --previous
```

Common causes:
- `DATABASE_URL` secret not created or incorrect
- ConfigMap not mounted correctly
- Insufficient memory (OOMKilled) — increase `resources.limits.memory`

### Rules not updating in pods after config change

Rules are pushed via gRPC when you update through the API or dashboard. Verify:
```bash
# Check gRPC connectivity from management to core pods
kubectl exec -n brainless-waf deploy/brainless-waf-management -- \
  grpc_health_probe -addr=brainless-waf-core:9091
```

### Ingress traffic not reaching WAF

If using Ingress Controller mode, check annotations:
```yaml
metadata:
  annotations:
    brainless.io/waf: "enabled"
    brainless.io/paranoia-level: "2"
```

---

## Log Reference

### Log levels

| Level | Meaning |
|-------|---------|
| `DEBUG` | Very verbose — request/response details, rule evaluations |
| `INFO` | Normal operation — startup, config changes, request summaries |
| `WARN` | Potential issues — backend timeouts, slow rules, config warnings |
| `ERROR` | Failures that affect requests — backend unavailable, rule parse errors |
| `FATAL` | Startup failures — config invalid, port bind failed |

### Common log messages

| Message | Meaning | Action |
|---------|---------|--------|
| `upstream health check failed` | Backend not reachable | Check backend is running |
| `rule evaluation timeout` | A rule took too long | Check for catastrophic backtracking in regex |
| `anomaly score exceeded threshold` | Request blocked | Normal — investigate if legitimate |
| `rule file reloaded` | Hot-reload succeeded | Informational |
| `certificate renewed` | TLS cert auto-renewed | Informational |
| `rate limit exceeded` | IP hit rate limit | Check if attacker or misconfigured client |
| `gRPC config push failed` | Management → core sync failed | Check gRPC connectivity |

### Audit log format

```json
{
  "timestamp": "2025-03-14T12:34:56.789Z",
  "event_id": "evt_abc123",
  "action": "block",
  "severity": "CRITICAL",
  "rule_id": 942200,
  "rule_msg": "Detects MySQL comments and SQL injection",
  "src_ip": "203.0.113.42",
  "src_country": "XX",
  "method": "GET",
  "uri": "/search?q=1+UNION+SELECT+1,2,3--",
  "host": "example.com",
  "user_agent": "Mozilla/5.0",
  "anomaly_score": 10,
  "matched_var": "ARGS:q",
  "matched_value": "1 UNION SELECT 1,2,3--",
  "upstream_response_time_ms": null,
  "tags": ["OWASP-A03", "SQLi", "UNION"]
}
```

---

*Still stuck? Join [Discord](https://discord.gg/brainless-waf) or [open an issue](https://github.com/brainless-security/brainless-waf/issues/new?template=bug_report.md).*
