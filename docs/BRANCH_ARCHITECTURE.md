# The Club — Branch Architecture & Design Document

> **Generated:** 2026-02-14
> **Repository:** `github.com/LittleClubFoot/TheClub`
> **Base Branch:** `main` / `master`

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Naming Convention](#naming-convention)
3. [Branch Lineage Map](#branch-lineage-map)
4. [Branch Dependency Chain](#branch-dependency-chain)
5. [Branch-by-Branch Analysis](#branch-by-branch-analysis)
   - [main / master](#1-main--master)
   - [claude/infra-monitoring-security-observability-automation](#2-claudeinfra-monitoring-security-observability-automation)
   - [claude/docs-synology-nas-deployment](#3-claudedocs-synology-nas-deployment)
   - [claude/infra-privacy-dns-search-rss](#4-claudeinfra-privacy-dns-search-rss)
   - [claude/ui-homer-dashboard-privacy-redesign](#5-claudeui-homer-dashboard-privacy-redesign)
   - [claude/docs-vaultwarden-synology-deployment](#6-claudedocs-vaultwarden-synology-deployment)
   - [claude/infra-frigate-nvr-distributed-architecture](#7-claudeinfra-frigate-nvr-distributed-architecture)
   - [claude/feat-notes-app-qr-label-printer](#8-claudefeat-notes-app-qr-label-printer)
6. [File-Level Dependency & Conflict Analysis](#file-level-dependency--conflict-analysis)
7. [Independent Merge Streams](#independent-merge-streams)
8. [Cumulative Service Inventory](#cumulative-service-inventory)
9. [Infrastructure Layer Diagram](#infrastructure-layer-diagram)
10. [Resource Budget](#resource-budget)
11. [Merge Strategy & Recommendations](#merge-strategy--recommendations)

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

## Naming Convention

Branches follow the pattern: **`claude/<type>-<description>-<session>`**

| Prefix | Meaning | Example |
|--------|---------|---------|
| `infra/` | Infrastructure stacks — docker-compose files, configs, deployment scripts | `infra-monitoring-security-observability-automation` |
| `docs/` | Primarily documentation — guides, architecture docs, reference material | `docs-synology-nas-deployment` |
| `ui/` | Dashboard and frontend changes — Homer config, theme, layout | `ui-homer-dashboard-privacy-redesign` |
| `feat/` | New application features — Go code, templates, frontend JS | `feat-notes-app-qr-label-printer` |

### Branch Name Mapping

| Old Name | New Name | Type |
|----------|----------|------|
| `claude/review-home-server-G6FQx` | `claude/infra-monitoring-security-observability-automation-WOqSv` | infra |
| `claude/synology-nas-docs-G6FQx` | `claude/docs-synology-nas-deployment-WOqSv` | docs |
| `claude/privacy-infrastructure-G6FQx` | `claude/infra-privacy-dns-search-rss-WOqSv` | infra |
| `claude/homer-enhancements-G6FQx` | `claude/ui-homer-dashboard-privacy-redesign-WOqSv` | ui |
| `claude/vaultwarden-docs-G6FQx` | `claude/docs-vaultwarden-synology-deployment-WOqSv` | docs |
| `claude/frigate-nvr-G6FQx` | `claude/infra-frigate-nvr-distributed-architecture-WOqSv` | infra |
| `claude/image-notes-qr-sharing-VFfib` | `claude/feat-notes-app-qr-label-printer-WOqSv` | feat |

---

## Branch Lineage Map

There are two distinct development lineages: a **linear chain** of six infrastructure/docs/UI branches that build incrementally on each other, and one **independent feature branch** with application code.

```
main (2 commits: LICENSE + initial)
│
├─── LINEAR CHAIN ────────────────────────────────────────────────────────────
│    │
│    ├─ claude/infra-monitoring-security-observability-automation  (+9 commits)
│    │   └─ claude/docs-synology-nas-deployment                   (+4 commits)
│    │       └─ claude/infra-privacy-dns-search-rss               (+4 commits)
│    │           └─ claude/ui-homer-dashboard-privacy-redesign     (+2 commits)
│    │               └─ claude/docs-vaultwarden-synology-deployment    (+3 commits)
│    │                   └─ claude/infra-frigate-nvr-distributed-architecture  (+8 commits)
│
├─── INDEPENDENT BRANCH ─────────────────────────────────────────────────────
│    │
│    └─ claude/feat-notes-app-qr-label-printer                    (+3 commits)
│
```

The linear chain is strictly additive — each branch contains all commits from the branch above it plus its own new commits. The feat branch is based directly on `main` and has no shared commits with the chain.

---

## Branch Dependency Chain

| # | Branch | Parent | Unique Commits | Type | Primary Focus |
|---|--------|--------|:-:|:---:|---|
| 1 | `main` / `master` | — | 2 | — | Bare repository with LICENSE |
| 2 | `claude/infra-monitoring-security-observability-automation` | main | 9 | infra | FOSS enhancement stacks (4 phases) |
| 3 | `claude/docs-synology-nas-deployment` | #2 | 4 | docs | Synology NAS deployment docs |
| 4 | `claude/infra-privacy-dns-search-rss` | #3 | 4 | infra | Pi-hole, Unbound, Searxng, FreshRSS |
| 5 | `claude/ui-homer-dashboard-privacy-redesign` | #4 | 2 | ui | Dashboard redesign (privacy theme) |
| 6 | `claude/docs-vaultwarden-synology-deployment` | #5 | 3 | docs | Self-hosted password manager |
| 7 | `claude/infra-frigate-nvr-distributed-architecture` | #6 | 8 | infra | AI security cameras + distributed arch |
| 8 | `claude/feat-notes-app-qr-label-printer` | main | 3 | feat | Notes app + QR codes + label printer |

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

### 2. claude/infra-monitoring-security-observability-automation

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

### 3. claude/docs-synology-nas-deployment

**Parent:** #2 infra-monitoring · **Unique commits:** 4 · **Lines added:** ~1,952

**Purpose:** Comprehensive documentation for deploying The Club alongside a Synology DS415+ NAS using subdomain-based routing under `pochita.synology.me`.

**New files:**
- `docs/NETWORK_ARCHITECTURE.md` (464 lines) — Network topology, DNS strategies (Cloudflare / DuckDNS / local), port forwarding options, SSL/TLS certificate strategy, physical deployment options
- `docs/SYNOLOGY_DEPLOYMENT.md` (855 lines) — 4-phase deployment guide: DNS prep → base stack → enhancement stacks → service protection & maintenance
- `docs/SYNOLOGY_QUICK_START.md` (289 lines) — Condensed 5-minute setup reference
- `Caddyfile.subdomain.example` (344 lines) — Production-ready Caddy config with subdomain routing for 8+ services, Authelia forward-auth, Let's Encrypt ACME

**Key architectural decision:** Subdomain routing (`status.pochita.synology.me`, `grafana.pochita.synology.me`, etc.) instead of path-based routing — cleaner URLs, per-service SSL, better isolation.

---

### 4. claude/infra-privacy-dns-search-rss

**Parent:** #3 docs-synology-nas · **Unique commits:** 4 · **Lines added:** ~1,418

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

### 5. claude/ui-homer-dashboard-privacy-redesign

**Parent:** #4 infra-privacy · **Unique commits:** 2 · **Lines added:** ~782

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

### 6. claude/docs-vaultwarden-synology-deployment

**Parent:** #5 ui-homer · **Unique commits:** 3 · **Lines added:** ~1,066

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

### 7. claude/infra-frigate-nvr-distributed-architecture

**Parent:** #6 docs-vaultwarden · **Unique commits:** 8 · **Lines added:** ~3,322

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

### 8. claude/feat-notes-app-qr-label-printer

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

## File-Level Dependency & Conflict Analysis

### File Overlap Matrix

This matrix shows which branches modify each shared file. Only files touched by **more than one branch** are listed. Each branch's changes are relative to its parent (unique modifications only).

| File | Br 2 (infra) | Br 3 (docs) | Br 4 (privacy) | Br 5 (ui) | Br 6 (vault) | Br 7 (frigate) | Br 8 (notes) | Conflict Risk |
|------|:---:|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| `Caddyfile` | **W** | | | | | | **W** | LOW |
| `Caddyfile.subdomain.example` | | **C** | **W** | | **W** | **W** | | NONE |
| `Makefile` | **W** | | **W** | | | | **W** | LOW |
| `docker-compose.prod.yml` | **W** | | | | | | **W** | NONE |
| `assets/config.yml` | | | | **W** | **W** | **W** | **W** | NONE |
| `authelia/configuration.yml` | **C** | | | | | | **C** | LOW |
| `authelia/users_database.yml` | **C** | | | | | | **C** | NONE |
| `docs/HOMER_CUSTOMIZATION.md` | | | | **C** | **W** | **W** | | NONE |

**C** = creates new file · **W** = modifies existing file

### Per-File Conflict Detail

#### `Caddyfile` — Branches 2 & 8 (LOW risk)

| Branch | Lines Modified | What's Added |
|--------|---------------|--------------|
| Br 2 (infra) | ~lines 27-49, 112-116 | 5 routes: `/status`, `/grafana`, `/auth`, `/dockge`, `/backup` + CSP enhancement |
| Br 8 (notes) | ~lines 27-37 | 2 routes: `/notes/*`, `/static/*` |

**Verdict:** Different routes appended to same block. Trivially mergeable — keep both route sets.

#### `Caddyfile.subdomain.example` — Branches 3, 4, 6, 7 (NONE)

Each branch appends new subdomain blocks to the end of the file. Branch 3 creates the base (344 lines), Branch 4 adds `dns.`/`search.`/`rss.` blocks (+70 lines), Branch 6 adds `vault.` (+59 lines), Branch 7 adds `nvr.` (+52 lines). Linear growth, no overlapping sections.

#### `Makefile` — Branches 2, 4, 8 (LOW risk)

| Branch | What's Added |
|--------|--------------|
| Br 2 (infra) | `stack-monitoring`, `stack-security`, `stack-observability`, `stack-automation`, `stack-all`, `stack-stop`, `backup` targets |
| Br 4 (privacy) | `stack-privacy` target, inserts into Br 2's `stack-all` and `stack-stop` |
| Br 8 (notes) | `notes-token` target, health check enhancement at line ~194 |

**Verdict:** Br 2 and Br 8 both touch the health check section (~line 194) with identical content. Br 4 cleanly extends Br 2's targets. Resolution: keep Br 2 health check, append Br 4 targets, append Br 8 `notes-token` target.

#### `docker-compose.prod.yml` — Branches 2 & 8 (NONE)

| Branch | What's Modified |
|--------|----------------|
| Br 2 (infra) | Caddy metrics port, caddy/homer resource limits |
| Br 8 (notes) | go-test-server volumes, env vars, increased resource limits, notes-data volume |

**Verdict:** Different service sections. No overlapping lines.

#### `assets/config.yml` — Branches 5, 6, 7, 8 (NONE)

| Branch | What's Modified |
|--------|----------------|
| Br 5 (ui) | Complete restructure: 8 categories, purple theme, 14+ services |
| Br 6 (vault) | +8 lines: Vaultwarden entry in Security section |
| Br 7 (frigate) | +18 lines: Home Security category with Frigate NVR |
| Br 8 (notes) | +6 lines: Notes entry in Development section |

**Verdict:** Br 5 sets the structure; Br 6, 7, 8 each insert into different sections. Clean sequential merge.

#### `authelia/configuration.yml` — Branches 2 & 8 (LOW risk)

Both branches create this file from scratch. Branch 8's version is a superset of Branch 2's — it includes the same base configuration plus additional notes-specific access control rules (`/notes/view/*` bypass, `/notes/*` one-factor). **Resolution:** Use Branch 8's version (it includes everything from Branch 2 plus notes rules).

#### `authelia/users_database.yml` — Branches 2 & 8 (NONE)

Both branches create this file with **identical content** (29 lines, same admin user hash).

---

## Independent Merge Streams

### Can the linear chain be split?

Even though branches 2–7 form a linear commit chain, we can analyze **file-level** and **content-level** dependencies to identify potential parallel streams.

**File-level independence** (do branches modify different files?):

| Branch pair | Shared files | Independent? |
|-------------|-------------|:---:|
| Br 2 (infra) ↔ Br 3 (docs) | None | **YES** |
| Br 2 (infra) ↔ Br 8 (notes) | Caddyfile, Makefile, docker-compose.prod.yml, authelia/* | Low conflict |
| Br 3 (docs) ↔ Br 8 (notes) | None | **YES** |
| Br 4 (privacy) ↔ Br 2 (infra) | Makefile | Depends on Br 2 |
| Br 4 (privacy) ↔ Br 3 (docs) | Caddyfile.subdomain.example | Depends on Br 3 |
| Br 5 (ui) ↔ Br 6 (vault) | assets/config.yml, HOMER_CUSTOMIZATION.md | Different sections |
| Br 5 (ui) ↔ Br 7 (frigate) | assets/config.yml, HOMER_CUSTOMIZATION.md | Different sections |
| Br 6 (vault) ↔ Br 7 (frigate) | Caddyfile.subdomain.example, assets/config.yml, HOMER_CUSTOMIZATION.md | Different sections |

**Content-level dependencies** (does branch content reference earlier services?):

| Branch | Cross-references to earlier branches |
|--------|--------------------------------------|
| Br 3 (docs) | 22 references to Br 2 services (uptime-kuma, grafana, authelia, etc.) in `Caddyfile.subdomain.example` |
| Br 4 (privacy) | Inserts `stack-privacy` into Br 2's `stack-all` Makefile target |
| Br 5 (ui) | Lists **all** services from Br 2, 3, 4 in Homer dashboard config |
| Br 6 (vault) | 38 references to earlier services in VAULTWARDEN_SYNOLOGY.md |
| Br 7 (frigate) | 189 references to all prior services in DISTRIBUTED_ARCHITECTURE.md |
| Br 8 (notes) | No content references to the chain. Fully independent. |

### Viable Stream Map

Based on the analysis above, the branches can be organized into **three parallel streams** off `main`, with a sequential merge phase afterward:

```
                              main
                           ┌────┼──────────────────────────┐
                           │    │                           │
                  Stream A │    │ Stream B                  │ Stream C
                  (infra)  │    │ (docs)                    │ (app feature)
                  no file  │    │ no file overlap           │ low-conflict
                  overlap  │    │ with A                    │ with A
                  with B   │    │                           │
                           ▼    ▼                           ▼
                        Br 2   Br 3                      Br 8
                  infra-monitoring  docs-synology       feat-notes-app
                        │    │                           (merge anytime)
                        └──┬─┘
                           │  merge A + B, then:
                           ▼
                         Br 4
                  infra-privacy-dns
                  (needs Br 2 Makefile +
                   Br 3 Caddyfile.subdomain)
                           │
                           ▼
                         Br 5
                  ui-homer-dashboard
                  (references all services
                   from Br 2, 3, 4)
                           │
                      ┌────┴────┐
                      ▼         ▼
                    Br 6      Br 7     ← PARALLEL (different sections
                    docs-     infra-      of same 3 files)
                    vault     frigate
```

### Stream Details

#### Stream A: Core Infrastructure (`infra-monitoring-security-observability-automation`)

- **Branch:** #2
- **From:** `main`
- **File touches:** `Caddyfile`, `Makefile`, `docker-compose.prod.yml`, `authelia/*`, + 11 new files
- **Can merge independently:** YES — no file overlap with Stream B

#### Stream B: Synology Documentation (`docs-synology-nas-deployment`)

- **Branch:** #3
- **From:** `main`
- **File touches:** `Caddyfile.subdomain.example` (new), + 3 new doc files
- **Can merge independently:** YES — all new files, no overlap with Stream A
- **Content note:** References Stream A services in docs/configs, but no structural dependency

#### Stream C: Notes Application (`feat-notes-app-qr-label-printer`)

- **Branch:** #8
- **From:** `main`
- **File touches:** `Caddyfile`, `Makefile`, `docker-compose.prod.yml`, `assets/config.yml`, `authelia/*`, `cmd/webserver/main.go`, `Dockerfile`, + 15 new files
- **Can merge independently:** YES — all conflicts with Stream A are additive (different routes, different Makefile targets, different docker-compose sections, compatible authelia rules)
- **Best timing:** Merge after Stream A for cleanest resolution, but can also merge first

#### After Streams A + B merge → Sequential phase:

| Order | Branch | Why sequential |
|-------|--------|---------------|
| 1st | Br 4 (privacy) | Modifies Br 2's `Makefile` (adds to `stack-all`) and Br 3's `Caddyfile.subdomain.example` |
| 2nd | Br 5 (homer UI) | `assets/config.yml` references services from Br 2, 3, 4 |
| 3rd | Br 6 (vaultwarden) AND Br 7 (frigate) — **parallel** | Both add to different sections of `Caddyfile.subdomain.example`, `assets/config.yml`, `HOMER_CUSTOMIZATION.md`. No overlapping lines. |

### Conflict Resolution Cheat Sheet

When merging Stream C after the linear chain:

| File | Resolution |
|------|-----------|
| `Caddyfile` | Keep chain routes (`/status`, `/grafana`, `/auth`, `/dockge`, `/backup`) + add notes routes (`/notes/*`, `/static/*`) |
| `Makefile` | Keep chain targets (all `stack-*`) + add `notes-token` target. Health check at line ~194 is identical in both — keep once. |
| `docker-compose.prod.yml` | Keep chain resource limits (caddy, homer) + add notes volumes, env vars, go-server limits |
| `assets/config.yml` | Use chain's restructured version (Br 5) + insert Notes entry into Development section |
| `authelia/configuration.yml` | Use Br 8's version — it's a superset that includes Br 2's base rules + notes-specific rules |
| `authelia/users_database.yml` | Identical in both — no conflict |

---

## Cumulative Service Inventory

The table below shows every service introduced across the branch chain, the branch that introduces it, and its access point at the `infra-frigate-nvr-distributed-architecture` branch tip (the most complete chain branch).

| Service | Introduced In | Subdomain / Path | Port |
|---------|--------------|-------------------|------|
| Homer Dashboard | main | `pochita.synology.me` | 8080 |
| Go/HTMX Server | main | `/test`, `/test/time`, `/test/api`, `/test/docs` | 8080 |
| Caddy Reverse Proxy | main | — (entry point) | 80/443 |
| Uptime Kuma | infra-monitoring | `status.pochita.synology.me` | 3001 |
| WireGuard VPN | infra-monitoring | — (UDP) | 51820 |
| Authelia SSO | infra-monitoring | `auth.pochita.synology.me` | 9091 |
| Redis | infra-monitoring | — (internal) | 6379 |
| Prometheus | infra-monitoring | — (internal) | 9090 |
| Grafana | infra-monitoring | `grafana.pochita.synology.me` | 3000 |
| Node Exporter | infra-monitoring | — (internal) | 9100 |
| cAdvisor | infra-monitoring | — (internal) | 8080 |
| Kopia Backup | infra-monitoring | `backup.pochita.synology.me` | 51515 |
| Watchtower | infra-monitoring | — (daemon) | — |
| Dockge | infra-monitoring | `dockge.pochita.synology.me` | 5001 |
| Pi-hole | infra-privacy | `dns.pochita.synology.me` | 53/8053 |
| Unbound | infra-privacy | — (internal, upstream of Pi-hole) | 5335 |
| Searxng | infra-privacy | `search.pochita.synology.me` | 8080 |
| FreshRSS | infra-privacy | `rss.pochita.synology.me` | 80 |
| Vaultwarden | docs-vaultwarden | `vault.pochita.synology.me` | 8100 |
| Jellyfin | docs-vaultwarden (config) | `media.pochita.synology.me` | 8096 |
| Synology DSM | docs-synology (config) | `nas.pochita.synology.me` | 5001 |
| Frigate NVR | infra-frigate | `nvr.pochita.synology.me` | 5000 |
| Notes App | feat-notes (independent) | `/notes/*` | 8080 |

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

### Cumulative RAM by branch (linear chain)

| Branch | New RAM | Running Total |
|--------|:-------:|:-------------:|
| main (Caddy + Homer + Go) | ~512 MB | ~512 MB |
| infra-monitoring (all 4 phases) | ~2,880 MB | ~3,392 MB |
| infra-privacy | ~1,664 MB | ~5,056 MB |
| docs-vaultwarden (on Synology) | ~256 MB* | ~5,312 MB |
| infra-frigate (on N97 + Pi 5) | ~2,048 MB* | ~7,360 MB |

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

### Recommended Merge Order

```
Phase 1 — Parallel (3 streams, zero conflicts)
├── git merge claude/infra-monitoring-security-observability-automation-WOqSv   (Stream A)
├── git merge claude/docs-synology-nas-deployment-WOqSv                         (Stream B)
└── git merge claude/feat-notes-app-qr-label-printer-WOqSv                     (Stream C)

Phase 2 — Sequential (depends on A + B)
└── git merge claude/infra-privacy-dns-search-rss-WOqSv

Phase 3 — Sequential (depends on Phase 2)
└── git merge claude/ui-homer-dashboard-privacy-redesign-WOqSv

Phase 4 — Parallel (different sections of same files)
├── git merge claude/docs-vaultwarden-synology-deployment-WOqSv
└── git merge claude/infra-frigate-nvr-distributed-architecture-WOqSv
```

**In practice**, since the linear chain already contains all commits at the tip, the simplest approach is:

1. Merge `claude/infra-frigate-nvr-distributed-architecture-WOqSv` into main (gets entire chain)
2. Merge `claude/feat-notes-app-qr-label-printer-WOqSv` and resolve the 6 low-risk file conflicts

### Branch cleanup

After merging, all branches can be safely deleted:
- The 5 intermediate chain branches are fully subsumed by the tip (`infra-frigate-nvr-distributed-architecture`)
- The feat branch will be merged independently
- Old-name branches (with G6FQx/VFfib suffixes) are aliases pointing to the same commits as the renamed branches

---

## Appendix: Commit Inventory

### Linear Chain (30 unique commits, cumulative at tip)

```
# infra-monitoring-security-observability-automation (9 commits)
59a4d84 Add resource limits and health monitoring to all production services
188b131 Add reverse proxy routes and enable CSP for enhancement services
2a05166 Add Uptime Kuma monitoring stack (Phase 1)
37589d5 Add automated configuration backup script (Phase 1)
f73f93c Add WireGuard VPN and Authelia SSO stack (Phase 2)
b7e19fe Add Prometheus and Grafana observability stack (Phase 3)
8f00f71 Add Kopia, Watchtower, and Dockge automation stack (Phase 4)
aa9a642 Add make commands for managing FOSS enhancement stacks
8c6b99f Add comprehensive implementation documentation

# docs-synology-nas-deployment (+4 commits)
a1b9347 Add network architecture documentation for Synology NAS setup
e20f7d0 Add comprehensive Synology NAS deployment guide
8977cfe Add Synology NAS quick start reference guide
1a4d9ec Add subdomain-based Caddyfile configuration example

# infra-privacy-dns-search-rss (+4 commits)
80b993a Add privacy infrastructure stack with Pi-hole, Unbound, Searxng, and FreshRSS
70b5729 Add subdomain routes for privacy infrastructure services
b1ffae8 Add comprehensive privacy infrastructure documentation
b3d5d2e Add stack-privacy make command for privacy infrastructure

# ui-homer-dashboard-privacy-redesign (+2 commits)
a867fd5 Enhance Homer dashboard with privacy-focused organization
b010893 Add comprehensive Homer dashboard customization guide

# docs-vaultwarden-synology-deployment (+3 commits)
a205002 Add comprehensive Vaultwarden deployment guide for Synology DS415+
f07acad Add Vaultwarden subdomain routing to Caddyfile example
c622a9e Update Homer dashboard to include Vaultwarden

# infra-frigate-nvr-distributed-architecture (+8 commits)
6ccb1ff Add comprehensive Frigate NVR deployment guide
f05e085 Add Frigate NVR subdomain routing to Caddyfile
eaf22d6 Update Homer dashboard to include Frigate NVR
accb783 Add comprehensive distributed architecture guide
f228cd2 Add Raspberry Pi 3 Model B/B+ support to distributed architecture
d6841c2 Add Raspberry Pi Zero 2 W analysis to distributed architecture
4e10125 Add Libre Computer Le Potato analysis to distributed architecture
2752da2 Add comprehensive Jellyfin server analysis for Le Potato
```

### Independent Feature Branch (3 unique commits)

```
542e22e Add image notes with QR code sharing feature
f745ee7 Add NIIMBOT B21 label printer integration via Web Bluetooth
f3267e4 Fix cross-branch integration for notes feature
```
