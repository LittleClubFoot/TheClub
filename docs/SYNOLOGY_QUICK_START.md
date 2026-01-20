# The Club - Synology NAS Quick Start Guide

## Your Setup

- **Domain**: pochita.synology.me (Synology DDNS)
- **NAS**: Synology running Jellyfin and other apps
- **Goal**: Subdomain-based routing for all services

---

## Quick Reference

### Subdomain Map

| URL | Service | Auth Required |
|-----|---------|---------------|
| pochita.synology.me | Homer Dashboard | No |
| status.pochita.synology.me | Uptime Kuma | No |
| grafana.pochita.synology.me | Grafana | Yes (Authelia) |
| auth.pochita.synology.me | Authelia SSO | No |
| dockge.pochita.synology.me | Docker Management | Yes (Authelia) |
| backup.pochita.synology.me | Kopia Backups | Yes (Authelia) |
| media.pochita.synology.me | Jellyfin (NAS) | Own Auth |
| nas.pochita.synology.me | Synology DSM | Yes (Authelia) |

---

## Prerequisites

### DNS Requirements
- [ ] Wildcard DNS configured: `*.pochita.synology.me` → Your public IP
- [ ] Option: Use Cloudflare for free wildcard DNS
- [ ] Option: Use DuckDNS as alternative

### Network Requirements
- [ ] Port 80 forwarded to Club server (HTTP)
- [ ] Port 443 forwarded to Club server (HTTPS)
- [ ] Port 51820 forwarded to Club server (WireGuard VPN)

### Hardware Requirements
- [ ] Separate machine for The Club (recommended)
  - Raspberry Pi 4 (4GB+), Mini PC, or old laptop
  - Or run on Synology directly (requires 4GB+ free RAM)

---

## 5-Minute Setup (After Prerequisites)

```bash
# 1. Clone repository
git clone https://github.com/LittleClubFoot/TheClub.git
cd TheClub

# 2. Configure domain
cp Caddyfile.subdomain.example Caddyfile
nano Caddyfile
# Replace [SYNOLOGY_IP] with your Synology's IP (e.g., 192.168.1.10)
# Update email address in global options

# 3. Generate Authelia secrets
docker run --rm authelia/authelia:latest authelia crypto rand --length 64
# Copy output twice, paste into authelia/configuration.yml

# 4. Create admin password
docker run --rm authelia/authelia:latest \
  authelia crypto hash generate argon2 --password 'YourPassword'
# Copy hash, paste into authelia/users_database.yml

# 5. Deploy base stack
make docker-prod

# 6. Deploy enhancements
make stack-all

# 7. Wait for startup (2-3 minutes for SSL certificates)
sleep 180

# 8. Check status
docker ps
curl -k https://pochita.synology.me
```

---

## Configuration Files to Edit

### 1. Caddyfile (Required)
```bash
cp Caddyfile.subdomain.example Caddyfile
nano Caddyfile
```

**Change:**
- Line 15: Update email address
- Throughout: Replace `[SYNOLOGY_IP]` with your Synology's local IP

### 2. authelia/configuration.yml (Required)
```bash
nano authelia/configuration.yml
```

**Change:**
- Line ~48: `domain: pochita.synology.me`
- Line ~51: `secret: [PASTE GENERATED SECRET]`
- Line ~85: `encryption_key: [PASTE GENERATED SECRET]`

### 3. authelia/users_database.yml (Required)
```bash
nano authelia/users_database.yml
```

**Change:**
- Line ~16: Your display name
- Line ~18: `password: [PASTE GENERATED HASH]`
- Line ~19: Your email

### 4. assets/config.yml (Optional - Homer Dashboard)
```bash
nano assets/config.yml
```

**Change:**
- Update service URLs to use subdomains
- Add your own services/links

---

## Testing Checklist

After deployment, verify:

### From Your Local Network
- [ ] http://[CLUB_SERVER_IP] loads Homer
- [ ] All containers running: `docker ps`
- [ ] No errors in logs: `docker logs caddy | grep -i error`

### From Internet (Mobile Data)
- [ ] https://pochita.synology.me loads with valid SSL
- [ ] https://status.pochita.synology.me shows Uptime Kuma
- [ ] https://grafana.pochita.synology.me redirects to auth
- [ ] https://auth.pochita.synology.me shows login page
- [ ] https://media.pochita.synology.me shows Jellyfin

### Authentication Flow
- [ ] Login at auth.pochita.synology.me works
- [ ] Protected services require auth
- [ ] After auth, services load correctly
- [ ] 2FA setup works in Authelia profile

### Monitoring
- [ ] Uptime Kuma shows all monitors green
- [ ] Grafana dashboards display metrics
- [ ] Prometheus targets all "UP"

---

## Common Issues

### "Can't access subdomains"
**Fix:** Check wildcard DNS is configured
```bash
nslookup status.pochita.synology.me
# Should return your public IP
```

### "SSL certificate errors"
**Fix:** Wait 2-3 minutes for Caddy to obtain certificates
```bash
docker logs caddy | grep acme
# Check for certificate issuance logs
```

### "Authelia not protecting services"
**Fix:** Check forward_auth configuration in Caddyfile
```bash
docker exec caddy caddy validate --config /etc/caddy/Caddyfile
docker logs authelia | tail -50
```

### "Can't reach Jellyfin/Synology services"
**Fix:** Verify Synology IP and connectivity
```bash
docker exec caddy wget -O- http://[SYNOLOGY_IP]:8096
# Should connect to Jellyfin
```

---

## Next Steps

1. **Set up Uptime Kuma monitors**
   - Add monitors for all your services
   - Configure notifications (email/Discord/Telegram)

2. **Import Grafana dashboards**
   - Node Exporter (ID: 1860)
   - Docker metrics (ID: 893)
   - Caddy metrics (ID: 14280)

3. **Configure WireGuard VPN**
   ```bash
   docker exec wireguard cat /config/peer1/peer1.conf
   # Import into WireGuard client
   ```

4. **Set up automated backups**
   - Configure Kopia at backup.pochita.synology.me
   - Create snapshot policies
   - Test restoration

5. **Add more services** (optional)
   - Immich for photos
   - PhotoPrism
   - Ollama for local AI
   - Whatever you need!

---

## Architecture Overview

```
Internet
    ↓
[pochita.synology.me + wildcard DNS]
    ↓
Your Router (Port Forward 80/443)
    ↓
The Club Server (Caddy)
    ├─→ Internal Services (Grafana, Uptime Kuma, etc.)
    └─→ Synology NAS (Jellyfin, DSM, etc.)
```

**Key Point:** All traffic goes through Caddy, which:
- Handles SSL certificates automatically
- Routes to appropriate services based on subdomain
- Enforces Authelia authentication where configured
- Proxies to your Synology for existing services

---

## Documentation Links

- **Full Deployment Guide**: [SYNOLOGY_DEPLOYMENT.md](SYNOLOGY_DEPLOYMENT.md)
- **Network Architecture**: [NETWORK_ARCHITECTURE.md](NETWORK_ARCHITECTURE.md)
- **Enhancement Stacks**: [../ENHANCEMENTS.md](../ENHANCEMENTS.md)
- **Implementation Plan**: [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md)

---

## Support

- Check logs: `docker logs <service-name>`
- Validate Caddy config: `docker exec caddy caddy validate --config /etc/caddy/Caddyfile`
- Test DNS: `nslookup subdomain.pochita.synology.me`
- Test ports: `nc -zv pochita.synology.me 443`

---

## Key Commands

```bash
# Deploy everything
make stack-all

# Stop everything
make stack-stop
docker compose -f docker-compose.prod.yml down

# View logs
docker logs caddy
docker logs authelia
docker logs uptime-kuma

# Restart a service
docker restart caddy

# Check resource usage
docker stats

# Run backup
make backup

# Check certificate status
docker exec caddy caddy list-certificates
```

---

**Ready to deploy?** Follow [SYNOLOGY_DEPLOYMENT.md](SYNOLOGY_DEPLOYMENT.md) for step-by-step instructions!
