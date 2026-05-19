package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	NodeTypePickup        = "pickup"
	NodeTypeDepart        = "depart"
	NodeTypeTransitArrive = "transit_arrive"
	NodeTypeTransitDepart = "transit_depart"
	NodeTypeArrive        = "arrive"
	NodeTypeSign          = "sign"
)

// TransportNode 运输节点记录
type TransportNode struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TaskID       uuid.UUID  `gorm:"type:uuid;not null" json:"taskId"`
	NodeType     string     `gorm:"size:32;not null" json:"nodeType"`
	NodeName     string     `gorm:"size:128" json:"nodeName"`
	LocationName string     `gorm:"size:256" json:"locationName"`
	Longitude    *float64   `gorm:"type:numeric(10,6)" json:"longitude"`
	Latitude     *float64   `gorm:"type:numeric(10,6)" json:"latitude"`
	ArrivedAt    *time.Time `json:"arrivedAt"`
	DepartedAt   *time.Time `json:"departedAt"`
	Remark       string     `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time  `json:"createdAt"`
	CreatedBy    *uuid.UUID `gorm:"type:uuid" json:"createdBy"`
}

func (TransportNode) TableName() string { return "tms_transport_node" }

// CreateNodeRequest 创建节点请求
type CreateNodeRequest struct {
	NodeType     string     `json:"nodeType"`
	NodeName     string     `json:"nodeName"`
	LocationName string     `json:"locationName"`
	Longitude    *float64   `json:"longitude"`
	Latitude     *float64   `json:"latitude"`
	ArrivedAt    *time.Time `json:"arrivedAt"`
	DepartedAt   *time.Time `json:"departedAt"`
	Remark       string     `json:"remark"`
}
