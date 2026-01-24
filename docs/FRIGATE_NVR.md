# Frigate NVR - Home Security System

## Overview

This guide covers deploying Frigate NVR (Network Video Recorder) on The Club home server for privacy-focused, local AI-powered home security monitoring.

**Architecture**: Distributed setup with:
- **N97 PC (The Club server)**: Frigate main instance, video recording, web interface
- **Raspberry Pi 5 + Hailo AI Hat**: Dedicated object detection accelerator node

**Why Frigate?**
- ✅ **Privacy-First**: All processing local, no cloud dependencies
- ✅ **AI Object Detection**: Real-time detection (person, car, dog, etc.)
- ✅ **Zero False Alerts**: Advanced motion detection + AI filtering
- ✅ **Open Source**: FOSS, no subscription fees
- ✅ **Hardware Accelerated**: Hailo-8L provides 100+ detections/sec
- ✅ **Rich Features**: Events, clips, snapshots, zones, notifications
- ✅ **Home Assistant Integration**: Native integration available

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Your Home Network (LAN)                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────┐         ┌──────────────────────────┐     │
│  │  EUFY/IP Cameras │         │ N97 PC (The Club)        │     │
│  │  (RTSP Streams)  │────────▶│                          │     │
│  │                  │         │ - Frigate Main Instance  │     │
│  │ camera1.local    │         │ - Video Recording        │     │
│  │ camera2.local    │         │ - Web Interface (5000)   │     │
│  │ camera3.local    │         │ - MQTT Broker (optional) │     │
│  └──────────────────┘         │ - Caddy Reverse Proxy    │     │
│                                │                          │     │
│                                └────────┬─────────────────┘     │
│                                         │                       │
│                                         │ Object Detection      │
│                                         │ via Network API       │
│                                         │                       │
│                                ┌────────▼─────────────────┐     │
│                                │ Raspberry Pi 5           │     │
│                                │ + Hailo AI Hat           │     │
│                                │                          │     │
│                                │ - Hailo-8L Accelerator   │     │
│                                │ - 100+ detections/sec    │     │
│                                │ - Network detector API   │     │
│                                │ - Low latency (~20ms)    │     │
│                                └──────────────────────────┘     │
│                                                                  │
│  External Access:                                               │
│  https://nvr.pochita.synology.me → Caddy → Frigate:5000        │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Component Roles

**N97 PC (The Club Server)**:
- Receives RTSP streams from cameras
- Manages recording (events, clips, continuous)
- Serves web interface (live view, events, settings)
- Sends frames to RPi5 for object detection
- Stores recordings (SSD/HDD)

**Raspberry Pi 5 + Hailo AI Hat**:
- Dedicated AI object detection accelerator
- Hailo-8L processes 100+ inferences per second
- Network-accessible detector API
- Minimal latency (~20ms per detection)
- No video storage (only processes frames)

**IP Cameras**:
- Provide RTSP/RTMP video streams
- Local network only (no cloud)
- Continuous streaming to Frigate

---

## Hardware Requirements

### N97 PC (The Club Server)

**Current Specifications**:
- **CPU**: Intel N97 (4-core, up to 3.6 GHz)
- **RAM**: Minimum 8GB recommended for Frigate
- **Storage**: SSD for recordings (250GB+ per camera month for 24/7)
- **Network**: Gigabit Ethernet
- **OS**: Linux (already running Docker)

**Storage Planning**:
```
Continuous Recording (1080p @ 20fps):
- ~3 GB/hour per camera
- ~72 GB/day per camera
- ~2.16 TB/month per camera

Event-Only Recording (typical):
- ~5-20 GB/month per camera
- Depends on motion/events
```

### Raspberry Pi 5 + Hailo AI Hat

**Required Hardware**:
- **Raspberry Pi 5**: 4GB or 8GB RAM model
- **Hailo AI Kit**: Includes Hailo-8L accelerator + M.2 HAT
- **Power Supply**: Official 27W USB-C power supply (Hailo requires extra power)
- **Storage**: 32GB+ microSD card (for OS and detector runtime)
- **Cooling**: Active cooling recommended (AI hat gets warm)
- **Network**: Gigabit Ethernet connection to LAN

**Hailo-8L Specifications**:
- **Performance**: 13 TOPS (Tera Operations Per Second)
- **Inference Speed**: ~15-25ms per detection
- **Power**: ~5W typical
- **Models**: Supports YOLO-based models optimized for Hailo

### IP Cameras

**Requirements**:
- **RTSP or RTMP support** (required for Frigate)
- **Local network access** (no cloud-only cameras)
- **Resolution**: 1080p minimum, 4K supported
- **Frame Rate**: 15-30 fps recommended
- **H.264 codec** (required), H.265 supported

---

## EUFY Camera Compatibility

### Important Considerations

⚠️ **EUFY Privacy Concerns**:
EUFY (Anker) has had privacy controversies:
- 2022: Unencrypted cloud uploads despite "local only" claims
- Cloud access without explicit user permission
- Privacy-focused users may want alternatives

### EUFY Models with RTSP Support

**Cameras with RTSP** (verify before purchasing):
- **Indoor Cam 2K** series (requires firmware update to enable RTSP)
- **EufyCam 2 Pro** (some models support RTSP via HomeBase)
- **Solo IndoorCam** series

⚠️ **RTSP Limitations**:
- May require beta firmware or manual enablement
- Not all EUFY models support RTSP
- Battery-powered cameras typically don't support continuous RTSP
- Verify RTSP support BEFORE purchasing

### Recommended Alternative Cameras

For privacy-focused, Frigate-optimized setups:

#### **Reolink** (Highly Recommended)
- **Models**: RLC-810A, RLC-820A, RLC-511WA
- **RTSP**: ✅ Native support out of the box
- **ONVIF**: ✅ Full support
- **Privacy**: Local only, optional cloud
- **Price**: $60-120 per camera
- **Quality**: Excellent 4K, good night vision
- **Frigate**: Widely used, well-documented

#### **Amcrest**
- **Models**: IP4M-1041, IP5M-T1179EW
- **RTSP**: ✅ Native support
- **ONVIF**: ✅ Full support
- **Privacy**: Local network, optional cloud
- **Price**: $50-100 per camera
- **Quality**: Good, reliable

#### **Dahua (OEM for many brands)**
- **Models**: IPC-HDW series
- **RTSP**: ✅ Native support
- **ONVIF**: ✅ Full support
- **Privacy**: Local network capable
- **Note**: Professional-grade, may need configuration

**Recommendation**: **Reolink RLC-810A** or **RLC-820A** for best Frigate compatibility and privacy.

---

## Deployment

### Phase 1: Raspberry Pi 5 Setup (Detector Node)

#### Step 1: Install Raspberry Pi OS

1. **Download Raspberry Pi Imager**:
   - https://www.raspberrypi.com/software/

2. **Flash OS**:
   - OS: **Raspberry Pi OS Lite (64-bit)** (Debian 12 based)
   - Configure hostname: `frigate-detector`
   - Enable SSH
   - Set username/password
   - Configure WiFi/Ethernet

3. **Boot and Update**:
   ```bash
   # SSH into RPi5
   ssh pi@frigate-detector.local

   # Update system
   sudo apt update && sudo apt upgrade -y

   # Install prerequisites
   sudo apt install -y git curl
   ```

#### Step 2: Install Hailo Drivers

```bash
# Install Hailo software stack
# Follow official Hailo installation guide for RPi5

# Clone Hailo RPi5 examples (if available)
git clone https://github.com/hailo-ai/hailo-rpi5-examples.git
cd hailo-rpi5-examples

# Run installation script
sudo ./install.sh

# Verify Hailo device
hailortcli scan
# Should show: Hailo-8L PCIe device detected
```

#### Step 3: Run Frigate Detector Container

Frigate supports running a standalone detector that other Frigate instances can use.

```bash
# Create directories
mkdir -p ~/frigate-detector
cd ~/frigate-detector

# Create config for detector mode
nano config.yml
```

**config.yml** (Detector Node):
```yaml
# Frigate Detector Node Configuration
# This RPi5 only runs object detection, no cameras/recording

detectors:
  hailo:
    type: hailo
    device: PCIe  # Hailo-8L via M.2 HAT

# API endpoint for remote detection
api:
  host: 0.0.0.0
  port: 5001

# No cameras configured (detector only)
cameras: {}
```

**Run Detector Container**:
```bash
# Pull Frigate image (ARM64)
docker pull ghcr.io/blakeblackshear/frigate:stable

# Run detector node
docker run -d \
  --name frigate-detector \
  --restart unless-stopped \
  --device /dev/hailo0 \
  -p 5001:5001 \
  -v ~/frigate-detector/config.yml:/config/config.yml:ro \
  -e FRIGATE_DETECTOR_ONLY=true \
  ghcr.io/blakeblackshear/frigate:stable
```

**Verify Detector**:
```bash
# Check logs
docker logs frigate-detector

# Test API (from N97 or other machine)
curl http://frigate-detector.local:5001/api/stats
```

#### Step 4: Set Static IP (Recommended)

```bash
# Edit dhcpcd config
sudo nano /etc/dhcpcd.conf
```

Add:
```
interface eth0
static ip_address=192.168.1.50/24
static routers=192.168.1.1
static domain_name_servers=192.168.1.1
```

Restart:
```bash
sudo systemctl restart dhcpcd
```

---

### Phase 2: N97 PC (The Club) - Frigate Main Instance

#### Step 1: Create Docker Compose Stack

```bash
cd ~/TheClub
nano docker-compose.frigate.yml
```

**docker-compose.frigate.yml**:
```yaml
version: "3.9"

services:
  frigate:
    image: ghcr.io/blakeblackshear/frigate:stable
    container_name: frigate
    restart: unless-stopped

    # Privileged for hardware access (if needed)
    privileged: false

    # Shared memory size (important for processing)
    shm_size: "256mb"

    # Ports
    ports:
      - "5000:5000"   # Web UI
      - "8554:8554"   # RTSP re-streaming
      - "8555:8555"   # WebRTC
      - "1984:1984"   # Go2RTC

    # Volumes
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - ./frigate/config.yml:/config/config.yml
      - ./frigate/storage:/media/frigate
      - type: tmpfs
        target: /tmp/cache
        tmpfs:
          size: 1G

    # Environment
    environment:
      - FRIGATE_RTSP_PASSWORD=${FRIGATE_RTSP_PASSWORD:-change_me}
      - TZ=America/New_York  # Adjust to your timezone

    # Resource limits
    deploy:
      resources:
        limits:
          memory: 4G
          cpus: '2.0'
        reservations:
          memory: 2G

    # Networking
    networks:
      - homelab

    # Health check
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:5000/api/version"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 60s

  # Optional: Mosquitto MQTT (for notifications)
  mqtt:
    image: eclipse-mosquitto:latest
    container_name: frigate-mqtt
    restart: unless-stopped
    ports:
      - "1883:1883"
    volumes:
      - ./frigate/mqtt/config:/mosquitto/config
      - ./frigate/mqtt/data:/mosquitto/data
      - ./frigate/mqtt/log:/mosquitto/log
    networks:
      - homelab
    deploy:
      resources:
        limits:
          memory: 128M
          cpus: '0.25'

networks:
  homelab:
    external: true
```

#### Step 2: Create Frigate Configuration

```bash
mkdir -p ~/TheClub/frigate
nano ~/TheClub/frigate/config.yml
```

**frigate/config.yml**:
```yaml
# ============================================================================
# Frigate NVR Configuration
# ============================================================================

# MQTT Configuration (optional - for Home Assistant integration)
mqtt:
  enabled: true
  host: mqtt
  port: 1883
  user: frigate
  password: change_me  # Update this

# Detectors Configuration
detectors:
  # Remote Hailo detector on Raspberry Pi 5
  hailo:
    type: hailo
    # Point to RPi5 detector API
    api_url: http://192.168.1.50:5001  # Update with RPi5 IP
    # For local Hailo (if running on N97):
    # type: hailo
    # device: PCIe

# Database
database:
  path: /media/frigate/frigate.db

# Model Configuration
model:
  path: /config/model_cache/yolov8n.tflite
  # Hailo-optimized models
  input_tensor: nhwc
  input_pixel_format: rgb
  labelmap_path: /config/model_cache/labelmap.txt
  width: 320
  height: 320

# Recording Configuration
record:
  enabled: true
  # Events + continuous segments
  retain:
    days: 7          # Keep recordings for 7 days
    mode: motion     # Record on motion/events (or 'all' for continuous)
  events:
    retain:
      default: 14    # Keep event clips for 14 days
      mode: motion
      objects:
        person: 30   # Keep person detections for 30 days

# Snapshots
snapshots:
  enabled: true
  clean_copy: true  # Save clean snapshot without bounding boxes
  retain:
    default: 14     # Keep snapshots for 14 days
    objects:
      person: 30    # Keep person snapshots for 30 days

# Object Detection
objects:
  # Track specific objects
  track:
    - person
    - car
    - dog
    - cat
    - bird
  # Object-specific filters
  filters:
    person:
      min_area: 5000        # Minimum pixel area
      max_area: 100000      # Maximum pixel area
      threshold: 0.7        # Confidence threshold (0-1)
    car:
      min_area: 10000
      threshold: 0.7

# Motion Detection
motion:
  threshold: 30            # Pixel difference threshold (lower = more sensitive)
  contour_area: 10         # Minimum contour area
  lightning_threshold: 0.8  # Ignore sudden brightness changes

# Go2RTC (live streaming)
go2rtc:
  streams:
    # Define camera streams here (added per camera below)

# ============================================================================
# CAMERAS
# ============================================================================

cameras:
  # Example Camera 1 (Front Door)
  front_door:
    enabled: true

    # FFMPEG inputs
    ffmpeg:
      inputs:
        # Main stream (high quality for recording)
        - path: rtsp://admin:password@192.168.1.100:554/h264Preview_01_main
          roles:
            - record
        # Sub stream (lower quality for detection)
        - path: rtsp://admin:password@192.168.1.100:554/h264Preview_01_sub
          roles:
            - detect

      # Hardware acceleration (if N97 supports)
      hwaccel_args: preset-vaapi

    # Detection settings
    detect:
      enabled: true
      width: 640
      height: 480
      fps: 10  # Lower FPS for detection (saves CPU)

    # Motion detection areas
    motion:
      mask:
        # Mask static areas (timestamp, logo, etc.)
        - 0,0,640,50  # Top banner

    # Object detection zones
    zones:
      driveway:
        coordinates: 100,480,300,480,300,200,100,200
        objects:
          - person
          - car
      porch:
        coordinates: 300,480,600,480,600,300,300,300
        objects:
          - person

    # Recording
    record:
      enabled: true
      retain:
        days: 7
        mode: motion
      events:
        retain:
          default: 14

    # Snapshots
    snapshots:
      enabled: true
      timestamp: false
      bounding_box: true
      crop: false
      required_zones:
        - driveway
        - porch

  # Example Camera 2 (Backyard)
  backyard:
    enabled: true

    ffmpeg:
      inputs:
        - path: rtsp://admin:password@192.168.1.101:554/h264Preview_01_main
          roles:
            - record
        - path: rtsp://admin:password@192.168.1.101:554/h264Preview_01_sub
          roles:
            - detect
      hwaccel_args: preset-vaapi

    detect:
      enabled: true
      width: 640
      height: 480
      fps: 10

    zones:
      yard:
        coordinates: 0,480,640,480,640,0,0,0
        objects:
          - person
          - dog
          - cat

    record:
      enabled: true
      retain:
        days: 7
        mode: motion

    snapshots:
      enabled: true

  # Add more cameras as needed...

# ============================================================================
# UI Configuration
# ============================================================================

ui:
  live_mode: mse  # Options: webrtc, mse, jsmpeg
  timezone: America/New_York
  use_experimental: false

# Logging
logger:
  default: info
  logs:
    frigate.mqtt: debug
    detector.hailo: debug
```

**Update Camera Streams**:
Replace `rtsp://admin:password@192.168.1.100:554/...` with your actual camera RTSP URLs.

#### Step 3: Create MQTT Configuration (Optional)

```bash
mkdir -p ~/TheClub/frigate/mqtt/config
nano ~/TheClub/frigate/mqtt/config/mosquitto.conf
```

```conf
listener 1883
allow_anonymous false
password_file /mosquitto/config/passwd
```

Create password file:
```bash
docker run -it --rm -v $(pwd)/frigate/mqtt/config:/mosquitto/config eclipse-mosquitto mosquitto_passwd -c /mosquitto/config/passwd frigate
# Enter password when prompted
```

#### Step 4: Create Storage Directory

```bash
mkdir -p ~/TheClub/frigate/storage
sudo chown -R 1000:1000 ~/TheClub/frigate
```

#### Step 5: Deploy Frigate

```bash
cd ~/TheClub

# Start Frigate stack
docker compose -f docker-compose.frigate.yml up -d

# Check logs
docker logs -f frigate

# Verify detector connection
docker logs frigate | grep hailo
```

#### Step 6: Access Web UI

Open browser:
```
http://[N97_IP]:5000
# Or via Caddy (after configuring):
https://nvr.pochita.synology.me
```

---

## Configuration

### Finding Camera RTSP URLs

#### Generic RTSP URL Format

```
rtsp://[username]:[password]@[camera_ip]:[port]/[stream_path]

Common ports:
- 554 (RTSP default)
- 8554 (alternative)

Common stream paths:
- /h264Preview_01_main  (Reolink)
- /cam/realmonitor?channel=1&subtype=0  (Dahua/Amcrest)
- /stream1  (Generic)
```

#### Finding Your Camera's RTSP URL

**Method 1: Check Camera Web Interface**
- Access camera via browser: `http://[camera_ip]`
- Look for "Network", "RTSP", or "Streaming" settings
- URL may be listed there

**Method 2: Use VLC Media Player**
1. Open VLC
2. Media → Open Network Stream
3. Try common URLs:
   ```
   rtsp://admin:password@192.168.1.100:554/h264Preview_01_main
   rtsp://admin:password@192.168.1.100:554/stream1
   ```

**Method 3: ONVIF Discovery** (if camera supports ONVIF)
```bash
# Install onvif-gui (Linux)
sudo apt install onvif-gui

# Or use ONVIF Device Manager (Windows)
```

### Camera-Specific RTSP URLs

**Reolink**:
```
Main: rtsp://admin:password@[IP]:554/h264Preview_01_main
Sub:  rtsp://admin:password@[IP]:554/h264Preview_01_sub
```

**Amcrest**:
```
Main: rtsp://admin:password@[IP]:554/cam/realmonitor?channel=1&subtype=0
Sub:  rtsp://admin:password@[IP]:554/cam/realmonitor?channel=1&subtype=1
```

**EUFY (if RTSP enabled)**:
```
rtsp://admin:password@[IP]:8554/live0
```

### Configuring Detection Zones

Zones allow you to define specific areas in the camera view.

**Steps**:
1. Access Frigate UI: `http://[N97_IP]:5000`
2. Navigate to **Camera** → **[Camera Name]**
3. Click **Mask & Zone Editor**
4. Draw zones by clicking points on the image
5. Copy coordinates
6. Update `config.yml`:

```yaml
zones:
  driveway:
    coordinates: 100,480,300,480,300,200,100,200  # Bottom-left → Top-right
    objects:
      - person
      - car
    filters:
      person:
        min_area: 5000
        threshold: 0.8
```

### Optimizing Detection

**Reduce False Positives**:
```yaml
objects:
  filters:
    person:
      min_area: 5000       # Ignore small detections
      max_area: 100000     # Ignore huge detections (shadows)
      threshold: 0.75      # Higher = more confident
      min_score: 0.5       # Minimum detection score
      mask:                # Mask static objects
        - 200,300,400,300,400,400,200,400  # Mask a tree
```

**Motion Masks** (ignore areas):
```yaml
motion:
  mask:
    - 0,0,640,50          # Timestamp area
    - 500,400,640,480     # Tree that sways
```

**FPS Tuning**:
```yaml
detect:
  fps: 10  # Lower saves CPU (5-10 is good for most)
```

---

## Caddy Integration

Add Frigate to The Club's reverse proxy.

### Update Caddyfile

Edit `Caddyfile` or `Caddyfile.subdomain.example`:

```caddyfile
# ============================================================================
# HOME SECURITY
# ============================================================================

# Frigate NVR - Home security with AI object detection
nvr.pochita.synology.me {
    # Protect with Authelia SSO
    forward_auth authelia:9091 {
        uri /api/verify?rd=https://auth.pochita.synology.me
        copy_headers Remote-User Remote-Groups Remote-Name Remote-Email
    }

    reverse_proxy frigate:5000 {
        # WebSocket support for live streaming
        header_up Connection {>Connection}
        header_up Upgrade {>Upgrade}
    }

    # Security headers
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        X-Content-Type-Options nosniff
        X-Frame-Options SAMEORIGIN
        Referrer-Policy strict-origin-when-cross-origin

        # Frigate CSP (allows WebSocket)
        Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; font-src 'self' data:; connect-src 'self' wss: https:; media-src 'self' blob:; worker-src 'self' blob:"
    }

    # Logging
    log {
        output file /var/log/caddy/nvr-access.log {
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

Restart Caddy:
```bash
docker restart caddy
```

---

## Homer Dashboard Integration

Add Frigate to Homer dashboard.

Edit `assets/config.yml`:

```yaml
services:
  # Add to Security section
  - name: "🔒 Privacy Infrastructure"
    icon: "fas fa-user-shield"
    items:
      # ... existing items ...

  # NEW: Home Security section
  - name: "📹 Home Security"
    icon: "fas fa-video"
    items:
      - name: "Frigate NVR"
        logo: "assets/icons/frigate.png"
        subtitle: "AI-powered security camera system"
        tag: "security"
        keywords: "nvr cameras security surveillance frigate ai detection"
        url: "https://nvr.pochita.synology.me"
        target: "_blank"
```

Download Frigate icon:
```bash
cd ~/TheClub/assets/icons
curl -o frigate.png https://raw.githubusercontent.com/blakeblackshear/frigate/dev/web/public/android-chrome-192x192.png
```

---

## Notifications

### MQTT + Home Assistant (Recommended)

If using Home Assistant:

1. **Add MQTT Integration** in HA
2. **Frigate Integration** auto-discovers cameras
3. **Create Automations**:

```yaml
# Example: Notify on person detection
automation:
  - alias: "Frigate - Person Detected Front Door"
    trigger:
      - platform: mqtt
        topic: frigate/events
        payload: front_door/person
    action:
      - service: notify.mobile_app
        data:
          title: "Person at Front Door"
          message: "Motion detected"
          data:
            image: "https://nvr.pochita.synology.me/api/events/{{ trigger.payload_json.id }}/snapshot.jpg"
```

### Webhook Notifications (Without Home Assistant)

Configure webhooks in `config.yml`:

```yaml
notifications:
  webhooks:
    - url: "https://ntfy.sh/your-topic"
      method: POST
      headers:
        Title: "Frigate Alert"
      events:
        - person_detected
```

Use services like:
- **ntfy.sh** (self-hosted push notifications)
- **Gotify** (FOSS notification server)
- **Apprise** (supports many notification services)

---

## Backup Strategy

### What to Backup

1. **Configuration**: `~/TheClub/frigate/config.yml`
2. **Database**: `~/TheClub/frigate/storage/frigate.db`
3. **Event Clips**: `~/TheClub/frigate/storage/clips/`
4. **Snapshots**: `~/TheClub/frigate/storage/snapshots/`

### Automated Backups with Kopia

If using The Club's Kopia stack:

```bash
# Add Frigate to Kopia policy
kopia policy set ~/TheClub/frigate \
  --compression=zstd \
  --keep-latest=7 \
  --keep-monthly=12
```

### Manual Backup Script

```bash
nano ~/TheClub/scripts/frigate-backup.sh
```

```bash
#!/bin/bash
# Frigate Backup Script

BACKUP_DIR="/backups/frigate"
DATE=$(date +%Y%m%d_%H%M%S)
SOURCE="/home/user/TheClub/frigate"

mkdir -p "$BACKUP_DIR"

# Backup config and database
tar -czf "$BACKUP_DIR/frigate_config_$DATE.tar.gz" \
    "$SOURCE/config.yml" \
    "$SOURCE/storage/frigate.db"

# Optional: Backup recent events (last 7 days)
find "$SOURCE/storage/clips" -type f -mtime -7 \
    -exec tar -czf "$BACKUP_DIR/frigate_clips_$DATE.tar.gz" {} +

# Keep only last 14 backups
ls -t "$BACKUP_DIR"/frigate_*.tar.gz | tail -n +15 | xargs -r rm

echo "Backup completed: $DATE"
```

```bash
chmod +x ~/TheClub/scripts/frigate-backup.sh

# Add to crontab (daily at 3 AM)
crontab -e
```

Add:
```
0 3 * * * /home/user/TheClub/scripts/frigate-backup.sh >> /var/log/frigate-backup.log 2>&1
```

---

## Maintenance

### Update Frigate

```bash
cd ~/TheClub

# Pull latest image
docker compose -f docker-compose.frigate.yml pull

# Recreate container
docker compose -f docker-compose.frigate.yml up -d

# Check logs
docker logs frigate
```

### Update Raspberry Pi 5 Detector

```bash
# SSH into RPi5
ssh pi@frigate-detector.local

# Update system
sudo apt update && sudo apt upgrade -y

# Update Hailo firmware (if available)
# Follow Hailo update procedures

# Update Frigate detector
docker pull ghcr.io/blakeblackshear/frigate:stable
docker restart frigate-detector
```

### Database Maintenance

```bash
# Vacuum database (optimize)
docker exec frigate sqlite3 /media/frigate/frigate.db "VACUUM;"

# Check database size
docker exec frigate du -sh /media/frigate/frigate.db
```

### Clear Old Recordings

```bash
# Manually delete old recordings
docker exec frigate find /media/frigate/recordings -type f -mtime +30 -delete

# Or adjust retention in config.yml:
record:
  retain:
    days: 7  # Adjust as needed
```

---

## Troubleshooting

### Issue: Cannot Connect to Hailo Detector

**Check RPi5 Detector**:
```bash
# SSH into RPi5
ssh pi@frigate-detector.local

# Check container
docker logs frigate-detector

# Check Hailo device
hailortcli scan

# Test API
curl http://localhost:5001/api/stats
```

**Check Network Connectivity**:
```bash
# From N97 PC
ping frigate-detector.local
curl http://192.168.1.50:5001/api/stats
```

**Update Frigate Config**:
Ensure `detectors.hailo.api_url` points to correct RPi5 IP.

### Issue: Camera Stream Not Working

**Test RTSP Stream**:
```bash
# Install ffmpeg
sudo apt install ffmpeg

# Test stream
ffmpeg -rtsp_transport tcp -i "rtsp://admin:password@192.168.1.100:554/h264Preview_01_main" -frames:v 1 test.jpg

# If successful, image saved to test.jpg
```

**Check Frigate Logs**:
```bash
docker logs frigate | grep camera_name
```

**Common Issues**:
- Wrong RTSP URL (check camera manual)
- Wrong credentials
- Camera firewall blocking RTSP
- Network connectivity

### Issue: High CPU Usage

**Reduce Detection FPS**:
```yaml
detect:
  fps: 5  # Lower FPS (was 10)
```

**Use Sub-Streams for Detection**:
```yaml
ffmpeg:
  inputs:
    - path: rtsp://.../main
      roles:
        - record       # High quality for recording
    - path: rtsp://.../sub
      roles:
        - detect       # Low quality for detection (saves CPU)
```

**Enable Hardware Acceleration**:
```yaml
ffmpeg:
  hwaccel_args: preset-vaapi  # For Intel iGPU
```

### Issue: Too Many False Detections

**Adjust Object Filters**:
```yaml
objects:
  filters:
    person:
      min_area: 8000     # Increase minimum size
      threshold: 0.85    # Increase confidence
```

**Add Motion Masks**:
Mask areas with constant motion (trees, flags):
```yaml
motion:
  mask:
    - 100,200,300,400  # Coordinates of area to ignore
```

**Use Required Zones**:
Only save events from specific zones:
```yaml
snapshots:
  required_zones:
    - driveway  # Only save if detected in driveway zone
```

### Issue: Recordings Taking Too Much Space

**Adjust Retention**:
```yaml
record:
  retain:
    days: 3  # Reduce from 7
    mode: motion  # Only record on motion (not continuous)
```

**Event-Only Recording**:
```yaml
record:
  enabled: false  # Disable continuous recording
  events:
    retain:
      default: 14  # Only keep event clips
```

---

## Performance Tuning

### N97 PC Optimizations

**Hardware Acceleration** (if N97 has Intel iGPU):
```yaml
ffmpeg:
  hwaccel_args: preset-vaapi
```

**Shared Memory** (increase if needed):
```yaml
services:
  frigate:
    shm_size: "512mb"  # Increase from 256mb
```

### Raspberry Pi 5 Optimizations

**Active Cooling**:
- Ensure active cooling is installed (Hailo generates heat)
- Monitor temperature: `vcgencmd measure_temp`

**Power Supply**:
- Use official 27W power supply (Hailo requires extra power)

**Overclock** (optional, careful):
```bash
# Edit config
sudo nano /boot/firmware/config.txt

# Add:
arm_freq=2800
gpu_freq=950
over_voltage=8
```

---

## Security Best Practices

### Network Segmentation

**VLAN for Cameras** (recommended):
- Place cameras on separate VLAN
- Block internet access for cameras
- Only allow access to Frigate server

**Firewall Rules**:
```bash
# Block camera internet access
iptables -A FORWARD -s 192.168.10.0/24 -j DROP  # Camera VLAN
iptables -A FORWARD -s 192.168.10.0/24 -d 192.168.1.50 -j ACCEPT  # Allow to N97
```

### Camera Hardening

- **Change default passwords** (use strong passwords)
- **Disable UPnP** on cameras
- **Disable cloud services** if available
- **Update firmware** regularly
- **Enable encryption** (if camera supports RTSPS)

### Frigate Security

- **Protect with Authelia** (already configured in Caddyfile)
- **Use strong RTSP password**
- **Regular updates** (Frigate + containers)
- **Limit network exposure** (only via Caddy reverse proxy)

---

## Privacy Considerations

### Why Frigate for Privacy?

✅ **Local Processing**:
- All AI processing on your hardware
- No cloud dependencies
- No third-party access to footage

✅ **No Telemetry**:
- Frigate sends no data externally
- Open source (auditable)

✅ **Full Control**:
- You own the recordings
- You control retention
- You decide who has access

### Privacy-Focused Configuration

**Disable External Access** (optional):
```caddyfile
# Remove from Caddyfile to make LAN-only
# nvr.pochita.synology.me { ... }
```

**Blur Faces** (experimental):
```yaml
# Future Frigate feature
objects:
  filters:
    person:
      blur_face: true
```

**Automatic Deletion**:
```yaml
record:
  retain:
    days: 3  # Short retention for privacy
snapshots:
  retain:
    default: 7
```

---

## Cost Estimate

### Hardware Costs

| Item | Price | Notes |
|------|-------|-------|
| Raspberry Pi 5 (8GB) | $80 | Already owned |
| Hailo AI Kit | $70 | Includes M.2 HAT + Hailo-8L |
| RPi5 Power Supply (27W) | $12 | Official |
| microSD Card (64GB) | $10 | For RPi5 OS |
| **Cameras** | | |
| Reolink RLC-810A (4K) × 3 | $240 | $80 each |
| Reolink RLC-820A (4K PTZ) × 1 | $120 | Optional |
| **Storage** | | |
| 2TB SSD (for recordings) | $150 | ~2 months 3-camera |
| **Total** | **$682** | For complete system |

### Ongoing Costs

- **Electricity**: ~50W total (~$4/month)
- **Storage**: Replace SSD every 3-5 years
- **No subscriptions** (unlike cloud NVR services)

**Comparison**:
- Cloud NVR (Nest, Ring): $10-30/month = $360/year
- Frigate: $0/month after initial hardware

---

## Next Steps

1. ✅ **Raspberry Pi 5 Setup**:
   - Install Raspberry Pi OS
   - Install Hailo drivers
   - Deploy Frigate detector container

2. ✅ **Purchase Cameras**:
   - Recommended: Reolink RLC-810A (RTSP native)
   - Verify RTSP support before buying

3. ✅ **Deploy Frigate on N97**:
   - Create docker-compose.frigate.yml
   - Configure cameras in config.yml
   - Deploy stack

4. ✅ **Configure Caddy**:
   - Add nvr.pochita.synology.me subdomain
   - Restart Caddy

5. ✅ **Set Up Notifications**:
   - Configure MQTT (optional)
   - Set up webhooks or Home Assistant

6. ✅ **Test & Tune**:
   - Verify object detection
   - Create zones
   - Adjust filters to reduce false positives

---

## Related Documentation

- **Frigate Official Docs**: https://docs.frigate.video/
- **Hailo AI Kit**: https://www.raspberrypi.com/documentation/accessories/ai-kit.html
- **Reolink Camera Setup**: Camera manual
- **The Club - Network Architecture**: `docs/NETWORK_ARCHITECTURE.md`
- **The Club - Homer Dashboard**: `docs/HOMER_CUSTOMIZATION.md`

---

## Quick Reference

```bash
# Start Frigate
cd ~/TheClub && docker compose -f docker-compose.frigate.yml up -d

# Stop Frigate
docker compose -f docker-compose.frigate.yml down

# View logs
docker logs -f frigate

# Restart detector (RPi5)
ssh pi@frigate-detector.local
docker restart frigate-detector

# Access UI
https://nvr.pochita.synology.me
# Or local: http://[N97_IP]:5000

# Check storage usage
docker exec frigate du -sh /media/frigate

# Backup config
cp ~/TheClub/frigate/config.yml ~/TheClub/frigate/config.yml.backup
```

---

**Privacy-focused, AI-powered home security. No cloud. No subscriptions. 📹🔒**
