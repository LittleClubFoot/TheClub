# Vaultwarden on Synology NAS DS415+

## Overview

This guide walks through deploying Vaultwarden (self-hosted Bitwarden-compatible password manager) on your Synology DS415+ NAS, integrated with The Club's subdomain routing.

**Why Vaultwarden on Synology?**
- **Privacy**: All passwords stored on your own hardware
- **Always Available**: Runs alongside Jellyfin on your NAS
- **Resource Efficient**: Vaultwarden is lightweight (perfect for DS415+)
- **Bitwarden Compatible**: Works with all official Bitwarden clients
- **Self-Hosted**: No third-party access to your credentials

---

## Synology DS415+ Compatibility

### Hardware Specifications
- **CPU**: Marvell Armada XP MV78230 (ARM dual-core 1.33GHz)
- **Architecture**: ARMv7 (32-bit)
- **RAM**: 1GB DDR3 (expandable to 2GB)
- **Release**: 2014

### Compatibility Check

✅ **Vaultwarden is compatible** with DS415+ via Docker:
- Vaultwarden provides ARM builds (`vaultwarden/server:latest` supports ARM)
- Docker is available for DSM 6.2+ and DSM 7.x
- Vaultwarden requires minimal resources (~100MB RAM)
- ARM32v7 architecture is supported

⚠️ **Important Considerations**:
- DS415+ uses ARMv7 (32-bit) architecture
- Ensure you're running DSM 6.2+ or DSM 7.x for best Docker support
- If running DSM 6.x, consider upgrading to DSM 7.x for security updates

---

## Prerequisites

### 1. Check DSM Version

```bash
# SSH into your Synology
ssh admin@nas.pochita.synology.me

# Check DSM version
cat /etc/VERSION
```

Required: **DSM 6.2 or newer** (DSM 7.x recommended)

### 2. Install Docker Package

If Docker isn't already installed:

1. **Open Package Center** in DSM web interface
2. **Search for "Docker"** (or "Container Manager" in DSM 7.2+)
3. **Install** the package
4. **Enable SSH** if not already enabled:
   - Control Panel → Terminal & SNMP → Enable SSH service

### 3. Create Shared Folder

Create a dedicated folder for Vaultwarden data:

1. **Control Panel → Shared Folder → Create**
2. **Name**: `docker`
3. **Description**: "Docker container data"
4. **Enable encryption**: Optional (recommended for passwords)

---

## Deployment Methods

Choose one of two deployment methods:

### Method 1: Docker CLI (Recommended for The Club Integration)
Best for integration with Caddy reverse proxy on The Club server.

### Method 2: Synology Docker GUI
Easier for beginners, less flexible for advanced configuration.

---

## Method 1: Docker CLI Deployment

### Step 1: SSH into Synology

```bash
ssh admin@nas.pochita.synology.me
# Or use your Synology's local IP
ssh admin@192.168.1.X
```

### Step 2: Create Directory Structure

```bash
# Create Vaultwarden data directory
sudo mkdir -p /volume1/docker/vaultwarden/data

# Set permissions
sudo chown -R 1000:1000 /volume1/docker/vaultwarden
```

### Step 3: Deploy Vaultwarden Container

```bash
sudo docker run -d \
  --name vaultwarden \
  --restart unless-stopped \
  -p 8100:80 \
  -v /volume1/docker/vaultwarden/data:/data \
  -e DOMAIN=https://vault.pochita.synology.me \
  -e SIGNUPS_ALLOWED=true \
  -e INVITATIONS_ALLOWED=true \
  -e SHOW_PASSWORD_HINT=false \
  -e LOG_FILE=/data/vaultwarden.log \
  -e LOG_LEVEL=info \
  -e EXTENDED_LOGGING=true \
  vaultwarden/server:latest
```

**Port Explanation**:
- `-p 8100:80`: Expose Vaultwarden on port 8100 (internal Synology network)
- The Club's Caddy will reverse proxy to `[SYNOLOGY_IP]:8100`

**Environment Variables**:
- `DOMAIN`: Your public URL (important for WebSocket support)
- `SIGNUPS_ALLOWED=true`: Allow new user registrations (disable after setup)
- `INVITATIONS_ALLOWED=true`: Allow inviting users via email
- `SHOW_PASSWORD_HINT=false`: Privacy/security best practice

### Step 4: Verify Container is Running

```bash
# Check container status
sudo docker ps | grep vaultwarden

# Check logs
sudo docker logs vaultwarden

# Test local access
curl http://localhost:8100
```

You should see the Vaultwarden web vault HTML response.

### Step 5: Disable Signups After Initial Setup

⚠️ **IMPORTANT SECURITY STEP**: After creating your admin account, disable public signups:

```bash
sudo docker stop vaultwarden
sudo docker rm vaultwarden

# Redeploy with signups disabled
sudo docker run -d \
  --name vaultwarden \
  --restart unless-stopped \
  -p 8100:80 \
  -v /volume1/docker/vaultwarden/data:/data \
  -e DOMAIN=https://vault.pochita.synology.me \
  -e SIGNUPS_ALLOWED=false \
  -e INVITATIONS_ALLOWED=true \
  -e SHOW_PASSWORD_HINT=false \
  -e LOG_FILE=/data/vaultwarden.log \
  -e LOG_LEVEL=info \
  -e EXTENDED_LOGGING=true \
  vaultwarden/server:latest
```

---

## Method 2: Synology Docker GUI Deployment

### Step 1: Open Container Manager

1. Open **Container Manager** (or **Docker** in older DSM versions)
2. Go to **Registry**
3. Search for **vaultwarden/server**
4. **Download** the image (select `latest` tag)

### Step 2: Create Container

1. Go to **Image** tab
2. Select **vaultwarden/server**
3. Click **Launch**

### Step 3: Configure Container

**General Settings**:
- **Container Name**: `vaultwarden`
- **Enable auto-restart**: ✅ Checked

**Port Settings**:
- **Local Port**: `8100` → **Container Port**: `80` (TCP)

**Volume Settings**:
- Click **Add Folder**
- Select or create: `/docker/vaultwarden/data`
- **Mount path**: `/data`

**Environment Variables** (click Advanced Settings → Environment):
```
DOMAIN=https://vault.pochita.synology.me
SIGNUPS_ALLOWED=true
INVITATIONS_ALLOWED=true
SHOW_PASSWORD_HINT=false
LOG_FILE=/data/vaultwarden.log
LOG_LEVEL=info
EXTENDED_LOGGING=true
```

**Resource Limits** (recommended):
- **Memory limit**: 256 MB
- **CPU priority**: Medium

### Step 4: Apply and Start

Click **Apply** → **Next** → **Apply** to start the container.

---

## The Club Caddy Integration

Now configure The Club's Caddy reverse proxy to route `vault.pochita.synology.me` to your Synology.

### Update Caddyfile

Add this block to `Caddyfile` (or `Caddyfile.subdomain.example`):

```caddyfile
# ============================================================================
# PASSWORD MANAGEMENT
# ============================================================================

# Vaultwarden - Self-hosted password manager (on Synology NAS)
vault.pochita.synology.me {
    # Optional: Protect with Authelia
    # Note: Vaultwarden has its own authentication
    # Only enable if you want SSO layer on top
    # forward_auth authelia:9091 {
    #     uri /api/verify?rd=https://auth.pochita.synology.me
    #     copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    # }

    reverse_proxy [SYNOLOGY_IP]:8100 {
        # Required headers for Vaultwarden
        header_up Host {upstream_hostport}
        header_up X-Real-IP {remote_host}
        header_up X-Forwarded-For {remote_host}
        header_up X-Forwarded-Proto {scheme}

        # WebSocket support (for live sync)
        header_up Connection {>Connection}
        header_up Upgrade {>Upgrade}
    }

    # Security headers
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options nosniff
        X-Frame-Options SAMEORIGIN
        X-XSS-Protection "1; mode=block"
        Referrer-Policy strict-origin-when-cross-origin

        # Vaultwarden-specific CSP
        Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' wss: https:; frame-ancestors 'self'"
    }

    # Logging
    log {
        output file /var/log/caddy/vault-access.log {
            roll_size 100mb
            roll_keep 5
            roll_keep_for 720h
        }
        format json
    }

    # Compression
    encode {
        gzip 6
        zstd
    }
}
```

**Replace `[SYNOLOGY_IP]`** with your Synology's local IP address (e.g., `192.168.1.10`).

### Apply Caddy Configuration

On The Club server:

```bash
cd ~/TheClub

# Edit Caddyfile
nano Caddyfile

# Add the vault.pochita.synology.me block above
# Replace [SYNOLOGY_IP] with actual IP

# Restart Caddy to apply changes
docker restart caddy

# Verify configuration
docker logs caddy
```

---

## DNS Configuration

Ensure your DNS is configured for the vault subdomain.

### Wildcard DNS (Recommended)

If you already have wildcard DNS configured:

```
*.pochita.synology.me → Your public IP
```

This automatically covers `vault.pochita.synology.me`.

### Specific A Record

Alternatively, add a specific DNS record:

```
vault.pochita.synology.me → Your public IP
```

**How to configure** (using Synology DDNS):
1. Log in to your DDNS provider (e.g., Synology DDNS, DuckDNS, Cloudflare)
2. Add wildcard record or specific subdomain
3. Wait for DNS propagation (can take up to 24 hours)

---

## Initial Setup

### 1. Access Vaultwarden

Open your browser and navigate to:

```
https://vault.pochita.synology.me
```

You should see the Vaultwarden web vault interface.

### 2. Create Admin Account

1. Click **Create Account**
2. Enter your email and master password
   - **Master Password**: Make it STRONG (20+ characters, random)
   - **Hint**: Leave empty (privacy best practice)
3. Click **Submit**

⚠️ **CRITICAL**: Your master password is NOT recoverable. Store it securely!

### 3. Disable Public Signups

After creating your account, disable public signups:

```bash
# SSH into Synology
ssh admin@nas.pochita.synology.me

# Update container with signups disabled
sudo docker stop vaultwarden
sudo docker rm vaultwarden

# Redeploy (see Method 1, Step 5 above)
```

### 4. Configure Vaultwarden Settings

Access admin panel (optional but recommended):

1. **Enable Admin Panel**:
   ```bash
   # SSH into Synology
   ssh admin@nas.pochita.synology.me

   # Generate admin token
   echo "MySecureAdminToken123" | docker exec -i vaultwarden /vaultwarden hash

   # Copy the hashed output
   ```

2. **Update Container with Admin Token**:
   ```bash
   sudo docker stop vaultwarden
   sudo docker rm vaultwarden

   sudo docker run -d \
     --name vaultwarden \
     --restart unless-stopped \
     -p 8100:80 \
     -v /volume1/docker/vaultwarden/data:/data \
     -e DOMAIN=https://vault.pochita.synology.me \
     -e SIGNUPS_ALLOWED=false \
     -e INVITATIONS_ALLOWED=true \
     -e SHOW_PASSWORD_HINT=false \
     -e ADMIN_TOKEN='$argon2id$...' \
     -e LOG_FILE=/data/vaultwarden.log \
     -e LOG_LEVEL=info \
     vaultwarden/server:latest
   ```

3. **Access Admin Panel**:
   - Navigate to: `https://vault.pochita.synology.me/admin`
   - Enter your plain-text admin token (e.g., `MySecureAdminToken123`)

---

## Client Setup

### Browser Extension

1. **Install Bitwarden Extension**:
   - Chrome: [Chrome Web Store](https://chrome.google.com/webstore/detail/bitwarden-free-password-m/nngceckbapebfimnlniiiahkandclblb)
   - Firefox: [Firefox Add-ons](https://addons.mozilla.org/en-US/firefox/addon/bitwarden-password-manager/)
   - Edge: [Microsoft Store](https://microsoftedge.microsoft.com/addons/detail/jbkfoedolllekgbhcbcoahefnbanhhlh)

2. **Configure Server**:
   - Click extension icon → **Settings** (gear icon)
   - **Server URL**: `https://vault.pochita.synology.me`
   - Click **Save**

3. **Log In**:
   - Enter your email and master password
   - Enable biometric unlock (optional)

### Mobile Apps

**iOS/Android**:
1. Install **Bitwarden** from App Store / Play Store
2. Tap **Settings** (gear icon)
3. **Server URL**: `https://vault.pochita.synology.me`
4. **Save** → **Log In**

### Desktop App

1. Download from [Bitwarden.com](https://bitwarden.com/download/)
2. **File → Settings → Self-hosted**
3. **Server URL**: `https://vault.pochita.synology.me`
4. **Save** → **Log In**

---

## Backup Strategy

### Automated Backups (Recommended)

**Option 1: Kopia Integration** (if using The Club's Kopia):

```bash
# On The Club server, configure Kopia to backup Synology
# This requires network access to Synology

# Add Synology path to Kopia policy
# See docs/IMPLEMENTATION_PLAN.md for Kopia setup
```

**Option 2: Synology Hyper Backup**:

1. Open **Hyper Backup** package in DSM
2. Create backup task:
   - **Source**: `/docker/vaultwarden/data`
   - **Destination**: External USB drive, another NAS, or cloud (encrypted)
   - **Schedule**: Daily at 3 AM
   - **Encryption**: ✅ Enable with strong passphrase

**Option 3: Manual Backup Script**:

Create a backup script on Synology:

```bash
# SSH into Synology
ssh admin@nas.pochita.synology.me

# Create backup script
sudo nano /usr/local/bin/backup-vaultwarden.sh
```

```bash
#!/bin/bash
# Vaultwarden Backup Script for Synology

BACKUP_DIR="/volume1/Backups/vaultwarden"
DATE=$(date +%Y%m%d_%H%M%S)
SOURCE="/volume1/docker/vaultwarden/data"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Create backup
tar -czf "$BACKUP_DIR/vaultwarden_backup_$DATE.tar.gz" -C "$SOURCE" .

# Keep only last 7 backups
ls -t "$BACKUP_DIR"/vaultwarden_backup_*.tar.gz | tail -n +8 | xargs -r rm

echo "Backup completed: vaultwarden_backup_$DATE.tar.gz"
```

```bash
# Make executable
sudo chmod +x /usr/local/bin/backup-vaultwarden.sh

# Add to crontab (daily at 2 AM)
sudo crontab -e
```

Add:
```
0 2 * * * /usr/local/bin/backup-vaultwarden.sh >> /var/log/vaultwarden-backup.log 2>&1
```

### Manual Export

From Vaultwarden web vault:
1. **Tools → Export Vault**
2. **File format**: JSON (encrypted) or CSV (unencrypted)
3. **Save** to secure location
4. ⚠️ If using CSV, encrypt the file immediately

---

## Security Best Practices

### 1. Strong Master Password

Your master password should be:
- **20+ characters** long
- **Random** combination of words (diceware method)
- **Unique** (not used anywhere else)
- **Stored offline** in a secure location

**Example good master password**:
```
correct-horse-battery-staple-mountain-purple-cloud-7
```

### 2. Enable Two-Factor Authentication

1. **Web Vault → Settings → Two-step Login**
2. **Enable Authenticator App** (TOTP):
   - Use Aegis (Android) or Raivo OTP (iOS) - both FOSS
   - Scan QR code
   - Save recovery code offline

3. **Optional**: Enable additional 2FA methods:
   - Email (requires SMTP configuration)
   - YubiKey (requires hardware key)

### 3. Regular Security Audits

In Vaultwarden admin panel:
- Review active sessions
- Check for compromised passwords (using Have I Been Pwned integration)
- Monitor failed login attempts

### 4. Firewall Rules

On your router:
- Only expose ports 80 and 443 from The Club server
- Block direct access to Synology ports from internet
- All traffic should go through Caddy reverse proxy

### 5. HTTPS Only

Vaultwarden should **NEVER** be accessed over HTTP:
- Caddy automatically handles HTTPS with Let's Encrypt
- Ensure `DOMAIN=https://vault.pochita.synology.me` is set

### 6. Admin Token Security

If you enabled the admin panel:
- Use a strong, random admin token (32+ characters)
- Only access `/admin` from trusted networks
- Consider disabling admin panel after initial configuration:
  ```bash
  # Remove ADMIN_TOKEN from container environment
  ```

---

## Integration with Authelia SSO (Optional)

You can add Authelia SSO protection on top of Vaultwarden's built-in auth.

### When to Use Authelia with Vaultwarden

**Pros**:
- Additional security layer (2FA before even reaching Vaultwarden)
- Consistent SSO experience across all services
- Centralized access control

**Cons**:
- Double authentication (Authelia + Vaultwarden)
- Can break Bitwarden clients (browser extension may not work)
- More complex setup

### How to Enable

Uncomment the `forward_auth` block in Caddyfile:

```caddyfile
vault.pochita.synology.me {
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.pochita.synology.me
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }

    reverse_proxy [SYNOLOGY_IP]:8100 {
        # ... rest of config
    }
}
```

⚠️ **Warning**: This may break Bitwarden browser extensions and mobile apps. Only enable if accessing exclusively via web vault.

**Recommended**: Skip Authelia integration for Vaultwarden. Use Vaultwarden's built-in auth + 2FA instead.

---

## Troubleshooting

### Issue: Can't Access vault.pochita.synology.me

**Check DNS**:
```bash
nslookup vault.pochita.synology.me
# Should return your public IP
```

**Check Caddy Logs**:
```bash
docker logs caddy | grep vault
```

**Check Vaultwarden Container**:
```bash
# SSH into Synology
sudo docker ps | grep vaultwarden
sudo docker logs vaultwarden
```

**Test Local Access**:
```bash
# From Synology
curl http://localhost:8100

# From The Club server
curl http://[SYNOLOGY_IP]:8100
```

### Issue: WebSocket Connection Failed

Ensure Caddy has WebSocket headers:
```caddyfile
header_up Connection {>Connection}
header_up Upgrade {>Upgrade}
```

Restart Caddy:
```bash
docker restart caddy
```

### Issue: Can't Create Account

**Verify Signups Enabled**:
```bash
# SSH into Synology
sudo docker inspect vaultwarden | grep SIGNUPS_ALLOWED
# Should show "SIGNUPS_ALLOWED=true"
```

**Check Logs**:
```bash
sudo docker logs vaultwarden | tail -50
```

### Issue: Browser Extension Not Working

**Verify Server URL**:
- Must include `https://` protocol
- Must NOT include trailing slash
- Example: `https://vault.pochita.synology.me` ✅
- NOT: `vault.pochita.synology.me` ❌
- NOT: `https://vault.pochita.synology.me/` ❌

**Clear Extension Data**:
1. Extension Settings → Options → Clear Cache
2. Log out and log back in

### Issue: Database Locked Error

This happens if multiple processes access the database:

```bash
# SSH into Synology
sudo docker restart vaultwarden

# If persists, check permissions
sudo chown -R 1000:1000 /volume1/docker/vaultwarden/data
```

### Issue: Out of Memory

DS415+ has limited RAM (1-2GB). If Vaultwarden crashes:

**Check Container Memory**:
```bash
sudo docker stats vaultwarden
```

**Add Memory Limit**:
```bash
sudo docker update --memory 256m --memory-swap 512m vaultwarden
```

**Restart Container**:
```bash
sudo docker restart vaultwarden
```

---

## Maintenance

### Update Vaultwarden

```bash
# SSH into Synology
ssh admin@nas.pochita.synology.me

# Pull latest image
sudo docker pull vaultwarden/server:latest

# Stop and remove old container
sudo docker stop vaultwarden
sudo docker rm vaultwarden

# Redeploy with latest image (use same command from Method 1, Step 3)
sudo docker run -d \
  --name vaultwarden \
  --restart unless-stopped \
  -p 8100:80 \
  -v /volume1/docker/vaultwarden/data:/data \
  -e DOMAIN=https://vault.pochita.synology.me \
  -e SIGNUPS_ALLOWED=false \
  -e INVITATIONS_ALLOWED=true \
  -e SHOW_PASSWORD_HINT=false \
  -e LOG_FILE=/data/vaultwarden.log \
  -e LOG_LEVEL=info \
  vaultwarden/server:latest

# Verify
sudo docker logs vaultwarden
```

### Monitor Logs

```bash
# SSH into Synology
ssh admin@nas.pochita.synology.me

# View real-time logs
sudo docker logs -f vaultwarden

# View recent errors
sudo docker logs vaultwarden 2>&1 | grep -i error
```

### Check Database Size

```bash
# SSH into Synology
du -sh /volume1/docker/vaultwarden/data
```

### Vacuum Database (Optimize)

```bash
# SSH into Synology
sudo docker exec -it vaultwarden sqlite3 /data/db.sqlite3 "VACUUM;"
```

---

## Migration from Bitwarden Cloud

If migrating from Bitwarden's cloud service:

### Step 1: Export from Bitwarden Cloud

1. Log in to [vault.bitwarden.com](https://vault.bitwarden.com)
2. **Tools → Export Vault**
3. **File format**: JSON (encrypted) or CSV
4. **Download** export file

### Step 2: Import to Vaultwarden

1. Log in to your Vaultwarden instance: `https://vault.pochita.synology.me`
2. **Tools → Import Data**
3. **Import file type**: Bitwarden (json) or CSV
4. **Select file** → **Import**

### Step 3: Verify Import

Check that all:
- Passwords
- Secure notes
- Identities
- Payment cards

...were imported successfully.

### Step 4: Update Client Apps

Update all browser extensions and mobile apps:
- **Settings → Server URL** → `https://vault.pochita.synology.me`
- Log out from Bitwarden cloud
- Log in to self-hosted instance

### Step 5: Delete Cloud Account (Optional)

If fully migrated:
1. Log in to [vault.bitwarden.com](https://vault.bitwarden.com)
2. **Settings → My Account → Delete Account**
3. Confirm deletion

---

## Advanced Configuration

### Email Configuration (for Invitations)

To enable email invitations and 2FA recovery:

```bash
sudo docker stop vaultwarden
sudo docker rm vaultwarden

sudo docker run -d \
  --name vaultwarden \
  --restart unless-stopped \
  -p 8100:80 \
  -v /volume1/docker/vaultwarden/data:/data \
  -e DOMAIN=https://vault.pochita.synology.me \
  -e SIGNUPS_ALLOWED=false \
  -e INVITATIONS_ALLOWED=true \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_FROM=yourname@gmail.com \
  -e SMTP_PORT=587 \
  -e SMTP_SECURITY=starttls \
  -e SMTP_USERNAME=yourname@gmail.com \
  -e SMTP_PASSWORD=your_app_password \
  vaultwarden/server:latest
```

**Gmail App Password**:
- Go to [Google Account Security](https://myaccount.google.com/security)
- Enable 2FA
- Generate App Password for Vaultwarden

### Organization Support

Vaultwarden supports Bitwarden Organizations (shared passwords):

1. **Web Vault → New Organization**
2. **Name**: "Family" or "Team"
3. **Invite members** via email
4. **Share collections** of passwords

### Monitoring with Uptime Kuma

Add Vaultwarden to Uptime Kuma:

1. Open **Uptime Kuma**: `https://status.pochita.synology.me`
2. **Add New Monitor**:
   - **Monitor Type**: HTTP(s)
   - **Friendly Name**: Vaultwarden
   - **URL**: `https://vault.pochita.synology.me`
   - **Heartbeat Interval**: 60 seconds
   - **Accepted Status Codes**: 200-299

---

## Performance Optimization

### Enable SQLite WAL Mode

Write-Ahead Logging improves performance:

```bash
# SSH into Synology
sudo docker exec -it vaultwarden sqlite3 /data/db.sqlite3 "PRAGMA journal_mode=WAL;"
```

### Adjust Resource Limits

For DS415+ with limited resources:

```bash
sudo docker update --memory 256m --memory-swap 512m --cpus 0.5 vaultwarden
```

---

## Privacy Considerations

### Data Stored on Vaultwarden

Your Vaultwarden instance stores:
- Encrypted vault data (passwords, notes, etc.)
- User account information (email, hashed master password)
- Device/session tokens
- Audit logs (login attempts, changes)

### Data NOT Sent to Third Parties

Unlike Bitwarden cloud:
- ✅ **No** data sent to Bitwarden servers
- ✅ **No** telemetry or analytics
- ✅ **No** third-party access
- ✅ **All data** stays on your Synology

### Network Privacy

- All traffic encrypted via HTTPS (Let's Encrypt)
- Passwords never transmitted in plain text
- End-to-end encryption (encrypted before leaving device)

---

## Summary

You now have:
- ✅ Vaultwarden running on Synology DS415+
- ✅ Accessible via `https://vault.pochita.synology.me`
- ✅ Integrated with The Club's Caddy reverse proxy
- ✅ Automatic HTTPS with Let's Encrypt
- ✅ Bitwarden-compatible (all official clients work)
- ✅ Privacy-focused (all data on your hardware)
- ✅ Automated backups (via Hyper Backup or Kopia)

**Next Steps**:
1. Install browser extensions and mobile apps
2. Import passwords from previous password manager
3. Enable 2FA on your account
4. Configure automated backups
5. Invite family members (optional)
6. Add to Homer dashboard

---

## Related Documentation

- **Vaultwarden Wiki**: https://github.com/dani-garcia/vaultwarden/wiki
- **Bitwarden Help**: https://bitwarden.com/help/
- **The Club - Synology Deployment**: `docs/SYNOLOGY_DEPLOYMENT.md`
- **The Club - Network Architecture**: `docs/NETWORK_ARCHITECTURE.md`
- **The Club - Homer Dashboard**: `docs/HOMER_CUSTOMIZATION.md`

---

## Quick Reference

```bash
# Check Vaultwarden status
sudo docker ps | grep vaultwarden

# View logs
sudo docker logs vaultwarden

# Restart container
sudo docker restart vaultwarden

# Update to latest version
sudo docker pull vaultwarden/server:latest && \
  sudo docker stop vaultwarden && \
  sudo docker rm vaultwarden && \
  # Then redeploy with same parameters

# Backup database
tar -czf vaultwarden_backup.tar.gz -C /volume1/docker/vaultwarden/data .

# Check database size
du -sh /volume1/docker/vaultwarden/data
```

---

**Your passwords are now truly yours. 🔒**
