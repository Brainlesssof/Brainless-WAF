# Migration Guide

## Migrating from ModSecurity to Brainless WAF

This guide walks you through migrating an existing ModSecurity deployment to Brainless WAF. The process typically takes 1–4 hours depending on the complexity of your custom rules.

---

## Overview

Brainless WAF's rule engine (BRF) is syntactically compatible with ModSecurity's SecLang / CRS 4.x. In most cases, your existing rules import without modification. The main effort is:

1. Importing your existing rules
2. Migrating your `modsecurity.conf` settings to Brainless WAF's `config.yaml`
3. Verifying behavior in detect mode before switching to block

---

## Step 1 — Audit Your Current Setup

Before migrating, document what you have:

```bash
# Find all ModSecurity rule files
find /etc/modsecurity /etc/nginx/modsec -name "*.conf" -o -name "*.rules" 2>/dev/null

# Count custom rules
grep -r "^SecRule" /path/to/custom/rules | wc -l

# Find all SecRuleUpdateTargetById (exceptions)
grep -r "SecRuleUpdateTargetById" /etc/modsecurity | sort -u

# Check which CRS version you're on
head -5 /usr/share/modsecurity-crs/rules/REQUEST-900-EXCLUSION-RULES-BEFORE-CRS.conf
```

---

## Step 2 — Install Brainless WAF in Parallel

Do not replace ModSecurity yet. Run Brainless WAF in parallel on a staging server or different port:

```bash
docker run -d \
  -p 9090:80 \
  -v ./config:/etc/brainless \
  brainlesssecurity/brainless-waf:1.0.0
```

---

## Step 3 — Import Your CRS Rules

```bash
# If you're using OWASP CRS, enable the built-in version (already bundled)
brainless-ctl config set rules.crs_enabled true

# Or import your specific CRS version from disk
brainless-ctl import-crs /usr/share/modsecurity-crs/rules/

# Output shows import results
# Imported: 4,521 rules
# Skipped:  12 rules (unsupported directives — see below)
# Errors:   0
```

---

## Step 4 — Import Your Custom Rules

```bash
# Import custom rule files
brainless-ctl import-rules /etc/modsecurity/custom/ \
  --output /etc/brainless/rules/custom/

# The tool will:
# 1. Parse each rule file
# 2. Report incompatible directives
# 3. Write compatible rules to the output directory
# 4. Write a migration report to migration-report.txt
```

Review `migration-report.txt` for any rules that need manual adjustment.

---

## Step 5 — Migrate Configuration

Map your `modsecurity.conf` settings to Brainless WAF's `config.yaml`:

| ModSecurity | Brainless WAF | Notes |
|-------------|--------------|-------|
| `SecRuleEngine On` | `detection.mode: block` | |
| `SecRuleEngine DetectionOnly` | `detection.mode: detect` | |
| `SecRuleEngine Off` | `detection.mode: learning` | |
| `SecRequestBodyAccess On` | `detection.mode: block` (default) | Phase 2 always active |
| `SecResponseBodyAccess On` | `detection.response_body_inspection: true` | Off by default |
| `SecResponseBodyLimit 1048576` | `detection.response_body_limit: 1048576` | |
| `SecRequestBodyLimit 13107200` | `advanced.request_body_limit: 13107200` | |
| `SecPcreMatchLimit 100000` | Not needed — Go regex doesn't have this limit | |
| `SecAuditLog /var/log/modsec_audit.log` | `logging.audit_log: /var/log/brainless/audit.log` | |
| `SecDefaultAction "phase:1,log,auditlog,pass"` | Default in Brainless WAF | |

**Paranoia Level mapping:**

If you're using CRS paranoia level variable:
```
# modsecurity CRS setup
setvar:'tx.paranoia_level=2'

# Brainless WAF equivalent
detection:
  paranoia_level: 2
```

---

## Step 6 — Migrate SecRuleUpdateTargetById Exceptions

Your existing exceptions import automatically via Step 4. Verify they appear correctly:

```bash
# List all exception rules
brainless-ctl rules list --type exception
```

If you have exceptions in a separate file, import them:
```bash
brainless-ctl import-rules /etc/modsecurity/crs-exclusion-rules.conf
```

---

## Step 7 — Verify in Detect Mode

Start Brainless WAF in detect mode and route a copy of your production traffic to it (or use a staging environment):

```yaml
detection:
  mode: detect
```

Run for 48 hours and compare:

```bash
# Compare block events between ModSecurity and Brainless WAF
brainless-ctl stats --period 48h --format summary

# Export detect-mode events for review
brainless-ctl events export --action detect --format csv > events.csv
```

Look for:
- False positives Brainless WAF catches that ModSecurity didn't (and vice versa)
- Rules that behave differently
- Any legitimate traffic that would be blocked

---

## Step 8 — Switch to Block Mode

Once satisfied:

```bash
brainless-ctl config set detection.mode block

# Update DNS / load balancer to point to Brainless WAF instead of ModSecurity
# Monitor for 1 hour
```

---

## Unsupported ModSecurity Directives

These ModSecurity directives are not supported in Brainless WAF v1.0:

| Directive | Status | Alternative |
|-----------|--------|-------------|
| `SecGeoLookupDb` | Not supported | Use GeoIP plugin (v1.1+) |
| `SecHttpBlKey` | Not supported | Use IP blocklist with threat feeds |
| `SecRemoteRules` | Not supported | Use `rules.auto_update` for remote rule updates |
| `SecUploadDir` | Not supported | File uploads are inspected in memory |
| `SecStreamInBodyInspection` | Not supported | Full body buffering only |
| `SecXmlExternalEntity` | Not applicable | XML external entity injection blocked by default |
| `SecUnicodeMapFile` | Not supported | Unicode normalization is automatic |

Rules using unsupported directives are flagged in the migration report with a suggested equivalent.

---

## Migrating from Nginx + ModSecurity

If you're also migrating away from NGINX as your proxy:

```bash
# Map nginx proxy settings to Brainless WAF config
# proxy_pass http://backend → server.upstream: http://backend
# proxy_set_header X-Forwarded-For → configured automatically
# proxy_read_timeout 30s → server.upstream_timeout: 30s

# TLS: copy your nginx TLS config
server:
  tls:
    cert: /etc/ssl/certs/your-cert.pem   # same cert as nginx
    key:  /etc/ssl/private/your-key.pem
```

After migration, NGINX is no longer needed — Brainless WAF handles both proxying and TLS.

---

## Rollback Plan

If issues arise after cutover:

```bash
# Option 1: Switch back to ModSecurity (DNS change or load balancer change)
# Option 2: Switch Brainless WAF to detect mode (no blocking) while investigating
brainless-ctl config set detection.mode detect
```

Keep ModSecurity running in parallel for at least 30 days before decommissioning.
