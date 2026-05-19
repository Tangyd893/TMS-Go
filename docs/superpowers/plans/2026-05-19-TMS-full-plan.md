# TMS 运输管理系统完整实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 从当前项目骨架阶段，按里程碑逐步完成 TMS 运输管理系统的全栈开发，形成完整的运输业务闭环。

**Architecture:** 后端采用标准库 `net/http` 模块化架构（已确立），前端 Vue 3 + Element Plus + Pinia，数据层 PostgreSQL + Redis + RabbitMQ + MinIO。后端按模块组织（`internal/modules/{module}/`），每个模块包含 api/service/repository 分层。

**Tech Stack:** Go 1.22 + net/http + GORM + PostgreSQL 16 + Redis 7 + RabbitMQ 3.13 + MinIO + Vue 3 + Element Plus + TypeScript + Pinia + Axios + Vite

---

## 里程碑概览

| 阶段 | 目标 | 预计文件数 | 状态 |
|------|------|-----------|------|
| M0 | 工程骨架 (已基本完成) | ~20 | ✅ |
| M1 | 认证权限系统 | ~30 | ✅ 已完成 |
| M2 | 基础资料管理 | ~25 | ✅ 已完成 |
| M3 | 业务主流程（订单+调度+运输任务） | ~35 | ✅ 已完成 |
| M4 | 运输执行+异常+签收 | ~25 | ✅ 已完成 |
| M5 | 费用结算+报表 | ~25 | ✅ 已完成 |
| M6 | 测试优化+权限补全+部署+前端样式完善 | ~20 | ❌ 未开始 |

---

## 项目现有约定（须严格遵守）

### 后端约定
- Web 框架：标准库 `net/http`（非 Gin），使用 Go 1.22 `mux.HandleFunc("METHOD /path", handler)`
- 模块结构：`internal/modules/{module}/api/handler.go` + `api/routes.go` + `service/service.go`
- 统一响应：`response.Success(w, r, data)` / `response.Fail(w, r, status, code, msg, data)`
- 错误类型：`errors.AppError{Code, Message, Cause, Details}`，需预先定义错误码
- 分页：`pagination.PageRequest` / `pagination.PageResult[T]` 泛型
- 配置：当前仅从环境变量加载（`config.Load()`），后续需扩展支持 YAML
- 日志：`slog.Logger`（JSON handler）
- 中间件：`func(http.Handler) http.Handler` 模式
- 模块名：`github.com/Tangyd893/TMS-Go/backend`
- 数据库迁移：`backend/migrations/` 下编号命名 `{seq}_{desc}.up.sql` / `{seq}_{desc}.down.sql`
- Repository 接口 + 实现模式：`service.go` 中定义接口，`service_impl.go` 实现

### 前端约定
- HTTP 封装：axios 实例，拦截器自动注入 token 并解包 `code===0` 响应
- 状态管理：Pinia stores（`stores/auth.ts` 等）
- 路由：`router/index.ts` 定义路由，`router/guards.ts` 定义守卫
- API 模块：`api/modules/{module}.ts`
- 页面：`views/{module}/` 下按功能拆分
- 组件：`components/` 下按类型组织

---

## 第一阶段：M1 — 认证权限系统 ✅ 已完成

本阶段已完成实现：
- 用户可登录/退出
- JWT 鉴权保护接口
- 用户、角色、权限、菜单的完整 CRUD
- 前端动态路由（基于菜单权限）
- 按钮权限控制
- 操作审计日志写入

### 技术决策
- 密码哈希：bcrypt
- JWT：`golang-jwt/jwt/v5`
- ORM：GORM（已有 `uuid-ossp` 扩展，主键用 UUID）
- 数据库连接：在 `bootstrap` 中初始化 `*gorm.DB`，通过依赖注入传递
- 配置扩展：支持 `config.example.yaml` 中的 database/redis/jwt 配置段
- Auth 中间件：从 Authorization header 提取 Bearer token，解析 JWT 获取 userId

### 详细任务清单

#### Task 1: 扩展配置与依赖
- 扩展 `internal/config/config.go` 支持 database/redis/jwt 配置
- 添加依赖到 go.mod：gorm、postgres driver、jwt、bcrypt
- 创建 `internal/config/config.yaml.go` 支持 YAML 文件加载

#### Task 2: 数据库连接初始化
- 创建 `internal/platform/database/database.go`：初始化 GORM 连接池
- 在 `main.go` 中初始化数据库连接
- 创建 `backend/migrations/000002_create_sys_user.up.sql`：用户表
- 创建 `backend/migrations/000003_create_sys_role.up.sql`：角色表
- 创建 `backend/migrations/000004_create_sys_permission.up.sql`：权限/菜单表
- 创建 `backend/migrations/000005_create_sys_user_role.up.sql`：用户角色关联
- 创建 `backend/migrations/000006_create_sys_role_permission.up.sql`：角色权限关联

#### Task 3: JWT 工具
- 创建 `internal/shared/jwt/jwt.go`：签发/验证 access_token 和 refresh_token

#### Task 4: 密码工具
- 创建 `internal/shared/crypto/password.go`：bcrypt 哈希与校验

#### Task 5: 认证模块 — Service 层
- 创建 `internal/modules/auth/service/service.go`：认证服务接口（Login、Logout、Refresh、ValidateToken）
- 创建 `internal/modules/auth/service/service_impl.go`：认证服务实现（校验用户名密码、签发JWT、黑名单管理）

#### Task 6: 认证模块 — Repository 层
- 创建 `internal/modules/auth/repository/user_repository.go`：用户查询接口
- 创建 `internal/modules/auth/repository/user_repository_pg.go`：PostgreSQL 实现

#### Task 7: 认证模块 — Model 层
- 创建 `internal/modules/auth/model/user.go`：sys_user 实体

#### Task 8: 认证模块 — API 层
- 创建 `internal/modules/auth/api/handler.go`：Login/Logout/Refresh handler
- 创建 `internal/modules/auth/api/routes.go`：注册认证路由

#### Task 9: 认证中间件
- 创建 `internal/middleware/auth.go`：JWT 校验中间件，将 userId 注入 context

#### Task 10: 系统管理模块 — 用户 CRUD
- 创建 `internal/modules/system/model/user.go`、`role.go`、`permission.go`
- 创建 `internal/modules/system/repository/*.go`：用户/角色/权限/菜单 CRUD
- 创建 `internal/modules/system/service/*.go`：业务编排
- 创建 `internal/modules/system/api/*.go`：Handler + 路由

#### Task 11: 操作审计日志
- 创建 `internal/modules/system/repository/operation_log_repository.go`
- 创建 `internal/middleware/audit.go`：审计日志中间件

#### Task 12: 前端对接 — 登录
- 重写 `LoginView.vue`：真实登录表单（用户名/密码），调用登录 API
- 创建 `api/modules/auth.ts`：登录/退出/刷新接口

#### Task 13: 前端对接 — 权限
- 扩展 `stores/auth.ts`：存储用户信息、角色、权限列表
- 扩展 `stores/permission.ts`：菜单权限、按钮权限
- 重写 `router/guards.ts`：白名单检查、token 校验、动态路由生成

#### Task 14: 前端对接 — 系统管理页面
- 创建 `views/system/UserList.vue`：用户列表页
- 创建 `views/system/RoleList.vue`：角色列表页
- 创建 `views/system/MenuList.vue`：菜单管理页
- 扩展 `layouts/BasicLayout.vue`：动态菜单渲染、用户信息/退出

#### Task 15: 集成测试
- 更新 `internal/bootstrap/http.go`：注册所有新路由
- 启动验证：登录 -> 获取 token -> 访问受保护接口
- `testing/smoke.http` 增加认证测试用例

---

## 第二阶段：M2 — 基础资料管理 ✅ 已完成

本阶段已完成实现：
- 客户/承运商/车辆/司机/线路/站点的完整 CRUD
- 字典管理（字典类型 + 字典项）
- 前端列表/表单/详情页面
- 所有操作受权限控制

### 详细任务清单

#### Task 16: 数据库迁移
- 创建 `migrations/000007_create_base_customer.up.sql`
- 创建 `migrations/000008_create_base_carrier.up.sql`
- 创建 `migrations/000009_create_base_vehicle.up.sql`
- 创建 `migrations/000010_create_base_driver.up.sql`
- 创建 `migrations/000011_create_base_route.up.sql`
- 创建 `migrations/000012_create_base_station.up.sql`
- 创建 `migrations/000013_create_sys_dict.up.sql`：字典类型+字典项

#### Task 17: 基础资料模块 — Model 层
- 创建 `internal/modules/base/model/customer.go`、`carrier.go`、`vehicle.go`、`driver.go`、`route.go`、`station.go`

#### Task 18: 基础资料模块 — Repository 层
- 创建 `internal/modules/base/repository/*.go`：各实体 CRUD 接口+实现

#### Task 19: 基础资料模块 — Service 层
- 创建 `internal/modules/base/service/*.go`：各实体业务编排

#### Task 20: 基础资料模块 — API 层
- 创建 `internal/modules/base/api/*.go`：Handler + 路由

#### Task 21: 字典模块
- 创建 `internal/modules/system/dict/`：字典类型+字典项 CRUD
- 添加 Redis 缓存支持（字典热数据缓存）

#### Task 22: 前端基础资料页面
- 创建 `api/modules/base.ts`、`api/modules/dict.ts`
- 创建 `views/base/CustomerList.vue`、`CarrierList.vue`、`VehicleList.vue`、`DriverList.vue`
- 创建 `views/base/RouteList.vue`、`StationList.vue`
- 创建 `views/system/DictList.vue`
- 创建通用组件 `components/table/DataTable.vue`、`components/form/SearchForm.vue`

#### Task 23: 集成测试
- 更新路由注册
- 前端菜单配置注册基础资料菜单项
- 冒烟测试验证所有 CRUD

---

## 第三阶段：M3 — 业务主流程（订单+调度+运输任务）✅ 已完成

本阶段已完成实现：
- 订单创建/编辑/提交/取消/列表/详情
- 订单状态机（draft→submitted→pending_dispatch→dispatched→in_transit→signed→settled）
- 调度计划创建/分配车辆司机/取消
- 运输任务自动生成（调度确认后）
- 前后端完整交互

### 详细任务清单

#### Task 24: 数据库迁移
- 创建 `migrations/000014_create_tms_order.up.sql`
- 创建 `migrations/000015_create_tms_order_cargo.up.sql`
- 创建 `migrations/000016_create_tms_status_history.up.sql`
- 创建 `migrations/000017_create_tms_dispatch_plan.up.sql`
- 创建 `migrations/000018_create_tms_dispatch_vehicle.up.sql`
- 创建 `migrations/000019_create_tms_transport_task.up.sql`

#### Task 25: 订单模块 — 领域层
- 创建 `internal/modules/order/domain/order_status.go`：订单状态枚举
- 创建 `internal/modules/order/domain/order_state_machine.go`：状态流转规则

#### Task 26: 订单模块 — Model/Repository/Service/API
- 创建 `internal/modules/order/model/order.go`、`order_cargo.go`
- 创建 `internal/modules/order/repository/*.go`
- 创建 `internal/modules/order/service/*.go`
- 创建 `internal/modules/order/api/*.go`

#### Task 27: 调度模块
- 创建 `internal/modules/dispatch/model/*.go`
- 创建 `internal/modules/dispatch/domain/dispatch_status.go`
- 创建 `internal/modules/dispatch/repository/*.go`
- 创建 `internal/modules/dispatch/service/*.go`（含生成运输任务逻辑）
- 创建 `internal/modules/dispatch/api/*.go`

#### Task 28: 运输任务模块
- 创建 `internal/modules/transport/model/transport_task.go`
- 创建 `internal/modules/transport/repository/*.go`
- 创建 `internal/modules/transport/service/*.go`
- 创建 `internal/modules/transport/api/*.go`

#### Task 29: 前端订单页面
- 创建 `api/modules/order.ts`
- 创建 `views/order/OrderList.vue`、`OrderCreate.vue`、`OrderDetail.vue`
- 创建 `components/status-tag/StatusTag.vue`：统一状态标签组件

#### Task 30: 前端调度页面
- 创建 `api/modules/dispatch.ts`
- 创建 `views/dispatch/PendingOrderList.vue`、`DispatchPlanList.vue`、`DispatchPlanCreate.vue`

#### Task 31: 集成测试
- 更新路由注册
- 验证订单创建→提交→调度→运输任务生成完整流程

---

## 第四阶段：M4 — 运输执行+异常+签收 ✅ 已完成

#### Task 32: 数据库迁移
- 创建 `migrations/000020_create_tms_transport_node.up.sql`
- 创建 `migrations/000021_create_tms_exception.up.sql`
- 创建 `migrations/000022_create_tms_exception_log.up.sql`
- 创建 `migrations/000023_create_tms_receipt.up.sql`

#### Task 33: 运输执行扩展
- 在 `internal/modules/transport/` 中添加：发车确认、在途节点回传、到达确认、签收确认、回单上传
- 扩展 transport Service 的运输任务状态机

#### Task 34: 异常模块
- 创建 `internal/modules/exception/` 完整模块：异常上报/处理/关闭
- 异常状态机：待处理→处理中→已解决/已关闭

#### Task 35: 文件模块
- 创建 `internal/modules/file/`：MinIO 文件上传/下载/元数据管理
- 文件安全校验（大小、类型、后缀）

#### Task 36: 前端运输+异常页面
- 创建 `views/transport/TaskList.vue`、`TaskDetail.vue`、`NodeTimeline.vue`
- 创建 `views/exception/ExceptionList.vue`、`ExceptionReport.vue`

---

## 第五阶段：M5 — 费用结算+报表 ✅ 已完成

> **已完成：** 数据库迁移（30对up/down全覆盖到fin_*表）、finance模块 全栈（model/repository/service/api）、report模块（service/api）、前端财务/报表页面全栈。

#### Task 37: 数据库迁移 ✅ 已完成
- ✅ `migrations/000024_create_fin_fee_item.up.sql` 已创建
- ✅ `migrations/000025_create_fin_receivable.up.sql` 已创建
- ✅ `migrations/000026_create_fin_payable.up.sql` 已创建
- ✅ `migrations/000027_create_fin_statement.up.sql` 已创建
- ✅ `migrations/000028_create_fin_settlement.up.sql` 已创建

#### Task 38: 财务模块 ✅ 已完成
- ✅ `internal/modules/finance/model/` — 5个模型（fee_item/receivable/payable/statement/settlement）
- ✅ `internal/modules/finance/repository/` — finance_repository.go
- ✅ `internal/modules/finance/service/` — finance_service.go（含 GenerateFees 自动生成应收应付逻辑）
- ✅ `internal/modules/finance/api/` — finance_handler.go（含路由注册 RegisterRoutes）
- ✅ 签收后自动生成费用的逻辑 — GenerateFees handler 已实现

#### Task 39: 报表模块 ✅ 已完成
- ✅ `internal/modules/report/service/` — report_service.go（Dashboard/OrderStats/TransportEfficiency）
- ✅ `internal/modules/report/api/` — report_handler.go（含路由注册 RegisterRoutes）

#### Task 40: 前端财务+报表页面 ✅ 已完成
- ✅ `views/finance/` — ReceivableList/PayableList/StatementList/SettlementList（含搜索/分页/状态标签/创建弹窗/确认操作）
- ✅ `views/report/` — OrderStatsView/TransportEfficiencyView（含统计图表/进度条/完成率环形图）
- ✅ `api/modules/finance.ts` — 13个API接口封装（类型定义完整）
- ✅ `api/modules/report.ts` — 3个API接口封装（类型定义完整）
- ✅ `views/dashboard/DashboardView.vue` — 对接后端真实Dashboard API，指标卡片可点击跳转
- ✅ `layouts/BasicLayout.vue` — 新增费用结算、报表中心菜单组
- ✅ `router/index.ts` — 新增5条财务路由 + 2条报表路由 + 1条调度管理路由

---

## 第六阶段：M6 — 测试优化+部署

#### Task 41: 后端单元测试
- 为关键业务逻辑编写测试：状态机、订单提交、调度校验等

#### Task 42: 前端类型检查与 Lint
- 运行 `npm run typecheck`、`npm run lint` 并修复问题

#### Task 43: 集成/E2E 测试
- 编写关键业务流程的集成测试

#### Task 44: 部署配置完善
- 完善 Dockerfile 多阶段构建
- 完善 CI/CD pipeline 脚本
- 完善部署文档

---

## 实施建议

1. **严格按阶段顺序执行**：每个阶段为下一阶段提供基础数据和能力
2. **每完成一个 Task 就 git commit**：使用 Conventional Commits 格式
3. **后端优先于前端**：先确保 API 可用，再开发前端页面
4. **每个模块完成后验证**：用 `testing/smoke.http` 或 curl 验证接口
5. **敏感信息不入库**：`.env.example` 提交，实际 `.env` 加入 `.gitignore`

## 下一步行动

**M5 已全部完成。当前应开始 M6 工作。** 优先级：

1. **后端单元测试**：状态机、订单提交、调度校验、费用生成等核心业务逻辑
2. **集成/E2E 测试**：关键业务流程端到端验证
3. **前端 Lint 检查**：运行 `npm run lint` 并修复问题
4. **部署配置完善**：Dockerfile 多阶段构建、CI/CD、部署文档

### 已知问题清单 — 全部已修复 ✅

| # | 位置 | 问题描述 | 修复方式 |
|---|------|----------|----------|
| 1 | `frontend/views/order/OrderCreate.vue:73` | `router.push('/orders')` 路径错误 | 修正为 `/order/list` |
| 2 | `frontend/views/exception/ExceptionList.vue:85` | `showDetail()` 函数体为空 | 实现完整详情弹窗（el-descriptions） |
| 3 | `frontend/api/modules/order.ts` ↔ `transport.ts` | TransportTask 类型重复定义 | 从 order.ts 删除，统一在 transport.ts |
| 4 | `frontend/views/dispatch/DispatchList.vue` | 文件存在但未注册到路由 | 注册到 `/dispatch/list` |
