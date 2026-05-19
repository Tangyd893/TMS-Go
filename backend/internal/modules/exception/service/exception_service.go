package service

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/exception/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrExceptionNotFound = errors.New(404001, "异常记录不存在")
)

type ExceptionService interface {
	List(ctx context.Context, req pagination.PageRequest, exceptionType, status string) (*pagination.PageResult[model.Exception], error)
	GetByID(ctx context.Context, id string) (*model.Exception, error)
	Create(ctx context.Context, req model.CreateExceptionRequest, userID, userName string) (*model.Exception, error)
	Handle(ctx context.Context, id, handlerID, handlerName, result string) error
	Close(ctx context.Context, id string) error
	GetLogs(ctx context.Context, exceptionID string) ([]model.ExceptionLog, error)
}

type exceptionService struct {
	repo repository.ExceptionRepository
}

func NewExceptionService(repo repository.ExceptionRepository) ExceptionService {
	return &exceptionService{repo: repo}
}

func (s *exceptionService) List(ctx context.Context, req pagination.PageRequest, exceptionType, status string) (*pagination.PageResult[model.Exception], error) {
	if req.Page < 1 { req.Page = 1 }
	if req.PageSize < 1 || req.PageSize > 100 { req.PageSize = 20 }
	items, total, err := s.repo.List(ctx, req.Page, req.PageSize, exceptionType, status)
	if err != nil { return nil, err }
	return &pagination.PageResult[model.Exception]{Items: items, Total: total, Page: req.Page, PageSize: req.PageSize}, nil
}

func (s *exceptionService) GetByID(ctx context.Context, id string) (*model.Exception, error) {
	uid, err := uuid.Parse(id)
	if err != nil { return nil, ErrExceptionNotFound }
	e, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrExceptionNotFound }
		return nil, err
	}
	return e, nil
}

func (s *exceptionService) Create(ctx context.Context, req model.CreateExceptionRequest, userID, userName string) (*model.Exception, error) {
	now := time.Now()
	exceptionNo := fmt.Sprintf("EX%s%04d", now.Format("20060102150405"), now.Nanosecond()%10000)

	e := &model.Exception{
		TaskID:        req.TaskID,
		OrderID:       req.OrderID,
		ExceptionNo:   exceptionNo,
		ExceptionType: req.ExceptionType,
		Description:   req.Description,
		Severity:      req.Severity,
		ReportBy:      uuid.MustParse(userID),
		ReportByName:  userName,
		Status:        model.ExceptionStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if e.Severity == "" { e.Severity = model.ExceptionSeverityNormal }

	if err := s.repo.Create(ctx, e); err != nil { return nil, err }
	return s.repo.FindByID(ctx, e.ID)
}

func (s *exceptionService) Handle(ctx context.Context, id, handlerID, handlerName, result string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return ErrExceptionNotFound }
	e, err := s.repo.FindByID(ctx, uid)
	if err != nil { return ErrExceptionNotFound }
	if e.Status != model.ExceptionStatusPending {
		return errors.New(640001, "仅待处理异常可进行处置")
	}
	hid := uuid.MustParse(handlerID)
	if err := s.repo.Handle(ctx, uid, hid, handlerName, result); err != nil { return err }
	return s.repo.AddLog(ctx, &model.ExceptionLog{
		ExceptionID: uid,
		Action:      "handle",
		Content:     result,
		OperatorID:  &hid,
		OperatorName: handlerName,
	})
}

func (s *exceptionService) Close(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil { return ErrExceptionNotFound }
	e, err := s.repo.FindByID(ctx, uid)
	if err != nil { return ErrExceptionNotFound }
	if e.Status != model.ExceptionStatusResolved {
		return errors.New(640002, "仅已解决异常可关闭")
	}
	return s.repo.Close(ctx, uid)
}

func (s *exceptionService) GetLogs(ctx context.Context, exceptionID string) ([]model.ExceptionLog, error) {
	uid, err := uuid.Parse(exceptionID)
	if err != nil { return nil, ErrExceptionNotFound }
	return s.repo.FindLogs(ctx, uid)
}
