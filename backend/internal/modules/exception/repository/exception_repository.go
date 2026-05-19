package repository

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExceptionRepository interface {
	List(ctx context.Context, page, pageSize int, exceptionType, status string) ([]model.Exception, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Exception, error)
	Create(ctx context.Context, e *model.Exception) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Handle(ctx context.Context, id uuid.UUID, handlerID uuid.UUID, handlerName, result string) error
	Close(ctx context.Context, id uuid.UUID) error
	AddLog(ctx context.Context, log *model.ExceptionLog) error
	FindLogs(ctx context.Context, exceptionID uuid.UUID) ([]model.ExceptionLog, error)
}

type exceptionRepository struct {
	db *gorm.DB
}

func NewExceptionRepository(db *gorm.DB) ExceptionRepository {
	return &exceptionRepository{db: db}
}

func (r *exceptionRepository) List(ctx context.Context, page, pageSize int, exceptionType, status string) ([]model.Exception, int64, error) {
	var items []model.Exception
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Exception{})
	if exceptionType != "" {
		query = query.Where("exception_type = ?", exceptionType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&items).Error
	return items, total, err
}

func (r *exceptionRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Exception, error) {
	var e model.Exception
	err := r.db.WithContext(ctx).First(&e, "id = ?", id).Error
	return &e, err
}

func (r *exceptionRepository) Create(ctx context.Context, e *model.Exception) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *exceptionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.Exception{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *exceptionRepository) Handle(ctx context.Context, id uuid.UUID, handlerID uuid.UUID, handlerName, result string) error {
	return r.db.WithContext(ctx).Model(&model.Exception{}).Where("id = ?", id).
		Updates(map[string]any{
			"status":        model.ExceptionStatusResolved,
			"handler_id":    handlerID,
			"handler_name":  handlerName,
			"handle_result": result,
			"handle_at":     gorm.Expr("now()"),
		}).Error
}

func (r *exceptionRepository) Close(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.Exception{}).Where("id = ?", id).
		Update("status", model.ExceptionStatusClosed).Error
}

func (r *exceptionRepository) AddLog(ctx context.Context, log *model.ExceptionLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *exceptionRepository) FindLogs(ctx context.Context, exceptionID uuid.UUID) ([]model.ExceptionLog, error) {
	var logs []model.ExceptionLog
	err := r.db.WithContext(ctx).Where("exception_id = ?", exceptionID).
		Order("created_at desc").Find(&logs).Error
	return logs, err
}
