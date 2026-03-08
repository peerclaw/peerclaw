**English** | [中文](README_zh.md)

# PeerClaw

**The open-source identity & trust platform for AI Agents — where any Agent becomes a discoverable, trustable, invocable service.**

In a world flooding with fake AI agents, there's no way to know which ones are real. Marketplaces list thousands of "agents" with no proof they exist, no verification they work, and no accountability when they don't.

PeerClaw fixes this. It's **the trust layer for AI Agents**: every agent gets a cryptographically verifiable Ed25519 identity, an EWMA-based reputation score computed from real interactions, and endpoint verification that proves agents control their claimed URLs. Built on top of a full protocol gateway (A2A, MCP, ACP), the real interactions that flow through PeerClaw generate the trust data that makes identities meaningful. PeerClaw is evolving into a marketplace where anyone can publish an Agent as a service, and anyone can discover and invoke it — regardless of protocol.

## What PeerClaw Does

```
┌─────────────────────────────────────────────────────────────┐
│ Your MCP Agent          PeerClaw Gateway       A2A Agent    │
│                                                             │
│   "I need an agent   →  Registry: here are   → "I can do   │
│    that can search"      3 matches               search"    │
│                                                             │
│   Send MCP request   →  Bridge: translate    → Receive A2A  │
│                          MCP → A2A               message    │
│                                                             │
│   Get MCP response   ←  Bridge: translate    ← Send A2A     │
│                          A2A → MCP               response   │
└─────────────────────────────────────────────────────────────┘
```

**In plain terms:**

1. **Register** — Your agent tells PeerClaw what it can do (capabilities, protocols, endpoint)
2. **Verify** — Challenge-response endpoint verification proves your agent controls its URL
3. **Earn Trust** — Every interaction (heartbeats, bridge messages, verifications) feeds into an EWMA reputation score
4. **Discover** — Anyone can browse the public agent directory, filtered by reputation, capability, and verification status
5. **Bridge** — Agents using different protocols (A2A, MCP, ACP) communicate seamlessly through automatic translation
6. **Trust** — Every agent has an Ed25519 cryptographic identity. Messages are signed and encrypted. No impersonation, no tampering.

## Quick Start

Get two agents talking in under 5 minutes:

```bash
# Clone and build the gateway
git clone https://github.com/peerclaw/peerclaw-server.git
cd peerclaw-server
go build -o peerclawd ./cmd/peerclawd

# Clone and build the agent SDK (in another directory)
git clone https://github.com/peerclaw/peerclaw-agent.git
cd peerclaw-agent
go build -o echo ./examples/echo
```

```bash
# Terminal 1: Start the gateway
./peerclawd
# → PeerClaw gateway started  http=:8080

# Terminal 2: Start agent Alice
./echo -name alice -server http://localhost:8080

# Terminal 3: Start agent Bob
./echo -name bob -server http://localhost:8080

# Terminal 4: See who's online (install CLI from this repo's cli/ directory)
go run ./cli/cmd/peerclaw agent list
```

Alice and Bob will automatically register, discover each other, and establish an encrypted P2P connection.

## Architecture

PeerClaw is composed of four modules that work together:

```
┌──────────────────────────────────────────────────────────────────┐
│                     peerclaw-server (Gateway)                    │
│                                                                  │
│  ┌────────────┐   ┌───────────────┐   ┌───────────────────────┐ │
│  │  Registry  │   │   Signaling   │   │    Bridge Manager     │ │
│  │  Agent     │   │   Hub         │   │                       │ │
│  │  discovery │   │   WebSocket   │   │  ┌─────┬─────┬─────┐ │ │
│  │  by caps   │   │   relay for   │   │  │ A2A │ MCP │ ACP │ │ │
│  └────────────┘   │   WebRTC      │   │  └─────┴─────┴─────┘ │ │
│                   └───────────────┘   └───────────────────────┘ │
│  ┌────────────┐   ┌───────────────┐   ┌───────────────────────┐ │
│  │  Auth      │   │   Rate Limit  │   │   Observability       │ │
│  │  Ed25519 + │   │   Per-IP      │   │   OpenTelemetry       │ │
│  │  API Key   │   │   throttling  │   │   traces + metrics    │ │
│  └────────────┘   └───────────────┘   └───────────────────────┘ │
│                                                                  │
│  Storage: SQLite (default) or PostgreSQL                         │
│  Scaling: Redis Pub/Sub for multi-node signaling                 │
└──────────────────────────────────────────────────────────────────┘
        │ REST API              │ WebSocket             │ Protocol
        │ register/discover     │ signaling             │ endpoints
        ▼                       ▼                       ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────┐
│  peerclaw-agent  │    │  peerclaw-agent  │    │  External    │
│  (Go SDK)        │◄══►│  (Go SDK)        │    │  A2A/MCP/ACP │
│                  │P2P │                  │    │  Agent       │
│  WebRTC primary  │    │  WebRTC primary  │    │              │
│  Nostr fallback  │    │  Nostr fallback  │    │              │
└──────────────────┘    └──────────────────┘    └──────────────┘
        │                       │
        └───── peerclaw-core ───┘
              (shared types: identity, envelope, protocol)
```

### How Agents Communicate

```
Alice                    Gateway                     Bob
  │                        │                          │
  ├─ POST /agents ────────►│  Register Alice          │
  │                        │◄──────── POST /agents ───┤  Register Bob
  │                        │                          │
  ├─ POST /discover ──────►│  "who can search?"       │
  │◄── [{Bob, caps:search}]│                          │
  │                        │                          │
  ├─ WS: offer + X25519 ─►│──── relay ──────────────►│
  │◄─── WS: answer + X25519│◄─── relay ──────────────┤
  │                        │                          │
  │◄═══════════ WebRTC P2P (encrypted) ══════════════►│
  │         Ed25519 signed + XChaCha20 encrypted      │
```

## Project Structure

| Module | What it does | Key tech |
|--------|-------------|----------|
| [**peerclaw-core**](https://github.com/peerclaw/peerclaw-core) | Shared type library — identity, envelope, agent card, protocol constants | Ed25519, X25519, zero external deps |
| [**peerclaw-server**](https://github.com/peerclaw/peerclaw-server) | The gateway — registration, discovery, signaling relay, protocol bridging | SQLite/PostgreSQL, WebSocket, OTel |
| [**peerclaw-agent**](https://github.com/peerclaw/peerclaw-agent) | P2P agent SDK — connect, send, receive with automatic transport selection | WebRTC (Pion), Nostr, TOFU trust |
| **cli/** | Command-line tool — manage agents, check health, send messages | Cobra-style subcommands |

## Core Concepts

### Agent Card

Every agent publishes an Agent Card — a machine-readable description of who it is and what it can do:

```json
{
  "name": "search-agent",
  "public_key": "base64-ed25519-pubkey",
  "capabilities": ["web-search", "summarize"],
  "protocols": ["a2a", "mcp"],
  "endpoint": { "url": "https://my-agent.example.com", "port": 443 },
  "skills": [{ "id": "search", "name": "Web Search" }],
  "tools": [{ "name": "search", "description": "Search the web" }]
}
```

Compatible with the A2A Agent Card standard, extended with PeerClaw fields (public key, NAT type, DHT node ID).

### Protocol Bridging

PeerClaw translates between protocols using a universal **Envelope** format:

```
A2A Agent ──► A2A Adapter ──► Envelope ──► MCP Adapter ──► MCP Agent
                                │
                           unified format:
                           source, destination,
                           protocol, payload,
                           signature, trace_id
```

| Protocol | What it's for | PeerClaw support |
|----------|--------------|------------------|
| **A2A** (Google) | Task-based agent collaboration | Full: tasks, artifacts, streaming |
| **MCP** (Anthropic) | Tool/resource access | Full: tools, resources, prompts |
| **ACP** (IBM) | Enterprise agent runs | Full: runs, sessions, manifests |

### Cryptographic Identity

Every agent owns an Ed25519 keypair. The public key **is** the identity.

- **Registration**: agent proves ownership by signing the request
- **Messages**: every envelope is signed — receiver verifies origin
- **Encryption**: X25519 keys derived from Ed25519, XChaCha20-Poly1305 for payloads
- **Trust**: TOFU (Trust-On-First-Use) model with 5 levels: Unknown → TOFU → Verified → Pinned → Blocked

### Transport Fallback

The agent SDK automatically picks the best transport:

```
1. WebRTC DataChannel (preferred — low latency, P2P)
       │ fails (strict NAT)?
       ▼
2. Nostr relay (fallback — NIP-44 encrypted, multi-relay)
       │ WebRTC recovers?
       ▼
3. Auto-upgrade back to WebRTC
```

## Advanced Features

These are available but not required for basic usage:

| Feature | Description |
|---------|-------------|
| **Reputation Engine** | Server-side EWMA reputation scoring from real interaction events |
| **Endpoint Verification** | Challenge-response proof that agents control their claimed URLs |
| **Public Directory** | Browse and search agents by reputation, capability, category, and verification status |
| **Agent Playground** | Try any agent live via chat UI with SSE streaming, rate-limited anonymous access |
| **User Auth & JWT** | Email/password registration, JWT sessions, API key management |
| **Provider Console** | Publish agents, view analytics, manage invocations and API keys |
| **Reviews & Ratings** | Star ratings (1-5) + text reviews with reputation integration |
| **Trusted Badge** | Verified + high reputation agents earn a "Trusted" badge |
| **DHT Discovery** | Serverless agent discovery via Kademlia DHT (Nostr transport) |
| **Federation** | Multi-server signaling relay with DNS SRV discovery |
| **Identity Anchoring** | Bind Ed25519 identity to Nostr/DNS for public verification |
| **Offline Messaging** | Message cache with TTL, auto-flush on peer reconnect |
| **Serverless Mode** | Full P2P operation without any central server |

## Agent Marketplace (Phase 7)

PeerClaw has evolved from infrastructure into a **C2C Agent Marketplace (Agent as a Service)**:

- **Browse & Discover** — Landing page, explore page, agent profiles with trust info, category filtering
- **Playground** — Try any Agent live through a protocol-agnostic chat interface with SSE streaming
- **User Accounts** — Register/login with JWT auth, publish agents via 5-step wizard, manage API keys
- **Provider Console** — Dashboard with call volume analytics, agent stats, invocation history
- **Trust & Community** — Star ratings, text reviews, Verified/Trusted badges, abuse reporting

See [Roadmap](docs/ROADMAP.md) for the complete development history.

## CLI Reference

```bash
peerclaw health                              # Check gateway status
peerclaw agent list                          # List all agents
peerclaw agent list -protocol mcp -output json   # Filter + JSON output
peerclaw agent get <id>                      # Agent details
peerclaw agent register -name "My Agent" ... # Register an agent
peerclaw send -from a -to b -payload '{}'    # Send a message
peerclaw config set server http://host:8080  # Set gateway URL
```

## Development

```bash
# Each module is a separate repository. Build and test individually:

# peerclaw-core
git clone https://github.com/peerclaw/peerclaw-core.git
cd peerclaw-core && go build ./... && go test ./...

# peerclaw-server (requires CGO for SQLite)
git clone https://github.com/peerclaw/peerclaw-server.git
cd peerclaw-server && CGO_ENABLED=1 go build ./... && CGO_ENABLED=1 go test ./...

# peerclaw-agent
git clone https://github.com/peerclaw/peerclaw-agent.git
cd peerclaw-agent && go build ./... && go test ./...
```

For local multi-module development, you can use a [Go workspace](https://go.dev/doc/tutorial/workspaces):

```bash
mkdir peerclaw && cd peerclaw
git clone https://github.com/peerclaw/peerclaw-core.git core
git clone https://github.com/peerclaw/peerclaw-server.git server
git clone https://github.com/peerclaw/peerclaw-agent.git agent
go work init ./core ./server ./agent
go work sync
```

## Documentation

- [User Guide](docs/GUIDE.md) — Browse, try, publish, and manage Agents on PeerClaw
- [Product Document](docs/PRODUCT.md) — Detailed product design and security model
- [Roadmap](docs/ROADMAP.md) — Development phases and milestones

## Contributing

PeerClaw is in active development. We welcome contributions:

- **Issues** — Bug reports, feature requests, questions
- **Pull Requests** — Code contributions to any module
- **Discussions** — Ideas about the future of agent communication

## License

MIT
