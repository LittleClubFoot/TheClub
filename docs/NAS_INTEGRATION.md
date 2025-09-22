# Synology NAS Integration Guide (Planned Enhancement)

This document outlines the planned integration of The Club dashboard with Synology NAS services.

## Overview

The Club is designed to work as a lightweight dashboard that will integrate with services running on your Synology NAS, providing unified access to media streaming, file management, gaming, and system administration.

**Status**: This integration is currently planned for future implementation. The configuration files contain commented-out sections ready for when these services are set up.

## Planned Synology NAS Services

### Media Streaming - Jellyfin
- **Service**: Jellyfin Media Server
- **Port**: 8096
- **Dashboard Route**: `/media`
- **Description**: Stream movies, TV shows, music, and other media content
- **Installation**: Install via Synology Package Center or Docker
- **Features**: Transcoding, remote access, mobile apps, subtitle support

### File Management - Copyparty
- **Service**: Copyparty File Browser
- **Port**: 8000
- **Dashboard Route**: `/drive`
- **Description**: Lightweight web-based file browser with upload/download capabilities
- **Installation**: Install via Docker on Synology NAS
- **Features**: File sharing, upload progress, thumbnail previews, mobile-friendly

### Gaming - RetroArch Web
- **Service**: RetroArch Web Server
- **Port**: 8010
- **Dashboard Route**: `/arcade`
- **Description**: Retro game emulation accessible through web browser
- **Installation**: Install via Docker on Synology NAS
- **Features**: Multiple console emulation, save states, controller support

### System Administration
- **Synology DSM**: Port 5000 (HTTP) / 5001 (HTTPS) - NAS management interface
- **Docker Management**: Portainer on port 9000 (optional)
- **Monitoring**: Grafana on port 3000 (optional)

## Configuration Steps

### 1. Enable Caddyfile Routes (When Ready)

The Caddyfile contains commented-out routes ready for your Synology NAS services:

```caddyfile
localhost {
    # Go/HTMX Test Server (Active)
    handle /test* {
        reverse_proxy go-test-server:8080
    }
    
    # FUTURE: Synology NAS Services (Uncomment when ready)
    # handle /media* {
    #     reverse_proxy synology.local:8096  # Jellyfin Media Server
    # }
    # handle /drive* {
    #     reverse_proxy synology.local:8000  # Copyparty File Browser
    # }
    # handle /arcade* {
    #     reverse_proxy synology.local:8010  # RetroArch Web Server
    # }
    
    # Homer Dashboard (default route)
    handle {
        reverse_proxy homer:8080
    }
}
```

**To activate**: Uncomment the desired routes and replace `synology.local` with your actual NAS IP address or hostname.

### 2. Enable Homer Dashboard Links (When Ready)

The Homer dashboard contains commented-out service links ready for activation:

```yaml
services:
  - name: "Synology NAS Services (Planned)"
    icon: "fas fa-server"
    items:
      # TODO: Uncomment when Synology NAS services are configured
      # - name: "Jellyfin Media"
      #   logo: "assets/icons/jellyfin.png"
      #   subtitle: "Movies, TV shows, and music streaming"
      #   tag: "media"
      #   url: "/media"
      #   target: "_self"
      # - name: "File Browser"
      #   logo: "assets/icons/files.png"
      #   subtitle: "Copyparty file management and sharing"
      #   tag: "files"
      #   url: "/drive"
      #   target: "_self"
      # - name: "RetroArch Arcade"
      #   logo: "assets/icons/games.png"
      #   subtitle: "Retro game emulation and arcade"
      #   tag: "games"
      #   url: "/arcade"
      #   target: "_self"
```

**To activate**: Uncomment the desired service entries in `assets/config.yml`.

### 3. Network Configuration

Ensure your NAS and dashboard server can communicate:

1. **Same Network**: Both should be on the same local network
2. **DNS Resolution**: Use IP addresses or configure local DNS
3. **Firewall**: Ensure required ports are open
4. **SSL**: Configure SSL certificates if needed

### 4. Security Considerations

- **Authentication**: Most NAS services have their own auth
- **Network Isolation**: Consider VLANs for security
- **Reverse Proxy**: Use Caddy for SSL termination
- **Access Control**: Configure service-specific permissions

## Synology NAS Service Installation

### Installing Jellyfin on Synology NAS

**Option 1: Package Center (Recommended)**
1. Open Synology Package Center
2. Search for "Jellyfin"
3. Install the official Jellyfin package
4. Configure media libraries through DSM

**Option 2: Docker Container**
1. Install Docker from Package Center
2. Pull Jellyfin image: `jellyfin/jellyfin:latest`
3. Configure volumes for media and config
4. Set port to 8096

### Installing Copyparty on Synology NAS

**Docker Installation (Recommended)**
1. Open Docker in DSM
2. Search for `copyparty/copyparty` image
3. Create container with port 8000
4. Mount volumes for file access
5. Configure authentication if needed

### Installing RetroArch Web on Synology NAS

**Docker Installation**
1. Use image: `inglebard/retroarch-web`
2. Set port to 8010
3. Mount ROM directories
4. Configure controller settings
5. Enable web interface

### Service Port Summary
```yaml
# Synology NAS Services
- DSM (Management): 5000 (HTTP) / 5001 (HTTPS)
- Jellyfin Media: 8096
- Copyparty Files: 8000
- RetroArch Gaming: 8010

# Optional Services
- Portainer (Docker): 9000
- Grafana (Monitoring): 3000
```

## Testing Integration

1. **Start The Club**: `make docker-dev`
2. **Access Dashboard**: `http://localhost`
3. **Test Links**: Click service links in Homer
4. **Check Routing**: Verify Caddy proxying works
5. **Monitor Logs**: Check Caddy logs for issues

## Troubleshooting

### Common Issues

1. **Service Unreachable**
   - Check NAS service is running
   - Verify network connectivity
   - Check firewall settings

2. **SSL Certificate Errors**
   - Configure proper certificates
   - Use HTTP for local testing
   - Check Caddy SSL configuration

3. **Authentication Problems**
   - Each service handles its own auth
   - Configure service-specific settings
   - Consider SSO solutions

### Debugging Commands

```bash
# Test connectivity
ping nas.local

# Check port accessibility  
telnet nas.local 8096

# View Caddy logs
docker logs caddy-container

# Test service directly
curl http://nas.local:8096
```

## Advanced Configuration

### SSL with Let's Encrypt

Update Caddyfile for automatic SSL:

```caddyfile
yourdomain.com {
    # Automatic HTTPS with Let's Encrypt
    handle / {
        reverse_proxy localhost:8080
    }
    
    handle /jellyfin* {
        reverse_proxy nas.local:8096
    }
}
```

### Authentication Proxy

Use Caddy for centralized auth:

```caddyfile
localhost {
    # Protect admin routes
    handle /admin* {
        basicauth {
            admin $2a$14$...  # bcrypt hash
        }
        reverse_proxy nas.local:9000
    }
}
```

This integration approach keeps The Club lightweight while leveraging your existing NAS infrastructure.