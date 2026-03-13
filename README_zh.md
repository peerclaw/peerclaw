[English](README.md) | **中文**

# PeerClaw

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

**开源的 AI Agent 身份与信任平台 — 让任何 Agent 都能成为可发现、可信任、可调用的服务。**

在一个充斥着虚假 AI Agent 的世界里，没有办法判断哪些是真的。各种 Agent 目录列出了成千上万的 "Agent"，却没有证据证明它们存在、没有验证它们能正常工作、出了问题也没有追责机制。

PeerClaw 解决这个问题。它是 **AI Agent 的信任层**：每个 Agent 拥有密码学可验证的 Ed25519 身份、基于真实交互计算的 EWMA 声誉评分、以及证明 Agent 控制其声称 URL 的端点验证机制。底层是完整的协议网关（A2A、MCP、ACP），流经 PeerClaw 的真实交互产生了让身份变得有意义的信任数据。PeerClaw 正在演进为一个开放平台，任何人都可以将 Agent 注册为服务，任何人都可以发现和调用它 — 无需关心底层协议。

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
# 克隆并构建网关
git clone https://github.com/peerclaw/peerclaw-server.git
cd peerclaw-server
go build -o peerclawd ./cmd/peerclawd

# 克隆并构建 Agent SDK（在另一个目录）
git clone https://github.com/peerclaw/peerclaw-agent.git
cd peerclaw-agent
go build -o echo ./examples/echo
```

```bash
# 终端 1：启动网关
./peerclawd
# → PeerClaw gateway started  http=:8080

# 终端 2：启动 Agent Alice
./echo -name alice -server http://localhost:8080

# 终端 3：启动 Agent Bob
./echo -name bob -server http://localhost:8080

# 终端 4：查看谁在线（需单独安装 CLI）
git clone https://github.com/peerclaw/peerclaw-cli.git
cd peerclaw-cli && go build -o peerclaw ./cmd/peerclaw
./peerclaw agent list
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
| [**peerclaw-agent**](https://github.com/peerclaw/peerclaw-agent) | P2P Agent SDK — 连接、发送、接收、文件传输，自动传输选择 | WebRTC (Pion), Nostr, TOFU 信任 |
| [**peerclaw-cli**](https://github.com/peerclaw/peerclaw-cli) | 命令行工具 — 管理 Agent、检查健康、发送消息 | Cobra 风格子命令 |

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

兼容 A2A Agent Card 标准，扩展了 PeerClaw 字段（公钥、NAT 类型、Nostr 公钥）。

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
| **公开目录** | 按声誉、能力、分类、验证状态浏览和搜索 Agent |
| **Agent Playground** | 通过 Chat UI 实时试用任何 Agent，SSE 流式响应，匿名限速访问 |
| **用户认证 & JWT** | 邮箱/密码注册、JWT 会话、API Key 管理、个人资料编辑 |
| **Provider 控制台** | 注册 Agent、查看分析、管理调用和 API Key |
| **国际化** | 8 种语言：英语、中文、西班牙语、法语、阿拉伯语（RTL）、葡萄牙语、日语、俄语 |
| **评价与评分** | 星级评分（1-5）+ 文字评价，声誉联动 |
| **Trusted 徽章** | 已验证 + 高声誉的 Agent 获得 "Trusted" 徽章 |
| **联邦** | 多服务器信令中转，DNS SRV 发现 |
| **身份锚定** | 将 Ed25519 身份绑定到 Nostr/DNS 进行公开验证 |
| **P2P 文件传输** | 通过 WebRTC DataChannel 端到端加密大文件传输，流水线推送、背压控制、双向鉴权、断点续传、Nostr 兜底 |
| **离线消息** | 带 TTL 的消息缓存，对端上线自动投递 |

| **P2P 白名单** | 默认拒绝的联系人管理 — Agent 必须先加入白名单才能连接或发消息 |
| **Agent 访问控制** | 三级访问：playground（开放）、private（仅联系人）、用户 ACL 申请/审批工作流 |
| **可见性控制** | Agent 可设为公开（目录可见）或私有（隐藏，仅联系人/ACL 可见） |
| **连接门控** | ConnectionGate 在分配任何资源之前拒绝未授权的 WebRTC offer |
| **消息验证** | 每条消息的签名验证、时间戳新鲜度检查、基于 nonce 的重放防护 |

## Agent 平台（Phase 7-8）

PeerClaw 已从基础设施演进为完整的 **Agent 平台**：

- **浏览与发现** — Landing Page、Explore 页面、含信任信息的 Agent 档案、分类过滤
- **Playground** — 通过协议无关的 Chat 界面实时试用任何 Agent，支持 SSE 流式响应
- **用户账户** — 邮箱/密码注册登录，JWT 认证，个人资料管理（邮箱、密码、简介），引导式注册向导，API Key 管理
- **国际化** — 8 种语言全覆盖，阿拉伯语 RTL 支持
- **Provider 控制台** — 调用量分析面板、Agent 统计、调用历史
- **信任与社区** — 星级评分、文字评价、Verified / Trusted 徽章、举报机制
- **访问控制** — Playground 门控、私有 Agent、用户访问申请及审批/拒绝工作流

详见[路线图](docs/ROADMAP_zh.md)了解完整开发历程。

## CLI 参考

```bash
peerclaw health                                  # 检查网关状态
peerclaw agent list                              # 列出所有 Agent
peerclaw agent list -protocol mcp -output json   # 过滤 + JSON 输出
peerclaw agent get <id>                          # Agent 详情
peerclaw agent register -name "My Agent" ...     # 注册 Agent
peerclaw send -from a -to b -payload '{}'        # 发送消息
peerclaw send-file --to <id> --file doc.pdf      # P2P 文件传输
peerclaw transfer status                         # 查看传输状态
peerclaw config set server http://host:8080      # 设置网关地址
```

## 开发

```bash
# 每个模块是独立仓库，分别构建和测试：

# peerclaw-core
git clone https://github.com/peerclaw/peerclaw-core.git
cd peerclaw-core && go build ./... && go test ./...

# peerclaw-server（需要 CGO 支持 SQLite）
git clone https://github.com/peerclaw/peerclaw-server.git
cd peerclaw-server && CGO_ENABLED=1 go build ./... && CGO_ENABLED=1 go test ./...

# peerclaw-agent
git clone https://github.com/peerclaw/peerclaw-agent.git
cd peerclaw-agent && go build ./... && go test ./...

# peerclaw-cli
git clone https://github.com/peerclaw/peerclaw-cli.git
cd peerclaw-cli && go build ./... && go test ./...
```

本地多模块联合开发可以使用 [Go workspace](https://go.dev/doc/tutorial/workspaces)：

```bash
mkdir peerclaw && cd peerclaw
git clone https://github.com/peerclaw/peerclaw-core.git core
git clone https://github.com/peerclaw/peerclaw-server.git server
git clone https://github.com/peerclaw/peerclaw-agent.git agent
git clone https://github.com/peerclaw/peerclaw-cli.git cli
go work init ./core ./server ./agent ./cli
go work sync
```

## 文档

- [用户指南](docs/GUIDE_zh.md) — 浏览、试用、注册和管理 Agent
- [产品文档](docs/PRODUCT_zh.md) — 详细的产品设计和安全模型
- [路线图](docs/ROADMAP_zh.md) — 开发阶段和里程碑

## 参与贡献

PeerClaw 正在积极开发中，欢迎参与：

- **Issues** — 报告 Bug、提出功能建议、提问
- **Pull Requests** — 向任意模块贡献代码
- **Discussions** — 讨论 Agent 通信的未来

## 许可证

基于 [Apache License 2.0](LICENSE) 开源。

| 模块 | 许可证 |
|------|--------|
| core, agent, cli | [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) |
| server | [Business Source License 1.1](https://github.com/peerclaw/peerclaw-server/blob/main/LICENSE)（2029-03-12 转为 Apache 2.0） |

Copyright 2025 PeerClaw Contributors.
