[English](PRODUCT.md) | **中文**

# PeerClaw 产品文档

## 产品愿景

构建一个去中心化优先的 AI Agent 通信基础设施，让任何 Agent 无论使用什么协议、部署在哪里，都能安全、高效地互相发现和通信。

## 目标用户

### AI Agent 开发者

- 需要让自己的 Agent 与其他 Agent 通信
- 不想被绑定在某一个协议生态（A2A / MCP / ACP）
- 希望 P2P 直连而非所有流量都经过中心服务器
- 需要开箱即用的安全方案（身份、签名、信任管理）

### 平台集成商

- 运营 Agent 平台，需要统一的注册和发现机制
- 需要桥接不同协议生态的 Agent
- 需要可观测性和审计能力
- 需要水平扩展和高可用

## 核心场景

### 场景 1: Agent 发现与连接

```
Alice (搜索 Agent)                    PeerClaw Server                    Bob (数据 Agent)
       │                                    │                                   │
       │  POST /api/v1/agents (注册)         │                                   │
       │──────────────────────────────────►  │                                   │
       │                                    │  ◄──── POST /api/v1/agents (注册)  │
       │                                    │                                   │
       │  POST /api/v1/discover             │                                   │
       │  {"capabilities": ["data"]}        │                                   │
       │──────────────────────────────────►  │                                   │
       │  ◄── [{id: "bob", pubkey: "..."}]  │                                   │
       │                                    │                                   │
```

Agent 通过 REST API 注册自身能力和公钥，其他 Agent 可按能力或协议搜索和发现。

### 场景 2: P2P 安全通信

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
  │  ◄═══════ WebRTC DataChannel (P2P 直连) ═══════►                   │
  │                                  │                                  │
  │  Envelope {encrypted_payload, signature}                           │
  │════════════════════════════════════════════════════════════════════►│
  │  ◄═════════════ Envelope {encrypted_payload, signature} ══════════│
```

信令服务器仅用于 WebRTC 握手，实际数据通过 P2P DataChannel 传输。每条消息附带 Ed25519 签名。

### 场景 3: 跨协议互操作

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

Bridge Manager 自动识别源协议和目标协议，通过 Envelope 中间格式完成无缝转换。

### 场景 4: 去中心化兜底

```
Alice                         Nostr Relay                          Bob
  │                               │                                 │
  │  (WebRTC 连接失败)             │                                 │
  │                               │                                 │
  │  NIP-44 加密 Envelope          │                                 │
  │──────────────────────────────►│────────────────────────────────►│
  │                               │  ◄── NIP-44 加密 Envelope       │
  │◄──────────────────────────────│◄────────────────────────────────│
```

当 WebRTC P2P 连接无法建立时（严格 NAT、防火墙），自动回退到 Nostr relay 传输。

## 系统架构

```
┌────────────────────────────────────────────────────────────────┐
│                      peerclaw-server                           │
│                                                                │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │           HTTP 入口 + 协议端点 + 联邦端点                  │   │
│  │  /api/v1/*  /a2a  /mcp  /acp/*  /api/v1/federation/*   │   │
│  └──────┬──────────────┬──────────────┬───────────────────┘   │
│         │              │              │                        │
│  ┌──────▼──────┐ ┌─────▼─────┐ ┌─────▼──────┐                │
│  │  Registry   │ │  Router   │ │  Signaling │                │
│  │  Service    │ │  Engine   │ │  Hub       │                │
│  │             │ │           │ │            │                │
│  │ - 注册/注销  │ │ - 路由表   │ │ - WS 连接  │                │
│  │ - 发现/查询  │ │ - 路由解析 │ │ - 消息转发  │                │
│  │ - 心跳管理  │ │ - 能力匹配 │ │ - Ping/Pong│                │
│  │ - 联邦发现  │ │           │ │- bridge_msg│                │
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
│  │ - 对端连接  │ │  └────────────┘  └──────────────────┘  │   │
│  │ - 信令转发  │ └────────────────────────────────────────┘   │
│  │ - DNS SRV   │                                              │
│  └─────────────┘                                              │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│                      peerclaw-agent (SDK)                      │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐ │
│  │   Agent API  │  │  Discovery   │  │     Signaling        │ │
│  │              │  │  (接口)       │  │     (接口)           │ │
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
│  │ - 路由表     │  │  Domain Verify + Recovery             │   │
│  │ - KV Store   │  └──────────────────────────────────────┘   │
│  │ - Bootstrap  │                                              │
│  └──────────────┘                                              │
│  ┌────────────────────────────────────────────────────────┐   │
│  │                    Transport                           │   │
│  │  WebRTC DataChannel ◄──► Nostr Relay ◄──► MessageCache │   │
│  └────────────────────────────────────────────────────────┘   │
└────────────────────────────────────────────────────────────────┘

┌────────────────────────────────────────────────────────────────┐
│                      peerclaw-core (类型库)                     │
│                                                                │
│  identity    envelope    agentcard    protocol    signaling    │
│  (Ed25519)   (消息信封)   (Agent 名片)  (协议常量)   (信令类型)    │
└────────────────────────────────────────────────────────────────┘
```

## 通信流程

完整的 Agent 间通信经历以下阶段：

### 1. 注册

Agent 启动时通过 REST API 向 Server 注册，提交 Agent Card（名称、公钥、能力、协议）。Server 存储到 SQLite 并更新路由表。

### 2. 发现

Agent A 通过 `/api/v1/discover` 按能力搜索目标 Agent。Server 返回匹配的 Agent Card 列表，包含公钥和连接信息。

### 3. 信令

Agent A 通过 WebSocket 连接 Server 的信令 Hub，发送 WebRTC offer SDP 给 Agent B。Server 转发信令消息。Agent B 回复 answer SDP，双方交换 ICE candidate。

### 4. P2P 连接

WebRTC ICE 协商完成后，Agent A 和 B 建立 DataChannel 直连。如果 ICE 失败（严格 NAT），回退到 Nostr relay。

### 5. 消息交换

消息封装在 Envelope 中，包含源/目标、协议类型、Payload、Ed25519 签名。接收方验证签名后处理消息。

## 安全模型

### 第一层：连接级 — TOFU (Trust-On-First-Use)

| 属性 | 说明 |
|------|------|
| 机制 | 首次连接记录对端公钥指纹 |
| 存储 | 本地 Trust Store 文件 |
| 验证 | 后续连接检查公钥是否与记录一致 |
| 威胁缓解 | 中间人攻击（首次之后） |
| 类比 | SSH known_hosts |

### 第二层：消息级 — Ed25519 签名

| 属性 | 说明 |
|------|------|
| 算法 | Ed25519（RFC 8032） |
| 签名对象 | Envelope Payload |
| 验证方 | 接收 Agent |
| 威胁缓解 | 消息篡改、身份伪造 |
| 性能 | ~76,000 签名/秒，~200,000 验证/秒 |

### 第三层：传输级 — 端到端加密

| 属性 | 说明 |
|------|------|
| 密钥交换 | X25519 ECDH（从 Ed25519 seed 派生） |
| 密钥派生 | HKDF-SHA256 |
| 对称加密 | XChaCha20-Poly1305（24 字节随机 nonce） |
| 会话建立 | 信令握手阶段交换 X25519 公钥 |
| Nostr 适配 | NIP-44 格式封装（secp256k1 会话密钥） |
| 威胁缓解 | 窃听、中间人、消息泄露 |

### 第四层：执行级 — 沙箱

| 属性 | 说明 |
|------|------|
| 机制 | 对外部请求实施权限约束 |
| 控制 | 资源限制、操作白名单 |
| 威胁缓解 | 恶意操作、资源耗尽 |

## 协议兼容矩阵

| 特性 | A2A | ACP | MCP |
|------|-----|-----|-----|
| 消息路由 | ✅ | ✅ | ✅ |
| 能力发现 | ✅ | ✅ | ✅ |
| Task 模型 | ✅ | — | — |
| Tool 调用 | — | — | ✅ |
| Artifact | ✅ | — | — |
| Resource | — | — | ✅ |
| Prompt | — | — | ✅ |
| Run 管理 | — | ✅ | — |
| 协议协商 | ✅ | ✅ | ✅ |
| 跨协议翻译 | ✅ | ✅ | ✅ |
| 双向流 | ✅ | ✅ | 🔲 Phase 4 |

## 部署架构

### 单节点

```
[ Agent A ] ──► [ PeerClaw Server (SQLite) ] ◄── [ Agent B ]
                         │
                   [ Agent A ] ◄═══ P2P ═══► [ Agent B ]
```

### 多节点（Phase 4 已实现）

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

### 去中心化（Phase 5 已实现）

```
[ Agent A ] ◄═══ P2P (WebRTC) ═══► [ Agent B ]
     │                                   │
     └──── DHT Discovery (Kademlia) ─────┘
     │                                   │
     └──── Nostr Signaling (kind 20006) ─┘
     │                                   │
     └──── Nostr Relay (kind 20004) ─────┘
```

### 联邦模式（Phase 5 已实现）

```
[ Agent A ] ──► [ Server 1 ] ◄══ Federation ══► [ Server 2 ] ◄── [ Agent B ]
                     │         (HTTP + Auth)          │
              [ DNS SRV 发现 ]                 [ DNS SRV 发现 ]
```

### 无 Server 模式（Phase 5 已实现）

```
[ Agent A ]                                    [ Agent B ]
     │  1. DHT Bootstrap (Nostr kind 20005)         │
     │──────────────►  Nostr Relays  ◄──────────────│
     │  2. Agent Card 存入 DHT                       │
     │  3. Nostr 信令 (kind 20006)                   │
     │──────────────►  Nostr Relays  ◄──────────────│
     │  4. WebRTC P2P 直连                           │
     │◄════════════ DataChannel ═══════════════════►│
     │  5. 离线消息缓存 (MessageCache)                │
```

## 信誉模型（Phase 5）

| 属性 | 说明 |
|------|------|
| 评分算法 | 指数加权移动平均（EWMA，alpha=0.1） |
| 评分范围 | 0.0 (恶意) ~ 1.0 (可信) |
| 恶意阈值 | < 0.15 自动隔离 |
| 行为类型 | success (+1.0), timeout (-0.3), error (-0.2), invalid_signature (-0.8), spam (-0.5), protocol_violation (-0.7) |
| 存储 | 本地 JSON 文件持久化 |
| Gossip | 可选 Nostr kind 30078，第二手权重 0.3x，仅接受 TrustVerified+ |

## 身份锚定（Phase 5）

| 属性 | 说明 |
|------|------|
| 接口 | IdentityAnchor (Publish/Verify/Resolve/RecoveryKeys) |
| 首选实现 | Nostr kind 10078 replaceable event |
| 密钥绑定 | 双向：Ed25519 签 Nostr key + Nostr key 签 Ed25519 key |
| 域名验证 | DNS TXT 记录 peerclaw-verify=<fingerprint> |
| 身份恢复 | threshold-of-n 多签 recovery keys |
