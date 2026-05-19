package repository

import (
	"context"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/dispatch/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispatchRepository interface {
	List(ctx context.Context, page, pageSize int) ([]model.DispatchPlan, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.DispatchPlan, error)
	Create(ctx context.Context, plan *model.DispatchPlan) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	UpdateDetail(ctx context.Context, id uuid.UUID, detail *model.DispatchDetail) error
	FindPendingOrders(ctx context.Context) ([]PendingOrder, error)
}

type PendingOrder struct {
	ID     uuid.UUID `json:"id"`
	OrderNo string   `json:"orderNo"`
	OriginName string `json:"originName"`
	DestName string `json:"destName"`
}

type dispatchRepository struct {
	db *gorm.DB
}

func NewDispatchRepository(db *gorm.DB) DispatchRepository {
	return &dispatchRepository{db: db}
}

func (r *dispatchRepository) List(ctx context.Context, page, pageSize int) ([]model.DispatchPlan, int64, error) {
	var plans []model.DispatchPlan
	var total int64
	query := r.db.WithContext(ctx).Model(&model.DispatchPlan{}).Preload("Details")
	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&plans).Error
	return plans, total, err
}

func (r *dispatchRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.DispatchPlan, error) {
	var plan model.DispatchPlan
	err := r.db.WithContext(ctx).Preload("Details").First(&plan, "id = ?", id).Error
	return &plan, err
}

func (r *dispatchRepository) Create(ctx context.Context, plan *model.DispatchPlan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(plan).Error; err != nil { return err }
		for i := range plan.Details {
			if plan.Details[i].ID == uuid.Nil {
				plan.Details[i].ID = uuid.New()
			}
			plan.Details[i].PlanID = plan.ID
			plan.Details[i].CreatedAt = time.Now()
		}
		return tx.Create(&plan.Details).Error
	})
}

func (r *dispatchRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).Model(&model.DispatchPlan{}).Where("id = ?", id).
		Update("status", status).Error
}

func (r *dispatchRepository) FindPendingOrders(ctx context.Context) ([]PendingOrder, error) {
	var orders []PendingOrder
	err := r.db.WithContext(ctx).Table("tms_order").
		Select("id, order_no, origin_name, dest_name").
		Where("status = ? AND deleted_at IS NULL", "pending_dispatch").
		Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *dispatchRepository) UpdateDetail(ctx context.Context, id uuid.UUID, detail *model.DispatchDetail) error {
	return r.db.WithContext(ctx).Model(&model.DispatchDetail{}).Where("id = ?", id).
		Updates(map[string]any{
			"vehicle_id": detail.VehicleID,
			"driver_id":  detail.DriverID,
			"carrier_id": detail.CarrierID,
		}).Error
}
