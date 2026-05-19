package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/model"
	"github.com/google/uuid"
)

type mockTransportRepo struct {
	findFn         func(ctx context.Context, id uuid.UUID) (*model.TransportTask, error)
	createFn       func(ctx context.Context, task *model.TransportTask) error
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
