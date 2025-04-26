package models

import (
	"fmt"
	"net/url"
	"strings"

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

func NewFileFromURL(config *config.Config, rawURL string) (*File, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	if url.Scheme != "https" {
		return nil, fmt.Errorf("invalid url scheme: %s", url.Scheme)
	}
	if strings.HasPrefix(url.Host, config.AWSS3Bucket) {
		return nil, fmt.Errorf("invalid url host: %s", url.Host)
	}
	var file *File
	if strings.HasPrefix(url.Path, config.AWSS3AvatarImageDirectory) {
		file = &File{
			Type:     FileTypeAvatarImage,
			FileName: strings.TrimPrefix(url.Path, config.AWSS3AvatarImageDirectory),
		}
	} else if strings.HasPrefix(url.Path, config.AWSS3ThumbnailDirectory) {
		file = &File{
			Type:     FileTypeThumbnailImage,
			FileName: strings.TrimPrefix(url.Path, config.AWSS3ThumbnailDirectory),
		}
	} else if strings.HasPrefix(url.Path, config.AWSSS3ContentImageDirectory) {
		file = &File{
			Type:     FileTypeBlogContentImage,
			FileName: strings.TrimPrefix(url.Path, config.AWSSS3ContentImageDirectory),
		}
	} else {
		return nil, fmt.Errorf("invalid url path: %s", url.Path)

	}
	if strings.Contains(file.FileName, "/") {
		return nil, fmt.Errorf("invalid file name: %s", file.FileName)
	}
	return file, nil
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

func (f *File) GetFileURL(config *config.Config) (string, error) {
	key, err := f.GetBucketKey(config)
	if err != nil {
		return "", fmt.Errorf("failed to get bucket key: %w", err)
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s", config.CdnDomain, key, f.FileName), nil

}
