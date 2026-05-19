package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payable struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PayableNo    string         `gorm:"size:32;not null" json:"payableNo"`
	OrderID      uuid.UUID      `gorm:"type:uuid;not null" json:"orderId"`
	TaskID       *uuid.UUID     `gorm:"type:uuid" json:"taskId"`
	CarrierID    uuid.UUID      `gorm:"type:uuid;not null" json:"carrierId"`
	CarrierName  string         `gorm:"size:128" json:"carrierName"`
	FeeItemID    uuid.UUID      `gorm:"type:uuid;not null" json:"feeItemId"`
	FeeItemName  string         `gorm:"size:128" json:"feeItemName"`
	Amount       float64        `gorm:"type:numeric(18,2);not null;default:0" json:"amount"`
	HasTax       bool           `gorm:"not null;default:false" json:"hasTax"`
	TaxRate      float64        `gorm:"type:numeric(5,4);default:0" json:"taxRate"`
	TaxAmount    float64        `gorm:"type:numeric(18,2);default:0" json:"taxAmount"`
	TotalAmount  float64        `gorm:"type:numeric(18,2);not null;default:0" json:"totalAmount"`
	Status       string         `gorm:"size:32;not null;default:pending" json:"status"`
	Remark       string         `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Payable) TableName() string { return "fin_payable" }
