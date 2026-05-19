# M6 测试优化+部署 — 设计文档

> **目标**：修 Bug → 后端单元测试 → 前端 Lint → CI/CD 部署配置，达到生产就绪标准。

---

## Phase 1: 业务 Bug 修复

### Bug 1: 调度 Assign 未持久化 Details

- **文件**：`backend/internal/modules/dispatch/service/dispatch_service.go`
- **问题**：`Assign` 方法遍历 `plan.Details` 设置 VehicleID/DriverID/CarrierID，但只调 `repo.UpdateStatus()` 更新状态，Details 修改丢失
- **修复**：新增 `repo.UpdateDetails()` 方法，在 Assign 中分别持久化每个 Detail

### Bug 2: 财务 Payable 状态错误

- **文件**：`backend/internal/modules/finance/service/finance_service.go`
- **问题**：`GenerateFees` 中 Payable.Status 直接设为 `"settled"`，跳过待确认→已确认→已结算的完整流程
- **修复**：Payable.Status 改为 `"pending"`

### Bug 3: 结算单 PartnerID 随机生成

- **文件**：`backend/internal/modules/finance/service/finance_service.go`
- **问题**：`CreateSettlement` 中 `PartnerID: uuid.New()` 是随机值，应从 Statement 中获取真实合作方 ID
- **修复**：先 `repo.FindStatementByID()` 查询对账单，取其 PartnerID

---

## Phase 2: 后端单元测试

### Mock 策略
所有 Service 通过构造函数注入 Repository 接口，手写 mock 结构体（同步方法，无需 goroutine/channel 同步）：

```go
type mockOrderRepo struct {
    findFn  func(ctx context.Context, id uuid.UUID) (*model.Order, error)
    updateFn func(ctx context.Context, id uuid.UUID, status string) error
    // ...
}
func (m *mockOrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
    return m.findFn(ctx, id)
}
```

### 测试文件清单

| 优先级 | 文件 | 测试内容 |
|--------|------|---------|
| P0 | `order/domain/order_status_test.go` | 7状态×8转换全部通过/非法路径/重复取消/已签收后取消 |
| P0 | `shared/crypto/password_test.go` | Hash成功/Hash失败/校验匹配/校验不匹配/空密码 |
| P0 | `shared/jwt/jwt_test.go` | 正常签发校验/过期Token/篡改Token/不同密钥/刷新Token |
| P1 | `order/service/order_service_test.go` | Create成功/Submit门控/Cancel门控/状态历史写入/mock repo错误 |
| P1 | `transport/service/transport_service_test.go` | Depart门控/Arrive门控/Sign门控/重复发车/mock错误 |
| P1 | `exception/service/exception_service_test.go` | Handle成功/Handle非pending拒绝/Close非resolved拒绝 |
| P2 | `dispatch/service/dispatch_service_test.go` | Create成功/Assign成功/Assign持久化Details/Cancel |
| P2 | `finance/service/finance_service_test.go` | GenerateFees生成应收应付/对账创建/结算创建/状态流转 |

### 测试命令
```bash
go test ./... -v -cover
go test -race ./...
```

---

## Phase 3: 前端 Lint

- 运行 `npm run lint` 检查 ESLint 问题
- 修复编译警告和格式问题
- 确认 `eslint.config.js` 配置存在

---

## Phase 4: CI/CD 部署配置

### Makefile (项目根目录)
```makefile
.PHONY: dev test lint build docker-up docker-down

dev:
    docker compose -f docker/docker-compose.yml up -d postgres redis rabbitmq minio

test:
    cd backend && go test ./... -v -cover

lint:
    cd backend && go vet ./...
    cd frontend && npm run typecheck && npm run lint

build:
    cd backend && go build ./cmd/server
    cd frontend && npm run build

docker-up:
    docker compose -f docker/docker-compose.yml up -d

docker-down:
    docker compose -f docker/docker-compose.yml down
```

### GitHub Actions CI (`.github/workflows/ci.yml`)
- 触发条件：push/PR to main, develop
- 步骤：checkout → go test → go vet → npm install → npm run typecheck → npm run lint → docker build

---

## 实施顺序

1. 修 Bug（先修再测，避免测到已知 Bug）
2. P0 单元测试（状态机、密码、JWT）
3. P1 单元测试（订单、运输、异常 service）
4. P2 单元测试（调度、财务 service）
5. 前端 Lint 修复
6. Makefile + GitHub Actions CI
7. 更新 TMS-full-plan.md M6 状态
