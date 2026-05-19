package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TaskStatusPending   = "pending"
	TaskStatusDeparted  = "departed"
	TaskStatusInTransit = "in_transit"
	TaskStatusArrived   = "arrived"
	TaskStatusSigned    = "signed"
	TaskStatusClosed    = "closed"
)

type TransportTask struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TaskNo          string         `gorm:"size:32;not null" json:"taskNo"`
	OrderID         uuid.UUID      `gorm:"type:uuid;not null" json:"orderId"`
	DispatchPlanID  *uuid.UUID     `gorm:"type:uuid" json:"dispatchPlanId"`
	CarrierID       *uuid.UUID     `gorm:"type:uuid" json:"carrierId"`
	VehicleID       *uuid.UUID     `gorm:"type:uuid" json:"vehicleId"`
	DriverID        *uuid.UUID     `gorm:"type:uuid" json:"driverId"`
	DriverName      string         `gorm:"size:64" json:"driverName"`
	PlateNo         string         `gorm:"size:32" json:"plateNo"`
	OriginName      string         `gorm:"size:128" json:"originName"`
	DestName        string         `gorm:"size:128" json:"destName"`
	PlanDepartTime  *time.Time     `json:"planDepartTime"`
	PlanArriveTime  *time.Time     `json:"planArriveTime"`
	ActualDepartTime *time.Time    `json:"actualDepartTime"`
	ActualArriveTime *time.Time    `json:"actualArriveTime"`
	Status          string         `gorm:"size:32;not null;default:pending" json:"status"`
	Remark          string         `gorm:"size:512" json:"remark"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TransportTask) TableName() string { return "tms_transport_task" }
