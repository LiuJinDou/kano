package service

import (
	"context"
	"kano/internal/provider/tencent"
	"kano/internal/repository/model"
	"mime/multipart"
)

// Uploader defines the interface for file upload services.
const UploadTypeTencent = "tencent"

// Uploader defines the interface for file upload services.
type Uploader interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
	GetCredential(ctx context.Context, application_code string) (interface{}, error)
	UploadRecord(ctx context.Context, uploadRecord []*model.UploadRecord) error
}

// NewUploader creates a new uploader based on the specified upload type.
// Currently supports "tencent" for Tencent Cloud COS.
func NewUploader(uploadType string) Uploader {
	switch uploadType {
	case UploadTypeTencent:
		return tencent.New()

	default:
		return nil
	}
}
