# Development Guide

This guide covers development setup, workflows, and best practices for The Club dashboard.

## 🛠️ **Development Environment Setup**

### Prerequisites

- **Docker & Docker Compose**: For containerized development
- **Go 1.25+**: For local Go development
- **Git**: For version control
- **Make**: For build automation (optional)

### Quick Setup

```bash
# Clone the repository
git clone <repository-url>
cd TheClub

# Start development environment
make docker-dev
# OR: docker compose -f docker-compose.dev.yml up

# Access services
open https://localhost        # Main dashboard
open https://localhost/test   # Go/HTMX server
```

## 🏗️ **Project Architecture**

### Service Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Caddy       │    │     Homer       │    │   Go Server     │
│  Reverse Proxy  │────│   Dashboard     │    │  HTMX/APIs      │
│   Port 80/443   │    │   Port 8080     │    │   Port 8082     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │  NAS Services   │
                    │ (External)      │
                    └─────────────────┘
```

### Directory Structure

```
TheClub/
├── cmd/webserver/          # Go application entry point
│   └── main.go            # HTTP server with HTMX routes
├── templates/             # HTML templates for Go server
│   ├── test.html         # Main test page
│   ├── api.html          # API documentation
│   └── docs.html         # General documentation
├── config/               # Configuration files
│   └── homer.yml         # Homer dashboard configuration
├── homer-config/         # Homer runtime configuration
│   └── config.yml        # Copied from config/homer.yml
├── assets/               # Static assets (icons, images)
├── docs/                 # Project documentation
├── logs/                 # Application logs (created at runtime)
├── Caddyfile            # Reverse proxy configuration
├── docker-compose.dev.yml  # Development environment
├── docker-compose.prod.yml # Production environment
├── Dockerfile           # Multi-stage container build
├── Makefile            # Build and deployment automation
└── .air.toml           # Live reload configuration
```

## 🔄 **Development Workflows**

### Local Go Development

For rapid Go development without Docker:

```bash
# Install Air for live reload
go install github.com/air-verse/air@latest

# Start with live reload
make dev
# OR: air

# Access Go server directly
curl http://localhost:8080/test
```

### Full Stack Development

For testing the complete system:

```bash
# Start all services
make docker-dev

# View logs
docker compose -f docker-compose.dev.yml logs -f

# Restart specific service
docker compose -f docker-compose.dev.yml restart go-test-server
```

### Configuration Changes

#### Homer Dashboard

```bash
# Edit Homer configuration
vim config/homer.yml

# Copy to runtime location
cp config/homer.yml homer-config/config.yml

# Restart Homer
docker compose -f docker-compose.dev.yml restart homer
```

#### Caddy Proxy

```bash
# Edit Caddyfile
vim Caddyfile

# Reload Caddy configuration (no restart needed)
docker compose -f docker-compose.dev.yml exec caddy caddy reload --config /etc/caddy/Caddyfile
```

#### Go Server

```bash
# Edit Go code
vim cmd/webserver/main.go

# Rebuild and restart (automatic with live reload)
docker compose -f docker-compose.dev.yml up -d --build go-test-server
```

## 🧪 **Testing**

### Manual Testing

```bash
# Test main dashboard
curl -k https://localhost/

# Test Go server endpoints
curl -k https://localhost/test
curl -k https://localhost/test/api
curl -k https://localhost/test/time

# Test health check
curl -k https://localhost/health
```

### HTMX Testing

The Go server includes HTMX demonstrations:

1. **Visit**: https://localhost/test
2. **Click**: "Get Current Time" button
3. **Observe**: Dynamic content update without page reload

### Load Testing

```bash
# Simple load test with curl
for i in {1..100}; do
  curl -k -s https://localhost/health > /dev/null &
done
wait

# Or use Apache Bench
ab -n 1000 -c 10 https://localhost/
```

## 🔧 **Adding New Features**

### Adding Go Routes

1. **Add handler function** in `cmd/webserver/main.go`:
   ```go
   func handleNewFeature(w http.ResponseWriter, r *http.Request, templates *template.Template) {
       // Implementation here
   }
   ```

2. **Register route** in `main()`:
   ```go
   mux.HandleFunc("/test/newfeature", func(w http.ResponseWriter, r *http.Request) {
       handleNewFeature(w, r, templates)
   })
   ```

3. **Create template** in `templates/newfeature.html`

4. **Test the endpoint**:
   ```bash
   curl -k https://localhost/test/newfeature
   ```

### Adding Homer Services

1. **Edit configuration**:
   ```yaml
   # config/homer.yml
   services:
     - name: "New Category"
       items:
         - name: "New Service"
           url: "http://service.local:8080"
           subtitle: "Description"
   ```

2. **Update runtime config**:
   ```bash
   cp config/homer.yml homer-config/config.yml
   ```

3. **Restart Homer**:
   ```bash
   docker compose -f docker-compose.dev.yml restart homer
   ```

### Adding Caddy Routes

1. **Edit Caddyfile**:
   ```caddyfile
   # Add new route
   handle /newservice* {
       reverse_proxy service.local:8080
   }
   ```

2. **Reload configuration**:
   ```bash
   docker compose -f docker-compose.dev.yml exec caddy caddy reload --config /etc/caddy/Caddyfile
   ```

## 🐛 **Debugging**

### Common Issues

#### Port Conflicts

```bash
# Check what's using ports
netstat -tulpn | grep :80
netstat -tulpn | grep :443
netstat -tulpn | grep :8080

# Kill conflicting processes
sudo pkill -f process-name
```

#### Container Issues

```bash
# Check container status
docker compose -f docker-compose.dev.yml ps

# View container logs
docker compose -f docker-compose.dev.yml logs caddy
docker compose -f docker-compose.dev.yml logs homer
docker compose -f docker-compose.dev.yml logs go-test-server

# Restart problematic container
docker compose -f docker-compose.dev.yml restart <service-name>
```

#### SSL Certificate Issues

```bash
# Check Caddy SSL status
docker compose -f docker-compose.dev.yml logs caddy | grep -i cert

# Force certificate renewal
docker compose -f docker-compose.dev.yml exec caddy caddy reload --config /etc/caddy/Caddyfile
```

### Debugging Tools

#### Container Shell Access

```bash
# Access Caddy container
docker compose -f docker-compose.dev.yml exec caddy sh

# Access Go server container
docker compose -f docker-compose.dev.yml exec go-test-server sh

# Access Homer container
docker compose -f docker-compose.dev.yml exec homer sh
```

#### Network Debugging

```bash
# Test internal connectivity
docker compose -f docker-compose.dev.yml exec caddy ping homer
docker compose -f docker-compose.dev.yml exec caddy ping go-test-server

# Check network configuration
docker network inspect theclub-dev-network
```

#### Log Analysis

```bash
# Follow all logs
docker compose -f docker-compose.dev.yml logs -f

# Filter logs by service
docker compose -f docker-compose.dev.yml logs -f caddy

# Search logs for errors
docker compose -f docker-compose.dev.yml logs | grep -i error
```

## 📝 **Code Style and Standards**

### Go Code Standards

- **Use gofmt**: `go fmt ./...`
- **Add comments**: Document exported functions and complex logic
- **Error handling**: Always handle errors appropriately
- **Logging**: Use structured logging with context

### Configuration Standards

- **Comments**: Document all configuration options
- **Validation**: Validate configuration on startup
- **Defaults**: Provide sensible defaults
- **Environment**: Support environment variable overrides

### Docker Standards

- **Multi-stage builds**: Separate development and production stages
- **Security**: Run containers as non-root users
- **Resources**: Set appropriate resource limits
- **Health checks**: Include health check endpoints

## 🚀 **Performance Optimization**

### Go Server Optimization

```go
// Use connection pooling
server := &http.Server{
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

### Caddy Optimization

```caddyfile
# Enable compression
encode gzip zstd

# Cache static assets
@static path *.css *.js *.png *.jpg
header @static Cache-Control "public, max-age=31536000"
```

### Docker Optimization

```yaml
# Set resource limits
deploy:
  resources:
    limits:
      memory: 256M
      cpus: '0.5'
```

## 📊 **Monitoring and Metrics**

### Health Checks

- **Endpoint**: `/health`
- **Expected Response**: `200 OK`
- **Monitoring**: Include in uptime monitoring

### Logging

- **Caddy**: JSON format in `/var/log/caddy/access.log`
- **Go Server**: Structured logging to stdout
- **Homer**: Standard web server logs

### Metrics Collection

```bash
# Container metrics
docker stats

# Resource usage
docker system df

# Network metrics
docker network inspect theclub-dev-network
```

## 🔐 **Security Considerations**

### Development Security

- **HTTPS**: Always use HTTPS, even in development
- **Secrets**: Never commit secrets to version control
- **Dependencies**: Keep dependencies updated
- **Scanning**: Regularly scan for vulnerabilities

### Container Security

- **Non-root**: Run containers as non-root users
- **Read-only**: Mount configuration as read-only
- **Networks**: Use internal networks for service communication
- **Updates**: Keep base images updated

## 📚 **Additional Resources**

- [Go Documentation](https://golang.org/doc/)
- [HTMX Documentation](https://htmx.org/docs/)
- [Caddy Documentation](https://caddyserver.com/docs/)
- [Homer Documentation](https://github.com/bastienwirtz/homer)
- [Docker Compose Documentation](https://docs.docker.com/compose/)