# Rule Engine Guide

This guide covers the Brainless Rule Format (BRF), all supported operators, variables, and actions, and how to write effective custom rules.

---

## Table of Contents

- [Overview](#overview)
- [Rule Syntax](#rule-syntax)
- [Phases](#phases)
- [Variables](#variables)
- [Operators](#operators)
- [Actions](#actions)
- [Anomaly Scoring](#anomaly-scoring)
- [Stateful Rules (Counters)](#stateful-rules-counters)
- [Rule Chaining](#rule-chaining)
- [Lua Scripting](#lua-scripting)
- [OWASP CRS Compatibility](#owasp-crs-compatibility)
- [Rule ID Ranges](#rule-id-ranges)
- [Testing Rules](#testing-rules)
- [Examples](#examples)

---

## Overview

Rules are the primary mechanism for defining what Brainless WAF blocks, allows, and logs. Rules are evaluated in order during specific pipeline phases. When a rule matches, it executes one or more actions (block, log, score, etc.).

Rules are loaded from `.rules` files in the configured rules directory. Files are loaded in alphabetical order; within a file, rules are evaluated top to bottom.

---

## Rule Syntax

```
SecRule VARIABLE "OPERATOR" "ACTION_LIST"
```

**SecRule** — the directive. Always `SecRule`.

**VARIABLE** — what to inspect (request URI, headers, body parameters, etc.)

**OPERATOR** — how to inspect it (regex match, string comparison, numeric comparison, etc.)

**ACTION_LIST** — comma-separated list of actions to take when the rule matches.

### Example

```
SecRule ARGS "@rx (?i)(union.*select|insert.*into)" \
    "id:10001,\
     phase:2,\
     deny,\
     status:403,\
     msg:'SQL Injection Detected',\
     tag:OWASP-A03,\
     severity:CRITICAL,\
     logdata:'Matched value: %{MATCHED_VAR}'"
```

---

## Phases

Phases control when a rule is evaluated in the request/response lifecycle.

| Phase | Name | When it runs | Typical use |
|-------|------|-------------|-------------|
| `phase:1` | Request Headers | After headers received, before body | IP checks, header validation |
| `phase:2` | Request Body | After full request body received | Parameter inspection, body scanning |
| `phase:3` | Response Headers | After upstream sends response headers | Response header enforcement |
| `phase:4` | Response Body | After full response body received | PII detection, error suppression |

> **Performance note:** Phase 4 (response body scanning) buffers the full response in memory. Enable only for endpoints where it's needed.

---

## Variables

### Request Variables

| Variable | Description |
|----------|-------------|
| `REQUEST_URI` | Raw request URI (path + query string) |
| `REQUEST_URI_RAW` | URI before URL decoding |
| `REQUEST_FILENAME` | Path component only (no query string) |
| `QUERY_STRING` | Raw query string |
| `REQUEST_METHOD` | HTTP method (GET, POST, etc.) |
| `REQUEST_PROTOCOL` | HTTP version (HTTP/1.1, HTTP/2) |
| `REQUEST_HEADERS` | All request headers |
| `REQUEST_HEADERS:Name` | Specific header (e.g., `REQUEST_HEADERS:Content-Type`) |
| `REQUEST_COOKIES` | All cookies |
| `REQUEST_COOKIES:name` | Specific cookie |
| `REQUEST_BODY` | Raw request body |
| `ARGS` | All GET and POST parameters |
| `ARGS:name` | Specific parameter (e.g., `ARGS:username`) |
| `ARGS_NAMES` | Just the parameter names |
| `FILES` | Uploaded file names |
| `FILES_SIZES` | Uploaded file sizes |
| `REMOTE_ADDR` | Client IP address |
| `REMOTE_HOST` | Client hostname (if reverse DNS available) |
| `SERVER_NAME` | Virtual host name |

### Response Variables

| Variable | Description |
|----------|-------------|
| `RESPONSE_STATUS` | HTTP response status code |
| `RESPONSE_HEADERS` | All response headers |
| `RESPONSE_HEADERS:Name` | Specific response header |
| `RESPONSE_BODY` | Full response body |

### Transaction Variables

| Variable | Description |
|----------|-------------|
| `TX:name` | Transaction variable (set by rules, scoped to this request) |
| `TX:ANOMALY_SCORE` | Current accumulated anomaly score |
| `IP:name` | Per-IP persistent variable |
| `SESSION:name` | Per-session variable (requires session tracking) |

### Variable Modifiers

Apply to any collection variable (ARGS, REQUEST_HEADERS, etc.):

```
ARGS                        # All parameters
ARGS:username               # Only "username" parameter
ARGS|ARGS_NAMES             # Union of two collections
!ARGS:safe_param            # Exclude "safe_param" from ARGS
```

---

## Operators

### String Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `@rx` | Regular expression match | `@rx (?i)select.*from` |
| `@streq` | Exact string equality | `@streq admin` |
| `@contains` | Substring match | `@contains <script>` |
| `@containsWord` | Word boundary match | `@containsWord union` |
| `@beginsWith` | Starts with string | `@beginsWith /admin` |
| `@endsWith` | Ends with string | `@endsWith .php` |
| `@within` | Value is in whitelist | `@within "GET POST"` |
| `@pm` | Multi-pattern match (fast) | `@pmFromFile sql_keywords.txt` |
| `@pmFromFile` | Load patterns from file | `@pmFromFile /etc/brainless/rules/patterns/sqli.txt` |

### Numeric Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `@eq` | Equal | `@eq 200` |
| `@ne` | Not equal | `@ne 0` |
| `@gt` | Greater than | `@gt 10` |
| `@lt` | Less than | `@lt 100` |
| `@ge` | Greater than or equal | `@ge 5` |
| `@le` | Less than or equal | `@le 1000` |

### IP / Network Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `@ipMatch` | Match IP or CIDR | `@ipMatch 10.0.0.0/8` |
| `@ipMatchFromFile` | Load IPs from file | `@ipMatchFromFile /etc/brainless/rules/tor_exits.txt` |

### Validation Operators

| Operator | Description |
|----------|-------------|
| `@validateByteRange 1-255` | Check each byte is in range |
| `@validateUrlEncoding` | Check URL encoding is valid |
| `@validateUtf8Encoding` | Check UTF-8 is valid |
| `@detectSQLi` | LibInjection SQL injection detector |
| `@detectXSS` | LibInjection XSS detector |

### Negation

Prefix any operator with `!` to negate:

```
SecRule REMOTE_ADDR "!@ipMatch 192.168.0.0/16" "id:...,deny"
```

---

## Actions

### Disruptive Actions (choose one)

| Action | Description |
|--------|-------------|
| `deny` | Block request, return error response |
| `allow` | Allow request, skip remaining rules |
| `pass` | Allow this rule to match but continue evaluating |
| `redirect:url` | Redirect client to URL |
| `proxy:url` | Forward to different upstream |

### Non-Disruptive Actions

| Action | Description |
|--------|-------------|
| `log` | Log this rule match (default: on) |
| `nolog` | Do not log this match |
| `auditlog` | Include in audit log |
| `noauditlog` | Exclude from audit log |
| `msg:'text'` | Human-readable message for the log |
| `logdata:'%{MATCHED_VAR}'` | Log the matched value |
| `tag:name` | Categorize the event (can use multiple) |
| `severity:level` | CRITICAL, ERROR, WARNING, NOTICE |
| `status:code` | HTTP status code for deny response |
| `setvar:name=value` | Set a transaction or IP variable |
| `expirevar:name=seconds` | Set TTL on a variable |
| `skipAfter:id` | Skip rules until rule `id` |
| `chain` | Chain with next rule (both must match) |

### Required Actions

Every rule **must** include:
- `id:NNNNN` — unique numeric ID
- `phase:N` — evaluation phase
- A disruptive action OR `pass`

---

## Anomaly Scoring

Instead of blocking on every rule match, Brainless WAF accumulates an anomaly score and blocks when the total exceeds a threshold. This reduces false positives.

```
SecRule ARGS "@detectSQLi" \
    "id:10100,phase:2,pass,\
     setvar:tx.anomaly_score=+%{tx.critical_anomaly_score},\
     tag:OWASP-A03,severity:CRITICAL"
```

Severity score values (configurable):

| Severity | Default score |
|----------|--------------|
| CRITICAL | 10 |
| ERROR | 5 |
| WARNING | 3 |
| NOTICE | 1 |

The final blocking rule (always last in the ruleset):

```
SecRule TX:ANOMALY_SCORE "@ge %{tx.inbound_anomaly_score_threshold}" \
    "id:19999,phase:2,deny,status:403,\
     msg:'Inbound anomaly score exceeded',\
     logdata:'Score: %{TX:ANOMALY_SCORE}'"
```

---

## Stateful Rules (Counters)

Use `IP:` variables for per-IP state that persists across requests:

```
# Track login attempts per IP
SecRule REQUEST_URI "@streq /api/auth/login" \
    "id:20001,phase:1,pass,\
     setvar:ip.login_attempts=+1,\
     expirevar:ip.login_attempts=300"

# Block after 10 attempts in 5 minutes
SecRule IP:LOGIN_ATTEMPTS "@gt 10" \
    "id:20002,phase:1,deny,status:429,\
     msg:'Too many login attempts',\
     setvar:ip.blocked=1,\
     expirevar:ip.blocked=900"
```

---

## Rule Chaining

Chain rules with `chain` — all rules in a chain must match for the action to trigger:

```
# Block only POST requests to /admin from non-internal IPs
SecRule REQUEST_METHOD "@streq POST" \
    "id:30001,phase:1,chain,\
     msg:'Unauthorized admin POST'"
    SecRule REQUEST_URI "@beginsWith /admin" \
        "chain"
        SecRule REMOTE_ADDR "!@ipMatch 10.0.0.0/8"  \
            "deny,status:403"
```

---

## Lua Scripting

For logic that's too complex for BRF syntax, embed Lua scripts:

```
SecRule REQUEST_HEADERS:Authorization "@rx ^Bearer\s+(.+)$" \
    "id:40001,phase:1,pass,capture,\
     setvar:tx.jwt_token=%{TX:1}"

SecRuleScript /etc/brainless/rules/scripts/validate_jwt.lua \
    "id:40002,phase:1,deny,status:401,\
     msg:'Invalid JWT token'"
```

```lua
-- /etc/brainless/rules/scripts/validate_jwt.lua
local jwt = require "brainless.jwt"
local token = m.getvar("TX:JWT_TOKEN")

if not token then
    return nil  -- no token, let other rules handle
end

local payload, err = jwt.verify(token, m.getConfig("jwt_secret"))
if err then
    m.log(4, "JWT validation failed: " .. err)
    return true  -- triggers the deny action
end

m.setvar("TX:USER_ID", payload.sub)
m.setvar("TX:USER_ROLE", payload.role)
return nil  -- JWT valid, continue
```

---

## OWASP CRS Compatibility

Brainless WAF is compatible with OWASP Core Rule Set 4.x. Import existing CRS rules:

```bash
# Import from a local CRS directory
brainless-ctl import-crs /path/to/coreruleset/rules/

# Or enable bundled CRS in config
rules:
  crs_enabled: true
  crs_paranoia_level: 2
```

Most ModSecurity directives are supported. Exceptions and workarounds are documented in `docs/crs-compatibility.md`.

---

## Rule ID Ranges

| Range | Owner | Notes |
|-------|-------|-------|
| 1–99,999 | OWASP CRS | Do not create rules in this range |
| 100,000–199,999 | Brainless built-in | Maintained by core team |
| 200,000–299,999 | Reserved | For future use |
| 500,000–599,999 | Community rules | Submit via PR |
| 900,000–999,999 | Local/custom rules | Your site-specific rules |

---

## Testing Rules

### Unit test a rule file

```bash
brainless-ctl rule-test \
  --rules rules/custom/my_rule.rules \
  --fixtures tests/fixtures/

# Expected output:
# ✓ sqli_basic.http → BLOCKED (rule 10001)
# ✓ xss_reflected.http → BLOCKED (rule 10050)
# ✓ normal_search.http → ALLOWED
# 3 tests passed, 0 failed
```

### Test with a raw HTTP request file

```bash
# tests/fixtures/attacks/sqli/union_select.http
GET /products?id=1+UNION+SELECT+1,2,3-- HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0

# Run
brainless-ctl rule-test \
  --rules rules/brainless/sqli.rules \
  --request tests/fixtures/attacks/sqli/union_select.http
```

### YAML test format

```yaml
# tests/rules/my_rule_test.yaml
rule_file: rules/custom/my_rule.rules

should_block:
  - name: "Basic UNION SELECT"
    request:
      method: GET
      uri: "/search?q=1 UNION SELECT 1,2,3--"
  - name: "Stacked queries"
    request:
      method: POST
      uri: "/login"
      body: "user=admin'; DROP TABLE users--&pass=x"
      headers:
        Content-Type: application/x-www-form-urlencoded

should_pass:
  - name: "Normal search"
    request:
      method: GET
      uri: "/search?q=blue+jeans"
  - name: "Legitimate SELECT word in product name"
    request:
      method: GET
      uri: "/products?name=select+combo+box"
```

---

## Examples

### Block requests from a country (GeoIP plugin required)

```
SecRule GEO:COUNTRY_CODE "@within CN RU KP" \
    "id:900001,phase:1,deny,status:403,\
     msg:'Request from blocked country',\
     tag:GeoIP"
```

### Enforce Content-Type on API endpoints

```
SecRule REQUEST_URI "@beginsWith /api/" \
    "id:900010,phase:1,chain"
    SecRule REQUEST_METHOD "@within POST PUT PATCH" \
        "chain"
        SecRule REQUEST_HEADERS:Content-Type "!@rx ^application/json" \
            "deny,status:415,\
             msg:'API requires Content-Type: application/json'"
```

### Block large file uploads

```
SecRule FILES_SIZES "@gt 10485760" \
    "id:900020,phase:2,deny,status:413,\
     msg:'File upload too large (max 10MB)'"
```

### Add security headers to all responses

```
SecRule RESPONSE_STATUS "@ge 0" \
    "id:900030,phase:3,pass,\
     setResponseHeader:'X-Frame-Options: DENY',\
     setResponseHeader:'X-Content-Type-Options: nosniff',\
     setResponseHeader:'Strict-Transport-Security: max-age=31536000'"
```

### Allowlist a monitoring endpoint

```
SecRule REQUEST_URI "@beginsWith /health" \
    "id:900040,phase:1,allow,\
     nolog,\
     msg:'Health check bypass'"
```
