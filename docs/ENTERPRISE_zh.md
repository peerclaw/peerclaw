[English](ENTERPRISE.md) | **中文**

# PeerClaw 企业内网部署指南

在内部网络中部署 PeerClaw，实现安全的 Agent 间通信 — 无需外部依赖。

## 架构

```
┌─────────────────────────────────────────────────────┐
│                     企业内网                          │
│                                                     │
│   ┌──────────────────┐                              │
│   │  peerclaw-server │  ← 单实例                     │
│   │  :8080           │    （发现 + 信令）              │
│   └────────┬─────────┘                              │
│            │                                        │
│   ┌────────┼────────────────────┐                   │
│   │        │        │           │                   │
│   ▼        ▼        ▼           ▼                   │
│ Agent A  Agent B  Agent C  Agent D                  │
│ (计费)   (审计)   (通知)    (发票)                    │
│                                                     │
│   所有通信通过 WebRTC P2P 或服务器中继                  │
│   无 Nostr、无 STUN/TURN                             │
└─────────────────────────────────────────────────────┘
```

## 快速开始（5 行 Go 代码）

```go
package main

import (
    "context"
    agent "github.com/peerclaw/peerclaw-agent"
    "github.com/peerclaw/peerclaw-core/envelope"
)

func main() {
    // 创建 Agent — 只需名称、服务器地址和能力列表。
    a, _ := agent.NewSimple("invoice-processor", "http://peerclaw.internal:8080",
        "process-invoice", "query-status",
    )

    // 预配置可信联系人。
    a.ImportContacts([]string{"agent-billing", "agent-audit", "agent-notify"})

    // 注册 Handler。
    a.Handle("process-invoice", func(ctx context.Context, env *envelope.Envelope) (*envelope.Envelope, error) {
        return &envelope.Envelope{Payload: []byte(`{"status":"ok"}`)}, nil
    })

    // 启动服务。
    ctx := context.Background()
    a.Start(ctx)
    select {} // 阻塞
}
```

### `NewSimple` 做了什么

`NewSimple(name, serverURL, capabilities...)` 等价于：

```go
agent.New(agent.Options{
    Name:         name,
    ServerURL:    serverURL,
    Capabilities: capabilities,
})
```

其余选项全部使用安全默认值：
- **自动生成 Ed25519 密钥对**（无需密钥文件管理）
- **仅服务器模式** — 无 Nostr relay、无 STUN/TURN
- **非 Serverless** — 依赖中心化 peerclaw-server

### 使用 `ImportContacts` 预配置信任

企业环境中 Agent 都是预知的。使用 `ImportContacts` 批量导入可信 Agent ID：

```go
a.ImportContacts([]string{"agent-billing", "agent-audit", "agent-notify"})
```

导入的 Agent 自动设为 `TrustVerified` 级别，可立即双向通信，无需 TOFU 握手。

## Docker Compose 部署

使用项目提供的 `deploy/docker-compose.yaml` 启动服务：

```bash
cd deploy

# 创建 .env 文件
cat > .env << 'EOF'
POSTGRES_USER=peerclaw
POSTGRES_PASSWORD=changeme
POSTGRES_DB=peerclaw
REDIS_PASSWORD=changeme
EOF

docker compose up -d
```

启动的服务：
- **peerclaw-server**（端口 8080）
- **PostgreSQL** 持久化存储
- **Redis** 跨节点信令（单节点可选）
- **Caddy** 反向代理（端口 80/443）

如只需最简部署，可修改配置使用 SQLite，移除 PostgreSQL/Redis 服务。

## 安全建议

### 网络

- peerclaw-server 部署在防火墙后 — 仅对内网开放
- Agent 与 Server 之间启用 TLS（通过 Caddy 或负载均衡器配置）
- 限制 Agent → Server 通信仅使用 peerclaw 端口

### 身份与密钥

- 每个 Agent 首次运行时自动生成 Ed25519 密钥对
- 重启后保持身份，使用 `Options` 中的 `KeypairPath`：
  ```go
  agent.New(agent.Options{
      Name:        "invoice-processor",
      ServerURL:   "http://peerclaw.internal:8080",
      KeypairPath: "/etc/peerclaw/agent.key",
  })
  ```
- 密钥文件设置严格权限（`chmod 600`）

### 信任管理

- 使用 `ImportContacts` 在启动时预配置已知 Agent
- 动态信任场景，注册 `ConnectionRequestHandler`：
  ```go
  a.OnConnectionRequest(func(ctx context.Context, req *agent.ConnectionRequest) bool {
      // 检查内部注册表或允许所有内部 Agent
      return isInternalAgent(req.FromAgentID)
  })
  ```
- 拉黑不信任的 Agent：`a.BlockAgent(agentID)`

## 完整示例

查看 [agent/examples/enterprise/main.go](../agent/examples/enterprise/main.go) 获取完整工作示例。
