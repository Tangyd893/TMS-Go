package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	FeeItemStatusActive   = "active"
	FeeItemStatusDisabled = "disabled"
)

type FeeItem struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code       string         `gorm:"size:64;not null" json:"code"`
	Name       string         `gorm:"size:128;not null" json:"name"`
	FeeType    string         `gorm:"size:16;not null;default:transport" json:"feeType"`
	UnitPrice  float64        `gorm:"type:numeric(18,2);not null;default:0" json:"unitPrice"`
	ChargeUnit string         `gorm:"size:32;not null;default:per_order" json:"chargeUnit"`
	Status     string         `gorm:"size:16;not null;default:active" json:"status"`
	Remark     string         `gorm:"size:512" json:"remark"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (FeeItem) TableName() string { return "fin_fee_item" }
