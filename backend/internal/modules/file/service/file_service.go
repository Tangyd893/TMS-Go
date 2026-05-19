package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/file/model"
	"github.com/Tangyd893/TMS-Go/backend/internal/modules/file/repository"
	"github.com/Tangyd893/TMS-Go/backend/internal/shared/errors"
	"github.com/google/uuid"
)

var (
	ErrFileNotFound = errors.New(404001, "文件不存在")
	ErrFileTooLarge = errors.New(400001, "文件大小超过限制")
	errFileTypeNotAllowed = errors.New(400001, "不支持的文件类型")

	maxFileSize   = int64(10 * 1024 * 1024) // 10MB
	allowedExts   = map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".csv": true, ".txt": true, ".zip": true, ".rar": true,
	}
)

type FileService interface {
	Upload(ctx context.Context, reader io.Reader, fileName string, bizType string, bizID string, userID string) (*model.FileObject, error)
	GetDownloadURL(ctx context.Context, id string) (io.ReadCloser, *model.FileObject, error)
	ListByBiz(ctx context.Context, bizType, bizID string) ([]model.FileObject, error)
	Delete(ctx context.Context, id string) error
}

type fileService struct {
	repo      repository.FileRepository
	bucket    string
}

func NewFileService(repo repository.FileRepository, cfg config.MinIOConfig) FileService {
	bucket := cfg.Bucket
	if bucket == "" { bucket = "tms-dev" }
	return &fileService{repo: repo, bucket: bucket}
}

func (s *fileService) Upload(ctx context.Context, reader io.Reader, fileName string, bizType string, bizIDStr string, userID string) (*model.FileObject, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if !allowedExts[ext] {
		return nil, errFileTypeNotAllowed
	}

	now := time.Now()
	uid := uuid.New()
	objectKey := fmt.Sprintf("%s/%s/%s_%s%s",
		bizType,
		now.Format("2006/01/02"),
		uid.String(),
		uuid.New().String()[:8],
		ext,
	)

	var bizIDVal *uuid.UUID
	if bizIDStr != "" {
		bizID := uuid.MustParse(bizIDStr)
		bizIDVal = &bizID
	}
	var uploadByVal *uuid.UUID
	if userID != "" {
		uid := uuid.MustParse(userID)
		uploadByVal = &uid
	}

	fileObj := &model.FileObject{
		ID:           uid,
		Bucket:       s.bucket,
		ObjectKey:    objectKey,
		OriginalName: fileName,
		BizType:      bizType,
		BizID:        bizIDVal,
		UploadedBy:   uploadByVal,
		CreatedAt:    now,
	}

	if err := s.repo.Create(ctx, fileObj); err != nil {
		return nil, err
	}
	return fileObj, nil
}

func (s *fileService) GetDownloadURL(ctx context.Context, id string) (io.ReadCloser, *model.FileObject, error) {
	uid := uuid.MustParse(id)
	f, err := s.repo.FindByID(ctx, uid)
	if err != nil { return nil, nil, ErrFileNotFound }
	return nil, f, nil
}

func (s *fileService) ListByBiz(ctx context.Context, bizType, bizIDStr string) ([]model.FileObject, error) {
	bizID := uuid.MustParse(bizIDStr)
	return s.repo.ListByBiz(ctx, bizType, bizID)
}

func (s *fileService) Delete(ctx context.Context, id string) error {
	uid := uuid.MustParse(id)
	return s.repo.Delete(ctx, uid)
}
