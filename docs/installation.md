# Installation Guide

This guide covers all supported installation methods for Brainless WAF.

---

## Table of Contents

- [System Requirements](#system-requirements)
- [Method 1: Docker Compose (Recommended)](#method-1-docker-compose-recommended)
- [Method 2: Kubernetes (Helm)](#method-2-kubernetes-helm)
- [Method 3: Binary (Bare Metal)](#method-3-binary-bare-metal)
- [Method 4: Build from Source](#method-4-build-from-source)
- [Post-Installation](#post-installation)
- [Upgrading](#upgrading)
- [Uninstalling](#uninstalling)

---

## System Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 2 cores | 8+ cores |
| RAM | 2 GB | 16 GB |
| Disk | 10 GB SSD | 100 GB NVMe |
| OS | Linux (kernel 4.19+) | Ubuntu 22.04 / Debian 12 |
| Architecture | amd64 | amd64 or arm64 |

**Ports needed:**
- `80` — HTTP ingress
- `443` — HTTPS ingress
- `8080` — Dashboard (can be changed)
- `8000` — Management API (should be restricted to internal network)

---

## Method 1: Docker Compose (Recommended)

Best for: single-server deployments, homelab, getting started quickly.

### Step 1 — Install Docker

```bash
# Ubuntu / Debian
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker
```

### Step 2 — Clone and configure

```bash
git clone https://github.com/brainless-security/brainless-waf.git
cd brainless-waf

# Create your config from the default template
cp config/default.yaml config/local.yaml
```

Edit `config/local.yaml` — at minimum, set your upstream server:

```yaml
server:
  upstream: http://your-backend-server:8080
```

### Step 3 — Start

```bash
docker compose up -d
```

### Step 4 — Verify

```bash
docker compose ps
# All services should show "running"

curl http://localhost/health
# {"status":"ok","version":"1.0.0"}
```

Open `http://your-server-ip:8080` for the dashboard.

**Default credentials:** `admin` / `changeme` — **change this immediately**.

### Step 5 — Secure the deployment

```bash
# Change admin password
docker exec -it brainless-waf-management brainless-ctl passwd admin

# Restrict Management API to localhost (edit docker-compose.yml)
# Change:  - "8000:8000"
# To:      - "127.0.0.1:8000:8000"
```

---

## Method 2: Kubernetes (Helm)

Best for: production deployments, cloud environments, teams already using Kubernetes.

### Prerequisites

- Kubernetes 1.24+
- Helm 3.10+
- `kubectl` configured for your cluster

### Step 1 — Add the Helm repository

```bash
helm repo add brainless https://charts.brainless-security.io
helm repo update
```

### Step 2 — Create a values file

```bash
helm show values brainless/brainless-waf > values.yaml
```

Edit `values.yaml` for your environment. Key settings:

```yaml
# values.yaml

config:
  upstream: http://my-backend-service:8080
  detection:
    mode: learning          # Start in learning mode
    paranoia_level: 2

ingress:
  enabled: true
  className: nginx
  hosts:
    - host: waf.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: waf-tls
      hosts:
        - waf.example.com

resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 2000m
    memory: 2Gi

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

postgresql:
  enabled: true    # Use bundled PostgreSQL (disable for external DB)

redis:
  enabled: true    # Required for multi-replica rate limiting
```

### Step 3 — Install

```bash
helm install brainless-waf brainless/brainless-waf \
  --namespace brainless-waf \
  --create-namespace \
  --values values.yaml
```

### Step 4 — Verify

```bash
kubectl get pods -n brainless-waf
# NAME                                    READY   STATUS    RESTARTS
# brainless-waf-core-7d9f8b-xxxx         1/1     Running   0
# brainless-waf-management-5c6d9-xxxx    1/1     Running   0
# brainless-waf-dashboard-8b7f4-xxxx     1/1     Running   0

kubectl logs -n brainless-waf deploy/brainless-waf-core -f
```

### Upgrading Helm release

```bash
helm repo update
helm upgrade brainless-waf brainless/brainless-waf \
  --namespace brainless-waf \
  --values values.yaml
```

---

## Method 3: Binary (Bare Metal)

Best for: environments without Docker, custom system integration.

### Step 1 — Download

```bash
# Get the latest version
VERSION=$(curl -s https://api.github.com/repos/brainless-security/brainless-waf/releases/latest | grep tag_name | cut -d'"' -f4)

# Linux amd64
curl -LO "https://github.com/brainless-security/brainless-waf/releases/download/${VERSION}/brainless-waf_linux_amd64.tar.gz"

# Verify checksum
curl -LO "https://github.com/brainless-security/brainless-waf/releases/download/${VERSION}/checksums.txt"
sha256sum --check checksums.txt
```

### Step 2 — Install

```bash
tar -xzf brainless-waf_linux_amd64.tar.gz
sudo mv brainless-waf brainless-ctl /usr/local/bin/
sudo chmod +x /usr/local/bin/brainless-waf /usr/local/bin/brainless-ctl
```

### Step 3 — Configure

```bash
sudo mkdir -p /etc/brainless-waf /var/log/brainless-waf /var/lib/brainless-waf
sudo cp config/default.yaml /etc/brainless-waf/config.yaml

# Edit the config
sudo nano /etc/brainless-waf/config.yaml
```

### Step 4 — PostgreSQL setup

```bash
# Install PostgreSQL (Ubuntu)
sudo apt install postgresql postgresql-contrib

# Create database and user
sudo -u postgres psql <<EOF
CREATE USER brainless WITH PASSWORD 'your-strong-password';
CREATE DATABASE brainless OWNER brainless;
EOF

# Run migrations
DATABASE_URL="postgresql://brainless:your-strong-password@localhost/brainless" \
  brainless-ctl db migrate
```

### Step 5 — Systemd service

```bash
# Install systemd service
sudo brainless-ctl install-service

# Start and enable
sudo systemctl enable --now brainless-waf

# Check status
sudo systemctl status brainless-waf
```

The service file is installed to `/etc/systemd/system/brainless-waf.service`.

---

## Method 4: Build from Source

Best for: development, customization, contributing to the project.

See [DEVELOPMENT.md](../DEVELOPMENT.md) for the full guide.

```bash
git clone https://github.com/brainless-security/brainless-waf.git
cd brainless-waf
make deps
make build
# Binaries output to ./bin/
```

---

## Post-Installation

### Required: Change default credentials

```bash
# Via CLI
brainless-ctl passwd admin

# Via API
curl -X POST http://localhost:8000/api/v1/auth/change-password \
  -H "Authorization: Bearer <token>" \
  -d '{"current_password": "changeme", "new_password": "your-strong-password"}'
```

### Recommended: Start in learning mode

```yaml
# config/local.yaml
detection:
  mode: learning
```

Let the WAF observe your traffic for 72 hours, then:
1. Review the learning report: `brainless-ctl learning-report`
2. Switch to `detect` mode and review false positives
3. Add exceptions for legitimate traffic patterns
4. Switch to `block` mode

### Enable TLS

```yaml
server:
  listen_tls: 0.0.0.0:443
  tls:
    acme: true                  # Automatic Let's Encrypt
    acme_email: you@example.com
    acme_domains:
      - waf.example.com
```

Or with your own certificate:

```yaml
server:
  tls:
    cert: /etc/ssl/certs/your-cert.pem
    key: /etc/ssl/private/your-key.pem
```

---

## Upgrading

### Docker Compose

```bash
docker compose pull
docker compose up -d
```

Migrations run automatically on startup.

### Helm

```bash
helm repo update
helm upgrade brainless-waf brainless/brainless-waf \
  --namespace brainless-waf \
  --values values.yaml
```

### Binary

```bash
# Download new version (same steps as install)
# Migrations
brainless-ctl db migrate
sudo systemctl restart brainless-waf
```

---

## Uninstalling

### Docker Compose

```bash
docker compose down -v    # -v removes volumes (deletes all data)
```

### Helm

```bash
helm uninstall brainless-waf -n brainless-waf
kubectl delete namespace brainless-waf
```

### Binary

```bash
sudo systemctl disable --now brainless-waf
sudo rm /usr/local/bin/brainless-waf /usr/local/bin/brainless-ctl
sudo rm -rf /etc/brainless-waf /var/log/brainless-waf /var/lib/brainless-waf
sudo systemctl daemon-reload
```
