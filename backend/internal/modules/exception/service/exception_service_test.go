package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/model"
	"github.com/google/uuid"
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
