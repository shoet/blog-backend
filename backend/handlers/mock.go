package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/shoet/blog/util"
)

func LoadMockData(jsonPath string, buf any) error {
	f, err := os.Open(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to open file in LoadMockData(): %w", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file in LoadMockData(): %w", err)
	}
	if err := json.Unmarshal(b, buf); err != nil {
		return fmt.Errorf("failed to unmarshal json in LoadMockData(): %w", err)
	}
	return nil
}

func RespondMockJSON(mockFileName string, buf any, w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		resp := ErrorResponse{
			Message: "failed to get current working directory",
		}
		RespondJSON(w, http.StatusInternalServerError, resp)
		return
	}
	projectRoot, err := util.GetProjectRoot(cwd)
	if err != nil {
		resp := ErrorResponse{
			Message: "failed to get current working directory",
		}
		RespondJSON(w, http.StatusInternalServerError, resp)
		return
	}
	if err := LoadMockData(
		filepath.Join(projectRoot, "handlers/data", mockFileName),
		buf,
	); err != nil {
		resp := ErrorResponse{
			Message: fmt.Sprintf("failed to load mock data: %s", err.Error()),
		}
		RespondJSON(w, http.StatusInternalServerError, resp)
		return
	}

	RespondJSON(w, http.StatusOK, buf)
	return
}
