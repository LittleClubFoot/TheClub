# The Club - FOSS Community Enhancements

## Quick Start

```bash
# View all enhancement commands
make help

# Deploy monitoring (Phase 1 - START HERE)
make stack-monitoring

# Deploy all enhancements at once
make stack-all

# Stop all enhancements
make stack-stop

# Run backup
make backup
```

---

## 🎯 What's New

The Club now includes battle-tested FOSS tools from the 2026 r/selfhosted community:

### Phase 1: Monitoring & Protection
- **Uptime Kuma** - Beautiful status monitoring at `/status`
- **Backup Script** - Automated config backups via `make backup`
- **Resource Limits** - All services have memory/CPU limits

### Phase 2: Security & Authentication
- **WireGuard** - Modern VPN (replaces OpenVPN)
- **Authelia** - Lightweight SSO (only 30MB RAM)
- **Redis** - Session storage

### Phase 3: Observability
- **Prometheus** - Time-series metrics database
- **Grafana** - Beautiful dashboards at `/grafana`
- **Node Exporter** - System metrics
- **cAdvisor** - Container metrics

### Phase 4: Automation
- **Kopia** - Encrypted backups at `/backup`
- **Watchtower** - Auto-updates (4 AM daily)
- **Dockge** - Compose stack management at `/dockge`

---

## 📊 Service URLs

| Service | URL | Credentials |
|---------|-----|-------------|
| Dashboard | https://localhost/ | - |
| Monitoring | https://localhost/status | Set on first visit |
| Metrics | https://localhost/grafana | admin/changeme |
| Auth | https://localhost/auth | See authelia/users_database.yml |
| Management | https://localhost/dockge | Set on first visit |
| Backup | https://localhost/backup | See compose file |

---

## 🚀 Quick Deployment

### Option 1: Phased (Recommended)
```bash
# Week 1: Monitoring
make stack-monitoring

# Week 2: Security
make stack-security

# Week 3: Observability
make stack-observability

# Week 4: Automation
make stack-automation
```

### Option 2: All at Once (Requires 4GB+ RAM)
```bash
make stack-all
```

---

## 📋 Pre-Deployment Checklist

Before deploying:

- [ ] Backup current config: `make backup`
- [ ] Update Authelia secrets (see authelia/configuration.yml)
- [ ] Update Authelia user passwords (see authelia/users_database.yml)
- [ ] Create backup destination: `mkdir -p /mnt/backup`
- [ ] Update domain in Caddyfile (replace 'localhost')
- [ ] Ensure 4GB+ RAM available

---

## 🔒 Security First Steps

After deployment, immediately:

1. **Change Default Passwords**
   ```bash
   # Grafana: https://localhost/grafana (admin/changeme)
   # Uptime Kuma: https://localhost/status (set on first visit)
   # Dockge: https://localhost/dockge (set on first visit)
   ```

2. **Generate Authelia Secrets**
   ```bash
   docker run --rm authelia/authelia:latest authelia crypto rand --length 64
   # Update authelia/configuration.yml with generated values
   ```

3. **Configure WireGuard**
   ```bash
   # Get VPN config
   docker exec wireguard cat /config/peer1/peer1.conf
   ```

4. **Set Up Notifications**
   - Uptime Kuma: Configure in web UI
   - Grafana: Alerting → Notification channels

---

## 📊 Resource Usage

| Phase | Services | RAM Usage |
|-------|----------|-----------|
| Base Stack | Caddy, Homer, Go | ~512MB |
| + Phase 1 | Uptime Kuma | +256MB |
| + Phase 2 | WireGuard, Authelia, Redis | +512MB |
| + Phase 3 | Prometheus, Grafana, Exporters | +1088MB |
| + Phase 4 | Kopia, Watchtower, Dockge | +896MB |
| **Total** | **All Services** | **~3.3GB** |

---

## 🛠️ Common Commands

```bash
# View all running services
docker ps

# View resource usage
docker stats

# View logs
docker logs <service_name>

# Restart a service
docker compose -f docker-compose.<stack>.yml restart <service>

# Update a stack
docker compose -f docker-compose.<stack>.yml pull
docker compose -f docker-compose.<stack>.yml up -d

# Backup configurations
make backup

# View backup history
ls -lh /backups/

# Health check
make health
```

---

## 📚 Documentation

- **Full Implementation Plan**: [docs/IMPLEMENTATION_PLAN.md](docs/IMPLEMENTATION_PLAN.md)
- **Architecture**: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
- **Troubleshooting**: See IMPLEMENTATION_PLAN.md → Troubleshooting section

---

## 🆘 Quick Troubleshooting

### Service won't start
```bash
docker logs <service_name>
docker compose -f docker-compose.<stack>.yml restart <service>
```

### Can't access service
```bash
# Check if running
docker ps | grep <service>

# Check Caddy logs
docker logs caddy | grep <service>
```

### Out of memory
```bash
# View usage
docker stats

# Stop non-essential stacks
make stack-stop
```

---

## 🎯 Next Steps

1. **Deploy Phase 1** (Monitoring)
   ```bash
   make stack-monitoring
   ```

2. **Configure Uptime Kuma**
   - Access: https://localhost/status
   - Create admin account
   - Add monitors for your services

3. **Set up first backup**
   ```bash
   make backup
   ```

4. **Review full plan**
   - Read: [docs/IMPLEMENTATION_PLAN.md](docs/IMPLEMENTATION_PLAN.md)
   - Follow phased approach

5. **Join the community**
   - [r/selfhosted](https://reddit.com/r/selfhosted)
   - [r/homelab](https://reddit.com/r/homelab)

---

## 🌟 Why These Tools?

All tools selected based on 2026 r/selfhosted community consensus:

- **Uptime Kuma** - Dominates status monitoring discussions
- **WireGuard** - Replaced OpenVPN as VPN standard
- **Authelia** - Lightweight SSO favorite (vs heavyweight Authentik)
- **Prometheus + Grafana** - Industry standard monitoring
- **Kopia** - Modern backup solution (encrypted, deduplicated)
- **Watchtower** - Community standard for auto-updates
- **Dockge** - Trending compose manager (by Uptime Kuma creator)

All are:
- ✅ 100% FOSS
- ✅ Actively maintained
- ✅ Well-documented
- ✅ Community-vetted
- ✅ Production-ready

---

**Questions?** See [docs/IMPLEMENTATION_PLAN.md](docs/IMPLEMENTATION_PLAN.md) for detailed guidance.
