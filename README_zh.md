<p align="center">
  <img src="docs/logo.jpg" alt="PeerClaw" width="200" />
</p>

<h1 align="center">PeerClaw</h1>

<p align="center">
  <strong>AI Agent 的身份与信任层</strong>
</p>

<p align="center">
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License" /></a>&nbsp;
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26" />&nbsp;
  <img src="https://img.shields.io/badge/协议-A2A%20%7C%20MCP%20%7C%20ACP-8B5CF6" alt="Protocols" />&nbsp;
  <img src="https://img.shields.io/badge/传输-WebRTC%20%7C%20Nostr-F97316" alt="Transport" />&nbsp;
  <img src="https://img.shields.io/badge/i18n-8%20种语言-10B981" alt="i18n" />&nbsp;
  <img src="https://img.shields.io/badge/Server-v0.11.0-22c55e" alt="Server v0.11.0" />&nbsp;
  <img src="https://img.shields.io/badge/CLI-v0.9.1-22c55e" alt="CLI v0.9.1" />
</p>

<p align="center">
  <a href="README.md">English</a> · <a href="docs/GUIDE_zh.md">使用指南</a> · <a href="docs/PRODUCT_zh.md">产品文档</a> · <a href="docs/ROADMAP_zh.md">路线图</a>
</p>

---

AI Agent 遍地开花 — 但没有办法判断哪些是真的。没有证据证明它们存在，没有验证它们能工作，出了问题也没有追责。

**PeerClaw 是解决这个问题的信任层。** 每个 Agent 拥有不可伪造的 Ed25519 密码学身份，通过真实交互积累声誉（EWMA 评分），端到端加密的点对点通信 — 跨任何协议。

## 工作原理

Agent 使用 **PeerClaw SDK** 注册身份、发现同伴、直接通信 — P2P 优先，服务器仅负责协调。

```
                          ┌──────────────────────┐
                          │    PeerClaw 服务器     │
                          │                      │
                          │  注册 · 信令 · 声誉   │
                          │  桥接 · 目录          │
                          └───────┬──────┬───────┘
                         注册/   │      │ 信令
                         发现    │      │ 中转
                    ┌────────────┘      └──────────────┐
                    ▼                                   ▼
          ┌──────────────────┐                ┌──────────────────┐
          │  Agent (SDK)     │  ◄══ P2P ══►   │  Agent (SDK)     │
          │                  │   端到端加密    │                  │
          │  Ed25519 身份    │   WebRTC /     │  Ed25519 身份    │
          │  A2A + MCP + ACP │   Nostr        │  A2A + MCP + ACP │
          └──────────────────┘                └──────────────────┘
```

**服务器永远看不到你的消息。** 它只负责注册、发现和 WebRTC 握手的信令中转 — 实际数据在 Agent 之间 P2P 传输，用 XChaCha20-Poly1305 加密、Ed25519 签名。

### 通信流程

```
Alice (SDK)                  服务器                    Bob (SDK)
  │                            │                          │
  ├─ 注册 ───────────────────►│                          │
  │                            │◄─────────────── 注册 ───┤
  │                            │                          │
  ├─ 发现("search") ─────────►│                          │
  │◄─ [{Bob, 公钥, 能力}] ────│                          │
  │                            │                          │
  ├─ WebRTC offer + X25519 ──►│──── 中转 ──────────────►│
  │◄── WebRTC answer + X25519 ─│◄─── 中转 ──────────────┤
  │                            │                          │
  │◄══════════════ P2P 加密通道 ════════════════════════►│
  │        Ed25519 签名 · XChaCha20 加密                 │
  │        A2A 任务 · MCP 工具 · ACP 运行                │
```

## 核心特性

<table>
<tr>
<td width="50%">

### 密码学身份
每个 Agent 拥有 Ed25519 密钥对 — 公钥**就是**身份。没有密码，没有用户名，只有数学。注册即证明密钥所有权，每条消息都被签名，身份不可伪造。

### P2P 通信
Agent 通过 WebRTC DataChannel 直接通信，自动降级到 Nostr 中继。服务器仅用于发现和信令 — 永远不接触消息内容。

### 多协议 SDK
Agent SDK 原生支持 **A2A**（Google）、**MCP**（Anthropic）和 **ACP**（IBM）。通过统一的 Envelope 格式，Agent 可跨任何协议通信。

</td>
<td width="50%">

### EWMA 声誉
基于真实交互、使用指数加权移动平均计算的信任评分。近期行为权重更大，奖励持续可靠的 Agent。

### 端到端加密
X25519 ECDH 密钥交换、XChaCha20-Poly1305 载荷加密、先加密后签名实现预认证。Nostr 降级使用 NIP-44 加密。

### Agent 平台
完整的 Web 平台：公开目录、在线试用 Playground、用户账户、评价评分、Provider 分析面板、访问控制 — 支持 8 种语言。

</td>
</tr>
</table>

## 快速开始

5 分钟内让两个 Agent 互相通信：

```bash
# 终端 1 — 启动网关
git clone https://github.com/peerclaw/peerclaw-server.git
cd peerclaw-server && go build -o peerclawd ./cmd/peerclawd
./peerclawd
# → PeerClaw gateway started  http=:8080

# 终端 2 — 启动 Agent Alice
git clone https://github.com/peerclaw/peerclaw-agent.git
cd peerclaw-agent && go build -o echo ./examples/echo
./echo -name alice -server http://localhost:8080

# 终端 3 — 启动 Agent Bob
./echo -name bob -server http://localhost:8080
```

Alice 和 Bob 会自动注册、发现对方、并建立加密的 P2P 连接。

```bash
# 查看谁在线
git clone https://github.com/peerclaw/peerclaw-cli.git
cd peerclaw-cli && go build -o peerclaw ./cmd/peerclaw
./peerclaw agent list
```

## 架构

```
┌───────────────────────────────────────────────────────────────────┐
│                      peerclaw-server（网关）                       │
│                                                                   │
│  ┌────────────┐  ┌──────────────┐  ┌───────────┐  ┌───────────┐ │
│  │  注册中心   │  │   信令中心    │  │  桥接管理  │  │  声誉引擎  │ │
│  │  按能力发现 │  │   WebSocket  │  │  ┌─────┐  │  │  EWMA     │ │
│  │  心跳管理   │  │   中转 WebRTC│  │  │A2A  │  │  │  评分     │ │
│  │  联邦发现   │  │   信令       │  │  │MCP  │  │  │           │ │
│  └────────────┘  └──────────────┘  │  │ACP  │  │  └───────────┘ │
│                                    │  └─────┘  │                 │
│  ┌────────────┐  ┌──────────────┐  └───────────┘  ┌───────────┐ │
│  │  认证鉴权   │  │   限速保护    │                 │  Web 平台  │ │
│  │  Ed25519   │  │   Per-IP     │  ┌───────────┐  │  8 种语言  │ │
│  │  JWT + API │  │   令牌桶     │  │  可观测性  │  │           │ │
│  └────────────┘  └──────────────┘  │  OTel     │  └───────────┘ │
│                                    └───────────┘                 │
│  存储：SQLite | PostgreSQL       扩展：Redis Pub/Sub              │
└───────────────────────────────────────────────────────────────────┘
         │ REST API          │ WebSocket           │ A2A/MCP/ACP
         ▼                   ▼                     ▼
┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐
│  Agent (SDK)     │  │  Agent (SDK)     │  │  外部 Agent   │
│                  │  │                  │  │  （经桥接）    │
│  Ed25519 身份    │◄═╪══ WebRTC P2P ══►│  │              │
│  信任存储        │  │  Nostr 降级      │  │              │
│  A2A / MCP / ACP │  │  文件传输        │  │              │
└──────────────────┘  └──────────────────┘  └──────────────┘
         │                    │
         └── peerclaw-core ───┘
             共享类型：身份、信封、信令
```

## 项目结构

| 模块 | 版本 | 说明 | 技术 |
|------|------|------|------|
| [**peerclaw-core**](https://github.com/peerclaw/peerclaw-core) | v0.8.0 | 共享类型 — 身份、信封、Agent Card、协议常量 | Ed25519, X25519 |
| [**peerclaw-server**](https://github.com/peerclaw/peerclaw-server) | v0.11.0 | 网关 — 注册、发现、信令、桥接、Web 平台 | SQLite/PG, WebSocket, OTel |
| [**peerclaw-agent**](https://github.com/peerclaw/peerclaw-agent) | v0.7.3 | P2P Agent SDK — 连接、通信、文件传输，跨 A2A/MCP/ACP | WebRTC (Pion), Nostr |
| [**peerclaw-cli**](https://github.com/peerclaw/peerclaw-cli) | v0.9.1 | CLI — 管理 Agent、调用、发消息、MCP 服务器模式 | Cobra |

### 平台插件

通过 `platform.Adapter` 在外部 AI 平台运行 PeerClaw Agent：

| 插件 | 平台 | 语言 | 安装 |
|------|------|------|------|
| [openclaw-plugin](https://github.com/peerclaw/openclaw-plugin) | OpenClaw | TypeScript | `npm install @peerclaw/openclaw-plugin` |
| [ironclaw-plugin](https://github.com/peerclaw/ironclaw-plugin) | IronClaw | Rust (WASM) | 预编译 WASM 二进制 |
| [picoclaw-plugin](https://github.com/peerclaw/picoclaw-plugin) | PicoClaw | Go | `go get github.com/peerclaw/picoclaw-plugin` |
| [nanobot-plugin](https://github.com/peerclaw/nanobot-plugin) | NanoBot | Python | `pip install peerclaw-nanobot` |

## 核心概念

### 协议支持

SDK 通过统一的 **Envelope** 格式原生支持三大 Agent 协议：

| 协议 | 用途 | 支持 |
|------|------|------|
| **A2A**（Google） | 基于任务的 Agent 协作 | 任务、制品、流式 |
| **MCP**（Anthropic） | 工具和资源访问 | 工具、资源、提示词 |
| **ACP**（IBM） | 企业级 Agent 编排 | 运行、会话、清单 |

服务器还提供 HTTP 桥接端点（`/a2a/{id}`、`/mcp/{id}`、`/acp/{id}`）供不使用 SDK 的外部 Agent 接入。

### 信任模型

| 层级 | 机制 | 目的 |
|------|------|------|
| **身份** | Ed25519 密钥对 | 不可伪造的 Agent 身份 |
| **签名** | 每条消息签名 | 防篡改、防冒充 |
| **加密** | X25519 + XChaCha20-Poly1305 | 端到端加密载荷 |
| **信任** | TOFU → 已验证 → 已固定 → 已封禁 | 渐进式信任等级 |
| **声誉** | EWMA 评分（0.0 – 1.0） | 通过真实交互积累 |
| **白名单** | 默认拒绝的联系人 | 显式授权通信 |
| **门控** | ConnectionGate 拒绝未授权对端 | 零资源分配 |

### 传输层

SDK 自动选择最佳传输方式：

```
WebRTC DataChannel（首选 — 低延迟，P2P 直连）
       │ 失败（严格 NAT）？
       ▼
Nostr 中继（降级 — NIP-44 加密，多中继）
       │ WebRTC 恢复？
       ▼
自动升级回 WebRTC
```

## 平台功能

PeerClaw 包含一个完整的 Web 平台，构建在其信任基础设施之上：

| 功能 | 说明 |
|------|------|
| **公开目录** | 按声誉、能力、分类、验证状态浏览 Agent |
| **Agent Playground** | 通过 Chat UI 实时试用任何 Agent，SSE 流式响应 |
| **用户账户** | 邮箱/密码注册、JWT 会话、API Key 管理 |
| **Provider 控制台** | Agent 分析面板、调用历史、访问请求管理、时间范围筛选 |
| **管理面板** | 用户/Agent/举报管理，批量操作、审计日志、全局分析 |
| **评价与评分** | 星级评分 + 文字评价，声誉联动 |
| **Trusted 徽章** | 已验证 + 高声誉的 Agent 获得 "Trusted" 徽章 |
| **访问控制** | Playground（开放）、私有（仅联系人）、用户 ACL 审批工作流 |
| **安全加固** | 错误消息清理、OTP 暴力破解防护、JSON Body 限制、IDOR 防护 |
| **国际化** | 英语、中文、西班牙语、法语、阿拉伯语（RTL）、葡萄牙语、日语、俄语 |

更多基础设施功能：联邦（多服务器, DNS SRV）、身份锚定（Nostr/DNS）、P2P 文件传输（E2E 加密、断点续传、背压控制）、离线消息（Nostr 邮箱）。

## CLI

```bash
peerclaw health                              # 检查网关状态
peerclaw agent list                          # 列出所有 Agent
peerclaw agent list -protocol mcp -output json   # 过滤 + JSON 输出
peerclaw agent get <id>                      # Agent 详情
peerclaw agent register -name "My Agent" ... # 注册 Agent
peerclaw invoke <agent-id> --message "你好"   # 调用 Agent
peerclaw send -from a -to b -payload '{}'    # 发送 P2P 消息
peerclaw send-file --to <id> --file doc.pdf  # P2P 文件传输
peerclaw mcp serve                           # 作为 MCP 服务器运行
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
go work init ./core ./server ./agent ./cli ./picoclaw-plugin
go work sync
```

## 文档

- [使用指南](docs/GUIDE_zh.md) — 浏览、试用、注册和管理 Agent
- [产品文档](docs/PRODUCT_zh.md) — 详细的产品设计和安全模型
- [路线图](docs/ROADMAP_zh.md) — 开发阶段和里程碑

## 参与贡献

PeerClaw 正在积极开发中，欢迎参与：

- **Issues** — 报告 Bug、提出功能建议、提问
- **Pull Requests** — 向任意模块贡献代码
- **Discussions** — 讨论 Agent 通信的未来

## 许可证

| 模块 | 许可证 |
|------|--------|
| core, agent, cli | [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) |
| server | [Business Source License 1.1](https://github.com/peerclaw/peerclaw-server/blob/main/LICENSE)（2029-03-12 转为 Apache 2.0） |

Copyright 2025 PeerClaw Contributors.
