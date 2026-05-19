package dict

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/google/uuid"
)

var (
	ErrDictTypeNotFound = errors.New(404001, "字典类型不存在")
	ErrDictItemNotFound = errors.New(404001, "字典项不存在")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListTypes(ctx context.Context) ([]DictType, error) {
	return s.repo.ListTypes(ctx)
}

func (s *Service) GetType(ctx context.Context, code string) (*DictType, error) {
	t, err := s.repo.FindTypeByCode(ctx, code)
	if err != nil {
		return nil, ErrDictTypeNotFound
	}
	return t, nil
}

func (s *Service) CreateType(ctx context.Context, t *DictType) (*DictType, error) {
	if err := s.repo.CreateType(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) UpdateType(ctx context.Context, t *DictType) (*DictType, error) {
	if err := s.repo.UpdateType(ctx, t); err != nil {
		return nil, err
	}
	return s.repo.FindTypeByCode(ctx, t.Code)
}

func (s *Service) DeleteType(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrDictTypeNotFound
	}
	return s.repo.DeleteType(ctx, uid)
}

func (s *Service) ListItems(ctx context.Context, typeCode string) ([]DictItem, error) {
	return s.repo.ListItems(ctx, typeCode)
}

func (s *Service) CreateItem(ctx context.Context, item *DictItem) (*DictItem, error) {
	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateItem(ctx context.Context, item *DictItem) (*DictItem, error) {
	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, err
	}
	items, _ := s.repo.ListItems(ctx, item.TypeCode)
	for i := range items {
		if items[i].ID == item.ID {
			return &items[i], nil
		}
	}
	return item, nil
}

func (s *Service) DeleteItem(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrDictItemNotFound
	}
	return s.repo.DeleteItem(ctx, uid)
}
