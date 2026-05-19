# M1 认证权限系统 — 详细实施计划 ✅ 已完成

> **状态：已完成。** 所有 Task (M1-1 ~ M1-15) 已于 2026-05-19 前实施完毕。本文档保留作为实施记录参考。

**Goal:** 实现完整的登录认证、JWT 鉴权、用户/角色/权限/菜单 CRUD、前端动态路由和权限守卫。

**Architecture:** 后端 `internal/modules/auth/` + `internal/modules/system/` 模块，共用 `internal/platform/database/` 数据库连接，前端 `stores/auth.ts` + `router/guards.ts` + `views/system/` 页面。

**Tech Stack:** Go 1.22 + net/http + GORM + PostgreSQL + bcrypt + golang-jwt + Redis（可选，用于 token 黑名单）

---

### Task M1-1: 添加 Go 依赖并扩展配置

**Files:**
- Modify: `backend/go.mod:1-3`
- Modify: `backend/internal/config/config.go:1-33`
- Create: `backend/internal/config/yaml.go`

- [ ] **Step 1: 添加依赖**

在 `backend/` 目录执行：
```bash
cd backend
go get gorm.io/gorm gorm.io/driver/postgres golang.org/x/crypto github.com/golang-jwt/jwt/v5 gopkg.in/yaml.v3 github.com/google/uuid
go mod tidy
```

- [ ] **Step 2: 扩展 Config 结构体**

修改 `backend/internal/config/config.go`：

```go
package config

import (
	"os"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
	RabbitMQ RabbitMQConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret    string
	RefreshSecret   string
	AccessTokenTTL  string
	RefreshTokenTTL string
	Issuer          string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type RabbitMQConfig struct {
	URL string
}

func Load() Config {
	return Config{
		App: AppConfig{
			Name: envOrDefault("APP_NAME", "tms-api"),
			Env:  envOrDefault("APP_ENV", "dev"),
			Port: envOrDefault("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:         envOrDefault("DB_HOST", "localhost"),
			Port:         envOrDefault("DB_PORT", "5432"),
			Name:         envOrDefault("DB_NAME", "tms"),
			User:         envOrDefault("DB_USER", "tms"),
			Password:     envOrDefault("DB_PASSWORD", "change_me"),
			SSLMode:      envOrDefault("DB_SSLMODE", "disable"),
			MaxOpenConns: 50,
			MaxIdleConns: 10,
		},
		Redis: RedisConfig{
			Host:     envOrDefault("REDIS_HOST", "localhost"),
			Port:     envOrDefault("REDIS_PORT", "6379"),
			Password: envOrDefault("REDIS_PASSWORD", ""),
			DB:       0,
		},
		JWT: JWTConfig{
			AccessSecret:    envOrDefault("JWT_ACCESS_SECRET", "tms-access-secret-change-me"),
			RefreshSecret:   envOrDefault("JWT_REFRESH_SECRET", "tms-refresh-secret-change-me"),
			AccessTokenTTL:  envOrDefault("JWT_ACCESS_TTL", "2h"),
			RefreshTokenTTL: envOrDefault("JWT_REFRESH_TTL", "168h"),
			Issuer:          envOrDefault("JWT_ISSUER", "tms"),
		},
		MinIO: MinIOConfig{
			Endpoint:  envOrDefault("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: envOrDefault("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: envOrDefault("MINIO_SECRET_KEY", "change_me"),
			Bucket:    envOrDefault("MINIO_BUCKET", "tms-dev"),
			UseSSL:    false,
		},
		RabbitMQ: RabbitMQConfig{
			URL: envOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		},
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
```

- [ ] **Step 3: 验证编译通过**

```bash
cd backend && go build ./...
```

---

### Task M1-2: 创建数据库连接初始化

**Files:**
- Create: `backend/internal/platform/database/database.go`

- [ ] **Step 1: 创建数据库初始化**

```go
package database

import (
	"fmt"
	"log/slog"

	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func New(cfg config.DatabaseConfig, log *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	log.Info("database connection established", slog.String("host", cfg.Host), slog.String("db", cfg.Name))

	return db, nil
}
```

- [ ] **Step 2: 更新 main.go 初始化数据库**

修改 `backend/cmd/server/main.go`：

```go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/bootstrap"
	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"github.com/Tangyd893/TMS-Go/backend/internal/platform/database"
	"github.com/Tangyd893/TMS-Go/backend/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.Env)

	db, err := database.New(cfg.Database, log)
	if err != nil {
		log.Error("failed to initialize database", slog.Any("error", err))
		os.Exit(1)
	}

	server := bootstrap.NewHTTPServer(cfg, log, db)

	go func() {
		log.Info("starting http server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server stopped unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("shutting down http server")
	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown http server", slog.Any("error", err))
		os.Exit(1)
	}
}
```

- [ ] **Step 3: 更新 bootstrap/http.go 签名**

```go
package bootstrap

import (
	"log/slog"
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	healthapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/api"
	healthsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/service"
	"gorm.io/gorm"
)

func NewHTTPServer(cfg config.Config, log *slog.Logger, db *gorm.DB) *http.Server {
	mux := http.NewServeMux()

	healthService := healthsvc.NewService()
	healthHandler := healthapi.NewHandler(healthService)
	healthapi.RegisterRoutes(mux, healthHandler)

	handler := middleware.RequestID(middleware.AccessLog(log)(mux))

	return &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: handler,
	}
}
```

- [ ] **Step 4: 编译验证**

```bash
cd backend && go build ./...
```

---

### Task M1-3: 创建数据库迁移文件

**Files:**
- Create: `backend/migrations/000002_create_sys_user.up.sql`
- Create: `backend/migrations/000002_create_sys_user.down.sql`
- Create: `backend/migrations/000003_create_sys_role.up.sql`
- Create: `backend/migrations/000003_create_sys_role.down.sql`
- Create: `backend/migrations/000004_create_sys_permission.up.sql`
- Create: `backend/migrations/000004_create_sys_permission.down.sql`
- Create: `backend/migrations/000005_create_sys_user_role.up.sql`
- Create: `backend/migrations/000005_create_sys_user_role.down.sql`
- Create: `backend/migrations/000006_create_sys_role_permission.up.sql`
- Create: `backend/migrations/000006_create_sys_role_permission.down.sql`
- Create: `backend/migrations/000007_create_sys_menu.up.sql`
- Create: `backend/migrations/000007_create_sys_menu.down.sql`

- [ ] **Step 1: 用户表迁移**

`backend/migrations/000002_create_sys_user.up.sql`：
```sql
create table if not exists sys_user (
  id uuid primary key default uuid_generate_v4(),
  username varchar(64) not null,
  password_hash varchar(256) not null,
  real_name varchar(64),
  phone varchar(32),
  email varchar(128),
  avatar_url varchar(512),
  status varchar(16) not null default 'active',
  org_id uuid,
  last_login_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_sys_user_username on sys_user(username) where deleted_at is null;
create index idx_sys_user_status on sys_user(status);
comment on table sys_user is '系统用户';
comment on column sys_user.status is 'active: 启用, disabled: 禁用';
```

`backend/migrations/000002_create_sys_user.down.sql`：
```sql
drop table if exists sys_user;
```

- [ ] **Step 2: 角色表迁移**

`backend/migrations/000003_create_sys_role.up.sql`：
```sql
create table if not exists sys_role (
  id uuid primary key default uuid_generate_v4(),
  code varchar(64) not null,
  name varchar(128) not null,
  remark varchar(512),
  status varchar(16) not null default 'active',
  sort_no int not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_sys_role_code on sys_role(code) where deleted_at is null;
comment on table sys_role is '系统角色';
```

`backend/migrations/000003_create_sys_role.down.sql`：
```sql
drop table if exists sys_role;
```

- [ ] **Step 3: 权限表迁移**

`backend/migrations/000004_create_sys_permission.up.sql`：
```sql
create table if not exists sys_permission (
  id uuid primary key default uuid_generate_v4(),
  code varchar(128) not null,
  name varchar(128) not null,
  type varchar(32) not null,
  parent_id uuid,
  sort_no int not null default 0,
  remark varchar(512),
  status varchar(16) not null default 'active',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create unique index uk_sys_permission_code on sys_permission(code) where deleted_at is null;
comment on table sys_permission is '系统权限';
comment on column sys_permission.type is 'menu: 菜单, button: 按钮, api: 接口';
```

`backend/migrations/000004_create_sys_permission.down.sql`：
```sql
drop table if exists sys_permission;
```

- [ ] **Step 4: 用户角色关联表迁移**

`backend/migrations/000005_create_sys_user_role.up.sql`：
```sql
create table if not exists sys_user_role (
  id uuid primary key default uuid_generate_v4(),
  user_id uuid not null,
  role_id uuid not null,
  created_at timestamptz not null default now()
);
create unique index uk_sys_user_role on sys_user_role(user_id, role_id);
create index idx_sys_user_role_user_id on sys_user_role(user_id);
create index idx_sys_user_role_role_id on sys_user_role(role_id);
comment on table sys_user_role is '用户角色关联';
```

`backend/migrations/000005_create_sys_user_role.down.sql`：
```sql
drop table if exists sys_user_role;
```

- [ ] **Step 5: 角色权限关联表迁移**

`backend/migrations/000006_create_sys_role_permission.up.sql`：
```sql
create table if not exists sys_role_permission (
  id uuid primary key default uuid_generate_v4(),
  role_id uuid not null,
  permission_id uuid not null,
  created_at timestamptz not null default now()
);
create unique index uk_sys_role_permission on sys_role_permission(role_id, permission_id);
comment on table sys_role_permission is '角色权限关联';
```

`backend/migrations/000006_create_sys_role_permission.down.sql`：
```sql
drop table if exists sys_role_permission;
```

- [ ] **Step 6: 菜单表迁移**

`backend/migrations/000007_create_sys_menu.up.sql`：
```sql
create table if not exists sys_menu (
  id uuid primary key default uuid_generate_v4(),
  parent_id uuid,
  name varchar(128) not null,
  path varchar(256),
  component varchar(256),
  icon varchar(64),
  permission_code varchar(128),
  type varchar(16) not null default 'menu',
  sort_no int not null default 0,
  visible boolean not null default true,
  status varchar(16) not null default 'active',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  deleted_at timestamptz
);
create index idx_sys_menu_parent_id on sys_menu(parent_id);
create index idx_sys_menu_sort_no on sys_menu(sort_no);
comment on table sys_menu is '系统菜单';
comment on column sys_menu.type is 'directory: 目录, menu: 菜单, button: 按钮';
```

`backend/migrations/000007_create_sys_menu.down.sql`：
```sql
drop table if exists sys_menu;
```

- [ ] **Step 7: 执行迁移**

```bash
# 使用 golang-migrate 工具执行迁移
# 或在 docker-compose 中连接 PostgreSQL 手动执行
docker compose -f docker/docker-compose.yml exec postgres psql -U tms -d tms -c "create extension if not exists \"uuid-ossp\";"
```

---

### Task M1-4: 创建 JWT 工具

**Files:**
- Create: `backend/internal/shared/jwt/jwt.go`

- [ ] **Step 1: JWT 工具实现**

```go
package jwt

import (
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwtlib.RegisteredClaims
}

type Manager struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	issuer        string
}

func NewManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration, issuer string) *Manager {
	return &Manager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		issuer:        issuer,
	}
}

func (m *Manager) GenerateAccessToken(userID, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.accessTTL)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.accessSecret))
}

func (m *Manager) GenerateRefreshToken(userID, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(m.refreshTTL)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.refreshSecret))
}

func (m *Manager) ParseAccessToken(tokenStr string) (*Claims, error) {
	return m.parseToken(tokenStr, m.accessSecret)
}

func (m *Manager) ParseRefreshToken(tokenStr string) (*Claims, error) {
	return m.parseToken(tokenStr, m.refreshSecret)
}

func (m *Manager) parseToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &Claims{}, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

---

### Task M1-5: 创建密码工具

**Files:**
- Create: `backend/internal/shared/crypto/password.go`

- [ ] **Step 1: 密码哈希实现**

```go
package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

---

### Task M1-6: 认证模块 — Model

**Files:**
- Create: `backend/internal/modules/auth/model/user.go`

- [ ] **Step 1: User 实体**

```go
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex:uk_sys_user_username,where:deleted_at is null" json:"username"`
	PasswordHash string         `gorm:"size:256;not null" json:"-"`
	RealName     string         `gorm:"size:64" json:"realName"`
	Phone        string         `gorm:"size:32" json:"phone"`
	Email        string         `gorm:"size:128" json:"email"`
	AvatarURL    string         `gorm:"size:512" json:"avatarUrl"`
	Status       string         `gorm:"size:16;not null;default:active" json:"status"`
	OrgID        *uuid.UUID     `gorm:"type:uuid" json:"orgId"`
	LastLoginAt  *time.Time     `json:"lastLoginAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Roles        []Role         `gorm:"many2many:sys_user_role;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "sys_user"
}

type Role struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code         string          `gorm:"size:64;not null;uniqueIndex:uk_sys_role_code,where:deleted_at is null" json:"code"`
	Name         string          `gorm:"size:128;not null" json:"name"`
	Remark       string          `gorm:"size:512" json:"remark"`
	Status       string          `gorm:"size:16;not null;default:active" json:"status"`
	SortNo       int             `gorm:"not null;default:0" json:"sortNo"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
	Permissions  []Permission    `gorm:"many2many:sys_role_permission;" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "sys_role"
}

type Permission struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code      string         `gorm:"size:128;not null;uniqueIndex:uk_sys_permission_code,where:deleted_at is null" json:"code"`
	Name      string         `gorm:"size:128;not null" json:"name"`
	Type      string         `gorm:"size:32;not null" json:"type"`
	ParentID  *uuid.UUID     `gorm:"type:uuid" json:"parentId"`
	SortNo    int            `gorm:"not null;default:0" json:"sortNo"`
	Remark    string         `gorm:"size:512" json:"remark"`
	Status    string         `gorm:"size:16;not null;default:active" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "sys_permission"
}

type Menu struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ParentID       *uuid.UUID     `gorm:"type:uuid" json:"parentId"`
	Name           string         `gorm:"size:128;not null" json:"name"`
	Path           string         `gorm:"size:256" json:"path"`
	Component      string         `gorm:"size:256" json:"component"`
	Icon           string         `gorm:"size:64" json:"icon"`
	PermissionCode string         `gorm:"size:128" json:"permissionCode"`
	Type           string         `gorm:"size:16;not null;default:menu" json:"type"`
	SortNo         int            `gorm:"not null;default:0" json:"sortNo"`
	Visible        bool           `gorm:"not null;default:true" json:"visible"`
	Status         string         `gorm:"size:16;not null;default:active" json:"status"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Children       []Menu         `gorm:"-" json:"children,omitempty"`
}

func (Menu) TableName() string {
	return "sys_menu"
}
```

---

### Task M1-7: 认证模块 — Repository

**Files:**
- Create: `backend/internal/modules/auth/repository/user_repository.go`

- [ ] **Step 1: UserRepository 接口与实现**

```go
package repository

import (
	"context"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("username = ?", username).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).
		Preload("Roles.Permissions").
		Where("id = ?", id).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, loginTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("last_login_at", loginTime).Error
}
```

---

### Task M1-8: 认证模块 — Service

**Files:**
- Create: `backend/internal/modules/auth/service/service.go`
- Create: `backend/internal/modules/auth/service/service_impl.go`

- [ ] **Step 1: 认证服务接口**

```go
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
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         *UserInfo   `json:"user"`
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
```

- [ ] **Step 2: 认证服务实现**

```go
package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/crypto"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/jwt"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/model"
	"github.com/google/uuid"
)

type authService struct {
	userRepo    repository.UserRepository
	jwtManager  *jwt.Manager
	log         *slog.Logger
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
```

- [ ] **Step 3: 编译验证**

```bash
cd backend && go build ./...
```

---

### Task M1-9: 认证模块 — API

**Files:**
- Create: `backend/internal/modules/auth/api/handler.go`
- Create: `backend/internal/modules/auth/api/dto.go`
- Create: `backend/internal/modules/auth/api/routes.go`

- [ ] **Step 1: DTO**

```go
package api

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}
```

- [ ] **Step 2: Handler**

```go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type Handler struct {
	authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}

	if req.Username == "" || req.Password == "" {
		response.Fail(w, r, http.StatusBadRequest, 400001, "用户名和密码不能为空", nil)
		return
	}

	result, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		appErr := unwrapAppError(err)
		response.Fail(w, r, http.StatusUnauthorized, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, result)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, r, http.StatusBadRequest, 400001, "请求参数错误", nil)
		return
	}

	result, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		appErr := unwrapAppError(err)
		response.Fail(w, r, http.StatusUnauthorized, appErr.Code, appErr.Message, nil)
		return
	}

	response.Success(w, r, result)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(string)

	user, err := h.authService.GetCurrentUser(r.Context(), userID)
	if err != nil {
		response.Fail(w, r, http.StatusUnauthorized, 401004, "获取用户信息失败", nil)
		return
	}

	response.Success(w, r, toUserResponse(user))
}

func unwrapAppError(err error) *errors.AppError {
	if e, ok := err.(*errors.AppError); ok {
		return e
	}
	return errors.New(500001, err.Error())
}

func toUserResponse(user interface{}) interface{} {
	return user
}
```

- [ ] **Step 3: Routes**

```go
package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /api/v1/auth/login", handler.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", handler.Refresh)
	mux.HandleFunc("GET /api/v1/auth/me", handler.GetCurrentUser)
}
```

---

### Task M1-10: 认证中间件

**Files:**
- Create: `backend/internal/middleware/auth.go`

- [ ] **Step 1: JWT 认证中间件**

```go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Tangyd893/TMS-Go/backend/internal/shared/jwt"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/response"
)

type userIDKey struct{}

func Auth(jwtManager *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Fail(w, r, http.StatusUnauthorized, 401001, "未登录或令牌无效", nil)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Fail(w, r, http.StatusUnauthorized, 401001, "未登录或令牌无效", nil)
				return
			}

			claims, err := jwtManager.ParseAccessToken(parts[1])
			if err != nil {
				response.Fail(w, r, http.StatusUnauthorized, 401003, "令牌已过期", nil)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey{}, claims.UserID)
			ctx = context.WithValue(ctx, "username", claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) string {
	value, ok := ctx.Value(userIDKey{}).(string)
	if !ok {
		return ""
	}
	return value
}
```

---

### Task M1-11: 更新 bootstrap — 注册所有新路由

**Files:**
- Modify: `backend/internal/bootstrap/http.go`

- [ ] **Step 1: 整合认证模块到 bootstrap**

```go
package bootstrap

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"github.com/Tangyd893/TMS-Go/backend/internal/middleware"
	authapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/api"
	authrepo "github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/repository"
	authsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/auth/service"
	healthapi "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/api"
	healthsvc "github.com/Tangyd893/TMS-Go/backend/internal/modules/health/service"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/jwt"
	"gorm.io/gorm"
)

func NewHTTPServer(cfg config.Config, log *slog.Logger, db *gorm.DB) *http.Server {
	accessTTL, _ := time.ParseDuration(cfg.JWT.AccessTokenTTL)
	if accessTTL == 0 {
		accessTTL = 2 * time.Hour
	}
	refreshTTL, _ := time.ParseDuration(cfg.JWT.RefreshTokenTTL)
	if refreshTTL == 0 {
		refreshTTL = 168 * time.Hour
	}

	jwtManager := jwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		accessTTL,
		refreshTTL,
		cfg.JWT.Issuer,
	)

	authMiddleware := middleware.Auth(jwtManager)

	userRepo := authrepo.NewUserRepository(db)
	authService := authsvc.NewAuthService(userRepo, jwtManager, log)
	authHandler := authapi.NewHandler(authService)

	mux := http.NewServeMux()

	healthService := healthsvc.NewService()
	healthHandler := healthapi.NewHandler(healthService)
	healthapi.RegisterRoutes(mux, healthHandler)

	authapi.RegisterRoutes(mux, authHandler)

	handler := middleware.RequestID(middleware.AccessLog(log)(mux))
	_ = authMiddleware // 后续 Task 会绑到受保护路由上

	return &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: handler,
	}
}
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```

---

### Task M1-12: 前端 — 重写登录页

**Files:**
- Modify: `frontend/src/views/login/LoginView.vue`
- Create: `frontend/src/api/modules/auth.ts`

- [ ] **Step 1: 创建 API 模块**

`frontend/src/api/modules/auth.ts`：
```ts
import { request } from '../request'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResult {
  accessToken: string
  refreshToken: string
  user: UserInfo
}

export interface UserInfo {
  id: string
  username: string
  realName: string
  avatarUrl: string
  permissions: string[]
  roles: string[]
}

export function login(data: LoginRequest): Promise<LoginResult> {
  return request.post('/auth/login', data)
}

export function refreshToken(refreshToken: string): Promise<LoginResult> {
  return request.post('/auth/refresh', { refreshToken })
}

export function getCurrentUser(): Promise<UserInfo> {
  return request.get('/auth/me')
}
```

- [ ] **Step 2: 重写登录页**

`frontend/src/views/login/LoginView.vue`：
```vue
<template>
  <main class="login">
    <el-card class="login__panel" shadow="never">
      <h1>TMS-Go</h1>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleLogin">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-button type="primary" class="login__button" :loading="loading" native-type="submit">
          {{ loading ? '登录中...' : '登录' }}
        </el-button>
      </el-form>
      <p v-if="errorMsg" class="login__error">{{ errorMsg }}</p>
    </el-card>
  </main>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { login } from '@/api/modules/auth'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

const router = useRouter()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const errorMsg = ref('')

const form = reactive({
  username: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  errorMsg.value = ''

  try {
    const result = await login({ username: form.username, password: form.password })
    authStore.setToken(result.accessToken)
    authStore.setUser(result.user)
    permissionStore.setPermissions(result.user.permissions)
    permissionStore.setRoles(result.user.roles)
    router.push('/')
  } catch (err: any) {
    errorMsg.value = err?.message || '登录失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login {
  min-height: 100vh;
  display: grid;
  place-items: center;
  background: #f6f8fb;
}

.login__panel {
  width: min(420px, calc(100vw - 32px));
}

.login__panel h1 {
  margin: 0 0 24px;
  font-size: 24px;
}

.login__button {
  width: 100%;
}

.login__error {
  color: #f56c6c;
  margin: 12px 0 0;
  text-align: center;
}
</style>
```

---

### Task M1-13: 前端 — 权限 Store + 路由守卫

**Files:**
- Modify: `frontend/src/stores/auth.ts`
- Create: `frontend/src/stores/permission.ts`
- Modify: `frontend/src/router/guards.ts`
- Modify: `frontend/src/router/index.ts`

- [ ] **Step 1: 扩展 auth store**

`frontend/src/stores/auth.ts`：
```ts
import { defineStore } from 'pinia'
import type { UserInfo } from '@/api/modules/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('access_token') || '',
    user: null as UserInfo | null,
  }),
  getters: {
    isLoggedIn: (state) => !!state.token,
    username: (state) => state.user?.username || '',
  },
  actions: {
    setToken(token: string) {
      this.token = token
      localStorage.setItem('access_token', token)
    },
    setUser(user: UserInfo) {
      this.user = user
    },
    clearAll() {
      this.token = ''
      this.user = null
      localStorage.removeItem('access_token')
    },
  },
})
```

- [ ] **Step 2: 创建 permission store**

`frontend/src/stores/permission.ts`：
```ts
import { defineStore } from 'pinia'

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    permissions: [] as string[],
    roles: [] as string[],
  }),
  getters: {
    hasPermission: (state) => (code: string) => state.permissions.includes(code),
  },
  actions: {
    setPermissions(permissions: string[]) {
      this.permissions = permissions
    },
    setRoles(roles: string[]) {
      this.roles = roles
    },
    clearAll() {
      this.permissions = []
      this.roles = []
    },
  },
})
```

- [ ] **Step 3: 重写路由守卫**

`frontend/src/router/guards.ts`：
```ts
import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const whiteList = ['/login']

export function setupRouterGuards(router: Router) {
  router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()

    if (whiteList.includes(to.path)) {
      if (authStore.isLoggedIn && to.path === '/login') {
        next('/')
        return
      }
      next()
      return
    }

    if (!authStore.isLoggedIn) {
      next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
      return
    }

    next()
  })
}
```

- [ ] **Step 4: 更新路由**

`frontend/src/router/index.ts`：
```ts
import { createRouter, createWebHistory } from 'vue-router'
import { setupRouterGuards } from './guards'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/login/LoginView.vue'),
      meta: { title: '登录' },
    },
    {
      path: '/',
      component: () => import('@/layouts/BasicLayout.vue'),
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue'),
          meta: { title: '工作台' },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

setupRouterGuards(router)

export default router
```

---

### Task M1-14: 前端 — 布局增加用户信息+退出

**Files:**
- Modify: `frontend/src/layouts/BasicLayout.vue`
- Create: `frontend/src/api/modules/menu.ts`

- [ ] **Step 1: 创建菜单 API**

`frontend/src/api/modules/menu.ts`：
```ts
import { request } from '../request'

export interface MenuItem {
  id: string
  parentId: string | null
  name: string
  path: string
  component: string
  icon: string
  permissionCode: string
  type: string
  sortNo: number
  visible: boolean
  children: MenuItem[]
}

export function getMenus(): Promise<MenuItem[]> {
  return request.get('/system/menus')
}
```

- [ ] **Step 2: 重写布局**

`frontend/src/layouts/BasicLayout.vue`：
```vue
<template>
  <el-container class="layout">
    <el-aside class="layout__aside" width="240px">
      <div class="layout__brand">TMS-Go</div>
      <el-menu router :default-active="activeMenu" class="layout__menu">
        <el-menu-item index="/dashboard">
          <el-icon><Monitor /></el-icon>
          <span>工作台</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="layout__header">
        <span>运输管理系统</span>
        <div class="layout__user">
          <el-dropdown @command="handleCommand">
            <span class="layout__username">
              {{ authStore.username || '管理员' }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="layout__main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Monitor, ArrowDown } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { usePermissionStore } from '@/stores/permission'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()

const activeMenu = computed(() => route.path)

function handleCommand(command: string) {
  if (command === 'logout') {
    authStore.clearAll()
    permissionStore.clearAll()
    router.push('/login')
  }
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
}

.layout__aside {
  border-right: 1px solid #e5e7eb;
  background: #ffffff;
}

.layout__brand {
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  font-size: 18px;
  font-weight: 700;
  border-bottom: 1px solid #e5e7eb;
}

.layout__menu {
  border-right: none;
}

.layout__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
}

.layout__user {
  display: flex;
  align-items: center;
}

.layout__username {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  color: #606266;
}

.layout__main {
  background: #f6f8fb;
}
</style>
```

---

### Task M1-15: 后端 — 系统管理模块（用户CRUD）

**Files:**
- Create: `backend/internal/modules/system/model/user.go`
- Create: `backend/internal/modules/system/repository/user_repository.go`
- Create: `backend/internal/modules/system/repository/user_repository_pg.go`
- Create: `backend/internal/modules/system/service/user_service.go`
- Create: `backend/internal/modules/system/service/user_service_impl.go`
- Create: `backend/internal/modules/system/api/user_handler.go`
- Create: `backend/internal/modules/system/api/user_dto.go`
- Create: `backend/internal/modules/system/api/user_routes.go`

- [ ] **Step 1: System Model**

`backend/internal/modules/system/model/user.go`：
```go
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex:uk_sys_user_username,where:deleted_at is null" json:"username"`
	PasswordHash string         `gorm:"size:256;not null" json:"-"`
	RealName     string         `gorm:"size:64" json:"realName"`
	Phone        string         `gorm:"size:32" json:"phone"`
	Email        string         `gorm:"size:128" json:"email"`
	AvatarURL    string         `gorm:"size:512" json:"avatarUrl"`
	Status       string         `gorm:"size:16;not null;default:active" json:"status"`
	OrgID        *uuid.UUID     `gorm:"type:uuid" json:"orgId"`
	LastLoginAt  *time.Time     `json:"lastLoginAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Roles        []Role         `gorm:"many2many:sys_user_role;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "sys_user"
}

type CreateUserParams struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	RealName  string    `json:"realName"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	RoleIDs   []string  `json:"roleIds"`
}
```

- [ ] **Step 2: System Repository**

`backend/internal/modules/system/repository/user_repository.go`：
```go
package repository

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	List(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User, roleIDs []string) error
	Update(ctx context.Context, user *model.User, roleIDs []string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
```

`backend/internal/modules/system/repository/user_repository_pg.go`：
```go
package repository

import (
	"context"
	"errors"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{}).Preload("Roles")
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Preload("Roles").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return r.syncRoles(tx, user.ID.String(), roleIDs)
	})
}

func (r *userRepository) Update(ctx context.Context, user *model.User, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(user).Select("real_name", "phone", "email", "status", "updated_at").
			Updates(user).Error; err != nil {
			return err
		}
		return r.syncRoles(tx, user.ID.String(), roleIDs)
	})
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// 软删除
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) syncRoles(tx *gorm.DB, userID string, roleIDs []string) error {
	if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}
	for _, roleID := range roleIDs {
		ur := &model.UserRole{
			UserID: uuid.MustParse(userID),
			RoleID: uuid.MustParse(roleID),
		}
		if err := tx.Create(ur).Error; err != nil {
			return err
		}
	}
	return nil
}

// UserRole 需要定义为 model
```

- [ ] **Step 3: 补充 model 中的 UserRole**

在 `backend/internal/modules/system/model/user.go` 末尾添加：
```go
type UserRole struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null" json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserRole) TableName() string {
	return "sys_user_role"
}
```

- [ ] **Step 4: System User Service**

`backend/internal/modules/system/service/user_service.go`：
```go
package service

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
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
```

`backend/internal/modules/system/service/user_service_impl.go`：
```go
package service

import (
	"context"

	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/system/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/crypto"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/pagination"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	roleIDs := params.RoleIDs
	if err := s.repo.Update(ctx, user, roleIDs); err != nil {
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
```

---

### Task M1-16: 编译验证 & 初次启动测试

- [ ] **Step 1: 编译所有代码**

```bash
cd backend && go build ./...
```

- [ ] **Step 2: 启动 Docker 中间件**

```bash
cd docker && docker compose up -d postgres redis rabbitmq minio
```

- [ ] **Step 3: 执行数据库迁移**

```bash
docker compose -f docker/docker-compose.yml exec postgres psql -U tms -d tms -f /docker-entrypoint-initdb.d/000001_init.up.sql
```

- [ ] **Step 4: 启动后端服务**

```bash
cd backend && go run ./cmd/server
```

- [ ] **Step 5: 测试健康检查**

```bash
curl http://localhost:8080/health
# 预期输出: {"code":0,"message":"success","data":{"status":"ok","service":"tms-api"}}
```

- [ ] **Step 6: 测试登录（需要先创建种子用户）**

```bash
# 通过种子数据插入管理员用户后再测试
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

---

### Task M1-17: 创建种子数据

**Files:**
- Create: `backend/migrations/000008_seed_admin.up.sql`
- Create: `backend/migrations/000008_seed_admin.down.sql`

- [ ] **Step 1: 种子数据迁移**

`backend/migrations/000008_seed_admin.up.sql`：
```sql
-- 管理员角色
insert into sys_role (id, code, name, remark, status, sort_no)
values (uuid_generate_v4(), 'admin', '超级管理员', '系统超级管理员角色', 'active', 1);

-- 管理员用户（密码: admin123，bcrypt hash）
-- 需要通过 Go 程序生成 hash 或使用预计算的 hash
-- bcrypt cost=12, password=admin123 的 hash:
insert into sys_user (id, username, password_hash, real_name, status)
values (
  uuid_generate_v4(),
  'admin',
  '$2a$12$LJ3m4ys3GZfnYMz8k7lJJeQWDmY/8Kj.SWkmHlKHR1JfGMEqX8Fy',
  '系统管理员',
  'active'
);

-- 关联管理员用户到管理员角色
insert into sys_user_role (id, user_id, role_id)
select uuid_generate_v4(), u.id, r.id
from sys_user u, sys_role r
where u.username = 'admin' and r.code = 'admin';

-- 基础权限
insert into sys_permission (id, code, name, type, sort_no, status)
values
  (uuid_generate_v4(), 'dashboard', '工作台', 'menu', 1, 'active'),
  (uuid_generate_v4(), 'system:user:list', '用户管理', 'menu', 2, 'active'),
  (uuid_generate_v4(), 'system:user:create', '创建用户', 'button', 1, 'active'),
  (uuid_generate_v4(), 'system:user:update', '编辑用户', 'button', 2, 'active'),
  (uuid_generate_v4(), 'system:user:delete', '删除用户', 'button', 3, 'active'),
  (uuid_generate_v4(), 'system:role:list', '角色管理', 'menu', 3, 'active'),
  (uuid_generate_v4(), 'system:role:create', '创建角色', 'button', 1, 'active'),
  (uuid_generate_v4(), 'system:menu:list', '菜单管理', 'menu', 4, 'active');

-- 关联管理员角色到所有权限
insert into sys_role_permission (id, role_id, permission_id)
select uuid_generate_v4(), r.id, p.id
from sys_role r, sys_permission p
where r.code = 'admin';

-- 基础菜单
insert into sys_menu (id, parent_id, name, path, component, icon, permission_code, type, sort_no, visible, status)
values
  (uuid_generate_v4(), null, '工作台', '/dashboard', 'views/dashboard/DashboardView', 'Monitor', 'dashboard', 'menu', 1, true, 'active'),
  (uuid_generate_v4(), null, '系统管理', '/system', null, 'Setting', null, 'directory', 9, true, 'active'),
  (uuid_generate_v4(), (select id from sys_menu where name = '系统管理' limit 1), '用户管理', '/system/user', 'views/system/UserList', null, 'system:user:list', 'menu', 1, true, 'active');
```

`backend/migrations/000008_seed_admin.down.sql`：
```sql
delete from sys_role_permission;
delete from sys_user_role;
delete from sys_permission;
delete from sys_menu;
delete from sys_role;
delete from sys_user;
```

---

## 验证清单

M1 阶段完成后，逐一验证：

- [ ] `POST /api/v1/auth/login` 使用 admin/admin123 成功返回 token
- [ ] 使用 token 访问 `GET /api/v1/auth/me` 返回用户信息
- [ ] 无 token 访问受保护接口返回 401
- [ ] `GET /api/v1/users?page=1&pageSize=20` 返回用户列表
- [ ] `POST /api/v1/users` 创建新用户成功
- [ ] 新用户可以登录
- [ ] 前端登录页可正常登录并跳转工作台
- [ ] 未登录访问前端其他页面自动跳转登录页
- [ ] 顶部栏显示用户名，可退出登录
