[English](README.md) | **中文**

# PeerClaw

**开源的 AI Agent 身份与信任平台 — 可验证身份、声誉评分、跨协议通信。**

在一个充斥着虚假 AI Agent 的世界里，没有办法判断哪些是真的。各种 Agent 市场列出了成千上万的 "Agent"，却没有证据证明它们存在、没有验证它们能正常工作、出了问题也没有追责机制。

PeerClaw 解决这个问题。它是 **AI Agent 的信任层**：每个 Agent 拥有密码学可验证的 Ed25519 身份、基于真实交互计算的 EWMA 声誉评分、以及证明 Agent 控制其声称 URL 的端点验证机制。底层是完整的协议网关（A2A、MCP、ACP），流经 PeerClaw 的真实交互产生了让身份变得有意义的信任数据。

## PeerClaw 做什么

```
┌─────────────────────────────────────────────────────────────┐
│ 你的 MCP Agent          PeerClaw 网关          A2A Agent    │
│                                                             │
│  "我需要一个能搜索   →  注册中心：找到       → "我能搜索"    │
│   的 Agent"              3 个匹配                           │
│                                                             │
│  发送 MCP 请求      →  桥接器：自动翻译     → 收到 A2A      │
│                          MCP → A2A              消息        │
│                                                             │
│  收到 MCP 响应      ←  桥接器：自动翻译     ← 发送 A2A      │
│                          A2A → MCP              响应        │
└─────────────────────────────────────────────────────────────┘
```

**简单来说：**

1. **注册** — Agent 告诉 PeerClaw 自己能做什么（能力、协议、端点）
2. **验证** — Challenge-Response 端点验证证明 Agent 控制其 URL
3. **积累信任** — 每次交互（心跳、桥接消息、验证）都会计入 EWMA 声誉评分
4. **发现** — 任何人都可以浏览公开的 Agent 目录，按声誉、能力、验证状态过滤
5. **桥接** — 使用不同协议（A2A、MCP、ACP）的 Agent 通过自动翻译无缝通信
6. **信任** — 每个 Agent 拥有 Ed25519 密码学身份，消息签名加密，无法冒充、无法篡改

## 快速开始

5 分钟内让两个 Agent 互相通信：

```bash
git clone https://github.com/peerclaw/peerclaw.git
cd peerclaw

# 克隆子项目
git clone https://github.com/peerclaw/peerclaw-core.git core
git clone https://github.com/peerclaw/peerclaw-server.git server
git clone https://github.com/peerclaw/peerclaw-agent.git agent

# 构建
cd server && go build -o peerclawd ./cmd/peerclawd && cd ..
cd agent && go build -o echo ./examples/echo && cd ..
cd cli && go build -o peerclaw ./cmd/peerclaw && cd ..
```

```bash
# 终端 1：启动网关
./server/peerclawd
# → PeerClaw gateway started  http=:8080

# 终端 2：启动 Agent Alice
./agent/echo -name alice -server http://localhost:8080

# 终端 3：启动 Agent Bob
./agent/echo -name bob -server http://localhost:8080

# 终端 4：查看谁在线
./cli/peerclaw agent list
```

Alice 和 Bob 会自动注册、发现对方、并建立加密的 P2P 连接。

## 架构

PeerClaw 由四个模块组成：

```
┌──────────────────────────────────────────────────────────────────┐
│                     peerclaw-server（网关）                       │
│                                                                  │
│  ┌────────────┐   ┌───────────────┐   ┌───────────────────────┐ │
│  │  注册中心   │   │    信令中心    │   │      协议桥接器        │ │
│  │  按能力     │   │   WebSocket   │   │                       │ │
│  │  发现 Agent │   │   中转 WebRTC │   │  ┌─────┬─────┬─────┐ │ │
│  │            │   │   信令        │   │  │ A2A │ MCP │ ACP │ │ │
│  └────────────┘   └───────────────┘   │  └─────┴─────┴─────┘ │ │
│                                       └───────────────────────┘ │
│  ┌────────────┐   ┌───────────────┐   ┌───────────────────────┐ │
│  │  认证鉴权   │   │    限速保护    │   │      可观测性         │ │
│  │  Ed25519 + │   │   Per-IP      │   │   OpenTelemetry       │ │
│  │  API Key   │   │   令牌桶      │   │   链路追踪 + 指标     │ │
│  └────────────┘   └───────────────┘   └───────────────────────┘ │
│                                                                  │
│  存储：SQLite（默认）或 PostgreSQL                                │
│  扩展：Redis Pub/Sub 多节点信令                                   │
└──────────────────────────────────────────────────────────────────┘
        │ REST API              │ WebSocket             │ 协议端点
        │ 注册/发现              │ 信令                  │ A2A/MCP/ACP
        ▼                       ▼                       ▼
┌──────────────────┐    ┌──────────────────┐    ┌──────────────┐
│  peerclaw-agent  │    │  peerclaw-agent  │    │   外部        │
│  (Go SDK)        │◄══►│  (Go SDK)        │    │  A2A/MCP/ACP │
│                  │P2P │                  │    │  Agent       │
│  WebRTC 首选     │    │  WebRTC 首选     │    │              │
│  Nostr 兜底      │    │  Nostr 兜底      │    │              │
└──────────────────┘    └──────────────────┘    └──────────────┘
        │                       │
        └───── peerclaw-core ───┘
              （共享类型：身份、信封、协议）
```

### Agent 如何通信

```
Alice                     网关                      Bob
  │                        │                          │
  ├─ POST /agents ────────►│  注册 Alice               │
  │                        │◄──────── POST /agents ───┤  注册 Bob
  │                        │                          │
  ├─ POST /discover ──────►│  "谁能搜索？"              │
  │◄── [{Bob, caps:search}]│                          │
  │                        │                          │
  ├─ WS: offer + X25519 ─►│──── 中转 ───────────────►│
  │◄─── WS: answer + X25519│◄─── 中转 ───────────────┤
  │                        │                          │
  │◄═══════════ WebRTC P2P（加密直连）═══════════════►│
  │        Ed25519 签名 + XChaCha20 加密              │
```

## 项目结构

| 模块 | 做什么 | 关键技术 |
|------|--------|---------|
| [**peerclaw-core**](https://github.com/peerclaw/peerclaw-core) | 共享类型库 — 身份、信封、Agent Card、协议常量 | Ed25519, X25519, 零外部依赖 |
| [**peerclaw-server**](https://github.com/peerclaw/peerclaw-server) | 网关 — 注册、发现、信令中转、协议桥接 | SQLite/PostgreSQL, WebSocket, OTel |
| [**peerclaw-agent**](https://github.com/peerclaw/peerclaw-agent) | P2P Agent SDK — 连接、发送、接收，自动传输选择 | WebRTC (Pion), Nostr, TOFU 信任 |
| **cli/** | 命令行工具 — 管理 Agent、检查健康、发送消息 | Cobra 风格子命令 |

## 核心概念

### Agent Card（Agent 名片）

每个 Agent 发布一张 Agent Card — 一份机器可读的能力描述：

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

兼容 A2A Agent Card 标准，扩展了 PeerClaw 字段（公钥、NAT 类型、DHT 节点 ID）。

### 协议桥接

PeerClaw 通过统一的 **Envelope（信封）** 格式在协议间翻译：

```
A2A Agent ──► A2A 适配器 ──► Envelope ──► MCP 适配器 ──► MCP Agent
                                │
                           统一格式：
                           source, destination,
                           protocol, payload,
                           signature, trace_id
```

| 协议 | 用途 | PeerClaw 支持 |
|------|------|-------------|
| **A2A**（Google） | 基于任务的 Agent 协作 | 完整：任务、制品、流式 |
| **MCP**（Anthropic） | 工具/资源访问 | 完整：工具、资源、提示词 |
| **ACP**（IBM） | 企业级 Agent 运行 | 完整：运行、会话、清单 |

### 密码学身份

每个 Agent 拥有一个 Ed25519 密钥对。公钥**就是**身份。

- **注册**：Agent 通过签名证明身份所有权
- **消息**：每条 Envelope 都被签名 — 接收方验证来源
- **加密**：从 Ed25519 派生 X25519 密钥，XChaCha20-Poly1305 加密载荷
- **信任**：TOFU（首次使用信任）模型，5 个等级：未知 → TOFU → 已验证 → 已固定 → 已封禁

### 传输层降级

Agent SDK 自动选择最佳传输方式：

```
1. WebRTC DataChannel（首选 — 低延迟，P2P 直连）
       │ 失败（严格 NAT）？
       ▼
2. Nostr 中继（兜底 — NIP-44 加密，多中继）
       │ WebRTC 恢复？
       ▼
3. 自动升级回 WebRTC
```

## 高级功能

这些功能已实现，但基础使用不需要：

| 功能 | 说明 |
|------|------|
| **声誉引擎** | 服务端 EWMA 声誉评分，基于真实交互事件 |
| **端点验证** | Challenge-Response 证明 Agent 控制其声称的 URL |
| **公开目录** | 按声誉、能力、验证状态浏览和搜索 Agent |
| **DHT 发现** | 通过 Kademlia DHT 无服务器发现 Agent（Nostr 传输） |
| **联邦** | 多服务器信令中转，DNS SRV 发现 |
| **身份锚定** | 将 Ed25519 身份绑定到 Nostr/DNS 进行公开验证 |
| **离线消息** | 带 TTL 的消息缓存，对端上线自动投递 |
| **无服务器模式** | 完全 P2P，无需任何中心服务器 |

## CLI 参考

```bash
peerclaw health                                  # 检查网关状态
peerclaw agent list                              # 列出所有 Agent
peerclaw agent list -protocol mcp -output json   # 过滤 + JSON 输出
peerclaw agent get <id>                          # Agent 详情
peerclaw agent register -name "My Agent" ...     # 注册 Agent
peerclaw send -from a -to b -payload '{}'        # 发送消息
peerclaw config set server http://host:8080      # 设置网关地址
```

## 开发

```bash
# 项目使用 Go workspace（go.work）管理多模块
go work sync

# 构建全部
cd core && go build ./... && cd ..
cd agent && go build ./... && cd ..
cd server && go build ./... && cd ..
cd cli && go build ./... && cd ..

# 测试全部
cd core && go test ./... && cd ..
cd agent && go test ./... && cd ..
cd server && CGO_ENABLED=1 go test ./... && cd ..
cd cli && go test ./... && cd ..
```

## 文档

- [产品文档](docs/PRODUCT_zh.md) — 详细的产品设计和安全模型
- [路线图](docs/ROADMAP_zh.md) — 开发阶段和里程碑

## 参与贡献

PeerClaw 正在积极开发中，欢迎参与：

- **Issues** — 报告 Bug、提出功能建议、提问
- **Pull Requests** — 向任意模块贡献代码
- **Discussions** — 讨论 Agent 通信的未来

## 许可证

MIT
