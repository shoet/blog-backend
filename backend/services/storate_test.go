package services

import (
	"fmt"
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
	url, err := s.GenerateThumbnailPutURL("test.jpg")
	if err != nil {
		t.Fatalf("failed to generate url: %v", err)
	}
	fmt.Println(url)
}
