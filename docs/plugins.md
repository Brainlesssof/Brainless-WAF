# Plugin System

> **Note:** The plugin system is planned for v1.1. This document describes the intended design and API. Implementation details may change before release.

---

## Overview

The Brainless WAF plugin system allows the community to extend the WAF with new capabilities without modifying the core. Plugins can add:

- New detection operators (e.g., `@geoip`, `@reputation`)
- New rule variables (e.g., `GEO:COUNTRY_CODE`, `THREAT:SCORE`)
- New actions (e.g., `challenge`, `captcha`, `tarpit`)
- New data sources (e.g., custom threat feeds, internal allow/blocklists)
- Integration hooks (e.g., webhooks on block events, Slack notifications)

---

## Plugin Types

### Operator Plugins

Add new `@operator` keywords to the rule engine:

```go
// Example: GeoIP operator plugin
type GeoIPOperator struct {
    db *geoip2.Reader
}

func (op *GeoIPOperator) Name() string { return "geoip" }

func (op *GeoIPOperator) Evaluate(value string, param string) (bool, error) {
    ip := net.ParseIP(value)
    record, err := op.db.Country(ip)
    if err != nil {
        return false, err
    }
    return strings.Contains(param, record.Country.IsoCode), nil
}
```

Register it:
```go
brainless.RegisterOperator(&GeoIPOperator{db: db})
```

Now available in rules:
```
SecRule REMOTE_ADDR "@geoip CN RU KP" "id:90001,phase:1,deny"
```

### Variable Plugins

Add new variables that rules can inspect:

```go
type ThreatScoreVariable struct {
    feed *ThreatFeed
}

func (v *ThreatScoreVariable) Name() string { return "THREAT:SCORE" }

func (v *ThreatScoreVariable) Extract(tx *brainless.Transaction) (string, error) {
    score := v.feed.GetScore(tx.RemoteAddr())
    return strconv.Itoa(score), nil
}
```

Now available in rules:
```
SecRule THREAT:SCORE "@gt 80" "id:90010,phase:1,deny,msg:'High threat score IP'"
```

### Action Plugins

Add new actions:

```go
type SlackNotifyAction struct {
    webhookURL string
}

func (a *SlackNotifyAction) Name() string { return "slack_notify" }

func (a *SlackNotifyAction) Execute(tx *brainless.Transaction, params string) error {
    // Send Slack notification
    return sendSlackAlert(a.webhookURL, tx.Event())
}
```

Now available in rules:
```
SecRule TX:ANOMALY_SCORE "@gt 15" \
    "id:90020,phase:2,deny,slack_notify:'High severity attack blocked'"
```

---

## Installing Plugins

### From the Plugin Registry

```bash
# List available plugins
brainless-ctl plugin list

# Install a plugin
brainless-ctl plugin install brainless-plugin-geoip

# List installed plugins
brainless-ctl plugin list --installed

# Update a plugin
brainless-ctl plugin update brainless-plugin-geoip

# Remove a plugin
brainless-ctl plugin remove brainless-plugin-geoip
```

### Manual Installation

```bash
# Copy plugin binary to the plugins directory
cp my-plugin.so /etc/brainless/plugins/

# Restart WAF to load
brainless-ctl reload
```

---

## Official Plugins

These plugins are maintained by the Brainless Security team:

| Plugin | Description | Install |
|--------|-------------|---------|
| `brainless-plugin-geoip` | Country/ASN-based blocking using MaxMind GeoLite2 | `brainless-ctl plugin install brainless-plugin-geoip` |
| `brainless-plugin-recaptcha` | Google reCAPTCHA v3 challenge for suspicious requests | `brainless-ctl plugin install brainless-plugin-recaptcha` |
| `brainless-plugin-cloudflare` | Forward real IPs from Cloudflare's proxy | `brainless-ctl plugin install brainless-plugin-cloudflare` |
| `brainless-plugin-crowdsec` | Integrate CrowdSec threat intelligence feed | `brainless-ctl plugin install brainless-plugin-crowdsec` |
| `brainless-plugin-slack` | Send block event notifications to Slack | `brainless-ctl plugin install brainless-plugin-slack` |
| `brainless-plugin-pagerduty` | Trigger PagerDuty incidents on critical events | `brainless-ctl plugin install brainless-plugin-pagerduty` |

---

## Writing a Plugin

### Project Structure

```
my-plugin/
├── main.go           # Plugin entry point
├── go.mod
├── plugin.yaml       # Plugin metadata
└── README.md
```

### plugin.yaml

```yaml
name: my-awesome-plugin
version: 1.0.0
description: Does something awesome
author: Your Name <you@example.com>
license: Apache-2.0
homepage: https://github.com/you/my-awesome-plugin
brainless_min_version: "1.1.0"

# Resources this plugin needs
permissions:
  - network_outbound    # Make outbound HTTP requests
  - file_read           # Read files from allowed paths

config_schema:
  api_key:
    type: string
    required: true
    description: API key for the external service
  timeout:
    type: duration
    default: 5s
    description: Request timeout
```

### main.go

```go
package main

import (
    "github.com/brainless-security/brainless-waf/pkg/plugin"
)

// Plugin entry point — must be named "Plugin"
var Plugin plugin.Plugin = &MyPlugin{}

type MyPlugin struct {
    config MyConfig
}

type MyConfig struct {
    APIKey  string        `yaml:"api_key"`
    Timeout time.Duration `yaml:"timeout"`
}

func (p *MyPlugin) Init(ctx plugin.Context) error {
    // Load config
    if err := ctx.Config(&p.config); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    // Register operators/variables/actions
    ctx.RegisterOperator(&MyOperator{config: p.config})
    return nil
}

func (p *MyPlugin) Shutdown() error {
    return nil
}
```

### Building

```bash
go build -buildmode=plugin -o my-plugin.so .
```

### Submitting to the Registry

1. Ensure your plugin has tests and documentation
2. Publish to GitHub with the `brainless-waf-plugin` topic
3. Open a PR to add it to the [community plugin registry](https://github.com/brainless-security/plugin-registry)

---

## Plugin Security

Plugins run in the WAF process with the same privileges. Only install plugins from trusted sources. The official registry reviews all plugins before listing them.

Plugin permissions are declared in `plugin.yaml` and enforced at runtime:
- `network_outbound` — required to make HTTP calls to external services
- `file_read` — required to read database files (e.g., GeoIP `.mmdb` files)
- `file_write` — required to write to disk (rarely needed)

Plugins that request `file_write` or unrestricted `network_outbound` receive additional review scrutiny.
