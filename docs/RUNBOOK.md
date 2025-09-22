# Operations Runbook

This runbook provides step-by-step procedures for operating and maintaining The Club dashboard.

## 🚀 **Deployment Procedures**

### Initial Deployment

1. **Prerequisites Check**
   ```bash
   # Verify Docker is running
   docker --version
   docker compose version
   
   # Check available resources
   docker system df
   free -h
   ```

2. **Environment Setup**
   ```bash
   # Clone repository
   git clone <repository-url>
   cd TheClub
   
   # Verify configuration files
   ls -la config/
   ls -la Caddyfile
   ```

3. **Development Deployment**
   ```bash
   # Start services
   make docker-dev
   
   # Verify services are running
   docker compose -f docker-compose.dev.yml ps
   
   # Check logs
   docker compose -f docker-compose.dev.yml logs
   ```

4. **Production Deployment**
   ```bash
   # Update Caddyfile with your domain
   vim Caddyfile
   
   # Deploy production stack
   make docker-prod
   
   # Verify deployment
   curl -k https://yourdomain.com/health
   ```

### Rolling Updates

1. **Update Application Code**
   ```bash
   # Pull latest changes
   git pull origin main
   
   # Rebuild and restart Go server
   docker compose -f docker-compose.prod.yml up -d --build go-test-server
   ```

2. **Update Configuration**
   ```bash
   # Update Homer config
   vim config/homer.yml
   cp config/homer.yml homer-config/config.yml
   
   # Restart Homer
   docker compose -f docker-compose.prod.yml restart homer
   ```

3. **Update Caddy Configuration**
   ```bash
   # Update Caddyfile
   vim Caddyfile
   
   # Reload Caddy configuration
   docker compose -f docker-compose.prod.yml exec caddy caddy reload --config /etc/caddy/Caddyfile
   ```

## 🔍 **Monitoring and Health Checks**

### Service Health Checks

```bash
# Check all container status
docker compose -f docker-compose.prod.yml ps

# Health check endpoints
curl -k https://yourdomain.com/health
curl -k https://yourdomain.com/test

# Check individual services
curl -k https://yourdomain.com/        # Homer
curl -k https://yourdomain.com/test    # Go server
```

### Log Monitoring

```bash
# View all logs
docker compose -f docker-compose.prod.yml logs

# Follow logs in real-time
docker compose -f docker-compose.prod.yml logs -f

# Service-specific logs
docker compose -f docker-compose.prod.yml logs caddy
docker compose -f docker-compose.prod.yml logs homer
docker compose -f docker-compose.prod.yml logs go-test-server

# Caddy access logs
docker compose -f docker-compose.prod.yml exec caddy tail -f /var/log/caddy/access.log
```

### Resource Monitoring

```bash
# Container resource usage
docker stats

# Disk usage
docker system df
du -sh logs/

# Network connectivity
docker network ls
docker network inspect theclub_club-network
```

## 🛠️ **Troubleshooting**

### Common Issues

#### 1. Services Not Starting

**Symptoms**: Containers exit immediately or fail to start

**Diagnosis**:
```bash
# Check container status
docker compose -f docker-compose.prod.yml ps -a

# Check logs for errors
docker compose -f docker-compose.prod.yml logs

# Check resource usage
docker system df
free -h
```

**Solutions**:
- Ensure sufficient disk space and memory
- Check for port conflicts: `netstat -tulpn | grep :80`
- Verify configuration files are valid
- Restart Docker daemon: `sudo systemctl restart docker`

#### 2. SSL Certificate Issues

**Symptoms**: SSL errors, certificate warnings

**Diagnosis**:
```bash
# Check Caddy logs
docker compose -f docker-compose.prod.yml logs caddy

# Test SSL certificate
openssl s_client -connect yourdomain.com:443 -servername yourdomain.com
```

**Solutions**:
- Ensure domain points to your server
- Check firewall allows ports 80/443
- Verify Caddyfile syntax: `docker compose -f docker-compose.prod.yml exec caddy caddy validate --config /etc/caddy/Caddyfile`

#### 3. Service Connectivity Issues

**Symptoms**: 502/503 errors, services unreachable

**Diagnosis**:
```bash
# Check network connectivity
docker compose -f docker-compose.prod.yml exec caddy ping homer
docker compose -f docker-compose.prod.yml exec caddy ping go-test-server

# Check service ports
docker compose -f docker-compose.prod.yml exec homer netstat -tulpn
```

**Solutions**:
- Verify service names in Caddyfile match Docker Compose services
- Check internal service ports (8080 for Homer, 8080 for Go server)
- Restart affected services

#### 4. Configuration Not Loading

**Symptoms**: Changes not reflected, default configuration shown

**Diagnosis**:
```bash
# Check file mounts
docker compose -f docker-compose.prod.yml exec homer ls -la /www/assets/
docker compose -f docker-compose.prod.yml exec caddy ls -la /etc/caddy/

# Verify file contents
docker compose -f docker-compose.prod.yml exec homer cat /www/assets/config.yml
```

**Solutions**:
- Ensure configuration files exist and are readable
- Restart services after configuration changes
- Check file permissions: `chmod 644 config/homer.yml`

## 🔧 **Maintenance Procedures**

### Regular Maintenance

#### Daily
- Monitor service health via `/health` endpoint
- Check disk space: `df -h`
- Review error logs

#### Weekly
- Update Docker images: `docker compose pull`
- Clean up unused containers: `docker system prune`
- Backup configuration files

#### Monthly
- Review and rotate logs
- Update base system packages
- Review security updates

### Backup Procedures

```bash
# Backup configuration
tar -czf backup-$(date +%Y%m%d).tar.gz \
  config/ \
  homer-config/ \
  Caddyfile \
  docker-compose.prod.yml

# Backup logs
tar -czf logs-backup-$(date +%Y%m%d).tar.gz logs/

# Backup to remote location
rsync -av backup-*.tar.gz user@backup-server:/backups/theclub/
```

### Recovery Procedures

#### Service Recovery
```bash
# Restart individual service
docker compose -f docker-compose.prod.yml restart <service-name>

# Restart all services
docker compose -f docker-compose.prod.yml restart

# Full redeployment
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d
```

#### Configuration Recovery
```bash
# Restore from backup
tar -xzf backup-YYYYMMDD.tar.gz

# Reload configuration
docker compose -f docker-compose.prod.yml restart
```

## 📊 **Performance Tuning**

### Resource Optimization

```bash
# Adjust container resource limits in docker-compose.prod.yml
deploy:
  resources:
    limits:
      memory: 256M
      cpus: '0.5'
    reservations:
      memory: 128M
      cpus: '0.25'
```

### Caddy Optimization

```caddyfile
# Add to Caddyfile for better performance
{
    # Enable compression
    encode gzip zstd
    
    # Cache static assets
    @static {
        path *.css *.js *.png *.jpg *.gif *.ico *.woff *.woff2
    }
    header @static Cache-Control "public, max-age=31536000"
}
```

## 🔐 **Security Procedures**

### Security Hardening

1. **Container Security**
   ```bash
   # Run containers as non-root
   # Already configured in docker-compose files
   
   # Enable security options
   security_opt:
     - no-new-privileges:true
   ```

2. **Network Security**
   ```bash
   # Restrict network access
   # Use internal networks for service communication
   
   # Monitor network connections
   netstat -tulpn
   ```

3. **SSL/TLS Security**
   ```caddyfile
   # Strong SSL configuration in Caddyfile
   tls {
     protocols tls1.2 tls1.3
   }
   ```

### Security Monitoring

```bash
# Check for security updates
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.CreatedAt}}"

# Scan for vulnerabilities (if available)
docker scout cves theclub-go-test-server

# Monitor access logs
tail -f logs/access.log | grep -E "(40[0-9]|50[0-9])"
```

## 📞 **Emergency Procedures**

### Service Outage

1. **Immediate Response**
   ```bash
   # Check service status
   docker compose -f docker-compose.prod.yml ps
   
   # Quick restart
   docker compose -f docker-compose.prod.yml restart
   ```

2. **If Restart Fails**
   ```bash
   # Full redeployment
   docker compose -f docker-compose.prod.yml down
   docker compose -f docker-compose.prod.yml up -d
   ```

3. **If Still Failing**
   ```bash
   # Check system resources
   free -h
   df -h
   
   # Check Docker daemon
   sudo systemctl status docker
   sudo systemctl restart docker
   ```

### Data Recovery

```bash
# If configuration is corrupted
git checkout HEAD -- config/
cp config/homer.yml homer-config/config.yml

# If containers are corrupted
docker compose -f docker-compose.prod.yml down
docker system prune -a
docker compose -f docker-compose.prod.yml up -d
```

## 📋 **Checklists**

### Pre-Deployment Checklist

- [ ] Configuration files reviewed and updated
- [ ] Domain DNS configured correctly
- [ ] Firewall ports 80/443 open
- [ ] Sufficient system resources available
- [ ] Backup of current configuration taken
- [ ] NAS services accessible from deployment server

### Post-Deployment Checklist

- [ ] All containers running and healthy
- [ ] SSL certificate obtained and valid
- [ ] All service URLs accessible
- [ ] Health check endpoint responding
- [ ] Logs show no errors
- [ ] Performance metrics within acceptable ranges

### Maintenance Checklist

- [ ] System updates applied
- [ ] Docker images updated
- [ ] Logs reviewed and rotated
- [ ] Backup completed successfully
- [ ] Security scan completed
- [ ] Performance metrics reviewed