# Web App Roadmap — Post v0.8.0

> Based on comprehensive audit of all 27 pages (Console 10, Admin 8, Public 9) and 5 layout/shared components.
> Last updated: 2026-03-14

---

## Overview

v0.8.0 addressed Admin Agents pagination, Analytics charts, Overview trends, Footer/404, shared hooks, and form protection. This roadmap covers all remaining findings, prioritized by user impact and implementation complexity.

---

## Phase 8: Mobile Navigation (P0)

**Problem**: PublicLayout, ConsoleLayout, and AppLayout have no mobile navigation. Sidebar and nav links overflow on small screens. PlaygroundPage is entirely desktop-only (fixed `w-72` sidebar).

### 8a. PublicLayout hamburger menu
**File**: `components/public/PublicLayout.tsx`

- Add responsive hamburger button (visible < `md`)
- Slide-in drawer with nav links, user menu, language switcher
- Close on route change via `useLocation`

### 8b. ConsoleLayout collapsible sidebar
**File**: `components/layout/ConsoleLayout.tsx`

- Add sidebar collapse toggle (icon rail mode on desktop)
- On mobile: hidden by default, slide-in drawer on hamburger tap
- Persist collapse preference to `localStorage`

### 8c. PlaygroundPage responsive layout
**File**: `pages/PlaygroundPage.tsx`

- Replace fixed `w-72` sidebar with collapsible panel
- On mobile: agent selector as top sheet, chat fills viewport
- Bottom input bar stays fixed on mobile

### 8d. AppLayout mobile sidebar
**File**: `components/layout/AppLayout.tsx`

- Mirror ConsoleLayout mobile pattern for admin

---

## Phase 9: Toast Notification System (P0)

**Problem**: All mutation errors use `alert()` or `window.confirm()`. No success feedback on saves/deletes.

### 9a. Add toast provider
**New**: `components/ui/toaster.tsx` (shadcn/ui Sonner integration)

- Install `sonner` package
- Wrap App with `<Toaster />` provider
- Export `toast()` function

### 9b. Replace alert() across all pages
**Files**: AgentsPage, UsersPage, ReportsPage, CategoriesPage, APIKeysPage, InvocationHistoryPage, ProviderAgentDetailPage

- Replace `alert(error.message)` with `toast.error(message)`
- Add `toast.success()` on mutations (delete, verify, role change, etc.)

### 9c. Replace window.confirm() with dialog
**Files**: ProviderAgentDetailPage, APIKeysPage

- Replace `window.confirm()` with `AlertDialog` (shadcn/ui)
- Consistent destructive action confirmation pattern

---

## Phase 10: Table Sorting (P1)

**Problem**: No admin table supports column sorting. Backend list endpoints lack sort parameters.

### 10a. Backend sort support
**File**: `internal/server/admin_handler.go`

Add `sort` and `order` query params to:
- `handleAdminListAgents` — sort by name, status, registered_at, last_heartbeat
- `handleAdminListUsers` — sort by email, role, created_at
- `handleAdminListReports` — sort by created_at, status
- `handleAdminListInvocations` — sort by duration_ms, created_at, status_code

### 10b. Sortable table header component
**New**: `components/ui/sortable-header.tsx`

- Reusable `<SortableHeader field="name" label={t('...')} />` component
- Arrow indicators for sort direction
- Click toggles asc/desc/none

### 10c. Apply sorting to all admin tables
**Files**: AgentsPage, UsersPage, ReportsPage, InvocationsPage

- Add sort state + URL param sync
- Update hooks to pass sort/order params

---

## Phase 11: Loading Skeletons (P1)

**Problem**: All pages show plain "Loading..." text. No skeleton placeholders for cards, tables, or charts.

### 11a. Skeleton components
**New**: `components/ui/skeleton.tsx` (likely already in shadcn/ui)

- `<TableSkeleton rows={5} cols={4} />`
- `<CardSkeleton />` for stat cards
- `<ChartSkeleton />` for Recharts placeholders

### 11b. Apply across admin pages
**Files**: OverviewPage, AgentsPage, UsersPage, AnalyticsPage, InvocationsPage, ReportsPage

### 11c. Apply across console pages
**Files**: ProviderDashboardPage, ProviderAgentsPage, ProviderAgentDetailPage, InvocationHistoryPage

### 11d. Apply across public pages
**Files**: DirectoryPage, PublicProfilePage, PlaygroundPage

---

## Phase 12: Data Export (P1)

**Problem**: No page supports CSV/JSON export. Admins need to extract data for reporting.

### 12a. Export utility
**New**: `lib/export.ts`

- `exportToCSV(data, columns, filename)` — generate and download CSV
- `exportToJSON(data, filename)` — download JSON file

### 12b. Add export buttons
**Files**: AgentsPage, UsersPage, InvocationsPage, AnalyticsPage

- "Export CSV" button in page header
- Uses current filter state for export scope

---

## Phase 13: Accessibility (P1)

**Problem**: Missing aria-labels on icon buttons, no skip-to-content links, no prefers-reduced-motion, no live regions for dynamic content.

### 13a. Icon button accessibility
**Files**: All pages with icon-only buttons

- Add `aria-label` to all `<Button>` with only icon children
- Applies to: copy, delete, verify, edit buttons

### 13b. Skip-to-content links
**Files**: PublicLayout, ConsoleLayout, AppLayout

- Add visually hidden "Skip to main content" link as first focusable element

### 13c. prefers-reduced-motion
**Files**: LandingPage (animated orbs), AboutPage (pulse animation), loading spinners

- Wrap animations in `@media (prefers-reduced-motion: no-preference)`
- Or use Tailwind's `motion-safe:` prefix

### 13d. ARIA live regions
**Files**: All pages with dynamic content

- Add `role="alert"` to error messages
- Add `aria-live="polite"` to loading states
- Add `aria-current="page"` to active nav links

### 13e. Form accessibility
**Files**: AgentEditPage, ProfilePage, RegisterPage, LoginPage

- Add `aria-describedby` linking help text to inputs
- Add `aria-invalid` on validation errors
- Add `aria-label` to form elements

---

## Phase 14: Provider Agents Enhancement (P1)

**Problem**: ProviderAgentsPage has no search, no sort, no pagination, and no filtering.

### 14a. Add search + filter bar
**File**: `pages/ProviderAgentsPage.tsx`

- Search by name
- Filter by status (online/offline/degraded)
- Filter by protocol

### 14b. Add pagination
Same file — match admin AgentsPage pattern (PAGE_SIZE=20, prev/next buttons)

### 14c. Improve empty state
- Add guidance text for new users
- Link to registration wizard
- Show claim token section inline

---

## Phase 15: Invocation Enhancements (P2)

**Problem**: Invocation pages (admin + console) lack date range filter, detail view, and status breakdown.

### 15a. Date range filter
**Files**: `pages/admin/InvocationsPage.tsx`, `pages/InvocationHistoryPage.tsx`

- Add date range selector (today, 7d, 30d, custom)
- Backend already supports `since` param on analytics, extend to invocations list

### 15b. Invocation detail modal
**Files**: Same pages

- Click on row to open modal with full error message, request/response metadata
- Backend: add `GET /api/v1/admin/invocations/{id}` if not exists

### 15c. Status summary bar
**Files**: Same pages

- Show mini stat bar above table: "150 total, 142 success, 8 errors, avg 45ms"

---

## Phase 16: Copy-to-Clipboard & Text Handling (P2)

**Problem**: Agent IDs, public keys, and endpoint URLs truncate without copy buttons throughout admin/public pages.

### 16a. CopyButton component
**New**: `components/ui/copy-button.tsx`

- Reusable inline copy button with success feedback (checkmark icon, 2s timeout)
- Props: `value: string`, `size?: "sm" | "default"`

### 16b. Apply across pages
**Files**: AdminAgentDetailPage, PublicProfilePage, InvocationsPage, AgentsPage

- Agent ID fields get copy button
- Public key fields get copy button
- Endpoint URL fields get copy button

---

## Phase 17: Admin Reports Search (P2)

**Problem**: Reports page has status tabs but no text search.

### 17a. Backend search support
**File**: `internal/server/admin_handler.go` — `handleAdminListReports()`

Add `search` query param for matching against `reason`, `target_id`, `reporter_id`.

### 17b. Frontend search input
**File**: `pages/admin/ReportsPage.tsx`

- Add search input above status tabs
- Debounce 300ms, reset page on change

---

## Phase 18: Confirmation Dialogs (P2)

**Problem**: Several pages use inline confirm/cancel button patterns. Agent detail uses `window.confirm()`.

### 18a. AlertDialog pattern
**Using**: shadcn/ui `AlertDialog`

Standardize destructive action confirmations across:
- AdminAgentDetailPage (delete agent, verify toggle)
- APIKeysPage (revoke key)
- ProviderAgentDetailPage (already has window.confirm — upgrade)
- CategoriesPage (delete category)

---

## Phase 19: Bulk Actions (P2)

**Problem**: No multi-select or bulk operations on any table.

### 19a. Selectable table component
**New**: `components/ui/selectable-table.tsx`

- Checkbox column with select all
- Floating action bar when items selected
- "N items selected" indicator

### 19b. Backend bulk endpoints
**File**: `internal/server/admin_handler.go`

- `POST /api/v1/admin/agents/bulk` — bulk verify/delete
- `PUT /api/v1/admin/reports/bulk` — bulk status update

### 19c. Apply to admin tables
**Files**: AgentsPage, UsersPage, ReportsPage

---

## Phase 20: Profile & Auth Improvements (P3)

### 20a. Password strength indicator
**Files**: RegisterPage, ProfilePage, ForgotPasswordPage

- Show strength bar (weak/fair/strong) as user types
- Show password requirements checklist

### 20b. Profile page improvements
**File**: `pages/ProfilePage.tsx`

- Add `role="alert"` on success/error messages
- Auto-dismiss success messages after 3s
- Show last login timestamp if available

### 20c. About page roadmap update
**File**: `pages/AboutPage.tsx`

- Update phases to reflect actual 16+ phases completed
- Add link to full docs/ROADMAP.md

---

## Phase 21: Admin Audit Log (P3)

### 21a. Backend audit logging
**New**: `internal/audit/` package

- Log admin actions (delete user, verify agent, update report, etc.)
- Store: action, admin_user_id, target_type, target_id, metadata, timestamp

### 21b. Admin Audit Log page
**New**: `pages/admin/AuditLogPage.tsx`

- Table with filters: admin user, action type, date range
- Add route `/admin/audit`

---

## Phase 22: Advanced Dashboard (P3)

### 22a. Overview quick links
**File**: `pages/OverviewPage.tsx`

- Each stat card clickable → navigates to respective admin page
- Add "Recent Activity" feed (last 10 admin actions)

### 22b. Provider dashboard trends
**File**: `pages/ProviderDashboardPage.tsx`

- Add time range selector (7d/30d)
- Add success rate trend sparkline per agent
- Add comparison with previous period

---

## Implementation Priority Summary

| Priority | Phases | Estimated Scope |
|----------|--------|----------------|
| **P0** — Must-have | 8 (Mobile nav), 9 (Toast) | ~20 files |
| **P1** — Important | 10 (Sorting), 11 (Skeletons), 12 (Export), 13 (a11y), 14 (Provider) | ~30 files |
| **P2** — Nice-to-have | 15-19 (Invocations, Copy, Reports, Dialogs, Bulk) | ~20 files |
| **P3** — Future | 20-22 (Auth, Audit, Dashboard) | ~10 files |

---

## Files Impact Matrix

| File | Phases |
|------|--------|
| `components/public/PublicLayout.tsx` | 8a, 13b |
| `components/layout/ConsoleLayout.tsx` | 8b, 13b |
| `pages/PlaygroundPage.tsx` | 8c |
| `pages/AgentsPage.tsx` | 9b, 10c, 11b, 12b, 13a, 19c |
| `pages/admin/UsersPage.tsx` | 9b, 10c, 11b, 12b, 19c |
| `pages/admin/ReportsPage.tsx` | 9b, 10c, 11b, 17b, 19c |
| `pages/admin/InvocationsPage.tsx` | 9b, 10c, 11b, 12b, 15a-c |
| `pages/admin/AnalyticsPage.tsx` | 11b, 12b |
| `pages/OverviewPage.tsx` | 11b, 22a |
| `pages/ProviderAgentsPage.tsx` | 9b, 11c, 14a-c |
| `pages/ProviderAgentDetailPage.tsx` | 9b, 9c, 11c, 16b, 18a |
| `pages/AgentEditPage.tsx` | 13e |
| `pages/ProfilePage.tsx` | 13e, 20b |
| `pages/LandingPage.tsx` | 13c |
| `pages/AboutPage.tsx` | 13c, 20c |
| `pages/DirectoryPage.tsx` | 11d |
| `pages/PublicProfilePage.tsx` | 11d, 16b |
| `internal/server/admin_handler.go` | 10a, 15a, 17a, 19b |
