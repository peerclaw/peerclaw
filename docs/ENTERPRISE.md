**English** | [中文](ENTERPRISE_zh.md)

# PeerClaw Enterprise Deployment Guide

Deploy PeerClaw on your internal network for secure agent-to-agent communication — no external dependencies required.

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Corporate Network                   │
│                                                     │
│   ┌──────────────────┐                              │
│   │  peerclaw-server │  ← Single instance           │
│   │  :8080           │    (discovery + signaling)    │
│   └────────┬─────────┘                              │
│            │                                        │
│   ┌────────┼────────────────────┐                   │
│   │        │        │           │                   │
│   ▼        ▼        ▼           ▼                   │
│ Agent A  Agent B  Agent C  Agent D                  │
│ (billing) (audit) (notify) (invoice)                │
│                                                     │
│   All communication via WebRTC P2P or               │
│   server relay — no Nostr, no STUN/TURN             │
└─────────────────────────────────────────────────────┘
```

## Quick Start (5 lines of Go)

```go
package main

import (
    "context"
    agent "github.com/peerclaw/peerclaw-agent"
    "github.com/peerclaw/peerclaw-core/envelope"
)

func main() {
    // Create agent — only name, server URL, and capabilities needed.
    a, _ := agent.NewSimple("invoice-processor", "http://peerclaw.internal:8080",
        "process-invoice", "query-status",
    )

    // Pre-provision trusted contacts.
    a.ImportContacts([]string{"agent-billing", "agent-audit", "agent-notify"})

    // Register handlers.
    a.Handle("process-invoice", func(ctx context.Context, env *envelope.Envelope) (*envelope.Envelope, error) {
        return &envelope.Envelope{Payload: []byte(`{"status":"ok"}`)}, nil
    })

    // Start serving.
    ctx := context.Background()
    a.Start(ctx)
    select {} // block forever
}
```

### What `NewSimple` does

`NewSimple(name, serverURL, capabilities...)` is equivalent to:

```go
agent.New(agent.Options{
    Name:         name,
    ServerURL:    serverURL,
    Capabilities: capabilities,
})
```

All other options use safe defaults:
- **Auto-generated Ed25519 keypair** (no key file management needed)
- **Server-only mode** — no Nostr relays, no STUN/TURN
- **No serverless** — relies on the central peerclaw-server

### Pre-provisioned Trust with `ImportContacts`

In enterprise environments, agents are known in advance. Use `ImportContacts` to bulk-import trusted agent IDs:

```go
a.ImportContacts([]string{"agent-billing", "agent-audit", "agent-notify"})
```

This sets all imported agents to `TrustVerified` level, allowing immediate bidirectional communication without TOFU handshakes.

## Docker Compose Deployment

Use the provided `deploy/docker-compose.yaml` to run the server:

```bash
cd deploy

# Create .env file
cat > .env << 'EOF'
POSTGRES_USER=peerclaw
POSTGRES_PASSWORD=changeme
POSTGRES_DB=peerclaw
REDIS_PASSWORD=changeme
EOF

docker compose up -d
```

This starts:
- **peerclaw-server** on port 8080
- **PostgreSQL** for persistence
- **Redis** for cross-node signaling (optional for single-node)
- **Caddy** as reverse proxy (ports 80/443)

For SQLite-only (simplest setup), modify the server config to use SQLite and remove the PostgreSQL/Redis services.

## Security Recommendations

### Network

- Deploy peerclaw-server behind a firewall — only expose to the internal network
- Use TLS between agents and server (configure Caddy or your load balancer)
- Restrict agent → server communication to the peerclaw port only

### Identity & Keys

- Each agent auto-generates an Ed25519 keypair on first run
- For persistent identity across restarts, use `KeypairPath` in `Options`:
  ```go
  agent.New(agent.Options{
      Name:        "invoice-processor",
      ServerURL:   "http://peerclaw.internal:8080",
      KeypairPath: "/etc/peerclaw/agent.key",
  })
  ```
- Store key files with restricted permissions (`chmod 600`)

### Trust Management

- Use `ImportContacts` to pre-provision known agents at startup
- For dynamic trust, register a `ConnectionRequestHandler`:
  ```go
  a.OnConnectionRequest(func(ctx context.Context, req *agent.ConnectionRequest) bool {
      // Check against an internal registry or approve all internal agents
      return isInternalAgent(req.FromAgentID)
  })
  ```
- Block untrusted agents: `a.BlockAgent(agentID)`

## Full Example

See [agent/examples/enterprise/main.go](../agent/examples/enterprise/main.go) for a complete working example.
