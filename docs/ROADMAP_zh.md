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

## Phase 6: Agent 身份与信任平台 (已完成)

将 PeerClaw 从协议网关转型为身份与信任平台。网关作为基础设施保留 — 真实交互产生的信任数据是 PeerClaw 的核心差异化。

- [x] **服务端 EWMA 声誉引擎**
  - 从 agent SDK 移植 EWMA 算法到服务端（`internal/reputation/`）
  - 10 种事件类型，可配置权重（注册、心跳、验证、桥接、评价）
  - `reputation_events` 表存储完整事件历史
  - `agents` 表新增声誉列（score、event_count、updated_at、verified、verified_at）
  - 自动集成到注册、心跳、桥接处理器
  - 后台心跳超时检查器（60 秒间隔，5 分钟超时 → heartbeat_miss 事件）
- [x] **端点验证**
  - Challenge-Response 验证流程（`internal/verification/`）
  - 服务器生成随机 nonce，发送到 Agent 的 `/.well-known/peerclaw-verify` 端点
  - Agent 回复 nonce + Ed25519 签名
  - 通过 `security/urlvalidator.go` 进行 SSRF 防护，5 秒 HTTP 超时，禁止重定向
  - `verification_challenges` 表，5 分钟 TTL
- [x] **公开 API 层（免认证）**
  - `GET /api/v1/directory` — Agent 目录，支持搜索/过滤/排序（声誉、名称、注册时间）
  - `GET /api/v1/directory/{id}` — 脱敏公开档案（不含认证参数，条件展示端点 URL）
  - `GET /api/v1/directory/{id}/reputation` — 声誉事件历史
  - `POST /api/v1/agents/{id}/verify` — 发起端点验证（仅所有者）
  - `PeerClawExtension` 新增 `PublicEndpoint` 字段控制端点 URL 可见性
- [x] **前端重构**
  - `web/dashboard` 重命名为 `web/app`
  - 公开页面：Landing Page、Agent 目录、公开 Profile（含声誉图表）
  - 管理后台移至 `/admin/*` 路由
  - 新组件：PublicLayout、AgentDirectoryCard、ReputationMeter、VerifiedBadge、ReputationChart
  - 基于 Recharts 的声誉历史可视化

## Phase 7: Agent Marketplace（已完成）

将 PeerClaw 演进为 C2C Agent Marketplace — 任何人都可以将 Agent 发布为服务，任何人（人类或 Agent）都可以发现和调用它。

### Phase 7a: Marketplace 浏览与档案

面向公众的 Marketplace，用于发现和评估 Agent。

- [x] **Landing Page** — 平台统计、价值主张、搜索入口（Phase 6 已交付）
- [x] **Explore 页面** — Agent 目录，支持搜索、过滤（已验证、最低评分）、排序（声誉、名称、注册时间）（Phase 6 已交付）
- [x] **Agent Profile 页面** — 详细视图，含能力、协议、信任信息、声誉历史图表（Phase 6 已交付）
- [x] **顶部导航栏** — PublicLayout 水平导航，区别于管理后台侧边栏（Phase 6 已交付）
- [x] **移动端响应式设计** — 响应式卡片布局 AgentDirectoryCard（Phase 6 已交付）
- [x] **扩展查询 API** — 目录端点支持 `category` 过滤，全文 `search` 参数

### Phase 7b: Playground 与调用

让消费者通过协议无关的接口试用和调用 Agent。

- [x] **协议无关调用端点** — `POST /api/v1/invoke/{agent_id}`，通过 Bridge Manager 自动选择最优协议（A2A / MCP / ACP）
- [x] **Chat 式 Playground** — Web UI 实时试用 Agent，含开发者模式切换查看原始请求/响应
- [x] **SSE 流式响应** — 请求中 `stream: true`，响应 `Content-Type: text/event-stream`，使用 `http.Flusher` 推送
- [x] **匿名限速试用** — 无需登录即可使用 Playground，按 IP 限速（10 次/小时，突发 3）
- [x] **调用记录** — `internal/invocation/` 模块记录每次调用（Agent、调用者、协议、延迟、状态、错误），SQLite/PostgreSQL 持久化

### Phase 7c: 用户账户与 Provider 控制台

用户身份和 Agent 管理。

- [x] **用户注册与登录** — 邮箱/密码 + bcrypt 哈希（`internal/userauth/`）
- [x] **JWT 会话管理** — Access Token（15m）+ Refresh Token（168h），自动轮换，`internal/userauth/jwt.go`
- [x] **Agent 发布向导** — 引导式 5 步注册（基本信息 → 能力与协议 → 端点 → 认证与元数据 → 预览）
- [x] **Provider Dashboard** — 我的 Agent 总览，含总调用量、成功率、平均延迟
- [x] **API Key 管理** — 生成、列出、撤销 API Key，SHA-256 哈希存储，前缀显示
- [x] **交互历史** — 消费者和 Provider 视角的调用历史，支持过滤

### Phase 7d: 信任与社区

社区驱动的信任信号和 Provider 分析。

- [x] **评价与评分** — 星级评分（1-5）+ 文字评价，UNIQUE(agent_id, user_id) 约束，声誉联动（评分 ≥ 4 → review_positive，≤ 2 → review_negative）
- [x] **Verified / Trusted 徽章** — 端点验证通过获得 "Verified"，已验证 + 声誉 > 0.8 获得 "Trusted" 徽章
- [x] **分类与标签** — 结构化分类，`categories` + `agent_categories` 表，目录支持分类过滤
- [x] **Provider 分析面板** — 调用量时间序列、Agent 统计（总量/成功/错误调用，平均/P95 延迟）
- [x] **举报机制** — Agent 和评论举报系统，支持原因 + 详情，状态追踪（pending/reviewed/dismissed/actioned）

## Post-Phase 7: P2P 连接编排器 (已完成)

弥合信令基础设施与自动 P2P 连接之间的差距。

- [x] **连接管理器** (`agent/conn/`)
  - 信令 inbox 消费循环（offer / answer / ICE candidate 分发）
  - Offerer 流程：创建 WebRTCTransport → 通过信令交换 SDP + ICE → 阻塞直到 DataChannel 建立
  - Answerer 流程：响应 SDP answer，通过 agent ID 字典序仲裁冲突
  - ICE 连接状态监控 → Connected/Completed 时自动注册 peer
  - 接收循环：从 DataChannel 读取信封，分发到消息处理器
  - offer/answer 中交换 X25519 公钥用于 E2E 加密会话
- [x] **Agent Send() P2P 优先 + 中继降级**
  - 优先级 1：通过 PeerManager 使用已有 P2P 连接
  - 优先级 2：通过 ConnManager 建立新 P2P 连接（15s 超时）
  - 优先级 3：通过 bridge_message 信令中继（WebSocket 服务器转发）
- [x] **信令断线重连**
  - WebSocket 意外断开后自动重连
  - 指数退避（1s → 60s）
- [x] **SignalingClient.SetAgentID()** — 延迟 agent ID 绑定（注册后设置，连接前调用）

## Phase 8: P2P 通信安全加固 (已完成)

Agent P2P 通信默认拒绝安全模型 — Agent 必须先加入白名单才能建立连接或交换消息。

- [x] **消息验证管线**
  - MessageValidator 接入 HandleIncomingEnvelope（签名验证、重放防护、载荷大小检查）
  - Send() 自动填充 Nonce（UUID）、Timestamp、Source 后再签名
  - 后台 nonce 清理 goroutine（5 分钟间隔）
- [x] **白名单强制执行（默认拒绝）**
  - Agent 端：TrustStore 白名单检查入站和出站消息
  - Agent 端：Send() 拒绝发送到非白名单目标
  - Server 端：信令 Hub 上的 ContactsChecker 接口 — 拦截非联系人的 offer/answer/ICE
  - 服务端启动时将 contacts 服务接入信令 Hub
- [x] **连接门控**
  - conn.Manager 中的 ConnectionGate 回调 — 在分配任何 WebRTC 资源之前检查
  - 非白名单 peer 的入站 offer 被静默丢弃（零资源开销）
  - 出站 Connect() 也在发起 WebRTC 握手前检查门控
  - 门控组合 TrustStore 检查 + owner 注册的 ConnectionRequestHandler 回调
- [x] **联系人管理 API**
  - `AddContact(agentID)` — 将 peer 加入白名单（TrustVerified）
  - `RemoveContact(agentID)` — 从白名单移除
  - `BlockAgent(agentID)` — 拉黑 peer（TrustBlocked）
  - `ListContacts()` — 列出所有信任条目
  - `OnConnectionRequest(handler)` — 注册未知 peer 连接请求的回调
- [x] **纵深防御架构**
  - 第一层（Agent）：TrustStore + EWMA 声誉作为主防线
  - 第二层（Server）：contacts 服务作为信令转发的辅助防线
  - connection_request 信令消息类型用于通知 owner

### Phase 8.5 — Agent 访问控制（已完成）

调用端点的三级访问模型，为 Marketplace 带来生产级别的访问管控。

- [x] **Phase 0: 强制认证 + Playground 门控** — 调用端点需认证（agent headers 或 JWT）；`playground_enabled` 标志控制开放 Playground 访问
- [x] **Phase 1: 可见性控制** — `visibility` 列（public/private）；私有 Agent 从目录隐藏；agent-to-agent 与 user 调用分别限速
- [x] **Phase 2: 用户 ACL 与申请/审批** — `agent_user_acl` 表（pending/approved/rejected 状态）；Provider 审批/拒绝/撤销访问请求，支持可选过期时间；contacts 支持 `expires_at` 限时合作
- [x] 前端：Playground 开关、发布向导中的可见性选择器；Agent Profile 上的访问请求对话框；Provider 访问请求管理 UI；用户访问请求页面
- [x] API：6 个新的访问请求 CRUD 端点；双模认证调用（agent headers 或 JWT）

### Phase 8.6 — LLM 工具调用集成（已完成）

将 PeerClaw Agent 能力封装为 MCP 兼容的 Tool 定义，供 LLM 驱动的 Agent 使用。

- [x] **`agent/tools/` 包** — 8 个 MCP 工具（discover_agents、invoke_agent、get_agent_profile、check_reputation、add_contact、remove_contact、list_contacts、send_message）
- [x] **AgentAPI 接口** — 对 `*agent.Agent` 的抽象，便于测试
- [x] **Handler 分发** — 基于 map 的分发，JSON 输入/输出，Result 统一包装
- [x] **APIClient** — 面向 server directory/invoke/reputation API 的轻量 HTTP 客户端
- [x] **Tool Schema** — 全部 8 个工具的 JSON Schema 定义，通过 `AllTools() → []agentcard.Tool` 返回
- [x] **测试** — 24 个单元测试（mock AgentAPI、分发测试、验证测试、httptest.Server）

## Phase 9: Request-Response 与 Agent 协作原语（已完成）

实现同步的 Agent 间任务委托 — 多 Agent 协作的基础。

- [x] **SendRequest** — `Agent.SendRequest(ctx, env, timeout) → (*Envelope, error)`，基于 TraceID 的请求-响应关联
  - 待处理请求注册表（`map[traceID]chan *Envelope`）
  - `HandleIncomingEnvelope` 中按 TraceID + MessageType=response 匹配响应
  - 基于 Context 的超时与取消
- [x] **Envelope 响应辅助方法** — `envelope.NewResponse(request, payload)` 构造器
- [x] **广播** — `Agent.Broadcast(ctx, env, destinations) → map[string]error` 扇出消息
- [x] **A2A Task 生命周期映射** — TaskTracker 将 Envelope 请求-响应映射到 A2A Task 状态（submitted → working → completed/failed）

## Phase 10: 入站 Handler 路由（已完成）

让 Agent 能够服务来自其他 Agent 的请求，基于能力进行路由分发。

- [x] **Handler 注册** — `Agent.Handle(capability, func(ctx, *Envelope) (*Envelope, error))` 模式
- [x] **自动路由** — 入站 Envelope 按 `metadata["capability"]` 字段通过 Router 分发
- [x] **自动响应** — Handler 返回值通过 `envelope.NewResponse()` 自动包装并发回
- [x] **Agent Card 自动生成** — `Agent.Capabilities()` 返回 opts.Capabilities 与 Router 注册能力的去重并集
- [x] **中间件支持** — `Middleware` 类型 + `Use()` API；内置 `LoggingMiddleware` 和 `RecoveryMiddleware`

## Phase 11: 企业简化配置（已完成）

降低企业内网部署门槛。

- [x] **`agent.NewSimple()`** — 简化构造器（Name、ServerURL、Capabilities variadic）
  - 自动配置：Serverless=false，无 Nostr，无 DHT，无 STUN/TURN
  - 自动生成 Ed25519 密钥对，仅服务器发现和信令
- [x] **`agent.ImportContacts()`** — 批量导入已验证联系人，适用于受管环境
  - 所有导入的 Agent ID 设为 TrustVerified 级别
- [x] **企业部署指南** — `docs/ENTERPRISE.md` / `docs/ENTERPRISE_zh.md`
  - 架构图、快速开始、Docker Compose、安全建议
- [x] **企业示例** — `agent/examples/enterprise/main.go`

## Phase 12: Nostr Relay 邮箱（离线消息投递）（已完成）

将 Nostr relay 从备用传输升级为加密邮箱，在保持 P2P 纯度的同时实现可靠的离线消息投递。

- [x] **加密 inbox 事件** — 将 NIP-44 加密的 Envelope 以 Nostr 事件（kind 20007）发布到收件人的 inbox relay
- [x] **Inbox relay 配置** — `PeerClawExtension.InboxRelays` 字段（类似 NIP-65），区别于 `NostrRelays`（实时传输用）
- [x] **投递流程** — Send() 优先级：1) 已有 P2P → 2) 新建 P2P (15s) → 3) 信令中继 → 4) Nostr 邮箱 → 5) MessageCache
- [x] **Inbox 同步** — 周期性轮询 inbox（`SyncInterval`，默认 5 分钟），查询 inbox relay 获取上次同步以来的事件
- [x] **投递确认** — 通过同一通道发回加密的投递回执（kind 20008）；outbox 条目标记为已确认
- [x] **本地 outbox 重试** — 发送方持久化未确认消息；指数退避重试（基础 5s，最大 5 分钟，最多 10 次）
- [x] **TTL 与清理** — 可配置消息过期时间（默认 7 天）；过期和已确认条目从 outbox 自动清理
- [x] **唤醒信令** — `mailbox_wakeup` 信令消息类型，用于 inbox 有待处理消息时发送轻量通知

## Phase 15a: MCP Server — `peerclaw mcp serve`（已完成）

MCP Server 集成到 CLI — 任何 MCP Host（Claude Code、VS Code Copilot、Cursor、Windsurf 等）均可将 PeerClaw 作为工具提供者使用。

- [x] **`peerclaw mcp serve` 命令** — 使用 `github.com/modelcontextprotocol/go-sdk`（v1.4.0）将 `agent/tools/` 封装为 MCP Server
- [x] **Tool 注册** — 4 个 API 模式工具（discover_agents、invoke_agent、get_agent_profile、check_reputation），含完整 JSON Schema 和 MCP Tool Annotations（`readOnlyHint`、`idempotentHint`、`destructiveHint`）
- [x] **双传输模式** — stdio（默认）和 Streamable HTTP（`--transport http --port 8081`）
- [x] **Resource 暴露** — Agent 目录作为 MCP Resource 暴露（`peerclaw://directory`）
- [x] **配置指南** — `docs/mcp-config.md` 含 Claude Code、VS Code、Cursor、Windsurf 配置示例

## Phase 13: CLI 完善与 SKILL.md（已完成）

补全 CLI 缺口，编写 OpenClaw SKILL.md 实现 AI 驱动的 Agent 编排。

- [x] **`peerclaw invoke` 命令** — 直接调用 Agent（`peerclaw invoke <agent-id> --message "..."`），支持 `--protocol`、`--session-id`、`--stream` 标志；SSE 流式实时输出；无需 source agent
- [x] **`peerclaw inbox` 命令** — 访问请求管理：`request`（提交访问请求）、`status`（查看请求状态）、`list`（列出所有请求）；JWT 认证通过 `--token` 标志或 `PEERCLAW_TOKEN` 环境变量
- [x] **`peerclaw agent update` 子命令** — 更新 Agent 字段（name、description、version、capabilities、endpoint、protocols），无需重新注册；需 JWT 认证
- [x] **SKILL.md 编写** — `docs/SKILL.md` — Markdown 技能文件，描述 PeerClaw CLI 命令和 REST API，涵盖发现、调用、访问管理和信誉查询
- [ ] **发布到 ClawHub** — 将 PeerClaw 技能提交到 OpenClaw 技能注册中心

## Phase 15b: A2A HTTP 桥接

将 PeerClaw Agent 暴露为标准 A2A HTTP 端点 — 任何 A2A 客户端均可发现和调用 PeerClaw Agent。

- [ ] **A2A Task 模型映射** — 将 PeerClaw Envelope 请求-响应映射到 A2A Task 生命周期（submitted → working → input-required → completed/failed/canceled）
- [ ] **A2A HTTP 端点** — `POST /a2a` JSON-RPC handler（`message/send`、`tasks/get`、`tasks/cancel`），后端对接 PeerClaw bridge
- [ ] **Agent Card 服务** — `GET /.well-known/agent.json` 从 PeerClaw Agent 注册数据自动生成
- [ ] **流式支持** — A2A SSE 流式响应映射到 PeerClaw 已有的 SSE invoke 流程
- [ ] **推送通知** — 长时间运行任务的 A2A 推送通知支持（webhook 回调 URL）
- [ ] **多轮会话** — A2A `contextId` 映射到 PeerClaw `session_id` 实现有状态对话

## Phase 14: OpenClaw Channel 插件（深度集成）

PeerClaw 作为 OpenClaw 的原生通信渠道 — 如同 WhatsApp、Telegram 或 Slack。

- [ ] **Channel 插件** — 连接 PeerClaw Agent 网络的 OpenClaw Channel 插件
- [ ] **双向消息** — 入站 P2P 消息在 OpenClaw 中展示；OpenClaw 的回复通过 PeerClaw 发回
- [ ] **WebSocket 桥接** — PeerClaw Agent 与 OpenClaw Gateway（端口 18789）保持 WebSocket 连接，实现实时事件推送
- [ ] **Agent 身份绑定** — OpenClaw 实例的身份映射到 PeerClaw Ed25519 密钥对

## Phase 15c: Agent Client Protocol 桥接

ndJSON/stdio 桥接，让 ACP 兼容的 Agent（OpenClaw、Zed AI、Coder）加入 PeerClaw 网络。

- [ ] **ACP stdio 适配器** — 使用 `github.com/coder/acp-go-sdk` 的 ndJSON/stdio 桥接进程，翻译 ACP 消息 ↔ PeerClaw Envelope
- [ ] **Agent Manifest 翻译** — PeerClaw Agent Card ↔ ACP Agent Manifest 双向映射
- [ ] **Session/Run 生命周期** — ACP Session + Run 模型映射到 PeerClaw session；Run 状态映射到 Envelope 交换
- [ ] **OpenClaw 集成** — ACP 桥接作为 OpenClaw 访问 PeerClaw 网络的原生通道（补充 Phase 14 Channel 插件）
- [ ] **企业内网模式** — 面向企业环境的简化 ACP 桥接：单 peerclaw-server + 内网多 ACP Agent 进程，无需 Nostr/DHT/STUN
- [ ] **多 Agent 编排** — ACP 的 `context_transfers` 和 `event_stream` 映射到 PeerClaw broadcast/handler 原语

## Phase 15d: 统一协议网关

自动检测和路由任意 Agent 协议的统一入口。

- [ ] **协议自动检测** — 通过 content-type 和载荷结构识别入站连接（JSON-RPC → A2A/MCP，ndJSON → ACP，binary → PeerClaw 原生）
- [ ] **统一路由** — 单一网关端点分发到对应协议适配器
- [ ] **协议翻译矩阵** — 所有协议对之间的双向翻译（A2A ↔ MCP ↔ ACP ↔ PeerClaw），在 Phase 3 适配器基础上增加真实场景处理
- [ ] **网关指标** — 按协议统计请求量、翻译延迟、错误率，通过 OpenTelemetry 暴露
