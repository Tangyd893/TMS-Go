package service

import (
	"context"
	"encoding/json"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/base/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
)

var (
	ErrEntityNotFound  = errors.New(404001, "记录不存在")
	ErrEntityDuplicate = errors.New(409001, "记录已存在")
)

type EntityService[T any] interface {
	List(ctx context.Context, req pagination.PageRequest, keyword string) (*pagination.PageResult[T], error)
	GetByID(ctx context.Context, id string) (*T, error)
	Create(ctx context.Context, js json.RawMessage) (*T, error)
	Update(ctx context.Context, id string, js json.RawMessage) (*T, error)
	Delete(ctx context.Context, id string) error
}

type entityService[T any] struct {
	repo    repository.CrudRepository[T]
	searchFields []string
}

func NewEntityService[T any](repo repository.CrudRepository[T], searchFields []string) EntityService[T] {
	return &entityService[T]{repo: repo, searchFields: searchFields}
}

func (s *entityService[T]) List(ctx context.Context, req pagination.PageRequest, keyword string) (*pagination.PageResult[T], error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	entities, total, err := s.repo.List(ctx, req.Page, req.PageSize, keyword, s.searchFields)
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[T]{
		Items:    entities,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *entityService[T]) GetByID(ctx context.Context, id string) (*T, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	entity, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return entity, nil
}

func (s *entityService[T]) Create(ctx context.Context, js json.RawMessage) (*T, error) {
	var entity T
	if err := json.Unmarshal(js, &entity); err != nil {
		return nil, errors.New(400001, "请求参数错误")
	}
	if err := s.repo.Create(ctx, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (s *entityService[T]) Update(ctx context.Context, id string, js json.RawMessage) (*T, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrEntityNotFound
	}

	var updates map[string]any
	if err := json.Unmarshal(js, &updates); err != nil {
		return nil, errors.New(400001, "请求参数错误")
	}

	if err := s.repo.Update(ctx, uid, updates); err != nil {
		return nil, err
	}

	entity, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, ErrEntityNotFound
	}
	return entity, nil
}

func (s *entityService[T]) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrEntityNotFound
	}
	return s.repo.Delete(ctx, uid)
}
