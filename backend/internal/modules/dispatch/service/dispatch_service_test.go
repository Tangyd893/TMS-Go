package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/repository"
	"github.com/google/uuid"
)

type mockDispatchRepo struct {
	findFn         func(ctx context.Context, id uuid.UUID) (*model.DispatchPlan, error)
	createFn       func(ctx context.Context, plan *model.DispatchPlan) error
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
	plan, err := svc.Create(context.Background(), []string{uuid.New().String()}, uuid.New().String())
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
