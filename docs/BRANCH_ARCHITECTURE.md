# The Club — Branch Architecture & Design Document

> **Generated:** 2026-02-14
> **Repository:** `github.com/LittleClubFoot/TheClub`
> **Base Branch:** `main` / `master`

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Branch Lineage Map](#branch-lineage-map)
3. [Branch Dependency Chain](#branch-dependency-chain)
4. [Branch-by-Branch Analysis](#branch-by-branch-analysis)
   - [main / master](#1-main--master)
   - [claude/review-home-server-G6FQx](#2-claudereview-home-server-g6fqx)
   - [claude/synology-nas-docs-G6FQx](#3-claudesynology-nas-docs-g6fqx)
   - [claude/privacy-infrastructure-G6FQx](#4-claudeprivacy-infrastructure-g6fqx)
   - [claude/homer-enhancements-G6FQx](#5-claudehomer-enhancements-g6fqx)
   - [claude/vaultwarden-docs-G6FQx](#6-claudevaultwarden-docs-g6fqx)
   - [claude/frigate-nvr-G6FQx](#7-claudefrigate-nvr-g6fqx)
   - [claude/image-notes-qr-sharing-VFfib](#8-claudeimage-notes-qr-sharing-vffib)
5. [Cumulative Service Inventory](#cumulative-service-inventory)
6. [Infrastructure Layer Diagram](#infrastructure-layer-diagram)
7. [Resource Budget](#resource-budget)
8. [Merge Strategy & Recommendations](#merge-strategy--recommendations)

---

## Project Overview

**The Club** is a containerized home lab dashboard built with Go, HTMX, Caddy, and Docker Compose, designed to run on a Synology NAS alongside dedicated hardware. It started as a lightweight Homer dashboard with a Go/HTMX test server and has evolved across seven Claude-assisted development branches into a comprehensive, privacy-first home infrastructure platform covering:

- Service monitoring & observability
- VPN & single sign-on authentication
- Network-level ad/tracker blocking
- Self-hosted password management
- AI-powered security cameras
- Distributed multi-device architecture
- Image notes with QR sharing & thermal label printing

**Tech Stack (base):** Go 1.21 · html/template · HTMX 1.9 · Caddy 2 · Homer · Docker Compose · Alpine Linux

---

## Branch Lineage Map

There are two distinct development lineages: a **linear chain** of six G6FQx-session branches that build incrementally on each other, and one **independent feature branch** from a separate VFfib session.

```
main (2 commits: LICENSE + initial)
│
├─── LINEAR CHAIN (G6FQx session) ────────────────────────────────────────────
│    │
│    ├─ claude/review-home-server-G6FQx        (+9 commits)
│    │   └─ claude/synology-nas-docs-G6FQx     (+4 commits)
│    │       └─ claude/privacy-infrastructure-G6FQx  (+4 commits)
│    │           └─ claude/homer-enhancements-G6FQx  (+2 commits)
│    │               └─ claude/vaultwarden-docs-G6FQx    (+3 commits)
│    │                   └─ claude/frigate-nvr-G6FQx     (+8 commits)
│
├─── INDEPENDENT BRANCH (VFfib session) ──────────────────────────────────────
│    │
│    └─ claude/image-notes-qr-sharing-VFfib    (+3 commits)
│
```

The G6FQx chain is strictly additive — each branch contains all commits from the branch above it plus its own new commits. The VFfib branch is based directly on `main` and has no shared commits with the G6FQx chain.

---

## Branch Dependency Chain

| # | Branch | Parent | Unique Commits | Primary Focus |
|---|--------|--------|:-:|---|
| 1 | `main` / `master` | — | 2 | Bare repository with LICENSE |
| 2 | `claude/review-home-server-G6FQx` | main | 9 | FOSS enhancement stacks (4 phases) |
| 3 | `claude/synology-nas-docs-G6FQx` | #2 | 4 | Synology NAS deployment docs |
| 4 | `claude/privacy-infrastructure-G6FQx` | #3 | 4 | Pi-hole, Unbound, Searxng, FreshRSS |
| 5 | `claude/homer-enhancements-G6FQx` | #4 | 2 | Dashboard redesign (privacy theme) |
| 6 | `claude/vaultwarden-docs-G6FQx` | #5 | 3 | Self-hosted password manager |
| 7 | `claude/frigate-nvr-G6FQx` | #6 | 8 | AI security cameras + distributed arch |
| 8 | `claude/image-notes-qr-sharing-VFfib` | main | 3 | Notes app + QR codes + label printer |

**Total unique commits across all branches:** 35 (excluding the 2 on main)

---

## Branch-by-Branch Analysis

### 1. main / master

**Commits:** 2
**State:** Bare project scaffold

Contains the foundational Go web server, Homer dashboard configuration, Caddy reverse proxy, Docker Compose dev/prod setups, Makefile, documentation (`docs/ARCHITECTURE.md`, `DEVELOPMENT.md`, `NAS_INTEGRATION.md`, `RUNBOOK.md`, `WEB_ARCHITECTURE.md`), and HTML templates.

**Services:** Homer dashboard, Go/HTMX test server, Caddy reverse proxy
**Routes:** `/` (Homer), `/test` (Go server), `/test/time`, `/test/api`, `/test/docs`, `/health`

---

### 2. claude/review-home-server-G6FQx

**Parent:** main · **Unique commits:** 9 · **Lines added:** ~2,035

**Purpose:** Transform The Club from a basic dashboard into a full-featured home server platform with monitoring, security, observability, and automation.

**Implementation is phased over 4 weeks:**

| Phase | Services Added | RAM Impact |
|-------|---------------|:----------:|
| Phase 1 — Monitoring | Uptime Kuma, config backup script | +384 MB |
| Phase 2 — Security | WireGuard VPN, Authelia SSO, Redis | +512 MB |
| Phase 3 — Observability | Prometheus, Grafana, Node Exporter, cAdvisor | +1,088 MB |
| Phase 4 — Automation | Kopia backup, Watchtower auto-update, Dockge UI | +896 MB |

**New files:**
- `docker-compose.monitoring.yml` — Uptime Kuma stack
- `docker-compose.security.yml` — WireGuard + Authelia + Redis
- `docker-compose.observability.yml` — Prometheus + Grafana + exporters
- `docker-compose.automation.yml` — Kopia + Watchtower + Dockge
- `authelia/configuration.yml` — SSO access control rules
- `authelia/users_database.yml` — User management template
- `monitoring/prometheus/prometheus.yml` — Scrape configuration
- `monitoring/grafana/provisioning/` — Datasources & dashboards
- `scripts/config-backup.sh` — Automated backup utility
- `docs/IMPLEMENTATION_PLAN.md` — Step-by-step deployment guide
- `ENHANCEMENTS.md` — Quick reference for all stacks

**Modified files:**
- `Caddyfile` — Added routes for `/status`, `/grafana`, `/auth`, `/dockge`, `/backup`; added CSP headers
- `Makefile` — Added `stack-monitoring`, `stack-security`, `stack-observability`, `stack-automation`, `stack-all`, `stack-stop`, `backup` targets
- `docker-compose.prod.yml` — Added resource limits, health checks, Caddy metrics port

---

### 3. claude/synology-nas-docs-G6FQx

**Parent:** #2 review-home-server · **Unique commits:** 4 · **Lines added:** ~1,952

**Purpose:** Comprehensive documentation for deploying The Club alongside a Synology DS415+ NAS using subdomain-based routing under `pochita.synology.me`.

**New files:**
- `docs/NETWORK_ARCHITECTURE.md` (464 lines) — Network topology, DNS strategies (Cloudflare / DuckDNS / local), port forwarding options, SSL/TLS certificate strategy, physical deployment options
- `docs/SYNOLOGY_DEPLOYMENT.md` (855 lines) — 4-phase deployment guide: DNS prep → base stack → enhancement stacks → service protection & maintenance
- `docs/SYNOLOGY_QUICK_START.md` (289 lines) — Condensed 5-minute setup reference
- `Caddyfile.subdomain.example` (344 lines) — Production-ready Caddy config with subdomain routing for 8+ services, Authelia forward-auth, Let's Encrypt ACME

**Key architectural decision:** Subdomain routing (`status.pochita.synology.me`, `grafana.pochita.synology.me`, etc.) instead of path-based routing — cleaner URLs, per-service SSL, better isolation.

---

### 4. claude/privacy-infrastructure-G6FQx

**Parent:** #3 synology-nas-docs · **Unique commits:** 4 · **Lines added:** ~1,418

**Purpose:** Network-level privacy protection stack: DNS-level ad/tracker blocking, private recursive DNS resolution, privacy-respecting search, and self-hosted RSS.

**New services:**

| Service | Port | Purpose | RAM |
|---------|------|---------|:---:|
| Pi-hole | 53 (DNS), 8053 (web) | Network-wide ad & tracker blocking | 512 MB |
| Unbound | 5335 | Private recursive DNS resolver (queries root servers directly) | 256 MB |
| Searxng | 8080 | Privacy meta-search engine (DuckDuckGo, Brave, Startpage) | 512 MB |
| FreshRSS | 80 | Self-hosted RSS feed reader | 256 MB |
| Redis | 6379 | Searxng cache | 128 MB |

**New files:**
- `docker-compose.privacy.yml` (274 lines) — Full orchestration for all 5 services
- `privacy/unbound/unbound.conf` (147 lines) — DNSSEC, QNAME minimization, 128 MB cache, client-subnet privacy
- `privacy/searxng/settings.yml` (238 lines) — Disabled Google/Bing, enabled DuckDuckGo/Brave/Startpage, image proxy, POST method
- `docs/PRIVACY_INFRASTRUCTURE.md` (676 lines) — Philosophy, deployment, device configuration (router/Windows/macOS/Linux/iOS/Android), verification, maintenance

**Modified files:**
- `Caddyfile.subdomain.example` — Added `dns.`, `search.`, `rss.` subdomain routes
- `Makefile` — Added `stack-privacy` target, updated `stack-all` and `stack-stop`

**Privacy outcomes:** 50–80% DNS query blocking, 30–40% bandwidth reduction, zero third-party DNS logging.

---

### 5. claude/homer-enhancements-G6FQx

**Parent:** #4 privacy-infrastructure · **Unique commits:** 2 · **Lines added:** ~782

**Purpose:** Rebrand the Homer dashboard from a generic home lab interface to a privacy-first showcase that reflects the full infrastructure built across prior branches.

**Changes to `assets/config.yml`:**
- Title → "The Club" with subtitle "Privacy-First Home Lab @ pochita.synology.me"
- Icon → `fas fa-shield-alt` (privacy shield)
- Color scheme → Purple (`#667eea` / `#764ba2`) for privacy identity
- Layout → Responsive `auto` columns (was fixed 3)
- Quick links → Added privacytools.io and EFF
- Service categories expanded from 2 to **7**: Privacy Infrastructure, Monitoring, Security, Management, Media & Files, Development, System
- All 14 services now have searchable keywords and color-coded tags

**New file:**
- `docs/HOMER_CUSTOMIZATION.md` (558 lines) — Complete customization guide: color schemes (Privacy Purple, Security Blue, Nature Green, Hacker Dark), icon sources, service management, responsive design, custom CSS, troubleshooting

---

### 6. claude/vaultwarden-docs-G6FQx

**Parent:** #5 homer-enhancements · **Unique commits:** 3 · **Lines added:** ~1,066

**Purpose:** Self-hosted password management via Vaultwarden on the Synology DS415+ NAS, proxied through Caddy.

**New file:**
- `docs/VAULTWARDEN_SYNOLOGY.md` (1,002 lines) — Complete deployment guide covering:
  - DS415+ hardware assessment (ARM dual-core, 1 GB DDR3, ARMv7 32-bit compatibility)
  - Docker CLI and Synology GUI deployment methods
  - Client setup (browser extensions, mobile apps, desktop)
  - Security hardening (disable public signups, admin token, 2FA, firewall rules)
  - Backup strategies (Kopia integration, Hyper Backup, manual scripts)
  - Bitwarden Cloud migration path
  - Performance optimization (SQLite WAL mode, 256 MB memory limit)
  - Troubleshooting (WebSocket issues, DNS, database locks, OOM on DS415+)

**Modified files:**
- `Caddyfile.subdomain.example` — Added `vault.pochita.synology.me` with WebSocket support, custom CSP for `wss:`, explicit note to NOT use Authelia (breaks browser extensions/mobile apps)
- `assets/config.yml` — Vaultwarden added to Security category
- `docs/HOMER_CUSTOMIZATION.md` — Updated service/icon references

---

### 7. claude/frigate-nvr-G6FQx

**Parent:** #6 vaultwarden-docs · **Unique commits:** 8 · **Lines added:** ~3,322

**Purpose:** AI-powered home security cameras (Frigate NVR) and a comprehensive distributed architecture strategy for distributing services across multiple hardware devices.

**New files:**
- `docs/FRIGATE_NVR.md` (1,417 lines) — Complete Frigate deployment guide:
  - N97 PC as main Frigate instance + Raspberry Pi 5 with Hailo-8L AI Hat for object detection (100+ detections/sec, ~20 ms latency)
  - EUFY/IP camera RTSP integration
  - Object detection zones, event filtering, recording strategies
  - Home Assistant integration, notifications
  - Privacy-first: all processing local, no cloud

- `docs/DISTRIBUTED_ARCHITECTURE.md` (1,840 lines) — Multi-device strategy:

  | Device | Role | RAM |
  |--------|------|-----|
  | N97 PC (Intel 4-core, 8 GB) | Main server: Caddy, Frigate, Authelia, Prometheus, Grafana, Searxng, FreshRSS, Homer | Primary |
  | Raspberry Pi 5 (8 GB) + Hailo Hat | AI object detection for Frigate | Dedicated |
  | Raspberry Pi 4 | Pi-hole + Unbound, WireGuard VPN, Uptime Kuma | ~300–400 MB |
  | Raspberry Pi 3 B/B+ | Single lightweight service (DNS only) | ~150 MB |
  | Raspberry Pi Zero 2 W | Testing only (512 MB, WiFi-only) | Not recommended |
  | Libre Computer Le Potato | 4K media decode candidate, 100 Mbps limit | Niche |
  | Synology DS415+ NAS | Vaultwarden, Jellyfin (storage-focused) | Existing |

  Expected savings by offloading to Pi 4: 500–600 MB RAM freed on N97, 15–20% CPU reduction, +6–9 W power.

**Modified files:**
- `Caddyfile.subdomain.example` — Added `nvr.pochita.synology.me` with Authelia protection and WebSocket support
- `assets/config.yml` — Added "Home Security" category with Frigate NVR
- `docs/HOMER_CUSTOMIZATION.md` — Updated with Frigate references

---

### 8. claude/image-notes-qr-sharing-VFfib

**Parent:** main (independent) · **Unique commits:** 3 · **Lines added:** ~3,321 across 22 files

**Purpose:** A complete image notes web application with QR code generation for sharing and direct-from-browser NIIMBOT B21 thermal label printing via Web Bluetooth.

This is the only branch that introduces **new Go application code** (as opposed to infrastructure/docs). It adds a full-stack feature to the Go server.

**New Go packages:**
- `internal/notes/handlers.go` (508 lines) — 10 HTTP handlers: list, create, view, edit, update, delete, QR generation, file serving, upload (HTMX), login
- `internal/notes/store.go` (288 lines) — File-based storage: `data/notes/{8-char-uuid}/note.md` with YAML frontmatter + attachments
- `internal/auth/middleware.go` (78 lines) — Multi-method auth: Authelia forward-auth, Bearer token, query param, session cookies

**New frontend:**
- `static/js/niimbot.js` (714 lines) — Complete NIIMBOT B21 label printer protocol over Web Bluetooth:
  - BLE service/characteristic UUIDs
  - Packet framing with XOR checksum
  - Image encoding: 90° rotation, 8-pixel bit-packing, row compression
  - Label renderer: title, QR code, tags, date, short URL
  - Configurable sizes: 40×30, 50×30, 40×60, 30×20 mm
  - PNG download fallback for non-BLE browsers

**New templates (6):**
- `templates/notes_list.html` — Dashboard listing all notes
- `templates/notes_new.html` — Creation form with camera capture
- `templates/notes_view.html` (611 lines) — View with QR display + print modal
- `templates/notes_edit.html` — Edit form
- `templates/notes_login.html` — Token-based login

**New dependencies (go.mod):**
- `github.com/skip2/go-qrcode` — QR code PNG generation (512×512, medium ECC)
- `github.com/google/uuid` — Short UUID generation
- `github.com/yuin/goldmark` + `goldmark-meta` — Markdown rendering with YAML frontmatter

**Modified files:** `cmd/webserver/main.go` (route registration), `Dockerfile` (multi-stage for new packages), `docker-compose.dev.yml` / `docker-compose.prod.yml` (notes data volume), `Caddyfile` (`/notes/*` and `/static/*` routes), `Makefile`, `.gitignore`, `.env.example`

**Routes added:** `/notes`, `/notes/new`, `/notes/create`, `/notes/upload`, `/notes/view/{id}`, `/notes/edit/{id}`, `/notes/update/{id}`, `/notes/delete/{id}`, `/notes/qr/{id}`, `/notes/files/{id}/{filename}`, `/notes/login`, `/static/*`

---

## Cumulative Service Inventory

The table below shows every service introduced across the branch chain, the branch that introduces it, and its access point at the `frigate-nvr` branch tip (the most complete G6FQx branch).

| Service | Introduced In | Subdomain / Path | Port |
|---------|--------------|-------------------|------|
| Homer Dashboard | main | `pochita.synology.me` | 8080 |
| Go/HTMX Server | main | `/test`, `/test/time`, `/test/api`, `/test/docs` | 8080 |
| Caddy Reverse Proxy | main | — (entry point) | 80/443 |
| Uptime Kuma | review-home-server | `status.pochita.synology.me` | 3001 |
| WireGuard VPN | review-home-server | — (UDP) | 51820 |
| Authelia SSO | review-home-server | `auth.pochita.synology.me` | 9091 |
| Redis | review-home-server | — (internal) | 6379 |
| Prometheus | review-home-server | — (internal) | 9090 |
| Grafana | review-home-server | `grafana.pochita.synology.me` | 3000 |
| Node Exporter | review-home-server | — (internal) | 9100 |
| cAdvisor | review-home-server | — (internal) | 8080 |
| Kopia Backup | review-home-server | `backup.pochita.synology.me` | 51515 |
| Watchtower | review-home-server | — (daemon) | — |
| Dockge | review-home-server | `dockge.pochita.synology.me` | 5001 |
| Pi-hole | privacy-infrastructure | `dns.pochita.synology.me` | 53/8053 |
| Unbound | privacy-infrastructure | — (internal, upstream of Pi-hole) | 5335 |
| Searxng | privacy-infrastructure | `search.pochita.synology.me` | 8080 |
| FreshRSS | privacy-infrastructure | `rss.pochita.synology.me` | 80 |
| Vaultwarden | vaultwarden-docs | `vault.pochita.synology.me` | 8100 |
| Jellyfin | vaultwarden-docs (config) | `media.pochita.synology.me` | 8096 |
| Synology DSM | synology-nas-docs (config) | `nas.pochita.synology.me` | 5001 |
| Frigate NVR | frigate-nvr | `nvr.pochita.synology.me` | 5000 |
| Notes App | image-notes-qr (independent) | `/notes/*` | 8080 |

**Total: 23 services/components**

---

## Infrastructure Layer Diagram

```
┌──────────────────────────────────────────────────────────────────────┐
│                         INTERNET                                     │
│                     pochita.synology.me                               │
└──────────────────────────┬───────────────────────────────────────────┘
                           │ :80/:443
                    ┌──────▼──────┐
                    │    Caddy    │  SSL termination, subdomain routing,
                    │  (main)    │  security headers, Authelia forward-auth
                    └──────┬──────┘
          ┌────────────────┼────────────────────────────────────┐
          │                │                                    │
  ┌───────▼────────┐ ┌────▼─────────────┐  ┌──────────────────▼──────────┐
  │  N97 PC        │ │  Synology DS415+ │  │  Raspberry Pi Fleet         │
  │  (Main Server) │ │  (NAS)           │  │                             │
  │                │ │                  │  │  Pi 5 + Hailo: AI detection │
  │  Homer         │ │  Vaultwarden     │  │  Pi 4: Pi-hole, Unbound,   │
  │  Go/HTMX       │ │  Jellyfin        │  │        WireGuard, Uptime   │
  │  Authelia+Redis│ │  DSM Admin       │  │        Kuma                │
  │  Prometheus    │ │  Storage volumes  │  │  Pi 3: DNS fallback        │
  │  Grafana       │ │                  │  │                             │
  │  Searxng       │ └──────────────────┘  └─────────────────────────────┘
  │  FreshRSS      │
  │  Frigate NVR   │
  │  Kopia         │
  │  Dockge        │
  │  Watchtower    │
  │  Node Exporter │
  │  cAdvisor      │
  └────────────────┘

  ┌──────────────────────────────────────────────────────────────┐
  │  DNS Flow: Device → Pi-hole (Pi 4) → Unbound → Root Servers │
  └──────────────────────────────────────────────────────────────┘
```

---

## Resource Budget

### Cumulative RAM by branch (G6FQx chain)

| Branch | New RAM | Running Total |
|--------|:-------:|:-------------:|
| main (Caddy + Homer + Go) | ~512 MB | ~512 MB |
| review-home-server (all 4 phases) | ~2,880 MB | ~3,392 MB |
| privacy-infrastructure | ~1,664 MB | ~5,056 MB |
| vaultwarden-docs (on Synology) | ~256 MB* | ~5,312 MB |
| frigate-nvr (on N97 + Pi 5) | ~2,048 MB* | ~7,360 MB |

*\* Runs on separate hardware; does not consume N97 RAM directly.*

### Distributed architecture target allocation

| Device | Services | RAM Budget |
|--------|----------|:----------:|
| N97 PC | Caddy, Homer, Go, Authelia, Redis, Prometheus, Grafana, Searxng, FreshRSS, Frigate, Kopia, Dockge, Watchtower, exporters | ~6 GB |
| Synology DS415+ | Vaultwarden, Jellyfin | ~512 MB |
| Raspberry Pi 5 | Hailo AI detection | Dedicated |
| Raspberry Pi 4 | Pi-hole, Unbound, WireGuard, Uptime Kuma | ~400 MB |

---

## Merge Strategy & Recommendations

### G6FQx Chain

Because the six G6FQx branches form a strict linear chain where each branch includes all commits from its parent, **the tip branch (`claude/frigate-nvr-G6FQx`) contains the complete state of all six branches**. Merging strategy options:

**Option A — Merge tip only (recommended)**
```bash
git checkout main
git merge claude/frigate-nvr-G6FQx
```
This brings in all 30 commits from the entire chain in one merge. Simple and complete.

**Option B — Sequential merges (for granular history)**
```bash
git checkout main
git merge claude/review-home-server-G6FQx
git merge claude/synology-nas-docs-G6FQx
git merge claude/privacy-infrastructure-G6FQx
git merge claude/homer-enhancements-G6FQx
git merge claude/vaultwarden-docs-G6FQx
git merge claude/frigate-nvr-G6FQx
```
Each merge is a no-op after the first conflict-free one since later branches already contain earlier commits. This approach only adds value if you want individual merge commits as historical markers.

### VFfib Branch

`claude/image-notes-qr-sharing-VFfib` is independent and modifies core application files (`cmd/webserver/main.go`, `Dockerfile`, `docker-compose.*.yml`, `Caddyfile`, `go.mod`). It will likely have **merge conflicts** with the G6FQx chain in these files:

| File | Conflict Source |
|------|----------------|
| `cmd/webserver/main.go` | Both branches modify route registration |
| `Caddyfile` | Both add new route blocks |
| `docker-compose.prod.yml` | Both add resource limits and services |
| `docker-compose.dev.yml` | Both modify service definitions |
| `Dockerfile` | Both modify build stages |
| `assets/config.yml` | Both add Homer service entries |
| `Makefile` | Both add new targets |
| `authelia/configuration.yml` | Both create this file with different rules |

**Recommended merge order:**
1. Merge `claude/frigate-nvr-G6FQx` into main first (infrastructure foundation)
2. Then merge `claude/image-notes-qr-sharing-VFfib` and resolve conflicts (application feature on top of infrastructure)

### Branch cleanup

After merging, all seven Claude branches can be safely deleted — the intermediate G6FQx branches are fully subsumed by the tip, and the VFfib branch will be merged independently.

---

## Appendix: Commit Inventory

### G6FQx Chain (30 unique commits, cumulative at tip)

```
# review-home-server (9 commits)
59a4d84 Add resource limits and health monitoring to all production services
188b131 Add reverse proxy routes and enable CSP for enhancement services
2a05166 Add Uptime Kuma monitoring stack (Phase 1)
37589d5 Add automated configuration backup script (Phase 1)
f73f93c Add WireGuard VPN and Authelia SSO stack (Phase 2)
b7e19fe Add Prometheus and Grafana observability stack (Phase 3)
8f00f71 Add Kopia, Watchtower, and Dockge automation stack (Phase 4)
aa9a642 Add make commands for managing FOSS enhancement stacks
8c6b99f Add comprehensive implementation documentation

# synology-nas-docs (+4 commits)
a1b9347 Add network architecture documentation for Synology NAS setup
e20f7d0 Add comprehensive Synology NAS deployment guide
8977cfe Add Synology NAS quick start reference guide
1a4d9ec Add subdomain-based Caddyfile configuration example

# privacy-infrastructure (+4 commits)
80b993a Add privacy infrastructure stack with Pi-hole, Unbound, Searxng, and FreshRSS
70b5729 Add subdomain routes for privacy infrastructure services
b1ffae8 Add comprehensive privacy infrastructure documentation
b3d5d2e Add stack-privacy make command for privacy infrastructure

# homer-enhancements (+2 commits)
a867fd5 Enhance Homer dashboard with privacy-focused organization
b010893 Add comprehensive Homer dashboard customization guide

# vaultwarden-docs (+3 commits)
a205002 Add comprehensive Vaultwarden deployment guide for Synology DS415+
f07acad Add Vaultwarden subdomain routing to Caddyfile example
c622a9e Update Homer dashboard to include Vaultwarden

# frigate-nvr (+8 commits)
6ccb1ff Add comprehensive Frigate NVR deployment guide
f05e085 Add Frigate NVR subdomain routing to Caddyfile
eaf22d6 Update Homer dashboard to include Frigate NVR
accb783 Add comprehensive distributed architecture guide
f228cd2 Add Raspberry Pi 3 Model B/B+ support to distributed architecture
d6841c2 Add Raspberry Pi Zero 2 W analysis to distributed architecture
4e10125 Add Libre Computer Le Potato analysis to distributed architecture
2752da2 Add comprehensive Jellyfin server analysis for Le Potato
```

### VFfib Branch (3 unique commits)

```
542e22e Add image notes with QR code sharing feature
f745ee7 Add NIIMBOT B21 label printer integration via Web Bluetooth
f3267e4 Fix cross-branch integration for notes feature
```
