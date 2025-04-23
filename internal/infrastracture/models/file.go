package models

import (
	"fmt"

	"github.com/shoet/blog/internal/config"
)

type FileType string

const (
	FileTypeAvatarImage      = FileType("avatar_image")
	FileTypeThumbnailImage   = FileType("thumbnail_image")
	FileTypeBlogContentImage = FileType("blog_content_image")
)

type File struct {
	Type     FileType
	FileName string
}

func NewFile(fileType FileType, name string) (*File, error) {
	switch fileType {
	case FileTypeAvatarImage:
		break
	default:
		return nil, fmt.Errorf("invalida FileType: %s", fileType)
	}
	return &File{
		Type: fileType, FileName: name,
	}, nil
}

func (f *File) GetBucketName(config *config.Config) (string, error) {
	switch f.Type {
	case FileTypeAvatarImage:
		return config.AWSS3Bucket, nil
	case FileTypeThumbnailImage:
		return config.AWSS3Bucket, nil
	case FileTypeBlogContentImage:
		return config.AWSS3Bucket, nil
	default:
		return "", fmt.Errorf("not found bucket")
	}
}

func (f *File) GetBucketKey(config *config.Config) (string, error) {
	switch f.Type {
	case FileTypeAvatarImage:
		return config.AWSS3AvatarImageDirectory, nil
	case FileTypeThumbnailImage:
		return config.AWSS3ThumbnailDirectory, nil
	case FileTypeBlogContentImage:
		return config.AWSSS3ContentImageDirectory, nil
	default:
		return "", fmt.Errorf("not found key")
	}
}
