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

## Post-Phase 7: P2P Connection Orchestrator (Complete)

Close the gap between signaling infrastructure and automatic P2P connectivity.

- [x] **Connection Manager** (`agent/conn/`)
  - Signaling inbox consumption loop (offer / answer / ICE candidate dispatch)
  - Offerer flow: create WebRTCTransport → exchange SDP + ICE via signaling → block until DataChannel opens
  - Answerer flow: respond with SDP answer, tie-breaking by agent ID
  - ICE connection state monitoring → auto-register peer on Connected/Completed
  - Receive loop: read envelopes from DataChannel, dispatch to message handler
  - X25519 public key exchange in offer/answer for E2E encrypted sessions
- [x] **Agent Send() P2P-first with relay fallback**
  - Priority 1: existing P2P connection via PeerManager
  - Priority 2: establish new P2P connection via ConnManager (15s timeout)
  - Priority 3: signaling relay via bridge_message (WebSocket server relay)
- [x] **Signaling reconnection**
  - Auto-reconnect on unexpected WebSocket disconnection
  - Exponential backoff (1s → 60s)
- [x] **SignalingClient.SetAgentID()** — Deferred agent ID binding (set after registration, before Connect)

## Phase 8: P2P Communication Security Hardening (Complete)

Default-deny security model for Agent P2P communication — Agents must be whitelisted before they can connect or exchange messages.

- [x] **Message validation pipeline**
  - MessageValidator integrated into HandleIncomingEnvelope (signature verification, replay protection, payload size check)
  - Send() auto-populates Nonce (UUID), Timestamp, and Source before signing
  - Background nonce cleanup goroutine (5-minute interval)
- [x] **Whitelist enforcement (default-deny)**
  - Agent-side: TrustStore-based whitelist check on both inbound and outbound messages
  - Agent-side: Outbound Send() rejects messages to non-whitelisted destinations
  - Server-side: ContactsChecker interface on signaling Hub — blocks offer/answer/ICE for non-contacts
  - Contacts service wired into signaling hub at server startup
- [x] **Connection gating**
  - ConnectionGate callback in conn.Manager — checked before any WebRTC resource allocation
  - Inbound offers from non-whitelisted peers are silently dropped (zero resource cost)
  - Outbound Connect() also checks gate before initiating WebRTC handshake
  - Gate combines TrustStore check + owner-registered ConnectionRequestHandler callback
- [x] **Contact management API**
  - `AddContact(agentID)` — whitelist a peer (TrustVerified)
  - `RemoveContact(agentID)` — remove from whitelist
  - `BlockAgent(agentID)` — block a peer (TrustBlocked)
  - `ListContacts()` — list all trust entries
  - `OnConnectionRequest(handler)` — register callback for unknown peer connection requests
- [x] **Defense-in-depth architecture**
  - Layer 1 (Agent): TrustStore + EWMA reputation as primary defense
  - Layer 2 (Server): contacts service as secondary defense on signaling relay
  - connection_request signaling message type for owner notification

### Phase 8.5 — Agent Access Control (Complete)

Three-tier access model for the invoke endpoint, bringing production-grade gating to the marketplace.

- [x] **Phase 0: Mandatory Auth + Playground Gating** — Invoke endpoint requires authentication (agent headers or JWT); `playground_enabled` flag per agent controls open playground access
- [x] **Phase 1: Visibility Control** — `visibility` column (public/private); private agents hidden from directory; rate limiting differentiates agent-to-agent vs user invocations
- [x] **Phase 2: User ACL with Application/Approval** — `agent_user_acl` table with pending/approved/rejected status; providers approve/reject/revoke access requests with optional expiry; contacts support `expires_at` for time-limited partnerships
- [x] Frontend: playground toggle, visibility selector in publish wizard; access request dialog on agent profiles; provider access request management UI; user access requests page
- [x] API: 6 new endpoints for access request CRUD; dual-auth invoke (agent headers OR JWT)

### Phase 8.6 — LLM Tool Calling Integration (Complete)

Wrap PeerClaw agent capabilities as MCP-compatible tool definitions for LLM-driven agents.

- [x] **`agent/tools/` package** — 8 MCP tools (discover_agents, invoke_agent, get_agent_profile, check_reputation, add_contact, remove_contact, list_contacts, send_message)
- [x] **AgentAPI interface** — Abstraction over `*agent.Agent` for testability
- [x] **Handler dispatch** — Map-based dispatch with JSON-in/JSON-out, Result wrapper
- [x] **APIClient** — Thin HTTP client for server directory/invoke/reputation APIs
- [x] **Tool schema** — JSON Schema definitions for all 8 tools via `AllTools() → []agentcard.Tool`
- [x] **Tests** — 24 unit tests (mock AgentAPI, dispatch, validation, httptest.Server)

## Phase 9: Request-Response & Agent Collaboration Primitives

Enable synchronous agent-to-agent task delegation — the foundation for multi-agent collaboration.

- [ ] **SendRequest** — `Agent.SendRequest(ctx, env, timeout) → (*Envelope, error)` with TraceID-based correlation
  - Pending request registry (`map[traceID]chan *Envelope`)
  - Response matching in `HandleIncomingEnvelope` by TraceID + MessageType=response
  - Context-based timeout and cancellation
- [ ] **Envelope response helper** — `envelope.NewResponse(request, payload)` constructor
- [ ] **Broadcast** — `Agent.Broadcast(ctx, env, destinations) → []error` for fan-out messaging
- [ ] **A2A Task lifecycle mapping** — Map Envelope exchanges to A2A Task states (submitted → working → completed/failed)

## Phase 10: Inbound Handler Router

Enable agents to serve requests from other agents with capability-based routing.

- [ ] **Handler registration** — `Agent.Handle(capability, func(ctx, *Envelope) (*Envelope, error))` pattern
- [ ] **Automatic routing** — Incoming envelopes dispatched by metadata capability/action field
- [ ] **Auto-response** — Handler return value automatically wrapped as response envelope and sent back
- [ ] **Agent Card auto-generation** — Registered handler names populate Card.Capabilities automatically
- [ ] **Middleware support** — Pre/post hooks for logging, auth, rate limiting on inbound handlers

## Phase 11: Enterprise Simplified Configuration

Lower the barrier for enterprise intranet deployments.

- [ ] **`agent.NewSimple()`** — Simplified constructor (Name, ServerURL, KeypairPath, Capabilities only)
  - Auto-configures: Serverless=false, no Nostr, no DHT, no STUN/TURN
- [ ] **Enterprise deployment guide** — Single peerclaw-server + multiple agents on internal network
- [ ] **Pre-provisioned trust** — Bulk import of verified contacts for managed environments

## Phase 12: Nostr Relay Mailbox (Offline Message Delivery)

Upgrade Nostr relay from fallback transport to encrypted mailbox, enabling reliable offline message delivery while preserving P2P purity.

- [ ] **Encrypted inbox events** — Publish XChaCha20-encrypted envelopes as Nostr gift-wrapped events (NIP-59/NIP-44)
- [ ] **Inbox relay configuration** — Agent Profile extended with `inbox_relays` field (analogous to NIP-65)
- [ ] **Delivery flow** — Send() priority: 1) WebRTC direct → 2) Signaling relay → 3) Nostr inbox relay
- [ ] **Inbox sync** — On agent startup, query inbox relays for events since last sync timestamp
- [ ] **Delivery confirmation** — Encrypted delivery receipts sent back through the same channel
- [ ] **Local outbox with retry** — Sender persists unconfirmed messages; exponential backoff retry
- [ ] **TTL and cleanup** — Configurable message expiry (default 7 days); auto-cleanup of delivered events
- [ ] **Wake-up signaling** — Lightweight notification via WebSocket / Webhook / Cron when inbox has pending messages

## Phase 13: OpenClaw SKILL.md Integration

Quick integration with the OpenClaw ecosystem via a SKILL.md skill file.

- [ ] **SKILL.md authoring** — Markdown skill file teaching OpenClaw to use PeerClaw CLI/API for discovery, messaging, and contact management
- [ ] **CLI enhancements for SKILL.md** — Ensure all operations needed by SKILL.md are available as CLI commands (inbox check, contact approve, etc.)
- [ ] **Publish to ClawHub** — Submit PeerClaw skill to the OpenClaw skill registry

## Phase 14: OpenClaw Channel Plugin (Deep Integration)

PeerClaw as a native OpenClaw communication channel — like WhatsApp, Telegram, or Slack.

- [ ] **Channel plugin** — OpenClaw channel plugin that connects to PeerClaw agent network
- [ ] **Bidirectional messaging** — Incoming P2P messages surfaced in OpenClaw; OpenClaw responses sent back via PeerClaw
- [ ] **WebSocket bridge** — PeerClaw agent maintains WebSocket connection to OpenClaw gateway (port 18789) for real-time event push
- [ ] **Agent identity binding** — OpenClaw instance's identity mapped to PeerClaw Ed25519 keypair

## Phase 15: Protocol Ecosystem Integration

Deep integration with the three major agent protocols, making PeerClaw a native participant in each ecosystem.

### Phase 15a: MCP Server (`peerclaw-mcp`)

Standalone MCP Server binary — any MCP Host (Claude Code, VS Code Copilot, Cursor, etc.) can use PeerClaw as a tool provider.

- [ ] **`peerclaw-mcp` stdio binary** — Wraps `agent/tools/` as a proper MCP Server using `github.com/modelcontextprotocol/go-sdk`
- [ ] **Tool registration** — 8 PeerClaw tools registered with JSON Schema, descriptions, and MCP Tool Annotations (`readOnlyHint`, `destructiveHint`, `idempotentHint`)
- [ ] **Streamable HTTP transport** — Optional HTTP mode (`--transport http --port 8081`) for remote MCP hosting, alongside default stdio
- [ ] **Resource exposure** — Agent Card, trust store entries, and reputation data exposed as MCP Resources (`resources/list`, `resources/read`)
- [ ] **Session management** — `initialize` handshake, `Mcp-Session-Id` header, capability negotiation
- [ ] **Configuration** — `claude_desktop_config.json` / VS Code `settings.json` example configs for one-click setup

### Phase 15b: A2A Task Model & HTTP Bridge

Expose PeerClaw agents as standard A2A HTTP endpoints — any A2A client can discover and invoke PeerClaw agents.

- [ ] **A2A Task model mapping** — Map PeerClaw Envelope request-response to A2A Task lifecycle (submitted → working → input-required → completed/failed/canceled)
- [ ] **A2A HTTP endpoints** — `POST /a2a` JSON-RPC handler (`message/send`, `tasks/get`, `tasks/cancel`) backed by PeerClaw bridge
- [ ] **Agent Card serving** — `GET /.well-known/agent.json` auto-generated from PeerClaw agent registration data
- [ ] **Streaming support** — A2A SSE streaming mapped to PeerClaw's existing SSE invoke flow
- [ ] **Push notifications** — A2A push notification support for long-running tasks (webhook callback URL)
- [ ] **Multi-turn sessions** — A2A `contextId` mapped to PeerClaw `session_id` for stateful conversations

### Phase 15c: Agent Client Protocol Bridge

ndJSON/stdio bridge enabling ACP-compatible agents (OpenClaw, Zed AI, Coder) to join the PeerClaw network.

- [ ] **ACP stdio adapter** — ndJSON/stdio bridge process using `github.com/coder/acp-go-sdk`, translates ACP messages ↔ PeerClaw Envelopes
- [ ] **Agent Manifest translation** — PeerClaw Agent Card ↔ ACP Agent Manifest bidirectional mapping
- [ ] **Session/Run lifecycle** — ACP Session + Run model mapped to PeerClaw sessions; Run states mapped to Envelope exchanges
- [ ] **OpenClaw integration** — ACP bridge as OpenClaw's native channel for PeerClaw network access (complements Phase 14 Channel Plugin)
- [ ] **Enterprise intranet mode** — Simplified ACP bridge for corporate environments: single peerclaw-server + multiple ACP agent processes on internal network, no Nostr/DHT/STUN
- [ ] **Multi-agent orchestration** — ACP's `context_transfers` and `event_stream` mapped to PeerClaw broadcast/handler primitives

### Phase 15d: Universal Protocol Gateway

Unified ingress that auto-detects and routes any agent protocol.

- [ ] **Protocol auto-detection** — Inbound connections identified by content-type and payload structure (JSON-RPC → A2A/MCP, ndJSON → ACP, binary → PeerClaw native)
- [ ] **Unified routing** — Single gateway endpoint that dispatches to appropriate protocol adapter
- [ ] **Protocol translation matrix** — Bidirectional translation between all protocol pairs (A2A ↔ MCP ↔ ACP ↔ PeerClaw), extending Phase 3 adapters with real-world handling
- [ ] **Gateway metrics** — Per-protocol request counts, translation latency, error rates exposed via OpenTelemetry
