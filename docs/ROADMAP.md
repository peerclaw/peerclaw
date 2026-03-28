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

Evolve toward decentralization with federation, reputation, and identity anchoring.

- [x] **Interface abstractions + Agent Card extensions**
  - Discovery interface (RegistryClient)
  - SignalingClient interface (WebSocket)
  - Agent struct refactored to use interfaces (backward-compatible)
  - PeerClawExtension new fields (nostr_pubkey / reputation_score / nostr_relays / identity_anchor)
  - New signaling message types (federation_forward)
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
- [x] **Offline message buffering**
  - MessageCache for offline message buffering (per-destination queues, TTL expiration, JSON persistence)
  - OnPeerAdded callback (flushes cached queue when a new peer connects)
- [x] **Nostr identity anchoring (optional)**
  - IdentityAnchor interface (Publish / Verify / Resolve)
  - NostrAnchor implementation (Nostr kind 10078 replaceable event, bidirectional key binding)
- [x] **CLI Phase 5 commands**
  - `peerclaw federation status|peers`
  - `peerclaw reputation show|list`
  - `peerclaw identity anchor|verify`
- [x] **Integration tests**
  - Malicious behavior triggering reputation isolation
  - Reputation gossip propagation across peers
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

## Phase 7: Agent Platform (Complete)

Evolve PeerClaw into an Agent Platform — where anyone can register an Agent and anyone (human or Agent) can discover and invoke it.

### Phase 7a: Directory & Profile

Public-facing directory for discovering and evaluating Agents.

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
- [x] **Agent registration wizard** — Guided 5-step registration (basic info → capabilities & protocols → endpoint → auth & metadata → preview)
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

Three-tier access model for the invoke endpoint, bringing production-grade gating to the platform.

- [x] **Phase 0: Mandatory Auth + Playground Gating** — Invoke endpoint requires authentication (agent headers or JWT); `playground_enabled` flag per agent controls open playground access
- [x] **Phase 1: Visibility Control** — `visibility` column (public/private); private agents hidden from directory; rate limiting differentiates agent-to-agent vs user invocations
- [x] **Phase 2: User ACL with Application/Approval** — `agent_user_acl` table with pending/approved/rejected status; providers approve/reject/revoke access requests with optional expiry; contacts support `expires_at` for time-limited partnerships
- [x] Frontend: playground toggle, visibility selector in registration wizard; access request dialog on agent profiles; provider access request management UI; user access requests page
- [x] API: 6 new endpoints for access request CRUD; dual-auth invoke (agent headers OR JWT)

### Phase 8.6 — LLM Tool Calling Integration (Complete)

Wrap PeerClaw agent capabilities as MCP-compatible tool definitions for LLM-driven agents.

- [x] **`agent/tools/` package** — 8 MCP tools (discover_agents, invoke_agent, get_agent_profile, check_reputation, add_contact, remove_contact, list_contacts, send_message)
- [x] **AgentAPI interface** — Abstraction over `*agent.Agent` for testability
- [x] **Handler dispatch** — Map-based dispatch with JSON-in/JSON-out, Result wrapper
- [x] **APIClient** — Thin HTTP client for server directory/invoke/reputation APIs
- [x] **Tool schema** — JSON Schema definitions for all 8 tools via `AllTools() → []agentcard.Tool`
- [x] **Tests** — 24 unit tests (mock AgentAPI, dispatch, validation, httptest.Server)

## Phase 9: Request-Response & Agent Collaboration Primitives (Complete)

Enable synchronous agent-to-agent task delegation — the foundation for multi-agent collaboration.

- [x] **SendRequest** — `Agent.SendRequest(ctx, env, timeout) → (*Envelope, error)` with TraceID-based correlation
  - Pending request registry (`map[traceID]chan *Envelope`)
  - Response matching in `HandleIncomingEnvelope` by TraceID + MessageType=response
  - Context-based timeout and cancellation
- [x] **Envelope response helper** — `envelope.NewResponse(request, payload)` constructor
- [x] **Broadcast** — `Agent.Broadcast(ctx, env, destinations) → map[string]error` for fan-out messaging
- [x] **A2A Task lifecycle mapping** — TaskTracker maps Envelope exchanges to A2A Task states (submitted → working → completed/failed)

## Phase 10: Inbound Handler Router (Complete)

Enable agents to serve requests from other agents with capability-based routing.

- [x] **Handler registration** — `Agent.Handle(capability, func(ctx, *Envelope) (*Envelope, error))` pattern
- [x] **Automatic routing** — Incoming envelopes dispatched by `metadata["capability"]` field via Router
- [x] **Auto-response** — Handler return value automatically wrapped via `envelope.NewResponse()` and sent back
- [x] **Agent Card auto-generation** — `Agent.Capabilities()` returns deduplicated union of opts.Capabilities + router-registered capabilities
- [x] **Middleware support** — `Middleware` type with `Use()` API; built-in `LoggingMiddleware` and `RecoveryMiddleware`

## Phase 11: Enterprise Simplified Configuration (Complete)

Lower the barrier for enterprise intranet deployments.

- [x] **`agent.NewSimple()`** — Simplified constructor (Name, ServerURL, Capabilities variadic)
  - Auto-configures: server-only discovery and signaling, no STUN/TURN
  - Auto-generates Ed25519 keypair, server-only discovery and signaling
- [x] **`agent.ImportContacts()`** — Bulk import of verified contacts for managed environments
  - Sets all imported agent IDs to TrustVerified level
- [x] **Enterprise deployment guide** — `docs/ENTERPRISE.md` / `docs/ENTERPRISE_zh.md`
  - Architecture diagram, quick start, Docker Compose, security recommendations
- [x] **Enterprise example** — `agent/examples/enterprise/main.go`

## Phase 12: Nostr Relay Mailbox (Offline Message Delivery) (Complete)

Upgrade Nostr relay from fallback transport to encrypted mailbox, enabling reliable offline message delivery while preserving P2P purity.

- [x] **Encrypted inbox events** — Publish NIP-44-encrypted envelopes as Nostr events (kind 20007) to recipient's inbox relays
- [x] **Inbox relay configuration** — `PeerClawExtension.InboxRelays` field (analogous to NIP-65), distinct from `NostrRelays` (real-time transport)
- [x] **Delivery flow** — Send() priority: 1) Existing P2P → 2) New P2P (15s) → 3) Signaling relay → 4) Nostr mailbox → 5) MessageCache
- [x] **Inbox sync** — Periodic inbox polling (`SyncInterval`, default 5 min) queries inbox relays for events since last sync timestamp
- [x] **Delivery confirmation** — Encrypted delivery receipts (kind 20008) sent back through the same channel; outbox entries marked as confirmed
- [x] **Local outbox with retry** — Sender persists unconfirmed messages; exponential backoff retry (5s base, 5 min max, 10 retries)
- [x] **TTL and cleanup** — Configurable message expiry (default 7 days); expired and confirmed entries auto-cleaned from outbox
- [x] **Wake-up signaling** — `mailbox_wakeup` signaling message type for lightweight notification when inbox has pending messages

## Phase 15a: MCP Server — `peerclaw mcp serve` (Complete)

MCP Server integrated into the CLI — any MCP Host (Claude Code, VS Code Copilot, Cursor, Windsurf, etc.) can use PeerClaw as a tool provider.

- [x] **`peerclaw mcp serve` command** — Wraps `agent/tools/` as a proper MCP Server using `github.com/modelcontextprotocol/go-sdk` (v1.4.0)
- [x] **Tool registration** — 4 API-mode tools (discover_agents, invoke_agent, get_agent_profile, check_reputation) with JSON Schema and MCP Tool Annotations (`readOnlyHint`, `idempotentHint`, `destructiveHint`)
- [x] **Dual transport** — stdio (default) and Streamable HTTP (`--transport http --port 8081`)
- [x] **Resource exposure** — Agent directory exposed as MCP Resource (`peerclaw://directory`)
- [x] **Configuration guide** — `docs/mcp-config.md` with examples for Claude Code, VS Code, Cursor, Windsurf

## Phase 13: CLI Completeness & SKILL.md (Complete)

Fill CLI gaps and author an OpenClaw SKILL.md for AI-driven agent orchestration.

- [x] **`peerclaw invoke` command** — Direct agent invocation (`peerclaw invoke <agent-id> --message "..."`) with `--protocol`, `--session-id`, `--stream` flags; SSE streaming with real-time output; no source agent required
- [x] **`peerclaw inbox` command** — Access request management: `request` (submit access request), `status` (check request status), `list` (list all my requests); JWT auth via `--token` flag or `PEERCLAW_TOKEN` env
- [x] **`peerclaw agent update` subcommand** — Update agent fields (name, description, version, capabilities, endpoint, protocols) without re-registration; JWT auth required
- [x] **SKILL.md authoring** — `docs/SKILL.md` — Markdown skill file describing PeerClaw CLI commands and REST API for discovery, invocation, access management, and reputation checking
- [ ] **Publish to ClawHub** — Submit PeerClaw skill to the OpenClaw skill registry

## Phase 15b: A2A HTTP Bridge (Complete)

Expose PeerClaw agents as standard A2A HTTP endpoints — any A2A client can discover and invoke PeerClaw agents.

- [x] **A2A Task model mapping** — Map PeerClaw Envelope request-response to A2A Task lifecycle (accepted → working → completed/failed/canceled)
- [x] **A2A HTTP endpoints** — `POST /a2a/{agent_id}` JSON-RPC handler (`message/send`, `message/send/subscribe`, `tasks/get`, `tasks/cancel`, `tasks/pushNotification/set|get`) backed by PeerClaw bridge
- [x] **Agent Card serving** — `GET /a2a/{agent_id}/.well-known/agent.json` auto-generated from PeerClaw agent registration data with endpoint, capabilities, and skills
- [x] **Streaming support** — A2A SSE streaming via `message/send/subscribe` or `Accept: text/event-stream`, each SSE event is a full JSON-RPC Response wrapping the Task
- [x] **Push notifications** — A2A push notification config storage (`tasks/pushNotification/set|get`) for long-running tasks
- [x] **Multi-turn sessions** — A2A `contextId` mapped to PeerClaw `session_id` for stateful conversations
- [x] **REST convenience endpoint** — `GET /a2a/{agent_id}/tasks/{task_id}` for task status polling
- [x] **Access control** — External A2A clients treated as anonymous users, gated by `playground_enabled` flag
- [x] **Rate limiting** — Per-IP rate limiting via `invokeRateLimiter` for A2A bridge requests
- [x] **Task cleanup** — Background goroutine cleans expired tasks (1 hour TTL)

## Phase 14: Multi-Platform Channel Integration (Complete)

PeerClaw as a native communication channel across AI orchestration platforms — platform adapter abstraction with plugins for 4 platforms.

- [x] **Platform Adapter interface** — `platform.Adapter` abstraction in agent SDK: `Connect()`, `SendChat()`, `InjectNotification()`, `SetOutboundHandler()` — platform-agnostic integration point
- [x] **OpenClaw adapter** — WebSocket gateway client (`agent/platform/openclaw/`) with req/res/event frame protocol, connect handshake, auto-reconnect
- [x] **IronClaw adapter** — HTTP/SSE gateway client (`agent/platform/ironclaw/`) with REST chat.send + SSE event streaming, bearer token auth
- [x] **Bridge adapter** — Generic local WebSocket bridge (`agent/platform/bridge/`) with simple JSON protocol for platforms without external APIs
- [x] **OpenClaw plugin** — TypeScript npm package (`@peerclaw/openclaw-plugin`) using `openclaw/plugin-sdk` external plugin API
- [x] **IronClaw plugin** — Rust WASM component (`peerclaw-ironclaw-plugin`) implementing `sandboxed-channel` WIT interface
- [x] **nanobot plugin** — Python pip package (`nanobot-channel-peerclaw`) implementing `BaseChannel` with entry-point auto-discovery
- [x] **PicoClaw plugin** — Go module (`peerclaw/picoclaw-plugin`) with `channels.RegisterFactory()` + `init()` self-registration
- [x] **Notification forwarding** — Server notifications pushed via signaling to agent, forwarded to platform conversations
- [x] **Bidirectional messaging** — Incoming P2P messages forwarded to platform AI; AI responses routed back via P2P

## Phase 15c: ACP HTTP Bridge (Complete)

Expose PeerClaw agents as standard ACP HTTP endpoints — any ACP client can discover and invoke PeerClaw agents.

- [x] **ACP Run model mapping** — Map PeerClaw Envelope request-response to ACP Run lifecycle (created → in-progress → completed/failed/cancelled)
- [x] **ACP HTTP endpoints** — `POST /acp/{agent_id}/runs` REST handler with sync/stream/async modes, `GET /acp/{agent_id}/runs/{run_id}` for status polling, `POST /acp/{agent_id}/runs/{run_id}/cancel` for cancellation
- [x] **Agent Manifest serving** — `GET /acp/{agent_id}/agents` auto-generated from PeerClaw agent registration data with name, description, capabilities, content types
- [x] **Streaming support** — ACP SSE streaming via `mode: "stream"`, each SSE event is `event: run_update\ndata: <Run JSON>`
- [x] **Async mode** — `mode: "async"` returns HTTP 202 immediately, background goroutine executes bridge call with 5-minute timeout
- [x] **Ping endpoint** — `GET /acp/{agent_id}/ping` for health checks
- [x] **Access control** — External ACP clients treated as anonymous users, gated by `playground_enabled` flag
- [x] **Rate limiting** — Per-IP rate limiting via `invokeRateLimiter` for ACP bridge requests
- [x] **Run cleanup** — Background goroutine cleans expired runs (1 hour TTL)

## Phase 15d: Universal Protocol Gateway (Complete)

Per-agent MCP bridge + unified protocol gateway with auto-detection and multi-format discovery.

- [x] **Per-agent MCP Bridge** — `POST /mcp/{agent_id}` JSON-RPC handler (`initialize`, `tools/list`, `tools/call`, `resources/list`, `prompts/list`) with session management, access control, rate limiting, and invocation logging
- [x] **MCP initialize** — Returns `InitializeResult` with `ServerInfo` from agent card, `Mcp-Session-Id` header for session tracking
- [x] **MCP tools mapping** — Agent card `Tools` automatically mapped to MCP `ToolDef` list; `tools/call` dispatches via bridge with `mcp.tool_name` envelope metadata
- [x] **MCP SSE endpoint** — `GET /mcp/{agent_id}` SSE placeholder for server-initiated notifications
- [x] **Universal Gateway invoke** — `POST /agent/{agent_id}` auto-detects protocol from request body and dispatches to A2A/MCP/ACP bridge handler
- [x] **Protocol auto-detection** — JSON-RPC `method` prefix matching (`message/` `tasks/` → A2A, `tools/` `resources/` `prompts/` `initialize` → MCP), `input`/`agent_name` fields → ACP, with params shape fallback
- [x] **Multi-format discovery** — `GET /agent/{agent_id}?format=a2a|mcp|acp` returns protocol-specific agent card (A2A AgentCard, MCP server info, ACP Manifest); default returns PeerClaw Card
- [x] **Gateway metrics** — `peerclaw.gateway.requests.total` OpenTelemetry counter with `protocol` attribute
- [x] **Session cleanup** — Background goroutine cleans expired MCP sessions (1 hour TTL)

## Phase 15e: ACP Stdio Bridge (Complete)

ndJSON/stdio bridge enabling ACP-compatible agents (OpenClaw, Zed AI, Coder) to join the PeerClaw network via local process communication.

- [x] **`peerclaw acp serve` command** — ndJSON/stdio bridge process, reads ACP requests from stdin and writes responses to stdout, proxying to PeerClaw server's ACP HTTP bridge
- [x] **ACP method dispatch** — 6 methods: `create_run` (POST /acp/{agent_id}/runs), `get_run` (with local cache), `cancel_run`, `list_agents` (directory → ACP manifest conversion), `get_agent`, `ping`
- [x] **Lightweight ACP types** — JSON-compatible copies of server ACP types (Run, Message, MessagePart, AgentManifest) in CLI package, no server module dependency
- [x] **Run caching** — Local sync.Map cache for created runs, reducing server round-trips on get_run
- [x] **Tests** — 10 unit tests (ping, invalid JSON, unknown method, create_run, get_run, get_run cached, cancel_run, list_agents, get_agent, blank lines) + 2 command tests

## Phase 16: P2P File Transfer (Complete)

Pure peer-to-peer large file transfer with E2E encryption — zero server dependency in the data path.

- [x] **WebRTC transport enhancement** — `CreateDataChannel()`, `RegisterDataChannelHandler()` for dedicated file transfer channels, backpressure control (1MB high-water, 256KB low-water) in `Send()`
- [x] **File transfer message types** — `file_offer`, `file_accept`, `file_reject`, `transfer_ready`, `transfer_complete`, `chunk_ack`, `resume_request`, `file_chunk` in `core/envelope/filetransfer.go`
- [x] **Binary frame protocol** — `[seq:4B][length:4B][flags:1B][encrypted_chunk]` with FlagData, FlagFIN, FlagACK; 64KB default chunk size
- [x] **Transfer state machine** — `Idle → Offered → Accepted → Transferring → Completing → Done/Failed/Cancelled` with per-state timeouts
- [x] **Challenge-response mutual auth** — 3-step Ed25519 handshake: FileOffer(challenge) → FileAccept(challenge_sig, counter_challenge) → TransferReady(counter_sig)
- [x] **Pipeline push sender** — Dedicated `ft-{file_id}` DataChannel (ordered, reliable), per-chunk XChaCha20-Poly1305 encryption with AAD = `file_id|seq`, FIN frame on completion
- [x] **Streaming receiver** — Binary frame decode → decrypt → write to file, periodic ChunkAck every 100 chunks, full-file SHA-256 verification on FIN
- [x] **Resume support** — `SaveResumeState()` / `LoadResumeState()` persist last-confirmed sequence to disk, `ResumeRequest` continues from `last_seq + 1`
- [x] **Nostr relay fallback** — When WebRTC ICE fails, chunks sent as encrypted Nostr events (~40KB/event) via standard `agent.Send()` path
- [x] **Mailbox wakeup** — FileOffer via mailbox sends `MessageTypeMailboxWakeup` to trigger immediate `SyncInbox()` instead of waiting for poll interval
- [x] **Agent integration** — `SendFile()`, `ListTransfers()`, `GetTransfer()`, `CancelTransfer()` public API; capability handler registered as `"file_transfer"`; `FileTransferDir` and `ResumeStatePath` options
- [x] **CLI commands** — `peerclaw send-file --to <id> --file <path>` with progress polling; `peerclaw transfer status [--transfer-id <id>]` for transfer listing
- [x] **Blob service removal** — Removed centralized `server/internal/blob/` package and all references; file transfer is now purely P2P

## Phase 17: Hardening & Developer Experience (Complete)

Security hardening, error handling consistency, and developer experience improvements based on security audit findings.

### Security Hardening

- [x] **Proof-of-Possession registration** — Agents must sign the request body with their Ed25519 key during `POST /agents` registration, preventing public key squatting
- [x] **Forward secrecy / session rekeying** — Automatic ephemeral X25519 rekey after 1000 messages or 1 hour; old key material is securely zeroed
- [x] **Encrypted trust store** — Trust store files encrypted at rest with XChaCha20-Poly1305 (key derived from Ed25519 seed via HKDF); transparent migration from plaintext
- [x] **Test coverage targets** — CI coverage gates: core 80%+, agent/server/cli 70%+

### Error Handling & API Quality

- [x] **Unified error types** — `core/errors` package with structured error codes (`not_found`, `validation_error`, `auth_error`, etc.) and `New`/`Wrap`/`Is` helpers
- [x] **Card.Validate() method** — Validation logic extracted from server to `Card.Validate()` in core; server delegates to it
- [x] **Structured error responses** — `jsonError` returns `{code, message}` format mapped from HTTP status codes

### CLI & Developer Experience

- [x] **Shell completion** — `peerclaw completion bash|zsh|fish` generates shell completion scripts with all commands, subcommands, and flags
- [x] **handleListAgents expansion** — Added `sort`, `search`, `min_score`, `page_size` query parameters (parity with public directory)

### Reputation Transparency

- [x] **Public reputation thresholds** — `ReputationLow` (0.3), `ReputationMedium` (0.7), `ReputationHigh` (0.8) constants in `core/agentcard/reputation.go`
- [x] **Reputation sorting & filtering** — `GET /api/v1/agents?sort=reputation&min_score=<float>` now supported

### Plugin Ecosystem

Plugin tier classification:

| Tier | Plugin | Language | Status |
|------|--------|----------|--------|
| Tier 1 | openclaw-plugin | TypeScript | Full tests + CI |
| Tier 1 | ironclaw-plugin | Rust WASM | Tests + CI |
| Tier 1 | zeroclaw-plugin | Rust | Tests + CI (NEW) |
| Community | picoclaw-plugin | Go | Community maintained |
| Community | nanobot-plugin | Python | Community maintained |

- [x] **zeroclaw-plugin** — New Rust crate implementing ZeroClaw's Channel trait via bridge WebSocket protocol
- [x] **openclaw-plugin tests** — vitest suite for config schema validation and account resolution
- [x] **ironclaw-plugin tests** — Unit tests for frame parsing, peer ID extraction, serialization
- [x] **Community plugin labels** — picoclaw and nanobot plugins marked as community maintained with CONTRIBUTING.md

## Phase 18: Web Dashboard Foundation (Complete)

Full-featured admin dashboard for the PeerClaw server with i18n support across 8 languages.

- [x] **Dashboard SPA** — React + TypeScript + Tailwind + shadcn/ui, embedded into server binary
- [x] **Admin panel** — Overview, users, agents, reports, categories, analytics, invocations pages
- [x] **Provider console** — Dashboard, agent management, discover, invocation history, API keys, notifications, profile
- [x] **Public pages** — Landing, directory, agent profiles, playground, about, login/register/forgot-password
- [x] **i18n** — 8 languages (en, zh, es, fr, ja, pt, ru, ar) with complete translation coverage
- [x] **Data export** — CSV/JSON export for admin tables
- [x] **Sortable tables** — Column sorting with sortable headers and loading skeletons
- [x] **Accessibility** — Skip-to-content, ARIA labels, keyboard navigation

## Phase 19-22: Dashboard Enhancements (Complete)

### Phase 19: Bulk Actions
- [x] **Bulk endpoints** — `POST /admin/{agents,reports,users}/bulk` with action-based dispatch
- [x] **SelectableTable** — Reusable component with checkbox selection and floating action bar
- [x] **Admin pages** — AgentsPage (verify/delete), UsersPage (delete), ReportsPage (review/dismiss/delete)

### Phase 20: Profile & Auth Improvements
- [x] **PasswordStrength** — 4-segment visual bar with requirements checklist, applied to register/profile/forgot-password
- [x] **Auto-dismiss** — Success messages on ProfilePage auto-dismiss after 3 seconds
- [x] **AboutPage** — Updated roadmap to Phase 4, added GitHub roadmap link

### Phase 21: Admin Audit Log
- [x] **adminaudit package** — Store/SQLite/Postgres/Service/Factory pattern, `admin_audit_log` table with indexes
- [x] **Audit recording** — All admin mutations logged (user.delete, agent.verify, report.update, category.create, etc.)
- [x] **AuditLogPage** — Filters (admin, action, target type, date range), pagination, sidebar nav item

### Phase 22: Advanced Dashboard
- [x] **Clickable stat cards** — Overview cards navigate to respective admin pages
- [x] **Recent Activity feed** — Last 10 audit events on overview dashboard
- [x] **Provider time range** — 7d/30d/All selector with backend `since` parameter support

## Phase 23-26: Security Hardening & Code Quality (Complete)

Systematic remediation of findings from the 2026-03-28 full project code audit.

### Phase 23: Critical Security Hardening
- [x] **SEC-C01** — Sanitize 90+ error message leaks across 13 handler files; zero 500-level leaks remain
- [x] **SEC-C02** — Fix X-Forwarded-For IP spoofing via `BridgeClientIP()` proxy trust validation
- [x] **SEC-H04** — CLI config file permissions restricted from 0644 to 0600
- [x] **SEC-M06** — Nil public key check added to `Verify()` to prevent panic

### Phase 24: Auth & Access Control Hardening
- [x] **SEC-H01** — Agent ID auth bypass closed (AgentExists registry check on pubKey fallback)
- [x] **SEC-H03** — owner_user_id IDOR prevention (strip from metadata in handleRegister)
- [x] **SEC-M01** — OTP strengthened from 6 to 8 digits + brute-force lockout (5 failures → 15min block)
- [x] **SEC-M03** — `decodeJSON()` helper with 64KB limit on all 11 auth endpoints

### Phase 25: Risk Code & Stability
- [x] **RISK-01** — A2A Bridge TOCTOU race fixed (atomic CompareAndSwap loop)
- [x] **RISK-02** — WebSocket read loop panic recovery (defer/recover)
- [x] **SEC-M05** — React ErrorBoundary wrapping all 3 route groups

### Phase 26: Code Quality & DRY
- [x] **DUP-05** — Shared status badge utility (`lib/status.ts`)
- [x] **DUP-04** — `useAsyncAction` hook for consistent async error handling
- [x] **DEAD-04** — Envelope `WithTTL()` and `GenerateNonce()` methods implemented
- [x] **Performance** — Admin routes lazy-loaded via React.lazy (6 pages code-split, -47KB main bundle)
