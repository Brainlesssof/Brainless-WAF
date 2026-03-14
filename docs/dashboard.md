# Dashboard Guide

The Brainless WAF Dashboard is a browser-based interface for monitoring traffic, managing rules, and configuring your WAF in real time. It is available at `http://your-waf-host:8080` by default.

---

## Table of Contents

- [Logging In](#logging-in)
- [Overview Page](#overview-page)
- [Events Page](#events-page)
- [Rules Page](#rules-page)
- [Rule Tester](#rule-tester)
- [IP Management](#ip-management)
- [Settings](#settings)
- [User Management](#user-management)
- [Roles & Permissions](#roles--permissions)
- [Keyboard Shortcuts](#keyboard-shortcuts)

---

## Logging In

Navigate to `http://your-waf-host:8080`. Enter your username and password.

> **Default credentials:** `admin` / `changeme`  
> Change this immediately after first login via **Settings → Account → Change Password**.

Sessions expire after 1 hour of inactivity. You will be redirected to the login page and asked to re-authenticate.

---

## Overview Page

The Overview page is your real-time command center. It shows a live snapshot of what the WAF is seeing and doing right now.

### Top Metrics Bar

| Metric | Description |
|--------|-------------|
| **Requests/sec** | Current inbound request rate (5-second rolling average) |
| **Block Rate** | Percentage of requests blocked in the last 5 minutes |
| **p50 Latency** | Median latency added by the WAF (microseconds to low milliseconds) |
| **p99 Latency** | 99th percentile latency — tail latency indicator |
| **Active Rules** | Number of rules currently enabled |
| **WAF Mode** | Current mode: Learning / Detect / Block (click to change) |

### Live Traffic Chart

A stacked area chart showing requests per second over the last 30 minutes, broken down by:
- **Green** — allowed requests
- **Red** — blocked requests
- **Orange** — rate-limited requests

Hover over any point to see exact counts. Click and drag to zoom into a time range.

### Recent Threat Events

A live feed of the last 20 security events. Each entry shows:
- Timestamp
- Severity badge (CRITICAL / ERROR / WARNING / NOTICE)
- Triggering rule message
- Source IP
- Blocked URI

Click any event to open the full event detail panel.

### Top Attack Types

A donut chart of attack categories detected in the last hour (SQLi, XSS, Path Traversal, etc.). Click a segment to filter the Events page to that category.

### Geographic Threat Map

A world map showing source countries of blocked requests. Darker shading = more blocks from that country. Hover for exact counts.

---

## Events Page

The Events page shows the full security event log with filtering, search, and export.

### Filtering

Use the filter bar at the top to narrow results:

- **Severity** — CRITICAL, ERROR, WARNING, NOTICE (multi-select)
- **Action** — Block, Detect, Allow
- **Time Range** — Last hour / 24h / 7 days / 30 days / Custom
- **Source IP** — exact IP or CIDR
- **Rule ID** — specific rule number
- **Attack Tag** — OWASP-A03, SQLi, XSS, etc.

### Event Detail Panel

Click any event row to open the detail panel (slides in from the right):

```
Event ID:     evt_abc123xyz
Timestamp:    2025-03-14 12:34:56 UTC
Action:       BLOCK
Severity:     CRITICAL
Anomaly Score: 10/10

Rule:         942200 — Detects MySQL comments and SQL injection
Matched Variable: ARGS:q
Matched Value:    1' UNION SELECT 1,2,3--

Request:
  Method:  GET
  URI:     /search?q=1'+UNION+SELECT+1,2,3--
  Host:    example.com
  IP:      203.0.113.42
  Country: XX
  UA:      Mozilla/5.0 (compatible; sqlmap/1.7)

Tags: OWASP-A03, SQLi, UNION
```

**Actions in the detail panel:**
- **Block IP** — add this IP to the blocklist (with optional expiry)
- **Create Exception** — generate an exception rule for this specific request pattern
- **Copy as cURL** — copies the request as a cURL command for testing
- **View in Rule Editor** — opens the triggering rule in the Rules page

### Exporting Events

Click **Export** → **CSV** or **JSON** to download the current filtered result set. Large exports (>10,000 events) are queued and emailed when ready.

---

## Rules Page

The Rules page lets you view, create, edit, enable/disable, and delete WAF rules.

### Rule List

Rules are grouped by source:
- **CRS** — OWASP Core Rule Set (read-only, but can be disabled or have exceptions added)
- **Brainless Built-in** — built-in rules maintained by the core team
- **Custom** — rules you have created

Use the search bar to find rules by ID, message, or tag. Filter by enabled/disabled or phase.

### Creating a Rule

Click **+ New Rule**. The rule editor opens with syntax highlighting, autocompletion, and inline documentation.

```
SecRule ARGS "@rx (?i)(union.*select)" \
    "id:50001,phase:2,deny,status:403,\
     msg:'Custom SQLi Rule',\
     tag:OWASP-A03,severity:CRITICAL"
```

As you type, the editor validates syntax in real time and highlights:
- ✅ Green — valid syntax
- ⚠️ Yellow — valid but potentially performance-impacting (e.g., catastrophic backtracking risk in regex)
- ❌ Red — syntax error, cannot save

Click **Test Rule** to run the rule against the Rule Tester before saving.

Click **Save** to save and immediately push the rule to the WAF core (no restart needed).

### Disabling a Rule

Toggle the **Enabled** switch in the rule list row. The change takes effect within 1 second. Useful for quickly disabling a rule that is causing false positives without deleting it.

### Adding an Exception to a Built-in Rule

You cannot edit CRS rules directly, but you can add exceptions:

1. Find the rule in the list
2. Click **⋮ → Add Exception**
3. Choose the exception type:
   - **Exclude Parameter** — exclude a specific GET/POST parameter from this rule
   - **Exclude URL Path** — skip this rule for a specific path prefix
   - **Exclude IP** — skip this rule for a specific IP or CIDR
4. The exception is saved as a supplementary rule

---

## Rule Tester

The Rule Tester lets you test any HTTP request against the currently active ruleset without sending real traffic to your backend.

**Access:** Rules page → **Rule Tester** tab, or from any Event's detail panel.

### How to use

1. Paste a raw HTTP request into the input box:
```
GET /search?q=1+UNION+SELECT+1,2,3-- HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0
Accept: text/html
```

2. Click **Test**

3. The result shows:
   - **Final action** (Block / Allow / Detect)
   - **Anomaly score** breakdown
   - **Every rule that matched**, with the matched variable and value
   - **Rules that almost matched** (score > 0 but not triggered)

This is the fastest way to diagnose false positives or verify that a new attack pattern is being caught.

---

## IP Management

Manage your IP blocklist and allowlist.

### Blocklist

IPs in the blocklist are rejected at the first pipeline stage before any rule evaluation. Add IPs that are:
- Known attack sources
- Scrapers causing load
- Temporary blocks during an incident

Each entry supports:
- Single IP: `203.0.113.42`
- CIDR range: `203.0.113.0/24`
- Optional reason (shown in event logs)
- Optional expiry time (auto-removes after that time)

### Allowlist

IPs in the allowlist bypass all rule evaluation. Use for:
- Internal monitoring tools
- Known good scanners (e.g., your own Nessus scanner)
- Partner services with unusual but legitimate traffic patterns
- Your own office/VPN IP ranges

> **Caution:** Allowlisted IPs bypass ALL rules including rate limiting. Only allowlist IPs you fully trust and control.

---

## Settings

### WAF Mode

Switch between **Learning**, **Detect**, and **Block** modes. The change is applied live.

A confirmation dialog is shown before switching to **Block** mode to prevent accidental lockouts.

### Paranoia Level

Slider from 1–4. Higher = more rules active = more false positives. See [Configuration Reference](configuration.md#detectionparanoia_level) for details.

### Anomaly Score Threshold

The score at which a request is blocked. Default: 10. Raise this if you're seeing excessive false positives; lower it for stricter protection.

### Rate Limiting

Configure default rate limits and per-endpoint overrides.

### Account

Change your password and manage your API keys.

---

## User Management

*Available to Admin role only.*

### Creating a User

**Settings → Users → + New User**

| Field | Description |
|-------|-------------|
| Username | Login name (lowercase, no spaces) |
| Email | For password reset notifications |
| Role | Admin, Analyst, or Read-only |
| Password | Temporary password (user must change on first login) |

### Roles

| Role | What they can do |
|------|----------------|
| **Admin** | Full access — all pages, create/delete users, change WAF mode |
| **Analyst** | View all pages, create/edit rules, manage IP lists, cannot manage users or change WAF mode |
| **Read-only** | View overview, events, and rules — no changes |

---

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `?` | Show keyboard shortcuts overlay |
| `g o` | Go to Overview |
| `g e` | Go to Events |
| `g r` | Go to Rules |
| `g i` | Go to IP Management |
| `g s` | Go to Settings |
| `n` | New rule (on Rules page) |
| `/` | Focus search bar |
| `Esc` | Close panel / modal |
| `Ctrl+S` | Save (in rule editor) |
| `Ctrl+Enter` | Test rule (in rule editor) |
