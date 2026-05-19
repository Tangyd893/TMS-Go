package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	SettlementStatusPending   = "pending"
	SettlementStatusCompleted = "completed"
)

type Settlement struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	SettlementNo     string         `gorm:"size:32;not null" json:"settlementNo"`
	StatementID      uuid.UUID      `gorm:"type:uuid;not null" json:"statementId"`
	PartnerID        uuid.UUID      `gorm:"type:uuid;not null" json:"partnerId"`
	PartnerName      string         `gorm:"size:128" json:"partnerName"`
	PartnerType      string         `gorm:"size:16;not null" json:"partnerType"`
	SettlementAmount float64        `gorm:"type:numeric(18,2);not null;default:0" json:"settlementAmount"`
	SettlementMethod string         `gorm:"size:32" json:"settlementMethod"`
	Status           string         `gorm:"size:32;not null;default:pending" json:"status"`
	SettledAt        *time.Time     `json:"settledAt"`
	Remark           string         `gorm:"size:512" json:"remark"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Settlement) TableName() string { return "fin_settlement" }
