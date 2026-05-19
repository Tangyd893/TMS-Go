package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/crypto"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/jwt"
	"github.com/google/uuid"
)

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.Manager
	log        *slog.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.Manager, log *slog.Logger) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
		log:        log,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		s.log.Warn("login failed: user not found", slog.String("username", username))
		return nil, ErrInvalidCredentials
	}

	if user.Status != "active" {
		return nil, ErrUserDisabled
	}

	if !crypto.CheckPassword(password, user.PasswordHash) {
		s.log.Warn("login failed: invalid password", slog.String("username", username))
		return nil, ErrInvalidCredentials
	}

	userIDStr := user.ID.String()
	accessToken, err := s.jwtManager.GenerateAccessToken(userIDStr, user.Username)
	if err != nil {
		s.log.Error("failed to generate access token", slog.Any("error", err))
		return nil, errors.New(500001, "系统错误")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userIDStr, user.Username)
	if err != nil {
		s.log.Error("failed to generate refresh token", slog.Any("error", err))
		return nil, errors.New(500001, "系统错误")
	}

	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, time.Now())

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         buildUserInfo(user),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshTokenStr string) (*LoginResult, error) {
	claims, err := s.jwtManager.ParseRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	if user.Status != "active" {
		return nil, ErrUserDisabled
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID.String(), user.Username)
	if err != nil {
		return nil, errors.New(500001, "系统错误")
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID.String(), user.Username)
	if err != nil {
		return nil, errors.New(500001, "系统错误")
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         buildUserInfo(user),
	}, nil
}

func (s *authService) GetCurrentUser(ctx context.Context, userID string) (*model.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrTokenInvalid
	}
	return s.userRepo.FindByID(ctx, id)
}

func buildUserInfo(user *model.User) *UserInfo {
	permissions := make(map[string]bool)
	roles := make(map[string]bool)

	for _, role := range user.Roles {
		roles[role.Code] = true
		for _, perm := range role.Permissions {
			permissions[perm.Code] = true
		}
	}

	permList := make([]string, 0, len(permissions))
	for p := range permissions {
		permList = append(permList, p)
	}

	roleList := make([]string, 0, len(roles))
	for r := range roles {
		roleList = append(roleList, r)
	}

	return &UserInfo{
		ID:          user.ID.String(),
		Username:    user.Username,
		RealName:    user.RealName,
		AvatarURL:   user.AvatarURL,
		Permissions: permList,
		Roles:       roleList,
	}
}
