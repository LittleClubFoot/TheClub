# The Club - Home Lab Dashboard

> A lightweight, containerized home lab dashboard designed to integrate seamlessly with your existing NAS infrastructure.

[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 🎯 **Overview**

The Club is a modern home lab dashboard that acts as a central hub for your Synology NAS services. It provides unified access to media streaming, file management, gaming, and system administration through a clean, responsive web interface.

### ✨ **Key Features**

- 🏠 **Homer Dashboard**: Beautiful, configurable service homepage
- 🧪 **Go/HTMX Server**: Custom applications with dynamic content
- 🔒 **Caddy Reverse Proxy**: Automatic SSL, intelligent routing
- 🖥️ **Synology NAS Integration**: Jellyfin, Copyparty, RetroArch integration
- 🐳 **Docker Compose**: One-command deployment
- 📱 **Responsive Design**: Works on desktop, tablet, and mobile

## 🚀 **Quick Start**

### Prerequisites

- Docker and Docker Compose
- 2GB+ RAM
- Network access to your Synology NAS

### 🏃 **Development Setup**

```bash
# Clone the repository
git clone <repository-url>
cd TheClub

# Start all services
make docker-dev
# OR: docker compose -f docker-compose.dev.yml up

# Access the dashboard
open https://localhost
```

### 🏭 **Production Deployment**

```bash
# Configure your domain in Caddyfile
vim Caddyfile

# Deploy with production settings
make docker-prod
# OR: docker compose -f docker-compose.prod.yml up -d
```

## 📁 **Project Structure**

```
The Club/
├── 📂 cmd/webserver/           # Go application entry point
├── 📂 templates/               # HTML templates for Go server
├── 📂 config/                  # Configuration files
│   └── homer.yml              # Homer dashboard config
├── 📂 homer-config/           # Homer assets and config
├── 📂 assets/                 # Static assets (icons, images)
├── 📂 docs/                   # Documentation
│   ├── ARCHITECTURE.md        # System architecture
│   ├── NAS_INTEGRATION.md     # NAS setup guide
│   └── RUNBOOK.md            # Operations guide
├── 📂 logs/                   # Application logs
├── 🐳 Dockerfile             # Multi-stage container build
├── 🐳 docker-compose.dev.yml # Development environment
├── 🐳 docker-compose.prod.yml# Production environment
├── 🔧 Caddyfile              # Reverse proxy configuration
├── 🔧 Makefile               # Build and deployment commands
└── 📋 .air.toml              # Live reload configuration
```

## 🌐 **Service Architecture**

The Club uses a microservices architecture with Caddy as the entry point:

- **Caddy** (Port 80/443): Reverse proxy with automatic SSL
- **Homer** (Port 8080): Static dashboard homepage
- **Go Server** (Port 8082): Custom HTMX applications
- **Synology NAS Services**: Jellyfin (media), Copyparty (files), RetroArch (gaming)

## 🔧 **Configuration**

### Homer Dashboard

The dashboard is prepared for future Synology NAS services integration:

```yaml
services:
  - name: "Synology NAS Services (Planned)"
    items:
      # TODO: Uncomment when services are configured
      # - name: "Jellyfin Media"
      #   url: "/media"          # Jellyfin on port 8096
      # - name: "File Browser"
      #   url: "/drive"          # Copyparty on port 8000
      # - name: "RetroArch Arcade"
      #   url: "/arcade"         # RetroArch on port 8010
```

**To activate**: Uncomment the desired services in `assets/config.yml` when your NAS services are ready.

### Caddy Reverse Proxy

The Caddyfile contains commented-out routes ready for your Synology NAS services:

```caddyfile
localhost {
    # FUTURE: Uncomment when NAS services are ready
    # handle /media* {
    #     reverse_proxy synology.local:8096  # Jellyfin
    # }
    # handle /drive* {
    #     reverse_proxy synology.local:8000  # Copyparty
    # }
    # handle /arcade* {
    #     reverse_proxy synology.local:8010  # RetroArch
    # }
}
```

**To activate**: Uncomment the desired routes and replace `synology.local` with your actual NAS IP address or hostname.

See [Synology NAS Integration Guide](docs/NAS_INTEGRATION.md) for detailed setup instructions.

## 📖 **Documentation**

- 🏗️ [Architecture Overview](docs/ARCHITECTURE.md)
- 🖥️ [Synology NAS Integration Guide](docs/NAS_INTEGRATION.md)
- 📋 [Operations Runbook](docs/RUNBOOK.md)

## 🛠️ **Available Commands**

```bash
# Development
make dev          # Start with live reload
make docker-dev   # Start full stack with Docker

# Production  
make docker-prod  # Deploy production stack
make docker-stop  # Stop all services

# Maintenance
make build        # Build Go application
make test         # Run tests
make clean        # Clean build artifacts
make fmt          # Format code
```

## 🔍 **Service URLs**

| Service | Development | Production |
|---------|-------------|------------|
| Dashboard | https://localhost | https://yourdomain.com |
| Go Server | https://localhost/test | https://yourdomain.com/test |
| Health Check | https://localhost/health | https://yourdomain.com/health |

## 🤝 **Contributing**

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run `make fmt` and `make test`
5. Submit a pull request

## 📄 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 **Support**

- 📖 Check the [documentation](docs/)
- 🐛 Report issues on GitHub
- 💬 Join the discussion in Issues

---

**Made with ❤️ for the self-hosted community**