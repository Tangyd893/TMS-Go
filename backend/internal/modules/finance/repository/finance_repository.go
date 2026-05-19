package repository

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/finance/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinanceRepository interface {
	ListReceivables(ctx context.Context, page, pageSize int, status string) ([]model.Receivable, int64, error)
	ListPayables(ctx context.Context, page, pageSize int, status string) ([]model.Payable, int64, error)
	ListFeeItems(ctx context.Context, page, pageSize int) ([]model.FeeItem, int64, error)
	CreateFeeItem(ctx context.Context, item *model.FeeItem) error
	GenerateReceivable(ctx context.Context, r *model.Receivable) error
	GeneratePayable(ctx context.Context, p *model.Payable) error
	ListStatements(ctx context.Context, page, pageSize int) ([]model.Statement, int64, error)
	CreateStatement(ctx context.Context, s *model.Statement) error
	ConfirmStatement(ctx context.Context, id uuid.UUID, confirmedBy uuid.UUID) error
	ListSettlements(ctx context.Context, page, pageSize int) ([]model.Settlement, int64, error)
	CreateSettlement(ctx context.Context, s *model.Settlement) error
	CompleteSettlement(ctx context.Context, id uuid.UUID) error
	UpdateReceivableStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdatePayableStatus(ctx context.Context, id uuid.UUID, status string) error
}

type financeRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) FinanceRepository {
	return &financeRepository{db: db}
}

func (r *financeRepository) ListReceivables(ctx context.Context, page, pageSize int, status string) ([]model.Receivable, int64, error) {
	var items []model.Receivable
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Receivable{})
	if status != "" { query = query.Where("status = ?", status) }
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&items).Error
	return items, total, err
}

func (r *financeRepository) ListPayables(ctx context.Context, page, pageSize int, status string) ([]model.Payable, int64, error) {
	var items []model.Payable
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Payable{})
	if status != "" { query = query.Where("status = ?", status) }
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&items).Error
	return items, total, err
}

func (r *financeRepository) ListFeeItems(ctx context.Context, page, pageSize int) ([]model.FeeItem, int64, error) {
	var items []model.FeeItem
	var total int64
	query := r.db.WithContext(ctx).Model(&model.FeeItem{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&items).Error
	return items, total, err
}

func (r *financeRepository) CreateFeeItem(ctx context.Context, item *model.FeeItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *financeRepository) GenerateReceivable(ctx context.Context, rv *model.Receivable) error {
	return r.db.WithContext(ctx).Create(rv).Error
}

func (r *financeRepository) GeneratePayable(ctx context.Context, p *model.Payable) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *financeRepository) ListStatements(ctx context.Context, page, pageSize int) ([]model.Statement, int64, error) {
	var items []model.Statement
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Statement{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&items).Error
	return items, total, err
}

func (r *financeRepository) CreateStatement(ctx context.Context, s *model.Statement) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *financeRepository) ConfirmStatement(ctx context.Context, id uuid.UUID, confirmedBy uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Statement{}).Where("id = ?", id).
		Updates(map[string]any{
			"status":       model.StatementStatusConfirmed,
			"confirmed_at": gorm.Expr("now()"),
			"confirmed_by": confirmedBy,
		}).Error
}

func (r *financeRepository) ListSettlements(ctx context.Context, page, pageSize int) ([]model.Settlement, int64, error) {
	var items []model.Settlement
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Settlement{})
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&items).Error
	return items, total, err
}

func (r *financeRepository) CreateSettlement(ctx context.Context, s *model.Settlement) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *financeRepository) CompleteSettlement(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Settlement{}).Where("id = ?", id).
		Updates(map[string]any{
			"status":     model.SettlementStatusCompleted,
			"settled_at": gorm.Expr("now()"),
		}).Error
}

func (r *financeRepository) UpdateReceivableStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Receivable{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *financeRepository) UpdatePayableStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Payable{}).Where("id = ?", id).
		Update("status", status).Error
}
