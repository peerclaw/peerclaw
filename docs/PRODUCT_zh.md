[English](PRODUCT.md) | **中文**

# PeerClaw 产品文档

## 产品愿景

构建 AI Agent 的信任基础设施，并将其演进为一个开放的 Marketplace，让任何 Agent 都能成为可发现、可信任、可调用的服务。

**第一层 — 基础设施（已完成）：** 去中心化优先的通信、可验证身份、声誉评分、跨协议桥接（A2A / MCP / ACP）。任何 Agent 无论使用什么协议、部署在哪里，都能安全、高效地互相发现和通信。

**第二层 — Marketplace（已完成）：** 在基础设施之上构建 C2C Agent Marketplace，任何人都可以将 Agent 发布为服务，任何人（人类或 Agent）都可以发现、评估和调用它 — 无需关心底层协议。

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

### Agent 服务消费者

- 需要发现和使用 Agent 服务，无需了解底层协议
- 希望在使用前评估 Agent 的可信度（声誉、评价、已验证身份）
- 需要简便的方式试用 Agent（Playground）并通过编程方式调用
- 可以是浏览 Marketplace 的人类，也可以是将子任务委托给专业 Agent 的 Agent

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

### 场景 5: Agent Marketplace

```
Provider                       PeerClaw Marketplace                    Consumer
    │                                   │                                  │
    │  注册 Agent                        │                                  │
    │  （能力、技能、端点）                │                                  │
    │──────────────────────────────────►│                                  │
    │                                   │                                  │
    │                                   │  ◄── 浏览 / 搜索 Agent           │
    │                                   │  ◄── 评估信任与声誉               │
    │                                   │──────────────────────────────────│
    │                                   │      Agent 档案 + 评价           │
    │                                   │                                  │
    │                                   │  ◄── 在 Playground 试用          │
    │                                   │──────────────────────────────────│
    │  ◄── 协议无关调用                  │      （限速试用）                 │
    │◄──────────────────────────────────│                                  │
    │  响应 ──────────────────────────►│──────────────────────────────────►│
    │                                   │                                  │
```

Provider 只需注册一次 Agent。Consumer 通过 Marketplace 发现它，借助 PeerClaw 内置的声誉和验证系统评估信任，在 Playground 试用，然后调用 — 全程无需关心 Agent 使用什么协议。PeerClaw 自动选择最优协议路径（A2A、MCP 或 ACP）。

## 系统架构

```
┌────────────────────────────────────────────────────────────────┐
│                   Marketplace 层（已完成）                        │
│                                                                │
│  Landing / Explore / Agent Profile / Playground                │
│  用户账户 / 评价评分 / Provider 分析面板                          │
└────────────────────────┬───────────────────────────────────────┘
                         │ 构建于
┌────────────────────────▼───────────────────────────────────────┐
│                    基础设施层（已完成）                           │
│  注册中心 / 信令 / 桥接 / 声誉 / 身份 / DHT                     │
└────────────────────────────────────────────────────────────────┘
```

### 基础设施详情

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
| 签名对象 | 完整 Envelope（头部 + 载荷；加密时签名覆盖密文） |
| 验证方 | 接收 Agent |
| 威胁缓解 | 消息篡改、身份伪造 |
| 性能 | ~76,000 签名/秒，~200,000 验证/秒 |

### 第三层：传输级 — 端到端加密

| 属性 | 说明 |
|------|------|
| 密钥交换 | X25519 ECDH（从 Ed25519 seed 派生） |
| 密钥派生 | HKDF-SHA256 |
| 对称加密 | XChaCha20-Poly1305（24 字节随机 nonce） |
| 顺序 | 先加密后签名（encrypt-then-sign）— 签名覆盖密文，接收端可在解密前验证发送者身份 |
| 会话建立 | 信令握手阶段交换 X25519 公钥 |
| Nostr 适配 | NIP-44 格式封装（secp256k1 会话密钥） |
| 威胁缓解 | 窃听、中间人、消息泄露、解密预言攻击 |

### 第四层：执行级 — 沙箱

| 属性 | 说明 |
|------|------|
| 机制 | 对外部请求实施权限约束 |
| 控制 | 资源限制、操作白名单 |
| 威胁缓解 | 恶意操作、资源耗尽 |

### 第五层：P2P 通信 — 白名单 + 消息验证（Phase 8）

| 属性 | 说明 |
|------|------|
| 默认策略 | 默认拒绝 — Agent 必须先加入白名单才能通信 |
| Agent 端白名单 | TrustStore 检查入站/出站消息和连接 |
| Server 端白名单 | 信令 Hub 上的 ContactsChecker 拦截非联系人的 offer/answer/ICE |
| 连接门控 | ConnectionGate 回调在分配任何 WebRTC 资源之前拒绝非白名单 offer |
| 消息验证 | 签名验证、时间戳新鲜度（±2 分钟）、基于 nonce 的重放防护、载荷大小限制（1MB） |
| 反重放 | 每条消息 UUID nonce，服务端 nonce 缓存 + 5 分钟清理 |
| 威胁缓解 | Prompt 注入、重放攻击、未授权连接的资源耗尽、信令泛洪 DDoS |
| 架构 | 纵深防御：Agent TrustStore（主防线）+ Server contacts 服务（辅助防线） |

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

## 产品演进

PeerClaw 经历了三个战略阶段的演进：

```
Phase 1-4                    Phase 5-6                       Phase 7+
通信基础设施                   身份与信任平台                    Agent Marketplace
                                                             (AaaS)

┌──────────────┐            ┌──────────────┐               ┌──────────────┐
│ 注册中心      │            │ 声誉系统      │               │ 浏览与发现    │
│ 信令          │──────────►│ 端点验证      │──────────────►│ Playground   │
│ 协议桥接      │            │ 公开目录      │               │ 用户账户      │
│ P2P / DHT    │            │ 身份锚定      │               │ 评价系统      │
│ 联邦          │            │              │               │ 分析面板      │
└──────────────┘            └──────────────┘               └──────────────┘
```

- **基础设施**（Phase 1-4）：核心协议网关 — 注册中心、信令、桥接、传输、生产就绪。
- **身份与信任平台**（Phase 5-6）：去中心化身份、声誉评分、端点验证、公开目录。真实交互产生的信任数据是 PeerClaw 的核心差异化。
- **Agent Marketplace**（Phase 7+）：C2C Marketplace，任何人都可以将 Agent 发布为服务，任何人都可以发现、评估、试用和调用它 — 无需关心底层协议。
- **P2P 通信安全加固**（Phase 8）：默认拒绝白名单强制执行、消息验证管线（签名、重放、时间戳）、连接门控 — Agent 端和 Server 端纵深防御。详见[路线图](ROADMAP_zh.md)。
