package repository

import (
	"context"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransportRepository interface {
	List(ctx context.Context, page, pageSize int, keyword, status string) ([]model.TransportTask, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.TransportTask, error)
	Create(ctx context.Context, task *model.TransportTask) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error

	// 节点
	CreateNode(ctx context.Context, node *model.TransportNode) error
	FindNodesByTask(ctx context.Context, taskID uuid.UUID) ([]model.TransportNode, error)

	// 回单
	CreateReceipt(ctx context.Context, receipt *model.Receipt) error
	FindReceiptByTask(ctx context.Context, taskID uuid.UUID) (*model.Receipt, error)
}

type transportRepository struct {
	db *gorm.DB
}

func NewTransportRepository(db *gorm.DB) TransportRepository {
	return &transportRepository{db: db}
}

func (r *transportRepository) List(ctx context.Context, page, pageSize int, keyword, status string) ([]model.TransportTask, int64, error) {
	var tasks []model.TransportTask
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TransportTask{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("task_no LIKE ? OR driver_name LIKE ? OR plate_no LIKE ?", like, like, like)
	}
	if status != "" { query = query.Where("status = ?", status) }

	query.Count(&total)
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&tasks).Error
	return tasks, total, err
}

func (r *transportRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.TransportTask, error) {
	var task model.TransportTask
	err := r.db.WithContext(ctx).First(&task, "id = ?", id).Error
	return &task, err
}

func (r *transportRepository) Create(ctx context.Context, task *model.TransportTask) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *transportRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	updates := map[string]any{"status": status, "updated_at": time.Now()}
	if status == model.TaskStatusDeparted {
		now := time.Now()
		updates["actual_depart_time"] = &now
	}
	if status == model.TaskStatusArrived {
		now := time.Now()
		updates["actual_arrive_time"] = &now
	}
	return r.db.WithContext(ctx).Model(&model.TransportTask{}).Where("id = ?", id).Updates(updates).Error
}

func (r *transportRepository) CreateNode(ctx context.Context, node *model.TransportNode) error {
	return r.db.WithContext(ctx).Create(node).Error
}

func (r *transportRepository) FindNodesByTask(ctx context.Context, taskID uuid.UUID) ([]model.TransportNode, error) {
	var nodes []model.TransportNode
	err := r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Order("created_at asc").
		Find(&nodes).Error
	return nodes, err
}

func (r *transportRepository) CreateReceipt(ctx context.Context, receipt *model.Receipt) error {
	return r.db.WithContext(ctx).Create(receipt).Error
}

func (r *transportRepository) FindReceiptByTask(ctx context.Context, taskID uuid.UUID) (*model.Receipt, error) {
	var receipt model.Receipt
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).First(&receipt).Error
	return &receipt, err
}
