package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderNo             string         `gorm:"size:32;not null" json:"orderNo"`
	CustomerID          uuid.UUID      `gorm:"type:uuid;not null" json:"customerId"`
	CustomerName        string         `gorm:"size:128" json:"customerName"`
	ShipperName         string         `gorm:"size:128" json:"shipperName"`
	ShipperPhone        string         `gorm:"size:32" json:"shipperPhone"`
	ShipperAddress      string         `gorm:"size:256" json:"shipperAddress"`
	ReceiverName        string         `gorm:"size:128" json:"receiverName"`
	ReceiverPhone       string         `gorm:"size:32" json:"receiverPhone"`
	ReceiverAddress     string         `gorm:"size:256" json:"receiverAddress"`
	OriginStationID     *uuid.UUID     `gorm:"type:uuid" json:"originStationId"`
	OriginName          string         `gorm:"size:128" json:"originName"`
	DestStationID       *uuid.UUID     `gorm:"type:uuid" json:"destStationId"`
	DestName            string         `gorm:"size:128" json:"destName"`
	PlanPickupTime      *time.Time     `json:"planPickupTime"`
	PlanDeliveryTime    *time.Time     `json:"planDeliveryTime"`
	ActualPickupTime    *time.Time     `json:"actualPickupTime"`
	ActualDeliveryTime  *time.Time     `json:"actualDeliveryTime"`
	CargoName           string         `gorm:"size:128" json:"cargoName"`
	CargoWeight         *float64       `gorm:"type:numeric(12,2)" json:"cargoWeight"`
	CargoVolume         *float64       `gorm:"type:numeric(12,3)" json:"cargoVolume"`
	CargoQuantity       int            `json:"cargoQuantity"`
	TransportRequirement string        `gorm:"size:512" json:"transportRequirement"`
	Status              string         `gorm:"size:32;not null;default:draft" json:"status"`
	Remark              string         `gorm:"size:512" json:"remark"`
	CreatedBy           *uuid.UUID     `gorm:"type:uuid" json:"createdBy"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
	CargoItems          []OrderCargo   `gorm:"foreignKey:OrderID" json:"cargoItems,omitempty"`
}

func (Order) TableName() string { return "tms_order" }

type OrderCargo struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"orderId"`
	CargoName string    `gorm:"size:128" json:"cargoName"`
	CargoCode string    `gorm:"size:64" json:"cargoCode"`
	CargoType string    `gorm:"size:32" json:"cargoType"`
	Quantity  int       `gorm:"default:1" json:"quantity"`
	Weight    *float64  `gorm:"type:numeric(12,2)" json:"weight"`
	Volume    *float64  `gorm:"type:numeric(12,3)" json:"volume"`
	Unit      string    `gorm:"size:16" json:"unit"`
	Remark    string    `gorm:"size:256" json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
}

func (OrderCargo) TableName() string { return "tms_order_cargo" }

type StatusHistory struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	BizType      string    `gorm:"size:32;not null" json:"bizType"`
	BizID        uuid.UUID `gorm:"type:uuid;not null" json:"bizId"`
	FromStatus   string    `gorm:"size:32" json:"fromStatus"`
	ToStatus     string    `gorm:"size:32;not null" json:"toStatus"`
	Action       string    `gorm:"size:32" json:"action"`
	OperatorID   *uuid.UUID `gorm:"type:uuid" json:"operatorId"`
	OperatorName string    `gorm:"size:64" json:"operatorName"`
	Remark       string    `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (StatusHistory) TableName() string { return "tms_status_history" }

type CreateOrderRequest struct {
	CustomerID          uuid.UUID              `json:"customerId"`
	ShipperName         string                 `json:"shipperName"`
	ShipperPhone        string                 `json:"shipperPhone"`
	ShipperAddress      string                 `json:"shipperAddress"`
	ReceiverName        string                 `json:"receiverName"`
	ReceiverPhone       string                 `json:"receiverPhone"`
	ReceiverAddress     string                 `json:"receiverAddress"`
	OriginStationID     *string                `json:"originStationId"`
	OriginName          string                 `json:"originName"`
	DestStationID       *string                `json:"destStationId"`
	DestName            string                 `json:"destName"`
	PlanPickupTime      *time.Time             `json:"planPickupTime"`
	PlanDeliveryTime    *time.Time             `json:"planDeliveryTime"`
	TransportRequirement string                `json:"transportRequirement"`
	Remark              string                 `json:"remark"`
	CargoItems          []CreateCargoRequest   `json:"cargoItems"`
}

type CreateCargoRequest struct {
	CargoName string  `json:"cargoName"`
	CargoType string  `json:"cargoType"`
	Quantity  int     `json:"quantity"`
	Weight    *float64 `json:"weight"`
	Volume    *float64 `json:"volume"`
	Unit      string  `json:"unit"`
}

type UpdateOrderRequest struct {
	ShipperName         *string `json:"shipperName"`
	ShipperPhone        *string `json:"shipperPhone"`
	ShipperAddress      *string `json:"shipperAddress"`
	ReceiverName        *string `json:"receiverName"`
	ReceiverPhone       *string `json:"receiverPhone"`
	ReceiverAddress     *string `json:"receiverAddress"`
	OriginStationID     *string `json:"originStationId"`
	OriginName          *string `json:"originName"`
	DestStationID       *string `json:"destStationId"`
	DestName            *string `json:"destName"`
	PlanPickupTime      *time.Time `json:"planPickupTime"`
	PlanDeliveryTime    *time.Time `json:"planDeliveryTime"`
	TransportRequirement *string `json:"transportRequirement"`
	Remark              *string `json:"remark"`
	CargoItems          []CreateCargoRequest `json:"cargoItems"`
}
