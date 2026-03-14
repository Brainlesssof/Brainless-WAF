# Observability Guide

This guide covers how to monitor, log, and trace Brainless WAF in production — from Prometheus metrics to SIEM integration to distributed tracing.

---

## Table of Contents

- [Prometheus Metrics](#prometheus-metrics)
- [Grafana Dashboard](#grafana-dashboard)
- [Structured Logging](#structured-logging)
- [Audit Log](#audit-log)
- [SIEM Integration](#siem-integration)
- [Distributed Tracing (OpenTelemetry)](#distributed-tracing-opentelemetry)
- [Alerting](#alerting)

---

## Prometheus Metrics

Brainless WAF exposes Prometheus metrics at `http://host:9113/metrics`.

### Traffic Metrics

```promql
# Total requests (label: method, status, upstream)
brainless_requests_total

# Requests per second (5m rate)
rate(brainless_requests_total[5m])

# Total blocked requests (label: rule_id, severity, attack_type)
brainless_blocked_total

# Block rate percentage
100 * rate(brainless_blocked_total[5m])
    / rate(brainless_requests_total[5m])

# Request latency histogram (WAF overhead only, not backend)
brainless_latency_seconds_bucket
brainless_latency_seconds_sum
brainless_latency_seconds_count

# p50, p95, p99 latency
histogram_quantile(0.50, rate(brainless_latency_seconds_bucket[5m]))
histogram_quantile(0.95, rate(brainless_latency_seconds_bucket[5m]))
histogram_quantile(0.99, rate(brainless_latency_seconds_bucket[5m]))

# Active connections
brainless_connections_active

# Connection rate
rate(brainless_connections_total[1m])
```

### Rule Metrics

```promql
# Rules currently loaded
brainless_rules_total{enabled="true"}

# Rule evaluation time
histogram_quantile(0.99, rate(brainless_rule_evaluation_seconds_bucket[5m]))

# Top triggered rules (use with topk)
topk(10, sum by (rule_id, msg) (rate(brainless_rule_matches_total[1h])))

# Rate limit hit rate
rate(brainless_rate_limit_exceeded_total[5m])
```

### System Metrics

```promql
# Memory usage
process_resident_memory_bytes{job="brainless-waf"}

# CPU usage
rate(process_cpu_seconds_total{job="brainless-waf"}[1m])

# Go GC pause time
go_gc_duration_seconds{quantile="0.99"}

# Goroutines (should be stable, not growing)
go_goroutines{job="brainless-waf"}

# TLS certificate expiry (seconds until expiry)
brainless_tls_cert_expiry_seconds{domain="example.com"}
```

### Upstream Metrics

```promql
# Upstream response time (includes backend latency)
histogram_quantile(0.99, rate(brainless_upstream_latency_seconds_bucket[5m]))

# Upstream error rate
rate(brainless_upstream_errors_total[5m])

# Upstream availability
1 - (rate(brainless_upstream_errors_total[5m])
     / rate(brainless_requests_total[5m]))
```

---

## Grafana Dashboard

Import the pre-built Grafana dashboard from `deploy/grafana/dashboard.json`.

**Or install from Grafana.com marketplace:** Dashboard ID `XXXXX` (published at v1.0 release)

The dashboard includes:
- **Overview row:** req/s, block rate, p99 latency, active connections
- **Traffic row:** request volume chart, method breakdown, status code distribution
- **Security row:** blocks over time, top attack types, top blocked IPs, geographic map
- **Rules row:** top triggered rules, rule evaluation time, anomaly score distribution
- **System row:** CPU, memory, GC pauses, goroutines, upstream health

---

## Structured Logging

Brainless WAF outputs structured JSON logs by default. Each log line is a complete JSON object.

```json
{
  "level": "info",
  "time": "2025-03-14T12:34:56.789Z",
  "caller": "proxy/handler.go:142",
  "msg": "request processed",
  "request_id": "req_abc123",
  "method": "GET",
  "uri": "/api/products",
  "host": "api.example.com",
  "remote_addr": "203.0.113.42",
  "status": 200,
  "bytes_sent": 4892,
  "duration_ms": 12.3,
  "upstream_duration_ms": 11.5,
  "waf_duration_ms": 0.8,
  "rules_evaluated": 2140,
  "anomaly_score": 0,
  "action": "allow"
}
```

For blocked requests:
```json
{
  "level": "warn",
  "time": "2025-03-14T12:34:57.001Z",
  "msg": "request blocked",
  "request_id": "req_xyz789",
  "event_id": "evt_def456",
  "method": "GET",
  "uri": "/search?q=1+UNION+SELECT+1,2,3--",
  "remote_addr": "203.0.113.99",
  "status": 403,
  "anomaly_score": 10,
  "anomaly_threshold": 10,
  "triggered_rules": [942200],
  "action": "block",
  "attack_type": "SQLi"
}
```

### Log Levels

Configure in `config.yaml`:
```yaml
logging:
  level: info     # debug | info | warn | error
  format: json    # json | text (use text for local dev)
  output: stdout  # stdout | /path/to/file
```

In production, always use `json` format for SIEM compatibility.

---

## Audit Log

The audit log is a separate, append-only file containing a detailed record of every security event. Unlike the access log, the audit log includes the full matched payload.

```yaml
logging:
  audit_log: /var/log/brainless/audit.log
```

Audit log entry format:
```json
{
  "timestamp": "2025-03-14T12:34:57.001Z",
  "event_id": "evt_def456",
  "type": "SECURITY_EVENT",
  "action": "BLOCK",
  "severity": "CRITICAL",
  "anomaly_score": 10,
  "source": {
    "ip": "203.0.113.99",
    "port": 54321,
    "country": "XX",
    "user_agent": "sqlmap/1.7.8"
  },
  "request": {
    "method": "GET",
    "uri": "/search?q=1+UNION+SELECT+1,2,3--",
    "host": "api.example.com",
    "protocol": "HTTP/2"
  },
  "matched_rules": [
    {
      "rule_id": 942200,
      "phase": 2,
      "msg": "Detects MySQL comments and SQL injection",
      "severity": "CRITICAL",
      "score": 10,
      "variable": "ARGS:q",
      "matched_value": "1 UNION SELECT 1,2,3--",
      "tags": ["OWASP-A03", "SQLi", "UNION"]
    }
  ]
}
```

> **Privacy note:** The audit log contains matched request parameters which may include user data. Ensure your data retention policy and GDPR/privacy obligations are satisfied. The `audit_log_redact_fields` config option can remove specific parameters from audit logs.

---

## SIEM Integration

### Elasticsearch / OpenSearch

```yaml
logging:
  siem:
    enabled: true
    type: elasticsearch
    endpoint: https://your-elk-host:9200
    index: brainless-waf-events
    username: brainless
    password: "${ELASTIC_PASSWORD}"
    tls_verify: true
    batch_size: 200
    flush_interval: 5s
```

Events are indexed as JSON documents. Use the Kibana dashboard template in `deploy/kibana/`.

### Splunk

```yaml
logging:
  siem:
    enabled: true
    type: splunk
    endpoint: https://your-splunk-hec:8088
    token: "${SPLUNK_HEC_TOKEN}"
    index: brainless_waf
    source: brainless-waf
    sourcetype: _json
```

Splunk app and dashboards available in `deploy/splunk/`.

### Syslog (RFC 5424)

```yaml
logging:
  siem:
    enabled: true
    type: syslog
    endpoint: udp://siem.internal:514
    # Or TCP: tcp://siem.internal:514
    # Or TLS: tls://siem.internal:6514
    facility: 16   # local0
```

### Generic Webhook

```yaml
logging:
  siem:
    enabled: true
    type: webhook
    endpoint: https://your-siem.example.com/ingest
    headers:
      Authorization: "Bearer ${WEBHOOK_TOKEN}"
      Content-Type: application/json
    batch_size: 100
    flush_interval: 10s
```

---

## Distributed Tracing (OpenTelemetry)

Brainless WAF propagates trace context (W3C `traceparent` header) and can export spans to any OpenTelemetry-compatible backend.

```yaml
tracing:
  enabled: true
  exporter: otlp             # otlp | jaeger | zipkin
  endpoint: http://otel-collector:4317
  sample_rate: 0.01          # Sample 1% of requests (use 1.0 for debug)
  service_name: brainless-waf
```

Each WAF trace includes spans for:
- `waf.tls_handshake` — TLS termination time
- `waf.request_parse` — request normalization
- `waf.rule_evaluate.phase1` — phase 1 rule evaluation
- `waf.rule_evaluate.phase2` — phase 2 rule evaluation
- `waf.upstream_proxy` — time waiting for backend response
- `waf.rule_evaluate.phase4` — response rule evaluation (if enabled)

The `traceparent` header from the client is honored and propagated to the upstream backend, enabling end-to-end trace correlation.

---

## Alerting

### Recommended Prometheus Alerts

```yaml
# prometheus-rules.yml
groups:
  - name: brainless-waf.critical
    rules:
      - alert: WAFDown
        expr: up{job="brainless-waf"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Brainless WAF instance {{ $labels.instance }} is down"

      - alert: WAFHighP99Latency
        expr: histogram_quantile(0.99, rate(brainless_latency_seconds_bucket[5m])) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "WAF p99 latency is {{ $value | humanizeDuration }} (threshold: 50ms)"

      - alert: WAFCertExpiringSoon
        expr: brainless_tls_cert_expiry_seconds < 86400 * 14
        labels:
          severity: warning
        annotations:
          summary: "TLS cert for {{ $labels.domain }} expires in {{ $value | humanizeDuration }}"

      - alert: WAFHighBlockRate
        expr: |
          rate(brainless_blocked_total[10m])
          / rate(brainless_requests_total[10m]) > 0.20
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "WAF is blocking {{ $value | humanizePercentage }} of traffic — possible attack or misconfiguration"

      - alert: WAFUpstreamErrors
        expr: rate(brainless_upstream_errors_total[5m]) > 5
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "WAF upstream error rate is {{ $value }}/s — backend may be down"

      - alert: WAFMemoryHigh
        expr: process_resident_memory_bytes{job="brainless-waf"} > 2e9
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "WAF memory usage is {{ $value | humanize1024 }}B — possible memory leak"
```
