package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username     string         `gorm:"size:64;not null" json:"username"`
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
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code        string         `gorm:"size:64;not null" json:"code"`
	Name        string         `gorm:"size:128;not null" json:"name"`
	Remark      string         `gorm:"size:512" json:"remark"`
	Status      string         `gorm:"size:16;not null;default:active" json:"status"`
	SortNo      int            `gorm:"not null;default:0" json:"sortNo"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Permissions []Permission   `gorm:"many2many:sys_role_permission;" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "sys_role"
}

type Permission struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code      string         `gorm:"size:128;not null" json:"code"`
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

type UserRole struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	RoleID    uuid.UUID `gorm:"type:uuid;not null" json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserRole) TableName() string {
	return "sys_user_role"
}

type CreateUserParams struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	RealName string   `json:"realName"`
	Phone    string   `json:"phone"`
	Email    string   `json:"email"`
	Status   string   `json:"status"`
	RoleIDs  []string `json:"roleIds"`
}
