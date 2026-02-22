# Gitea + Woodpecker CI Setup Guide

Self-hosted Git service and CI/CD for The Club dashboard.

## Architecture

```
                    ┌─────────────────────────────┐
                    │          Caddy               │
                    │    (Reverse Proxy + TLS)     │
                    │                              │
                    │  /gitea/* → gitea:3000       │
                    │  /ci/*    → woodpecker:8000  │
                    └──────┬──────────┬────────────┘
                           │          │
              ┌────────────┘          └─────────────┐
              ▼                                     ▼
     ┌─────────────────┐              ┌──────────────────────┐
     │     Gitea        │  OAuth2     │  Woodpecker Server   │
     │  (Git Service)   │◄───────────►│    (CI/CD Server)    │
     │  Port 3000       │  Webhooks   │    Port 8000         │
     │  SSH  22→2222    │────────────►│                      │
     └─────────────────┘              └──────────┬───────────┘
              │                                  │ gRPC :9000
              │                                  ▼
              │                       ┌──────────────────────┐
              │                       │  Woodpecker Agent    │
              │                       │  (Pipeline Runner)   │
              │                       │  Docker Socket       │
              │                       └──────────────────────┘
              │
         SQLite DB
         (embedded)
```

**Resource estimates** (added to existing stack):

| Service            | RAM   | CPU   |
|--------------------|-------|-------|
| Gitea              | ~200MB | 0.25 |
| Woodpecker Server  | ~100MB | 0.10 |
| Woodpecker Agent   | ~100MB | 0.25 |
| **Total added**    | **~400MB** | **0.60** |

## Quick Start

### 1. Generate Secrets

```bash
make gitea-setup
```

This creates a `.env` file from `.env.example` and generates the `WOODPECKER_AGENT_SECRET`.

### 2. Start the Stack

```bash
make docker-dev
```

### 3. Configure Gitea

1. Open `https://localhost/gitea` in your browser
2. Register the first user account (this becomes the admin)
3. Create a test repository to verify everything works

### 4. Create OAuth2 Application for Woodpecker

1. In Gitea, go to **Site Administration** > **Applications**
2. Click **Create a new OAuth2 Application**
3. Fill in:
   - **Application Name**: `Woodpecker CI`
   - **Redirect URI**: `https://localhost/ci/authorize`
4. Click **Create Application**
5. Copy the **Client ID** and **Client Secret**

### 5. Configure Woodpecker

Edit your `.env` file with the OAuth2 credentials:

```bash
WOODPECKER_GITEA_CLIENT=<paste-client-id>
WOODPECKER_GITEA_SECRET=<paste-client-secret>
```

### 6. Restart the Stack

```bash
make docker-stop
make docker-dev
```

### 7. Activate Woodpecker

1. Open `https://localhost/ci`
2. Log in with your Gitea account (OAuth2 redirect)
3. Activate your repository in the Woodpecker dashboard
4. Woodpecker will automatically create a webhook in Gitea

## Git over SSH

Clone repositories using SSH on port 2222:

```bash
git clone ssh://git@localhost:2222/username/repo.git
```

Add your SSH key in Gitea: **Settings** > **SSH / GPG Keys** > **Add Key**.

## CI Pipeline

The `.woodpecker.yml` in the repository root defines the CI pipeline. It runs on every push and pull request:

| Step       | Description                        |
|------------|------------------------------------|
| build      | Compiles the Go binary             |
| test       | Runs tests with race detection     |
| vet        | Static analysis with `go vet`      |
| fmt-check  | Verifies `gofmt` formatting        |

### Custom Pipeline Steps

Add steps to `.woodpecker.yml`:

```yaml
steps:
  - name: my-step
    image: alpine:latest
    commands:
      - echo "running custom step"
```

See [Woodpecker CI docs](https://woodpecker-ci.org/docs/usage/pipeline-syntax) for full syntax.

## Makefile Commands

| Command          | Description                              |
|------------------|------------------------------------------|
| `make gitea-setup` | Generate secrets and show setup steps  |
| `make gitea-logs`  | View Gitea container logs              |
| `make ci-logs`     | View Woodpecker server + agent logs    |
| `make health`      | Health check all services incl. Gitea  |

## Service URLs

| Service          | Dev URL                        | Direct Port |
|------------------|--------------------------------|-------------|
| Gitea Web UI     | `https://localhost/gitea`      | `3000`      |
| Gitea SSH        | `ssh://localhost:2222`         | `2222`      |
| Woodpecker CI    | `https://localhost/ci`         | `8000`      |
| Homer Dashboard  | `https://localhost`            | `8080`      |

## Production Configuration

For production, update `.env`:

```bash
GITEA_ROOT_URL=https://yourdomain.com/gitea/
GITEA_DOMAIN=yourdomain.com
GITEA_SSH_DOMAIN=yourdomain.com
WOODPECKER_HOST=https://yourdomain.com/ci
```

Production differences from dev:
- Registration disabled (`DISABLE_REGISTRATION=true`)
- Sign-in required to view repos (`REQUIRE_SIGNIN_VIEW=true`)
- Open signup disabled in Woodpecker (`WOODPECKER_OPEN=false`)
- Resource limits enforced (512MB Gitea, 256MB each Woodpecker component)
- Health checks enabled
- `no-new-privileges` security option set

Deploy with:

```bash
make docker-prod
```

## Troubleshooting

### Gitea returns 502 Bad Gateway
Gitea may need a moment to start. Check logs:
```bash
make gitea-logs
```

### Woodpecker shows "unauthorized" or can't connect to Gitea
- Verify OAuth2 credentials in `.env` match what Gitea shows
- Ensure `WOODPECKER_GITEA_URL=http://gitea:3000` (internal Docker network)
- Restart after changing `.env`: `make docker-stop && make docker-dev`

### Woodpecker pipelines aren't triggering
- Check that the repository is **activated** in Woodpecker
- Verify the webhook was created in Gitea (repo Settings > Webhooks)
- Check agent logs: `make ci-logs`

### SSH clone fails
- Verify SSH key is added in Gitea user settings
- Use port 2222: `git clone ssh://git@localhost:2222/user/repo.git`
- Check Gitea SSH is running: `ssh -T git@localhost -p 2222`
