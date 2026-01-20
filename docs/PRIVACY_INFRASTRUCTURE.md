# The Club - Privacy Infrastructure Guide

## Overview

This guide covers the deployment and configuration of privacy-focused services that protect your entire network from tracking, profiling, and surveillance.

---

## Privacy Philosophy

**Problem**: Modern internet is filled with tracking, profiling, and surveillance:
- Google/Facebook tracking pixels on 90% of websites
- ISPs logging and selling your DNS queries
- Search engines building profiles of your interests
- RSS readers tracking what news you read

**Solution**: Self-hosted privacy infrastructure
- Block tracking at DNS level (protects ALL devices)
- Private DNS resolution (no third-party logging)
- Private search (no profiling or tracking)
- Self-hosted RSS (read news without exposing habits)

---

## Services Included

### 1. Pi-hole - Network-Wide Ad & Tracker Blocking
**Privacy Benefit**: Blocks ads and trackers BEFORE they reach your devices

- DNS-level blocking (works on ALL devices, no per-device configuration)
- Blocks ads in apps, not just browsers
- Prevents tracking pixels from loading
- Reduces bandwidth usage
- Speeds up browsing (fewer requests)

**What it blocks**:
- Advertising networks (Google Ads, DoubleClick, etc.)
- Analytics and tracking (Google Analytics, Facebook Pixel, etc.)
- Telemetry (Windows, Apple, smart TVs)
- Malware and phishing domains

**Access**: `https://dns.pochita.synology.me/admin`

### 2. Unbound - Private Recursive DNS Resolver
**Privacy Benefit**: Your DNS queries never leave your network to third parties

- Queries root DNS servers directly (no Google, Cloudflare, ISP)
- No logging of your DNS queries to external parties
- DNSSEC validation (prevents DNS spoofing)
- QNAME minimization (sends minimal info to servers)

**Why it matters**:
- Your ISP can't see what websites you visit
- Google/Cloudflare can't build profiles from your DNS queries
- Reduces metadata leakage
- True privacy at the DNS level

**Access**: Internal service (no web interface, used by Pi-hole)

### 3. Searxng - Privacy-Respecting Meta-Search Engine
**Privacy Benefit**: Search without Google/Bing profiling you

- Aggregates results from multiple search engines
- No tracking cookies or JavaScript
- No search history logging
- No IP-based profiling
- Image proxy (don't load images directly from search engines)

**Search engines used** (privacy-friendly only):
- DuckDuckGo
- Brave Search
- Startpage
- Qwant
- (Google/Bing disabled by default for privacy)

**Access**: `https://search.pochita.synology.me`

### 4. FreshRSS - Self-Hosted RSS Reader
**Privacy Benefit**: Read news without exposing browsing habits

- Self-hosted (no cloud service tracking what you read)
- Direct feed fetching (websites don't see your IP for every visit)
- No third-party analytics
- Aggregates all your news in one place

**Why it matters**:
- News sites can't track your reading habits
- No ads or tracking in RSS feeds
- Read content without visiting tracking-heavy websites
- Control your own data

**Access**: `https://rss.pochita.synology.me`

---

## Deployment

### Prerequisites

- The Club base stack deployed
- Wildcard DNS configured (`*.pochita.synology.me`)
- At least 1GB free RAM

### Step 1: Generate Secrets

```bash
# Generate Pi-hole password
openssl rand -base64 32

# Generate Searxng secret
openssl rand -base64 48
```

### Step 2: Update Configuration

Edit `docker-compose.privacy.yml`:

```yaml
# Line 72: Update Pi-hole password
- WEBPASSWORD=your_generated_password_here

# Line 73: Update to your Club server IP
- FTLCONF_LOCAL_IPV4=192.168.1.100  # Change this

# Line 129: Update Searxng secret
- SEARXNG_SECRET=your_generated_secret_here
```

### Step 3: Deploy Privacy Stack

```bash
cd ~/TheClub

# Deploy all privacy services
docker compose -f docker-compose.privacy.yml up -d

# Wait for services to start (60 seconds)
sleep 60

# Check status
docker ps | grep -E 'pihole|unbound|searxng|freshrss'
```

### Step 4: Configure Pi-hole

1. **Access Pi-hole**:
   - Visit: `https://dns.pochita.synology.me/admin`
   - Login with password from Step 2

2. **Add Blocklists** (Settings → Blocklists):
   ```
   https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
   https://raw.githubusercontent.com/anudeepND/blacklist/master/adservers.txt
   https://dbl.oisd.nl/
   ```

3. **Update Gravity** (Tools → Update Gravity):
   - Click "Update" to download blocklists

4. **Verify Unbound Integration**:
   - Settings → DNS → Upstream DNS Servers
   - Should show custom DNS (Unbound)

### Step 5: Configure FreshRSS

1. **Initial Setup**:
   - Visit: `https://rss.pochita.synology.me`
   - Follow setup wizard
   - Create admin account

2. **Add Feeds**:
   - Recommended starting feeds (privacy-focused news):
     - EFF: `https://www.eff.org/rss/updates.xml`
     - PrivacyTools: `https://blog.privacytools.io/feed/`
     - Techdirt: `https://www.techdirt.com/techdirt_rss.xml`

3. **Configure Auto-Update**:
   - Already configured in docker-compose (every 15 minutes)

### Step 6: Configure Searxng

1. **Access Searxng**:
   - Visit: `https://search.pochita.synology.me`

2. **Configure Preferences** (top-right menu → Preferences):
   - **General**:
     - Language: English
     - Autocomplete: Off (for privacy)
   - **Privacy**:
     - Method: POST
     - Image proxy: Enabled
   - **Search**:
     - Default categories: General
     - Language: Auto-detect

3. **Set as Default Search Engine**:
   - **Firefox**:
     1. Visit search.pochita.synology.me
     2. Right-click address bar → "Add The Club Search"
     3. Settings → Search → Default Search Engine → The Club Search

   - **Chrome/Edge**:
     1. Settings → Search Engine → Manage
     2. Add: `search.pochita.synology.me`, shortcut: `@club`
     3. Set as default

---

## Configure Devices to Use Pi-hole

**CRITICAL**: You must configure your devices to use Pi-hole as their DNS server for network-wide blocking to work.

### Option 1: Router Configuration (Recommended)
**Benefit**: All devices automatically protected

1. Access your router admin page
2. Find DHCP settings (usually under "LAN" or "Network")
3. Set Primary DNS to your Club server IP (e.g., `192.168.1.100`)
4. Set Secondary DNS to `1.1.1.1` (Cloudflare, fallback only)
5. Save and reboot router
6. Restart devices to get new DNS settings

### Option 2: Manual Device Configuration
**Benefit**: Selective protection (e.g., work devices keep default DNS)

**Windows**:
1. Settings → Network & Internet → Ethernet/WiFi → Properties
2. IP settings → Edit → Manual
3. IPv4: On
4. Preferred DNS: `192.168.1.100` (your Club server)
5. Alternate DNS: `1.1.1.1`
6. Save

**macOS**:
1. System Preferences → Network → Advanced
2. DNS tab → Add (+)
3. Add: `192.168.1.100`
4. Add: `1.1.1.1` (fallback)
5. OK → Apply

**Linux**:
```bash
# Edit /etc/resolv.conf
sudo nano /etc/resolv.conf

# Add:
nameserver 192.168.1.100
nameserver 1.1.1.1

# Or use NetworkManager
nmcli con mod "Your Connection" ipv4.dns "192.168.1.100 1.1.1.1"
nmcli con down "Your Connection" && nmcli con up "Your Connection"
```

**iPhone/iPad**:
1. Settings → WiFi → (i) next to your network
2. Configure DNS → Manual
3. Remove existing servers
4. Add Server: `192.168.1.100`
5. Add Server: `1.1.1.1`
6. Save

**Android**:
1. Settings → WiFi → Long press network → Modify
2. Advanced options → IP settings: Static
3. DNS 1: `192.168.1.100`
4. DNS 2: `1.1.1.1`
5. Save

### Option 3: DNS Profile (iOS/macOS)
**Benefit**: Easy deployment across multiple Apple devices

Use a DNS profile generator or manually create one for Pi-hole.

---

## Verification

### Test Pi-hole Blocking

```bash
# Test from a device using Pi-hole DNS

# This should be blocked
nslookup doubleclick.net

# Should return 0.0.0.0 or your Pi-hole IP

# This should work
nslookup google.com

# Should return Google's IP
```

Or visit Pi-hole dashboard → Query Log to see blocked domains in real-time.

### Test Search Engine

1. Visit: `https://search.pochita.synology.me`
2. Search for: "privacy tools"
3. Check that results appear
4. Verify no tracking (browser console should show no external requests)

### Test DNS Privacy

```bash
# Check which DNS server is being used
nslookup -type=txt whoami.akamai.net

# Should show your Club server IP, not ISP's DNS
```

---

## Privacy Dashboard

View your privacy statistics:

### Pi-hole Dashboard
- **Access**: `https://dns.pochita.synology.me/admin`
- **Shows**:
  - Queries blocked (%)
  - Top blocked domains
  - Query types
  - Client activity

### Grafana Integration (Optional)

Add Pi-hole metrics to Grafana:

1. Install Pi-hole Exporter:
   ```bash
   docker run -d \
     --name pihole-exporter \
     -e PIHOLE_HOSTNAME=pihole \
     -e PIHOLE_PASSWORD=your_password \
     --network theclub_club-network \
     ekofr/pihole-exporter
   ```

2. Update Prometheus config (`monitoring/prometheus/prometheus.yml`):
   ```yaml
   - job_name: 'pihole'
     static_configs:
       - targets: ['pihole-exporter:9617']
         labels:
           service: 'dns'
   ```

3. Import Pi-hole dashboard in Grafana (ID: 10176)

---

## Privacy Best Practices

### 1. Regular Blocklist Updates
```bash
# Update Pi-hole blocklists weekly
docker exec pihole pihole -g
```

Or configure cron:
```bash
# Weekly updates every Sunday at 3 AM
0 3 * * 0 docker exec pihole pihole -g
```

### 2. Review Blocked Queries

Occasionally check Pi-hole Query Log for:
- False positives (legitimate sites blocked)
- New tracking domains (add to blocklists)
- Suspicious activity

### 3. Whitelist Necessary Domains

Some services may break with aggressive blocking:

```bash
# Whitelist a domain
docker exec pihole pihole -w domain.com

# Or via web interface
# Pi-hole → Whitelist → Add domain
```

**Common whitelists needed**:
- `s.youtube.com` (YouTube history)
- `www.googleadservices.com` (Google Ads - if you need them for work)
- Smart TV update domains (varies by manufacturer)

### 4. Disable Query Logging (Maximum Privacy)

Edit `docker-compose.privacy.yml`:
```yaml
- QUERY_LOGGING=false  # No query logs (more private, loses statistics)
```

Restart Pi-hole:
```bash
docker restart pihole
```

### 5. Use Searxng for All Searches

Replace browser default search:
- Firefox: Settings → Search → Default Search Engine
- Chrome: Settings → Search engine → Manage search engines

Set to: `https://search.pochita.synology.me`

### 6. Configure Browser Privacy Settings

Even with privacy infrastructure, configure browsers:

**Firefox**:
- Settings → Privacy & Security
- Enhanced Tracking Protection: Strict
- Do Not Track: Always
- HTTPS-Only Mode: Enable

**Chrome/Brave**:
- Settings → Privacy and security
- Cookies: Block third-party
- Send "Do Not Track": On

---

## Advanced Privacy Enhancements

### 1. DNS-over-HTTPS (DoH) Blocking

Prevent apps from bypassing Pi-hole:

Edit Pi-hole dnsmasq config:
```bash
docker exec -it pihole bash
echo "address=/dns.google/0.0.0.0" >> /etc/dnsmasq.d/02-custom.conf
echo "address=/cloudflare-dns.com/0.0.0.0" >> /etc/dnsmasq.d/02-custom.conf
exit

docker restart pihole
```

### 2. CNAME Uncloaking

Block first-party tracking (e.g., Facebook Pixel on websites):

Pi-hole Settings → DNS → Advanced DNS settings:
- Enable "Never forward reverse lookups for private IP ranges"
- Enable "Use DNSSEC"

### 3. Conditional Forwarding

For local network domain resolution:

Pi-hole Settings → DNS → Advanced DNS settings:
- Enable "Conditional Forwarding"
- Local network: `192.168.1.0/24`
- Router IP: `192.168.1.1`
- Domain: `lan`

### 4. Custom Block Lists

Create custom blocklists for specific threats:

```bash
# Create custom list
docker exec -it pihole bash
nano /etc/pihole/custom.list

# Add domains (one per line)
facebook.com
facebook.net
fbcdn.net
# etc...

pihole restartdns
exit
```

---

## Troubleshooting

### Pi-hole Not Blocking

**Check DNS configuration**:
```bash
# From client device
nslookup pi.hole

# Should resolve to your Pi-hole IP
```

**If not working**:
1. Verify device is using Pi-hole DNS
2. Clear DNS cache: `ipconfig /flushdns` (Windows) or `sudo dscacheutil -flushcache` (macOS)
3. Restart network connection

### Websites Breaking

**Symptoms**: Sites not loading, videos not playing, logins failing

**Solution**: Check Pi-hole Query Log
1. Find blocked domain causing issue
2. Whitelist it:
   ```bash
   docker exec pihole pihole -w problematic-domain.com
   ```

### Slow DNS Resolution

**Cause**: First queries to new domains are slower (recursive resolution)

**Solutions**:
1. Increase cache size in Unbound config
2. Enable prefetching (already enabled in our config)
3. Wait - subsequent queries will be fast (cached)

### Searxng Not Returning Results

**Check**:
1. Searxng logs: `docker logs searxng`
2. Redis connectivity: `docker logs searxng-redis`
3. Search engine settings in `privacy/searxng/settings.yml`

**Fix**: Enable more search engines if some are down

### FreshRSS Feeds Not Updating

**Check**:
1. FreshRSS logs: `docker logs freshrss`
2. Feed update cron: Should run every 15 minutes

**Manual update**:
```bash
docker exec freshrss php ./app/actualize_script.php
```

---

## Privacy Impact Metrics

After deployment, you should see:

### Network-Wide
- **50-80% of DNS queries blocked** (ads/trackers)
- **Faster page loads** (fewer resources to download)
- **Reduced bandwidth** (30-40% reduction typical)

### Search Privacy
- **Zero Google profiling** (searches don't go to Google directly)
- **No search history logging** (Searxng doesn't log)
- **Anonymous results** (aggregated from multiple sources)

### Data Privacy
- **No ISP DNS logging** (queries stay local with Unbound)
- **No third-party tracking** (blocked by Pi-hole)
- **Private RSS reading** (FreshRSS doesn't expose habits)

---

## Maintenance

### Daily
- None (automated)

### Weekly
- Review Pi-hole statistics
- Check for false positives in Query Log

### Monthly
- Update Pi-hole blocklists: `docker exec pihole pihole -g`
- Update containers: `docker compose -f docker-compose.privacy.yml pull`
- Review Searxng search engine status

### Quarterly
- Review and update blocklists (add new sources)
- Check for Unbound updates
- Audit privacy settings

---

## Cost-Benefit Analysis

### Resources Used
- **RAM**: ~1.5GB total (Pi-hole: 512MB, Unbound: 256MB, Searxng: 512MB, FreshRSS: 256MB)
- **CPU**: <5% average
- **Storage**: ~2GB

### Privacy Benefits
- Network-wide ad/tracker blocking
- Private DNS (no third-party logging)
- Private search (no profiling)
- Private news reading

**Verdict**: Minimal resource cost for maximum privacy benefit

---

## Comparison with Alternatives

### Pi-hole vs. Browser Ad Blockers
| Feature | Pi-hole | uBlock Origin |
|---------|---------|---------------|
| Blocks ads in apps | ✅ Yes | ❌ No (browser only) |
| Blocks ads on all devices | ✅ Yes | ❌ No (per-device) |
| Blocks DNS tracking | ✅ Yes | ❌ No |
| Reduces bandwidth | ✅ Yes | ⚠️ Partial |
| Statistics/monitoring | ✅ Yes | ⚠️ Limited |

**Recommendation**: Use BOTH (Pi-hole + browser ad blocker)

### Unbound vs. Cloudflare DNS
| Feature | Unbound | Cloudflare |
|---------|---------|------------|
| Privacy (no logging) | ✅ Yes | ⚠️ Claims no logging |
| Data stays local | ✅ Yes | ❌ No |
| No third-party | ✅ Yes | ❌ No (Cloudflare sees queries) |
| DNSSEC | ✅ Yes | ✅ Yes |
| Speed (first query) | ⚠️ Slower | ✅ Faster |
| Speed (cached) | ✅ Fast | ✅ Fast |

**Recommendation**: Unbound for maximum privacy, Cloudflare if speed is priority

### Searxng vs. DuckDuckGo
| Feature | Searxng | DuckDuckGo |
|---------|---------|------------|
| Privacy | ✅ Self-hosted | ⚠️ Trust-based |
| No tracking | ✅ Guaranteed | ⚠️ Their policy |
| Multiple sources | ✅ Yes | ❌ No (own index) |
| Self-hosted | ✅ Yes | ❌ No |
| Convenience | ⚠️ Setup required | ✅ Just works |

**Recommendation**: Searxng for maximum privacy, DuckDuckGo for convenience

---

## Next Steps

1. **Deploy privacy stack**: Follow deployment section above
2. **Configure devices**: Point DNS to Pi-hole
3. **Set default search**: Use Searxng for all searches
4. **Add RSS feeds**: Set up FreshRSS with your favorite sources
5. **Monitor privacy**: Check Pi-hole dashboard regularly
6. **Educate users**: Tell family/roommates about privacy benefits

---

## Additional Resources

- **Pi-hole Documentation**: https://docs.pi-hole.net/
- **Unbound Documentation**: https://unbound.docs.nlnetlabs.nl/
- **Searxng Documentation**: https://docs.searxng.org/
- **FreshRSS Documentation**: https://freshrss.github.io/FreshRSS/
- **Privacy Tools**: https://www.privacytools.io/
- **EFF Privacy Guides**: https://ssd.eff.org/

---

## Congratulations! 🎉

You now have a privacy-focused infrastructure that:
- Blocks ads and trackers network-wide
- Keeps your DNS queries private
- Allows private searching without profiling
- Enables private news reading

All while being:
- 100% FOSS (Free and Open Source Software)
- Self-hosted (you control the data)
- Transparent (you can audit the code)
- Privacy-first (designed with privacy as core principle)

**Your network is now significantly more private than 99% of internet users!**
