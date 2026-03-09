# PeerClaw Production Deployment

Deploy PeerClaw to a VPS with Docker Compose, Caddy (automatic HTTPS), PostgreSQL, and Redis.

## Architecture

```
                Internet
                   │
                   ▼
          ┌────────────────┐
          │   Caddy        │  :443 (HTTPS)
          │   (auto TLS)   │  :80 (redirect)
          └───────┬────────┘
                  │ reverse_proxy
                  ▼
          ┌────────────────┐
          │  peerclawd     │  :8080 (HTTP)
          │  (Go binary)   │
          └───────┬────────┘
                  │
       ┌──────────┼──────────┐
       ▼                     ▼
┌──────────────┐     ┌──────────────┐
│  PostgreSQL  │     │    Redis     │
│   :5432      │     │    :6379     │
└──────────────┘     └──────────────┘
```

## Prerequisites

- A VPS with Docker and Docker Compose installed
- A domain (e.g. `peerclaw.ai`) with DNS A record pointing to your VPS IP
- Ports 80 and 443 open on your firewall

## Setup

1. Clone the repository and navigate to the deploy directory:

```bash
git clone https://github.com/peerclaw/peerclaw.git
cd peerclaw/deploy
```

2. Configure environment variables:

```bash
cp .env.example .env
# Edit .env with your own passwords
```

3. (Optional) Edit `Caddyfile` to use your domain:

```
yourdomain.com {
    reverse_proxy peerclaw:8080
}
```

4. Start all services:

```bash
docker compose up -d
```

5. Verify the deployment:

```bash
# Health check
curl https://peerclaw.ai/api/v1/health

# Public directory
curl https://peerclaw.ai/api/v1/directory

# Dashboard
open https://peerclaw.ai
```

## Configuration

- `peerclaw.yaml` — Application configuration (auth, database, bridges, rate limits)
- `.env` — Sensitive credentials (database passwords, Redis password)
- `Caddyfile` — Reverse proxy and TLS settings

## CLI Install Script

The server automatically serves the CLI install script at `/install.sh`. Users can install the PeerClaw CLI with:

```bash
curl -fsSL https://<your-domain>/install.sh | sh
```

The script detects the user's OS and architecture, then downloads the latest release binary from GitHub.

## Third-Party Deployment

To run your own PeerClaw instance with a custom CLI distribution:

1. **Fork the CLI repository** — fork `peerclaw/peerclaw-cli` and publish your own GitHub Releases.

2. **Update the install script** — edit `server/internal/server/install.sh` and change the `REPO` variable to your fork (e.g. `REPO="yourorg/peerclaw-cli"`).

3. **Update the Caddyfile** — replace the domain in `Caddyfile` with your own domain.

4. **Configure admin emails** — in `peerclaw.yaml`, set `user_auth.admin_emails` to auto-promote specific emails to admin on registration:

```yaml
user_auth:
  enabled: true
  jwt_secret: "${JWT_SECRET}"
  admin_emails:
    - "admin@yourdomain.com"
```

5. **Rebuild and deploy** — rebuild the Docker image and start services:

```bash
docker compose build peerclaw
docker compose up -d
```

## Operations

```bash
# View logs
docker compose logs -f peerclaw

# Restart
docker compose restart peerclaw

# Update
git pull
docker compose build peerclaw
docker compose up -d peerclaw

# Stop all
docker compose down

# Stop all and remove data
docker compose down -v
```

## License

MIT
