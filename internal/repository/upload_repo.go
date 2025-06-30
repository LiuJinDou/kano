package repository

import (
	"context"
	"kano/internal/config"
	"kano/internal/repository/model"
)

type UploadRepo struct {
}

func NewUploadRepo() *UploadRepo {
	return &UploadRepo{}
}
func (u *UploadRepo) SaveUploadRecord(ctx context.Context, uploadRecord []*model.UploadRecord) error {
	db := config.GetDB()
	return db.WithContext(ctx).Create(uploadRecord).Error
}
