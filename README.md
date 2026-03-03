# PeerClaw

**让 AI Agent 像人一样自由通信。**

PeerClaw 是一个去中心化优先的 AI Agent 通信框架。Agent 通过密码学身份互相识别，通过 WebRTC 直连对话，在 NAT 穿越失败时回退到 Nostr relay，并通过协议桥接实现 A2A / ACP / MCP 互操作。

## 愿景

当前 AI Agent 生态面临严重的通信碎片化问题：

- **协议割裂** — A2A、ACP、MCP 各自为政，Agent 无法跨协议交流
- **中心化依赖** — Agent 通信必须经过平台服务器中转，增加延迟和单点故障
- **身份缺失** — Agent 缺少统一的密码学身份，无法验证消息来源
- **安全薄弱** — 大多数 Agent 通信方案缺少端到端的安全保障

PeerClaw 的回答：

- **去中心化优先** — WebRTC P2P 直连，Nostr relay 兜底，不依赖任何单一服务
- **协议桥接** — 内置 A2A / ACP / MCP 适配器，统一转换为 PeerClaw Envelope
- **密码学身份** — 每个 Agent 拥有 Ed25519 密钥对，公钥即身份
- **四层安全** — 连接级 TOFU + 消息级签名 + 端到端加密 (XChaCha20-Poly1305) + 执行级沙箱

## 架构

```
┌─────────────────────────────────────────────────────────────┐
│                        peerclaw-server                      │
│                                                             │
│   ┌────────────┐  ┌──────────────┐  ┌───────────────────┐  │
│   │  Registry  │  │   Signaling  │  │  Bridge Manager   │  │
│   │  (发现)     │  │   Hub (信令)  │  │  A2A / ACP / MCP  │  │
│   └────────────┘  └──────────────┘  └───────────────────┘  │
│         │                │                    │             │
└─────────┼────────────────┼────────────────────┼─────────────┘
          │                │                    │
    ┌─────┴─────┐    ┌─────┴─────┐        ┌────┴────┐
    │  Agent A  │◄──►│  Agent B  │   外部  │ A2A/MCP │
    │ (SDK)     │P2P │ (SDK)     │  Agent  │ Agent   │
    └───────────┘    └───────────┘        └─────────┘
         │                │
    WebRTC DataChannel / Nostr relay
```

**通信流程：** 注册 → 发现 → 信令握手（含 X25519 密钥交换） → P2P 连接（WebRTC 优先，自动降级 Nostr） → 加密签名消息交换

## 子项目

| 仓库 | 说明 | 状态 |
|------|------|------|
| [peerclaw-core](https://github.com/peerclaw/peerclaw-core) | 核心共享类型库（身份、信封、协议常量） | Active |
| [peerclaw-server](https://github.com/peerclaw/peerclaw-server) | 中心化平台（注册/发现/信令/桥接） | Active |
| [peerclaw-agent](https://github.com/peerclaw/peerclaw-agent) | P2P Agent SDK（WebRTC + Nostr + 安全） | Active |

## 快速体验

5 分钟跑通 P2P 通信 demo：

### 1. 克隆并构建

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
```

### 2. 启动 Server

```bash
./server/peerclawd
# 输出: PeerClaw gateway started  http=:8080  grpc=:9090
```

### 3. 启动两个 Echo Agent

```bash
# 终端 1
./agent/echo -name alice -server http://localhost:8080

# 终端 2
./agent/echo -name bob -server http://localhost:8080
```

两个 Agent 将自动注册到 Server，通过信令建立 WebRTC P2P 连接。

## 本地开发

```bash
# 项目使用 Go workspace 管理多模块
# 确保三个子仓库在正确位置：core/ server/ agent/

# 同步 workspace
go work sync

# 构建所有模块
cd core && go build ./... && cd ..
cd server && go build ./... && cd ..
cd agent && go build ./... && cd ..

# 运行测试
cd server && CGO_ENABLED=1 go test ./... && cd ..
cd agent && go test ./... && cd ..
```

## 文档

- [产品文档](docs/PRODUCT.md) — 详细的产品设计、架构和安全模型
- [路线图](docs/ROADMAP.md) — 从当前到去中心化演进的五阶段计划

## 社区与贡献

PeerClaw 正处于早期阶段，欢迎参与：

- 提交 Issue 报告问题或建议功能
- 提交 Pull Request 贡献代码
- 参与讨论 Agent 通信的未来

## License

MIT
