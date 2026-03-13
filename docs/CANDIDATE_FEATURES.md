# Candidate Features

Features and ideas that have been considered or prototyped but are **not currently implemented**. These are maintained here for future reference — inclusion in this list does not imply a commitment to implement.

Each entry records the idea, rationale for deferral, and conditions under which it might be reconsidered.

---

## Kademlia DHT Discovery

**What:** Fully decentralized agent discovery using a Kademlia distributed hash table, with Nostr relays as the transport layer for DHT messages (kind 20005).

**Was prototyped?** Yes — a full implementation existed (k-bucket routing, Ping/Store/FindNode/FindValue RPC, DHT coordinator, CompositeDiscovery with server-first + DHT fallback). It was removed during a project slimdown.

**Why deferred:**
- AI Agent ecosystem is early stage; centralized Directory provides better UX and discoverability
- Nostr relay already provides decentralized signaling/transport without DHT complexity
- Multi-server Federation achieves fault tolerance with far less complexity
- Security audit found 6 issues (Eclipse attacks, Sybil attacks, unsigned messages, SHA-1 NodeID) — maintenance cost not justified at current scale
- Cold-start problem: DHT needs critical mass of nodes to be useful
- Multi-hop discovery latency (seconds) vs single REST call (milliseconds)

**Reconsider when:**
- Network grows to 10,000+ active agents where centralized registry becomes a bottleneck
- Censorship-resistant discovery becomes a hard requirement
- Nostr relay infrastructure proves unreliable for signaling

**References:**
- Security findings: `internal-docs/AUDIT_REPORT.md` (H-12, M-20, M-21, M-22, L-04, L-08)
- Original design: `internal-docs/architecture/user-capabilities-analysis.md`

---

## Identity Anchoring (DNS / Nostr)

**What:** Bind an agent's Ed25519 identity to external systems (DNS TXT records, Nostr profile metadata) so third parties can independently verify agent ownership without trusting the PeerClaw server.

**Was prototyped?** Partially — the `identity_anchor` field exists in the agent card spec, and Nostr pubkey binding is functional. DNS anchoring was designed but not fully implemented.

**Why deferred:**
- Nostr pubkey binding already works for the decentralized use case
- DNS TXT verification adds operational complexity for agent operators
- Current Ed25519 signature verification is sufficient for trust establishment

**Reconsider when:**
- Enterprise customers need domain-based identity proof (similar to DKIM for email)
- Cross-platform agent identity portability becomes important

---

## Plugin System

**What:** Allow third-party developers to extend PeerClaw server with custom bridges, reputation algorithms, or discovery backends via a Go plugin interface.

**Was prototyped?** No — only discussed.

**Why deferred:**
- Go plugin system has significant limitations (same Go version, same OS, no Windows)
- Current bridge/protocol system is extensible through code contributions
- gRPC-based plugin model (like HashiCorp's go-plugin) is more viable but adds complexity

**Reconsider when:**
- Community contributors want to add protocol bridges without forking
- PeerClaw is used as embedded infrastructure in diverse environments

---

## Agent Marketplace / Monetization

**What:** Enable agent providers to charge for API access, with PeerClaw handling billing, usage metering, and payment settlement.

**Was prototyped?** Designed only — see `internal-docs/marketplace/`.

**Why deferred:**
- Premature to add payment infrastructure before establishing user base
- Regulatory complexity (payment processing, tax compliance)
- Current focus is on trust/discovery infrastructure, not commerce

**Reconsider when:**
- Significant number of providers request monetization features
- Clear revenue model validated through user research

---

## Multi-Signature Identity Recovery

**What:** Allow agents to pre-register a set of recovery keys (threshold-of-n) so that if the primary Ed25519 private key is lost or compromised, the identity can be rotated to a new keypair with approval from enough recovery signers.

**Was prototyped?** Partially — a `RecoveryManager` with threshold validation existed in the agent SDK. It was never integrated into the agent lifecycle, server API, or CLI.

**Why deferred:**
- Current keypair model is simple and sufficient — most agents are ephemeral or operator-managed
- Recovery adds significant UX complexity (recovery key distribution, secure storage)
- No user requests for key recovery; key rotation via re-registration is the current fallback

**Reconsider when:**
- Long-lived agent identities become common and key loss is a real operational risk
- Enterprise customers require key rotation without identity change

---

*Last updated: 2026-03-13*
