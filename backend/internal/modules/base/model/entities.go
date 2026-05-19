package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code         string         `gorm:"size:64;not null" json:"code"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	ShortName    string         `gorm:"size:64" json:"shortName"`
	ContactName  string         `gorm:"size:64" json:"contactName"`
	ContactPhone string         `gorm:"size:32" json:"contactPhone"`
	ContactEmail string         `gorm:"size:128" json:"contactEmail"`
	Address      string         `gorm:"size:256" json:"address"`
	Status       string         `gorm:"size:16;not null;default:active" json:"status"`
	Remark       string         `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Customer) TableName() string { return "base_customer" }

type Carrier struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code            string         `gorm:"size:64;not null" json:"code"`
	Name            string         `gorm:"size:128;not null" json:"name"`
	ShortName       string         `gorm:"size:64" json:"shortName"`
	ContactName     string         `gorm:"size:64" json:"contactName"`
	ContactPhone    string         `gorm:"size:32" json:"contactPhone"`
	ContactEmail    string         `gorm:"size:128" json:"contactEmail"`
	Address         string         `gorm:"size:256" json:"address"`
	PaymentMethod   string         `gorm:"size:32" json:"paymentMethod"`
	SettlementCycle string         `gorm:"size:32" json:"settlementCycle"`
	Status          string         `gorm:"size:16;not null;default:active" json:"status"`
	Remark          string         `gorm:"size:512" json:"remark"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Carrier) TableName() string { return "base_carrier" }

type Vehicle struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PlateNo    string         `gorm:"size:32;not null" json:"plateNo"`
	Type       string         `gorm:"column:vehicle_type;size:32" json:"type"`
	BrandModel string         `gorm:"size:64" json:"brandModel"`
	Color      string         `gorm:"size:16" json:"color"`
	MaxLoad    *float64       `gorm:"type:numeric(12,2)" json:"maxLoad"`
	MaxVolume  *float64       `gorm:"type:numeric(12,3)" json:"maxVolume"`
	OwnerType  string         `gorm:"size:16" json:"ownerType"`
	CarrierID  *uuid.UUID     `gorm:"type:uuid" json:"carrierId"`
	DriverID   *uuid.UUID     `gorm:"type:uuid" json:"driverId"`
	Status     string         `gorm:"size:16;not null;default:active" json:"status"`
	Remark     string         `gorm:"size:512" json:"remark"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Vehicle) TableName() string { return "base_vehicle" }

type Driver struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code             string         `gorm:"size:64;not null" json:"code"`
	Name             string         `gorm:"size:64;not null" json:"name"`
	IDCard           string         `gorm:"column:id_card;size:32" json:"idCard"`
	Phone            string         `gorm:"size:32" json:"phone"`
	LicenseType      string         `gorm:"size:16" json:"licenseType"`
	LicenseNo        string         `gorm:"size:32" json:"licenseNo"`
	LicenseExpireDate *time.Time    `gorm:"type:date" json:"licenseExpireDate"`
	Status           string         `gorm:"size:16;not null;default:active" json:"status"`
	Remark           string         `gorm:"size:512" json:"remark"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Driver) TableName() string { return "base_driver" }

type Route struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code           string         `gorm:"size:64;not null" json:"code"`
	Name           string         `gorm:"size:128;not null" json:"name"`
	OriginStationID *uuid.UUID    `gorm:"type:uuid" json:"originStationId"`
	OriginName     string         `gorm:"size:128" json:"originName"`
	DestStationID  *uuid.UUID     `gorm:"type:uuid" json:"destStationId"`
	DestName       string         `gorm:"size:128" json:"destName"`
	DistanceKM     *float64       `gorm:"type:numeric(10,2)" json:"distanceKm"`
	EstimatedHours *float64       `gorm:"type:numeric(6,1)" json:"estimatedHours"`
	Status         string         `gorm:"size:16;not null;default:active" json:"status"`
	Remark         string         `gorm:"size:512" json:"remark"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Route) TableName() string { return "base_route" }

type Station struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Code         string         `gorm:"size:64;not null" json:"code"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	Province     string         `gorm:"size:64" json:"province"`
	City         string         `gorm:"size:64" json:"city"`
	District     string         `gorm:"size:64" json:"district"`
	Address      string         `gorm:"size:256" json:"address"`
	ContactName  string         `gorm:"size:64" json:"contactName"`
	ContactPhone string         `gorm:"size:32" json:"contactPhone"`
	Latitude     *float64       `gorm:"type:numeric(10,7)" json:"latitude"`
	Longitude    *float64       `gorm:"type:numeric(10,7)" json:"longitude"`
	Status       string         `gorm:"size:16;not null;default:active" json:"status"`
	Remark       string         `gorm:"size:512" json:"remark"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Station) TableName() string { return "base_station" }
