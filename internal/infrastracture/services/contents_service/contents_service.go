package contents_service

type S3StorageAdapter interface {
	GeneratePreSignedURL(destinationPath string, fileName string) (presignedUrl, objectUrl string, err error)
}

type ContentsService struct {
	s3adapter             S3StorageAdapter
	thumbnailDirectory    string
	contentImageDirectory string
}

func NewContentsService(
	s3adapter S3StorageAdapter,
	thumnailDirectory string,
	contentImageDirectory string,
) (*ContentsService, error) {
	return &ContentsService{
		s3adapter:             s3adapter,
		thumbnailDirectory:    thumnailDirectory,
		contentImageDirectory: contentImageDirectory,
	}, nil
}

// GeneratePutURL generates a signed url for put object.
func (c *ContentsService) GenerateThumbnailPutURL(fileName string) (presignedUrl, objectUrl string, err error) {
	return c.s3adapter.GeneratePreSignedURL(c.thumbnailDirectory, fileName)
}

func (c *ContentsService) GenerateContentImagePutURL(fileName string) (presignedUrl, objectUrl string, err error) {
	return c.s3adapter.GeneratePreSignedURL(c.contentImageDirectory, fileName)
}
