package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ExceptionStatusPending    = "pending"
	ExceptionStatusProcessing = "processing"
	ExceptionStatusResolved   = "resolved"
	ExceptionStatusClosed     = "closed"
)

const (
	ExceptionSeverityNormal   = "normal"
	ExceptionSeveritySerious  = "serious"
	ExceptionSeverityCritical = "critical"
)

// Exception 运输异常
type Exception struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TaskID        uuid.UUID      `gorm:"type:uuid;not null" json:"taskId"`
	OrderID       uuid.UUID      `gorm:"type:uuid;not null" json:"orderId"`
	ExceptionNo   string         `gorm:"size:32;not null" json:"exceptionNo"`
	ExceptionType string         `gorm:"size:64;not null" json:"exceptionType"`
	Description   string         `gorm:"type:text;not null" json:"description"`
	Severity      string         `gorm:"size:16;not null;default:normal" json:"severity"`
	ReportBy      uuid.UUID      `gorm:"type:uuid;not null" json:"reportBy"`
	ReportByName  string         `gorm:"size:64" json:"reportByName"`
	HandlerID     *uuid.UUID     `gorm:"type:uuid" json:"handlerId"`
	HandlerName   string         `gorm:"size:64" json:"handlerName"`
	HandleResult  string         `gorm:"type:text" json:"handleResult"`
	HandleAt      *time.Time     `json:"handleAt"`
	Status        string         `gorm:"size:32;not null;default:pending" json:"status"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Exception) TableName() string { return "tms_exception" }

// ExceptionLog 异常处理日志
type ExceptionLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ExceptionID uuid.UUID `gorm:"type:uuid;not null" json:"exceptionId"`
	Action      string    `gorm:"size:32;not null" json:"action"`
	Content     string    `gorm:"type:text" json:"content"`
	OperatorID  *uuid.UUID `gorm:"type:uuid" json:"operatorId"`
	OperatorName string   `gorm:"size:64" json:"operatorName"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (ExceptionLog) TableName() string { return "tms_exception_log" }

// CreateExceptionRequest 创建异常请求
type CreateExceptionRequest struct {
	TaskID        uuid.UUID `json:"taskId"`
	OrderID       uuid.UUID `json:"orderId"`
	ExceptionType string    `json:"exceptionType"`
	Description   string    `json:"description"`
	Severity      string    `json:"severity"`
}

// HandleExceptionRequest 处理异常请求
type HandleExceptionRequest struct {
	HandleResult string `json:"handleResult"`
}

// CloseExceptionRequest 关闭异常请求
type CloseExceptionRequest struct {
	Remark string `json:"remark"`
}
