package dict

import (
	"time"

	"github.com/google/uuid"
)

type DictType struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code      string    `gorm:"size:64;not null" json:"code"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Status    string    `gorm:"size:16;not null;default:active" json:"status"`
	Remark    string    `gorm:"size:512" json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
	Items     []DictItem `gorm:"foreignKey:TypeCode;references:Code" json:"items,omitempty"`
}

func (DictType) TableName() string { return "sys_dict_type" }

type DictItem struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TypeCode string    `gorm:"size:64;not null" json:"typeCode"`
	ItemCode string    `gorm:"size:64;not null" json:"itemCode"`
	ItemName string    `gorm:"size:128;not null" json:"itemName"`
	SortNo   int       `gorm:"not null;default:0" json:"sortNo"`
	Status   string    `gorm:"size:16;not null;default:active" json:"status"`
	Remark   string    `gorm:"size:512" json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (DictItem) TableName() string { return "sys_dict_item" }
