package adapter_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastructure/adapter"
	"github.com/shoet/blog/internal/testutil"
)

func Test_AWSStorageService_GenerateThumbnailPutURL(t *testing.T) {
	testutil.LoadDotenvForTest(t)
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}
	s, err := adapter.NewS3Adapter(cfg)
	if err != nil {
		t.Fatalf("failed to create aws storage service: %v", err)
	}
	wantFileName := "test.jpg"
	wantDirectoryName := "test"
	wantPath := filepath.Join(wantDirectoryName, wantFileName)
	signedUrl, objectUrl, err := s.GeneratePreSignedURL(wantDirectoryName, wantFileName)
	if err != nil {
		t.Fatalf("failed to generate url: %v", err)
	}
	if !strings.Contains(objectUrl, wantPath) {
		t.Fatalf("object url is not expected: %v", objectUrl)
	}
	if !strings.Contains(signedUrl, wantPath) {
		t.Fatalf("signed url is not expected: %v", signedUrl)
	}
}
