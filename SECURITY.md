# Security Policy

## Supported Versions

We actively maintain security patches for the following versions:

| Version | Supported | End of Support |
|---------|-----------|----------------|
| 1.x (latest) | ✅ Active | TBD |
| 0.9.x | ✅ Security fixes only | 6 months after v1.1 release |
| < 0.9 | ❌ Not supported | — |

We strongly recommend always running the latest stable release.

---

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues, discussions, or Discord.**

### How to Report

Send your report to: **`security@brainless-security.io`**

Encrypt your report using our PGP key if the vulnerability is sensitive:

```
Key ID: 0xABCD1234EFGH5678
Fingerprint: XXXX XXXX XXXX XXXX XXXX  XXXX XXXX XXXX XXXX XXXX
Download: https://brainless-security.io/pgp-key.asc
```

### What to Include

A good vulnerability report helps us respond faster. Please include:

- **Description** — what is the vulnerability and what does it affect?
- **Impact** — what can an attacker do by exploiting it?
- **Affected versions** — which versions are vulnerable?
- **Steps to reproduce** — a minimal, reliable reproduction case
- **Proof of concept** — code or commands that demonstrate the issue (optional but helpful)
- **Suggested fix** — if you have one (optional)

### What Happens Next

| Timeline | Action |
|----------|--------|
| Within 48 hours | We acknowledge receipt of your report |
| Within 7 days | We confirm whether the issue is valid and our initial assessment |
| Within 14 days | We develop and test a fix (complex issues may take longer) |
| Within 30 days | We release a patch and publish a security advisory |
| After patch release | You may publish details (coordinated disclosure) |

We will keep you informed throughout the process. If you do not receive an acknowledgement within 48 hours, please follow up via Discord (DM a Core Maintainer directly).

---

## Coordinated Disclosure Policy

We follow responsible / coordinated disclosure:

1. Reporter notifies us privately
2. We work together to understand and fix the issue
3. We release a patch
4. We publish a [GitHub Security Advisory](https://github.com/brainless-security/brainless-waf/security/advisories)
5. Reporter may publish their own write-up 30 days after the patch is released (or earlier with our agreement)

We will credit you in the security advisory unless you prefer to remain anonymous. We do not offer a bug bounty program at this time, but we deeply appreciate responsible disclosure.

---

## Scope

### In Scope

- Vulnerabilities in the Brainless WAF core engine (`/core`)
- Vulnerabilities in the Management API (`/management`)
- Vulnerabilities in the Dashboard (`/dashboard`)
- WAF bypass techniques (requests that should be blocked but are not)
- Authentication or authorization flaws in the management API
- Remote code execution, privilege escalation, or data exposure
- Supply chain vulnerabilities (malicious dependencies)
- Docker image vulnerabilities introduced by our configuration

### Out of Scope

- Vulnerabilities in third-party dependencies (report to the upstream project; notify us too if Brainless WAF is specifically affected)
- Vulnerabilities in your own backend that Brainless WAF is protecting
- Social engineering attacks
- Denial of service via resource exhaustion on intentionally unprotected endpoints
- Issues in demo/test environments with intentionally weak configurations
- Missing security headers on endpoints that are documented as not requiring them

---

## Security Best Practices for Deployers

When deploying Brainless WAF in production:

- **Change default credentials** immediately (`admin` / `changeme` must not be used in production)
- **Restrict API access** — the Management API should not be publicly accessible; place it behind a VPN or private network
- **Enable TLS** on all interfaces including the dashboard and API
- **Start in `learning` mode** for 72 hours before switching to `block`
- **Keep rules updated** — enable `auto_update: true` to receive the latest threat signatures
- **Monitor logs** — forward logs to a SIEM and set up alerts for high-severity events
- **Pin the Docker image** to a specific version tag, not `:latest`, in production
- **Rotate API keys** regularly and revoke any that are no longer needed
- **Review allowlists** periodically — overly permissive allowlists defeat the purpose of the WAF

---

## Security Architecture

### Core Security Principles

- **Zero trust** — every request is treated as potentially malicious until proven safe
- **Defense in depth** — multiple detection layers (signatures, anomaly scoring, ML)
- **Fail safe** — in the event of an internal error, Brainless WAF returns a 502 rather than allowing the request through
- **Minimal attack surface** — the Management API and Dashboard are separate processes and can be disabled entirely in high-security deployments

### Internal Security Controls

- All inter-component communication uses mTLS
- Database credentials are never stored in config files (use environment variables or Vault)
- The WAF core runs as a non-root user (`uid: 1000`) in Docker
- Docker image is built from `scratch` (no shell, no package manager) in production builds
- All build artifacts are signed with cosign

---

## Known Limitations

- Encrypted traffic (end-to-end encrypted payloads within HTTPS) cannot be inspected beyond the outer TLS layer. This is by design.
- WebSocket inspection is available from v1.1+ only; earlier versions pass WebSocket traffic without deep inspection.
- HTTP/3 QUIC traffic inspection is in beta; behavior may differ from HTTP/1.1 and HTTP/2.

---

## Hall of Fame

We thank the following researchers for responsibly disclosing security issues:

| Researcher | Issue | Version Fixed |
|------------|-------|---------------|
| *(Be the first!)* | — | — |

---

*This security policy was last updated: 2025.*
