package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DispatchStatusPending  = "pending"
	DispatchStatusAssigned = "assigned"
	DispatchStatusDeparted = "departed"
	DispatchStatusCompleted = "completed"
	DispatchStatusCancelled = "cancelled"
)

type DispatchPlan struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PlanNo    string         `gorm:"size:32;not null" json:"planNo"`
	Status    string         `gorm:"size:32;not null;default:pending" json:"status"`
	Remark    string         `gorm:"size:512" json:"remark"`
	CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"createdBy"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Details   []DispatchDetail `gorm:"foreignKey:PlanID" json:"details,omitempty"`
}

func (DispatchPlan) TableName() string { return "tms_dispatch_plan" }

type DispatchDetail struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null" json:"planId"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"orderId"`
	CarrierID *uuid.UUID `gorm:"type:uuid" json:"carrierId"`
	VehicleID *uuid.UUID `gorm:"type:uuid" json:"vehicleId"`
	DriverID  *uuid.UUID `gorm:"type:uuid" json:"driverId"`
	Seq       int        `gorm:"default:1" json:"seq"`
	CreatedAt time.Time  `json:"createdAt"`
}

func (DispatchDetail) TableName() string { return "tms_dispatch_detail" }
