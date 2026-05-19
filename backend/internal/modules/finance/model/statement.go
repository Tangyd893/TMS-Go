package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatementStatusDraft     = "draft"
	StatementStatusConfirmed = "confirmed"
	StatementStatusSettled   = "settled"

	PartnerTypeCustomer = "customer"
	PartnerTypeCarrier  = "carrier"
)

type Statement struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	StatementNo         string         `gorm:"size:32;not null" json:"statementNo"`
	PartnerID           uuid.UUID      `gorm:"type:uuid;not null" json:"partnerId"`
	PartnerName         string         `gorm:"size:128" json:"partnerName"`
	PartnerType         string         `gorm:"size:16;not null" json:"partnerType"`
	StatementPeriodStart string        `gorm:"type:date;not null" json:"statementPeriodStart"`
	StatementPeriodEnd  string         `gorm:"type:date;not null" json:"statementPeriodEnd"`
	TotalReceivable     float64        `gorm:"type:numeric(18,2);default:0" json:"totalReceivable"`
	TotalPayable        float64        `gorm:"type:numeric(18,2);default:0" json:"totalPayable"`
	Status              string         `gorm:"size:32;not null;default:draft" json:"status"`
	ConfirmedAt         *time.Time     `json:"confirmedAt"`
	ConfirmedBy         *uuid.UUID     `gorm:"type:uuid" json:"confirmedBy"`
	Remark              string         `gorm:"size:512" json:"remark"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Statement) TableName() string { return "fin_statement" }
