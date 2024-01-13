package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/shoet/blog/util"
)

func LoadDotenvForTest(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed get cwd: %v", err)
	}
	rootDir, err := util.GetProjectRoot(cwd)
	if err != nil {
		t.Fatalf("failed get project root: %v", err)
	}
	ci, ok := os.LookupEnv("CI")
	if ok && ci == "true" {
		return
	}
	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		t.Fatalf("failed load .env: %v", err)
	}
}
