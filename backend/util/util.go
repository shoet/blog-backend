package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectRoot(cwd string) (string, error) {
	if cwd == "/" {
		return "", fmt.Errorf("go.mod not found")
	}
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		return cwd, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to stat go.mod: %w", err)
	}
	return GetProjectRoot(filepath.Dir(cwd))
}
