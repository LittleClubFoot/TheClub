# Self-Hosted Home Security System — Ring Replacement Guide

## Overview

This document covers the complete replacement of a Ring home security ecosystem (doorbell camera, security cameras, and door/window sensors) with a **privacy-first, self-hosted** alternative. All video processing, AI detection, storage, alarm logic, and notifications happen entirely on local hardware — no cloud subscriptions, no third-party data access, no internet dependency.

**Ring Components Being Replaced:**

| Ring Product | Self-Hosted Replacement | Protocol |
|---|---|---|
| Ring Video Doorbell | Reolink PoE Doorbell + Frigate NVR | RTSP/ONVIF |
| Ring Stick Up / Spotlight Cam | Reolink RLC-810A/820A IP cameras | RTSP/ONVIF |
| Ring Alarm Contact Sensors | Zigbee door/window sensors (Aqara/Sonoff) | Zigbee 3.0 |
| Ring Alarm Base Station | Home Assistant + Zigbee coordinator | Z-Wave/Zigbee |
| Ring Alarm Keypad | Home Assistant Alarmo + wall tablet | Zigbee/WiFi |
| Ring App (notifications) | HA Companion App + ntfy (self-hosted push) | Local/HTTPS |
| Ring Protect (cloud storage) | Frigate NVR → local NAS/SSD | Local only |

**Why self-host?**
- Ring uploads footage to Amazon servers — you don't control who sees it
- Ring has participated in law enforcement data-sharing programs without warrants
- Ring subscriptions cost $100-200/yr with no local fallback
- Cloud-dependent: no internet = no security
- Self-hosted means: your cameras, your storage, your rules

---

## Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         Home Network (LAN)                               │
├──────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────┐     ┌────────────────────────────────────────┐ │
│  │  PoE Switch          │     │  N97 PC (The Club Server)              │ │
│  │  (Camera Network)    │     │                                        │ │
│  │                      │     │  ┌──────────────┐  ┌───────────────┐  │ │
│  │  ┌─────────────────┐ │     │  │ Frigate NVR  │  │ Home Assistant│  │ │
│  │  │ Reolink Doorbell│─┼────▶│  │ :8971 web    │  │ :8123 web    │  │ │
│  │  │ (front door)    │ │     │  │ :8554 RTSP   │  │ Alarmo       │  │ │
│  │  └─────────────────┘ │     │  │ :8555 WebRTC │  │ Automations  │  │ │
│  │  ┌─────────────────┐ │     │  └──────┬───────┘  └──────┬────────┘  │ │
│  │  │ Reolink 810A    │─┼────▶│         │                 │           │ │
│  │  │ (backyard)      │ │     │  ┌──────▼─────────────────▼────────┐  │ │
│  │  └─────────────────┘ │     │  │         Mosquitto MQTT          │  │ │
│  │  ┌─────────────────┐ │     │  │         :1883                   │  │ │
│  │  │ Reolink 820A    │─┼────▶│  └──────┬─────────────────────────┘  │ │
│  │  │ (driveway)      │ │     │         │                             │ │
│  │  └─────────────────┘ │     │  ┌──────▼──────┐  ┌───────────────┐  │ │
│  │  ┌─────────────────┐ │     │  │ go2rtc      │  │ ntfy          │  │ │
│  │  │ Reolink 810A    │─┼────▶│  │ :1984 WebRTC│  │ :2586 push   │  │ │
│  │  │ (garage)        │ │     │  │ two-way audio│  │ notifications│  │ │
│  │  └─────────────────┘ │     │  └─────────────┘  └───────────────┘  │ │
│  └─────────────────────┘     │                                        │ │
│                               │  Zigbee Coordinator (USB stick)       │ │
│  ┌─────────────────────┐     │  └─→ Zigbee2MQTT :8080                │ │
│  │  Zigbee Sensors     │     └────────────────────────────────────────┘ │
│  │  (mesh network)     │                    │                           │
│  │                     │     ┌──────────────▼─────────────────────────┐ │
│  │  ◉ Front door      │────▶│  Raspberry Pi 5 + Hailo AI Hat         │ │
│  │  ◉ Back door       │     │  (Object Detection Node)               │ │
│  │  ◉ Garage door     │     │  - Person/vehicle/animal detection     │ │
│  │  ◉ Windows (x4)    │     │  - 13 TOPS, ~20ms inference            │ │
│  │  ◉ Motion sensors  │     └─────────────────────────────────────────┘ │
│  └─────────────────────┘                                                │
│                                                                          │
│  Storage:  N97 SSD (events/clips) + Synology NAS (long-term archive)   │
│  Access:   https://nvr.pochita.synology.me (Caddy → Frigate)           │
│            https://ha.pochita.synology.me  (Caddy → Home Assistant)    │
│                                                                          │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Hardware Bill of Materials

### Cameras

| Item | Model | Qty | Est. Price | Notes |
|---|---|---|---|---|
| Doorbell Camera | Reolink Video Doorbell PoE | 1 | $80 | ONVIF, RTSP, 2K+, two-way audio, no subscription |
| Outdoor Camera | Reolink RLC-810A (4K PoE) | 2 | $60 ea | 4K, RTSP native, excellent Frigate compat |
| Outdoor Camera | Reolink RLC-820A (4K PoE) | 1 | $70 | 4K with spotlight, same protocol support |
| PoE Switch | TP-Link TL-SG1005P (5-port) | 1 | $40 | Powers cameras over Ethernet, no wall warts |

**Why Reolink?**
- Native RTSP/ONVIF out of the box — no firmware hacks, no cloud dependency
- Officially certified "Works with Home Assistant"
- Best price-to-quality ratio for Frigate setups
- Local storage to SD card as fallback
- Block internet access entirely and cameras still function

**Alternative cameras** (also RTSP/ONVIF compatible):
- **Amcrest AD410** — doorbell, great HA integration, ~$80
- **Amcrest IP4M-1041** — outdoor, solid ONVIF, ~$55
- **Hikvision DS-2CD series** — pro-grade, configurable, ~$80-150
- **Dahua IPC-HDW series** — OEM for many brands, reliable, ~$60-100

### Zigbee Sensors (Ring Alarm Replacement)

| Item | Model | Qty | Est. Price | Notes |
|---|---|---|---|---|
| Zigbee USB Coordinator | Sonoff Zigbee 3.0 USB Dongle Plus (CC2652P) | 1 | $25 | Recommended coordinator for Zigbee2MQTT |
| Door/Window Sensor | Aqara Door & Window Sensor (MCCGQ11LM) | 6 | $12 ea | Zigbee 3.0, CR1632 battery (~2yr life) |
| Motion Sensor | Aqara Motion Sensor P1 (MS-S02) | 2 | $22 ea | Zigbee 3.0, adjustable sensitivity |
| Water Leak Sensor | Aqara Water Leak Sensor (SJCGQ11LM) | 2 | $15 ea | Bonus: Ring doesn't even offer this |
| Siren | Heiman Zigbee Smart Siren (HS2WD-E) | 1 | $30 | 95dB alarm, Zigbee 3.0 |
| Keypad (optional) | Ring Alarm Keypad (2nd Gen) | 1 | $30 | Z-Wave, works with HA without Ring subscription |

**Why Aqara sensors?**
- Zigbee 3.0 mesh networking — sensors relay through each other, extending range
- Tiny form factor, adhesive mount, 2+ year battery life
- Extremely well-supported in Zigbee2MQTT and Home Assistant ZHA
- ~$12/sensor vs Ring's ~$20/sensor, and no subscription required

**Alternative sensor brands:**
- **Sonoff SNZB-04** — door/window, ~$8 (budget pick)
- **Third Reality** — Zigbee sensors, good HA support, ~$10-15
- **Zooz Z-Wave sensors** — if you prefer Z-Wave over Zigbee

### Total Estimated Cost

| Category | Cost |
|---|---|
| Cameras (4) + PoE switch | ~$310 |
| Zigbee coordinator + sensors (11) | ~$220 |
| Siren | ~$30 |
| **Total hardware** | **~$560** |
| **Ring equivalent (annual subscription)** | **$0/yr saved** vs $200/yr Ring Protect Pro |

Pays for itself in under 3 years, then it's free forever.

---

## Software Stack

### Core Services

| Service | Purpose | Port | Image |
|---|---|---|---|
| **Frigate NVR** | Camera recording, AI detection, event management | 8971 (web), 8554 (RTSP), 8555 (WebRTC) | `ghcr.io/blakeblackshear/frigate:stable` |
| **Home Assistant** | Automation hub, alarm panel, sensor management | 8123 | `ghcr.io/home-assistant/home-assistant:stable` |
| **Mosquitto** | MQTT broker (Frigate ↔ HA communication) | 1883 | `eclipse-mosquitto:2` |
| **Zigbee2MQTT** | Zigbee sensor bridge to MQTT | 8080 | `koenkk/zigbee2mqtt:latest` |
| **go2rtc** | WebRTC streaming, two-way audio for doorbell | 1984 | Built into Frigate |
| **ntfy** | Self-hosted push notifications (phone alerts) | 2586 | `binber/ntfy:latest` |

### How They Connect

```
Cameras ──RTSP──▶ Frigate NVR ──MQTT──▶ Mosquitto ◀──MQTT── Home Assistant
                      │                     ▲                      │
                      │                     │                      │
                      ▼                     │                      ▼
                 go2rtc (WebRTC)      Zigbee2MQTT ◀─Zigbee─ Door Sensors
                 two-way audio              │                      │
                                            │                      ▼
                                            │              Alarmo (alarm panel)
                                            │                      │
                                            │                      ▼
                                            └──────────────── ntfy (push alerts)
                                                               │
                                                               ▼
                                                          Phone App
```

1. **Cameras** stream RTSP to **Frigate**, which records and runs AI detection
2. **Frigate** publishes detection events (person at front door) to **Mosquitto** (MQTT)
3. **Home Assistant** subscribes to MQTT and receives all Frigate events
4. **Zigbee sensors** report state changes through **Zigbee2MQTT** → **Mosquitto** → **Home Assistant**
5. **Alarmo** (HA integration) manages arm/disarm states and triggers the siren
6. **ntfy** sends push notifications to your phone — self-hosted, no Google/Apple push dependency

---

## Deployment

### Phase 1: Core Infrastructure (Day 1)

Deploy Home Assistant, Mosquitto, and Zigbee2MQTT using `docker-compose.security.yml`.

```bash
# From The Club project root
make security-up
# or directly:
docker compose -f docker-compose.security.yml up -d
```

**Post-deploy checklist:**
- [ ] Home Assistant accessible at `http://<server-ip>:8123`
- [ ] Create HA admin account during onboarding
- [ ] Mosquitto broker running (test with `mosquitto_pub -t test -m hello`)
- [ ] Zigbee2MQTT accessible at `http://<server-ip>:8080`
- [ ] Plug in Zigbee USB coordinator and verify Zigbee2MQTT detects it

### Phase 2: Pair Zigbee Sensors (Day 1-2)

1. Open Zigbee2MQTT web UI → enable "Permit Join"
2. For each sensor: press the pairing button (small pin-hole reset) for 5 seconds
3. Sensor appears in Zigbee2MQTT → automatically discovered in Home Assistant
4. Rename each sensor meaningfully: `front_door_contact`, `back_door_contact`, etc.
5. Disable "Permit Join" when done

### Phase 3: Install Alarmo (Day 2)

1. Install HACS (Home Assistant Community Store) in HA
2. Install **Alarmo** via HACS → restart HA
3. Configure Alarmo:
   - Add all door/window sensors as entry points
   - Add motion sensors as interior sensors
   - Configure arm modes (Home, Away, Night)
   - Set entry/exit delays (30s default)
   - Add siren as alarm output
   - Add ntfy notification on alarm trigger

### Phase 4: Camera Network + Frigate (Day 2-3)

1. Mount cameras, connect to PoE switch
2. Configure each camera's RTSP stream URL via its web interface
3. Block camera internet access at router level (optional but recommended)
4. Update `security/frigate/config.yml` with your camera RTSP URLs
5. Start Frigate: `docker compose -f docker-compose.security.yml up -d frigate`
6. Verify live feeds at `http://<server-ip>:8971`

### Phase 5: Notifications (Day 3)

1. Install HA Companion App on phone
2. Configure ntfy server in HA for rich push notifications
3. Create automations:
   - Person detected at front door → snapshot + push notification
   - Doorbell pressed → two-way audio prompt + notification
   - Alarm triggered → siren + notification to all family members
   - Door opened while armed → entry delay countdown notification

### Phase 6: Remote Access (Day 3-4)

Configure Caddy reverse proxy for secure remote access:
- `nvr.pochita.synology.me` → Frigate web UI
- `ha.pochita.synology.me` → Home Assistant
- All behind HTTPS with Authelia SSO (if configured)

---

## Frigate NVR Configuration

The Frigate config template lives at `security/frigate/config.yml`. Key sections:

**MQTT** — connects Frigate to Home Assistant via Mosquitto
**Cameras** — defines RTSP stream URLs, detection zones, recording rules
**Detectors** — configures the Hailo AI Hat on RPi5 (or CPU fallback)
**Recording** — event-based + optional continuous, retention policies
**Snapshots** — saves best detection frame per event
**Objects** — which objects to detect (person, car, dog, cat, package)

### Detection Zones

Define zones per camera so you only get alerts for relevant areas:

```yaml
# Example: front door camera
cameras:
  front_door:
    zones:
      porch:
        coordinates: 0.1,0.3,0.9,0.3,0.9,1.0,0.1,1.0
        objects:
          - person
          - package
      sidewalk:
        coordinates: 0.0,0.0,1.0,0.0,1.0,0.3,0.0,0.3
        objects:
          - person
          - car
```

### Storage Planning

| Mode | Storage per Camera/Month | Recommended For |
|---|---|---|
| Events only | 5-20 GB | Most users (recommended) |
| Events + continuous (low-res) | 50-100 GB | Paranoid mode |
| Full continuous (1080p) | ~2 TB | Overkill for most |

**Retention defaults in config:**
- Events/clips: 30 days on local SSD
- Snapshots: 30 days
- Continuous recordings: 7 days (if enabled)
- Archive: Synology NAS for long-term (rsync cron job)

---

## Security Hardening

### Network Isolation

```
VLAN 10 (IoT/Cameras)          VLAN 1 (Trusted LAN)
┌───────────────────┐           ┌──────────────────────┐
│ Cameras           │           │ N97 PC               │
│ Zigbee sensors    │◀─────────▶│ Phones/laptops       │
│ (no internet)     │  firewall │ (internet access)    │
└───────────────────┘  rules    └──────────────────────┘
```

**Recommended firewall rules:**
- Cameras: LAN access only, **block all internet**
- Zigbee devices: mesh-only (no IP connectivity)
- Frigate/HA: LAN access + specific external ports via Caddy
- ntfy: outbound HTTPS only (for push delivery)

### Camera Hardening
- Change default camera admin passwords immediately
- Disable UPnP on cameras
- Disable cloud/P2P features in camera firmware
- Block camera MAC addresses from internet at router
- Use dedicated PoE switch (isolated from main network if possible)

### Home Assistant Security
- Use strong admin password + 2FA (TOTP)
- Expose only through Caddy reverse proxy (never direct port)
- Keep HA updated — security patches are frequent
- Review installed integrations periodically

### Physical Security
- PoE switch + server on UPS (uninterruptible power supply)
- If power is cut, UPS keeps cameras recording for 15-30 min
- Zigbee sensors are battery-powered — work during power outage
- Consider cellular backup for critical alerts (optional)

---

## Home Assistant Automations

Key automations for the security system (defined in `security/homeassistant/automations.yaml`):

### Doorbell Press → Notification with Snapshot
When someone presses the doorbell button, Frigate captures a snapshot and HA sends a rich push notification with the image to your phone. Tap to open two-way audio via go2rtc/WebRTC.

### Person Detected → Smart Alert
Frigate detects a person in a camera zone → HA checks:
- Is alarm armed? If yes → trigger alarm + siren
- Is alarm disarmed? → send notification only (no siren)
- Is it a known face? (future: face recognition) → suppress alert

### Door Opened While Armed → Entry Delay
Contact sensor opens while Alarmo is armed → start entry delay countdown (30s). If not disarmed within delay → trigger alarm + siren + notification.

### Night Mode Auto-Arm
At 11 PM, automatically arm in "Night" mode (perimeter sensors active, interior motion sensors inactive so you can walk around inside).

### Morning Auto-Disarm
At 7 AM on weekdays, auto-disarm (configurable per household schedule).

---

## Comparison: Ring vs Self-Hosted

| Feature | Ring | Self-Hosted |
|---|---|---|
| Monthly cost | $10-20/mo ($200/yr Pro) | $0 |
| Video storage | Amazon cloud (30-180 days) | Local SSD/NAS (unlimited) |
| AI detection | Cloud-processed | Local (Hailo/Coral, ~20ms) |
| Internet required | Yes (useless offline) | No (fully local) |
| Data privacy | Amazon sees everything | Only you |
| Law enforcement access | Ring Neighbors / subpoenas | Requires physical access |
| Two-way audio | Yes (cloud relay) | Yes (go2rtc, local WebRTC) |
| Sensor support | Ring Z-Wave ecosystem | Any Zigbee/Z-Wave/WiFi |
| Customization | Minimal | Unlimited (HA automations) |
| Works during outage | No | Yes (UPS + battery sensors) |
| Face recognition | No | Possible (CompreFace/DoubleTake) |
| Integration | Ring ecosystem only | Anything with an API |

---

## Relationship to Existing Branches

This security setup builds on and references work from these branches:

| Branch | Relevance |
|---|---|
| `claude/infra-frigate-nvr-distributed-architecture` | Frigate NVR deployment, RPi5 Hailo detector setup, camera config |
| `claude/infra-monitoring-security-observability-automation` | Prometheus/Grafana monitoring for security services |
| `claude/infra-privacy-dns-search-rss` | Pi-hole DNS blocking (block camera telemetry), WireGuard VPN |
| `claude/docs-vaultwarden-synology-deployment` | Vaultwarden for storing camera/sensor credentials |
| `claude/ui-homer-dashboard-privacy-redesign` | Homer dashboard integration for security services |

---

## Resources

- [Frigate NVR Documentation](https://docs.frigate.video/)
- [Frigate GitHub](https://github.com/blakeblackshear/frigate)
- [Frigate-Compatible Doorbells Discussion](https://github.com/blakeblackshear/frigate/discussions/3572)
- [Home Assistant Reolink Integration](https://www.home-assistant.io/integrations/reolink/)
- [Home Assistant Amcrest Integration](https://www.home-assistant.io/integrations/amcrest/)
- [Zigbee2MQTT Supported Devices](https://www.zigbee2mqtt.io/supported-devices/)
- [Alarmo (HACS Custom Alarm Panel)](https://github.com/nielsfaber/alarmo)
- [ntfy Self-Hosted Push Notifications](https://ntfy.sh/)
- [go2rtc WebRTC Streaming](https://github.com/AlexxIT/go2rtc)
- [Frigate + Coral TPU Docker Guide](https://github.com/blakeblackshear/frigate/discussions/16650)
- [Frigate + HA Guide (2026)](https://thissmart.house/2026/02/11/make-your-existing-cameras-smarter-with-frigate-and-home-assistant/)
- [Local Control Video Doorbells Comparison](https://www.thesmarthomehookup.com/local-control-video-doorbells-reolink-unifi-amcrest-hikvision-dahua/)
- [Self-Hosted Show: Ring Doorbell Alternative](https://selfhosted.show/18)
