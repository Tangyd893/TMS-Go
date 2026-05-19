package model

import (
	"time"

	"github.com/google/uuid"
)

// Receipt 签收回单
type Receipt struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TaskID      uuid.UUID `gorm:"type:uuid;not null" json:"taskId"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null" json:"orderId"`
	ReceiptNo   string    `gorm:"size:32;not null" json:"receiptNo"`
	SignBy      string    `gorm:"size:64" json:"signBy"`
	SignAt      *time.Time `json:"signAt"`
	SignImageURL string   `gorm:"size:512" json:"signImageUrl"`
	Remark      string    `gorm:"size:512" json:"remark"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid" json:"createdBy"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Receipt) TableName() string { return "tms_receipt" }

// CreateReceiptRequest 创建回单请求
type CreateReceiptRequest struct {
	SignBy       string `json:"signBy"`
	SignImageURL string `json:"signImageUrl"`
	Remark       string `json:"remark"`
}
