package services

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/shoet/blog/config"
)

type AWSS3StorageService struct {
	config        *config.Config
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}

func NewAWSS3StorageService(cfg *config.Config) (*AWSS3StorageService, error) {
	sdkConfig, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	return &AWSS3StorageService{
		config:        cfg,
		S3Client:      s3Client,
		PresignClient: s3.NewPresignClient(s3Client),
	}, nil
}

func (s *AWSS3StorageService) GenerateThumbnailPutURL(fileName string) (string, string, error) {
	bucketName := s.config.AWSS3Bucket
	objectKey := filepath.Join(s.config.AWSS3ThumbnailDirectory, fileName)
	request, err := s.GenerateSignedURL(bucketName, objectKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate signed url: %w", err)
	}
	objectURL := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s/%s",
		bucketName,
		s.config.AWSS3Region,
		s.config.AWSS3ThumbnailDirectory,
		fileName,
	)
	return request.URL, objectURL, nil
}

func (s *AWSS3StorageService) GenerateSignedURL(
	bucketName string, objectKey string,
) (*v4.PresignedHTTPRequest, error) {
	request, err := s.PresignClient.PresignPutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(s.config.AWSS3PresignPutExpiresSec * int64(time.Second))
		})
	if err != nil {
		fmt.Errorf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}
