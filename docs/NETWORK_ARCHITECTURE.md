# The Club - Synology NAS Network Architecture

## Overview

This document describes the production network architecture for The Club home server running alongside a Synology NAS at **pochita.synology.me**.

---

## Network Architecture

### High-Level Topology

```
Internet
    ↓
[Synology DDNS: pochita.synology.me]
    ↓
Home Router (Port Forwarding)
    ↓
    ├─→ Synology NAS (192.168.x.x)
    │   ├─→ Jellyfin (8096)
    │   ├─→ DSM (5000/5001)
    │   └─→ Other Synology Apps
    │
    └─→ The Club Server (192.168.x.x)
        └─→ Caddy (80/443) → All Club Services
```

---

## Subdomain Structure

The Club uses **subdomain-based routing** under your Synology DDNS domain. Each service is accessible via its own subdomain.

### Subdomain Mapping

| Subdomain | Service | Port (Internal) | Purpose |
|-----------|---------|-----------------|---------|
| **pochita.synology.me** | Homer Dashboard | 8080 | Main landing page |
| **status.pochita.synology.me** | Uptime Kuma | 3001 | Service monitoring |
| **grafana.pochita.synology.me** | Grafana | 3000 | Metrics dashboards |
| **auth.pochita.synology.me** | Authelia | 9091 | SSO authentication |
| **dockge.pochita.synology.me** | Dockge | 5001 | Docker management |
| **backup.pochita.synology.me** | Kopia | 51515 | Backup management |
| **media.pochita.synology.me** | Jellyfin (NAS) | 8096 | Media streaming |
| **nas.pochita.synology.me** | Synology DSM | 5001 | NAS admin interface |

### Optional Additional Services

| Subdomain | Service | Port | Purpose |
|-----------|---------|------|---------|
| **photos.pochita.synology.me** | Immich/PhotoPrism | 2342 | Photo gallery |
| **files.pochita.synology.me** | Copyparty | 8000 | File browser |
| **arcade.pochita.synology.me** | RetroArch | 8010 | Retro gaming |
| **ai.pochita.synology.me** | Ollama/Open WebUI | 3000 | Local AI/LLM |
| **prom.pochita.synology.me** | Prometheus | 9090 | Metrics (internal) |

---

## DNS Configuration

### Synology DDNS Setup

Your Synology NAS already provides:
- **Primary domain**: pochita.synology.me
- **Automatic IP updates**: Synology DDNS service

### Wildcard DNS (Required for Subdomains)

To enable subdomain routing, you need **wildcard DNS** pointing to your home IP:

```
pochita.synology.me           → Your Home IP (handled by Synology)
*.pochita.synology.me         → Your Home IP (wildcard - needs manual setup)
```

#### Option 1: Synology DDNS Service (If Supported)
Check if Synology DDNS supports wildcard records. Some providers do, some don't.

#### Option 2: Custom DNS Provider
If Synology DDNS doesn't support wildcards, use:
- **Cloudflare DNS** (free, supports wildcards)
- **DuckDNS** (free, supports wildcards)
- **No-IP** (free tier available)

**Cloudflare Setup Example:**
1. Add `pochita.synology.me` as CNAME to your Cloudflare account
2. Add wildcard record: `*.pochita.synology.me` → CNAME → `pochita.synology.me`
3. Update your router/Synology to update Cloudflare IP (via API)

#### Option 3: Local DNS Override (Development Only)
For local testing without public DNS:
```bash
# Add to /etc/hosts on your client machines
192.168.1.100  pochita.synology.me
192.168.1.100  status.pochita.synology.me
192.168.1.100  grafana.pochita.synology.me
# ... etc
```

---

## Network Configuration

### Port Forwarding (Router Configuration)

#### Required Ports

Forward these ports from your router to **The Club server**:

| Protocol | External Port | Internal Port | Destination | Service |
|----------|---------------|---------------|-------------|---------|
| TCP | 80 | 80 | Club Server | Caddy HTTP |
| TCP | 443 | 443 | Club Server | Caddy HTTPS |
| UDP | 51820 | 51820 | Club Server | WireGuard VPN |

#### Existing Synology Ports (Keep As-Is)

Your Synology should already have:
- Port 80/443 forwarded for DSM/Jellyfin (may need to change - see below)
- Synology-specific ports

**IMPORTANT: Port 80/443 Conflict Resolution**

You have two options:

**Option A: The Club on Standard Ports (Recommended)**
```
Router Port 80/443 → Club Server (Caddy handles all routing)
Caddy → Proxies to Synology services internally
```
- Pros: Clean subdomain architecture, all HTTPS handled by Caddy
- Cons: Requires changing Synology port forwards

**Option B: Synology on Standard Ports (Keep Current Setup)**
```
Router Port 80/443 → Synology NAS
Router Port 8080/8443 → Club Server
```
- Pros: No changes to existing Synology setup
- Cons: Club services accessed via ports (e.g., pochita.synology.me:8443)

**Recommendation**: Use Option A for clean URLs. Configure in next section.

---

## Physical Deployment Options

### Option 1: Separate Machine (Recommended)
- Run The Club on a separate device (Raspberry Pi, Mini PC, old laptop)
- Dedicated resources for Club services
- Clean separation from NAS
- Easier troubleshooting

### Option 2: Synology Docker
- Run The Club containers directly on Synology NAS
- Single device, less hardware
- May impact NAS performance
- Requires sufficient RAM (4GB+ free)

### Option 3: VM on Synology
- Run The Club in a VM on Synology (if Virtual Machine Manager installed)
- Better isolation than pure Docker
- Moderate resource overhead
- Good middle ground

**This guide assumes Option 1 (separate machine).**

---

## SSL/TLS Certificate Strategy

### Let's Encrypt via Caddy (Automatic)

Caddy will automatically obtain SSL certificates for all subdomains:

```
✓ pochita.synology.me
✓ status.pochita.synology.me
✓ grafana.pochita.synology.me
✓ auth.pochita.synology.me
... etc
```

**Requirements:**
1. Wildcard DNS properly configured
2. Port 80 and 443 forwarded to Caddy
3. Subdomains publicly resolvable

**How it works:**
- Caddy uses HTTP-01 challenge for each subdomain
- Certificates auto-renewed every 60 days
- No manual intervention needed

### Alternative: Wildcard Certificate

If you prefer a single wildcard cert:

```caddyfile
*.pochita.synology.me {
    tls {
        dns cloudflare {env.CLOUDFLARE_API_TOKEN}
    }

    @status host status.pochita.synology.me
    handle @status {
        reverse_proxy uptime-kuma:3001
    }

    # ... other subdomain handlers
}
```

Requires DNS provider API access (Cloudflare recommended).

---

## Integration with Existing Synology Services

### Jellyfin Integration

Your existing Jellyfin on Synology can be proxied through Caddy:

```caddyfile
media.pochita.synology.me {
    reverse_proxy synology.local:8096
}
```

**Where `synology.local` is:**
- Your Synology's local IP (e.g., 192.168.1.10)
- Or hostname if resolvable on local network

### Synology DSM Access

```caddyfile
nas.pochita.synology.me {
    reverse_proxy https://synology.local:5001 {
        transport http {
            tls_insecure_skip_verify  # Only if using self-signed cert
        }
    }
}
```

### Docker Network Considerations

**If Club and NAS on separate machines:**
- Use IP addresses for reverse proxy targets
- Ensure machines can communicate on local network
- May need to adjust firewall rules

**If Club runs on Synology (Docker):**
- Use Docker network names
- Create shared network for inter-container communication

---

## Security Considerations

### Exposing Services to Internet

**Public-facing** (no auth required):
- Homer dashboard (main landing page)
- Uptime Kuma status page (optional - can protect with Authelia)

**Protected with Authelia** (require login):
- Grafana (metrics)
- Dockge (management)
- Kopia (backups)
- Jellyfin (optional - has own auth)
- DSM admin (highly recommended)

### Authelia Protection Example

```caddyfile
grafana.pochita.synology.me {
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.pochita.synology.me
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }
    reverse_proxy grafana:3000
}
```

See [docs/NAS_INTEGRATION.md](NAS_INTEGRATION.md) for full Authelia setup.

---

## Deployment Architecture Diagram

```
                    ┌─────────────────────┐
                    │   Internet          │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │  pochita.synology.me│
                    │  (Synology DDNS)    │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │   Home Router       │
                    │  Port Forwarding:   │
                    │  80/443 → Club      │
                    │  51820  → WireGuard │
                    └─────────┬───────────┘
                              │
                 ┌────────────┴────────────┐
                 │                         │
      ┌──────────▼─────────┐    ┌─────────▼────────┐
      │  Club Server       │    │  Synology NAS     │
      │  192.168.1.100     │    │  192.168.1.10     │
      │                    │    │                   │
      │  ┌──────────────┐  │    │  ┌─────────────┐ │
      │  │ Caddy :80/443│◄─┼────┼──│ Jellyfin    │ │
      │  └──┬───────────┘  │    │  │ :8096       │ │
      │     │              │    │  └─────────────┘ │
      │     │              │    │                   │
      │  ┌──▼───────────┐  │    │  ┌─────────────┐ │
      │  │ Uptime Kuma  │  │    │  │ DSM         │ │
      │  │ :3001        │  │    │  │ :5001       │ │
      │  └──────────────┘  │    │  └─────────────┘ │
      │                    │    │                   │
      │  ┌──────────────┐  │    │  ┌─────────────┐ │
      │  │ Grafana      │  │    │  │ Other Apps  │ │
      │  │ :3000        │  │    │  │             │ │
      │  └──────────────┘  │    │  └─────────────┘ │
      │                    │    │                   │
      │  ┌──────────────┐  │    └───────────────────┘
      │  │ Authelia     │  │
      │  │ :9091        │  │
      │  └──────────────┘  │
      │                    │
      │  ┌──────────────┐  │
      │  │ Prometheus   │  │
      │  │ :9090        │  │
      │  └──────────────┘  │
      │                    │
      │  ┌──────────────┐  │
      │  │ WireGuard    │  │
      │  │ :51820       │  │
      │  └──────────────┘  │
      └────────────────────┘
```

---

## Traffic Flow Examples

### External User → Status Page

```
User → status.pochita.synology.me
  → Router (443)
  → Club Server Caddy (443)
  → Uptime Kuma (3001)
  → Response back through chain
```

### External User → Jellyfin (on NAS)

```
User → media.pochita.synology.me
  → Router (443)
  → Club Server Caddy (443)
  → Synology NAS Jellyfin (8096)
  → Response back through chain
```

### VPN User → Internal Access

```
User → WireGuard VPN
  → Encrypted tunnel (51820)
  → Club Server WireGuard
  → Full access to internal network (192.168.1.x)
  → Can access services directly without port forwarding
```

---

## Next Steps

1. **Review DNS requirements** - Ensure wildcard DNS is possible
2. **Plan port forwarding changes** - Decide between Option A or B
3. **Review security** - Determine which services need Authelia protection
4. **Follow deployment guide** - See [SYNOLOGY_DEPLOYMENT.md](SYNOLOGY_DEPLOYMENT.md)

---

## Related Documentation

- [SYNOLOGY_DEPLOYMENT.md](SYNOLOGY_DEPLOYMENT.md) - Step-by-step deployment guide
- [NAS_INTEGRATION.md](NAS_INTEGRATION.md) - Integrating with Synology services
- [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) - Phased enhancement deployment
- [Caddyfile](../Caddyfile) - Main reverse proxy configuration

---

## Troubleshooting

### Can't Access Subdomains

**Check DNS:**
```bash
# Test DNS resolution
nslookup status.pochita.synology.me

# Should return your public IP
```

**Check Port Forwarding:**
```bash
# From outside your network
curl -I https://status.pochita.synology.me

# Should get HTTP response, not timeout
```

### SSL Certificate Errors

**Check Let's Encrypt logs:**
```bash
docker logs caddy | grep -i acme
```

**Common issues:**
- Port 80 not forwarded (required for HTTP-01 challenge)
- DNS not resolving publicly
- Rate limit hit (5 certs per 7 days per domain)

### Can't Reach Synology Services Through Caddy

**Check connectivity from Club server:**
```bash
# From Club server
docker exec caddy wget -O- http://192.168.1.10:8096

# Should connect to Jellyfin
```

**Check firewall:**
- Synology firewall may block Club server
- Add allow rule for Club server IP

---

## FAQ

**Q: Can I use this without exposing to the internet?**
A: Yes! Use WireGuard VPN for remote access, no port forwarding needed except UDP 51820.

**Q: Do I need to move Jellyfin off the Synology?**
A: No! Caddy will proxy to it. Jellyfin stays on Synology.

**Q: What if I can't get wildcard DNS?**
A: Use path-based routing instead (pochita.synology.me/status, /grafana, etc). See Caddyfile examples.

**Q: Can Club run on the Synology itself?**
A: Yes, via Docker on Synology. Requires 4GB+ free RAM and Docker installed.

**Q: How do I secure DSM access through Caddy?**
A: Add Authelia forward_auth to nas.pochita.synology.me. See security section above.
