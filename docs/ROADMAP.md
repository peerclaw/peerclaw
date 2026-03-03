# PeerClaw Roadmap

## Phase 1: Foundation (已完成)

奠定核心基础设施，验证端到端通信流程。

- [x] **peerclaw-core** — 共享类型库
  - Ed25519 身份（密钥对生成/加载/保存/签名/验证）
  - Envelope 统一消息信封
  - Agent Card 定义（兼容 A2A + PeerClaw 扩展）
  - 协议常量（A2A / ACP / MCP）与传输类型
  - 信令消息类型
- [x] **peerclaw-server** — 中心化平台
  - Agent 注册与注销（REST API）
  - 按能力和协议发现 Agent
  - 心跳与状态管理
  - WebSocket 信令 Hub（offer / answer / ICE candidate）
  - 路由引擎（能力匹配、协议路由）
  - 协议桥接框架（A2A / ACP / MCP 适配器骨架）
  - SQLite 持久化
  - YAML 配置
- [x] **peerclaw-agent** — P2P Agent SDK
  - WebRTC DataChannel 传输
  - Nostr relay 传输（基础）
  - TOFU Trust Store
  - 消息签名与验证
  - Peer Manager
  - Discovery Client
  - Signaling Client
  - Echo Agent 示例

## Phase 2: Transport & Security Hardening

加固传输层和安全机制，提升连接成功率。

- [ ] **Nostr relay 完整实现**
  - NIP-44 加密
  - 多 relay 支持与故障切换
  - relay 健康检查
- [ ] **NAT 穿越优化**
  - TURN server 集成
  - ICE candidate 筛选与排序优化
  - 连接质量监控
- [ ] **自动传输选择**
  - WebRTC → Nostr 自动降级
  - 连接恢复后自动升级
  - 传输健康评分
- [ ] **Trust Store 增强**
  - CLI 管理工具（列出/撤销/导出信任）
  - 信任等级（untrusted / tofu / verified / pinned）
  - 信任事件通知
- [ ] **端到端加密**
  - X25519 密钥交换
  - ChaCha20-Poly1305 消息加密
  - 前向保密（可选）

## Phase 3: Protocol Ecosystem

完整实现三大协议桥接，让 PeerClaw 成为协议互操作中枢。

- [ ] **A2A 完整适配**
  - Task 生命周期（create / update / complete / cancel）
  - Artifact 模型
  - Streaming 支持
  - Agent Card 标准兼容
- [ ] **MCP 完整适配**
  - Tool 模型（list / call / result）
  - Resource 模型（list / read）
  - Prompt 模型
  - Stdio 传输适配
- [ ] **ACP 完整适配**
  - 完整消息类型
  - 会话管理
- [ ] **协议自动协商**
  - Agent 间自动选择最佳协议
  - 能力对齐与兼容检测
- [ ] **Agent 能力声明增强**
  - 结构化能力描述（JSON Schema）
  - 能力版本管理
  - 运行时能力更新

## Phase 4: Production Readiness

面向生产环境的稳定性、可观测性和运维能力。

- [ ] **可观测性**
  - OpenTelemetry traces
  - OpenTelemetry metrics（连接数、消息吞吐、延迟分位数）
  - 结构化日志增强
  - Grafana dashboard 模板
- [ ] **水平扩展**
  - 多 Server 节点部署
  - Redis Pub/Sub 跨节点信令
  - PostgreSQL 后端（替代 SQLite）
  - Session affinity
- [ ] **审计日志**
  - Agent 注册/注销记录
  - 消息路由审计
  - 安全事件记录
- [ ] **速率限制与防护**
  - 请求速率限制
  - 连接数限制
  - 消息大小限制
  - 异常流量检测
- [ ] **CLI 工具 (peerclaw-cli)**
  - `peerclaw agent list` — 查看在线 Agent
  - `peerclaw agent register` — 手动注册
  - `peerclaw send` — 发送消息
  - `peerclaw trust` — 管理信任
  - `peerclaw config` — 配置管理

## Phase 5: Decentralized Evolution

向完全去中心化演进，最终实现无 Server 的 Agent 通信。

- [ ] **DHT 去中心化发现**
  - Kademlia DHT 实现
  - Agent Card 分布式存储
  - DHT bootstrap 节点
- [ ] **多信令节点联邦**
  - Server 间信令消息转发
  - 跨节点 Agent 发现
  - 联邦协议
- [ ] **链上身份锚定（可选）**
  - 公钥哈希上链
  - 域名绑定验证
  - 身份恢复机制
- [ ] **Agent 信誉系统**
  - 通信行为评分
  - 信誉传播算法
  - 恶意 Agent 隔离
- [ ] **无 Server 纯 P2P 模式**
  - DHT 发现 + Nostr 信令
  - 完全无中心节点运行
  - 离线消息缓存
