# The Club Enhancement Plan - FOSS Community Stack Implementation

## 🎯 Overview

This implementation plan adds battle-tested, community-vetted FOSS tools to The Club home server while preserving the existing Caddy/Homer/Go architecture. All new services are based on r/selfhosted and r/homelab 2026 community standards.

**Implementation Time**: ~4 weeks (phased rollout)
**Total Additional Services**: 11 FOSS tools
**Resource Impact**: ~2-3GB RAM total (all phases)

---

## 📋 Prerequisites

Before starting, ensure you have:

- [ ] Docker and Docker Compose installed
- [ ] The Club base stack running (`make docker-prod`)
- [ ] Basic understanding of Docker Compose
- [ ] Backup of current configuration (`make backup`)
- [ ] At least 4GB RAM available for full stack

---

## 🚀 PHASE 1: Monitoring & Protection (Week 1)

**Goal**: Get visibility into your services and protect configurations
**Time**: ~2 hours
**Priority**: CRITICAL
**Resource Impact**: +384MB RAM

### Services Added

1. **Uptime Kuma** - Status monitoring dashboard (256MB)
2. **Backup Script** - Automated configuration backups
3. **Resource Limits** - Prevent resource exhaustion

### Implementation Steps

#### Step 1.1: Deploy Monitoring Stack

```bash
# Deploy Uptime Kuma
make stack-monitoring

# Verify deployment
docker ps | grep uptime-kuma

# Access the dashboard
# Open: https://localhost/status
# Create your admin account on first visit
```

#### Step 1.2: Configure Uptime Kuma Monitors

1. Access https://localhost/status
2. Create admin account
3. Add monitors:
   - **HTTP Monitor**: https://localhost/health (Main Health)
   - **HTTP Monitor**: https://localhost/test (Go Server)
   - **HTTP Monitor**: https://localhost/ (Homer Dashboard)
   - **Docker Container**: caddy (if enabled)
   - **Docker Container**: homer (if enabled)
   - **Docker Container**: go-test-server (if enabled)

4. Configure notifications:
   - Settings → Notifications
   - Add notification method (Email, Discord, Telegram, etc.)
   - Test notifications

#### Step 1.3: Run Your First Backup

```bash
# Run backup script
make backup

# Verify backup created
ls -lh /backups/theclub-*/

# View backup manifest
cat /backups/theclub-*/manifest.txt
```

#### Step 1.4: Schedule Automated Backups

```bash
# Add to crontab
crontab -e

# Add this line (runs daily at 3 AM)
0 3 * * * cd /home/user/TheClub && make backup >> /var/log/theclub-backup.log 2>&1
```

### Validation

- [ ] Uptime Kuma accessible at https://localhost/status
- [ ] All monitors showing green status
- [ ] Backup created successfully
- [ ] Resource limits applied (verify with `docker stats`)

---

## 🔐 PHASE 2: Security & Authentication (Week 2)

**Goal**: Secure remote access and centralized authentication
**Time**: ~3 hours
**Priority**: HIGH
**Resource Impact**: +512MB RAM

### Services Added

1. **WireGuard** - Modern VPN for secure remote access (256MB)
2. **Authelia** - Lightweight SSO authentication (128MB)
3. **Redis** - Session storage for Authelia (128MB)

### Implementation Steps

#### Step 2.1: Configure Authelia

```bash
# Generate secure secrets
docker run --rm authelia/authelia:latest authelia crypto rand --length 64

# Edit authelia/configuration.yml
# Update these values with generated secrets:
#   - session.secret
#   - storage.encryption_key

# Update domain
# Change "yourdomain.com" to your actual domain throughout the file
```

#### Step 2.2: Create User Account

```bash
# Generate password hash
docker run --rm authelia/authelia:latest \
  authelia crypto hash generate argon2 --password 'your_secure_password'

# Edit authelia/users_database.yml
# Replace the hash for the admin user with your generated hash
# Change email and displayname as needed
```

#### Step 2.3: Deploy Security Stack

```bash
# Deploy WireGuard and Authelia
make stack-security

# Wait for services to start
sleep 30

# Check status
docker ps | grep -E 'wireguard|authelia|redis'
```

#### Step 2.4: Configure WireGuard

```bash
# Get WireGuard configuration for peer 1
docker exec wireguard cat /config/peer1/peer1.conf

# Or get QR code for mobile
docker exec wireguard cat /config/peer1/peer1.png

# Download config to your device
docker cp wireguard:/config/peer1/peer1.conf ~/Downloads/

# Import into WireGuard client:
#   - Desktop: https://www.wireguard.com/install/
#   - Mobile: App store (WireGuard)
```

#### Step 2.5: Test Authelia

```bash
# Access Authelia
# Open: https://localhost/auth

# Login with credentials from users_database.yml
# Default: admin / changeme (if you didn't change it)

# Enable 2FA
# Follow the setup wizard to configure TOTP
```

#### Step 2.6: Protect NAS Services (Optional - When Ready)

Edit `authelia/configuration.yml` and uncomment the access control rules:

```yaml
access_control:
  rules:
    - domain: "*.yourdomain.com"
      policy: two_factor
      resources:
        - "^/media.*"
        - "^/drive.*"
        - "^/arcade.*"
        - "^/grafana.*"
```

Then update Caddyfile to use Authelia forward_auth (example):

```caddyfile
# Protect NAS services
@protected {
    path /media* /drive* /arcade* /grafana*
}
handle @protected {
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.yourdomain.com
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }
}
```

### Validation

- [ ] Authelia accessible at https://localhost/auth
- [ ] Can login with configured credentials
- [ ] 2FA working (if enabled)
- [ ] WireGuard VPN connects successfully
- [ ] Can access services through VPN

---

## 📊 PHASE 3: Observability (Week 3)

**Goal**: Deep metrics and beautiful dashboards
**Time**: ~2 hours
**Priority**: MEDIUM
**Resource Impact**: +1088MB RAM

### Services Added

1. **Prometheus** - Time-series metrics database (512MB)
2. **Grafana** - Visualization dashboards (256MB)
3. **Node Exporter** - System metrics (64MB)
4. **cAdvisor** - Container metrics (256MB)

### Implementation Steps

#### Step 3.1: Deploy Observability Stack

```bash
# Deploy Prometheus and Grafana
make stack-observability

# Wait for services to start
sleep 60

# Check status
docker ps | grep -E 'prometheus|grafana|node-exporter|cadvisor'
```

#### Step 3.2: Access Grafana

```bash
# Open Grafana
# URL: https://localhost/grafana
# Username: admin
# Password: changeme

# Change password on first login
# Click: Admin (gear icon) → Profile → Change password
```

#### Step 3.3: Import Community Dashboards

In Grafana web UI:

1. **Import Node Exporter Dashboard**:
   - Dashboards → New → Import
   - Dashboard ID: `1860`
   - Click "Load"
   - Select "Prometheus" as datasource
   - Click "Import"

2. **Import Docker Dashboard**:
   - Dashboards → New → Import
   - Dashboard ID: `893`
   - Click "Load"
   - Select "Prometheus" as datasource
   - Click "Import"

3. **Import Caddy Dashboard** (if available):
   - Dashboards → New → Import
   - Dashboard ID: `14280`
   - Click "Load"
   - Select "Prometheus" as datasource
   - Click "Import"

#### Step 3.4: Configure Alerts (Optional)

1. In Grafana, go to Alerting → Notification channels
2. Add your preferred notification method:
   - Email
   - Slack
   - Discord
   - Telegram
   - Webhook

3. Set up alert rules:
   - High CPU usage (>80% for 5m)
   - High memory usage (>90%)
   - Disk space low (<10% free)
   - Container down

#### Step 3.5: Verify Metrics Collection

```bash
# Check Prometheus targets
# Open: http://localhost:9090/targets
# All targets should show "UP" status

# Query some metrics in Prometheus
# Open: http://localhost:9090/graph
# Try queries:
#   - up
#   - node_cpu_seconds_total
#   - container_memory_usage_bytes
#   - caddy_http_requests_total
```

### Validation

- [ ] Grafana accessible at https://localhost/grafana
- [ ] All dashboards showing data
- [ ] Prometheus targets all "UP"
- [ ] Node Exporter metrics visible
- [ ] Container metrics visible

---

## 🤖 PHASE 4: Automation (Week 4)

**Goal**: Automated backups, updates, and management
**Time**: ~2 hours
**Priority**: MEDIUM
**Resource Impact**: +896MB RAM

### Services Added

1. **Kopia** - Automated encrypted backups (512MB)
2. **Watchtower** - Automatic container updates (128MB)
3. **Dockge** - Docker Compose management UI (256MB)

### Implementation Steps

#### Step 4.1: Prepare Backup Destination

```bash
# Create backup directory (or mount external drive)
sudo mkdir -p /mnt/backup
sudo chown $(whoami):$(whoami) /mnt/backup

# Or use NAS mount (example)
# sudo mount -t nfs synology.local:/volume1/backups /mnt/backup
```

#### Step 4.2: Configure Kopia

```bash
# Edit docker-compose.automation.yml
# Update KOPIA_PASSWORD with a secure password

# Update backup destination if not using /mnt/backup
# Change the volume mount: - /your/path:/repository
```

#### Step 4.3: Deploy Automation Stack

```bash
# Create stacks directory for Dockge
sudo mkdir -p /opt/stacks
sudo chown -R $(whoami):$(whoami) /opt/stacks

# Deploy automation services
make stack-automation

# Wait for services to start
sleep 60

# Check status
docker ps | grep -E 'kopia|watchtower|dockge'
```

#### Step 4.4: Initialize Kopia Repository

```bash
# Access Kopia UI
# Open: https://localhost/backup

# Initial setup:
# 1. Create repository at /repository
# 2. Set encryption password (matches KOPIA_PASSWORD)
# 3. Configure backup policy:
#    - Frequency: Daily at 3 AM
#    - Retention: 7 daily, 4 weekly, 6 monthly
#    - Compression: zstd
```

#### Step 4.5: Create Backup Snapshots

In Kopia UI:

1. Navigate to "Snapshots"
2. Click "New Snapshot"
3. Select directories to backup:
   - `/data/theclub` (configuration files)
   - `/data/caddy_data` (SSL certificates)
   - `/data/uptime-kuma-data` (monitoring config)
   - `/data/grafana-data` (dashboards)

4. Create backup policy
5. Schedule regular snapshots

#### Step 4.6: Configure Watchtower Notifications (Optional)

Edit `docker-compose.automation.yml` and uncomment email notification section:

```yaml
environment:
  - WATCHTOWER_NOTIFICATIONS=email
  - WATCHTOWER_NOTIFICATION_EMAIL_FROM=watchtower@theclub.local
  - WATCHTOWER_NOTIFICATION_EMAIL_TO=your@email.com
  # Configure SMTP settings
```

Restart Watchtower:
```bash
docker compose -f docker-compose.automation.yml restart watchtower
```

#### Step 4.7: Access Dockge

```bash
# Open Dockge
# URL: https://localhost/dockge

# Create admin account on first visit

# Explore features:
# - View all compose stacks
# - Start/stop/restart services
# - View logs in real-time
# - Edit compose files
# - Terminal access to containers
```

### Validation

- [ ] Kopia accessible at https://localhost/backup
- [ ] Kopia repository initialized
- [ ] First backup snapshot created
- [ ] Watchtower running (check logs: `docker logs watchtower`)
- [ ] Dockge accessible at https://localhost/dockge
- [ ] All stacks visible in Dockge

---

## 🎉 PHASE 5: Full Stack Deployment (Optional)

Deploy all enhancements at once (if you have sufficient resources):

```bash
# Deploy entire enhancement stack
make stack-all

# This deploys:
# - Monitoring (Uptime Kuma)
# - Security (WireGuard + Authelia + Redis)
# - Observability (Prometheus + Grafana + Exporters)
# - Automation (Kopia + Watchtower + Dockge)

# Total resource usage: ~2.9GB RAM
```

---

## 📊 Final Architecture

After all phases, your architecture will be:

```
Internet/LAN
    ↓
[WireGuard VPN] ←─ Remote Access
    ↓
Caddy (SSL + Security Headers)
    ├─→ / (Homer Dashboard)
    ├─→ /test (Go/HTMX Server)
    ├─→ /status (Uptime Kuma)
    ├─→ /grafana (Metrics Dashboards)
    ├─→ /auth (Authelia SSO) ←─ Protects services
    ├─→ /dockge (Stack Management)
    ├─→ /backup (Kopia Backups)
    └─→ [Future: /media, /drive, /arcade - NAS Services]

Background Services:
    ├─→ Prometheus (Metrics Collection)
    ├─→ Node Exporter (System Metrics)
    ├─→ cAdvisor (Container Metrics)
    ├─→ Redis (Authelia Sessions)
    └─→ Watchtower (Auto Updates)
```

---

## 📖 Service Access URLs

After full deployment:

| Service | URL | Default Credentials |
|---------|-----|---------------------|
| **Homer Dashboard** | https://localhost/ | N/A |
| **Go/HTMX Server** | https://localhost/test | N/A |
| **Uptime Kuma** | https://localhost/status | Set on first visit |
| **Grafana** | https://localhost/grafana | admin/changeme |
| **Authelia** | https://localhost/auth | See users_database.yml |
| **Dockge** | https://localhost/dockge | Set on first visit |
| **Kopia** | https://localhost/backup | See env var |
| **Prometheus** | http://localhost:9090 | N/A (internal) |

---

## 🔧 Maintenance Tasks

### Daily (Automated)
- Watchtower checks for updates (4 AM)
- Kopia creates backup snapshots (3 AM)
- Uptime Kuma monitors services (continuous)
- Prometheus collects metrics (15s interval)

### Weekly
- Review Uptime Kuma alerts
- Check Grafana dashboards for anomalies
- Review Watchtower update logs

### Monthly
- Test backup restoration
- Review and update Authelia users
- Update Grafana dashboards
- Check disk usage for metrics/logs

### Quarterly
- Security audit (review Authelia logs)
- Update WireGuard peer configurations
- Review and optimize Prometheus retention
- Test disaster recovery procedures

---

## 🆘 Troubleshooting

### Service Won't Start

```bash
# Check logs
docker logs <container_name>

# Check resources
docker stats

# Check network
docker network inspect theclub_club-network

# Restart service
docker compose -f docker-compose.<stack>.yml restart <service>
```

### Can't Access Service

```bash
# Check if service is running
docker ps | grep <service_name>

# Check Caddy routing
docker logs caddy | grep <service_name>

# Test internal connectivity
docker exec caddy wget -O- http://<service>:<port>
```

### Out of Memory

```bash
# Check memory usage
docker stats --no-stream

# Identify memory hogs
docker stats --no-stream --format "table {{.Name}}\t{{.MemUsage}}" | sort -k2 -h

# Adjust resource limits in compose files
# Edit deploy.resources.limits.memory values
```

### Backup Failed

```bash
# Check Kopia logs
docker logs kopia

# Verify backup destination is writable
ls -la /mnt/backup

# Check Kopia repository status
docker exec kopia kopia repository status

# Recreate repository if corrupted
docker exec kopia kopia repository create filesystem --path=/repository
```

---

## 🔐 Security Hardening Checklist

- [ ] Change all default passwords
  - [ ] Grafana admin password
  - [ ] Kopia encryption password
  - [ ] Authelia admin user password
  - [ ] Uptime Kuma admin password
  - [ ] Dockge admin password

- [ ] Generate secure secrets
  - [ ] Authelia session.secret
  - [ ] Authelia storage.encryption_key

- [ ] Configure SSL for production
  - [ ] Update Caddyfile with real domain
  - [ ] Verify Let's Encrypt certificates

- [ ] Enable 2FA
  - [ ] Authelia TOTP
  - [ ] Consider hardware keys (YubiKey)

- [ ] Configure backup encryption
  - [ ] Kopia repository password
  - [ ] Off-site backup location

- [ ] Set up monitoring alerts
  - [ ] Uptime Kuma notifications
  - [ ] Grafana alert rules
  - [ ] Email/Discord/Telegram

- [ ] Firewall configuration
  - [ ] Close unused ports
  - [ ] Allow only WireGuard (51820/udp)
  - [ ] Internal services not exposed

---

## 📚 Additional Resources

### Official Documentation
- [Uptime Kuma](https://github.com/louislam/uptime-kuma)
- [Authelia](https://www.authelia.com/)
- [WireGuard](https://www.wireguard.com/)
- [Prometheus](https://prometheus.io/docs/)
- [Grafana](https://grafana.com/docs/)
- [Kopia](https://kopia.io/docs/)
- [Dockge](https://github.com/louislam/dockge)

### Community Resources
- [r/selfhosted](https://www.reddit.com/r/selfhosted/)
- [r/homelab](https://www.reddit.com/r/homelab/)
- [Awesome Self-Hosted](https://github.com/awesome-selfhosted/awesome-selfhosted)
- [Perfect Media Server](https://perfectmediaserver.com/)

### Grafana Dashboards
- [Node Exporter Full (1860)](https://grafana.com/grafana/dashboards/1860)
- [Docker & Host Metrics (893)](https://grafana.com/grafana/dashboards/893)
- [Caddy Metrics (14280)](https://grafana.com/grafana/dashboards/14280)

---

## 🎯 Success Criteria

You'll know the implementation is successful when:

- ✅ All services accessible via HTTPS with valid certificates
- ✅ Uptime Kuma shows all monitors green
- ✅ Grafana dashboards displaying metrics for all services
- ✅ Automated backups running daily
- ✅ WireGuard VPN allows secure remote access
- ✅ Authelia protecting designated services
- ✅ Watchtower keeping containers updated
- ✅ Dockge providing easy stack management
- ✅ Resource usage within acceptable limits (<4GB RAM total)
- ✅ All default passwords changed
- ✅ Monitoring alerts configured and working

---

**Next Steps**: Start with Phase 1 (Monitoring & Protection) and proceed through the phases systematically. Don't rush - take time to understand each service before moving to the next phase.

**Questions?** See the troubleshooting section or consult the official documentation for each service.
