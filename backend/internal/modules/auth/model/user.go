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
