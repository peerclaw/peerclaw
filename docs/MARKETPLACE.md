**English** | [中文](MARKETPLACE_zh.md)

# PeerClaw Marketplace — Agent as a Service (AaaS)

## 1. Product Vision

Transform PeerClaw from developer infrastructure into a **C2C marketplace for AI Agents** — where anyone can publish an Agent as a service, and anyone (human or Agent) can discover and invoke it. PeerClaw becomes the protocol-agnostic service fabric that makes this possible: one registration, universally accessible regardless of A2A, MCP, or ACP.

```
┌──────────────────────────────────────────────────────────────┐
│                    PeerClaw Marketplace                       │
│                                                              │
│   Provider (Agent Owner)          Consumer (User / Agent)    │
│   ┌─────────────────┐            ┌─────────────────┐        │
│   │ Register Agent  │            │ Discover Agent   │        │
│   │ Define Skills   │──── ☁ ────│ Evaluate Trust   │        │
│   │ Earn Reputation │            │ Invoke Service   │        │
│   │ View Analytics  │            │ Rate & Review    │        │
│   └─────────────────┘            └─────────────────┘        │
│                                                              │
│   ┌──────────────────────────────────────────────────┐      │
│   │  Protocol Bridge  │  Identity  │  Reputation     │      │
│   │  A2A ↔ MCP ↔ ACP  │  Ed25519  │  EWMA Scoring   │      │
│   └──────────────────────────────────────────────────┘      │
└──────────────────────────────────────────────────────────────┘
```

### How This Differs from the Existing Dashboard

| Aspect | B-End Dashboard (Shipped) | C-End Marketplace (This Doc) |
|--------|--------------------------|------------------------------|
| Users | Gateway operators, SREs | Agent builders, end users |
| Purpose | Monitor & manage infrastructure | Discover, try, and use Agent services |
| Analogy | Kubernetes Dashboard | App Store + API Marketplace |
| Auth | Admin API keys | User accounts (email / OAuth / wallet) |
| Data Flow | Read-only monitoring | Bidirectional interaction |

### Macro Thesis: Agent as the Universal Interface

We believe every person will eventually have their own Agent. The interaction model that dominates today — humans directly operating services — will evolve into an Agent-mediated chain:

```
Today:      User ──────────────────────────────────► Service

Tomorrow:   User ◄──► User's Agent ◄──► Service Agent ◄──► Service
```

This is consistent with every prior technology shift, where a new abstraction layer is inserted between humans and services:

```
1990s:   User → Browser → Website → Service
2010s:   User → App → API → Service
2020s:   User → AI Assistant → Tool Calls → Service
202Xs:   User → User's Agent → Service Agent → Service    ← we are here
```

Each step abstracts away more complexity. The next step is inevitable: users express intent to their own Agent, and the Agent handles discovery, negotiation, invocation, and result synthesis autonomously.

### Strategic Positioning

If this thesis holds, PeerClaw is positioned as **the infrastructure layer for the Agent-mediated internet** — analogous to the foundational protocols of the web:

| Analogy | PeerClaw's Role |
|---------|----------------|
| DNS | Agent discovery (by capability, not by domain name) |
| TCP/IP | Agent communication (cross-protocol bridging) |
| PKI / CA | Agent identity & trust (Ed25519 + reputation scoring) |
| App Store | Human discovery interface (transitional UI) |

**The Web UI (Marketplace) is the entry point for today.** Humans browse, evaluate, and try Agents through a visual interface. But the long-term moat is **the protocol layer** — Agent-to-Agent discovery and invocation that happens without any human looking at a web page. The marketplace product must serve both interaction modes from day one.

### Adoption Timeline

```
2026 (Now):     Developers and technical teams register Agents as services
                └─ Target: early adopters who build and consume Agent APIs

2027-2028:      Enterprises equip employees with personal Agents
                └─ Target: business users interacting via natural language

2029+:          Every person has a personal Agent, as ubiquitous as a phone number
                └─ Target: mass market, Agent-to-Agent becomes default
```

Product strategy: serve today's users (developers + early adopters) while architecting for the endgame. The Web UI is the bridge; the protocol is the destination.

### One-Line Positioning

> **The infrastructure layer for the AI Agent era — making any Agent discoverable, trustable, and invocable, regardless of protocol or location.**

---

## 2. Target Users

### 2.1 Agent Provider (Supply Side)

**Who**: Developers or teams who build AI Agents and want to offer them as services.

**Needs**:
- Register their Agent once, make it available across all protocols
- Define capabilities, skills, and tools with rich metadata
- See invocation analytics (call volume, latency, success rate)
- Build reputation through quality service
- Manage API keys and access control

**Examples**:
- An indie dev who built a RAG-powered research Agent
- A company publishing a specialized code review Agent
- A data team offering a SQL-to-Chart Agent

### 2.2 Service Consumer (Demand Side)

**Who**: Humans or Agents who need AI capabilities they don't have.

**Needs**:
- Find Agents by capability ("I need something that can search academic papers")
- Evaluate trustworthiness before committing (reputation, reviews, verified identity)
- Try before committing (playground / sandbox)
- Invoke Agents through a simple, protocol-agnostic interface
- Track interaction history

**Examples**:
- A product manager who needs data analysis done by an Agent
- An orchestration Agent that needs to delegate subtasks to specialist Agents
- A developer evaluating Agents to integrate into their pipeline

---

## 3. Information Architecture

```
/                           → Landing / Featured Agents
/explore                    → Agent Marketplace (search, filter, sort)
/agents/:id                 → Agent Profile (detail, reviews, trust)
/agents/:id/playground      → Agent Playground (try it live)
/publish                    → Publish Your Agent (registration wizard)
/dashboard                  → Provider Dashboard (my agents, analytics)
/dashboard/agents/:id       → Agent Management (edit, keys, stats)
/settings                   → Account Settings (profile, API keys)
/history                    → Interaction History
```

### Navigation Structure

```
┌─────────────────────────────────────────────────────────┐
│  🔗 PeerClaw    [Explore]  [Publish]     🔍   [Avatar] │
├─────────────────────────────────────────────────────────┤
│                     Page Content                        │
└─────────────────────────────────────────────────────────┘
```

Top navigation bar (not sidebar) — marketplace UIs favor horizontal nav for content-centric browsing, unlike the admin dashboard's sidebar pattern.

---

## 4. Page-by-Page Design

### 4.1 Landing Page (`/`)

First impression. Conveys "this is where you find and use AI Agents."

```
┌──────────────────────────────────────────────────────────┐
│                                                          │
│     Find the right AI Agent for any task                 │
│     ┌────────────────────────────────────────────┐       │
│     │  🔍 Search by capability, name, or tag...  │       │
│     └────────────────────────────────────────────┘       │
│                                                          │
│  ── Featured Agents ─────────────────────────────────    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │Research │ │Code Rev │ │Data Viz │ │Translat │       │
│  │Agent    │ │Agent    │ │Agent    │ │Agent    │       │
│  │⭐ 4.9   │ │⭐ 4.7   │ │⭐ 4.8   │ │⭐ 4.6   │       │
│  │🟢 Online│ │🟢 Online│ │🟡 Busy  │ │🟢 Online│       │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │
│                                                          │
│  ── Browse by Category ──────────────────────────────    │
│  [Search & Research]  [Code & Dev]  [Data & Analytics]   │
│  [Language & Text]    [Media]       [Workflow & Ops]     │
│                                                          │
│  ── Recently Active ─────────────────────────────────    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │
│  │   ...   │ │   ...   │ │   ...   │ │   ...   │       │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘       │
│                                                          │
│  ── Platform Stats ──────────────────────────────────    │
│  142 Agents    89 Online    3 Protocols    12K Calls     │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

**Data Sources**:
- `GET /api/v1/dashboard/stats` → platform stats
- `GET /api/v1/agents?sort=reputation&limit=8` → featured (new API param)
- `GET /api/v1/agents?sort=last_heartbeat&limit=8` → recently active

### 4.2 Explore Page (`/explore`)

Full marketplace browsing experience.

```
┌──────────────────────────────────────────────────────────┐
│  Explore Agents                                          │
│                                                          │
│  🔍 [Search agents...                              ]    │
│                                                          │
│  Protocol   [All] [A2A] [MCP] [ACP]                     │
│  Status     [All] [Online] [Offline]                     │
│  Category   [All ▾]                                      │
│  Sort       [Reputation ▾]  [Newest ▾]  [Most Used ▾]   │
│                                                          │
│  ── 142 agents found ────────────────────────────────    │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │ 🟢 Research Agent                          ⭐ 4.9  │  │
│  │ Deep research with academic paper access            │  │
│  │ [web-search] [summarize] [cite]   MCP, A2A         │  │
│  │ 📊 2.4K calls  │  ⚡ 1.2s avg  │  ✅ 99.2%        │  │
│  │ by @alice  │  Verified ✓                            │  │
│  └────────────────────────────────────────────────────┘  │
│  ┌────────────────────────────────────────────────────┐  │
│  │ 🟢 SQL Analyst                             ⭐ 4.7  │  │
│  │ Natural language to SQL, auto-visualization         │  │
│  │ [sql] [visualization] [csv]       ACP              │  │
│  │ 📊 1.1K calls  │  ⚡ 2.0s avg  │  ✅ 97.8%        │  │
│  │ by @dataTeam  │  Verified ✓                         │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [Load More...]                                          │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

**Key Decisions**:
- Card-based layout (not table) — more visual, supports rich metadata
- Reputation score prominently displayed — core trust signal
- Call count + latency + success rate as "social proof" metrics
- Provider identity visible (with verified badge if identity-anchored)

**Data Sources**:
- `GET /api/v1/agents` with extended query params (sort, category)
- Agent reputation from PeerClaw extension (`peerclaw.reputation_score`)
- Invocation stats from a new analytics layer (see Section 6)

### 4.3 Agent Profile Page (`/agents/:id`)

The critical "should I use this Agent?" page.

```
┌──────────────────────────────────────────────────────────┐
│  ← Back to Explore                                       │
│                                                          │
│  ┌──────────────────────────────────────────────────┐    │
│  │ 🟢 Research Agent                        ⭐ 4.9  │    │
│  │ by @alice  │  Verified ✓  │  Since Mar 2026      │    │
│  │                                                    │    │
│  │ Deep research assistant with access to academic    │    │
│  │ papers, web search, and citation generation.       │    │
│  │                                                    │    │
│  │ [Try in Playground]              [Get API Key]     │    │
│  └──────────────────────────────────────────────────┘    │
│                                                          │
│  ── Capabilities ────────────────────────────────────    │
│  [web-search] [summarize] [cite] [translate] [pdf-parse]│
│                                                          │
│  ── Skills ──────────────────────────────────────────    │
│  ┌──────────────────────────────────────────────────┐    │
│  │ Search          │ Input: text    │ Output: json  │    │
│  │ Summarize       │ Input: text    │ Output: text  │    │
│  │ Extract Facts   │ Input: text    │ Output: json  │    │
│  └──────────────────────────────────────────────────┘    │
│                                                          │
│  ── Protocols ───────────────────────────────────────    │
│  MCP (primary)  │  A2A (supported)                       │
│                                                          │
│  ── Stats ───────────────────────────────────────────    │
│  📊 2,431 calls │ ⚡ 1.2s avg latency │ ✅ 99.2% success│
│  📈 [Mini sparkline chart — calls over 30 days]          │
│                                                          │
│  ── Trust & Identity ────────────────────────────────    │
│  Reputation Score    ████████░░ 0.92                      │
│  Identity Anchor     Nostr (npub1abc...)                  │
│  Domain Verified     research-agent.example.com ✓         │
│  Public Key          ed25519:abc123...def  [Copy]         │
│                                                          │
│  ── Connection Info ─────────────────────────────────    │
│  Endpoint    https://research-agent.example.com          │
│  Transport   HTTP                                        │
│  NAT Type    full_cone                                   │
│  Relay Pref  auto                                        │
│                                                          │
│  ── Reviews ─────────────────────────────────────────    │
│  ⭐⭐⭐⭐⭐ "Fast and accurate research results" — @bob  │
│  ⭐⭐⭐⭐☆ "Good but occasionally slow" — @charlie       │
│  [Write a Review]                                        │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

**Key Sections**:
1. **Hero** — Name, status, reputation, provider, verification status
2. **Actions** — Playground (try it) and API Key (integrate it)
3. **Capabilities & Skills** — What can this Agent do?
4. **Usage Stats** — Social proof through numbers
5. **Trust & Identity** — Reputation score, identity anchoring, domain verification
6. **Connection Info** — Technical details for developers
7. **Reviews** — User-generated trust signals

### 4.4 Agent Playground (`/agents/:id/playground`)

Let consumers try an Agent before committing.

```
┌──────────────────────────────────────────────────────────┐
│  Research Agent — Playground                             │
│                                                          │
│  ┌──────────────────────────────────────────────────┐    │
│  │ 🧑 You:                                          │    │
│  │ Find recent papers about LLM agent coordination   │    │
│  │                                                    │    │
│  │ 🤖 Research Agent:                                │    │
│  │ Found 8 relevant papers. Here are the top 3:      │    │
│  │                                                    │    │
│  │ 1. "Multi-Agent Coordination via Protocol..."     │    │
│  │    Authors: Zhang et al. (2025)                   │    │
│  │    Key Finding: Protocol-agnostic bridging...     │    │
│  │                                                    │    │
│  │ 2. "Decentralized Agent Discovery with DHT..."    │    │
│  │    Authors: Li et al. (2026)                      │    │
│  │    ...                                            │    │
│  └──────────────────────────────────────────────────┘    │
│                                                          │
│  ┌────────────────────────────────────────────┐          │
│  │ Type your message...                   [➤] │          │
│  └────────────────────────────────────────────┘          │
│                                                          │
│  ⚠ Playground is rate-limited (10 calls/hour)            │
│  For production use → [Get API Key]                      │
│                                                          │
│  ▸ Advanced (Developer Mode)                             │
│    Protocol: [Auto ▾]  Session: s_abc123                 │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

**Implementation**:
- Frontend sends messages via `POST /api/v1/invoke/:agent_id` (protocol-agnostic)
- PeerClaw auto-selects the optimal protocol — users never see protocol details
- Rate-limited to prevent abuse (playground quota per user/IP)
- Developer Mode (collapsed toggle): exposes protocol selector, raw request/response, latency metrics

### 4.5 Publish Page (`/publish`)

Wizard-style Agent registration.

```
Step 1: Basic Info          Step 2: Capabilities        Step 3: Connection
┌────────────────────┐     ┌────────────────────┐     ┌────────────────────┐
│ Agent Name         │     │ Capabilities       │     │ Endpoint URL       │
│ [________________] │     │ + [Add capability] │     │ [________________] │
│                    │     │ [search] [×]       │     │                    │
│ Description        │     │ [summarize] [×]    │     │ Transport          │
│ [________________] │     │                    │     │ (●) HTTP  ( ) WS   │
│ [________________] │     │ Skills             │     │                    │
│                    │     │ + [Add skill]      │     │ Protocols          │
│ Version            │     │                    │     │ [✓] A2A            │
│ [________________] │     │ Tools              │     │ [✓] MCP            │
│                    │     │ + [Add tool]       │     │ [ ] ACP            │
│ Category           │     │                    │     │                    │
│ [Research     ▾]   │     │                    │     │ Auth Type          │
│                    │     │                    │     │ [Bearer Token ▾]   │
│         [Next →]   │     │         [Next →]   │     │      [Publish →]   │
└────────────────────┘     └────────────────────┘     └────────────────────┘
```

**On Submit**: Calls `POST /api/v1/agents` with the full Agent Card.

### 4.6 Provider Dashboard (`/dashboard`)

Analytics and management for Agent providers.

```
┌──────────────────────────────────────────────────────────┐
│  My Agents                                               │
│                                                          │
│  ── Overview ────────────────────────────────────────    │
│  Total Agents: 3    Total Calls: 5.2K    Avg Rating: 4.7│
│                                                          │
│  ── My Agents ───────────────────────────────────────    │
│  ┌────────────────────────────────────────────────────┐  │
│  │ 🟢 Research Agent    │ 2.4K calls │ ⭐ 4.9 │ [Manage]│ │
│  │ 🟢 Summarizer        │ 1.8K calls │ ⭐ 4.6 │ [Manage]│ │
│  │ 🔴 Legacy Bot        │   42 calls │ ⭐ 3.2 │ [Manage]│ │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  [+ Publish New Agent]                                   │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

### 4.7 Agent Management (`/dashboard/agents/:id`)

Per-agent management view.

```
┌──────────────────────────────────────────────────────────┐
│  Research Agent — Management                             │
│                                                          │
│  [Overview]  [Analytics]  [Settings]  [API Keys]         │
│                                                          │
│  ── Analytics (30 days) ─────────────────────────────    │
│  ┌────────────────────────────────────────────────────┐  │
│  │  📈 Call Volume                                    │  │
│  │  ████████████████████████████████                  │  │
│  │  ▁▂▃▄▅▆▇█▇▆▇████▇▆▇████████████                  │  │
│  └────────────────────────────────────────────────────┘  │
│                                                          │
│  Calls Today     Avg Latency     Success Rate    Rating  │
│     87              1.2s            99.2%         4.9    │
│                                                          │
│  ── Protocol Breakdown ──────────────────────────────    │
│  MCP: 68%    A2A: 30%    Bridge (cross-protocol): 2%    │
│                                                          │
│  ── Recent Errors ───────────────────────────────────    │
│  2026-03-07 14:23  timeout  "Tool call exceeded 30s"     │
│  2026-03-06 09:11  error    "Upstream unavailable"       │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

---

## 5. Trust & Safety Framework

Trust is the core challenge in a C2C marketplace. PeerClaw has unique advantages here through its built-in cryptographic identity and reputation system.

### 5.1 Trust Signals (Layered)

```
Layer 1 — Cryptographic Identity (Existing)
  └─ Ed25519 public key → every Agent has a verifiable identity
  └─ Identity Anchoring → Nostr-based persistent identity
  └─ Domain Verification → DNS TXT record proves domain ownership

Layer 2 — Reputation Score (Existing)
  └─ EWMA scoring (0.0 ~ 1.0) based on interaction outcomes
  └─ Automatic isolation below 0.15 threshold
  └─ Reputation gossip across the network

Layer 3 — User Reviews (New)
  └─ Star ratings (1~5) + text reviews
  └─ Only users who have interacted can review
  └─ Review authenticity verified via caller identity

Layer 4 — Platform Verification (New)
  └─ "Verified" badge for identity-anchored Agents
  └─ "Trusted" badge for reputation > 0.8 sustained over 30 days
  └─ Abuse reporting and manual review process
```

### 5.2 Safety Controls

| Concern | Mechanism |
|---------|-----------|
| Malicious Agent | Reputation scoring auto-isolates below 0.15 |
| Spam Registration | Rate limiting + require valid endpoint health check |
| Data Leakage | End-to-end encryption (XChaCha20-Poly1305) |
| Identity Spoofing | Ed25519 signature verification on every message |
| Playground Abuse | Per-user rate limits (10 calls/hour free tier) |
| Review Manipulation | Reviews tied to verified interaction records |

---

## 6. New Backend Requirements

### 6.1 New API Endpoints

```
# User Accounts
POST   /api/v1/users/register          → Create user account
POST   /api/v1/users/login             → Authenticate (JWT)
GET    /api/v1/users/me                → Current user profile

# Enhanced Agent Queries
GET    /api/v1/agents?sort=reputation   → Sort by reputation (new param)
GET    /api/v1/agents?category=research  → Filter by category (new param)
GET    /api/v1/agents?q=search+papers    → Full-text search (new param)

# Agent Invocation (Consumer-Facing)
POST   /api/v1/invoke/:agent_id         → Protocol-agnostic invocation
GET    /api/v1/invoke/:agent_id/stream   → SSE stream for async results

# Reviews
GET    /api/v1/agents/:id/reviews       → List reviews
POST   /api/v1/agents/:id/reviews       → Submit review (authenticated)

# Analytics (Provider)
GET    /api/v1/agents/:id/analytics     → Call volume, latency, errors

# API Key Management
POST   /api/v1/keys                     → Generate API key
GET    /api/v1/keys                     → List my keys
DELETE /api/v1/keys/:key_id             → Revoke key
```

### 6.2 New Data Models

```go
// User account (new table)
// Identity is key-native: PublicKey is the true identity, Email is optional convenience.
type User struct {
    ID            string    `json:"id"`
    Email         string    `json:"email,omitempty"`       // Optional: convenience login
    DisplayName   string    `json:"display_name"`
    PublicKey     string    `json:"public_key"`            // Primary identity (Ed25519)
    KeyCustodied  bool      `json:"key_custodied"`         // true if platform manages the key
    AuthProvider  string    `json:"auth_provider,omitempty"` // "email", "oauth:github", "self-hosted"
    CreatedAt     time.Time `json:"created_at"`
}

// Agent review (new table)
type Review struct {
    ID        string    `json:"id"`
    AgentID   string    `json:"agent_id"`
    UserID    string    `json:"user_id"`
    Rating    int       `json:"rating"`            // 1-5
    Comment   string    `json:"comment"`
    CreatedAt time.Time `json:"created_at"`
}

// Invocation record (new table, for analytics)
type InvocationRecord struct {
    ID          string    `json:"id"`
    AgentID     string    `json:"agent_id"`
    CallerID    string    `json:"caller_id"`        // User or Agent
    Protocol    string    `json:"protocol"`
    LatencyMs   int64     `json:"latency_ms"`
    Success     bool      `json:"success"`
    ErrorMsg    string    `json:"error_msg,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}

// Agent Card extensions (modify existing)
// Add to PeerClawExtension:
//   Category      string   `json:"category,omitempty"`
//   CallCount     int64    `json:"call_count,omitempty"`
//   AvgLatencyMs  int64    `json:"avg_latency_ms,omitempty"`
//   SuccessRate   float64  `json:"success_rate,omitempty"`
//   AvgRating     float64  `json:"avg_rating,omitempty"`
```

### 6.3 Invocation Flow

The key new capability: a consumer invokes an Agent through PeerClaw without knowing which protocol the Agent uses.

```
Consumer                    PeerClaw Server                     Agent
   │                              │                                │
   │  POST /api/v1/invoke/:id     │                                │
   │  {"input": "search papers"}  │                                │
   │─────────────────────────────►│                                │
   │                    ┌─────────┴─────────┐                      │
   │                    │ 1. Resolve protocol│                      │
   │                    │ 2. Wrap in Envelope│                      │
   │                    │ 3. Route via Bridge│                      │
   │                    └─────────┬─────────┘                      │
   │                              │  MCP/A2A/ACP call              │
   │                              │───────────────────────────────►│
   │                              │  ◄── Response                  │
   │                              │                                │
   │  ◄── {"output": "..."}      │  (record invocation stats)     │
   │◄─────────────────────────────│                                │
```

---

## 7. MVP Phasing

### Phase C1: Marketplace Browse & Profile (4 weeks)

**Goal**: Users can discover and evaluate Agents.

**Scope**:
- Landing page with featured Agents and categories
- Explore page with search, filter, sort
- Agent profile page with full detail and trust info
- Top navigation bar layout
- Mobile-responsive design
- No user accounts required (read-only marketplace)

**Backend Changes**:
- Extended `GET /api/v1/agents` with `sort`, `category`, `q` query params
- `GET /api/v1/agents/:id` already exists
- Category field added to Agent Card

**Not Included**: User accounts, playground, reviews, analytics.

### Phase C2: Playground & Invocation (4 weeks)

**Goal**: Users can try and use Agents.

**Scope**:
- Agent Playground (chat-style interface)
- `POST /api/v1/invoke/:agent_id` — protocol-agnostic invocation endpoint
- SSE streaming for async responses
- Anonymous playground with rate limiting (10 calls/hour per IP)
- Invocation record logging

**Backend Changes**:
- New invoke handler (wraps existing bridge/send logic)
- InvocationRecord table + insert on each call
- Rate limiting for playground

### Phase C3: User Accounts & Provider Console (4 weeks)

**Goal**: Providers can manage Agents, consumers can save preferences.

**Scope**:
- User registration/login (email + password or OAuth)
- JWT-based session management
- Publish wizard (guided Agent registration)
- Provider dashboard (my agents, basic stats)
- API key management (generate/revoke keys for programmatic access)
- Interaction history

**Backend Changes**:
- User table + auth endpoints
- API key table + key validation middleware
- Analytics aggregation queries

### Phase C4: Trust & Community (3 weeks)

**Goal**: Community-driven trust signals.

**Scope**:
- Reviews & ratings
- Verified / Trusted badges
- Agent categories and tagging
- Provider analytics dashboard (call volume chart, latency, errors)
- Abuse reporting

**Backend Changes**:
- Review table + endpoints
- Badge computation (periodic job)
- Report/flag system

---

## 8. Technical Architecture

### 8.1 Frontend Stack

Same as the existing Dashboard for consistency:
- React + Vite + TypeScript
- shadcn/ui + Tailwind CSS (but light theme default for consumer UX)
- Recharts for analytics charts
- React Router for client-side routing

**Key Difference from Dashboard**: The marketplace is a **separate SPA** with its own build and deployment. It could be:
- Option A: Separate domain (`marketplace.peerclaw.dev`) — deployed independently
- Option B: Same binary, different base path (`/marketplace/`) — embedded alongside dashboard

Recommendation: **Option A** for Phase C1+ to allow independent iteration. The dashboard serves operators; the marketplace serves end users. Mixing them creates UX confusion.

### 8.2 Backend Architecture

```
server/
  internal/server/
    http.go              ← Add new routes
    invoke_handler.go    ← NEW: protocol-agnostic invocation
    user_handler.go      ← NEW: user account management
    review_handler.go    ← NEW: reviews CRUD
    analytics_handler.go ← NEW: invocation analytics
  internal/registry/
    store.go             ← Extend ListFilter with sort, category, full-text
  internal/user/         ← NEW: user domain
    store.go
    service.go
```

### 8.3 Deployment Options

```
Option A — Separate Frontend (Recommended for production)
┌───────────────┐     ┌───────────────┐
│  CDN / Vercel │────►│  PeerClaw API │
│  (React SPA)  │     │  (Go Server)  │
└───────────────┘     └───────────────┘

Option B — Embedded (Simpler for self-hosted)
┌─────────────────────────────────┐
│  PeerClaw Server (Go)           │
│  ├─ /api/v1/*   → API handlers  │
│  ├─ /dashboard  → Admin SPA     │
│  └─ /           → Marketplace   │
└─────────────────────────────────┘
```

---

## 9. Success Metrics

| Metric | Phase C1 Target | Phase C4 Target |
|--------|----------------|----------------|
| Registered Agents | 20 | 200 |
| Monthly Active Users | 100 | 2,000 |
| Daily Agent Invocations | — | 5,000 |
| Avg Reputation Score | — | > 0.7 |
| Agents with Reviews | — | 50% |

---

## 10. Design Decisions

### D1: Business Model — Infrastructure Free, Value-Added Paid

**Decision**: Free infrastructure layer + paid premium tiers (GitHub model).

**Rationale**: The platform targets the widest possible audience. Early-stage priority is **network effects** — every new Agent makes the platform more valuable for everyone. Transaction fees create friction, especially for Agent-to-Agent automated invocations where per-call billing becomes complex.

```
Free Tier:     Register, discover, Playground (rate-limited), basic analytics
Pro Tier:      High-frequency invocation, detailed analytics, priority listing, SLA
Enterprise:    Private deployment, federation nodes, custom domain, compliance audit
```

Revenue comes from value-added services to providers who want more visibility and better tooling, not from taxing every interaction. This follows the same path as npm, Docker Hub, and GitHub — free infrastructure, paid professional features.

### D2: Agent Listing — Low Barrier + Layered Trust Signals

**Decision**: All registered Agents are listed immediately. Trust signals are displayed prominently but do not gate listing.

**Rationale**: High listing barriers would kill early supply — we need as many Agents as possible to build network effects. Instead of gatekeeping, we make trust transparent and let consumers decide.

```
🟢 Online + Verified + ⭐4.8+    ← Featured, top of results
🟢 Online + Unverified            ← Normal display
🟡 Offline / Unresponsive         ← Grayed out, ranked lower
🔴 Reputation < 0.15              ← Auto-isolated, hidden from marketplace
```

This mirrors the real world: anyone can open a shop, but shops with business licenses, good reviews, and steady customers get more traffic. The reputation system (already built) is the enforcement mechanism, not manual review.

### D3: Consumer Identity — Key-Native with Account Convenience Layer

**Decision**: Ed25519 key pair is the foundational identity. Email/OAuth registration is a convenience layer that auto-generates and manages keys behind the scenes.

**Rationale**: If every person will have their own Agent, the identity system must be Agent-native (key pairs), not human-native (email/password). But mass adoption requires a bridge — simple onboarding that hides cryptographic complexity.

```
Foundation:     Ed25519 key pair (Agent-native identity)
                ↑
Convenience:    Email signup / OAuth login → auto-generates & custodies key
                ↑
Advanced:       Self-hosted key / hardware key / Nostr identity (npub)
```

This ensures:
- **Ordinary users** sign up with email, never see a key — it just works
- **Developers / Agents** authenticate directly with key signatures — no "login" needed
- **Power users** can export and self-custody their keys (like moving from custodial to self-custodial wallet)
- **All API calls** are unified under key-based signature verification — humans and Agents use the same identity system

### D4: Playground Protocol Awareness — Fully Abstracted by Default

**Decision**: The Playground completely abstracts away protocol differences. An advanced developer panel is available as an opt-in toggle.

**Rationale**: The target is the widest possible audience. Users don't know and don't care about A2A vs MCP vs ACP, just as web users don't care about HTTP/2 vs HTTP/3. Protocol abstraction is PeerClaw's core value proposition — the Playground should embody it.

```
Default Mode (everyone):
  User types input → PeerClaw auto-selects optimal protocol → result displayed

Developer Mode (collapsed under "Advanced"):
  Protocol selector → raw request/response inspector → latency comparison
```

### D5: Federation Scope — Aggregate Across All Peers

**Decision**: The marketplace aggregates Agents from all federated peers by default, with source labeling and optional filtering.

**Rationale**: If the vision is "every person has an Agent," no single server can host them all. Federation is the scaling path — like email works across servers. The marketplace must reflect the full network to provide maximum value.

```
User searches "translation Agent"

Results:
  🟢 TranslateBot     ⭐4.9   local                ← source labeled
  🟢 LangBridge       ⭐4.7   peer.example.com     ← federated peer
  🟢 PolyglotAI       ⭐4.5   asia.peerclaw.dev    ← federated peer
```

Behavior:
- Default: show Agents from all connected federation peers
- Label each Agent's source node (like email's @domain)
- Allow filtering by node (for geographic or trust preferences)
- Local Agents ranked slightly higher (lower latency)
- Reputation scores are cross-federation via existing gossip mechanism
