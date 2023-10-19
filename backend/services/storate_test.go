package services

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shoet/blog/config"
)

func Test_AWSStorageService_GenerateThumbnailPutURL(t *testing.T) {
	// TODO: 自動テストの際にサービス連携部分はどうするのか考える
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	s, err := NewAWSS3StorageService(cfg)
	if err != nil {
		t.Fatalf("failed to create aws storage service: %v", err)
	}
	wantFileName := "test.jpg"
	signedUrl, objectUrl, err := s.GenerateThumbnailPutURL(wantFileName)
	if err != nil {
		t.Fatalf("failed to generate url: %v", err)
	}
	if strings.HasSuffix(signedUrl, wantFileName) {
		t.Fatalf("signed url is not expected: %v", signedUrl)
	}
	fmt.Println(signedUrl)
	fmt.Println(objectUrl)
}
