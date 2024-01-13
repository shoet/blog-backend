package services

import (
	"strings"
	"testing"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/testutil"
)

func Test_AWSStorageService_GenerateThumbnailPutURL(t *testing.T) {
	testutil.LoadDotenvForTest(t)
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
	if !strings.Contains(objectUrl, wantFileName) {
		t.Fatalf("object url is not expected: %v", objectUrl)
	}
	if !strings.Contains(signedUrl, wantFileName) {
		t.Fatalf("signed url is not expected: %v", signedUrl)
	}
}
