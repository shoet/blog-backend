package adapter

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/shoet/blog/internal/config"
)

type AWSS3StorageAdapter struct {
	config        *config.Config
	S3Client      *s3.Client
	PresignClient *s3.PresignClient
}

func NewAWSS3StorageAdapter(cfg *config.Config) (*AWSS3StorageAdapter, error) {
	sdkConfig, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}
	s3Client := s3.NewFromConfig(sdkConfig)
	return &AWSS3StorageAdapter{
		config:        cfg,
		S3Client:      s3Client,
		PresignClient: s3.NewPresignClient(s3Client),
	}, nil
}

func (s *AWSS3StorageAdapter) ExistObject(ctx context.Context, bucketName string, key string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}
	if _, err := s.S3Client.HeadObject(ctx, input); err != nil {
		var genericError *smithy.GenericAPIError
		if errors.As(err, &genericError) {
			if genericError.ErrorCode() == "404" {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to HeadObject: %w", err)
	}
	return true, nil
}

func (s *AWSS3StorageAdapter) GeneratePreSignedURL(destinationPath string, fileName string) (presignedUrl, objectUrl string, err error) {
	bucketName := s.config.AWSS3Bucket
	objectKey := filepath.Join(destinationPath, fileName)
	request, err := s.GenerateSignedURL(bucketName, objectKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate signed url: %w", err)
	}
	objectURL := fmt.Sprintf(
		"https://%s/%s/%s",
		s.config.CdnDomain,
		destinationPath,
		fileName,
	)
	return request.URL, objectURL, nil
}

func (s *AWSS3StorageAdapter) GenerateSignedURL(
	bucketName string, objectKey string,
) (presignedRequest *v4.PresignedHTTPRequest, err error) {
	request, err := s.PresignClient.PresignPutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(s.config.AWSS3PresignPutExpiresSec * int64(time.Second))
		})
	if err != nil {
		return nil, fmt.Errorf(
			"couldn't get a presigned request to put %s:%s. Here's why: %v", bucketName, objectKey, err)
	}
	return request, err
}
