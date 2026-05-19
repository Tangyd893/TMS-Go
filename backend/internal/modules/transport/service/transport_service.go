package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/transport/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound = errors.New(404001, "运输任务不存在")
)

type TransportService interface {
	List(ctx context.Context, req pagination.PageRequest, keyword, status string) (*pagination.PageResult[model.TransportTask], error)
	GetByID(ctx context.Context, id string) (*model.TransportTask, error)
	Create(ctx context.Context, task *model.TransportTask) (*model.TransportTask, error)
	Depart(ctx context.Context, id string) (*model.TransportTask, error)
	Arrive(ctx context.Context, id string) (*model.TransportTask, error)
	Sign(ctx context.Context, id string) (*model.TransportTask, error)

	// 节点
	AddNode(ctx context.Context, taskID string, req model.CreateNodeRequest, userID string) (*model.TransportNode, error)
	GetNodes(ctx context.Context, taskID string) ([]model.TransportNode, error)

	// 回单
	CreateReceipt(ctx context.Context, taskID string, req model.CreateReceiptRequest, userID string) (*model.Receipt, error)
	GetReceipt(ctx context.Context, taskID string) (*model.Receipt, error)
}

type transportService struct {
	repo repository.TransportRepository
}

func NewTransportService(repo repository.TransportRepository) TransportService {
	return &transportService{repo: repo}
}

func (s *transportService) List(ctx context.Context, req pagination.PageRequest, keyword, status string) (*pagination.PageResult[model.TransportTask], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	tasks, total, err := s.repo.List(ctx, req.Page, req.PageSize, keyword, status)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.TransportTask]{Items: tasks, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *transportService) GetByID(ctx context.Context, id string) (*model.TransportTask, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, ErrTaskNotFound }
	task, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrTaskNotFound }
		return nil, err
	}
	return task, nil
}

func (s *transportService) Create(ctx context.Context, task *model.TransportTask) (*model.TransportTask, error) {
	now := time.Now()
	task.TaskNo = fmt.Sprintf("TT%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)
	task.Status = model.TaskStatusPending
	task.CreatedAt = now
	task.UpdatedAt = now
	if err := s.repo.Create(ctx, task); err != nil { return nil, err }
	return s.repo.FindByID(ctx, task.ID)
}

func (s *transportService) Depart(ctx context.Context, id string) (*model.TransportTask, error) {
	uid, _ := uuid.Parse(id)
	task, err := s.repo.FindByID(ctx, uid)
	if err != nil { return nil, ErrTaskNotFound }
	if task.Status != model.TaskStatusPending {
		return nil, errors.New(630001, "运输任务不能重复发车")
	}
	if err := s.repo.UpdateStatus(ctx, uid, model.TaskStatusDeparted); err != nil { return nil, err }
	task.Status = model.TaskStatusDeparted
	return task, nil
}

func (s *transportService) Arrive(ctx context.Context, id string) (*model.TransportTask, error) {
	uid, _ := uuid.Parse(id)
	task, err := s.repo.FindByID(ctx, uid)
	if err != nil { return nil, ErrTaskNotFound }
	if task.Status != model.TaskStatusDeparted && task.Status != model.TaskStatusInTransit {
		return nil, errors.New(630002, "运输任务状态不允许到达")
	}
	if err := s.repo.UpdateStatus(ctx, uid, model.TaskStatusArrived); err != nil { return nil, err }
	task.Status = model.TaskStatusArrived
	return task, nil
}

func (s *transportService) Sign(ctx context.Context, id string) (*model.TransportTask, error) {
	uid, _ := uuid.Parse(id)
	task, err := s.repo.FindByID(ctx, uid)
	if err != nil { return nil, ErrTaskNotFound }
	if task.Status != model.TaskStatusArrived {
		return nil, errors.New(630003, "任务未到达，不能签收")
	}
	if err := s.repo.UpdateStatus(ctx, uid, model.TaskStatusSigned); err != nil { return nil, err }
	task.Status = model.TaskStatusSigned
	return task, nil
}

func (s *transportService) AddNode(ctx context.Context, taskID string, req model.CreateNodeRequest, userID string) (*model.TransportNode, error) {
	taskUID, err := uuid.Parse(taskID)
	if err != nil { return nil, ErrTaskNotFound }
	_, err = s.repo.FindByID(ctx, taskUID)
	if err != nil { return nil, ErrTaskNotFound }

	node := &model.TransportNode{
		TaskID:       taskUID,
		NodeType:     req.NodeType,
		NodeName:     req.NodeName,
		LocationName: req.LocationName,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		ArrivedAt:    req.ArrivedAt,
		DepartedAt:   req.DepartedAt,
		Remark:       req.Remark,
		CreatedAt:    time.Now(),
	}

	if userID != "" {
		uid := uuid.MustParse(userID)
		node.CreatedBy = &uid
	}

	if err := s.repo.CreateNode(ctx, node); err != nil { return nil, err }
	return node, nil
}

func (s *transportService) GetNodes(ctx context.Context, taskID string) ([]model.TransportNode, error) {
	taskUID, err := uuid.Parse(taskID)
	if err != nil { return nil, ErrTaskNotFound }
	return s.repo.FindNodesByTask(ctx, taskUID)
}

func (s *transportService) CreateReceipt(ctx context.Context, taskID string, req model.CreateReceiptRequest, userID string) (*model.Receipt, error) {
	taskUID, err := uuid.Parse(taskID)
	if err != nil { return nil, ErrTaskNotFound }
	task, err := s.repo.FindByID(ctx, taskUID)
	if err != nil { return nil, ErrTaskNotFound }

	now := time.Now()
	receiptNo := fmt.Sprintf("RT%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)

	receipt := &model.Receipt{
		TaskID:       taskUID,
		OrderID:      task.OrderID,
		ReceiptNo:    receiptNo,
		SignBy:       req.SignBy,
		SignImageURL: req.SignImageURL,
		Remark:       req.Remark,
		SignAt:       &now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if userID != "" {
		uid := uuid.MustParse(userID)
		receipt.CreatedBy = &uid
	}

	if err := s.repo.CreateReceipt(ctx, receipt); err != nil { return nil, err }
	return receipt, nil
}

func (s *transportService) GetReceipt(ctx context.Context, taskID string) (*model.Receipt, error) {
	taskUID, err := uuid.Parse(taskID)
	if err != nil { return nil, ErrTaskNotFound }
	return s.repo.FindReceiptByTask(ctx, taskUID)
}
