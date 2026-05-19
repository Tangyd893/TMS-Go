package service

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
)

var (
	ErrInvalidCredentials = errors.New(401001, "用户名或密码错误")
	ErrUserDisabled       = errors.New(401002, "用户已被禁用")
	ErrTokenExpired       = errors.New(401003, "令牌已过期")
	ErrTokenInvalid       = errors.New(401004, "无效令牌")
)

type LoginResult struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	RealName    string   `json:"realName"`
	AvatarURL   string   `json:"avatarUrl"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error)
	GetCurrentUser(ctx context.Context, userID string) (*model.User, error)
}
