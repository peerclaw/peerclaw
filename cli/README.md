# peerclaw-cli

PeerClaw 命令行工具。通过 REST API 与 PeerClaw Server 交互，管理 Agent、发送消息、检查服务状态。

## 安装

```bash
cd cli
go build -o peerclaw ./cmd/peerclaw
```

## 使用

### 配置

默认连接 `http://localhost:8080`。可通过环境变量或配置文件修改：

```bash
# 环境变量
export PEERCLAW_SERVER=http://my-server:8080

# 或配置文件
peerclaw config set server http://my-server:8080
peerclaw config show
```

### Agent 管理

```bash
# 列出所有 Agent
peerclaw agent list

# 按协议过滤
peerclaw agent list -protocol a2a

# 查看 Agent 详情
peerclaw agent get <agent-id>

# 注册 Agent
peerclaw agent register -name "MyAgent" -url http://localhost:3000 -protocols a2a,mcp

# 删除 Agent
peerclaw agent delete <agent-id>
```

### 发送消息

```bash
peerclaw send -from agent-a -to agent-b -protocol a2a -payload '{"message": "hello"}'
```

### 健康检查

```bash
peerclaw health

# JSON 输出
peerclaw health -output json
```

### 输出格式

所有列表命令支持 `-output` 参数：

- `table`（默认）：表格格式
- `json`：JSON 格式

```bash
peerclaw agent list -output json
```
