#!/bin/bash
# ============================================================================
# The Club - Configuration Backup Script
# ============================================================================
# Backs up all critical configuration files and Docker volumes
#
# Usage:
#   ./scripts/config-backup.sh
#   ./scripts/config-backup.sh /custom/backup/path
#
# Backups are stored in: /backups/theclub-YYYYMMDD-HHMMSS/
# ============================================================================

set -e  # Exit on error
set -u  # Exit on undefined variable

# Configuration
BACKUP_BASE_DIR="${1:-/backups}"
TIMESTAMP="$(date +%Y%m%d-%H%M%S)"
BACKUP_DIR="${BACKUP_BASE_DIR}/theclub-${TIMESTAMP}"
PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🏴‍☠️ The Club - Configuration Backup${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# Create backup directory
echo -e "${YELLOW}Creating backup directory...${NC}"
mkdir -p "$BACKUP_DIR"
echo -e "${GREEN}✓ Created: $BACKUP_DIR${NC}"
echo ""

# Backup configuration files
echo -e "${YELLOW}Backing up configuration files...${NC}"
cd "$PROJECT_ROOT"

tar -czf "$BACKUP_DIR/configs.tar.gz" \
    Caddyfile \
    docker-compose.prod.yml \
    docker-compose.dev.yml \
    docker-compose.monitoring.yml \
    assets/config.yml \
    Makefile \
    .air.toml \
    go.mod \
    go.sum \
    2>/dev/null || true

echo -e "${GREEN}✓ Configuration files backed up${NC}"
echo ""

# Backup Docker volumes
echo -e "${YELLOW}Backing up Docker volumes...${NC}"

# Caddy data volume
if docker volume inspect caddy_data >/dev/null 2>&1; then
    echo "  - Backing up caddy_data..."
    docker run --rm \
        -v caddy_data:/source:ro \
        -v "$BACKUP_DIR":/backup \
        alpine tar -czf /backup/caddy_data.tar.gz -C /source . 2>/dev/null || true
    echo -e "${GREEN}  ✓ caddy_data backed up${NC}"
fi

# Caddy config volume
if docker volume inspect caddy_config >/dev/null 2>&1; then
    echo "  - Backing up caddy_config..."
    docker run --rm \
        -v caddy_config:/source:ro \
        -v "$BACKUP_DIR":/backup \
        alpine tar -czf /backup/caddy_config.tar.gz -C /source . 2>/dev/null || true
    echo -e "${GREEN}  ✓ caddy_config backed up${NC}"
fi

# Uptime Kuma data (if exists)
if docker volume inspect uptime-kuma-data >/dev/null 2>&1; then
    echo "  - Backing up uptime-kuma-data..."
    docker run --rm \
        -v uptime-kuma-data:/source:ro \
        -v "$BACKUP_DIR":/backup \
        alpine tar -czf /backup/uptime-kuma-data.tar.gz -C /source . 2>/dev/null || true
    echo -e "${GREEN}  ✓ uptime-kuma-data backed up${NC}"
fi

echo ""

# Backup templates and source code
echo -e "${YELLOW}Backing up source code...${NC}"
tar -czf "$BACKUP_DIR/source.tar.gz" \
    cmd/ \
    templates/ \
    docs/ \
    2>/dev/null || true
echo -e "${GREEN}✓ Source code backed up${NC}"
echo ""

# Create backup manifest
echo -e "${YELLOW}Creating backup manifest...${NC}"
cat > "$BACKUP_DIR/manifest.txt" <<EOF
The Club Backup Manifest
========================
Date: $(date)
Hostname: $(hostname)
Backup Directory: $BACKUP_DIR

Contents:
- configs.tar.gz        : Configuration files
- caddy_data.tar.gz     : Caddy SSL certificates and data
- caddy_config.tar.gz   : Caddy configuration
- source.tar.gz         : Application source code
- uptime-kuma-data.tar.gz : Monitoring data (if exists)

Restore Instructions:
1. Extract configs: tar -xzf configs.tar.gz
2. Extract volumes: docker run --rm -v volume_name:/dest -v $(pwd):/backup alpine tar -xzf /backup/volume_name.tar.gz -C /dest
3. Restart services: make docker-prod

Total Size: $(du -sh "$BACKUP_DIR" | cut -f1)
EOF
echo -e "${GREEN}✓ Manifest created${NC}"
echo ""

# Calculate total size
TOTAL_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)

# Summary
echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}✅ Backup Complete!${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo -e "Location: ${BLUE}$BACKUP_DIR${NC}"
echo -e "Size: ${BLUE}$TOTAL_SIZE${NC}"
echo -e "Files: ${BLUE}$(ls -1 "$BACKUP_DIR" | wc -l)${NC}"
echo ""
echo -e "To restore this backup, see: ${BLUE}$BACKUP_DIR/manifest.txt${NC}"
echo ""

# Optional: Clean up old backups (keep last 7 days)
if [ -d "$BACKUP_BASE_DIR" ]; then
    echo -e "${YELLOW}Cleaning up old backups (keeping last 7 days)...${NC}"
    find "$BACKUP_BASE_DIR" -type d -name "theclub-*" -mtime +7 -exec rm -rf {} + 2>/dev/null || true
    echo -e "${GREEN}✓ Cleanup complete${NC}"
fi

exit 0
