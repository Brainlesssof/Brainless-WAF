# Production Deployment Guide

This guide covers everything you need to run Brainless WAF in a production environment safely and reliably.

---

## Table of Contents

- [Pre-Deployment Checklist](#pre-deployment-checklist)
- [Deployment Modes](#deployment-modes)
- [Docker Compose — Production](#docker-compose--production)
- [Kubernetes — Production](#kubernetes--production)
- [High Availability Setup](#high-availability-setup)
- [Network Architecture](#network-architecture)
- [TLS in Production](#tls-in-production)
- [Secrets Management](#secrets-management)
- [Monitoring Stack](#monitoring-stack)
- [Backup & Recovery](#backup--recovery)
- [Incident Response](#incident-response)
- [Maintenance Procedures](#maintenance-procedures)

---

## Pre-Deployment Checklist

Complete this checklist before going live:

### Security
- [ ] Default `admin` password changed
- [ ] Management API (`port 8000`) is NOT accessible from the public internet
- [ ] Dashboard (`port 8080`) is behind authentication (or on internal network only)
- [ ] TLS enabled on all public-facing ports
- [ ] `JWT_SECRET` and `API_KEY_ENCRYPTION_KEY` are strong random values (≥32 chars), stored in secrets manager
- [ ] `fail_open: false` in configuration (default)
- [ ] Firewall rules block direct access to backend (all traffic must go through WAF)

### Configuration
- [ ] `upstream` points to the correct backend server(s)
- [ ] `trusted_proxies` set correctly (if behind a load balancer)
- [ ] `real_ip_header` set correctly
- [ ] Detection mode set to `learning` for initial deployment
- [ ] Rate limits configured for your traffic patterns
- [ ] Log output configured (file or SIEM)

### Operations
- [ ] Health check configured in load balancer: `GET /health`
- [ ] Prometheus scraping configured (or log-based alerting set up)
- [ ] Database backups scheduled
- [ ] Runbook written for your team
- [ ] On-call rotation established for security alerts

### Testing
- [ ] Verified that legitimate traffic passes through cleanly
- [ ] Verified that basic SQLi and XSS payloads are blocked
- [ ] Load tested at 150% of expected peak traffic
- [ ] Failover tested (kill one WAF instance — traffic should shift)

---

## Deployment Modes

### Reverse Proxy (Most Common)

The WAF sits in front of your web server. DNS points to the WAF, not the backend.

```
Internet → WAF → Backend
```

Backend server should:
- Accept connections from WAF IP only (firewall rule)
- Trust the `X-Forwarded-For` header from the WAF for real client IPs

### Inline Transparent Mode

The WAF intercepts traffic at the network level. No DNS or application changes needed.

```
Internet → [Network] → WAF (transparent) → Backend
```

Requires network-level routing changes. Contact your network team.

### API Gateway Mode

The WAF acts as an API gateway with WAF capabilities built in. Routes different paths to different backends.

```yaml
server:
  routes:
    - path: /api/v1/
      upstream: http://api-service:8080
    - path: /static/
      upstream: http://cdn-origin:80
    - path: /
      upstream: http://web-service:3000
```

### Kubernetes Sidecar

Deployed as a sidecar container in Kubernetes pods. The WAF intercepts traffic destined for the main application container.

See [Kubernetes — Production](#kubernetes--production) below.

---

## Docker Compose — Production

Use the production Compose file, which differs from the dev file in key ways:
- No debug ports exposed
- Resource limits applied
- Health checks configured
- Restart policies set
- Secrets via Docker secrets (not environment variables)

```yaml
# docker-compose.prod.yml
version: "3.9"

services:
  waf-core:
    image: brainlesssecurity/brainless-waf:1.0.0   # Pin version, never :latest in prod
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
      # Do NOT expose 9091 (gRPC) or 6060 (pprof) publicly
    volumes:
      - ./config/production.yaml:/etc/brainless/config.yaml:ro
      - ./rules:/etc/brainless/rules:ro
      - tls_certs:/var/lib/brainless/acme
      - waf_logs:/var/log/brainless
    environment:
      - BRAINLESS_LOG_LEVEL=info
    secrets:
      - jwt_secret
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost/health"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 30s
    deploy:
      resources:
        limits:
          cpus: "4"
          memory: 2G
        reservations:
          cpus: "1"
          memory: 512M
    logging:
      driver: "json-file"
      options:
        max-size: "100m"
        max-file: "5"

  management:
    image: brainlesssecurity/brainless-waf-mgmt:1.0.0
    restart: unless-stopped
    ports:
      - "127.0.0.1:8000:8000"  # Bind to localhost ONLY
    environment:
      - DATABASE_URL_FILE=/run/secrets/db_url
    secrets:
      - jwt_secret
      - db_url
      - api_key_encryption_key
    depends_on:
      postgres:
        condition: service_healthy
    deploy:
      resources:
        limits:
          cpus: "1"
          memory: 512M

  dashboard:
    image: brainlesssecurity/brainless-waf-dashboard:1.0.0
    restart: unless-stopped
    ports:
      - "127.0.0.1:8080:80"   # Bind to localhost, put nginx in front
    deploy:
      resources:
        limits:
          cpus: "0.5"
          memory: 128M

  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    volumes:
      - pg_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=brainless
      - POSTGRES_USER=brainless
      - POSTGRES_PASSWORD_FILE=/run/secrets/db_password
    secrets:
      - db_password
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "brainless"]
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          memory: 1G

secrets:
  jwt_secret:
    external: true
  api_key_encryption_key:
    external: true
  db_url:
    external: true
  db_password:
    external: true

volumes:
  pg_data:
  tls_certs:
  waf_logs:
```

**Create Docker secrets before starting:**
```bash
echo "$(openssl rand -base64 48)" | docker secret create jwt_secret -
echo "$(openssl rand -base64 32)" | docker secret create api_key_encryption_key -
echo "your-db-password" | docker secret create db_password -
echo "postgresql://brainless:your-db-password@postgres:5432/brainless" | docker secret create db_url -
```

---

## Kubernetes — Production

### Recommended Resource Requests/Limits

```yaml
# values.yaml for Helm chart
wafCore:
  replicaCount: 2            # Minimum 2 for HA
  resources:
    requests:
      cpu: 500m
      memory: 512Mi
    limits:
      cpu: 4000m
      memory: 2Gi

management:
  replicaCount: 1            # Can be 2 for HA (stateless)
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 1000m
      memory: 512Mi

dashboard:
  replicaCount: 1
  resources:
    requests:
      cpu: 50m
      memory: 64Mi
    limits:
      cpu: 500m
      memory: 128Mi
```

### Pod Disruption Budget

Ensure at least 1 WAF core instance is always running during deployments:

```yaml
# deploy/k8s/pdb.yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: brainless-waf-core-pdb
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: brainless-waf-core
```

### Network Policies

Restrict which pods can talk to the WAF:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: brainless-waf-core
spec:
  podSelector:
    matchLabels:
      app: brainless-waf-core
  ingress:
    - ports:
        - port: 80
        - port: 443
      # Allow from anywhere (it's a public-facing WAF)
  egress:
    - ports:
        - port: 8080   # Upstream backend
    - ports:
        - port: 9091   # gRPC to management
    - ports:
        - port: 6379   # Redis (rate limiting)
```

### Rolling Updates

Configure your deployment strategy to ensure zero downtime:

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 1          # Spin up 1 new pod before terminating old
    maxUnavailable: 0    # Never terminate before replacement is ready
```

---

## High Availability Setup

For production systems where downtime is unacceptable:

```
                 ┌─────────────────────────┐
    Internet     │   Cloud Load Balancer    │
    ────────────►│   (health check /health) │
                 └─────────┬───────────────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
         ┌────▼───┐   ┌────▼───┐   ┌───▼────┐
         │ WAF-1  │   │ WAF-2  │   │ WAF-3  │
         │ (AZ-a) │   │ (AZ-b) │   │ (AZ-c) │
         └────┬───┘   └────┬───┘   └───┬────┘
              │            │            │
              └────────────┼────────────┘
                           │ gRPC (config sync)
                    ┌──────▼───────┐
                    │  Mgmt API    │
                    │  (2 replicas)│
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
         ┌────▼───┐   ┌────▼───┐   ┌───▼────┐
         │PostgreSQL  │ Redis   │   │ Redis  │
         │(primary)│  │(primary)│  │(replica)│
         └────────┘   └────────┘   └────────┘
```

**Key HA requirements:**
- WAF core instances are stateless — deploy 3+ across availability zones
- Rate limit counters in Redis (not in-memory) — shared across all WAF instances
- PostgreSQL with streaming replication (or a managed DB service)
- Load balancer health check on `GET /health` (remove unhealthy instances within 10s)

---

## Network Architecture

### Recommended Firewall Rules

| Source | Destination | Port | Action | Reason |
|--------|------------|------|--------|--------|
| Internet | WAF | 80, 443 | ALLOW | Public traffic |
| Internet | WAF | Any other | DENY | Lock down attack surface |
| WAF | Backend | 8080 (or your port) | ALLOW | WAF → app |
| Backend | Internet | Any | DENY | Backend should not be internet-facing |
| Admin IPs | WAF | 8080 (dashboard) | ALLOW | Dashboard access |
| Admin IPs | WAF | 8000 (mgmt API) | ALLOW | API access |
| Internet | WAF | 8080, 8000 | DENY | Management not public |
| WAF | Redis | 6379 | ALLOW | Rate limit counters |
| Mgmt API | WAF | 9091 | ALLOW | gRPC config push |
| Prometheus | WAF | 9113 | ALLOW | Metrics scraping |

---

## TLS in Production

### Option A: ACME / Let's Encrypt (Recommended for internet-facing)

```yaml
server:
  tls:
    acme: true
    acme_email: ops@yourcompany.com
    acme_domains:
      - api.yourcompany.com
      - www.yourcompany.com
    acme_cache_dir: /var/lib/brainless/acme
```

Certificates auto-renew 30 days before expiry. Ensure port 80 is accessible from the internet (ACME HTTP-01 challenge).

### Option B: Corporate CA / Custom Certificates

```yaml
server:
  tls:
    cert: /etc/ssl/certs/waf.pem
    key: /etc/ssl/private/waf.key
```

Set up a cron job or cert-manager to rotate certificates before expiry and send `SIGHUP` to reload:

```bash
# Rotate certificate
cp new-cert.pem /etc/ssl/certs/waf.pem
cp new-key.pem /etc/ssl/private/waf.key
kill -HUP $(pgrep brainless-waf)   # Graceful reload
```

### Mutual TLS (for API clients)

If your API clients support mTLS, require client certificates for the Management API:

```yaml
management:
  tls:
    require_client_cert: true
    client_ca: /etc/ssl/certs/client-ca.pem
```

---

## Secrets Management

Never put secrets in config files or environment variables that are visible in process lists.

### HashiCorp Vault (Recommended)

```yaml
# config/production.yaml
secrets:
  provider: vault
  vault_addr: https://vault.internal:8200
  vault_role: brainless-waf
  paths:
    jwt_secret: secret/brainless-waf/jwt_secret
    db_password: secret/brainless-waf/db_password
    api_key_encryption_key: secret/brainless-waf/api_key_encryption_key
```

### Kubernetes Secrets

```bash
kubectl create secret generic brainless-waf-secrets \
  --from-literal=jwt-secret="$(openssl rand -base64 48)" \
  --from-literal=api-key-encryption-key="$(openssl rand -base64 32)" \
  --from-literal=db-password="$(openssl rand -base64 24)" \
  -n brainless-waf
```

---

## Monitoring Stack

### Prometheus + Grafana (Recommended)

```yaml
# prometheus.yml scrape config
scrape_configs:
  - job_name: 'brainless-waf'
    static_configs:
      - targets: ['waf-host:9113']
    metrics_path: /metrics
    scrape_interval: 15s
```

Import the Brainless WAF Grafana dashboard:
```bash
# Dashboard ID for Grafana.com marketplace
# Or import from: deploy/grafana/dashboard.json
```

### Alerting Rules

```yaml
# prometheus-alerts.yml
groups:
  - name: brainless-waf
    rules:
      - alert: WAFHighLatency
        expr: histogram_quantile(0.99, brainless_latency_seconds_bucket) > 0.05
        for: 5m
        annotations:
          summary: "WAF p99 latency > 50ms"

      - alert: WAFHighBlockRate
        expr: rate(brainless_blocked_total[5m]) / rate(brainless_requests_total[5m]) > 0.1
        for: 10m
        annotations:
          summary: "WAF blocking >10% of traffic — possible attack or misconfiguration"

      - alert: WAFInstanceDown
        expr: up{job="brainless-waf"} == 0
        for: 1m
        annotations:
          summary: "A WAF instance is down"

      - alert: WAFCertExpiringSoon
        expr: brainless_tls_cert_expiry_seconds < 86400 * 14
        annotations:
          summary: "TLS certificate expires in less than 14 days"
```

---

## Backup & Recovery

### What to Back Up

| Data | Location | Frequency | Method |
|------|---------|-----------|--------|
| WAF configuration | PostgreSQL `config` table | Continuous (WAL) | pg_basebackup |
| Custom rules | PostgreSQL `rules` table | Continuous (WAL) | pg_basebackup |
| IP lists | PostgreSQL | Continuous (WAL) | pg_basebackup |
| Event logs | PostgreSQL or log files | Daily snapshot | pg_dump or rsync |
| TLS certificates | `/var/lib/brainless/acme` | Weekly | File backup |

### PostgreSQL Backup

```bash
# Daily dump
pg_dump -U brainless brainless | gzip > /backups/brainless-$(date +%Y%m%d).sql.gz

# Verify backup
zcat /backups/brainless-20250314.sql.gz | psql -U brainless brainless_restore

# Retention: keep 30 days of daily backups
find /backups -name "*.sql.gz" -mtime +30 -delete
```

### Recovery Procedure

```bash
# 1. Stop WAF
docker compose stop waf-core management

# 2. Restore database
zcat /backups/brainless-20250314.sql.gz | psql -U brainless brainless

# 3. Restart
docker compose start management
sleep 10  # Wait for management API to be ready
docker compose start waf-core
```

---

## Incident Response

### Playbook: Active Attack

1. **Identify source** — Dashboard → Events, filter by `action=block`, group by `src_ip`
2. **Block the attacker** — IP Management → Blocklist → Add IP/CIDR
3. **Check for successful requests** — Look for events with same IP that were NOT blocked
4. **Escalate if needed** — If attack is DDoS-scale, contact upstream provider for null routing

### Playbook: High False Positive Rate

1. **Switch to detect mode** — Settings → WAF Mode → Detect (no blocking)
2. **Identify the rule** — Events → filter `action=detect`, find the rule causing most FPs
3. **Add exception** — Rules → find rule → Add Exception
4. **Switch back to block** — Settings → WAF Mode → Block
5. **Monitor** — verify FP rate drops

### Playbook: WAF Instance Down

```bash
# Check container status
docker compose ps

# Check for OOM kill
dmesg | grep -i "oom" | tail -20

# Restart
docker compose restart waf-core

# If keeps crashing — check logs and open GitHub issue
docker compose logs waf-core --tail=100 > incident-$(date +%Y%m%d-%H%M).log
```

---

## Maintenance Procedures

### Updating Brainless WAF

```bash
# 1. Check the changelog for breaking changes
open https://github.com/brainless-security/brainless-waf/blob/main/CHANGELOG.md

# 2. Update in staging first, test for 24 hours

# 3. Production update (Docker)
docker compose pull
docker compose up -d --no-deps --build waf-core management dashboard

# 4. Verify health
curl http://localhost/health
docker compose logs waf-core --tail=20
```

### Rotating Secrets

```bash
# 1. Generate new JWT secret
NEW_SECRET=$(openssl rand -base64 48)

# 2. Update in secrets manager
echo "$NEW_SECRET" | docker secret create jwt_secret_new -
docker service update --secret-rm jwt_secret --secret-add jwt_secret_new brainless_management

# 3. All existing JWT tokens are now invalid — users must log in again
# 4. Remove old secret after confirming everything works
docker secret rm jwt_secret
```

### Rule Updates

```bash
# Check for available updates
brainless-ctl rules check-updates

# Apply updates (in detect mode first to check for new FPs)
brainless-ctl config set detection.mode detect
brainless-ctl rules update
# Monitor for 1 hour
brainless-ctl config set detection.mode block
```
