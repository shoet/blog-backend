package upload_file

import (
	"context"

	"github.com/shoet/blog/internal/infrastructure/models"
)

type FileRepository interface {
	GetUploadURL(ctx context.Context, file *models.File) (uploadURL string, destinationURL string, err error)
}

type Usecase struct {
	FileRepository FileRepository
}

func NewUsecase(
	fileRepository FileRepository,
) *Usecase {
	return &Usecase{
		FileRepository: fileRepository,
	}
}

type UploadFileInput struct {
	FileType string
	FileName string
}

func (u *Usecase) Run(ctx context.Context, input UploadFileInput) (string, string, error) {
	file, err := models.NewFile(models.FileType(input.FileType), input.FileName)
	if err != nil {
		return "", "", err
	}
	uploadURL, destinationURL, err := u.FileRepository.GetUploadURL(ctx, file)
	if err != nil {
		return "", "", err
	}
	return uploadURL, destinationURL, nil
}
