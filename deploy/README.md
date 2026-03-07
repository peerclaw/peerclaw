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
