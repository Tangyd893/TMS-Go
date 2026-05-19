package service

import (
	"context"
	"errors"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/crypto"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) List(ctx context.Context, req pagination.PageRequest, keyword string) (*pagination.PageResult[model.User], error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	users, total, err := s.repo.List(ctx, req.Page, req.PageSize, keyword)
	if err != nil {
		return nil, err
	}

	return &pagination.PageResult[model.User]{
		Items:    users,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	user, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, repository.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Create(ctx context.Context, params model.CreateUserParams) (*model.User, error) {
	existing, err := s.repo.FindByUsername(ctx, params.Username)
	if err == nil && existing != nil {
		return nil, ErrUsernameDuplicate
	}

	hash, err := crypto.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     params.Username,
		PasswordHash: hash,
		RealName:     params.RealName,
		Phone:        params.Phone,
		Email:        params.Email,
		Status:       params.Status,
	}

	if user.Status == "" {
		user.Status = "active"
	}

	if err := s.repo.Create(ctx, user, params.RoleIDs); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, user.ID)
}

func (s *userService) Update(ctx context.Context, id string, params model.CreateUserParams) (*model.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user, err := s.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if params.Username != "" && params.Username != user.Username {
		existing, _ := s.repo.FindByUsername(ctx, params.Username)
		if existing != nil {
			return nil, ErrUsernameDuplicate
		}
		user.Username = params.Username
	}

	user.RealName = params.RealName
	user.Phone = params.Phone
	user.Email = params.Email
	if params.Status != "" {
		user.Status = params.Status
	}

	if err := s.repo.Update(ctx, user, params.RoleIDs); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, uid)
}

func (s *userService) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return ErrUserNotFound
	}
	return s.repo.Delete(ctx, uid)
}
