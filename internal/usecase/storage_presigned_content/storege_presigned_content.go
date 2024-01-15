package storage_presigned_content

import "context"

type ContentsService interface {
	GenerateContentImagePutURL(fileName string) (presignedUrl, objectUrl string, err error)
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
	return u.contentsService.GenerateContentImagePutURL(fileName)
}
