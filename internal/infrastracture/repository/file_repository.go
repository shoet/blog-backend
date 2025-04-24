package repository

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastracture/adapter"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type FileRepository struct {
	Config    *config.Config
	S3Adapter adapter.S3Adapter
}

func (r *FileRepository) ExistsFile(ctx context.Context, file *models.File) (bool, error) {
	bucketName, err := file.GetBucketName(r.Config)
	if err != nil {
		return false, fmt.Errorf("failed to get bucket name")
	}
	key, err := file.GetBucketKey(r.Config)
	if err != nil {
		return false, fmt.Errorf("failed to get file key")
	}
	fileKey := fmt.Sprintf("%s/%s", key, file.FileName)
	exists, err := r.S3Adapter.ExistObject(ctx, bucketName, fileKey)
	if err != nil {
		return false, fmt.Errorf("failed to check exist object")
	}
	return exists, nil
}
