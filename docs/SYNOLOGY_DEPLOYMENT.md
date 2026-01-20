# The Club - Synology NAS Deployment Guide

## Overview

This guide provides step-by-step instructions for deploying The Club alongside your existing Synology NAS at **pochita.synology.me** with subdomain-based routing.

---

## Prerequisites

Before starting, ensure you have:

### Hardware
- [ ] Synology NAS running at pochita.synology.me (already configured)
- [ ] Separate machine for The Club (Raspberry Pi, Mini PC, or old laptop)
  - **Minimum**: 2GB RAM, 2 CPU cores, 20GB storage
  - **Recommended**: 4GB+ RAM, 4 CPU cores, 40GB storage
- [ ] Both devices on same local network

### Network Access
- [ ] Access to router admin (for port forwarding)
- [ ] Ability to configure DNS (Cloudflare, DuckDNS, or similar)
- [ ] Your public IP address

### Software
- [ ] Docker and Docker Compose installed on Club server
- [ ] SSH access to Club server
- [ ] Git installed

### Knowledge
- [ ] Basic Linux command line
- [ ] Understanding of Docker
- [ ] Basic networking (ports, DNS, port forwarding)

---

## Phase 0: DNS and Network Preparation

### Step 1: Set Up Wildcard DNS

Your Synology DDNS provides `pochita.synology.me`, but you need wildcard support for subdomains.

#### Option A: Cloudflare DNS (Recommended)

**Why Cloudflare?**
- Free tier includes wildcard DNS
- API for automated updates
- Built-in DDoS protection
- Fast propagation

**Setup:**

1. **Create Cloudflare Account**
   - Visit https://cloudflare.com
   - Sign up for free account

2. **Add Domain Delegation** (if using custom domain)
   - Or skip if using Cloudflare for DNS only

3. **Add DNS Records**
   ```
   Type: A
   Name: pochita.synology.me
   Content: [Your Public IP]
   Proxy: Off (DNS only)
   TTL: Auto

   Type: A
   Name: *.pochita.synology.me
   Content: [Your Public IP]
   Proxy: Off (DNS only)
   TTL: Auto
   ```

4. **Get API Token**
   - My Profile → API Tokens → Create Token
   - Use "Edit zone DNS" template
   - Save token securely (needed for Caddy)

5. **Update IP Automatically** (optional)
   - Use ddclient or Cloudflare's Dynamic DNS
   - Or let Synology DDNS handle it with CNAME

#### Option B: DuckDNS (Simpler Alternative)

1. Visit https://www.duckdns.org
2. Create account (via GitHub, Google, etc.)
3. Create domain: `pochita` (becomes pochita.duckdns.org)
4. Note your token
5. Update Synology to use DuckDNS (via DDNS settings)

**Note:** If keeping Synology DDNS, you may need to manually point *.pochita.synology.me to your IP via Synology's DNS provider if they support it.

### Step 2: Configure Router Port Forwarding

**Important:** This will change how you access Synology services. Back up your current configuration first!

#### Identify Current Port Forwards

Check your router for existing forwards:
```
Port 80   → Synology (typically)
Port 443  → Synology (typically)
Port 5000 → Synology DSM (optional)
Port 5001 → Synology DSM HTTPS (optional)
```

#### New Port Forwarding Configuration

**Option A: Club Takes 80/443 (Recommended for subdomain architecture)**

| External Port | Protocol | Internal IP | Internal Port | Service |
|---------------|----------|-------------|---------------|---------|
| 80 | TCP | [Club Server IP] | 80 | Caddy HTTP |
| 443 | TCP | [Club Server IP] | 443 | Caddy HTTPS |
| 51820 | UDP | [Club Server IP] | 51820 | WireGuard VPN |
| 8096 | TCP | [Synology IP] | 8096 | Jellyfin (temp, until proxied) |
| 5001 | TCP | [Synology IP] | 5001 | DSM HTTPS (temp, until proxied) |

**Result:** All traffic goes through Caddy, which proxies to Synology services internally.

**Option B: Keep Synology on 80/443 (Minimal changes)**

| External Port | Protocol | Internal IP | Internal Port | Service |
|---------------|----------|-------------|---------------|---------|
| 80 | TCP | [Synology IP] | 80 | Synology HTTP |
| 443 | TCP | [Synology IP] | 443 | Synology HTTPS |
| 8080 | TCP | [Club Server IP] | 80 | Caddy HTTP |
| 8443 | TCP | [Club Server IP] | 443 | Caddy HTTPS |
| 51820 | UDP | [Club Server IP] | 51820 | WireGuard VPN |

**Result:** Access Club services via port (e.g., pochita.synology.me:8443/status)

**This guide assumes Option A for clean subdomain URLs.**

### Step 3: Verify Network Setup

Test from outside your network (use mobile data or VPN):

```bash
# Test DNS resolution
nslookup pochita.synology.me
nslookup status.pochita.synology.me

# Both should return your public IP

# Test ports (will fail until Club is deployed, but shouldn't timeout)
nc -zv pochita.synology.me 443
```

---

## Phase 1: Deploy The Club Base Stack

### Step 1: Clone Repository on Club Server

```bash
# SSH into your Club server
ssh user@[club-server-ip]

# Clone the repository
cd ~
git clone https://github.com/LittleClubFoot/TheClub.git
cd TheClub

# Checkout the enhancement branch
git checkout claude/review-home-server-G6FQx
```

### Step 2: Configure Domain in Caddyfile

```bash
# Edit the Caddyfile
nano Caddyfile
```

**Change line 15** from:
```caddyfile
localhost {
```

**To:**
```caddyfile
pochita.synology.me {
```

**Also update subdomains** (around line 30-50):

```caddyfile
pochita.synology.me {
    # Main homepage - Homer dashboard
    handle {
        reverse_proxy homer:8080
    }
}

status.pochita.synology.me {
    reverse_proxy uptime-kuma:3001
}

grafana.pochita.synology.me {
    reverse_proxy grafana:3000
}

auth.pochita.synology.me {
    reverse_proxy authelia:9091
}

dockge.pochita.synology.me {
    reverse_proxy dockge:5001
}

backup.pochita.synology.me {
    reverse_proxy kopia:51515
}

# Jellyfin on Synology NAS
media.pochita.synology.me {
    reverse_proxy [SYNOLOGY_IP]:8096
}

# Synology DSM
nas.pochita.synology.me {
    reverse_proxy https://[SYNOLOGY_IP]:5001 {
        transport http {
            tls_insecure_skip_verify
        }
    }
}
```

**Replace `[SYNOLOGY_IP]`** with your Synology's local IP (e.g., 192.168.1.10).

Save and exit (Ctrl+O, Enter, Ctrl+X).

### Step 3: Update Homer Dashboard Config

```bash
# Edit Homer config
nano assets/config.yml
```

Update the title and subtitle:

```yaml
title: "The Club"
subtitle: "pochita.synology.me"
```

Update service URLs to use subdomains:

```yaml
services:
  - name: "Monitoring"
    icon: "fas fa-heartbeat"
    items:
      - name: "Status Dashboard"
        subtitle: "Uptime monitoring"
        url: "https://status.pochita.synology.me"
        icon: "fas fa-chart-line"

      - name: "Metrics"
        subtitle: "Grafana dashboards"
        url: "https://grafana.pochita.synology.me"
        icon: "fas fa-chart-area"

  - name: "Management"
    icon: "fas fa-tools"
    items:
      - name: "Docker Stacks"
        subtitle: "Dockge management"
        url: "https://dockge.pochita.synology.me"
        icon: "fas fa-cubes"

      - name: "Backups"
        subtitle: "Kopia backup system"
        url: "https://backup.pochita.synology.me"
        icon: "fas fa-database"

  - name: "Media"
    icon: "fas fa-film"
    items:
      - name: "Jellyfin"
        subtitle: "Media server"
        url: "https://media.pochita.synology.me"
        icon: "fas fa-play-circle"

      - name: "Synology NAS"
        subtitle: "NAS management"
        url: "https://nas.pochita.synology.me"
        icon: "fas fa-server"
```

Save and exit.

### Step 4: Deploy Base Stack

```bash
# Run initial backup
make backup

# Deploy production stack
make docker-prod

# Wait for containers to start (30 seconds)
sleep 30

# Check status
docker ps
```

You should see:
- caddy
- homer
- go-test-server

### Step 5: Test Base Stack

From any device on your network:

```bash
# Test Homer dashboard
curl -k https://pochita.synology.me

# Check Caddy logs
docker logs caddy | tail -20

# Check SSL certificate status
docker exec caddy caddy list-certificates
```

**From outside your network** (mobile data):

Visit https://pochita.synology.me

You should see:
- ✅ Valid SSL certificate (Let's Encrypt)
- ✅ Homer dashboard loads
- ✅ No certificate warnings

**If you see certificate errors:**
- Wait 2-3 minutes (Caddy obtaining certs)
- Check Caddy logs: `docker logs caddy | grep acme`
- Ensure port 80 is forwarded (required for HTTP-01 challenge)

---

## Phase 2: Deploy Enhancement Stacks

Now that the base is working, add monitoring, security, and automation.

### Step 1: Deploy Monitoring (Uptime Kuma)

```bash
cd ~/TheClub

# Deploy monitoring stack
make stack-monitoring

# Wait for startup
sleep 30

# Check status
docker ps | grep uptime-kuma
```

**Configure Uptime Kuma:**

1. Visit https://status.pochita.synology.me
2. Create admin account on first visit
3. Add monitors:
   - **HTTP**: https://pochita.synology.me (Homer)
   - **HTTP**: https://media.pochita.synology.me (Jellyfin)
   - **HTTP**: https://nas.pochita.synology.me (DSM)
   - **Docker**: caddy
   - **Docker**: homer
   - **Docker**: uptime-kuma

4. Set up notifications:
   - Settings → Notifications
   - Add email/Discord/Telegram
   - Test notification

### Step 2: Deploy Security Stack (WireGuard + Authelia)

**Before deploying, configure Authelia:**

```bash
# Generate secure secrets
docker run --rm authelia/authelia:latest authelia crypto rand --length 64

# Copy output, you'll need two of these
```

Edit Authelia config:

```bash
nano authelia/configuration.yml
```

Update:
```yaml
session:
  domain: pochita.synology.me  # Change from yourdomain.com
  secret: [PASTE SECRET 1 HERE]

storage:
  encryption_key: [PASTE SECRET 2 HERE - min 32 chars]
```

Create user password:

```bash
# Generate password hash
docker run --rm authelia/authelia:latest \
  authelia crypto hash generate argon2 --password 'YourSecurePassword123!'

# Copy the hash
```

Edit users database:

```bash
nano authelia/users_database.yml
```

Update:
```yaml
users:
  admin:
    displayname: "Your Name"
    password: "[PASTE HASH HERE]"
    email: "your.email@example.com"
    groups:
      - admins
      - users
```

**Deploy security stack:**

```bash
make stack-security

# Wait for startup
sleep 60

# Check status
docker ps | grep -E 'wireguard|authelia|redis'

# Test Authelia
curl -k https://auth.pochita.synology.me
```

**Configure WireGuard:**

```bash
# Get VPN config for first peer
docker exec wireguard cat /config/peer1/peer1.conf

# Or get QR code for mobile
docker exec wireguard cat /config/peer1/peer1.png | base64

# Copy config to your device
docker cp wireguard:/config/peer1/peer1.conf ~/wireguard-client.conf
```

**Install WireGuard client:**
- Desktop: https://www.wireguard.com/install/
- Mobile: App store (WireGuard app)

Import the config and connect!

### Step 3: Deploy Observability (Prometheus + Grafana)

```bash
make stack-observability

# Wait for startup
sleep 60

# Check status
docker ps | grep -E 'prometheus|grafana|node-exporter|cadvisor'
```

**Configure Grafana:**

1. Visit https://grafana.pochita.synology.me
2. Login: admin / changeme
3. Change password immediately
4. Import dashboards:
   - Dashboards → New → Import → Enter ID `1860` (Node Exporter Full)
   - Select Prometheus datasource → Import
   - Repeat for IDs: `893` (Docker metrics), `14280` (Caddy metrics)

### Step 4: Deploy Automation (Kopia + Watchtower + Dockge)

**Prepare backup destination:**

```bash
# Create backup directory (or use external drive/NAS mount)
sudo mkdir -p /mnt/backup
sudo chown $(whoami):$(whoami) /mnt/backup

# Or mount Synology NAS share
# sudo mount -t nfs [SYNOLOGY_IP]:/volume1/backups /mnt/backup
```

**Configure Kopia password:**

```bash
nano docker-compose.automation.yml
```

Update line ~55:
```yaml
- KOPIA_PASSWORD=YourSecureBackupPassword123!  # Change this!
```

**Deploy automation stack:**

```bash
# Create stacks directory for Dockge
sudo mkdir -p /opt/stacks
sudo chown -R $(whoami):$(whoami) /opt/stacks

make stack-automation

# Wait for startup
sleep 60

# Check status
docker ps | grep -E 'kopia|watchtower|dockge'
```

**Configure Kopia:**

1. Visit https://backup.pochita.synology.me
2. Set up repository at `/repository`
3. Enter password (matches KOPIA_PASSWORD above)
4. Create snapshot policies:
   - Daily at 3 AM
   - Keep: 7 daily, 4 weekly, 6 monthly
5. Create snapshots for:
   - `/data/theclub` (configs)
   - `/data/caddy_data` (SSL certs)
   - `/data/uptime-kuma-data` (monitoring)
   - `/data/grafana-data` (dashboards)

**Access Dockge:**

1. Visit https://dockge.pochita.synology.me
2. Create admin account
3. View all your compose stacks in web UI

---

## Phase 3: Protect Services with Authelia

Now that everything is running, add SSO protection to sensitive services.

### Step 1: Update Authelia Access Control

```bash
nano authelia/configuration.yml
```

Add under `access_control.rules`:

```yaml
access_control:
  default_policy: deny

  rules:
    # Public services (no auth)
    - domain: "pochita.synology.me"
      policy: bypass

    - domain: "status.pochita.synology.me"
      policy: bypass  # Or 'one_factor' if you want to protect it

    # Protected services (require login)
    - domain:
        - "grafana.pochita.synology.me"
        - "dockge.pochita.synology.me"
        - "backup.pochita.synology.me"
        - "nas.pochita.synology.me"
      policy: two_factor  # Requires 2FA

    # Media can have its own auth
    - domain: "media.pochita.synology.me"
      policy: bypass  # Jellyfin has own auth
```

### Step 2: Update Caddyfile with Authelia Forward Auth

```bash
nano Caddyfile
```

Add before the protected service blocks:

```caddyfile
# Grafana with Authelia protection
grafana.pochita.synology.me {
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.pochita.synology.me
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }
    reverse_proxy grafana:3000
}

# Dockge with Authelia protection
dockge.pochita.synology.me {
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.pochita.synology.me
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }
    reverse_proxy dockge:5001
}

# Kopia with Authelia protection
backup.pochita.synology.me {
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.pochita.synology.me
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }
    reverse_proxy kopia:51515
}

# Synology DSM with Authelia protection (highly recommended!)
nas.pochita.synology.me {
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.pochita.synology.me
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }
    reverse_proxy https://[SYNOLOGY_IP]:5001 {
        transport http {
            tls_insecure_skip_verify
        }
    }
}
```

### Step 3: Restart Services

```bash
# Restart Caddy to apply changes
docker restart caddy

# Restart Authelia to apply access rules
docker compose -f docker-compose.security.yml restart authelia

# Test protected service
curl -I https://grafana.pochita.synology.me
# Should get 302 redirect to auth.pochita.synology.me
```

### Step 4: Enable 2FA in Authelia

1. Visit https://auth.pochita.synology.me
2. Login with your credentials
3. Click your profile
4. Set up two-factor authentication
5. Scan QR code with authenticator app (Google Authenticator, Authy, etc.)
6. Enter code to verify

Now accessing protected services requires:
1. Username + password
2. TOTP code from authenticator app

---

## Phase 4: Verify Complete Setup

### Check All Services

Visit each subdomain and verify:

| URL | Expected Result |
|-----|-----------------|
| https://pochita.synology.me | ✅ Homer dashboard, no auth |
| https://status.pochita.synology.me | ✅ Uptime Kuma, all monitors green |
| https://grafana.pochita.synology.me | 🔒 Redirects to auth, then shows Grafana |
| https://auth.pochita.synology.me | ✅ Authelia login page |
| https://dockge.pochita.synology.me | 🔒 Auth required, shows stacks |
| https://backup.pochita.synology.me | 🔒 Auth required, shows Kopia |
| https://media.pochita.synology.me | ✅ Jellyfin (own auth) |
| https://nas.pochita.synology.me | 🔒 Auth required, shows DSM |

### Check Resource Usage

```bash
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}"
```

All services should be within limits.

### Check Logs

```bash
# No errors in Caddy
docker logs caddy --tail=50 | grep -i error

# Authelia working
docker logs authelia --tail=50

# Prometheus collecting metrics
docker logs prometheus --tail=50
```

### Test Backups

```bash
# Run manual backup
make backup

# Check backup created
ls -lh /backups/theclub-*/

# Test Kopia snapshot
docker exec kopia kopia snapshot list
```

### Test WireGuard VPN

1. Connect via WireGuard
2. Access services via internal IPs (bypass internet)
3. Should be faster (local network speed)

---

## Maintenance Tasks

### Daily (Automated)
- ✅ Watchtower checks for updates (4 AM)
- ✅ Kopia creates backups (3 AM)
- ✅ Uptime Kuma monitors services
- ✅ Prometheus collects metrics

### Weekly
- Check Uptime Kuma for any downtime
- Review Grafana dashboards for anomalies
- Check Watchtower logs for updates

### Monthly
- Test backup restoration
- Review Authelia logs for failed logins
- Update WireGuard peer configs if needed
- Check disk usage

---

## Troubleshooting

### Can't Access Subdomains from Internet

**Check DNS:**
```bash
nslookup status.pochita.synology.me 8.8.8.8
# Should return your public IP
```

**Check port forwarding:**
- Ensure 80/443 forwarded to Club server
- Test from outside network (mobile data)

**Check Caddy:**
```bash
docker logs caddy | grep -i error
docker exec caddy caddy list-certificates
```

### Authelia Not Protecting Services

**Check Authelia logs:**
```bash
docker logs authelia | tail -50
```

**Check Caddyfile syntax:**
```bash
docker exec caddy caddy validate --config /etc/caddy/Caddyfile
```

**Test forward auth directly:**
```bash
curl -I http://localhost:9091/api/verify
# Should return 401 Unauthorized
```

### Jellyfin Not Accessible

**Check connectivity from Club server:**
```bash
docker exec caddy wget -O- http://[SYNOLOGY_IP]:8096
```

**Check Synology firewall:**
- Control Panel → Security → Firewall
- Add allow rule for Club server IP

### SSL Certificate Errors

**Check Let's Encrypt rate limits:**
- 50 certificates per domain per week
- 5 certificates per hostname per week

**Use staging mode temporarily:**
```caddyfile
{
    acme_ca https://acme-staging-v02.api.letsencrypt.org/directory
}
```

**Check DNS propagation:**
```bash
dig status.pochita.synology.me +trace
```

---

## Next Steps

1. **Create pull request** to merge enhancements
2. **Add more services** as needed (Immich, PhotoPrism, etc.)
3. **Configure email notifications** for Watchtower and Uptime Kuma
4. **Set up off-site backups** (Kopia to cloud storage)
5. **Add custom dashboards** in Grafana
6. **Document your specific setup** for future reference

---

## Getting Help

- **Documentation**: See [docs/](.) directory
- **Logs**: `docker logs <service-name>`
- **Community**: r/selfhosted, r/homelab
- **Authelia Docs**: https://www.authelia.com/
- **Caddy Docs**: https://caddyserver.com/docs/

---

## Congratulations! 🎉

Your home server is now fully deployed with:
- ✅ Beautiful dashboard at pochita.synology.me
- ✅ Monitoring with Uptime Kuma
- ✅ Metrics with Prometheus & Grafana
- ✅ Secure access with Authelia & WireGuard
- ✅ Automated backups with Kopia
- ✅ Auto-updates with Watchtower
- ✅ Easy management with Dockge
- ✅ Full integration with your Synology NAS

All accessible via clean subdomain URLs with automatic SSL certificates!
