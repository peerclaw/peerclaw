**English** | [中文](PRODUCT_zh.md)

# PeerClaw Product Document

## Product Vision

Build the trust infrastructure for AI Agents and evolve it into an open platform where any Agent becomes a discoverable, trustable, invocable service.

**Layer 1 — Infrastructure (Complete):** Decentralization-first communication, verifiable identity, reputation scoring, and cross-protocol bridging (A2A / MCP / ACP). Any Agent, regardless of protocol or deployment location, can discover and communicate with others securely and efficiently.

**Layer 2 — Platform (Complete):** An Agent Platform built on top of this infrastructure, where anyone can register an Agent as a service and anyone (human or Agent) can discover, evaluate, and invoke it — regardless of protocol.

## Target Users

### AI Agent Developers

- Need their Agent to communicate with other Agents
- Don't want to be locked into a single protocol ecosystem (A2A / MCP / ACP)
- Prefer P2P direct connections over routing all traffic through a central server
- Need out-of-the-box security (identity, signatures, trust management)

### Platform Integrators

- Operate Agent platforms and need a unified registration and discovery mechanism
- Need to bridge Agents across different protocol ecosystems
- Need observability and auditing capabilities
- Need horizontal scaling and high availability

### Agent Service Consumers

- Need to discover and use Agent services without knowing the underlying protocol
- Want to evaluate Agent trustworthiness before committing (reputation, reviews, verified identity)
- Need a simple way to try Agents (Playground) and invoke them programmatically
- Could be humans browsing the platform or Agents delegating subtasks to specialist Agents

## Core Scenarios

### Scenario 1: Agent Discovery and Connection

```
Alice (Search Agent)                  PeerClaw Server                    Bob (Data Agent)
       │                                    │                                   │
       │  POST /api/v1/agents (Register)    │                                   │
       │──────────────────────────────────►  │                                   │
       │                                    │  ◄── POST /api/v1/agents (Register)│
       │                                    │                                   │
       │  POST /api/v1/discover             │                                   │
       │  {"capabilities": ["data"]}        │                                   │
       │──────────────────────────────────►  │                                   │
       │  ◄── [{id: "bob", pubkey: "..."}]  │                                   │
       │                                    │                                   │
```

Agents register their capabilities and public keys via the REST API. Other Agents can search and discover them by capability or protocol.

### Scenario 2: Secure P2P Communication

```
Alice                           Signaling Hub                          Bob
  │                                  │                                  │
  │  WS: {type:"offer", to:"bob",   │                                  │
  │       sdp:"..."}                 │                                  │
  │─────────────────────────────────►│─────────────────────────────────►│
  │                                  │  WS: {type:"answer", to:"alice", │
  │                                  │       sdp:"..."}                 │
  │◄─────────────────────────────────│◄─────────────────────────────────│
  │                                  │                                  │
  │  ◄═══════ WebRTC DataChannel (P2P Direct) ═══════►                 │
  │                                  │                                  │
  │  Envelope {encrypted_payload, signature}                           │
  │════════════════════════════════════════════════════════════════════►│
  │  ◄═════════════ Envelope {encrypted_payload, signature} ══════════│
```

The signaling server is only used for the WebRTC handshake. Actual data flows through the P2P DataChannel. Every message carries an Ed25519 signature.

### Scenario 3: Cross-Protocol Interoperability

```
A2A Agent                    PeerClaw Server                    MCP Agent
    │                              │                                │
    │  A2A Task Request            │                                │
    │─────────────────────────────►│                                │
    │                    ┌─────────┴─────────┐                      │
    │                    │ A2A Adapter       │                      │
    │                    │ → PeerClaw Envelope│                      │
    │                    │ → MCP Adapter      │                      │
    │                    └─────────┬─────────┘                      │
    │                              │  MCP Tool Call                  │
    │                              │───────────────────────────────►│
    │                              │  ◄── MCP Tool Result           │
    │                              │───────────────────────────────►│
    │  ◄── A2A Task Result         │                                │
    │◄─────────────────────────────│                                │
```

The Bridge Manager automatically identifies the source and target protocols and performs seamless translation through the Envelope intermediate format.

### Scenario 4: Decentralized Fallback

```
Alice                         Nostr Relay                          Bob
  │                               │                                 │
  │  (WebRTC connection failed)   │                                 │
  │                               │                                 │
  │  NIP-44 Encrypted Envelope    │                                 │
  │──────────────────────────────►│────────────────────────────────►│
  │                               │  ◄── NIP-44 Encrypted Envelope  │
  │◄──────────────────────────────│◄────────────────────────────────│
```

When a WebRTC P2P connection cannot be established (strict NAT, firewalls), the system automatically falls back to transport via Nostr relay.

### Scenario 5: Agent Platform

```
Provider                       PeerClaw Platform                       Consumer
    │                                   │                                  │
    │  Register Agent                   │                                  │
    │  (capabilities, skills, endpoint) │                                  │
    │──────────────────────────────────►│                                  │
    │                                   │                                  │
    │                                   │  ◄── Browse / Search Agents      │
    │                                   │  ◄── Evaluate Trust & Reputation │
    │                                   │──────────────────────────────────│
    │                                   │      Agent Profile + Reviews     │
    │                                   │                                  │
    │                                   │  ◄── Try in Playground           │
    │                                   │──────────────────────────────────│
    │  ◄── Protocol-Agnostic Invocation │      (rate-limited trial)       │
    │◄──────────────────────────────────│                                  │
    │  Response ──────────────────────►│──────────────────────────────────►│
    │                                   │                                  │
```

The Provider registers their Agent once. Consumers discover it through the platform, evaluate trust through PeerClaw's built-in reputation and verification system, try it in the Playground, and invoke it — all without needing to know which protocol the Agent uses. PeerClaw automatically selects the optimal protocol path (A2A, MCP, or ACP).

## System Architecture

```
┌────────────────────────────────────────────────────────────────┐
│                     Platform Layer (Complete)                     │
│                                                                │
│  Landing / Explore / Agent Profile / Playground                │
│  User Accounts / Reviews & Ratings / Provider Analytics        │
└────────────────────────┬───────────────────────────────────────┘
                         │ Built on
┌────────────────────────▼───────────────────────────────────────┐
│                 Infrastructure Layer (Complete)                 │
│  Registry / Signaling / Bridge / Reputation / Identity / DHT   │
└────────────────────────────────────────────────────────────────┘
```

### Infrastructure Detail

```
┌────────────────────────────────────────────────────────────────┐
│                      peerclaw-server                           │
│                                                                │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │       HTTP Entry + Protocol Endpoints + Federation      │   │
│  │  /api/v1/*  /a2a  /mcp  /acp/*  /api/v1/federation/*   │   │
│  └──────┬──────────────┬──────────────┬───────────────────┘   │
│         │              │              │                        │
│  ┌──────▼──────┐ ┌─────▼─────┐ ┌─────▼──────┐                │
│  │  Registry   │ │  Router   │ │  Signaling │                │
│  │  Service    │ │  Engine   │ │  Hub       │                │
│  │             │ │           │ │            │                │
│  │ - Register/ │ │ - Routing │ │ - WS Conn  │                │
│  │   Deregister│ │   Table   │ │ - Message  │                │
│  │ - Discover/ │ │ - Route   │ │   Relay    │                │
│  │   Query     │ │   Resolve │ │ - Ping/Pong│                │
│  │ - Heartbeat │ │ - Capabil.│ │- bridge_msg│                │
│  │ - Federated │ │   Match   │ │            │                │
│  │   Discovery │ │           │ │            │                │
│  └──────┬──────┘ └───────────┘ └─────┬──────┘                │
│         │                            │                        │
│  ┌──────▼──────┐ ┌───────────────────▼────────────────────┐   │
│  │   SQLite/   │ │         Bridge Manager                 │   │
│  │  PostgreSQL │ │  ┌──────┐  ┌──────┐  ┌──────┐         │   │
│  └─────────────┘ │  │ A2A  │  │ ACP  │  │ MCP  │         │   │
│                  │  │Adapter│  │Adapter│  │Adapter│         │   │
│  ┌─────────────┐ │  └──────┘  └──────┘  └──────┘         │   │
│  │ Federation  │ │  ┌────────────┐  ┌──────────────────┐  │   │
│  │  Service    │ │  │ Negotiator │  │ Bridge Forwarder │  │   │
│  │ - Peer Conn │ │  └────────────┘  └──────────────────┘  │   │
│  │ - Signaling │ └────────────────────────────────────────┘   │
│  │   Relay     │                                              │
│  │ - DNS SRV   │                                              │
│  └─────────────┘                                              │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│                      peerclaw-agent (SDK)                      │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐ │
│  │   Agent API  │  │  Discovery   │  │     Signaling        │ │
│  │              │  │  (Interface) │  │     (Interface)      │ │
│  │              │  │ Registry     │  │  WebSocket Client    │ │
│  │              │  │ DHTDiscovery │  │  NostrSignaling      │ │
│  │              │  │ Composite    │  │  Composite           │ │
│  └──────────────┘  └──────────────┘  └──────────────────────┘ │
│  ┌──────────────┐  ┌──────────────────────────────────────┐   │
│  │ Peer Manager │  │            Security                  │   │
│  │ - OnPeerAdded│  │  Trust Store + Message Validator +    │   │
│  │              │  │  Reputation + Sandbox                 │   │
│  └──────────────┘  └──────────────────────────────────────┘   │
│  ┌──────────────┐  ┌──────────────────────────────────────┐   │
│  │     DHT      │  │            Identity                  │   │
│  │  (Kademlia)  │  │  IdentityAnchor + NostrAnchor +      │   │
│  │ - Routing    │  │  Domain Verify + Recovery             │   │
│  │   Table      │  └──────────────────────────────────────┘   │
│  │ - KV Store   │                                              │
│  │ - Bootstrap  │                                              │
│  └──────────────┘                                              │
│  ┌────────────────────────────────────────────────────────┐   │
│  │                    Transport                           │   │
│  │  WebRTC DataChannel ◄──► Nostr Relay ◄──► MessageCache │   │
│  └────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│                      peerclaw-core (Type Library)              │
│                                                                │
│  identity    envelope    agentcard    protocol    signaling    │
│  (Ed25519)   (Message    (Agent       (Protocol   (Signaling   │
│               Envelope)   Card)        Constants)  Types)      │
└────────────────────────────────────────────────────────────────┘
```

## Communication Flow

The full Agent-to-Agent communication goes through the following stages:

### 1. Registration

When an Agent starts up, it registers with the Server via the REST API, submitting an Agent Card (name, public key, capabilities, protocols). The Server stores this in SQLite and updates the routing table.

### 2. Discovery

Agent A searches for target Agents by capability via `/api/v1/discover`. The Server returns a list of matching Agent Cards containing public keys and connection information.

### 3. Signaling

Agent A connects to the Server's Signaling Hub via WebSocket and sends a WebRTC offer SDP to Agent B. The Server relays signaling messages. Agent B replies with an answer SDP, and both sides exchange ICE candidates.

### 4. P2P Connection

After WebRTC ICE negotiation completes, Agent A and B establish a direct DataChannel connection. If ICE fails (strict NAT), the system falls back to Nostr relay.

### 5. Message Exchange

Messages are wrapped in an Envelope containing source/destination, protocol type, payload, and an Ed25519 signature. The receiver verifies the signature before processing the message.

## Security Model

### Layer 1: Connection-Level -- TOFU (Trust-On-First-Use)

| Property | Description |
|----------|-------------|
| Mechanism | Records the peer's public key fingerprint on first connection |
| Storage | Local Trust Store file |
| Verification | Subsequent connections check whether the public key matches the record |
| Threat Mitigation | Man-in-the-middle attacks (after first use) |
| Analogy | SSH known_hosts |

### Layer 2: Message-Level -- Ed25519 Signatures

| Property | Description |
|----------|-------------|
| Algorithm | Ed25519 (RFC 8032) |
| Signed Object | Full Envelope (headers + payload; ciphertext when encrypted) |
| Verifier | Receiving Agent |
| Threat Mitigation | Message tampering, identity spoofing |
| Performance | ~76,000 signatures/sec, ~200,000 verifications/sec |

### Layer 3: Transport-Level -- End-to-End Encryption

| Property | Description |
|----------|-------------|
| Key Exchange | X25519 ECDH (derived from Ed25519 seed) |
| Key Derivation | HKDF-SHA256 |
| Symmetric Encryption | XChaCha20-Poly1305 (24-byte random nonce) |
| Ordering | Encrypt-then-sign — signature covers ciphertext, enabling pre-authentication before decryption |
| Session Establishment | X25519 public keys exchanged during signaling handshake |
| Nostr Adaptation | NIP-44 format wrapping (secp256k1 session key) |
| Threat Mitigation | Eavesdropping, man-in-the-middle, message leakage, decryption-oracle attacks |

### Layer 4: Execution-Level -- Sandbox

| Property | Description |
|----------|-------------|
| Mechanism | Enforces permission constraints on external requests |
| Controls | Resource limits, operation allowlists |
| Threat Mitigation | Malicious operations, resource exhaustion |

### Layer 5: P2P Communication -- Whitelist + Message Validation (Phase 8)

| Property | Description |
|----------|-------------|
| Default Policy | Default-deny — Agents must be whitelisted before communication |
| Agent-Side Whitelist | TrustStore-based check on inbound/outbound messages and connections |
| Server-Side Whitelist | ContactsChecker on signaling Hub blocks offer/answer/ICE for non-contacts |
| Connection Gating | ConnectionGate callback rejects offers before any WebRTC resource allocation |
| Message Validation | Signature verification, timestamp freshness (±2min), nonce-based replay protection, payload size limit (1MB) |
| Anti-Replay | UUID nonce per message, server-side nonce cache with 5-minute cleanup |
| Threat Mitigation | Prompt injection, replay attacks, resource exhaustion from unauthorized connections, DDoS via signaling flood |
| Architecture | Defense-in-depth: Agent TrustStore (primary) + Server contacts service (secondary) |

### Agent Access Control

PeerClaw enforces three tiers of access for agent invocation:

1. **Playground Access** — Agents with `playground_enabled=true` can be invoked by any authenticated user from the playground. This is opt-in per agent.
2. **Private Agents** — Agents with `visibility=private` are hidden from the public directory and can only be invoked by whitelisted contacts or approved users.
3. **User ACL** — For agents that are neither playground-enabled nor have the caller in contacts, users can submit access requests. Providers approve/reject with optional expiry dates.

This model gives providers full control over who can invoke their agents, from fully open to fully gated.

## Protocol Compatibility Matrix

| Feature | A2A | ACP | MCP |
|---------|-----|-----|-----|
| Message Routing | ✅ | ✅ | ✅ |
| Capability Discovery | ✅ | ✅ | ✅ |
| Task Model | ✅ | — | — |
| Tool Invocation | — | — | ✅ |
| Artifact | ✅ | — | — |
| Resource | — | — | ✅ |
| Prompt | — | — | ✅ |
| Run Management | — | ✅ | — |
| Protocol Negotiation | ✅ | ✅ | ✅ |
| Cross-Protocol Translation | ✅ | ✅ | ✅ |
| Bidirectional Streaming | ✅ | ✅ | 🔲 Phase 4 |

## Deployment Architecture

### Single Node

```
[ Agent A ] ──► [ PeerClaw Server (SQLite) ] ◄── [ Agent B ]
                         │
                   [ Agent A ] ◄═══ P2P ═══► [ Agent B ]
```

### Multi-Node (Phase 4 -- Implemented)

```
                        ┌── OTel Collector ── Grafana
                        │
[ CLI ] ──► [ Server 1 ] ◄── Redis Pub/Sub ──► [ Server 2 ] ◄── [ Agent ]
   │            │   │                               │   │
   │     Rate Limiter │                        Rate Limiter │
   │          Audit Log                           Audit Log
   │            │                                   │
[ Agent ] [ PostgreSQL / SQLite ]           [ PostgreSQL / SQLite ]
```

Middleware chain (per request):
`Recovery → RequestID → Tracing → Logging → Metrics → RateLimit → MaxBody`

### Decentralized (Phase 5 -- Implemented)

```
[ Agent A ] ◄═══ P2P (WebRTC) ═══► [ Agent B ]
     │                                   │
     └──── DHT Discovery (Kademlia) ─────┘
     │                                   │
     └──── Nostr Signaling (kind 20006) ─┘
     │                                   │
     └──── Nostr Relay (kind 20004) ─────┘
```

### Federated Mode (Phase 5 -- Implemented)

```
[ Agent A ] ──► [ Server 1 ] ◄══ Federation ══► [ Server 2 ] ◄── [ Agent B ]
                     │         (HTTP + Auth)          │
              [ DNS SRV Discovery ]            [ DNS SRV Discovery ]
```

### Serverless Mode (Phase 5 -- Implemented)

```
[ Agent A ]                                    [ Agent B ]
     │  1. DHT Bootstrap (Nostr kind 20005)         │
     │──────────────►  Nostr Relays  ◄──────────────│
     │  2. Agent Card stored in DHT                  │
     │  3. Nostr Signaling (kind 20006)              │
     │──────────────►  Nostr Relays  ◄──────────────│
     │  4. WebRTC P2P Direct Connection              │
     │◄════════════ DataChannel ═══════════════════►│
     │  5. Offline Message Cache (MessageCache)      │
```

## Reputation Model (Phase 5)

| Property | Description |
|----------|-------------|
| Scoring Algorithm | Exponentially Weighted Moving Average (EWMA, alpha=0.1) |
| Score Range | 0.0 (malicious) ~ 1.0 (trusted) |
| Malicious Threshold | < 0.15 triggers automatic isolation |
| Behavior Types | success (+1.0), timeout (-0.3), error (-0.2), invalid_signature (-0.8), spam (-0.5), protocol_violation (-0.7) |
| Storage | Persisted in local JSON file |
| Gossip | Optional Nostr kind 30078, second-hand weight 0.3x, accepts only TrustVerified+ |

## Identity Anchoring (Phase 5)

| Property | Description |
|----------|-------------|
| Interface | IdentityAnchor (Publish/Verify/Resolve/RecoveryKeys) |
| Primary Implementation | Nostr kind 10078 replaceable event |
| Key Binding | Bidirectional: Ed25519 signs Nostr key + Nostr key signs Ed25519 key |
| Domain Verification | DNS TXT record peerclaw-verify=<fingerprint> |
| Identity Recovery | threshold-of-n multisig recovery keys |

## Product Evolution

PeerClaw has evolved through three strategic stages:

```
Phase 1-4                    Phase 5-6                       Phase 7+
Communication                Identity & Trust                Agent Platform
Infrastructure               Platform

┌──────────────┐            ┌──────────────┐               ┌──────────────┐
│ Registry     │            │ Reputation   │               │ Browse &     │
│ Signaling    │──────────►│ Verification │──────────────►│ Discover     │
│ Bridging     │            │ Public       │               │ Playground   │
│ P2P / DHT   │            │ Directory    │               │ User Accounts│
│ Federation   │            │ Identity     │               │ Reviews      │
│              │            │ Anchoring    │               │ Analytics    │
└──────────────┘            └──────────────┘               └──────────────┘
```

- **Infrastructure** (Phase 1-4): Core protocol gateway — registry, signaling, bridging, transport, production readiness.
- **Identity & Trust Platform** (Phase 5-6): Decentralized identity, reputation scoring, endpoint verification, public directory. The real interactions generate trust data that differentiates PeerClaw.
- **Agent Platform** (Phase 7+): An open platform where anyone can register an Agent as a service and anyone can discover, evaluate, try, and invoke it — regardless of protocol.
- **P2P Security Hardening** (Phase 8): Default-deny whitelist enforcement, message validation pipeline (signature, replay, timestamp), and connection gating — defense-in-depth at both Agent and Server layers. See [Roadmap](ROADMAP.md) for details.
