package storage_presigned_thumbnail

import "context"

type ContentsService interface {
	GenerateThumbnailPutURL(fileName string) (presignedUrl, objectUrl string, err error)
}

type Usecase struct {
	contentsService ContentsService
}

func NewUsecase(contentsService ContentsService) *Usecase {
	return &Usecase{
		contentsService: contentsService,
	}
}

func (u *Usecase) Run(ctx context.Context, fileName string) (presignedUrl, objectUrl string, err error) {
	return u.contentsService.GenerateThumbnailPutURL(fileName)
}
