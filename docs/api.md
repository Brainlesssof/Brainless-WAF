# REST API Reference

The Brainless WAF Management API provides full programmatic control over rules, configuration, events, and users.

**Base URL:** `https://your-waf-host/api/v1`  
**Interactive docs:** `https://your-waf-host/api/v1/docs` (OpenAPI 3.1 / Swagger UI)

---

## Authentication

All endpoints (except `/health` and `/auth/token`) require authentication.

### Get a token

```bash
curl -X POST https://waf.example.com/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "yourpassword"}'
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "bearer",
  "expires_in": 3600
}
```

### Use the token

```bash
curl https://waf.example.com/api/v1/rules \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### API Keys

For programmatic/automation use, create an API key (does not expire):

```bash
# Create key
curl -X POST https://waf.example.com/api/v1/auth/api-keys \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "CI/CD Pipeline", "scopes": ["rules:read", "rules:write"]}'

# Response includes the key — save it, it won't be shown again
# Use it just like a Bearer token:
curl -H "Authorization: Bearer bwaf_live_xxxxxxxxxxxx" ...
```

---

## Endpoints

### Health

#### `GET /health`

Check WAF status. No authentication required.

```bash
curl https://waf.example.com/api/v1/health
```

Response `200 OK`:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime_seconds": 86400,
  "rules_loaded": 4523,
  "mode": "block",
  "upstream_healthy": true
}
```

---

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/auth/token` | Get JWT access token |
| `POST` | `/auth/refresh` | Refresh an expiring token |
| `POST` | `/auth/change-password` | Change your password |
| `GET` | `/auth/api-keys` | List your API keys |
| `POST` | `/auth/api-keys` | Create a new API key |
| `DELETE` | `/auth/api-keys/{id}` | Revoke an API key |

---

### Rules

#### `GET /rules`

List all rules.

```bash
curl https://waf.example.com/api/v1/rules \
  -H "Authorization: Bearer <token>"

# Filter by enabled/disabled
curl ".../rules?enabled=true&tag=OWASP-A03&limit=50&offset=0"
```

Response:
```json
{
  "total": 4523,
  "items": [
    {
      "id": "rule_abc123",
      "rule_id": 942200,
      "enabled": true,
      "source": "crs",
      "phase": 2,
      "msg": "Detects MySQL comments and SQL injection",
      "tags": ["OWASP-A03", "SQLi"],
      "severity": "CRITICAL",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

#### `POST /rules`

Create a new rule.

```bash
curl -X POST https://waf.example.com/api/v1/rules \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "rule_text": "SecRule ARGS \"@rx (?i)union.*select\" \"id:50001,phase:2,deny,status:403,msg:'\''SQLi detected'\'',severity:CRITICAL\"",
    "enabled": true,
    "description": "Custom SQL injection detection"
  }'
```

Response `201 Created`:
```json
{
  "id": "rule_xyz789",
  "rule_id": 50001,
  "enabled": true,
  "source": "custom",
  "created_at": "2025-03-14T12:00:00Z"
}
```

#### `PUT /rules/{id}`

Update a rule.

```bash
curl -X PUT https://waf.example.com/api/v1/rules/rule_xyz789 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"enabled": false}'
```

#### `DELETE /rules/{id}`

Delete a custom rule (built-in rules can only be disabled, not deleted).

#### `POST /rules/test`

Test a raw HTTP request against active rules without sending it to the backend.

```bash
curl -X POST https://waf.example.com/api/v1/rules/test \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "raw_request": "GET /search?q=1+UNION+SELECT+1,2,3-- HTTP/1.1\r\nHost: example.com\r\n\r\n"
  }'
```

Response:
```json
{
  "action": "block",
  "anomaly_score": 10,
  "triggered_rules": [
    {
      "rule_id": 942200,
      "msg": "Detects MySQL comments and SQL injection",
      "severity": "CRITICAL",
      "score": 10,
      "matched_var": "ARGS:q",
      "matched_value": "1 UNION SELECT 1,2,3--"
    }
  ]
}
```

---

### Events

#### `GET /events`

Query security events.

```bash
curl "https://waf.example.com/api/v1/events?\
  severity=CRITICAL&\
  action=block&\
  from=2025-03-01T00:00:00Z&\
  to=2025-03-14T23:59:59Z&\
  limit=100&\
  offset=0" \
  -H "Authorization: Bearer <token>"
```

Query parameters:

| Parameter | Type | Description |
|-----------|------|-------------|
| `severity` | string | Filter: CRITICAL, ERROR, WARNING, NOTICE |
| `action` | string | Filter: block, detect, allow |
| `rule_id` | integer | Filter by specific rule ID |
| `src_ip` | string | Filter by source IP or CIDR |
| `from` | ISO 8601 | Start of time range |
| `to` | ISO 8601 | End of time range |
| `limit` | integer | Results per page (max 1000, default 50) |
| `offset` | integer | Pagination offset |
| `export` | string | Set to `csv` to download as CSV |

#### `GET /events/{event_id}`

Get a single event by ID.

---

### Statistics

#### `GET /stats`

Get traffic statistics.

```bash
curl "https://waf.example.com/api/v1/stats?period=24h" \
  -H "Authorization: Bearer <token>"
```

Response:
```json
{
  "period": "24h",
  "requests_total": 1482930,
  "requests_blocked": 342,
  "block_rate_pct": 0.023,
  "latency_p50_ms": 0.8,
  "latency_p95_ms": 2.1,
  "latency_p99_ms": 4.7,
  "top_attack_types": [
    {"type": "SQLi", "count": 156},
    {"type": "XSS", "count": 89},
    {"type": "PathTraversal", "count": 43}
  ],
  "top_blocked_ips": [
    {"ip": "203.0.113.42", "count": 87}
  ]
}
```

---

### IP Lists

#### `GET /blocklist`

List all blocked IPs/CIDRs.

#### `POST /blocklist`

Add an IP or CIDR to the blocklist.

```bash
curl -X POST https://waf.example.com/api/v1/blocklist \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cidr": "203.0.113.0/24",
    "reason": "Attacking our login endpoint",
    "expires_at": "2025-04-01T00:00:00Z"
  }'
```

#### `DELETE /blocklist/{id}`

Remove an entry from the blocklist.

Same endpoints exist for `/allowlist`.

---

### Configuration

#### `GET /config`

Get current WAF configuration.

#### `PATCH /config`

Update configuration. Changes are applied immediately without restart.

```bash
curl -X PATCH https://waf.example.com/api/v1/config \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"detection": {"mode": "block", "paranoia_level": 2}}'
```

---

### Users (Admin only)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/users` | List all users |
| `POST` | `/users` | Create a user |
| `PUT` | `/users/{id}` | Update a user |
| `DELETE` | `/users/{id}` | Delete a user |
| `POST` | `/users/{id}/reset-password` | Force password reset |

---

## Error Responses

All errors follow this format:

```json
{
  "error": "short_error_code",
  "message": "Human-readable description of what went wrong",
  "details": {}
}
```

| HTTP Status | Error Code | Meaning |
|-------------|------------|---------|
| 400 | `invalid_request` | Request body or parameters are invalid |
| 401 | `unauthorized` | Missing or invalid authentication |
| 403 | `forbidden` | Authenticated but insufficient permissions |
| 404 | `not_found` | Resource does not exist |
| 409 | `conflict` | Resource already exists (e.g., duplicate rule ID) |
| 422 | `validation_error` | Input passes parsing but fails validation |
| 429 | `rate_limited` | Too many requests to the API |
| 500 | `internal_error` | Unexpected server error (report as bug) |

---

## Rate Limits

The Management API applies rate limiting to prevent abuse:

| Endpoint | Limit |
|----------|-------|
| `POST /auth/token` | 10 requests/minute per IP |
| `POST /rules/test` | 60 requests/minute per user |
| All other endpoints | 300 requests/minute per token |

Rate limit headers are included in every response:
```
X-RateLimit-Limit: 300
X-RateLimit-Remaining: 298
X-RateLimit-Reset: 1710432000
```

---

## Pagination

All list endpoints support pagination:

```bash
# Page 1 (first 50)
GET /events?limit=50&offset=0

# Page 2
GET /events?limit=50&offset=50
```

Response always includes `total` count for building pagination UI.

---

## Scopes (API Keys)

When creating API keys, specify the minimum scopes needed:

| Scope | Permits |
|-------|---------|
| `rules:read` | List and view rules |
| `rules:write` | Create, update, delete rules |
| `events:read` | View security events |
| `config:read` | View configuration |
| `config:write` | Update configuration |
| `ip_lists:write` | Manage blocklist/allowlist |
| `users:write` | Manage users (Admin only) |
| `*` | All permissions (use with caution) |
