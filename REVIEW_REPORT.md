# PeerClaw 全项目代码审计报告 (二次审计)

**审计日期**: 2026-03-28
**审计范围**: peerclaw-core v0.8.0, peerclaw-server v0.11.0, peerclaw-agent v0.7.3, peerclaw-cli v0.9.1, 5 个平台插件, React Dashboard
**审计类型**: 安全修复后复审 + 产品价值分析 + DDD 合规深度评估
**前次审计**: 2026-03-28 (Phase 23-26 修复已应用)

---

## 目录

1. [修复验证](#1-修复验证)
2. [新发现漏洞](#2-新发现漏洞)
3. [风险代码](#3-风险代码)
4. [死代码与无效代码](#4-死代码与无效代码)
5. [重复代码](#5-重复代码)
6. [模块集成度](#6-模块集成度)
7. [DDD 合规性深度评估](#7-ddd-合规性深度评估)
8. [产品设计与价值分析](#8-产品设计与价值分析)
9. [可优化项与精简建议](#9-可优化项与精简建议)
10. [各模块质量评分](#10-各模块质量评分)
11. [修复优先级](#11-修复优先级)

---

## 1. 修复验证

### Phase 23-26 修复已确认生效

| 修复项 | 状态 | 验证方式 |
|--------|:----:|---------|
| SEC-C01 错误信息泄露 | ✅ | grep 确认 0 个 500 级泄露，15 个残留为安全的 auth 消息 |
| SEC-C02 X-Forwarded-For | ✅ | user_auth_handler.go 已使用 BridgeClientIP() |
| SEC-H01 Agent ID 绕过 | ✅ | auth.go 使用 AgentExists 回调，http.go 已注入 |
| SEC-H03 IDOR | ✅ | handleRegister 中 `delete(req.Metadata, "owner_user_id")` |
| SEC-H04 CLI 权限 | ✅ | config.go 使用 0600 |
| SEC-M01 OTP 8 位 | ✅ | generateOTP 4 字节，OTPInput 默认 8 |
| SEC-M03 Body 限制 | ✅ | decodeJSON + maxAuthBodyBytes 64KB |
| SEC-M06 Verify nil | ✅ | sign.go 顶部 pubKey nil 检查 |
| RISK-01 A2A 竞态 | ✅ | CompareAndSwap 原子循环 |
| RISK-02 WS panic | ✅ | defer/recover 在 read loop 内 |
| SEC-M05 ErrorBoundary | ✅ | 3 个路由组均已包裹 |
| DEAD-04 Envelope 方法 | ✅ | WithTTL + GenerateNonce 已实现 |
| OTP 暴力防护 | ✅ | isOTPBlocked/recordOTPFailure/clearOTPFailures 已集成 |
| Lazy Loading | ✅ | 6 个 admin 页面 React.lazy 代码分割 |

---

## 2. 新发现漏洞

### 2.1 CRITICAL

#### NEW-C01: 批量操作无审计记录

| 属性 | 值 |
|------|------|
| 严重性 | **CRITICAL** |
| 文件 | server/internal/server/admin_handler.go |
| 行号 | handleAdminBulkAgents, handleAdminBulkReports, handleAdminBulkUsers |

**问题**: Phase 19 添加的 3 个批量操作 handler 未调用 `recordAdminAudit()`。管理员可通过批量操作执行破坏性操作（删除用户/Agent/举报）**而不留下任何审计轨迹**。

对比：单条操作 handler 均有审计（如 `handleAdminDeleteAgent:266` 调用 `recordAdminAudit()`），但批量 handler 完全遗漏。

---

#### NEW-C02: Agent 回调使用 context.Background()

| 属性 | 值 |
|------|------|
| 严重性 | **CRITICAL** |
| 文件 | agent/agent.go:432, 608, 631, 640 |

**问题**: Agent SDK 的回调函数使用 `context.Background()` 而非 Agent 生命周期 context。当 `Agent.Stop()` 被调用时，回调中的长时间操作（网络请求、文件 I/O）无法被取消，导致进程无法正常退出。

```go
// agent.go:432
a.HandleIncomingEnvelope(context.Background(), &env)  // 无法取消
```

---

#### NEW-C03: 文件传输 LastConfirmedSeq TOCTOU 竞态

| 属性 | 值 |
|------|------|
| 严重性 | **CRITICAL** |
| 文件 | agent/filetransfer/sender.go:79,112 |

**问题**: Sender 读取 `transfer.LastConfirmedSeq` 时无同步，另一 goroutine (manager.go:240) 可同时更新该字段。读取到过期值后 Seek 到错误的文件偏移量，导致数据损坏。

---

#### NEW-C04: Nostr 重连 goto 导致资源泄露

| 属性 | 值 |
|------|------|
| 严重性 | **CRITICAL** |
| 文件 | agent/transport/nostr.go:216-246 |

**问题**: `subscribeLoop()` 使用 `goto reconnect` 跳过 `sub.Unsub()` 清理，导致订阅对象泄露。长时间运行的 Agent 会累积未释放的 Nostr 订阅。

---

### 2.2 HIGH

#### NEW-H01: 服务端 Goroutine 缺少 WaitGroup 追踪

| 属性 | 值 |
|------|------|
| 严重性 | **HIGH** |
| 文件 | server/cmd/peerclawd/main.go:478,488,532,560 |

**问题**: 5 个后台 goroutine 中仅 1 个有 `wg.Add(1)` 追踪。SIGINT/SIGTERM 时其余 4 个 goroutine 可能在数据操作中途被强制终止。

---

#### NEW-H02: X25519PrivateKey() 缺少 nil 检查

| 属性 | 值 |
|------|------|
| 严重性 | **HIGH** |
| 文件 | core/identity/x25519.go:23 |

**问题**: `kp.PrivateKey.Seed()` 在 kp 或 kp.PrivateKey 为 nil 时会 panic。

---

#### NEW-H03: Trust Store 加密降级为明文

| 属性 | 值 |
|------|------|
| 严重性 | **HIGH** |
| 文件 | agent/security/trust.go:353-376 |

**问题**: 当加密密钥已配置但 trust store 文件不以 magic 前缀开头时，代码 fallback 到明文 JSON 解析，绕过加密保护。

---

### 2.3 MEDIUM

| ID | 描述 | 文件 |
|----|------|------|
| NEW-M01 | Provider Dashboard 聚合统计忽略 `since` 参数 | server provider_handler.go:503 |
| NEW-M02 | 前端 ~15 个硬编码英文字符串未 i18n | OverviewPage, AgentTable, EventFeed, AuditLogPage |
| NEW-M03 | Agent 模块 50+ 处 `_ = err` 静默错误 | agent/agent.go:608,631,640 |
| NEW-M04 | 各插件 WebSocket 读取限制不统一 (256KB/512KB/无限制) | picoclaw, openclaw, ironclaw |

---

## 3. 风险代码

| ID | 类型 | 位置 | 严重性 |
|----|------|------|--------|
| NEW-C02 | Context 泄露 | agent/agent.go 回调 | CRITICAL |
| NEW-C03 | 数据竞态 | filetransfer/sender.go | CRITICAL |
| NEW-C04 | 资源泄露 | transport/nostr.go | CRITICAL |
| NEW-H01 | 进程退出不安全 | cmd/peerclawd/main.go | HIGH |
| NEW-M03 | 静默失败 | agent 模块全局 | MEDIUM |

---

## 4. 死代码与无效代码

### 已确认死代码

| ID | 位置 | 描述 | 建议 |
|----|------|------|------|
| DEAD-01 | core/agentcard/card.go | `HasCapability()`, `HasSkill()`, `HasTool()`, `SupportsProtocol()` — 4 个方法在整个代码库中从未被调用 | 删除或添加调用方 |
| DEAD-02 | core/api/ | 目录存在但仅含 proto 子目录，无 Go 代码 | 保留 (protobuf 定义) |

### 已确认非死代码（上次报告误判）

| ID | 位置 | 说明 |
|----|------|------|
| ~~DEAD-01~~ | core/identity/keypair.go | `SaveKeypair()` / `LoadKeypair()` 在 agent 初始化和 CLI claim 中使用 — **非死代码** |

---

## 5. 重复代码

### 5.1 前端重复 (仍然存在)

| 模式 | 影响范围 | 行数浪费 | 状态 |
|------|---------|---------|------|
| 分页逻辑 (useState+useEffect) | 7 个页面 | ~200 行 | 未修复 |
| 错误处理回调 (try/catch/toast) | 10+ 处 | ~150 行 | 部分修复 (useAsyncAction 已创建但未迁移) |
| 时间格式化 ("just now", "m ago") | 3 个文件 | ~40 行 | 未修复 |
| 状态 Badge 映射 | 4+ 页面 | ~50 行 | 部分修复 (lib/status.ts 已创建但未迁移) |

### 5.2 后端重复 (仍然存在)

| 模式 | 影响范围 | 状态 |
|------|---------|------|
| 服务可用性检查 `if s.xxx == nil` | 42+ 处 | 未修复 |
| Store 五件套样板代码 | 9 个模块 | 结构性问题，建议长期迁移 |

### 5.3 前端分页索引不一致

InvocationHistoryPage 使用 1-indexed `page`，而 UsersPage/ReportsPage 使用 0-indexed `page`。统一 hook 可解决此问题。

---

## 6. 模块集成度

### 6.1 版本对齐 ✅

| 模块 | go.mod 中 core 版本 | go.mod 中 agent 版本 |
|------|:---:|:---:|
| server | v0.8.0 ✅ | — |
| agent | v0.8.0 ✅ | — |
| cli | v0.8.0 ✅ | v0.7.3 ✅ |

### 6.2 前后端类型同步 ✅

AdminAuditEvent 8 个字段前后端完全匹配。BulkActionResponse 匹配。所有 API 类型已同步。

### 6.3 集成缺陷

| ID | 问题 | 影响 |
|----|------|------|
| INT-01 | Provider Dashboard `since` 参数仅影响 per-agent 统计，不影响聚合统计 | 前端时间筛选器功能不完整 |
| INT-02 | 批量操作无审计记录 | 审计日志不完整 |
| INT-03 | 前后端类型手动同步 | 长期维护风险 |

---

## 7. DDD 合规性深度评估

### 7.1 Service 层定量分析

对 3 个典型 Service 进行方法级分析：

| Service | 总方法数 | 纯 Store 委托 | 含业务逻辑 | 委托率 |
|---------|:-------:|:------------:|:---------:|:-----:|
| invocation/service.go | 11 | **11** | 0 | **100%** |
| adminaudit/service.go | 2 | 1 | 1 | 50% |
| contacts/service.go | 5 | 2 | 3 | 40% |
| **加权平均** | **18** | **14** | **4** | **78%** |

**结论**: 78% 的 Service 方法是纯 Store 代理 — 经典贫血模型反模式。

### 7.2 Handler 层业务逻辑量化

`handleAdminDashboard` 分析：
- 总行数: 63 行
- 业务逻辑: 35 行 (55%) — 聚合统计、健康检查、趋势计算
- HTTP 处理: 28 行 (45%) — 参数解析、响应序列化

**应该**: `resp := s.adminService.GetDashboard(ctx)` — 1 行调用。

### 7.3 领域事件: 完全缺失

搜索整个代码库未发现任何 event bus / publish-subscribe 模式。所有副作用（审计记录、声誉更新、通知）通过同步直接调用实现。

### 7.4 DDD 合规评分 (更新)

| 原则 | 上次 | 本次 | 变化 | 说明 |
|------|:----:|:----:|:----:|------|
| 领域模型丰富性 | 4 | 4 | → | Service 层仍为贫血模型 |
| 限界上下文划分 | 7 | 7 | → | 模块边界清晰 |
| 仓储模式 | 8 | 8 | → | Store 接口 + 多数据库实现 |
| Application Service | 3 | 3 | → | 78% 纯委托 |
| 领域事件 | 1 | 1 | → | 完全缺失 |
| 聚合根 | 2 | 2 | → | 无定义 |
| 值对象 | 7 | 7.5 | ↑ | WithTTL/GenerateNonce 增强 |
| 反腐层 | 6 | 6 | → | Plugin Adapter 设计合理 |

**DDD 综合评分: 4.8/10** — 仓储模式执行良好，但缺乏领域驱动的核心设计。

### 7.5 DDD 改进建议

**短期可行** (不改变架构):
1. 将 `handleAdminDashboard` 的聚合逻辑抽取到 `AdminDashboardService.GetOverview()`
2. 批量操作 handler 中的循环逻辑移入 Service 方法
3. 移除纯委托 Service（如 invocation），Handler 直接调用 Store

**中期改进**:
4. 引入简单事件总线（`EventBus` + subscriber 模式）用于审计、通知
5. 定义 Agent 聚合根，将 Registry + Reputation + Review 操作统一

---

## 8. 产品设计与价值分析

### 8.1 产品定位

**声明**: "The Identity & Trust Layer for AI Agents"
**核心价值**:

| 价值点 | 实现程度 | 说明 |
|--------|:-------:|------|
| Ed25519 身份 | ✅ 完整 | core/identity/, 签名/验证全链路 |
| EWMA 声誉 | ⚠️ 部分 | 仅在 server 端计算，Agent SDK 运行时无法获取 |
| E2E 加密通信 | ✅ 完整 | XChaCha20-Poly1305 + X25519 ECDH |
| 多协议桥接 | ✅ 完整 | A2A/MCP/ACP + Envelope 统一格式 |
| P2P 传输 | ✅ 完整 | WebRTC + Nostr 降级 |
| Web 平台 | ✅ 完整 | 8 语言 i18n, 管理面板, 供应商控制台 |

### 8.2 产品差异化分析

| 竞品方向 | PeerClaw 优势 | 风险 |
|---------|-------------|------|
| MCP 生态 (Anthropic) | PeerClaw 支持 MCP + A2A + ACP，而 MCP 仅支持工具调用 | 若 MCP 成为行业标准，多协议优势削弱 |
| A2A 生态 (Google) | PeerClaw 加入了 P2P 加密和声誉，A2A 原生不含 | A2A 可能独立实现信任层 |
| 中心化平台 | PeerClaw 支持联邦和 P2P，不依赖单一服务商 | 去中心化增加运维复杂度 |

### 8.3 可精简的模块

| 模块 | 当前状态 | 建议 |
|------|---------|------|
| ACP Bridge | 完整实现 | 保留 — IBM ACP 生态有一定采用 |
| Nostr 传输 | 完整实现 | 保留 — 核心差异化功能 |
| 联邦 (Federation) | 基础实现 | **考虑降级** — 无实际用户使用案例，增加代码复杂度 |
| 文件传输 | 完整但有竞态 | 修复后保留 — 独特能力 |
| picoclaw/nanobot 插件 | 社区维护 | **可归档** — 长期无活跃贡献 |

### 8.4 产品优化建议

1. **声誉同步**: Agent SDK 应能定期从 Server 同步信任分数，用于运行时信任决策
2. **Claim Token 流程**: 已经非常流畅（一键复制 prompt），建议增加 QR 码扫描注册
3. **Dashboard 分析**: 增加时间趋势图（调用量、成功率）— 当前仅有静态统计
4. **SDK 文档**: core 包缺少 godoc 示例，新开发者上手困难

---

## 9. 可优化项与精简建议

### 9.1 立即可优化

| 项目 | 预估收益 | 工作量 |
|------|---------|--------|
| 批量 handler 添加审计记录 | 审计完整性 | 0.5h |
| Provider Dashboard since 参数修复 | 数据一致性 | 1h |
| X25519 nil 检查 | 防止 panic | 0.5h |
| 硬编码字符串 i18n 化 | 完整国际化 | 2h |

### 9.2 结构性优化

| 项目 | 描述 | ROI |
|------|------|-----|
| 抽取前端 usePagination Hook | 消除 7 页面重复 | 高 |
| 迁移 lib/status.ts 到使用方 | 消除 4+ 页面 Badge 映射 | 中 |
| 简化 Service 层 | 移除纯委托 Service 或充实业务逻辑 | 中 |
| 前后端类型自动生成 | 消除手动同步风险 | 高 (长期) |

### 9.3 架构性优化 (长期)

| 项目 | 描述 |
|------|------|
| Token 迁移 httpOnly Cookie | 解决 XSS 窃取风险 |
| 事件总线 | 解耦审计/通知/声誉更新 |
| Agent 聚合根 | 统一 Agent 生命周期管理 |
| 声誉同步 API | Agent 运行时获取信任分数 |

---

## 10. 各模块质量评分

| 模块 | 安全 | 质量 | 测试 | DDD | 综合 | 评级 | 趋势 |
|------|:----:|:----:|:----:|:---:|:----:|:----:|:----:|
| **peerclaw-core** | 8 | 8 | 7 | 7.5 | **7.6** | B+ | ↑ |
| **peerclaw-server (handlers)** | 7 | 7 | 6 | 4 | **6.0** | B- | ↑ |
| **peerclaw-server (stores)** | 8 | 8 | 7 | 8 | **7.8** | B+ | → |
| **peerclaw-server (middleware)** | 8 | 8 | 6 | — | **7.3** | B+ | ↑ |
| **peerclaw-agent** | 5 | 6 | 6 | 6 | **5.8** | C+ | ↓ |
| **peerclaw-cli** | 7 | 7 | 5 | — | **6.3** | B- | ↑ |
| **openclaw-plugin** | 7 | 8 | 5 | — | **6.7** | B- | → |
| **ironclaw-plugin** | 5 | 6 | 4 | — | **5.0** | C | → |
| **nanobot-plugin** | 7 | 7 | 5 | — | **6.3** | B- | → |
| **picoclaw-plugin** | 7 | 7 | 5 | — | **6.3** | B- | → |
| **Dashboard (React)** | 5.5 | 7 | 3 | — | **5.2** | C+ | ↑ |
| **adminaudit** | 8 | 8 | 4 | 7 | **6.8** | B | → |

**综合评分: 6.4/10** (上次 6.8 → 本次 6.4，因新发现 Agent 模块关键问题拉低)

### 评级变化说明

- **core ↑ 7.3→7.6**: WithTTL, GenerateNonce, Verify nil check 修复
- **server handlers ↑ 5.5→6.0**: 错误泄露修复，auth 加固，但批量审计遗漏
- **server middleware ↑ 7.0→7.3**: BridgeClientIP 统一使用
- **agent ↓ 6.5→5.8**: 新发现 context 泄露、竞态条件、goto 资源泄露（此前未深入审计）
- **cli ↑ 6.0→6.3**: 权限修复
- **Dashboard ↑ 5.0→5.2**: ErrorBoundary, lazy loading, 但 token 存储仍未修复

---

## 11. 修复优先级

### P0 — 立即修复

| # | 问题 | 模块 | 工作量 |
|---|------|------|--------|
| 1 | NEW-C01: 批量操作添加审计记录 | server | 0.5h |
| 2 | NEW-M01: Provider Dashboard since 参数修复 | server | 1h |
| 3 | NEW-H02: X25519 nil 检查 | core | 0.5h |

### P1 — 高优先级

| # | 问题 | 模块 | 工作量 |
|---|------|------|--------|
| 4 | NEW-C02: Agent 回调 Context 替换 | agent | 4h |
| 5 | NEW-C03: 文件传输 TOCTOU 修复 | agent | 2h |
| 6 | NEW-C04: Nostr goto 重构 | agent | 4h |
| 7 | NEW-H01: main.go WaitGroup 补全 | server | 1h |
| 8 | NEW-H03: Trust Store 严格模式 | agent | 2h |

### P2 — 中优先级

| # | 问题 | 模块 | 工作量 |
|---|------|------|--------|
| 9 | NEW-M02: 硬编码字符串 i18n | frontend | 2h |
| 10 | NEW-M03: Agent 静默错误改为 warn 日志 | agent | 4h |
| 11 | 前端分页 Hook 抽取 | frontend | 3h |
| 12 | Service 层业务逻辑充实 | server | 8h |

### P3 — 长期改进

| # | 问题 | 模块 | 工作量 |
|---|------|------|--------|
| 13 | Token 迁移 httpOnly Cookie | server + frontend | 16h |
| 14 | 事件总线引入 | server | 16h |
| 15 | Agent 聚合根定义 | core + server | 8h |

---

## 附录: 审计统计

| 指标 | 上次 | 本次 |
|------|:----:|:----:|
| CRITICAL 漏洞 | 3 | **4** (新发现，原 3 个已修复) |
| HIGH 漏洞 | 5 | **3** (新发现，原 5 个已修复) |
| MEDIUM 漏洞 | 8 | **4** (新发现，原 8 个已修复) |
| 已修复项 | 0 | **14** |
| 死代码项 | 7 | **1** (6 个已修复或确认非死代码) |
| 重复代码模式 | 5 | **5** (部分创建了工具但未迁移) |
| DDD 综合评分 | 6.0 | **4.8** (深入评估后下调) |
| 综合评分 | 6.8 | **6.4** (Agent 模块问题拉低) |

**关键改进**: Server 模块安全评分从 5 分升至 7-8 分。
**新发现重点**: Agent SDK 模块存在 context 管理、竞态条件、资源泄露等系统性问题，需要专项修复。

---

*本报告为 Phase 23-26 安全修复后的复审报告。建议对 Agent SDK 模块进行专项 Sprint 修复 (P1 #4-#8)。*
