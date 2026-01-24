# Distributed Architecture Guide

## Overview

This guide explains how to distribute The Club services across multiple devices to optimize resource usage and performance.

**Available Hardware**:
- **N97 PC (The Club Server)**: Main server, Docker host
- **Raspberry Pi 5 + Hailo AI Hat**: Frigate object detection (already allocated)
- **Raspberry Pi 4**: Available for service offloading
- **Synology DS415+ NAS**: Network storage, specific services

---

## Service Distribution Strategy

### Device Capabilities

#### N97 PC
- **CPU**: Intel N97 (4-core, 3.6 GHz)
- **RAM**: 8GB+ recommended
- **Storage**: SSD/HDD
- **Best For**: CPU-intensive tasks, video processing, heavy I/O
- **Keep**: Services requiring high CPU, storage, or integration

#### Raspberry Pi 5 (8GB)
- **CPU**: ARM Cortex-A76 (4-core, 2.4 GHz)
- **RAM**: 8GB
- **Storage**: microSD/NVMe (via HAT)
- **Hailo AI Hat**: AI acceleration (13 TOPS)
- **Current**: Frigate object detection
- **Can Also Run**: Additional lightweight services

#### Raspberry Pi 4
- **CPU**: ARM Cortex-A72 (4-core, 1.5-1.8 GHz)
- **RAM**: Varies (1GB/2GB/4GB/8GB models)
- **Storage**: microSD/USB SSD
- **Network**: Gigabit Ethernet, USB 3.0
- **Best For**: Lightweight services, network services, always-on tasks
- **Available**: Fully available for offloading

#### Raspberry Pi 3 Model B/B+
- **CPU**: ARM Cortex-A53 (4-core, 1.2-1.4 GHz)
- **RAM**: 1GB LPDDR2
- **Storage**: microSD only
- **Network**: 100 Mbps Ethernet (shared with USB 2.0 bus)
- **Released**: 2016-2018
- **Best For**: Single lightweight service, DNS, basic monitoring
- **Limitations**: Limited RAM (1GB), slower network, USB 2.0 bottleneck
- **Available**: Can be used for specific services (see below)

#### Raspberry Pi Zero 2 W
- **CPU**: ARM Cortex-A53 (4-core, 1.0 GHz) - same as Pi 3 but slower
- **RAM**: 512MB LPDDR2 (half of Pi 3!)
- **Storage**: microSD only
- **Network**: WiFi 802.11n (2.4 GHz), Bluetooth 4.2 - **NO Ethernet port**
- **Power**: ~0.4W idle, ~1W active (most power-efficient)
- **Released**: 2021
- **Price**: ~$15 new
- **Best For**: Ultra-lightweight single service, testing, non-critical tasks
- **Limitations**: Only 512MB RAM (very limiting), WiFi-only (latency/reliability), slower CPU
- **Not Recommended For**: Critical services (DNS, VPN), production use
- **Available**: Can be used for experimental/testing purposes only

#### Libre Computer Le Potato (AML-S905X-CC)
- **CPU**: Amlogic S905X (4-core ARM Cortex-A53, 1.5 GHz)
- **RAM**: 2GB DDR3-2133 (same as Pi 4 2GB model)
- **GPU**: ARM Mali-450 with **hardware 4K video decode** (H.265, VP9, H.264)
- **Storage**: microSD + eMMC module support (more reliable than microSD)
- **Network**: **100 Mbps Ethernet** (NOT Gigabit) - No built-in WiFi
- **USB**: USB 2.0 ports
- **Power**: ~2-3W typical (more efficient than Pi 4)
- **Released**: 2017 (still sold in 2026)
- **Price**: ~$35 new (similar to Pi 4 2GB)
- **Best For**: Media applications (4K hardware decode), Pi-hole/DNS, low-power server
- **Advantages**: Better video transcoding than Pi 3/4, eMMC support, faster than Pi 3 B+
- **Limitations**: 100 Mbps Ethernet (vs Pi 4's Gigabit), no WiFi, smaller community than Pi
- **Performance**: 50-60% of Pi 4 CPU performance, but uses half the power
- **Available**: Good alternative to Raspberry Pi for specific use cases

#### Synology DS415+ NAS
- **CPU**: Marvell Armada XP (ARM dual-core, 1.33 GHz)
- **RAM**: 1GB (expandable to 2GB)
- **Storage**: 4-bay NAS (large storage)
- **Current**: Jellyfin, Vaultwarden
- **Best For**: Storage-dependent services, media

---

## Recommended Service Distribution

### Tier 1: Perfect for Raspberry Pi 4 (Offload These)

#### 🥇 **Pi-hole + Unbound** (HIGHLY RECOMMENDED)
**Current Location**: N97 PC
**Move To**: Raspberry Pi 4

**Why**:
- ✅ Pi-hole is literally designed for Raspberry Pi
- ✅ Extremely lightweight (~100MB RAM total)
- ✅ Network DNS should be on dedicated hardware (resilience)
- ✅ Won't impact other services if issues occur
- ✅ Always-on service (Pi consumes less power)
- ✅ Frees up ~150MB RAM on N97

**Resource Usage**:
- **RAM**: 50-100MB (Pi-hole) + 20-50MB (Unbound)
- **CPU**: <5% on Pi 4
- **Network**: DNS queries only (minimal bandwidth)

**How to Move**: See "Migration: Pi-hole + Unbound to RPi4" section below

---

#### 🥇 **WireGuard VPN** (HIGHLY RECOMMENDED)
**Current Location**: N97 PC
**Move To**: Raspberry Pi 4

**Why**:
- ✅ Extremely lightweight kernel module
- ✅ Perfect for always-on VPN gateway
- ✅ Low latency, minimal CPU usage
- ✅ Separates VPN from main server (security)
- ✅ Can route through Pi-hole for ad-blocking over VPN
- ✅ Frees up VPN overhead on N97

**Resource Usage**:
- **RAM**: 20-50MB
- **CPU**: <5% (even during active VPN connections)
- **Network**: VPN throughput (Pi 4 can handle 100+ Mbps)

**Network Architecture**:
```
Internet → Router → Raspberry Pi 4 (WireGuard) → Internal LAN
                         ↓
                    Pi-hole DNS
```

---

#### 🥇 **Uptime Kuma** (RECOMMENDED)
**Current Location**: N97 PC
**Move To**: Raspberry Pi 4

**Why**:
- ✅ Lightweight monitoring tool
- ✅ Benefits from being on separate hardware (monitors N97)
- ✅ If N97 crashes, monitoring still works
- ✅ Can monitor Pi-hole, WireGuard, and all N97 services
- ✅ Frees up ~100-150MB RAM on N97

**Resource Usage**:
- **RAM**: 100-150MB
- **CPU**: <10% (periodic health checks)
- **Network**: HTTP(S) requests to monitored services

---

#### 🥈 **Homer Dashboard** (OPTIONAL)
**Current Location**: N97 PC
**Move To**: Raspberry Pi 4 OR keep on N97

**Why**:
- ✅ Static site, extremely lightweight
- ✅ Benefits from being always-on
- ✅ Frees up ~30-50MB RAM on N97
- ⚠️ But requires reverse proxy configuration

**Resource Usage**:
- **RAM**: 30-50MB
- **CPU**: <1%

**Recommendation**: Keep on N97 unless Pi 4 becomes central entry point

---

#### 🥈 **MQTT Broker (Mosquitto)** (RECOMMENDED for Frigate)
**Current Location**: N97 PC (with Frigate)
**Move To**: Raspberry Pi 4 OR Raspberry Pi 5

**Why**:
- ✅ Extremely lightweight
- ✅ Central message broker for IoT/Frigate
- ✅ Can run on either Pi
- ✅ Frees up ~50-100MB RAM on N97

**Resource Usage**:
- **RAM**: 50-100MB
- **CPU**: <5%

**Recommendation**: Run on RPi4 if setting up as IoT hub, or RPi5 alongside Frigate detector

---

### Tier 2: Could Run on Raspberry Pi (Moderate)

#### 🟡 **FreshRSS** (OPTIONAL)
**Current Location**: N97 PC
**Move To**: Raspberry Pi 4 (if 4GB+ RAM)

**Why**:
- ✅ PHP-based, relatively lightweight
- ⚠️ Requires database (SQLite or PostgreSQL)
- ⚠️ RSS fetching can be CPU-intensive with many feeds

**Resource Usage**:
- **RAM**: 150-300MB (depends on feed count)
- **CPU**: 10-20% (during feed updates)

**Recommendation**: Move if Pi 4 has 4GB+ RAM

---

#### 🟡 **Authelia + Redis** (NOT RECOMMENDED)
**Current Location**: N97 PC
**Keep On**: N97 PC

**Why Keep on N97**:
- ⚠️ Central authentication service (critical)
- ⚠️ Benefits from being on same host as Caddy (low latency)
- ⚠️ Session storage (Redis) should be fast
- ⚠️ High availability important

**Resource Usage**:
- **RAM**: 150-250MB (Authelia + Redis)
- **CPU**: 5-10%

**Recommendation**: Keep on N97 for performance and reliability

---

### Tier 3: Keep on N97 (Resource-Intensive)

#### ❌ **Frigate NVR (Main Instance)**
**Keep On**: N97 PC

**Why**:
- ❌ Video processing requires CPU power
- ❌ Recording storage (SSD I/O)
- ❌ Multiple camera streams
- ❌ Web interface serving video

**Resource Usage**:
- **RAM**: 2-4GB (with 3+ cameras)
- **CPU**: 20-50% (depends on camera count)
- **Storage**: High I/O (continuous recording)

---

#### ❌ **Prometheus + Grafana + cAdvisor + Node Exporter**
**Keep On**: N97 PC

**Why**:
- ❌ Prometheus time-series database (storage-intensive)
- ❌ Grafana rendering dashboards (CPU-intensive)
- ❌ cAdvisor monitoring Docker (needs access to Docker)
- ❌ Better performance on more powerful hardware

**Resource Usage**:
- **RAM**: 500MB-1GB (Prometheus DB grows over time)
- **CPU**: 10-20%
- **Storage**: Time-series data (SSD preferred)

**Recommendation**: Keep on N97, adjust retention if needed

---

#### ❌ **Kopia Backup**
**Keep On**: N97 PC

**Why**:
- ❌ Backup operations are CPU/I/O intensive
- ❌ Deduplication requires RAM and CPU
- ❌ Compression is CPU-heavy
- ❌ Needs access to data being backed up

**Resource Usage**:
- **RAM**: 200-500MB (during backup)
- **CPU**: 20-50% (during backup operations)
- **I/O**: High disk I/O

---

#### ❌ **Caddy Reverse Proxy** (Keep on N97)
**Keep On**: N97 PC

**Why Keep on N97**:
- ⚠️ Central entry point for all services
- ⚠️ SSL termination (Let's Encrypt)
- ⚠️ Most services are on N97 (low latency)
- ⚠️ Routing to multiple backends

**Alternative Architecture**:
If you want distributed setup, you could:
- Run Caddy on N97 (proxies to N97 services + Pi services)
- Or run Nginx/Caddy on Pi 4 as main entry point (proxies to N97 + other Pis)

**Recommendation**: Keep on N97 for simplicity

---

#### ❌ **Searxng** (Keep on N97)
**Keep On**: N97 PC

**Why**:
- ⚠️ Search aggregation is CPU-intensive
- ⚠️ Fetches from multiple search engines simultaneously
- ⚠️ Benefits from faster CPU

**Resource Usage**:
- **RAM**: 150-300MB
- **CPU**: 10-30% (during searches)

**Recommendation**: Keep on N97 unless you have RPi 4 with 8GB RAM

---

#### ❌ **Dockge** (Keep on N97)
**Keep On**: N97 PC

**Why**:
- ❌ Manages Docker Compose stacks on N97
- ❌ Needs access to Docker socket
- ❌ Should be on same host it's managing

---

#### ❌ **Watchtower** (Keep on N97)
**Keep On**: N97 PC

**Why**:
- ❌ Updates Docker containers on N97
- ❌ Needs access to Docker socket
- ❌ Should be on same host it's managing

**Note**: You would need separate Watchtower instances on each Pi

---

## Recommended Distribution Plan

### Option A: Conservative (Recommended)

Move only the most obvious candidates to Pi 4:

**Raspberry Pi 4**:
- ✅ Pi-hole + Unbound (DNS)
- ✅ WireGuard VPN
- ✅ Uptime Kuma (monitoring)
- ✅ MQTT Broker (for Frigate/IoT)

**Raspberry Pi 5**:
- ✅ Frigate Detector (Hailo AI)
- ✅ Optional: Additional lightweight services if needed

**N97 PC**:
- ✅ Caddy (reverse proxy)
- ✅ Homer (dashboard)
- ✅ Authelia + Redis (SSO)
- ✅ Frigate NVR (main instance)
- ✅ Prometheus + Grafana (metrics)
- ✅ Kopia + Watchtower + Dockge (management)
- ✅ Searxng + FreshRSS (privacy apps)
- ✅ Go/HTMX test server

**Synology NAS**:
- ✅ Vaultwarden (passwords)
- ✅ Jellyfin (media)

**Savings on N97**: ~400-500MB RAM, reduced CPU load (DNS, VPN, monitoring offloaded)

---

### Option B: Aggressive (Maximum Offload)

Move everything possible to Pi 4:

**Raspberry Pi 4** (requires 4GB+ RAM model):
- ✅ Pi-hole + Unbound (DNS)
- ✅ WireGuard VPN
- ✅ Uptime Kuma (monitoring)
- ✅ MQTT Broker
- ✅ FreshRSS (RSS reader)
- ✅ Homer (dashboard)
- ⚠️ Optional: Caddy (as main reverse proxy)

**Raspberry Pi 5**:
- ✅ Frigate Detector (Hailo AI)

**N97 PC**:
- ✅ Frigate NVR (main instance)
- ✅ Authelia + Redis
- ✅ Prometheus + Grafana
- ✅ Kopia + Watchtower + Dockge
- ✅ Searxng
- ✅ Go/HTMX test server
- ⚠️ Optional: Caddy (if keeping as reverse proxy)

**Synology NAS**:
- ✅ Vaultwarden
- ✅ Jellyfin

**Savings on N97**: ~700-800MB RAM, significant CPU reduction

---

### Option C: Raspberry Pi 3 Model B/B+ (Limited Hardware)

If you only have a Pi 3 available, here's what it can realistically handle:

**Best Choice: Single Service Deployment**

Due to 1GB RAM limitation, run **ONE** of these services:

#### Option C1: Dedicated DNS Server (RECOMMENDED)
**Raspberry Pi 3**:
- ✅ Pi-hole + Unbound (DNS only)

**Why This Works**:
- ✅ Pi-hole was originally designed for Pi 3 (and even Pi Zero)
- ✅ Combined RAM usage: ~120-150MB (well within 1GB)
- ✅ DNS queries are not CPU-intensive
- ✅ Network: 100 Mbps is plenty for DNS (queries are tiny)
- ✅ Most reliable use case for Pi 3

**Resource Usage**:
- **RAM**: 120-150MB / 1GB (leaves 850MB buffer)
- **CPU**: 5-10% average
- **Network**: <1 Mbps (DNS queries)

**Verdict**: ✅ **EXCELLENT** - This is what many people use Pi 3 for

---

#### Option C2: Dedicated VPN Server (ACCEPTABLE)
**Raspberry Pi 3**:
- ✅ WireGuard VPN only

**Why This Works**:
- ✅ WireGuard is very lightweight
- ✅ RAM usage: ~50-80MB
- ⚠️ Network limited to ~50-80 Mbps (USB 2.0 + 100 Mbps Ethernet bottleneck)
- ⚠️ Slower than Pi 4, but functional

**Resource Usage**:
- **RAM**: 50-80MB / 1GB
- **CPU**: 15-25% under load
- **Network**: Limited to ~50-80 Mbps throughput

**Verdict**: ⚠️ **ACCEPTABLE** - Works, but limited by 100 Mbps Ethernet

---

#### Option C3: Dedicated Monitoring (MARGINAL)
**Raspberry Pi 3**:
- ⚠️ Uptime Kuma only

**Why This Is Marginal**:
- ⚠️ Uptime Kuma: ~150-200MB RAM (tight on 1GB)
- ⚠️ Can become slow with many monitors (>20)
- ⚠️ Works, but not ideal

**Resource Usage**:
- **RAM**: 150-200MB / 1GB
- **CPU**: 10-15%

**Verdict**: ⚠️ **MARGINAL** - Better on Pi 4, but works

---

#### What Pi 3 CANNOT Handle Well

❌ **Do NOT Run on Pi 3**:
- ❌ Multiple services simultaneously (1GB RAM too limiting)
- ❌ FreshRSS (needs 200-300MB + database)
- ❌ Homer + other services (wastes resources)
- ❌ Any video processing
- ❌ Database-heavy applications
- ❌ Prometheus/Grafana

---

### Raspberry Pi 3 vs Pi 4: Comparison

| Capability | Pi 3 Model B/B+ | Pi 4 (2GB) | Pi 4 (4GB+) |
|------------|----------------|------------|-------------|
| **Pi-hole + Unbound** | ✅ Excellent | ✅ Excellent | ✅ Excellent |
| **WireGuard VPN** | ⚠️ 50-80 Mbps | ✅ 100-200 Mbps | ✅ 200+ Mbps |
| **Uptime Kuma** | ⚠️ Single service only | ✅ With others | ✅ With others |
| **MQTT Broker** | ✅ Alone, marginal with others | ✅ With others | ✅ With others |
| **Homer Dashboard** | ⚠️ Alone only | ✅ With others | ✅ With others |
| **Multiple Services** | ❌ Not recommended | ⚠️ 2-3 services | ✅ 4-5 services |
| **FreshRSS** | ❌ Too limited | ⚠️ Marginal | ✅ Works well |
| **Power Draw** | 2-3W | 3-6W | 4-8W |
| **Cost (Used)** | $15-25 | $35-50 | $45-60 |

---

### Recommended: Pi 3 as Dedicated DNS Server

**Best Use Case for Raspberry Pi 3**:

```yaml
# Raspberry Pi 3: DNS Server Only
services:
  unbound:
    image: mvance/unbound-rpi:latest
    # ... (same config as Pi 4 guide)

  pihole:
    image: pihole/pihole:latest
    # ... (same config as Pi 4 guide)
```

**Configuration**:
- **Static IP**: `192.168.1.39` (or any available IP)
- **Services**: Pi-hole + Unbound only
- **Memory**: 150MB / 1GB used (~85% free)
- **Router DNS**: Point to `192.168.1.39`

**Benefits**:
- ✅ Dedicated DNS on separate hardware (resilient)
- ✅ Low power consumption (2-3W always-on)
- ✅ Frees up N97 resources
- ✅ Pi 3 runs cool and stable with just DNS
- ✅ Perfect use case for older hardware

**If You Have Both Pi 3 and Pi 4**:

**Raspberry Pi 3**:
- ✅ Pi-hole + Unbound (DNS)
- Static IP: `192.168.1.39`

**Raspberry Pi 4**:
- ✅ WireGuard VPN
- ✅ Uptime Kuma
- ✅ MQTT Broker
- ✅ Homer Dashboard (optional)
- Static IP: `192.168.1.40`

This gives you maximum offload with both devices!

---

### Pi 3 Setup Optimizations

To get the best performance from Pi 3:

#### 1. Use USB SSD Boot (Optional but Recommended)

**Why**:
- microSD cards are slow and wear out
- USB SSD provides faster I/O
- More reliable long-term

**How**:
```bash
# Enable USB boot on Pi 3 B+ (newer firmware)
# Flash Raspberry Pi OS to USB SSD
# Boot from USB instead of microSD
```

Note: Only Pi 3 B+ supports USB boot. Pi 3 Model B needs microSD.

#### 2. Disable Unnecessary Services

```bash
# Disable Bluetooth (saves RAM)
sudo systemctl disable bluetooth
sudo systemctl disable hciuart

# Disable WiFi if using Ethernet
sudo rfkill block wifi

# Reduce GPU memory (headless)
sudo raspi-config
# Advanced Options → Memory Split → Set to 16MB
```

#### 3. Use Lightweight Docker Images

For Pi 3, use ARM32v7 images when available:
- `pihole/pihole:latest` - Official, ARM-compatible
- `mvance/unbound-rpi:latest` - Optimized for Raspberry Pi

#### 4. Monitor Resources

```bash
# Check memory
free -h

# Check temperature
vcgencmd measure_temp

# Monitor Docker containers
docker stats
```

**Safe Operating Ranges for Pi 3**:
- **RAM**: Keep usage below 800MB (leave 200MB free)
- **Temp**: Keep below 70°C (add heatsink if needed)
- **CPU**: <50% average is healthy

---

### Raspberry Pi Zero 2 W: Ultra-Limited Use Cases

The Pi Zero 2 W is significantly more limited than even the Pi 3.

#### Hardware Constraints

**Critical Limitations**:
- ❌ **Only 512MB RAM** (half of Pi 3)
- ❌ **WiFi only** - No Ethernet port (less reliable, higher latency)
- ❌ **Slower CPU** (1.0 GHz vs Pi 3's 1.4 GHz)
- ❌ **WiFi bandwidth** limited (~40 Mbps real-world)

#### What It CAN Handle (Barely)

**Pi-hole ONLY** (without Unbound):
- ⚠️ Pi-hole alone: ~80-100MB RAM
- ⚠️ **Cannot** run Unbound alongside (would exceed 512MB)
- ⚠️ Must use external DNS upstream (Google, Cloudflare)
- ⚠️ WiFi adds latency to DNS queries (~5-10ms extra)
- ⚠️ Less reliable than Ethernet-based DNS

**Resource Usage**:
- **RAM**: 80-120MB / 512MB (leaves 400MB buffer, but...)
- **CPU**: 10-15%
- **Network**: WiFi 2.4 GHz only

**Verdict**: ⚠️ **NOT RECOMMENDED** for production DNS

---

#### What Pi Zero 2 W CANNOT Handle

❌ **Do NOT use Pi Zero 2 W for**:
- ❌ Pi-hole + Unbound (exceeds 512MB RAM)
- ❌ VPN (WiFi bottleneck, unreliable)
- ❌ Uptime Kuma (needs 150-200MB, too tight)
- ❌ Any critical network service (WiFi unreliable)
- ❌ Multiple services (RAM too limited)
- ❌ Production use (too constrained)

---

#### When Pi Zero 2 W Makes Sense

**✅ Good Use Cases**:
1. **Testing/Development**: Test Docker configurations before deploying to real hardware
2. **IoT Sensors**: Temperature monitoring, motion sensors (non-critical)
3. **Display Dashboard**: Kiosk mode showing Homer dashboard on a screen
4. **Learning**: Experiment with Docker/Linux on cheap hardware
5. **Backup Pi-hole**: Secondary DNS server (failover only)

**Power Advantage**:
- **0.4-1W** power consumption (lowest of all Pis)
- Perfect for battery-powered or solar setups
- Runs very cool (rarely needs heatsink)

---

#### Pi Zero 2 W vs Other Pis

| Feature | Pi Zero 2 W | Pi 3 B+ | Pi 4 (4GB) |
|---------|-------------|---------|------------|
| **RAM** | 512MB ❌ | 1GB ⚠️ | 4GB ✅ |
| **Network** | WiFi only ❌ | 100 Mbps Ethernet ✅ | Gigabit Ethernet ✅ |
| **Pi-hole + Unbound** | ❌ No (RAM) | ✅ Yes | ✅ Yes |
| **Pi-hole only** | ⚠️ Marginal | ✅ Yes | ✅ Yes |
| **VPN** | ❌ No (WiFi) | ⚠️ 50-80 Mbps | ✅ 100-200 Mbps |
| **Multiple Services** | ❌ No | ❌ No | ✅ Yes |
| **Power Draw** | 0.4-1W ⭐ | 2-3W | 4-6W |
| **Reliability** | ⚠️ WiFi | ✅ Wired | ✅ Wired |
| **Price** | $15 | $15-25 | $45-60 |

---

#### Recommendation for Pi Zero 2 W

**If You Have a Pi Zero 2 W**:
- ✅ Use it for: Testing, IoT projects, display kiosk
- ❌ Don't use it for: DNS, VPN, or any critical service
- 💡 Better to get a Pi 3 or Pi 4 for homelab services

**Why Not for Homelab?**:
1. 512MB RAM is too limiting for most services
2. WiFi-only is unreliable for critical network services (DNS, VPN)
3. No Ethernet means higher latency and potential dropouts
4. Cannot run Pi-hole + Unbound together (the recommended setup)

**Verdict**: ❌ **Skip Pi Zero 2 W for The Club homelab** - Use Pi 3 or Pi 4 instead

---

### Libre Computer Le Potato (AML-S905X-CC): Raspberry Pi Alternative

The Le Potato is a Raspberry Pi alternative with some unique strengths, particularly for media applications.

#### Hardware Strengths

**Superior Video Capabilities**:
- ✅ **Hardware 4K video decode** (H.265/HEVC, VP9, H.264)
- ✅ ARM Mali-450 GPU with OpenGL and OpenVG
- ✅ Better media transcoding than Pi 3/4 (hardware accelerated)
- ✅ Ideal for media servers (Jellyfin, Plex, Emby)

**Other Advantages**:
- ✅ **2GB DDR3 RAM** (same as Pi 4 2GB model)
- ✅ **eMMC module support** (more reliable than microSD for 24/7 operation)
- ✅ **Faster than Pi 3 B+** with half the power consumption
- ✅ **Lower price**: ~$35 (competitive with Pi 4 2GB)
- ✅ **Compatible with Docker** and Docker Compose

---

#### Hardware Limitations

**Critical Limitations**:
- ❌ **Only 100 Mbps Ethernet** (not Gigabit like Pi 4)
- ❌ **No built-in WiFi** (requires USB WiFi dongle)
- ❌ **50-60% slower CPU** than Raspberry Pi 4
- ❌ **Smaller community** and less software support than Raspberry Pi
- ❌ **Slower microSD I/O** (15-20 MB/s sequential, 2.5-4.5 MB/s random)

---

#### What Le Potato CAN Handle

**✅ Good Use Cases for The Club**:

1. **Pi-hole + Unbound (DNS)** - Excellent
   - RAM: 120-150MB / 2GB (plenty of headroom)
   - CPU: 5-10% usage
   - Network: 100 Mbps plenty for DNS queries
   - eMMC support = very reliable for 24/7 DNS
   - **Verdict**: ✅ Better than Pi 3, similar to Pi 4 for this use case

2. **Media Transcoding** (if needed)
   - Hardware 4K decode significantly better than Pi 3/4
   - Could handle Jellyfin transcoding (but Jellyfin already on Synology)
   - **Verdict**: ✅ Excellent if you need media processing

3. **Single Service** (DNS OR VPN OR Monitoring)
   - 2GB RAM enough for one service
   - 100 Mbps Ethernet limiting factor
   - **Verdict**: ✅ Works well for single service deployment

---

#### What Le Potato CANNOT Handle Well

**⚠️ Limitations for The Club**:

1. **WireGuard VPN** - Limited by 100 Mbps Ethernet
   - Throughput: ~80-100 Mbps max
   - CPU: 15-25% under load
   - **Verdict**: ⚠️ Acceptable but limited (Pi 4 is better with Gigabit)

2. **Multiple Services** - 100 Mbps bottleneck
   - DNS + VPN + Monitoring = network congestion
   - All traffic shares 100 Mbps Ethernet
   - **Verdict**: ⚠️ Not recommended (Pi 4 better for multi-service)

3. **High Network I/O** - 100 Mbps ceiling
   - Cannot saturate Gigabit home network
   - Backup operations will be slow
   - **Verdict**: ❌ Skip for high I/O tasks

---

#### Le Potato vs Raspberry Pi Comparison

| Feature | Le Potato | Pi 3 B+ | Pi 4 (2GB) |
|---------|-----------|---------|------------|
| **RAM** | 2GB DDR3 ✅ | 1GB ⚠️ | 2GB ✅ |
| **Network** | 100 Mbps ⚠️ | 100 Mbps ⚠️ | Gigabit ✅ |
| **WiFi** | None ❌ (USB dongle) | Built-in ✅ | Built-in ✅ |
| **4K Video Decode** | Hardware ✅ | Software ❌ | Software ❌ |
| **eMMC Support** | Yes ✅ | No ❌ | No ❌ |
| **CPU Performance** | 150-200% of Pi 3 | Baseline | 300-350% of Pi 3 ✅ |
| **Pi-hole + Unbound** | ✅ Excellent | ✅ Excellent | ✅ Excellent |
| **VPN Throughput** | ~80-100 Mbps ⚠️ | ~50-80 Mbps | ~200+ Mbps ✅ |
| **Multiple Services** | ⚠️ Limited | ❌ No | ✅ Yes |
| **Power Draw** | 2-3W ✅ | 2-3W | 4-6W |
| **Community Support** | Small ⚠️ | Large ✅ | Large ✅ |
| **Price** | $35 | $15-25 used | $35-45 used |

---

#### Real-World Performance for Homelab

**Pi-hole + Unbound Performance**:
- **DNS Latency**: 8-12ms (similar to Pi 3/4)
- **Throughput**: 6,000+ queries/sec (plenty for home use)
- **CPU Usage**: 5-10% average
- **Reliability**: Excellent with eMMC (better than microSD)
- **Verdict**: ✅ **Rock-solid** for DNS according to user reports

**Docker Support**:
- ✅ Docker and Docker Compose work well
- ✅ Armbian provides good Linux support (Ubuntu 22.04 LTS)
- ✅ Pi-hole installation "runs without hiccups"
- ⚠️ Smaller image ecosystem than Raspberry Pi

---

#### Recommendation for Le Potato

**When to Choose Le Potato**:
- ✅ If you need hardware 4K video decode
- ✅ If you want eMMC reliability (better than microSD)
- ✅ If you're running a single service (DNS)
- ✅ If you want better media transcoding than Pi 3/4

**When to Choose Pi 4 Instead**:
- ✅ If you need Gigabit Ethernet (VPN, backups, multi-service)
- ✅ If you need WiFi built-in
- ✅ If you want larger community support
- ✅ If you're running multiple services

**For The Club Homelab**:

**Best Use Case**: **Dedicated DNS Server with eMMC**
- Le Potato + eMMC module = extremely reliable DNS server
- 2GB RAM provides headroom
- 100 Mbps Ethernet is plenty for DNS
- Hardware video decode not needed for DNS (wasted capability)

**Verdict**: ⚠️ **Good but not ideal** for The Club
- **For DNS only**: Similar to Pi 3/4 (all three work great)
- **For multiple services**: Pi 4 is better (Gigabit Ethernet)
- **For media**: Excellent (but you already have Jellyfin on Synology)

**Recommendation**:
- **If you already own Le Potato**: Great for dedicated DNS server with eMMC
- **If buying new**: Get **Raspberry Pi 4 (4GB)** for flexibility (Gigabit + WiFi + larger community)

---

### Cost-Benefit: Pi Zero 2 W vs Pi 3 vs Pi 4 vs Le Potato

**Raspberry Pi Zero 2 W**:
- **New Price**: $15
- **Power**: 0.4-1W (ultra-low)
- **Best For**: Testing, IoT sensors, non-critical tasks
- **Verdict**: ❌ Not suitable for homelab critical services (WiFi-only, 512MB RAM)

**Raspberry Pi 3 Model B/B+**:
- **Used Price**: $15-25
- **Power**: 2-3W
- **Best For**: Single service (DNS with Pi-hole + Unbound)
- **Verdict**: ✅ Great value if you already own one, perfect for dedicated DNS

**Libre Computer Le Potato (2GB)**:
- **New Price**: $35
- **Power**: 2-3W
- **Best For**: DNS server with eMMC reliability, media transcoding
- **Advantages**: Hardware 4K decode, eMMC support, faster than Pi 3
- **Limitations**: 100 Mbps Ethernet (not Gigabit), no WiFi, smaller community
- **Verdict**: ⚠️ Good for DNS with eMMC, but Pi 4 better for multi-service/VPN

**Raspberry Pi 4 (4GB)**:
- **New Price**: $55
- **Used Price**: $35-45
- **Power**: 4-6W
- **Best For**: Multiple services (DNS + VPN + monitoring)
- **Verdict**: ✅ Better investment if buying new, most flexible

**Recommendation**:
- **Have Pi Zero 2 W?** Use for testing/IoT only, not for homelab services
- **Have Pi 3 already?** Perfect for Pi-hole + Unbound (dedicated DNS server)
- **Have Le Potato?** Excellent for DNS with eMMC module (very reliable)
- **Buying new for homelab?** Get **Pi 4 (4GB)** for flexibility (Gigabit + WiFi + community)
- **Buying new for media?** Le Potato has better 4K hardware decode than Pi 4
- **Budget constrained?** Pi 3 B+ used (~$20) best value for DNS only

---

**Savings on N97**: ~700-800MB RAM, significant CPU reduction

---

## Migration Guides

### Migration: Pi-hole + Unbound to Raspberry Pi 4

#### Prerequisites

1. **Raspberry Pi 4** with:
   - Raspberry Pi OS Lite (64-bit)
   - Static IP configured (e.g., `192.168.1.40`)
   - SSH enabled
   - 2GB+ RAM recommended

#### Step 1: Prepare Raspberry Pi 4

```bash
# SSH into Pi 4
ssh pi@192.168.1.40

# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Reboot
sudo reboot
```

#### Step 2: Copy Configuration from N97

On N97 PC:
```bash
# Export Pi-hole configuration
cd ~/TheClub
docker exec pihole pihole -a -t

# Copy config files
scp -r privacy/pihole pi@192.168.1.40:~/
scp -r privacy/unbound pi@192.168.1.40:~/
scp docker-compose.privacy.yml pi@192.168.1.40:~/
```

#### Step 3: Deploy on Raspberry Pi 4

On Raspberry Pi 4:
```bash
# Create directory structure
mkdir -p ~/dns-stack
cd ~/dns-stack

# Move configs
mv ~/pihole ~/dns-stack/
mv ~/unbound ~/dns-stack/
mv ~/docker-compose.privacy.yml ~/dns-stack/

# Edit docker-compose to only include Pi-hole and Unbound
nano docker-compose.privacy.yml
```

**Minimal docker-compose.privacy.yml** (Pi 4):
```yaml
version: "3.9"

services:
  unbound:
    image: mvance/unbound-rpi:latest
    container_name: unbound
    restart: unless-stopped
    hostname: unbound
    ports:
      - "5335:5335/tcp"
      - "5335:5335/udp"
    volumes:
      - ./unbound:/opt/unbound/etc/unbound
    healthcheck:
      test: ["CMD", "drill", "@127.0.0.1", "-p", "5335", "cloudflare.com"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      dns-network:
        ipv4_address: 172.20.0.2

  pihole:
    image: pihole/pihole:latest
    container_name: pihole
    restart: unless-stopped
    hostname: pihole
    ports:
      - "53:53/tcp"
      - "53:53/udp"
      - "8080:80/tcp"  # Web interface
    environment:
      - TZ=America/New_York
      - WEBPASSWORD=change_me
      - PIHOLE_DNS_=172.20.0.2#5335
      - DNSSEC=true
      - REV_SERVER=true
      - REV_SERVER_TARGET=192.168.1.1
      - REV_SERVER_CIDR=192.168.1.0/24
    volumes:
      - ./pihole/etc-pihole:/etc/pihole
      - ./pihole/etc-dnsmasq.d:/etc/dnsmasq.d
    depends_on:
      unbound:
        condition: service_healthy
    networks:
      dns-network:
        ipv4_address: 172.20.0.3

networks:
  dns-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/24
```

Deploy:
```bash
# Start services
docker compose up -d

# Check logs
docker logs pihole
docker logs unbound

# Access Pi-hole web interface
# http://192.168.1.40:8080/admin
```

#### Step 4: Update Router DNS

1. Log into your router
2. Change DNS server to: `192.168.1.40` (Pi 4 IP)
3. Save and reboot router

#### Step 5: Update Caddy Reverse Proxy

On N97 PC, update `Caddyfile`:

```caddyfile
# Pi-hole - Network-wide ad/tracker blocking (on Raspberry Pi 4)
dns.pochita.synology.me {
    reverse_proxy 192.168.1.40:8080  # Point to Pi 4

    log {
        output file /var/log/caddy/dns-access.log {
            roll_size 100mb
            roll_keep 5
            roll_keep_for 720h
        }
        format json
    }
}
```

Restart Caddy:
```bash
docker restart caddy
```

#### Step 6: Stop Services on N97

```bash
cd ~/TheClub
docker compose -f docker-compose.privacy.yml stop pihole unbound
```

#### Step 7: Verify

```bash
# Test DNS resolution
nslookup google.com 192.168.1.40

# Check Pi-hole stats
curl http://192.168.1.40:8080/admin/api.php

# Access via subdomain
curl https://dns.pochita.synology.me
```

---

### Migration: WireGuard VPN to Raspberry Pi 4

#### Step 1: Prepare Pi 4

```bash
# SSH into Pi 4
ssh pi@192.168.1.40

# Enable IP forwarding
sudo sysctl -w net.ipv4.ip_forward=1
sudo sysctl -w net.ipv6.conf.all.forwarding=1

# Make permanent
echo "net.ipv4.ip_forward=1" | sudo tee -a /etc/sysctl.conf
echo "net.ipv6.conf.all.forwarding=1" | sudo tee -a /etc/sysctl.conf
```

#### Step 2: Copy WireGuard Config

On N97:
```bash
cd ~/TheClub
scp -r security/wireguard pi@192.168.1.40:~/wireguard
```

#### Step 3: Deploy WireGuard on Pi 4

On Pi 4:
```bash
mkdir -p ~/vpn-stack
cd ~/vpn-stack

nano docker-compose.yml
```

```yaml
version: "3.9"

services:
  wireguard:
    image: linuxserver/wireguard:latest
    container_name: wireguard
    restart: unless-stopped
    cap_add:
      - NET_ADMIN
      - SYS_MODULE
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=America/New_York
      - SERVERURL=pochita.synology.me
      - SERVERPORT=51820
      - PEERS=laptop,phone,tablet
      - PEERDNS=172.20.0.3  # Pi-hole IP if on same Pi
      - INTERNAL_SUBNET=10.13.13.0
      - ALLOWEDIPS=0.0.0.0/0
    volumes:
      - ./wireguard/config:/config
      - /lib/modules:/lib/modules:ro
    ports:
      - "51820:51820/udp"
    sysctls:
      - net.ipv4.conf.all.src_valid_mark=1
```

Deploy:
```bash
docker compose up -d
docker logs wireguard
```

#### Step 4: Port Forwarding

Update router port forwarding:
- **Old**: Port 51820 → N97 PC
- **New**: Port 51820 → Raspberry Pi 4 (`192.168.1.40`)

#### Step 5: Stop on N97

```bash
cd ~/TheClub
docker compose -f docker-compose.security.yml stop wireguard
```

---

### Migration: Uptime Kuma to Raspberry Pi 4

#### Step 1: Deploy on Pi 4

```bash
ssh pi@192.168.1.40
mkdir -p ~/monitoring
cd ~/monitoring

nano docker-compose.yml
```

```yaml
version: "3.9"

services:
  uptime-kuma:
    image: louislam/uptime-kuma:1
    container_name: uptime-kuma
    restart: unless-stopped
    ports:
      - "3001:3001"
    volumes:
      - ./uptime-kuma-data:/app/data
    deploy:
      resources:
        limits:
          memory: 256M
          cpus: '0.5'
```

```bash
docker compose up -d
```

#### Step 2: Export/Import Monitors

1. **On N97**: Access Uptime Kuma → Settings → Backup
   - Download backup JSON
2. **On Pi 4**: Access `http://192.168.1.40:3001` → Settings → Restore
   - Upload backup JSON

#### Step 3: Update Caddy

On N97, update `Caddyfile`:

```caddyfile
status.pochita.synology.me {
    reverse_proxy 192.168.1.40:3001  # Point to Pi 4

    log {
        output file /var/log/caddy/status-access.log {
            roll_size 100mb
            roll_keep 5
            roll_keep_for 720h
        }
        format json
    }
}
```

```bash
docker restart caddy
```

#### Step 4: Stop on N97

```bash
docker compose -f docker-compose.monitoring.yml stop uptime-kuma
```

---

## Network Architecture: Distributed Setup

### Option A: N97 as Central Reverse Proxy

```
Internet
  ↓
Router (Port 80/443 → N97)
  ↓
N97 PC (Caddy Reverse Proxy)
  ├─→ Local services on N97 (Frigate, Grafana, etc.)
  ├─→ Raspberry Pi 4 services (Pi-hole, Uptime Kuma, etc.)
  ├─→ Raspberry Pi 5 (Frigate detector)
  └─→ Synology NAS (Jellyfin, Vaultwarden)
```

**Caddyfile on N97**:
```caddyfile
dns.pochita.synology.me {
    reverse_proxy 192.168.1.40:8080  # Pi 4 - Pi-hole
}

status.pochita.synology.me {
    reverse_proxy 192.168.1.40:3001  # Pi 4 - Uptime Kuma
}

nvr.pochita.synology.me {
    reverse_proxy localhost:5000  # N97 - Frigate
}

metrics.pochita.synology.me {
    reverse_proxy localhost:3000  # N97 - Grafana
}
```

**Pros**:
- ✅ Simple (one entry point)
- ✅ Single SSL certificate management
- ✅ N97 handles all external traffic

**Cons**:
- ⚠️ N97 is single point of failure for web access
- ⚠️ All traffic flows through N97

---

### Option B: Raspberry Pi 4 as Central Reverse Proxy

```
Internet
  ↓
Router (Port 80/443 → Pi 4)
  ↓
Raspberry Pi 4 (Caddy Reverse Proxy)
  ├─→ Local services on Pi 4 (Pi-hole, Uptime Kuma)
  ├─→ N97 PC services (Frigate, Grafana, etc.)
  ├─→ Raspberry Pi 5 (Frigate detector)
  └─→ Synology NAS (Jellyfin, Vaultwarden)
```

**Pros**:
- ✅ Offloads reverse proxy from N97
- ✅ Pi 4 consumes less power (always-on entry point)
- ✅ N97 only handles backend services

**Cons**:
- ⚠️ More complex setup
- ⚠️ Pi 4 becomes single point of failure
- ⚠️ Network latency for N97 services

**Recommendation**: **Option A** (N97 as reverse proxy) is simpler and more reliable

---

## Power Consumption & Cost

### Power Usage

| Device | Idle Power | Active Power | Monthly Cost (24/7) |
|--------|-----------|--------------|---------------------|
| N97 PC | 15W | 30W | $3-6/month |
| Raspberry Pi 4 | 3W | 6W | $0.50-1/month |
| Raspberry Pi 5 | 4W | 8W | $0.70-1.50/month |
| Raspberry Pi 3 Model B/B+ | 1.5W | 2.5W | $0.25-0.50/month |
| Libre Computer Le Potato | 1.5W | 3W | $0.25-0.60/month |
| Raspberry Pi Zero 2 W | 0.4W | 1W | $0.10-0.20/month |
| Synology DS415+ | 20W | 35W | $4-7/month |

**Scenario 1: Move Pi-hole + Unbound to Pi 3**
- N97 Power Savings: ~3W (reduced load)
- Pi 3 Additional: +2.5W (active)
- Net Change: -0.5W (saves power!)
- **Benefit**: Frees N97 resources, Pi 3 extremely power-efficient

**Scenario 2: Move Pi-hole, VPN, Uptime Kuma to Pi 4**
- N97 Power Savings: ~5W (reduced load)
- Pi 4 Additional: +6W (active)
- Net Change: +1W (minimal)
- **Benefit**: Frees N97 resources, Pi 4 handles multiple services efficiently

**Scenario 3: Use Both Pi 3 and Pi 4**
- Pi 3: DNS (2.5W)
- Pi 4: VPN + Monitoring (6W)
- N97 Savings: ~8W
- Net Change: +0.5W
- **Benefit**: Maximum N97 offload, distributed resilience

---

## Performance Comparison

### DNS Queries (Pi-hole + Unbound)

| Location | Query Latency | Throughput | CPU Usage | Notes |
|----------|---------------|------------|-----------|-------|
| N97 PC | 5-10ms | 10,000+ queries/sec | 2-5% | Wired |
| Pi 4 | 5-10ms | 8,000+ queries/sec | 5-10% | Wired |
| Le Potato | 8-12ms | 6,000+ queries/sec | 5-10% | Wired, 2GB RAM |
| Pi 3 B+ | 8-15ms | 5,000+ queries/sec | 8-15% | Wired |
| Pi Zero 2 W | 15-25ms | 2,000+ queries/sec | 15-25% | **WiFi + 512MB RAM** |

**Verdict**:
- N97, Pi 4, Le Potato, Pi 3: All excellent for home use (typical: 100-200 queries/minute)
- Le Potato: ✅ Similar performance to Pi 3/4, extra reliability with eMMC option
- Pi Zero 2 W: ❌ **Cannot run Pi-hole + Unbound together** (exceeds 512MB RAM)
  - Can only run Pi-hole alone with external DNS (less private)
  - WiFi adds 5-10ms latency vs Ethernet
  - Not recommended for production DNS

### VPN Throughput (WireGuard)

| Location | Throughput | CPU Usage | Latency | Network Type |
|----------|-----------|-----------|---------|--------------|
| N97 PC | 500+ Mbps | 5-10% | <1ms | Wired |
| Pi 4 | 100-200 Mbps | 10-20% | 1-2ms | Gigabit Ethernet |
| Le Potato | 80-100 Mbps | 15-25% | 1-2ms | **100 Mbps Ethernet** |
| Pi 3 B+ | 50-80 Mbps | 20-35% | 2-3ms | 100 Mbps Ethernet |
| Pi Zero 2 W | 20-40 Mbps | 35-50% | 5-10ms | **WiFi 2.4 GHz** |

**Verdict**:
- **N97**: Overkill for most home use
- **Pi 4**: Excellent for gigabit home internet
- **Le Potato**: ⚠️ Acceptable (~80-100 Mbps max), limited by 100 Mbps Ethernet, faster CPU than Pi 3
- **Pi 3**: Sufficient for <100 Mbps internet, limited by 100 Mbps Ethernet
- **Pi Zero 2 W**: ❌ Not recommended - WiFi bottleneck (~40 Mbps max), unreliable for VPN

---

## Monitoring Distributed Setup

### Prometheus Configuration

Update Prometheus to scrape all devices:

```yaml
scrape_configs:
  # N97 PC
  - job_name: 'n97-node-exporter'
    static_configs:
      - targets: ['localhost:9100']

  - job_name: 'n97-cadvisor'
    static_configs:
      - targets: ['localhost:8080']

  # Raspberry Pi 3 (if using for DNS)
  - job_name: 'pi3-node-exporter'
    static_configs:
      - targets: ['192.168.1.39:9100']

  # Raspberry Pi 4 (if using for VPN/monitoring)
  - job_name: 'pi4-node-exporter'
    static_configs:
      - targets: ['192.168.1.40:9100']

  # Raspberry Pi 5 (Frigate detector)
  - job_name: 'pi5-node-exporter'
    static_configs:
      - targets: ['192.168.1.50:9100']

  # Synology NAS
  - job_name: 'synology-snmp'
    static_configs:
      - targets: ['192.168.1.10:9116']
```

Install Node Exporter on Pis:
```bash
# On Pi 3, Pi 4, and Pi 5
docker run -d \
  --name node-exporter \
  --restart unless-stopped \
  --net host \
  --pid host \
  -v /:/host:ro,rslave \
  prom/node-exporter:latest \
  --path.rootfs=/host
```

---

## Summary Recommendation

### Scenario 1: If You Have Raspberry Pi 3 Model B/B+

**Move to Raspberry Pi 3** (Single Service):
1. ✅ **Pi-hole + Unbound ONLY** - Perfect fit, designed for Pi
   - Static IP: `192.168.1.39`
   - RAM: 120-150MB / 1GB
   - Power: 2-3W

**Keep on Raspberry Pi 5**:
1. ✅ **Frigate Detector** (Hailo AI)

**Keep on N97 PC**:
1. ✅ Everything else (Caddy, WireGuard, Authelia, Frigate NVR, Grafana, etc.)

**Keep on Synology NAS**:
1. ✅ **Vaultwarden** - Passwords
2. ✅ **Jellyfin** - Media

**Expected Results**:
- **N97 RAM Freed**: ~150MB (DNS offloaded)
- **N97 CPU Reduced**: 3-5%
- **Power**: Pi 3 uses only 2-3W (very efficient)
- **Resilience**: DNS on separate dedicated hardware

---

### Scenario 2: If You Have Both Pi 3 and Pi 4

**Move to Raspberry Pi 3**:
1. ✅ **Pi-hole + Unbound** - Dedicated DNS server
   - Static IP: `192.168.1.39`

**Move to Raspberry Pi 4**:
1. ✅ **WireGuard VPN** - Secure remote access
2. ✅ **Uptime Kuma** - Monitoring
3. ✅ **MQTT Broker** - IoT message broker
   - Static IP: `192.168.1.40`

**Keep on Raspberry Pi 5**:
1. ✅ **Frigate Detector** (Hailo AI)

**Keep on N97 PC**:
1. ✅ **Caddy** - Reverse proxy
2. ✅ **Authelia + Redis** - Central auth
3. ✅ **Frigate NVR** - Video processing
4. ✅ **Prometheus + Grafana** - Metrics
5. ✅ **Kopia + Dockge + Watchtower** - Management
6. ✅ **Searxng + FreshRSS** - Privacy apps
7. ✅ **Homer** - Dashboard

**Keep on Synology NAS**:
1. ✅ **Vaultwarden** - Passwords
2. ✅ **Jellyfin** - Media

**Expected Results**:
- **N97 RAM Freed**: ~500-600MB
- **N97 CPU Reduced**: 15-20%
- **Power**: Pi 3 (2-3W) + Pi 4 (4-6W) = 6-9W total
- **Maximum Offload**: Best resource distribution

---

### Scenario 3: If You Only Have Raspberry Pi 4

**Move to Raspberry Pi 4** (Recommended):
1. ✅ **Pi-hole + Unbound** - Perfect fit, designed for Pi
2. ✅ **WireGuard VPN** - Lightweight, always-on
3. ✅ **Uptime Kuma** - Monitors from separate device
4. ✅ **MQTT Broker** - Lightweight IoT hub

**Keep on Raspberry Pi 5**:
1. ✅ **Frigate Detector** (Hailo AI) - Already planned

**Keep on N97 PC**:
1. ✅ **Caddy** - Reverse proxy (central entry point)
2. ✅ **Frigate NVR** - Video processing
3. ✅ **Authelia + Redis** - Central auth
4. ✅ **Prometheus + Grafana** - Metrics
5. ✅ **Kopia + Dockge + Watchtower** - Management
6. ✅ **Searxng + FreshRSS** - Privacy apps
7. ✅ **Homer** - Dashboard

**Keep on Synology NAS**:
1. ✅ **Vaultwarden** - Passwords
2. ✅ **Jellyfin** - Media

### Expected Results

**N97 PC**:
- **RAM Freed**: ~400-500MB
- **CPU Reduced**: 10-15% average load reduction
- **Focus**: Video processing, metrics, heavy workloads

**Raspberry Pi 4**:
- **RAM Used**: ~300-400MB total
- **CPU Used**: 10-20% average
- **Always-On**: DNS, VPN, monitoring

**Benefits**:
- ✅ Better resource allocation
- ✅ More resilient (DNS/VPN on separate device)
- ✅ Lower power consumption for always-on services
- ✅ N97 freed up for CPU-intensive tasks (Frigate, etc.)

---

## Quick Start: Pi 4 Setup Script

```bash
#!/bin/bash
# Raspberry Pi 4 - Quick Setup for DNS/VPN/Monitoring

# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Enable IP forwarding (for VPN)
echo "net.ipv4.ip_forward=1" | sudo tee -a /etc/sysctl.conf
sudo sysctl -p

# Create directory structure
mkdir -p ~/services/{dns,vpn,monitoring}

# Set static IP (edit as needed)
# Use router DHCP reservation instead (recommended)

echo "Raspberry Pi 4 setup complete!"
echo "Next steps:"
echo "1. Configure DHCP reservation for static IP"
echo "2. Deploy docker-compose stacks for each service"
echo "3. Update N97 Caddyfile to point to Pi 4 services"
echo "4. Update router DNS to Pi 4 IP"
```

---

## Quick Start: Pi 3 Setup Script

```bash
#!/bin/bash
# Raspberry Pi 3 - Quick Setup for DNS Only (Pi-hole + Unbound)

# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER

# Reduce GPU memory (headless server)
echo "gpu_mem=16" | sudo tee -a /boot/firmware/config.txt

# Disable unnecessary services (save RAM)
sudo systemctl disable bluetooth
sudo systemctl disable hciuart

# Disable WiFi if using Ethernet (save power)
sudo rfkill block wifi

# Create directory for DNS stack
mkdir -p ~/dns-stack

# Set static IP (edit as needed)
# Use router DHCP reservation instead (recommended)

echo "Raspberry Pi 3 setup complete!"
echo "Next steps:"
echo "1. Configure DHCP reservation for static IP (e.g., 192.168.1.39)"
echo "2. Reboot: sudo reboot"
echo "3. Deploy Pi-hole + Unbound docker-compose stack"
echo "4. Update router DNS to Pi 3 IP"
echo "5. Update N97 Caddyfile to point dns.pochita.synology.me to Pi 3"
echo ""
echo "Note: Pi 3 has 1GB RAM - run DNS ONLY for best performance"
```

---

## Next Steps

1. ✅ Decide which services to move (recommend Pi-hole + WireGuard + Uptime Kuma)
2. ✅ Set up Raspberry Pi 4 with static IP
3. ✅ Deploy services on Pi 4 using migration guides
4. ✅ Update Caddy reverse proxy configuration
5. ✅ Update router settings (DNS, port forwarding)
6. ✅ Monitor performance with Prometheus/Grafana
7. ✅ Stop services on N97 after verification

---

**Distributed. Efficient. Resilient. 🚀**
