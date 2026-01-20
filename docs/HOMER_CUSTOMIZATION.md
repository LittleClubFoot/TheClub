# Homer Dashboard Customization Guide

## Overview

This guide explains the enhanced Homer dashboard configuration for The Club, showcasing your privacy-focused home server infrastructure.

---

## What's New

### Enhanced Configuration
- **Privacy-First Organization**: Services grouped by function with privacy infrastructure highlighted
- **Updated Color Scheme**: Purple accent colors representing privacy focus
- **Comprehensive Service Catalog**: All services (privacy, monitoring, security, management, media)
- **Helpful Subtitles**: Each service explains its purpose
- **Keywords for Search**: Built-in search functionality
- **Privacy-Focused Links**: Quick access to EFF, Privacy Tools, and documentation

### Service Categories

#### 🔒 Privacy Infrastructure
- **Pi-hole DNS**: Network-wide ad and tracker blocking
- **Searxng**: Privacy-respecting meta-search engine
- **FreshRSS**: Self-hosted RSS reader

#### 📊 Monitoring
- **Uptime Kuma**: Service health and status monitoring
- **Grafana**: Metrics dashboards and visualization

#### 🔐 Security
- **Authelia SSO**: Single sign-on with 2FA
- **WireGuard VPN**: Secure remote access
- **Vaultwarden**: Self-hosted password manager (Synology)

#### 🛠️ Management
- **Dockge**: Docker Compose stack management
- **Kopia**: Encrypted backup management

#### 🎬 Media & Files
- **Jellyfin**: Media streaming (Synology NAS)
- **Synology NAS**: Network storage and DSM admin

#### 💻 Development
- **Test Server**: Go/HTMX test application
- **API Docs**: API documentation

#### ⚙️ System
- **Health Check**: System health endpoint
- **Prometheus**: Metrics collection (internal)

---

## Customization

### Update Domain

Change all `pochita.synology.me` URLs to your actual domain:

```bash
cd ~/TheClub
nano assets/config.yml

# Replace all instances of pochita.synology.me with your domain
:%s/pochita.synology.me/yourdomain.com/g
```

### Change Colors

Edit the color scheme in `assets/config.yml` (lines 24-48):

```yaml
colors:
  light:
    highlight-primary: "#667eea"  # Change to your preferred color
    highlight-secondary: "#764ba2"
    # ... etc
```

**Popular Color Schemes:**

**Privacy Purple (Current)**:
```yaml
highlight-primary: "#667eea"
highlight-secondary: "#764ba2"
```

**Security Blue**:
```yaml
highlight-primary: "#3367d6"
highlight-secondary: "#4285f4"
```

**Nature Green**:
```yaml
highlight-primary: "#10b981"
highlight-secondary: "#059669"
```

**Hacker Dark**:
```yaml
highlight-primary: "#00ff00"
highlight-secondary: "#00cc00"
```

### Add Custom Icons

Homer supports custom icons for services. Place PNG files in `assets/icons/`:

#### Required Icons (Reference in config)

```
assets/icons/
├── pihole.png      # Pi-hole logo
├── search.png      # Searxng/search icon
├── rss.png         # RSS icon
├── uptime.png      # Uptime Kuma logo
├── grafana.png     # Grafana logo
├── auth.png        # Authelia/lock icon
├── vpn.png         # WireGuard/VPN icon
├── vault.png       # Vaultwarden/password manager icon
├── docker.png      # Docker/Dockge logo
├── backup.png      # Backup/Kopia icon
├── jellyfin.png    # Jellyfin logo
├── nas.png         # Synology/NAS icon
├── dev.png         # Development icon
├── api.png         # API icon
├── health.png      # Health check icon
└── prometheus.png  # Prometheus logo
```

#### Icon Sources

1. **Official Logos**:
   - [Jellyfin](https://jellyfin.org/images/logo.svg)
   - [Grafana](https://grafana.com/static/assets/img/grafana_icon.svg)
   - [Prometheus](https://prometheus.io/assets/prometheus_logo.svg)
   - [Vaultwarden](https://github.com/bitwarden/brand/blob/master/icons/icon.png) (use Bitwarden logo)

2. **Icon Packs**:
   - [Font Awesome](https://fontawesome.com/) (built-in)
   - [Dashboard Icons](https://github.com/walkxcode/dashboard-icons)
   - [Homer Icons](https://github.com/NX211/homer-icons)

3. **Create Your Own**:
   ```bash
   # Recommended size: 128x128px or 256x256px
   # Format: PNG with transparency
   # Style: Flat, modern icons work best
   ```

#### Using Font Awesome Instead of Custom Icons

If you don't want to manage custom icons, use Font Awesome:

```yaml
- name: "Pi-hole DNS"
  icon: "fas fa-shield-alt"  # Instead of logo
  subtitle: "Network-wide ad & tracker blocking"
  # ... rest of config
```

Available Font Awesome icons:
- `fas fa-shield-alt` (Privacy/Security)
- `fas fa-search` (Search)
- `fas fa-rss` (RSS)
- `fas fa-chart-line` (Monitoring)
- `fas fa-lock` (Security)
- `fas fa-key` (VPN)
- `fas fa-cubes` (Docker)
- `fas fa-database` (Backup)
- `fas fa-film` (Media)
- `fas fa-server` (NAS)

### Add/Remove Services

#### Add a New Service

```yaml
- name: "Your Service"
  logo: "assets/icons/yourservice.png"
  subtitle: "Brief description"
  tag: "category"
  keywords: "search terms"
  url: "https://service.pochita.synology.me"
  target: "_blank"
```

#### Remove a Service

Comment out or delete the service block:

```yaml
# - name: "Service to Remove"
#   logo: "assets/icons/service.png"
#   # ... rest commented out
```

### Add New Category

```yaml
- name: "🆕 New Category"
  icon: "fas fa-star"
  items:
    - name: "Service 1"
      # ... service config
```

### Change Footer

Edit line 14 in `assets/config.yml`:

```yaml
footer: '<p>Your custom footer text here</p>'
```

**Examples:**

```yaml
# Privacy-focused
footer: '<p>Privacy-First Home Lab 🔒</p>'

# Hacker-style
footer: '<p>[ root@homelab ]# uptime</p>'

# Minimalist
footer: '<p>The Club</p>'

# No footer
footer: false
```

---

## Themes

Homer supports different visual themes.

### Available Themes

**Default (Current)**:
```yaml
theme: default
```

**SUI (Semantic UI)**:
```yaml
theme: sui
```

### Custom CSS

Create custom styles in `assets/custom.css`:

```css
/* Privacy-focused custom styles */
.card {
  border-left: 3px solid #667eea;
}

.tag.privacy {
  background-color: #667eea;
}

.tag.security {
  background-color: #10b981;
}
```

Then reference in config:

```yaml
stylesheet:
  - "assets/custom.css"
```

---

## Search Functionality

Homer includes built-in search. Use keywords to make services discoverable:

```yaml
- name: "Pi-hole DNS"
  keywords: "dns ad-block privacy tracking pihole"
  # User can search "ad" or "dns" or "privacy" to find this
```

**Best Practices:**
- Include service name variations
- Add function keywords
- Add technology keywords
- Include common misspellings

---

## Responsive Design

### Column Configuration

Homer automatically adjusts columns based on screen size:

```yaml
columns: "auto"  # Recommended - auto-adjusts
# OR
columns: "3"     # Fixed 3 columns
# OR
columns: "4"     # Fixed 4 columns
```

**Auto Mode Breakpoints:**
- Desktop (>1408px): 4 columns
- Laptop (1024-1407px): 3 columns
- Tablet (768-1023px): 2 columns
- Mobile (<768px): 1 column

### Mobile Optimization

Services are already optimized for mobile with:
- Touch-friendly tap targets
- Responsive images
- Readable font sizes
- Collapsible categories

---

## Advanced Features

### Connectivity Check

Automatically checks if services are online:

```yaml
connectivityCheck: true  # Enabled by default
```

Shows a green/red indicator on each service card.

### External Links

Control whether links open in new tabs:

```yaml
target: "_blank"  # Opens in new tab
target: "_self"   # Opens in same tab
```

**Recommended:**
- `_blank`: External services (Grafana, Pi-hole, etc.)
- `_self`: Internal pages (Test server, Health check)

### Service Tags

Visual tags to categorize services:

```yaml
tag: "privacy"  # Shows colored tag on service card
```

**Current Tags:**
- `privacy` (purple)
- `monitor` (blue)
- `security` (green)
- `admin` (orange)
- `media` (red)
- `dev` (gray)
- `system` (yellow)

---

## Deployment

### Apply Changes

```bash
# Homer auto-reloads config changes
# No restart needed!

# If changes don't appear:
docker restart homer

# Or rebuild:
docker compose -f docker-compose.prod.yml up -d --force-recreate homer
```

### Validate Configuration

```bash
# Check YAML syntax
yamllint assets/config.yml

# Or use online validator
# https://www.yamllint.com/
```

---

## Troubleshooting

### Services Not Showing

**Problem**: Added service but doesn't appear

**Solutions:**
1. Check YAML indentation (must be consistent)
2. Verify no syntax errors
3. Restart Homer: `docker restart homer`
4. Check browser cache (hard refresh: Ctrl+Shift+R)

### Icons Not Loading

**Problem**: Custom icons show placeholder

**Solutions:**
1. Verify icon path: `assets/icons/filename.png`
2. Check file permissions: `chmod 644 assets/icons/*.png`
3. Verify PNG format (not JPEG, WebP, etc.)
4. Check icon size (recommend 128x128 or 256x256)
5. Use Font Awesome as fallback

### Colors Not Applying

**Problem**: Custom colors don't show

**Solutions:**
1. Check color format: Must be hex (#RRGGBB)
2. Clear browser cache
3. Verify both light and dark themes defined
4. Hard refresh browser (Ctrl+Shift+R)

### Search Not Working

**Problem**: Can't find services via search

**Solutions:**
1. Add `keywords` field to services
2. Ensure search bar visible (top-right)
3. Check browser console for JavaScript errors
4. Try different search terms

---

## Examples

### Minimal Configuration

```yaml
title: "Home Lab"
subtitle: "Services"
services:
  - name: "Services"
    items:
      - name: "Service 1"
        url: "https://service1.example.com"
      - name: "Service 2"
        url: "https://service2.example.com"
```

### Maximum Privacy Configuration

```yaml
title: "Private Server"
subtitle: "Self-Hosted & Encrypted"
theme: default
colors:
  dark:
    background: "#000000"
    card-background: "#1a1a1a"
    highlight-primary: "#00ff00"
services:
  - name: "🔒 Privacy Tools"
    items:
      - name: "Pi-hole"
        url: "https://dns.example.com"
        tag: "privacy"
      - name: "VPN"
        url: "#"
        tag: "privacy"
```

---

## Best Practices

### Organization
- Group related services together
- Use emoji in category names (🔒, 📊, 🔐)
- Order by frequency of use
- Separate public vs. protected services

### Naming
- Keep service names short (1-2 words)
- Use descriptive subtitles
- Include technology name if helpful
- Be consistent with naming style

### Performance
- Use lightweight icons (<100KB each)
- Limit to 20-30 services max
- Enable connectivity check only for critical services
- Use Font Awesome for frequently accessed services (faster)

### Security
- Don't expose internal URLs publicly
- Use authentication for sensitive dashboards
- Keep services behind VPN when possible
- Regularly review exposed services

### Maintenance
- Document custom changes
- Keep config in version control
- Test after updates
- Back up config regularly (`make backup`)

---

## Related Documentation

- **Homer Official Docs**: https://github.com/bastienwirtz/homer
- **Font Awesome Icons**: https://fontawesome.com/icons
- **YAML Syntax**: https://yaml.org/
- **Color Picker**: https://htmlcolorcodes.com/

---

## Quick Reference

```bash
# Edit configuration
nano assets/config.yml

# Restart Homer
docker restart homer

# View logs
docker logs homer

# Validate YAML
yamllint assets/config.yml

# Backup config
make backup
```

---

## Summary

Your enhanced Homer dashboard now:
- ✅ Showcases all privacy infrastructure
- ✅ Organized by logical categories
- ✅ Privacy-focused color scheme
- ✅ Comprehensive service descriptions
- ✅ Search functionality with keywords
- ✅ Mobile-responsive design
- ✅ Easy to customize and extend

**Next Steps:**
1. Add custom icons (optional)
2. Adjust colors to your preference
3. Add/remove services as needed
4. Share with family/roommates!
