package model

import (
	"time"

	"github.com/google/uuid"
)

// FileObject 文件对象元数据
type FileObject struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Bucket       string     `gorm:"size:64;not null" json:"bucket"`
	ObjectKey    string     `gorm:"size:512;not null" json:"objectKey"`
	OriginalName string     `gorm:"size:256;not null" json:"originalName"`
	ContentType  string     `gorm:"size:128" json:"contentType"`
	Size         int64      `gorm:"not null;default:0" json:"size"`
	SHA256       string     `gorm:"size:64" json:"sha256"`
	BizType      string     `gorm:"size:32" json:"bizType"`
	BizID        *uuid.UUID `gorm:"type:uuid" json:"bizId"`
	UploadedBy   *uuid.UUID `gorm:"type:uuid" json:"uploadedBy"`
	CreatedAt    time.Time  `json:"createdAt"`
}

func (FileObject) TableName() string { return "file_object" }
