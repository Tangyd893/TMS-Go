package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/domain"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/order/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockOrderRepo struct {
	findFn          func(ctx context.Context, id uuid.UUID) (*model.Order, error)
	createFn        func(ctx context.Context, order *model.Order) error
	updateStatusFn  func(ctx context.Context, id uuid.UUID, status string) error
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
	result, err := svc.Submit(context.Background(), orderID.String(), uuid.New().String())
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
	_, err := svc.Submit(context.Background(), uuid.New().String(), uuid.New().String())
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
	_, err := svc.Submit(context.Background(), orderID.String(), uuid.New().String())
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
	result, err := svc.Cancel(context.Background(), orderID.String(), uuid.New().String())
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
	w100 := float64(100)
	v5 := float64(5)
	req := model.CreateOrderRequest{
		ShipperName:  "发货方",
		ReceiverName: "收货方",
		OriginName:   "上海",
		DestName:     "北京",
		CargoItems: []model.CreateCargoRequest{
			{CargoName: "货物A", Quantity: 10, Weight: &w100, Volume: &v5},
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
