package service

import (
	"context"
	"testing"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/model"
	"github.com/google/uuid"
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
