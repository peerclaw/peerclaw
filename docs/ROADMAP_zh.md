[English](ROADMAP.md) | **中文**

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

## Phase 2: Transport & Security Hardening (已完成)

加固传输层和安全机制，提升连接成功率。

- [x] **Nostr relay 完整实现**
  - NIP-44 加密（基于 `fiatjaf.com/nostr` 库）
  - 多 relay 支持与故障切换（发布到全部，订阅去重）
  - relay 健康检查（指数退避重连，3 次失败移出活跃集）
- [x] **NAT 穿越优化**
  - TURN server 集成（信令连接后 Server 推送 ICE config）
  - ICE candidate 筛选与排序优化（host > srflx > relay）
  - 连接质量监控（RTT、丢包率、吞吐统计）
- [x] **自动传输选择**
  - WebRTC → Nostr 自动降级（Transport Selector）
  - 连接恢复后自动升级（后台探测）
  - 传输健康评分（滚动窗口成功/失败计数）
- [x] **Trust Store 增强**
  - CLI 管理工具 `peerclaw-trust`（list / verify / pin / revoke / export / import）
  - 信任等级（Unknown / TOFU / Verified / Blocked / Pinned）
  - 信任事件通知（OnTrustChange 回调）
- [x] **端到端加密**
  - X25519 密钥交换（从 Ed25519 seed 派生，信令阶段交换公钥）
  - XChaCha20-Poly1305 消息加密（HKDF-SHA256 密钥派生）
  - Nostr 传输额外使用 NIP-44 格式封装

## Phase 3: Protocol Ecosystem (已完成)

完整实现三大协议桥接，让 PeerClaw 成为协议互操作中枢。

- [x] **JSON-RPC 2.0 共享库**
  - Request / Response / Error / Notification 类型
  - ParseMessage 自动识别消息类型
  - 标准错误码（-32700 ~ -32600）
  - A2A 和 MCP 共用
- [x] **A2A 完整适配**
  - Task 生命周期（create / working / complete / cancel / fail）
  - Artifact 模型（多 Part 支持：text / file / structured data）
  - Agent Card 标准兼容（GET /.well-known/agent.json）
  - JSON-RPC 入站 Handler（message/send / tasks/get / tasks/cancel）
  - 出站 SendMessage → Task 响应 → Envelope 转换
- [x] **MCP 完整适配（Streamable HTTP）**
  - Tool 模型（tools/list / tools/call）
  - Resource 模型（resources/list / resources/read）
  - Prompt 模型（prompts/list / prompts/get）
  - Session 管理（initialize 握手 + Mcp-Session-Id）
  - Streamable HTTP 传输（POST + GET SSE 端点）
- [x] **ACP 完整适配**
  - 完整消息类型（MessagePart 含 content_type / content / content_url）
  - 会话管理（Session + Run 生命周期追踪）
  - Agent Manifest 查询（GET /acp/agents/{name}）
  - Run 管理（create / get / cancel）
- [x] **协议自动协商**
  - Negotiator 自动选择最佳协议路径
  - 同协议直连 > 跨协议翻译（优先级：A2A > MCP > ACP）
  - 跨协议翻译（A2A ↔ MCP ↔ ACP）
- [x] **Agent 能力声明增强**
  - 结构化 Skills 声明（A2A 兼容：name / description / input_modes / output_modes）
  - 结构化 Tools 声明（MCP 兼容：name / description / input_schema）
  - HasSkill() / HasTool() 查询方法
  - SQLite 持久化（JSON 序列化存储）
- [x] **Server 路由集成**
  - 协议端点路由（POST /a2a, POST /mcp, GET/POST /acp/*）
  - 通用桥接发送端点（POST /api/v1/bridge/send）
  - Bridge Forwarder（bridge inbox → signaling hub → agent）
  - bridge_message 信令消息类型

## Phase 4: Production Readiness (已完成)

面向生产环境的稳定性、可观测性和运维能力。

- [x] **可观测性**
  - OpenTelemetry traces（可选启用，OTLP gRPC 推送）
  - OpenTelemetry metrics（HTTP 请求率/延迟、WebSocket 连接数、Agent 注册数、Bridge 消息吞吐）
  - 结构化日志增强（中间件链：Recovery → RequestID → Tracing → Logging → RateLimit → MaxBody）
  - Grafana dashboard 模板（docs/grafana/peerclaw-overview.json）
- [x] **水平扩展**
  - Redis Pub/Sub 跨节点信令（Broker 接口 + LocalBroker + RedisBroker）
  - PostgreSQL 后端（JSONB + GIN 索引，通过 factory.go 按配置选择驱动）
  - 增强健康检查（组件级状态：database、signaling）
- [x] **审计日志**
  - Agent 注册/注销记录
  - 消息路由审计（Bridge send）
  - 安全事件记录（速率限制、信令连接/断开）
  - 独立 slog.Logger 实例（stdout 或文件输出）
- [x] **速率限制与防护**
  - Per-IP token bucket 请求速率限制（golang.org/x/time/rate）
  - WebSocket 连接数限制（Hub.maxConns）
  - 消息大小限制（http.MaxBytesReader 中间件）
- [x] **CLI 工具 (peerclaw-cli)**
  - `peerclaw agent list|get|register|delete` — Agent 管理
  - `peerclaw send` — 通过 Bridge 发送消息
  - `peerclaw health` — 服务器健康检查
  - `peerclaw config show|set` — CLI 配置管理

## Phase 5: Decentralized Evolution (已完成)

向完全去中心化演进，实现无 Server 的 Agent 通信。

- [x] **接口抽象 + Agent Card 扩展**
  - Discovery 接口（RegistryClient / DHTDiscovery / CompositeDiscovery）
  - SignalingClient 接口（WebSocket / NostrSignaling / CompositeSignaling）
  - Agent 结构体改用接口（向后兼容）
  - PeerClawExtension 新字段（nostr_pubkey / dht_node_id / reputation_score / nostr_relays / identity_anchor）
  - 新信令消息类型（federation_forward / dht_ping/store/find_node/find_value）
- [x] **DHT 去中心化发现**
  - 最小 Kademlia DHT（160-bit NodeID，K=20 k-bucket 路由表）
  - DHT 协议消息（Ping / Store / FindNode / FindValue RPC）
  - DHT 传输层（Nostr event kind 20005 + InMemory 测试传输）
  - DHT 协调器（Bootstrap / Put / Get / FindNode / bucket refresh / data republish）
  - DHTDiscovery 实现 Discovery 接口（主键 + 能力索引）
  - CompositeDiscovery（Server 优先 + DHT fallback）
- [x] **多信令节点联邦**
  - FederationConfig（静态配置 + DNS SRV 发现）
  - FederationService（连接对端、ForwardSignal、QueryAgents）
  - DNS SRV 发现（_peerclaw._tcp.<domain>）
  - FederationBroker（本地有目标则 local.Publish，否则 federation.ForwardSignal）
  - Hub.HasAgent() 判断本地连接
  - 联邦 HTTP 端点（/api/v1/federation/signal、/api/v1/federation/discover）
  - DiscoverFederated 方法
- [x] **Agent 信誉系统**
  - ReputationStore：per-peer EWMA 评分（0.0-1.0），6 种行为类型
  - RecordEvent / GetScore / IsMalicious（< 0.15 阈值）
  - TrustStore 集成（SetReputationStore / IsAllowedWithReputation）
  - 信誉 Gossip（Nostr kind 30078，第二手信誉权重 0.3x，仅接受 TrustVerified+ 的 peer）
  - JSON 文件持久化
- [x] **无 Server 纯 P2P 模式**
  - NostrSignaling（event kind 20006，NIP-44 加密）
  - CompositeSignaling（WebSocket 优先 + Nostr fallback）
  - MessageCache 离线消息缓存（per-destination 队列，TTL 过期，JSON 持久化）
  - OnPeerAdded 回调（新 peer 连接时刷新缓存队列）
  - Serverless 模式 Options（DHTEnabled / Serverless / ICEServers / MessageCachePath）
- [x] **链上身份锚定（可选）**
  - IdentityAnchor 接口（Publish / Verify / Resolve / RecoveryKeys）
  - NostrAnchor 实现（Nostr kind 10078 replaceable event，双向密钥绑定）
  - 域名绑定验证（DNS TXT 记录 peerclaw-verify=<fingerprint>）
  - 多签恢复（threshold-of-n recovery keys）
- [x] **CLI Phase 5 命令**
  - `peerclaw dht bootstrap|lookup`
  - `peerclaw federation status|peers`
  - `peerclaw reputation show|list`
  - `peerclaw identity anchor|verify`
- [x] **集成测试**
  - DHT-only 发现通信（无 Server）
  - Server + DHT 混合模式降级
  - 恶意行为触发信誉隔离
  - 信誉 Gossip 跨 peer 传播
  - DHT Agent Card 存取
  - 离线消息缓存投递
