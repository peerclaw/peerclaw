**English** | [中文](ROADMAP_zh.md)

# PeerClaw Roadmap

## Phase 1: Foundation (Complete)

Lay the core infrastructure and validate end-to-end communication flows.

- [x] **peerclaw-core** — Shared type library
  - Ed25519 identity (key pair generation / loading / saving / signing / verification)
  - Envelope: unified message wrapper
  - Agent Card definition (A2A-compatible + PeerClaw extensions)
  - Protocol constants (A2A / ACP / MCP) and transport types
  - Signaling message types
- [x] **peerclaw-server** — Centralized platform
  - Agent registration and deregistration (REST API)
  - Discover agents by capability and protocol
  - Heartbeat and status management
  - WebSocket signaling Hub (offer / answer / ICE candidate)
  - Routing engine (capability matching, protocol routing)
  - Protocol bridging framework (A2A / ACP / MCP adapter scaffolding)
  - SQLite persistence
  - YAML configuration
- [x] **peerclaw-agent** — P2P Agent SDK
  - WebRTC DataChannel transport
  - Nostr relay transport (basic)
  - TOFU Trust Store
  - Message signing and verification
  - Peer Manager
  - Discovery Client
  - Signaling Client
  - Echo Agent example

## Phase 2: Transport & Security Hardening (Complete)

Harden the transport layer and security mechanisms to improve connection reliability.

- [x] **Full Nostr relay implementation**
  - NIP-44 encryption (based on the `fiatjaf.com/nostr` library)
  - Multi-relay support with failover (publish to all, deduplicate subscriptions)
  - Relay health checks (exponential backoff reconnection, removed from active set after 3 failures)
- [x] **NAT traversal optimization**
  - TURN server integration (server pushes ICE config after signaling connection)
  - ICE candidate filtering and priority ordering (host > srflx > relay)
  - Connection quality monitoring (RTT, packet loss, throughput metrics)
- [x] **Automatic transport selection**
  - WebRTC to Nostr automatic fallback (Transport Selector)
  - Automatic upgrade when connectivity recovers (background probing)
  - Transport health scoring (rolling window success/failure counters)
- [x] **Trust Store enhancements**
  - CLI management tool `peerclaw-trust` (list / verify / pin / revoke / export / import)
  - Trust levels (Unknown / TOFU / Verified / Blocked / Pinned)
  - Trust event notifications (OnTrustChange callback)
- [x] **End-to-end encryption**
  - X25519 key exchange (derived from Ed25519 seed, public keys exchanged during signaling)
  - XChaCha20-Poly1305 message encryption (HKDF-SHA256 key derivation)
  - Additional NIP-44 format wrapping for Nostr transport

## Phase 3: Protocol Ecosystem (Complete)

Fully implement protocol bridging across all three protocols, making PeerClaw the interoperability hub.

- [x] **JSON-RPC 2.0 shared library**
  - Request / Response / Error / Notification types
  - ParseMessage for automatic message type detection
  - Standard error codes (-32700 ~ -32600)
  - Shared by A2A and MCP
- [x] **Full A2A adapter**
  - Task lifecycle (create / working / complete / cancel / fail)
  - Artifact model (multi-Part support: text / file / structured data)
  - Agent Card standard compliance (GET /.well-known/agent.json)
  - JSON-RPC inbound Handler (message/send / tasks/get / tasks/cancel)
  - Outbound SendMessage to Task response to Envelope conversion
- [x] **Full MCP adapter (Streamable HTTP)**
  - Tool model (tools/list / tools/call)
  - Resource model (resources/list / resources/read)
  - Prompt model (prompts/list / prompts/get)
  - Session management (initialize handshake + Mcp-Session-Id)
  - Streamable HTTP transport (POST + GET SSE endpoints)
- [x] **Full ACP adapter**
  - Complete message types (MessagePart with content_type / content / content_url)
  - Session management (Session + Run lifecycle tracking)
  - Agent Manifest queries (GET /acp/agents/{name})
  - Run management (create / get / cancel)
- [x] **Automatic protocol negotiation**
  - Negotiator automatically selects the optimal protocol path
  - Same-protocol direct connection preferred over cross-protocol translation (priority: A2A > MCP > ACP)
  - Cross-protocol translation (A2A ↔ MCP ↔ ACP)
- [x] **Enhanced agent capability declarations**
  - Structured Skills declarations (A2A-compatible: name / description / input_modes / output_modes)
  - Structured Tools declarations (MCP-compatible: name / description / input_schema)
  - HasSkill() / HasTool() query methods
  - SQLite persistence (JSON-serialized storage)
- [x] **Server routing integration**
  - Protocol endpoint routing (POST /a2a, POST /mcp, GET/POST /acp/*)
  - Unified bridge send endpoint (POST /api/v1/bridge/send)
  - Bridge Forwarder (bridge inbox → signaling hub → agent)
  - bridge_message signaling message type

## Phase 4: Production Readiness (Complete)

Stability, observability, and operational capabilities for production environments.

- [x] **Observability**
  - OpenTelemetry traces (opt-in, OTLP gRPC export)
  - OpenTelemetry metrics (HTTP request rate/latency, WebSocket connections, agent registrations, bridge message throughput)
  - Enhanced structured logging (middleware chain: Recovery → RequestID → Tracing → Logging → RateLimit → MaxBody)
  - Grafana dashboard template (docs/grafana/peerclaw-overview.json)
- [x] **Horizontal scaling**
  - Redis Pub/Sub for cross-node signaling (Broker interface + LocalBroker + RedisBroker)
  - PostgreSQL backend (JSONB + GIN indexes, driver selected via factory.go based on config)
  - Enhanced health checks (component-level status: database, signaling)
- [x] **Audit logging**
  - Agent registration/deregistration records
  - Message routing audit (Bridge send)
  - Security event logging (rate limiting, signaling connect/disconnect)
  - Dedicated slog.Logger instance (stdout or file output)
- [x] **Rate limiting and protection**
  - Per-IP token bucket request rate limiting (golang.org/x/time/rate)
  - WebSocket connection limits (Hub.maxConns)
  - Message size limits (http.MaxBytesReader middleware)
- [x] **CLI tooling (peerclaw-cli)**
  - `peerclaw agent list|get|register|delete` — Agent management
  - `peerclaw send` — Send messages via Bridge
  - `peerclaw health` — Server health check
  - `peerclaw config show|set` — CLI configuration management

## Phase 5: Decentralized Evolution (Complete)

Evolve toward full decentralization, enabling serverless agent communication.

- [x] **Interface abstractions + Agent Card extensions**
  - Discovery interface (RegistryClient / DHTDiscovery / CompositeDiscovery)
  - SignalingClient interface (WebSocket / NostrSignaling / CompositeSignaling)
  - Agent struct refactored to use interfaces (backward-compatible)
  - PeerClawExtension new fields (nostr_pubkey / dht_node_id / reputation_score / nostr_relays / identity_anchor)
  - New signaling message types (federation_forward / dht_ping/store/find_node/find_value)
- [x] **DHT-based decentralized discovery**
  - Minimal Kademlia DHT (160-bit NodeID, K=20 k-bucket routing table)
  - DHT protocol messages (Ping / Store / FindNode / FindValue RPC)
  - DHT transport layer (Nostr event kind 20005 + InMemory test transport)
  - DHT coordinator (Bootstrap / Put / Get / FindNode / bucket refresh / data republish)
  - DHTDiscovery implementing Discovery interface (primary key + capability index)
  - CompositeDiscovery (Server-first + DHT fallback)
- [x] **Multi-node signaling federation**
  - FederationConfig (static configuration + DNS SRV discovery)
  - FederationService (connect to peers, ForwardSignal, QueryAgents)
  - DNS SRV discovery (_peerclaw._tcp.<domain>)
  - FederationBroker (local.Publish if target is local, otherwise federation.ForwardSignal)
  - Hub.HasAgent() to check local connections
  - Federation HTTP endpoints (/api/v1/federation/signal, /api/v1/federation/discover)
  - DiscoverFederated method
- [x] **Agent reputation system**
  - ReputationStore: per-peer EWMA scoring (0.0-1.0), 6 behavior types
  - RecordEvent / GetScore / IsMalicious (< 0.15 threshold)
  - TrustStore integration (SetReputationStore / IsAllowedWithReputation)
  - Reputation gossip (Nostr kind 30078, second-hand reputation weighted at 0.3x, only accepts peers at TrustVerified+)
  - JSON file persistence
- [x] **Serverless pure P2P mode**
  - NostrSignaling (event kind 20006, NIP-44 encryption)
  - CompositeSignaling (WebSocket-first + Nostr fallback)
  - MessageCache for offline message buffering (per-destination queues, TTL expiration, JSON persistence)
  - OnPeerAdded callback (flushes cached queue when a new peer connects)
  - Serverless mode Options (DHTEnabled / Serverless / ICEServers / MessageCachePath)
- [x] **On-chain identity anchoring (optional)**
  - IdentityAnchor interface (Publish / Verify / Resolve / RecoveryKeys)
  - NostrAnchor implementation (Nostr kind 10078 replaceable event, bidirectional key binding)
  - Domain binding verification (DNS TXT record peerclaw-verify=<fingerprint>)
  - Multi-sig recovery (threshold-of-n recovery keys)
- [x] **CLI Phase 5 commands**
  - `peerclaw dht bootstrap|lookup`
  - `peerclaw federation status|peers`
  - `peerclaw reputation show|list`
  - `peerclaw identity anchor|verify`
- [x] **Integration tests**
  - DHT-only discovery and communication (no server)
  - Server + DHT hybrid mode with fallback
  - Malicious behavior triggering reputation isolation
  - Reputation gossip propagation across peers
  - DHT Agent Card storage and retrieval
  - Offline message cache delivery

## Phase 6: Agent Identity & Trust Platform (Complete)

Transform PeerClaw from a protocol gateway into an identity & trust platform. The gateway remains as infrastructure — real interactions generate the trust data that differentiates PeerClaw.

- [x] **Server-side EWMA reputation engine**
  - Ported EWMA algorithm from agent SDK to server (`internal/reputation/`)
  - 10 event types with configurable weights (registration, heartbeat, verification, bridge, review)
  - `reputation_events` table for full event history
  - Reputation columns on `agents` table (score, event_count, updated_at, verified, verified_at)
  - Auto-integrated into registration, heartbeat, and bridge handlers
  - Background heartbeat timeout checker (60s interval, 5m timeout → heartbeat_miss event)
- [x] **Endpoint verification**
  - Challenge-response flow (`internal/verification/`)
  - Server generates random nonce, sends to agent's `/.well-known/peerclaw-verify` endpoint
  - Agent responds with nonce + Ed25519 signature
  - SSRF protection via `security/urlvalidator.go`, 5s HTTP timeout, no redirects
  - `verification_challenges` table with 5-minute TTL
- [x] **Public API layer (no auth required)**
  - `GET /api/v1/directory` — Agent directory with search/filter/sort (reputation, name, registered_at)
  - `GET /api/v1/directory/{id}` — Sanitized public profile (no auth params, conditional endpoint URL)
  - `GET /api/v1/directory/{id}/reputation` — Reputation event history
  - `POST /api/v1/agents/{id}/verify` — Initiate endpoint verification (owner only)
  - `PublicEndpoint` opt-in field on `PeerClawExtension` for endpoint URL visibility
- [x] **Frontend restructure**
  - Renamed `web/dashboard` to `web/app`
  - Public pages: Landing Page, Agent Directory, Public Profile (with reputation chart)
  - Admin dashboard moved to `/admin/*` routes
  - New components: PublicLayout, AgentDirectoryCard, ReputationMeter, VerifiedBadge, ReputationChart
  - Recharts-based reputation history visualization

## Phase 7: Agent Marketplace (Complete)

Evolve PeerClaw into a C2C Agent Marketplace — where anyone can publish an Agent as a service, and anyone (human or Agent) can discover and invoke it.

### Phase 7a: Marketplace Browse & Profile

Public-facing marketplace for discovering and evaluating Agents.

- [x] **Landing page** — Platform stats, value propositions, search entry point (shipped in Phase 6)
- [x] **Explore page** — Agent Directory with search, filter (verified, min_score), and sort (reputation, name, registered_at) (shipped in Phase 6)
- [x] **Agent Profile page** — Detailed view with capabilities, protocols, trust info, reputation history chart (shipped in Phase 6)
- [x] **Top navigation bar** — PublicLayout with horizontal nav, distinct from admin sidebar (shipped in Phase 6)
- [x] **Mobile-responsive design** — Responsive card-based layout with AgentDirectoryCard (shipped in Phase 6)
- [x] **Extended query API** — `category` filter on directory endpoint, full-text `search` param

### Phase 7b: Playground & Invocation

Let consumers try and invoke Agents through a protocol-agnostic interface.

- [x] **Protocol-agnostic invocation endpoint** — `POST /api/v1/invoke/{agent_id}`, auto-selects optimal protocol (A2A / MCP / ACP) via Bridge Manager
- [x] **Chat-style Playground** — Web UI for trying Agents live, with developer mode toggle for raw request/response inspection
- [x] **SSE streaming** — `stream: true` in invoke request, `Content-Type: text/event-stream` response with `http.Flusher`
- [x] **Anonymous rate-limited trial** — Playground access without login, rate-limited per IP (10 calls/hour, burst 3)
- [x] **Invocation record logging** — `internal/invocation/` module records every call (agent, caller, protocol, latency, status, error) with SQLite/PostgreSQL persistence

### Phase 7c: User Accounts & Provider Console

User identity and Agent management for providers.

- [x] **User registration & login** — Email/password with bcrypt hashing (`internal/userauth/`)
- [x] **JWT session management** — Access token (15m) + refresh token (168h) with automatic rotation, `internal/userauth/jwt.go`
- [x] **Agent publish wizard** — Guided 5-step registration (basic info → capabilities & protocols → endpoint → auth & metadata → preview)
- [x] **Provider Dashboard** — My agents overview with total calls, success rate, average latency
- [x] **API key management** — Generate, list, and revoke API keys with SHA-256 hashing, prefix display
- [x] **Interaction history** — Consumer and provider views of past invocations with filtering

### Phase 7d: Trust & Community

Community-driven trust signals and provider analytics.

- [x] **Reviews & ratings** — Star ratings (1-5) + text reviews with UNIQUE(agent_id, user_id) constraint, reputation integration (rating ≥ 4 → review_positive, ≤ 2 → review_negative)
- [x] **Verified / Trusted badges** — "Verified" via endpoint verification, "Trusted" badge for verified + reputation > 0.8
- [x] **Categories & tagging** — Structured categorization with `categories` + `agent_categories` tables, category filter on directory
- [x] **Provider analytics dashboard** — Call volume time series, agent stats (total/success/error calls, avg/p95 duration)
- [x] **Abuse reporting** — Report system for agents and reviews with reason + details, status tracking (pending/reviewed/dismissed/actioned)
