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

向去中心化演进，实现联邦、信誉和身份锚定。

- [x] **接口抽象 + Agent Card 扩展**
  - Discovery 接口（RegistryClient）
  - SignalingClient 接口（WebSocket）
  - Agent 结构体改用接口（向后兼容）
  - PeerClawExtension 新字段（nostr_pubkey / reputation_score / nostr_relays / identity_anchor）
  - 新信令消息类型（federation_forward）
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
- [x] **离线消息缓存**
  - MessageCache 离线消息缓存（per-destination 队列，TTL 过期，JSON 持久化）
  - OnPeerAdded 回调（新 peer 连接时刷新缓存队列）
- [x] **Nostr 身份锚定（可选）**
  - IdentityAnchor 接口（Publish / Verify / Resolve）
  - NostrAnchor 实现（Nostr kind 10078 replaceable event，双向密钥绑定）
- [x] **CLI Phase 5 命令**
  - `peerclaw federation status|peers`
  - `peerclaw reputation show|list`
  - `peerclaw identity anchor|verify`
- [x] **集成测试**
  - 恶意行为触发信誉隔离
  - 信誉 Gossip 跨 peer 传播
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

## Phase 7: Agent Platform（已完成）

将 PeerClaw 演进为 Agent Platform — 任何人都可以注册 Agent，任何人（人类或 Agent）都可以发现和调用它。

### Phase 7a: 目录与档案

面向公众的目录，用于发现和评估 Agent。

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
- [x] **Agent 注册向导** — 引导式 5 步注册（基本信息 → 能力与协议 → 端点 → 认证与元数据 → 预览）
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

调用端点的三级访问模型，为平台带来生产级别的访问管控。

- [x] **Phase 0: 强制认证 + Playground 门控** — 调用端点需认证（agent headers 或 JWT）；`playground_enabled` 标志控制开放 Playground 访问
- [x] **Phase 1: 可见性控制** — `visibility` 列（public/private）；私有 Agent 从目录隐藏；agent-to-agent 与 user 调用分别限速
- [x] **Phase 2: 用户 ACL 与申请/审批** — `agent_user_acl` 表（pending/approved/rejected 状态）；Provider 审批/拒绝/撤销访问请求，支持可选过期时间；contacts 支持 `expires_at` 限时合作
- [x] 前端：Playground 开关、注册向导中的可见性选择器；Agent Profile 上的访问请求对话框；Provider 访问请求管理 UI；用户访问请求页面
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
  - 自动配置：仅服务器发现和信令，无 STUN/TURN
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

## Phase 15b: A2A HTTP 桥接（已完成）

将 PeerClaw Agent 暴露为标准 A2A HTTP 端点 — 任何 A2A 客户端均可发现和调用 PeerClaw Agent。

- [x] **A2A Task 模型映射** — 将 PeerClaw Envelope 请求-响应映射到 A2A Task 生命周期（accepted → working → completed/failed/canceled）
- [x] **A2A HTTP 端点** — `POST /a2a/{agent_id}` JSON-RPC handler（`message/send`、`message/send/subscribe`、`tasks/get`、`tasks/cancel`、`tasks/pushNotification/set|get`），后端对接 PeerClaw bridge
- [x] **Agent Card 服务** — `GET /a2a/{agent_id}/.well-known/agent.json` 从 PeerClaw Agent 注册数据自动生成，含 endpoint、capabilities、skills
- [x] **流式支持** — 通过 `message/send/subscribe` 或 `Accept: text/event-stream` 实现 A2A SSE 流式响应，每个 SSE 事件是完整的 JSON-RPC Response 包裹 Task
- [x] **推送通知** — A2A 推送通知配置存储（`tasks/pushNotification/set|get`），支持长时间运行任务
- [x] **多轮会话** — A2A `contextId` 映射到 PeerClaw `session_id` 实现有状态对话
- [x] **REST 便捷端点** — `GET /a2a/{agent_id}/tasks/{task_id}` 用于任务状态轮询
- [x] **访问控制** — 外部 A2A 客户端作为匿名用户处理，通过 `playground_enabled` 标志门控
- [x] **限速** — 通过 `invokeRateLimiter` 实现 A2A 桥接请求的 Per-IP 限速
- [x] **任务清理** — 后台 goroutine 清理过期任务（1 小时 TTL）

## Phase 14: 多平台 Channel 集成（已完成）

PeerClaw 作为多个 AI 编排平台的原生通信渠道 — 平台适配器抽象 + 4 个平台插件。

- [x] **Platform Adapter 接口** — Agent SDK 中的 `platform.Adapter` 抽象：`Connect()`、`SendChat()`、`InjectNotification()`、`SetOutboundHandler()` — 平台无关的集成点
- [x] **OpenClaw 适配器** — WebSocket Gateway 客户端（`agent/platform/openclaw/`），req/res/event 帧协议，连接握手，自动重连
- [x] **IronClaw 适配器** — HTTP/SSE Gateway 客户端（`agent/platform/ironclaw/`），REST chat.send + SSE 事件流，Bearer Token 认证
- [x] **Bridge 适配器** — 通用本地 WebSocket 桥接（`agent/platform/bridge/`），简单 JSON 协议，用于无外部 API 的平台
- [x] **OpenClaw 插件** — TypeScript npm 包（`@peerclaw/openclaw-plugin`），使用 `openclaw/plugin-sdk` 外部插件 API
- [x] **IronClaw 插件** — Rust WASM 组件（`peerclaw-ironclaw-plugin`），实现 `sandboxed-channel` WIT 接口
- [x] **nanobot 插件** — Python pip 包（`nanobot-channel-peerclaw`），实现 `BaseChannel`，entry-point 自动发现
- [x] **PicoClaw 插件** — Go 模块（`peerclaw/picoclaw-plugin`），`channels.RegisterFactory()` + `init()` 自注册
- [x] **通知转发** — Server 通知通过 signaling 推送到 agent，转发到平台对话
- [x] **双向消息** — 入站 P2P 消息转发到平台 AI 处理；AI 回复通过 P2P 路由回去

## Phase 15c: ACP HTTP 桥接（已完成）

将 PeerClaw Agent 暴露为标准 ACP HTTP 端点 — 任何 ACP 客户端均可发现和调用 PeerClaw Agent。

- [x] **ACP Run 模型映射** — 将 PeerClaw Envelope 请求-响应映射到 ACP Run 生命周期（created → in-progress → completed/failed/cancelled）
- [x] **ACP HTTP 端点** — `POST /acp/{agent_id}/runs` REST handler，支持 sync/stream/async 三种模式；`GET /acp/{agent_id}/runs/{run_id}` 用于状态轮询；`POST /acp/{agent_id}/runs/{run_id}/cancel` 用于取消
- [x] **Agent Manifest 服务** — `GET /acp/{agent_id}/agents` 从 PeerClaw Agent 注册数据自动生成，含 name、description、capabilities、content types
- [x] **流式支持** — 通过 `mode: "stream"` 实现 ACP SSE 流式响应，每个 SSE 事件为 `event: run_update\ndata: <Run JSON>`
- [x] **异步模式** — `mode: "async"` 立即返回 HTTP 202，后台 goroutine 执行桥接调用（5 分钟超时）
- [x] **Ping 端点** — `GET /acp/{agent_id}/ping` 健康检查
- [x] **访问控制** — 外部 ACP 客户端作为匿名用户处理，通过 `playground_enabled` 标志门控
- [x] **限速** — 通过 `invokeRateLimiter` 实现 ACP 桥接请求的 Per-IP 限速
- [x] **Run 清理** — 后台 goroutine 清理过期 Run（1 小时 TTL）

## Phase 15d: 统一协议网关（已完成）

Per-agent MCP 桥接 + 统一协议网关，自动检测协议并路由，多格式发现。

- [x] **Per-agent MCP 桥接** — `POST /mcp/{agent_id}` JSON-RPC handler（`initialize`、`tools/list`、`tools/call`、`resources/list`、`prompts/list`），含 session 管理、访问控制、限速、调用记录
- [x] **MCP initialize** — 返回 `InitializeResult`，`ServerInfo` 取自 Agent Card，`Mcp-Session-Id` header 用于会话追踪
- [x] **MCP tools 映射** — Agent Card 的 `Tools` 自动映射为 MCP `ToolDef` 列表；`tools/call` 通过 bridge 分发，使用 `mcp.tool_name` envelope metadata
- [x] **MCP SSE 端点** — `GET /mcp/{agent_id}` SSE 占位端点，用于服务端主动通知
- [x] **统一网关调用** — `POST /agent/{agent_id}` 从请求体自动检测协议，分发到 A2A/MCP/ACP 桥接 handler
- [x] **协议自动检测** — JSON-RPC `method` 前缀匹配（`message/` `tasks/` → A2A，`tools/` `resources/` `prompts/` `initialize` → MCP），`input`/`agent_name` 字段 → ACP，含 params 形状 fallback
- [x] **多格式发现** — `GET /agent/{agent_id}?format=a2a|mcp|acp` 返回协议专属 Agent Card（A2A AgentCard、MCP server info、ACP Manifest）；默认返回 PeerClaw Card
- [x] **网关指标** — `peerclaw.gateway.requests.total` OpenTelemetry 计数器，含 `protocol` 属性
- [x] **Session 清理** — 后台 goroutine 清理过期 MCP session（1 小时 TTL）

## Phase 15e: ACP Stdio 桥接（已完成）

ndJSON/stdio 桥接，让 ACP 兼容的 Agent（OpenClaw、Zed AI、Coder）通过本地进程通信加入 PeerClaw 网络。

- [x] **`peerclaw acp serve` 命令** — ndJSON/stdio 桥接进程，从 stdin 读取 ACP 请求，向 stdout 写入响应，代理到 PeerClaw server 的 ACP HTTP 桥接
- [x] **ACP 方法分发** — 6 种方法：`create_run`（POST /acp/{agent_id}/runs）、`get_run`（含本地缓存）、`cancel_run`、`list_agents`（目录 → ACP manifest 转换）、`get_agent`、`ping`
- [x] **轻量 ACP 类型** — 在 CLI 包中创建与 server ACP 类型 JSON 兼容的副本（Run、Message、MessagePart、AgentManifest），不依赖 server 模块
- [x] **Run 缓存** — 本地 sync.Map 缓存已创建的 Run，减少 get_run 的服务器往返
- [x] **测试** — 10 个单元测试（ping、非法 JSON、未知方法、create_run、get_run、get_run 缓存、cancel_run、list_agents、get_agent、空行）+ 2 个命令测试

## Phase 16: P2P 文件传输（已完成）

纯点对点大文件传输，端到端加密 — 数据路径零服务器依赖。

- [x] **WebRTC 传输增强** — `CreateDataChannel()`、`RegisterDataChannelHandler()` 用于专用文件传输通道，`Send()` 中添加背压控制（1MB 高水位、256KB 低水位）
- [x] **文件传输消息类型** — `file_offer`、`file_accept`、`file_reject`、`transfer_ready`、`transfer_complete`、`chunk_ack`、`resume_request`、`file_chunk`，定义在 `core/envelope/filetransfer.go`
- [x] **二进制帧协议** — `[seq:4B][length:4B][flags:1B][encrypted_chunk]`，含 FlagData、FlagFIN、FlagACK 标志；默认 64KB 分块
- [x] **传输状态机** — `Idle → Offered → Accepted → Transferring → Completing → Done/Failed/Cancelled`，每个状态有超时控制
- [x] **Challenge-Response 双向鉴权** — 三步 Ed25519 握手：FileOffer(challenge) → FileAccept(challenge_sig, counter_challenge) → TransferReady(counter_sig)
- [x] **流水线推送发送器** — 专用 `ft-{file_id}` DataChannel（有序、可靠），逐块 XChaCha20-Poly1305 加密（AAD = `file_id|seq`），完成后发送 FIN 帧
- [x] **流式接收器** — 二进制帧解码 → 解密 → 写文件，每 100 块发送 ChunkAck，FIN 时校验全文 SHA-256
- [x] **断点续传** — `SaveResumeState()` / `LoadResumeState()` 将最后确认序号持久化到磁盘，`ResumeRequest` 从 `last_seq + 1` 继续
- [x] **Nostr relay 兜底** — WebRTC ICE 失败时，通过标准 `agent.Send()` 路径将加密分块作为 Nostr event 发送（~40KB/event）
- [x] **Mailbox 唤醒** — FileOffer 经由 mailbox 发送时额外发送 `MessageTypeMailboxWakeup`，触发即时 `SyncInbox()` 而非等待轮询
- [x] **Agent 集成** — `SendFile()`、`ListTransfers()`、`GetTransfer()`、`CancelTransfer()` 公开 API；能力 handler 注册为 `"file_transfer"`；`FileTransferDir` 和 `ResumeStatePath` 选项
- [x] **CLI 命令** — `peerclaw send-file --to <id> --file <path>` 含进度轮询；`peerclaw transfer status [--transfer-id <id>]` 传输列表
- [x] **移除 Blob 服务** — 移除中心化的 `server/internal/blob/` 包及所有引用；文件传输现为纯 P2P

## Phase 17: 加固与开发者体验（已完成）

基于安全审计发现的安全加固、错误处理一致性和开发者体验改进。

### 安全加固

- [x] **注册时的持有证明（Proof-of-Possession）** — Agent 在 `POST /agents` 注册时必须用 Ed25519 密钥签署请求体，防止公钥抢注
- [x] **前向保密 / 会话密钥轮换** — 每 1000 条消息或 1 小时后自动进行临时 X25519 密钥轮换；旧密钥材料安全清零
- [x] **加密信任存储** — 信任存储文件使用 XChaCha20-Poly1305 静态加密（密钥通过 HKDF 从 Ed25519 种子派生）；支持从明文透明迁移
- [x] **测试覆盖率目标** — CI 覆盖率门禁：core 80%+，agent/server/cli 70%+

### 错误处理与 API 质量

- [x] **统一错误类型** — `core/errors` 包提供结构化错误码（`not_found`、`validation_error`、`auth_error` 等）和 `New`/`Wrap`/`Is` 辅助函数
- [x] **Card.Validate() 方法** — 验证逻辑从 server 提取到 core 的 `Card.Validate()`；server 委托调用
- [x] **结构化错误响应** — `jsonError` 返回 `{code, message}` 格式，从 HTTP 状态码映射错误码

### CLI 与开发者体验

- [x] **Shell 补全** — `peerclaw completion bash|zsh|fish` 生成包含所有命令、子命令和标志的 shell 补全脚本
- [x] **handleListAgents 扩展** — 新增 `sort`、`search`、`min_score`、`page_size` 查询参数（与公开目录保持一致）

### 声誉透明度

- [x] **公开声誉阈值** — `ReputationLow`（0.3）、`ReputationMedium`（0.7）、`ReputationHigh`（0.8）常量定义在 `core/agentcard/reputation.go`
- [x] **声誉排序与筛选** — `GET /api/v1/agents?sort=reputation&min_score=<float>` 已支持

### 插件生态

插件分级分类：

| 级别 | 插件 | 语言 | 状态 |
|------|------|------|------|
| Tier 1 | openclaw-plugin | TypeScript | 完整测试 + CI |
| Tier 1 | ironclaw-plugin | Rust WASM | 测试 + CI |
| Tier 1 | zeroclaw-plugin | Rust | 测试 + CI（新建） |
| 社区维护 | picoclaw-plugin | Go | 社区维护 |
| 社区维护 | nanobot-plugin | Python | 社区维护 |

- [x] **zeroclaw-plugin** — 新建 Rust crate，通过 bridge WebSocket 协议实现 ZeroClaw 的 Channel trait
- [x] **openclaw-plugin 测试** — vitest 测试套件，覆盖配置 schema 验证和账户解析
- [x] **ironclaw-plugin 测试** — 帧解析、peer ID 提取、序列化单元测试
- [x] **社区插件标记** — picoclaw 和 nanobot 插件标记为社区维护，添加 CONTRIBUTING.md

## Phase 18: Web 仪表盘基础建设 (已完成)

全功能管理仪表盘，支持 8 种语言 i18n。

- [x] **Dashboard SPA** — React + TypeScript + Tailwind + shadcn/ui，嵌入服务端二进制
- [x] **管理面板** — 总览、用户、Agent、举报、分类、分析、调用记录页面
- [x] **供应商控制台** — 仪表盘、Agent 管理、发现、调用历史、API Key、通知、个人资料
- [x] **公共页面** — 首页、目录、Agent 详情、Playground、关于、登录/注册/找回密码
- [x] **国际化** — 8 种语言 (en, zh, es, fr, ja, pt, ru, ar) 完整翻译覆盖
- [x] **数据导出** — 管理表格 CSV/JSON 导出
- [x] **可排序表格** — 列排序 + 加载骨架屏
- [x] **无障碍** — 跳至内容、ARIA 标签、键盘导航

## Phase 19-22: 仪表盘增强 (已完成)

### Phase 19: 批量操作
- [x] **批量端点** — `POST /admin/{agents,reports,users}/bulk`
- [x] **SelectableTable** — 可复用的复选框选择 + 浮动操作栏组件
- [x] **管理页面** — AgentsPage (验证/删除)、UsersPage (删除)、ReportsPage (审核/驳回/删除)

### Phase 20: 个人资料与认证改进
- [x] **密码强度** — 4 段可视化条 + 需求清单，应用于注册/个人资料/找回密码
- [x] **自动隐藏** — ProfilePage 成功消息 3 秒后自动消失
- [x] **关于页面** — 更新路线图至 Phase 4，添加 GitHub 路线图链接

### Phase 21: 管理员审计日志
- [x] **adminaudit 包** — Store/SQLite/Postgres/Service/Factory 模式
- [x] **审计记录** — 所有管理操作记录 (user.delete, agent.verify, report.update 等)
- [x] **AuditLogPage** — 过滤器 (管理员、操作、目标类型、日期范围) + 分页 + 侧栏导航

### Phase 22: 高级仪表盘
- [x] **可点击统计卡片** — 总览卡片点击导航到对应管理页面
- [x] **最近活动** — 总览页显示最近 10 条审计事件
- [x] **供应商时间范围** — 7天/30天/全部选择器，后端 `since` 参数支持

## Phase 23-26: 安全加固与代码质量 (已完成)

基于 2026-03-28 全项目代码审计报告的系统性修复。

### Phase 23: 关键安全加固
- [x] **SEC-C01** — 清理 90+ 处错误信息泄露，零 500 级泄露
- [x] **SEC-C02** — 修复 X-Forwarded-For IP 欺骗，使用 `BridgeClientIP()` 代理信任验证
- [x] **SEC-H04** — CLI 配置文件权限从 0644 限制为 0600
- [x] **SEC-M06** — `Verify()` 添加 nil 公钥检查防止 panic

### Phase 24: 认证与访问控制加固
- [x] **SEC-H01** — 关闭 Agent ID 认证绕过 (pubKey 回退时 AgentExists 注册检查)
- [x] **SEC-H03** — owner_user_id IDOR 防护 (handleRegister 中剥离元数据)
- [x] **SEC-M01** — OTP 从 6 位增强至 8 位 + 暴力破解锁定 (5 次失败 → 15 分钟锁定)
- [x] **SEC-M03** — `decodeJSON()` 辅助函数，认证端点 64KB 限制

### Phase 25: 风险代码与稳定性
- [x] **RISK-01** — A2A Bridge TOCTOU 竞态修复 (原子 CompareAndSwap)
- [x] **RISK-02** — WebSocket 读循环 panic 恢复 (defer/recover)
- [x] **SEC-M05** — React ErrorBoundary 包裹所有 3 个路由组

### Phase 26: 代码质量与 DRY
- [x] **DUP-05** — 共享状态 Badge 工具函数 (`lib/status.ts`)
- [x] **DUP-04** — `useAsyncAction` Hook 统一异步错误处理
- [x] **DEAD-04** — Envelope `WithTTL()` 和 `GenerateNonce()` 方法实现
- [x] **性能** — Admin 路由 React.lazy 懒加载 (6 页面代码分割，主包减小 47KB)
