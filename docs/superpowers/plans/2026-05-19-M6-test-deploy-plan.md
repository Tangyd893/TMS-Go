# M6 测试优化+部署 — 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复3个业务Bug、为后端核心模块编写单元测试、前端Lint修复、添加CI/CD配置，完成M6。

**Architecture:** 修Bug→P0纯逻辑测试→P1核心service测试(mock repo)→P2边缘service测试→前端lint→部署配置。Mock策略：手写mock结构体实现Repository接口注入Service。

**Tech Stack:** Go 1.22 + testing + go vet + ESLint + GitHub Actions + GNU Make

---

### Task 1: 修复调度Assign未持久化Details

**Files:**
- Modify: `backend/internal/modules/dispatch/repository/dispatch_repository.go`
- Modify: `backend/internal/modules/dispatch/service/dispatch_service.go`

- [ ] **Step 1: 在DispatchRepository接口中添加UpdateDetail方法**

在 `dispatch_repository.go` 第18行后添加：

```go
	UpdateDetail(ctx context.Context, id uuid.UUID, detail *model.DispatchDetail) error
```

- [ ] **Step 2: 在dispatchRepository中实现UpdateDetail**

在 `dispatchRepository` 的 `FindPendingOrders` 方法后（第77行后）添加：

```go
func (r *dispatchRepository) UpdateDetail(ctx context.Context, id uuid.UUID, detail *model.DispatchDetail) error {
	return r.db.WithContext(ctx).Model(&model.DispatchDetail{}).Where("id = ?", id).
		Updates(map[string]any{
			"vehicle_id": detail.VehicleID,
			"driver_id":  detail.DriverID,
			"carrier_id": detail.CarrierID,
		}).Error
}
```

- [ ] **Step 3: 在Assign方法中持久化Details**

修改 `dispatch_service.go` 第93-108行，将for循环内的修改改为立即持久化：

```go
	for i := range plan.Details {
		if plan.Details[i].ID.String() == detailID {
			d := &plan.Details[i]
			if vehicleID != nil {
				vid := uuid.MustParse(*vehicleID)
				d.VehicleID = &vid
			}
			if driverID != nil {
				did := uuid.MustParse(*driverID)
				d.DriverID = &did
			}
			if carrierID != nil {
				cid := uuid.MustParse(*carrierID)
				d.CarrierID = &cid
			}
			if err := s.repo.UpdateDetail(ctx, d.ID, d); err != nil {
				return err
			}
		}
	}
```

- [ ] **Step 4: 编译验证**

```bash
cd backend && go build ./...
```

---

### Task 2: 修复财务GenerateFees中Payable状态

**Files:**
- Modify: `backend/internal/modules/finance/service/finance_service.go`

- [ ] **Step 1: 修改Payable状态为pending**

修改 `finance_service.go` 第113行，将：
```go
			Status:      string(ReceivableStatusSettled),
```
改为：
```go
			Status:      "pending",
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

---

### Task 3: 修复财务CreateSettlement中PartnerID随机生成

**Files:**
- Modify: `backend/internal/modules/finance/repository/finance_repository.go`
- Modify: `backend/internal/modules/finance/service/finance_service.go`

- [ ] **Step 1: 在FinanceRepository接口中添加FindStatementByID**

在 `finance_repository.go` 第20行后（`ConfirmStatement` 和 `ListSettlements` 之间）添加：

```go
	FindStatementByID(ctx context.Context, id uuid.UUID) (*model.Statement, error)
```

- [ ] **Step 2: 实现FindStatementByID**

在 `financeRepository` 的 `CreateStatement` 方法后添加：

```go
func (r *financeRepository) FindStatementByID(ctx context.Context, id uuid.UUID) (*model.Statement, error) {
	var s model.Statement
	err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error
	return &s, err
}
```

- [ ] **Step 3: 修改CreateSettlement使用真实PartnerID**

修改 `finance_service.go` 第155行起的 `CreateSettlement` 方法，从：

```go
func (s *financeService) CreateSettlement(ctx context.Context, statementID string) (*model.Settlement, error) {
	now := time.Now()
	settleNo := fmt.Sprintf("STL%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
	stmtUID := uuid.MustParse(statementID)

	settle := &model.Settlement{
		SettlementNo:  settleNo,
		StatementID:   stmtUID,
		PartnerID:     uuid.New(),
		PartnerType:   "customer",
		Status:        model.SettlementStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
```

改为：

```go
func (s *financeService) CreateSettlement(ctx context.Context, statementID string) (*model.Settlement, error) {
	now := time.Now()
	settleNo := fmt.Sprintf("STL%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
	stmtUID := uuid.MustParse(statementID)

	stmt, err := s.repo.FindStatementByID(ctx, stmtUID)
	if err != nil {
		return nil, ErrStatementNotFound
	}

	settle := &model.Settlement{
		SettlementNo:  settleNo,
		StatementID:   stmtUID,
		PartnerID:     stmt.PartnerID,
		PartnerName:   stmt.PartnerName,
		PartnerType:   stmt.PartnerType,
		Status:        model.SettlementStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
```

- [ ] **Step 4: 编译验证并提交**

```bash
cd backend && go build ./...
git add -A && git commit -m "fix: 修复调度Assign/财务Payable状态/结算单PartnerID三个Bug"
```

---

### Task 4: P0 — 密码工具单元测试

**Files:**
- Create: `backend/internal/shared/crypto/password_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("my-secret-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if hash == "my-secret-password" {
		t.Fatal("HashPassword() returned plaintext")
	}
}

func TestCheckPasswordMatch(t *testing.T) {
	password := "correct-password"
	hash, _ := HashPassword(password)
	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword() should return true for matching password")
	}
}

func TestCheckPasswordMismatch(t *testing.T) {
	hash, _ := HashPassword("correct-password")
	if CheckPassword("wrong-password", hash) {
		t.Fatal("CheckPassword() should return false for wrong password")
	}
}

func TestCheckPasswordEmpty(t *testing.T) {
	hash, _ := HashPassword("")
	if !CheckPassword("", hash) {
		t.Fatal("CheckPassword() should work with empty password")
	}
}

func TestHashPasswordProducesUniqueHashes(t *testing.T) {
	hash1, _ := HashPassword("same-password")
	hash2, _ := HashPassword("same-password")
	if hash1 == hash2 {
		t.Fatal("HashPassword() should produce different hashes due to salt")
	}
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/shared/crypto/ -v
```

预期: 5个测试全部PASS

---

### Task 5: P0 — JWT工具单元测试

**Files:**
- Create: `backend/internal/shared/jwt/jwt_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package jwt

import (
	"testing"
	"time"
)

func newTestManager() *Manager {
	return NewManager("access-secret", "refresh-secret", 2*time.Hour, 168*time.Hour, "tms-test")
}

func TestGenerateAccessToken(t *testing.T) {
	m := newTestManager()
	token, err := m.GenerateAccessToken("user-1", "admin")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateAccessToken() returned empty token")
	}
}

func TestParseAccessToken(t *testing.T) {
	m := newTestManager()
	token, _ := m.GenerateAccessToken("user-1", "admin")
	claims, err := m.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("ParseAccessToken() UserID = %s, want user-1", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Fatalf("ParseAccessToken() Username = %s, want admin", claims.Username)
	}
	if claims.Issuer != "tms-test" {
		t.Fatalf("ParseAccessToken() Issuer = %s, want tms-test", claims.Issuer)
	}
}

func TestParseAccessTokenWrongSecret(t *testing.T) {
	m := newTestManager()
	token, _ := m.GenerateAccessToken("user-1", "admin")
	m2 := NewManager("wrong-secret", "refresh-secret", 2*time.Hour, 168*time.Hour, "tms-test")
	_, err := m2.ParseAccessToken(token)
	if err == nil {
		t.Fatal("ParseAccessToken() should fail with wrong secret")
	}
}

func TestParseAccessTokenExpired(t *testing.T) {
	m := NewManager("access-secret", "refresh-secret", -1*time.Hour, 168*time.Hour, "tms-test")
	token, _ := m.GenerateAccessToken("user-1", "admin")
	_, err := m.ParseAccessToken(token)
	if err == nil {
		t.Fatal("ParseAccessToken() should fail for expired token")
	}
}

func TestParseAccessTokenTampered(t *testing.T) {
	m := newTestManager()
	token, _ := m.GenerateAccessToken("user-1", "admin")
	tampered := token + "x"
	_, err := m.ParseAccessToken(tampered)
	if err == nil {
		t.Fatal("ParseAccessToken() should fail for tampered token")
	}
}

func TestGenerateAndParseRefreshToken(t *testing.T) {
	m := newTestManager()
	token, err := m.GenerateRefreshToken("user-1", "admin")
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}
	claims, err := m.ParseRefreshToken(token)
	if err != nil {
		t.Fatalf("ParseRefreshToken() error = %v", err)
	}
	if claims.UserID != "user-1" {
		t.Fatalf("ParseRefreshToken() UserID = %s, want user-1", claims.UserID)
	}
}

func TestAccessTokenCannotBeParsedAsRefresh(t *testing.T) {
	m := newTestManager()
	accessToken, _ := m.GenerateAccessToken("user-1", "admin")
	_, err := m.ParseRefreshToken(accessToken)
	if err == nil {
		t.Fatal("Access token should not be parseable as refresh token")
	}
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/shared/jwt/ -v
```

预期: 7个测试全部PASS

---

### Task 6: P0 — 订单状态机单元测试

**Files:**
- Create: `backend/internal/modules/order/domain/order_status_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package domain

import (
	"testing"
)

func TestStateMachineDraftToSubmitted(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusDraft)
	if err := sm.CanSubmit(); err != nil {
		t.Fatalf("CanSubmit() from draft should succeed, got: %v", err)
	}
	status := sm.Submit()
	if status != OrderStatusSubmitted {
		t.Fatalf("Submit() = %s, want %s", status, OrderStatusSubmitted)
	}
}

func TestStateMachineSubmittedToPendingDispatch(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSubmitted)
	if err := sm.CanMoveToPendingDispatch(); err != nil {
		t.Fatalf("CanMoveToPendingDispatch() from submitted should succeed")
	}
	status := sm.MoveToPendingDispatch()
	if status != OrderStatusPendingDispatch {
		t.Fatalf("MoveToPendingDispatch() = %s, want %s", status, OrderStatusPendingDispatch)
	}
}

func TestStateMachinePendingDispatchToDispatched(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusPendingDispatch)
	if err := sm.CanDispatch(); err != nil {
		t.Fatalf("CanDispatch() from pending_dispatch should succeed")
	}
	status := sm.Dispatch()
	if status != OrderStatusDispatched {
		t.Fatalf("Dispatch() = %s, want %s", status, OrderStatusDispatched)
	}
}

func TestStateMachineDispatchedToInTransit(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusDispatched)
	if err := sm.CanStartTransit(); err != nil {
		t.Fatalf("CanStartTransit() from dispatched should succeed")
	}
	status := sm.StartTransit()
	if status != OrderStatusInTransit {
		t.Fatalf("StartTransit() = %s, want %s", status, OrderStatusInTransit)
	}
}

func TestStateMachineInTransitToSigned(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusInTransit)
	if err := sm.CanSign(); err != nil {
		t.Fatalf("CanSign() from in_transit should succeed")
	}
	status := sm.Sign()
	if status != OrderStatusSigned {
		t.Fatalf("Sign() = %s, want %s", status, OrderStatusSigned)
	}
}

func TestStateMachineSignedToSettled(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSigned)
	if err := sm.CanSettle(); err != nil {
		t.Fatalf("CanSettle() from signed should succeed")
	}
	status := sm.Settle()
	if status != OrderStatusSettled {
		t.Fatalf("Settle() = %s, want %s", status, OrderStatusSettled)
	}
}

func TestStateMachineInvalidTransition(t *testing.T) {
	tests := []struct {
		name   string
		status string
		action func(*OrderStateMachine) error
	}{
		{"Submit from submitted", OrderStatusSubmitted, (*OrderStateMachine).CanSubmit},
		{"Dispatch from draft", OrderStatusDraft, (*OrderStateMachine).CanDispatch},
		{"Sign from draft", OrderStatusDraft, (*OrderStateMachine).CanSign},
		{"StartTransit from submitted", OrderStatusSubmitted, (*OrderStateMachine).CanStartTransit},
		{"Settle from draft", OrderStatusDraft, (*OrderStateMachine).CanSettle},
		{"MoveToPendingDispatch from draft", OrderStatusDraft, (*OrderStateMachine).CanMoveToPendingDispatch},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewOrderStateMachine(tt.status)
			if err := tt.action(sm); err == nil {
				t.Fatalf("%s should fail with ErrInvalidStatusTransition", tt.name)
			}
		})
	}
}

func TestStateMachineCancel(t *testing.T) {
	cancelableStatuses := []string{
		OrderStatusDraft, OrderStatusSubmitted, OrderStatusPendingDispatch,
		OrderStatusDispatched, OrderStatusInTransit,
	}
	for _, status := range cancelableStatuses {
		t.Run("Cancel from "+status, func(t *testing.T) {
			sm := NewOrderStateMachine(status)
			if err := sm.CanCancel(); err != nil {
				t.Fatalf("CanCancel() from %s should succeed, got: %v", status, err)
			}
			if s := sm.Cancel(); s != OrderStatusCancelled {
				t.Fatalf("Cancel() = %s from %s, want %s", s, status, OrderStatusCancelled)
			}
		})
	}
}

func TestStateMachineCannotCancelSigned(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSigned)
	if err := sm.CanCancel(); err != ErrOrderAlreadySigned {
		t.Fatalf("CanCancel() from signed should return ErrOrderAlreadySigned, got: %v", err)
	}
}

func TestStateMachineCannotCancelSettled(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusSettled)
	if err := sm.CanCancel(); err != ErrOrderAlreadySigned {
		t.Fatalf("CanCancel() from settled should return ErrOrderAlreadySigned, got: %v", err)
	}
}

func TestStateMachineDoubleCancel(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusDraft)
	sm.Cancel()
	if sm.Status() != OrderStatusCancelled {
		t.Fatal("First cancel should succeed")
	}
	// 第二次Cancel的CanCancel应通过（已取消状态再次取消在CanCancel中没有拒绝逻辑）
	// 根据代码，CanCancel只检查Signed/Settled，不检查已取消
	if err := sm.CanCancel(); err != nil {
		t.Fatalf("CanCancel() from cancelled should succeed, got: %v", err)
	}
}

func TestStateMachineCannotSubmitFromCancelled(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusDraft)
	sm.Cancel()
	if err := sm.CanSubmit(); err != ErrOrderAlreadyCancelled {
		t.Fatalf("CanSubmit() from cancelled should return ErrOrderAlreadyCancelled, got: %v", err)
	}
}

func TestStateMachineCannotSignTwice(t *testing.T) {
	sm := NewOrderStateMachine(OrderStatusInTransit)
	sm.Sign()
	if err := sm.CanSign(); err != ErrOrderAlreadySigned {
		t.Fatalf("CanSign() on already signed should return ErrOrderAlreadySigned, got: %v", err)
	}
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/modules/order/domain/ -v
```

预期: 13个测试全部PASS

- [ ] **Step 3: 提交P0测试**

```bash
git add -A && git commit -m "test: P0 单元测试 — 密码/JWT/订单状态机"
```

---

### Task 7: P1 — 订单服务单元测试

**Files:**
- Create: `backend/internal/modules/order/service/order_service_test.go`

- [ ] **Step 1: 创建mock repo和测试**

```go
package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/domain"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockOrderRepo struct {
	findFn        func(ctx context.Context, id uuid.UUID) (*model.Order, error)
	createFn      func(ctx context.Context, order *model.Order) error
	updateStatusFn func(ctx context.Context, id uuid.UUID, status string) error
	createHistoryFn func(ctx context.Context, h *model.StatusHistory) error
}

func (m *mockOrderRepo) List(ctx context.Context, page, pageSize int, keyword, status string) ([]model.Order, int64, error) {
	return nil, 0, nil
}
func (m *mockOrderRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	return m.findFn(ctx, id)
}
func (m *mockOrderRepo) FindByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *mockOrderRepo) Create(ctx context.Context, order *model.Order) error {
	order.ID = uuid.New()
	return m.createFn(ctx, order)
}
func (m *mockOrderRepo) Update(ctx context.Context, order *model.Order) error { return nil }
func (m *mockOrderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return m.updateStatusFn(ctx, id, status)
}
func (m *mockOrderRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockOrderRepo) CreateStatusHistory(ctx context.Context, h *model.StatusHistory) error {
	return m.createHistoryFn(ctx, h)
}

func TestSubmitOrderSuccess(t *testing.T) {
	orderID := uuid.New()
	order := &model.Order{ID: orderID, Status: domain.OrderStatusDraft}

	repo := &mockOrderRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
			return order, nil
		},
		updateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
			order.Status = status
			return nil
		},
		createHistoryFn: func(ctx context.Context, h *model.StatusHistory) error {
			return nil
		},
	}

	svc := NewOrderService(repo)
	result, err := svc.Submit(context.Background(), orderID.String(), "user-1")
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Status != domain.OrderStatusSubmitted {
		t.Fatalf("Submit() status = %s, want %s", result.Status, domain.OrderStatusSubmitted)
	}
}

func TestSubmitOrderNotFound(t *testing.T) {
	repo := &mockOrderRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := NewOrderService(repo)
	_, err := svc.Submit(context.Background(), uuid.New().String(), "user-1")
	if err != ErrOrderNotFound {
		t.Fatalf("Submit() error = %v, want ErrOrderNotFound", err)
	}
}

func TestSubmitOrderAlreadySubmitted(t *testing.T) {
	orderID := uuid.New()
	order := &model.Order{ID: orderID, Status: domain.OrderStatusSubmitted}

	repo := &mockOrderRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
			return order, nil
		},
	}

	svc := NewOrderService(repo)
	_, err := svc.Submit(context.Background(), orderID.String(), "user-1")
	if err == nil {
		t.Fatal("Submit() on already submitted order should fail")
	}
}

func TestCancelOrderSuccess(t *testing.T) {
	orderID := uuid.New()
	order := &model.Order{ID: orderID, Status: domain.OrderStatusDraft}

	repo := &mockOrderRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
			return order, nil
		},
		updateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
			order.Status = status
			return nil
		},
		createHistoryFn: func(ctx context.Context, h *model.StatusHistory) error {
			return nil
		},
	}

	svc := NewOrderService(repo)
	result, err := svc.Cancel(context.Background(), orderID.String(), "user-1")
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	if result.Status != domain.OrderStatusCancelled {
		t.Fatalf("Cancel() status = %s, want %s", result.Status, domain.OrderStatusCancelled)
	}
}

func TestCreateOrderSuccess(t *testing.T) {
	repo := &mockOrderRepo{
		createFn: func(ctx context.Context, order *model.Order) error {
			order.ID = uuid.New()
			return nil
		},
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
			return &model.Order{ID: id, Status: domain.OrderStatusDraft}, nil
		},
	}
	svc := NewOrderService(repo)
	req := model.CreateOrderRequest{
		CustomerID:      "cust-1",
		ShipperName:     "发货方",
		ReceiverName:    "收货方",
		OriginName:      "上海",
		DestName:        "北京",
		CargoItems: []model.CargoItemInput{
			{CargoName: "货物A", Quantity: 10, Weight: 100, Volume: 5},
		},
	}
	result, err := svc.Create(context.Background(), req, "")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.Status != domain.OrderStatusDraft {
		t.Fatalf("Create() status = %s, want %s", result.Status, domain.OrderStatusDraft)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := &mockOrderRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Order, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := NewOrderService(repo)
	_, err := svc.GetByID(context.Background(), uuid.New().String())
	if err != ErrOrderNotFound {
		t.Fatalf("GetByID() error = %v, want ErrOrderNotFound", err)
	}
}

func assert(t *testing.T, cond bool, msg string) {
	t.Helper()
	if !cond {
		t.Fatal(msg)
	}
}

func init() {
	_ = errors.New(0, "")
	_ = repository.NewOrderRepository(nil)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/modules/order/service/ -v
```

预期: 6个测试全部PASS

---

### Task 8: P1 — 运输服务单元测试

**Files:**
- Create: `backend/internal/modules/transport/service/transport_service_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockTransportRepo struct {
	findFn        func(ctx context.Context, id uuid.UUID) (*model.TransportTask, error)
	createFn      func(ctx context.Context, task *model.TransportTask) error
	updateStatusFn func(ctx context.Context, id uuid.UUID, status string) error
}

func (m *mockTransportRepo) List(ctx context.Context, page, pageSize int, keyword, status string) ([]model.TransportTask, int64, error) {
	return nil, 0, nil
}
func (m *mockTransportRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.TransportTask, error) {
	return m.findFn(ctx, id)
}
func (m *mockTransportRepo) Create(ctx context.Context, task *model.TransportTask) error {
	task.ID = uuid.New()
	return m.createFn(ctx, task)
}
func (m *mockTransportRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return m.updateStatusFn(ctx, id, status)
}
func (m *mockTransportRepo) CreateNode(ctx context.Context, node *model.TransportNode) error { return nil }
func (m *mockTransportRepo) FindNodesByTask(ctx context.Context, taskID uuid.UUID) ([]model.TransportNode, error) {
	return nil, nil
}
func (m *mockTransportRepo) CreateReceipt(ctx context.Context, receipt *model.Receipt) error { return nil }
func (m *mockTransportRepo) FindReceiptByTask(ctx context.Context, taskID uuid.UUID) (*model.Receipt, error) {
	return nil, nil
}

func TestDepartTaskSuccess(t *testing.T) {
	taskID := uuid.New()
	task := &model.TransportTask{ID: taskID, Status: model.TaskStatusPending}
	repo := &mockTransportRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.TransportTask, error) { return task, nil },
		updateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error { return nil },
	}
	svc := NewTransportService(repo)
	result, err := svc.Depart(context.Background(), taskID.String())
	if err != nil {
		t.Fatalf("Depart() error = %v", err)
	}
	if result.Status != model.TaskStatusDeparted {
		t.Fatalf("Depart() status = %s, want %s", result.Status, model.TaskStatusDeparted)
	}
}

func TestDepartTaskTwice(t *testing.T) {
	taskID := uuid.New()
	task := &model.TransportTask{ID: taskID, Status: model.TaskStatusDeparted}
	repo := &mockTransportRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.TransportTask, error) { return task, nil },
	}
	svc := NewTransportService(repo)
	_, err := svc.Depart(context.Background(), taskID.String())
	if err == nil {
		t.Fatal("Depart() on departed task should fail")
	}
}

func TestArriveTaskSuccess(t *testing.T) {
	taskID := uuid.New()
	task := &model.TransportTask{ID: taskID, Status: model.TaskStatusDeparted}
	repo := &mockTransportRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.TransportTask, error) { return task, nil },
		updateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error { return nil },
	}
	svc := NewTransportService(repo)
	result, err := svc.Arrive(context.Background(), taskID.String())
	if err != nil {
		t.Fatalf("Arrive() error = %v", err)
	}
	if result.Status != model.TaskStatusArrived {
		t.Fatalf("Arrive() status = %s, want %s", result.Status, model.TaskStatusArrived)
	}
}

func TestSignTaskBeforeArrival(t *testing.T) {
	taskID := uuid.New()
	task := &model.TransportTask{ID: taskID, Status: model.TaskStatusPending}
	repo := &mockTransportRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.TransportTask, error) { return task, nil },
	}
	svc := NewTransportService(repo)
	_, err := svc.Sign(context.Background(), taskID.String())
	if err == nil {
		t.Fatal("Sign() on pending task should fail")
	}
}

func init() {
	_ = errors.New(0, "")
	_ = repository.NewTransportRepository(nil)
	_ = gorm.ErrRecordNotFound
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/modules/transport/service/ -v
```

预期: 4个测试全部PASS

---

### Task 9: P1 — 异常服务单元测试

**Files:**
- Create: `backend/internal/modules/exception/service/exception_service_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockExceptionRepo struct {
	findFn   func(ctx context.Context, id uuid.UUID) (*model.Exception, error)
	handleFn func(ctx context.Context, id uuid.UUID, handlerID uuid.UUID, handlerName, result string) error
	closeFn  func(ctx context.Context, id uuid.UUID) error
	addLogFn func(ctx context.Context, log *model.ExceptionLog) error
}

func (m *mockExceptionRepo) List(ctx context.Context, page, pageSize int, t, s string) ([]model.Exception, int64, error) {
	return nil, 0, nil
}
func (m *mockExceptionRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Exception, error) {
	return m.findFn(ctx, id)
}
func (m *mockExceptionRepo) Create(ctx context.Context, e *model.Exception) error { return nil }
func (m *mockExceptionRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error { return nil }
func (m *mockExceptionRepo) Handle(ctx context.Context, id uuid.UUID, hid uuid.UUID, hn, result string) error {
	return m.handleFn(ctx, id, hid, hn, result)
}
func (m *mockExceptionRepo) Close(ctx context.Context, id uuid.UUID) error {
	return m.closeFn(ctx, id)
}
func (m *mockExceptionRepo) AddLog(ctx context.Context, log *model.ExceptionLog) error {
	return m.addLogFn(ctx, log)
}
func (m *mockExceptionRepo) FindLogs(ctx context.Context, id uuid.UUID) ([]model.ExceptionLog, error) {
	return nil, nil
}

func TestHandleExceptionSuccess(t *testing.T) {
	excID := uuid.New()
	exc := &model.Exception{ID: excID, Status: model.ExceptionStatusPending}
	repo := &mockExceptionRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Exception, error) { return exc, nil },
		handleFn: func(ctx context.Context, id uuid.UUID, hid uuid.UUID, hn, result string) error {
			return nil
		},
		addLogFn: func(ctx context.Context, log *model.ExceptionLog) error { return nil },
	}
	svc := NewExceptionService(repo)
	if err := svc.Handle(context.Background(), excID.String(), uuid.New().String(), "处理人", "已处理"); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandleExceptionNotPending(t *testing.T) {
	excID := uuid.New()
	exc := &model.Exception{ID: excID, Status: model.ExceptionStatusResolved}
	repo := &mockExceptionRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Exception, error) { return exc, nil },
	}
	svc := NewExceptionService(repo)
	err := svc.Handle(context.Background(), excID.String(), uuid.New().String(), "处理人", "已处理")
	if err == nil {
		t.Fatal("Handle() on resolved exception should fail")
	}
}

func TestCloseExceptionSuccess(t *testing.T) {
	excID := uuid.New()
	exc := &model.Exception{ID: excID, Status: model.ExceptionStatusResolved}
	repo := &mockExceptionRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.Exception, error) { return exc, nil },
		closeFn: func(ctx context.Context, id uuid.UUID) error { return nil },
	}
	svc := NewExceptionService(repo)
	if err := svc.Close(context.Background(), excID.String()); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func init() {
	_ = errors.New(0, "")
	_ = repository.NewExceptionRepository(nil)
	_ = gorm.ErrRecordNotFound
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/modules/exception/service/ -v
```

预期: 3个测试全部PASS

---

### Task 10: P2 — 调度服务单元测试

**Files:**
- Create: `backend/internal/modules/dispatch/service/dispatch_service_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockDispatchRepo struct {
	findFn        func(ctx context.Context, id uuid.UUID) (*model.DispatchPlan, error)
	createFn      func(ctx context.Context, plan *model.DispatchPlan) error
	updateStatusFn func(ctx context.Context, id uuid.UUID, status string) error
	updateDetailFn func(ctx context.Context, id uuid.UUID, detail *model.DispatchDetail) error
}

func (m *mockDispatchRepo) List(ctx context.Context, page, pageSize int) ([]model.DispatchPlan, int64, error) {
	return nil, 0, nil
}
func (m *mockDispatchRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.DispatchPlan, error) {
	return m.findFn(ctx, id)
}
func (m *mockDispatchRepo) Create(ctx context.Context, plan *model.DispatchPlan) error {
	plan.ID = uuid.New()
	return m.createFn(ctx, plan)
}
func (m *mockDispatchRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return m.updateStatusFn(ctx, id, status)
}
func (m *mockDispatchRepo) FindPendingOrders(ctx context.Context) ([]repository.PendingOrder, error) {
	return nil, nil
}
func (m *mockDispatchRepo) UpdateDetail(ctx context.Context, id uuid.UUID, detail *model.DispatchDetail) error {
	return m.updateDetailFn(ctx, id, detail)
}

func TestCreateDispatchPlanSuccess(t *testing.T) {
	repo := &mockDispatchRepo{
		createFn: func(ctx context.Context, plan *model.DispatchPlan) error {
			plan.ID = uuid.New()
			return nil
		},
		findFn: func(ctx context.Context, id uuid.UUID) (*model.DispatchPlan, error) {
			return &model.DispatchPlan{ID: id, Status: model.DispatchStatusPending}, nil
		},
	}
	svc := NewDispatchService(repo)
	plan, err := svc.Create(context.Background(), []string{uuid.New().String()}, "user-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if plan.Status != model.DispatchStatusPending {
		t.Fatalf("Create() status = %s, want %s", plan.Status, model.DispatchStatusPending)
	}
}

func TestAssignDispatchDetail(t *testing.T) {
	planID := uuid.New()
	detailID := uuid.New()
	plan := &model.DispatchPlan{
		ID:     planID,
		Status: model.DispatchStatusPending,
		Details: []model.DispatchDetail{
			{ID: detailID, PlanID: planID, OrderID: uuid.New()},
		},
	}
	detailUpdated := false
	repo := &mockDispatchRepo{
		findFn: func(ctx context.Context, id uuid.UUID) (*model.DispatchPlan, error) {
			return plan, nil
		},
		updateStatusFn: func(ctx context.Context, id uuid.UUID, status string) error {
			return nil
		},
		updateDetailFn: func(ctx context.Context, id uuid.UUID, detail *model.DispatchDetail) error {
			detailUpdated = true
			return nil
		},
	}
	svc := NewDispatchService(repo)
	vid := uuid.New().String()
	did := uuid.New().String()
	err := svc.Assign(context.Background(), planID.String(), detailID.String(), &vid, &did, nil)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	if !detailUpdated {
		t.Fatal("Assign() should persist Detail changes via UpdateDetail")
	}
}

func init() {
	_ = errors.New(0, "")
	_ = repository.NewDispatchRepository(nil)
	_ = gorm.ErrRecordNotFound
}
```

- [ ] **Step 2: 编译确保UpdateDetail接口匹配**

mockDispatchRepo需要包含UpdateDetail方法，这是因为我们在Task 1中添加了接口方法。需要确保 `repository.pendingOrder` 的字段名 `PendingOrder` 可导出。检查 `dispatch_repository.go` 中 `pendingOrder` 是小写，需要改用空接口或导出类型。

修改 `dispatch_service_test.go` 中 `FindPendingOrders` 返回类型：

```go
func (m *mockDispatchRepo) FindPendingOrders(ctx context.Context) (any, error) {
	return nil, nil
}
```

并将 `dispatch_service_test.go` 的import中移除 `dispatch/repository`（不再需要引用 `pendingOrder`）。

- [ ] **Step 3: 运行测试**

```bash
cd backend && go test ./internal/modules/dispatch/service/ -v
```

预期: 2个测试全部PASS

---

### Task 11: P2 — 财务服务单元测试

**Files:**
- Create: `backend/internal/modules/finance/service/finance_service_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockFinanceRepo struct {
	genReceivableFn func(ctx context.Context, r *model.Receivable) error
	genPayableFn    func(ctx context.Context, p *model.Payable) error
	findStmtFn      func(ctx context.Context, id uuid.UUID) (*model.Statement, error)
	createStmtFn    func(ctx context.Context, s *model.Statement) error
	createSettleFn  func(ctx context.Context, s *model.Settlement) error
}

func (m *mockFinanceRepo) ListReceivables(ctx context.Context, page, pageSize int, status string) ([]model.Receivable, int64, error) {
	return nil, 0, nil
}
func (m *mockFinanceRepo) ListPayables(ctx context.Context, page, pageSize int, status string) ([]model.Payable, int64, error) {
	return nil, 0, nil
}
func (m *mockFinanceRepo) ListFeeItems(ctx context.Context, page, pageSize int) ([]model.FeeItem, int64, error) {
	return nil, 0, nil
}
func (m *mockFinanceRepo) CreateFeeItem(ctx context.Context, item *model.FeeItem) error { return nil }
func (m *mockFinanceRepo) GenerateReceivable(ctx context.Context, r *model.Receivable) error {
	return m.genReceivableFn(ctx, r)
}
func (m *mockFinanceRepo) GeneratePayable(ctx context.Context, p *model.Payable) error {
	return m.genPayableFn(ctx, p)
}
func (m *mockFinanceRepo) ListStatements(ctx context.Context, page, pageSize int) ([]model.Statement, int64, error) {
	return nil, 0, nil
}
func (m *mockFinanceRepo) CreateStatement(ctx context.Context, s *model.Statement) error {
	return m.createStmtFn(ctx, s)
}
func (m *mockFinanceRepo) ConfirmStatement(ctx context.Context, id uuid.UUID, confirmedBy uuid.UUID) error {
	return nil
}
func (m *mockFinanceRepo) FindStatementByID(ctx context.Context, id uuid.UUID) (*model.Statement, error) {
	return m.findStmtFn(ctx, id)
}
func (m *mockFinanceRepo) ListSettlements(ctx context.Context, page, pageSize int) ([]model.Settlement, int64, error) {
	return nil, 0, nil
}
func (m *mockFinanceRepo) CreateSettlement(ctx context.Context, s *model.Settlement) error {
	return m.createSettleFn(ctx, s)
}
func (m *mockFinanceRepo) CompleteSettlement(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockFinanceRepo) UpdateReceivableStatus(ctx context.Context, id uuid.UUID, status string) error {
	return nil
}
func (m *mockFinanceRepo) UpdatePayableStatus(ctx context.Context, id uuid.UUID, status string) error {
	return nil
}

func TestGenerateFeesCreatesReceivable(t *testing.T) {
	receivableCreated := false
	repo := &mockFinanceRepo{
		genReceivableFn: func(ctx context.Context, r *model.Receivable) error {
			receivableCreated = true
			if r.Status != "pending" {
				t.Fatalf("Receivable status = %s, want pending", r.Status)
			}
			return nil
		},
		genPayableFn: func(ctx context.Context, p *model.Payable) error {
			return nil
		},
	}
	svc := NewFinanceService(repo)
	err := svc.GenerateFees(context.Background(),
		uuid.New().String(), uuid.New().String(), "",
		"客户A", "", 1000.0)
	if err != nil {
		t.Fatalf("GenerateFees() error = %v", err)
	}
	if !receivableCreated {
		t.Fatal("GenerateFees() should create receivable for customer")
	}
}

func TestGenerateFeesCreatesPayable(t *testing.T) {
	payableCreated := false
	repo := &mockFinanceRepo{
		genReceivableFn: func(ctx context.Context, r *model.Receivable) error { return nil },
		genPayableFn: func(ctx context.Context, p *model.Payable) error {
			payableCreated = true
			if p.Status != "pending" {
				t.Fatalf("Payable status = %s, want pending", p.Status)
			}
			return nil
		},
	}
	svc := NewFinanceService(repo)
	err := svc.GenerateFees(context.Background(),
		uuid.New().String(), "", uuid.New().String(),
		"", "承运商B", 1000.0)
	if err != nil {
		t.Fatalf("GenerateFees() error = %v", err)
	}
	if !payableCreated {
		t.Fatal("GenerateFees() should create payable for carrier")
	}
}

func TestCreateSettlement(t *testing.T) {
	stmtID := uuid.New()
	stmt := &model.Statement{
		ID:          stmtID,
		PartnerID:   uuid.New(),
		PartnerName: "客户A",
		PartnerType: "customer",
	}
	repo := &mockFinanceRepo{
		findStmtFn: func(ctx context.Context, id uuid.UUID) (*model.Statement, error) {
			return stmt, nil
		},
		createSettleFn: func(ctx context.Context, s *model.Settlement) error {
			return nil
		},
	}
	svc := NewFinanceService(repo)
	result, err := svc.CreateSettlement(context.Background(), stmtID.String())
	if err != nil {
		t.Fatalf("CreateSettlement() error = %v", err)
	}
	if result.PartnerID != stmt.PartnerID {
		t.Fatalf("CreateSettlement() PartnerID = %s, want %s", result.PartnerID, stmt.PartnerID)
	}
	if result.PartnerName != stmt.PartnerName {
		t.Fatalf("CreateSettlement() PartnerName = %s, want %s", result.PartnerName, stmt.PartnerName)
	}
}

func init() {
	_ = repository.NewFinanceRepository(nil)
	_ = gorm.ErrRecordNotFound
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/modules/finance/service/ -v
```

预期: 3个测试全部PASS

- [ ] **Step 3: 提交所有P1+P2测试**

```bash
cd backend && go test ./... -v
git add -A && git commit -m "test: P1+P2 单元测试 — 订单/运输/异常/调度/财务 service"
```

---

### Task 12: 前端Lint修复

**Files:**
- `/mnt/d/workspace/coding/TMS-Go/frontend/` 下所有被 ESLint 标记的文件

- [ ] **Step 1: 运行 ESLint 检查**

```bash
npm run lint
```

查看输出，记录所有 error 和 warning。

- [ ] **Step 2: 修复 ESLint 报错**

根据 lint 输出逐文件修复。常见问题：
- `any` 类型替换为具体类型
- 未使用变量删除或加 `_` 前缀
- props 添加显式类型注解
- 组件 emits 显式声明

- [ ] **Step 3: 确认修复**

```bash
npm run lint
```

预期: 零 error / 零 warning

- [ ] **Step 4: 提交**

```bash
git add -A && git commit -m "fix: 前端 ESLint 修复"
```

---

### Task 13: 创建 Makefile

**Files:**
- Create: `Makefile`

- [ ] **Step 1: 创建 Makefile**

```makefile
.PHONY: dev test lint typecheck build docker-up docker-down clean

dev:
	docker compose -f docker/docker-compose.yml up -d postgres redis rabbitmq minio

test:
	cd backend && go test ./... -v -cover

test-race:
	cd backend && go test -race ./... -v

lint:
	cd backend && go vet ./...
	cd frontend && npm run typecheck && npm run lint

typecheck:
	cd frontend && npm run typecheck

build:
	cd backend && go build ./cmd/server
	cd frontend && npm run build

docker-up:
	docker compose -f docker/docker-compose.yml up -d

docker-down:
	docker compose -f docker/docker-compose.yml down

clean:
	docker compose -f docker/docker-compose.yml down -v
	rm -rf backend/tmp/ frontend/dist/
```

- [ ] **Step 2: 验证 Makefile**

```bash
make test
```

- [ ] **Step 3: 提交**

```bash
git add Makefile && git commit -m "ci: 添加 Makefile 构建/测试/lint 目标"
```

---

### Task 14: 创建 GitHub Actions CI

**Files:**
- Create: `.github/workflows/ci.yml`

- [ ] **Step 1: 创建 CI 工作流**

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  backend:
    name: Backend Tests
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16-alpine
        env:
          POSTGRES_DB: tms
          POSTGRES_USER: tms
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
    env:
      DB_HOST: localhost
      DB_PORT: 5432
      DB_NAME: tms
      DB_USER: tms
      DB_PASSWORD: test
      DB_SSLMODE: disable
      REDIS_HOST: localhost
      REDIS_PORT: 6379
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Create uuid-ossp extension
        run: PGPASSWORD=test psql -h localhost -U tms -d tms -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp"'
      - name: Run tests
        run: cd backend && go test ./... -v -cover
      - name: Run go vet
        run: cd backend && go vet ./...

  frontend:
    name: Frontend Checks
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
      - name: Install dependencies
        run: cd frontend && npm ci
      - name: Type check
        run: cd frontend && npm run typecheck
      - name: Lint
        run: cd frontend && npm run lint
```

- [ ] **Step 2: 提交**

```bash
mkdir -p .github/workflows
git add .github/ && git commit -m "ci: 添加 GitHub Actions CI 流水线"
```

---

### Task 15: 更新计划文档

**Files:**
- Modify: `docs/superpowers/plans/2026-05-19-TMS-full-plan.md`

- [ ] **Step 1: 更新M6状态**

将M6行改为：

```markdown
| M6 | 测试优化+部署 | ~20 | ✅ 已完成 |
```

更新 M6 下各 Task 状态为已完成。

- [ ] **Step 2: 提交**

```bash
git add docs/ && git commit -m "docs: 更新 M6 为已完成状态"
```

---

### 最终验证

- [ ] **运行全部后端测试并生成覆盖率报告**

```bash
cd backend && go test ./... -cover -coverprofile=coverage.out && go tool cover -func=coverage.out
```

- [ ] **确认 git 状态干净**

```bash
git status && git log --oneline -5
```
