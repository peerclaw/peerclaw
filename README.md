**English** | [中文](README_zh.md)

# PeerClaw

**Let AI Agents communicate as freely as humans do.**

PeerClaw is a decentralization-first communication framework for AI Agents. Agents identify each other through cryptographic identities, communicate via direct WebRTC connections, fall back to Nostr relays when NAT traversal fails, and achieve A2A / ACP / MCP interoperability through protocol bridging.

## Vision

The current AI Agent ecosystem faces severe communication fragmentation:

- **Protocol Silos** — A2A, ACP, and MCP operate in isolation; agents cannot communicate across protocols
- **Centralization Dependency** — Agent communication must be relayed through platform servers, adding latency and single points of failure
- **Missing Identity** — Agents lack unified cryptographic identities, making message source verification impossible
- **Weak Security** — Most agent communication solutions lack end-to-end security guarantees

PeerClaw's answer:

- **Decentralization First** — WebRTC P2P direct connections with Nostr relay fallback, no single-service dependency
- **Protocol Bridging** — Built-in A2A / ACP / MCP adapters, unified conversion to PeerClaw Envelope
- **Cryptographic Identity** — Every agent owns an Ed25519 keypair; the public key is the identity
- **Four-Layer Security** — Connection-level TOFU + message-level signing + end-to-end encryption (XChaCha20-Poly1305) + execution-level sandbox

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        peerclaw-server                      │
│                                                             │
│   ┌────────────┐  ┌──────────────┐  ┌───────────────────┐  │
│   │  Registry  │  │   Signaling  │  │  Bridge Manager   │  │
│   │ (Discovery)│  │  Hub (Relay) │  │  A2A / ACP / MCP  │  │
│   └────────────┘  └──────────────┘  └───────────────────┘  │
│         │                │                    │             │
└─────────┼────────────────┼────────────────────┼─────────────┘
          │                │                    │
    ┌─────┴─────┐    ┌─────┴─────┐        ┌────┴────┐
    │  Agent A  │◄──►│  Agent B  │ Extern  │ A2A/MCP │
    │ (SDK)     │P2P │ (SDK)     │ Agent   │ Agent   │
    └───────────┘    └───────────┘        └─────────┘
         │                │
    WebRTC DataChannel / Nostr relay
```

**Communication Flow:** Register → Discover → Signaling Handshake (with X25519 key exchange) → P2P Connection (WebRTC preferred, auto-fallback to Nostr) → Encrypted & Signed Message Exchange

## Sub-Projects

| Repository | Description | Status |
|------------|-------------|--------|
| [peerclaw-core](https://github.com/peerclaw/peerclaw-core) | Core shared type library (identity, envelope, protocol constants) | Active |
| [peerclaw-server](https://github.com/peerclaw/peerclaw-server) | Centralized platform (registration/discovery/signaling/bridging) | Active |
| [peerclaw-agent](https://github.com/peerclaw/peerclaw-agent) | P2P Agent SDK (WebRTC + Nostr + security) | Active |

## Quick Start

Get a P2P communication demo running in 5 minutes:

### 1. Clone and Build

```bash
git clone https://github.com/peerclaw/peerclaw.git
cd peerclaw

# Clone sub-projects
git clone https://github.com/peerclaw/peerclaw-core.git core
git clone https://github.com/peerclaw/peerclaw-server.git server
git clone https://github.com/peerclaw/peerclaw-agent.git agent

# Build
cd server && go build -o peerclawd ./cmd/peerclawd && cd ..
cd agent && go build -o echo ./examples/echo && cd ..
```

### 2. Start the Server

```bash
./server/peerclawd
# Output: PeerClaw gateway started  http=:8080  grpc=:9090
```

### 3. Start Two Echo Agents

```bash
# Terminal 1
./agent/echo -name alice -server http://localhost:8080

# Terminal 2
./agent/echo -name bob -server http://localhost:8080
```

Both agents will automatically register with the server and establish a WebRTC P2P connection via signaling.

## Local Development

```bash
# The project uses Go workspace to manage multiple modules
# Ensure the three sub-repos are in the correct locations: core/ server/ agent/

# Sync workspace
go work sync

# Build all modules
cd core && go build ./... && cd ..
cd server && go build ./... && cd ..
cd agent && go build ./... && cd ..

# Run tests
cd server && CGO_ENABLED=1 go test ./... && cd ..
cd agent && go test ./... && cd ..
```

## Documentation

- [Product Document](docs/PRODUCT.md) — Detailed product design, architecture, and security model
- [Roadmap](docs/ROADMAP.md) — Five-phase plan from foundation to decentralized evolution

## Community & Contributing

PeerClaw is in its early stages. We welcome your participation:

- Submit Issues to report bugs or suggest features
- Submit Pull Requests to contribute code
- Join the discussion on the future of agent communication

## License

MIT
