package service

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
)

var (
	ErrUserNotFound      = errors.New(404001, "用户不存在")
	ErrUsernameDuplicate = errors.New(409001, "用户名已存在")
)

type UserService interface {
	List(ctx context.Context, req pagination.PageRequest, keyword string) (*pagination.PageResult[model.User], error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Create(ctx context.Context, params model.CreateUserParams) (*model.User, error)
	Update(ctx context.Context, id string, params model.CreateUserParams) (*model.User, error)
	Delete(ctx context.Context, id string) error
}
