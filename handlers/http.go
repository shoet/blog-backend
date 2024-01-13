package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shoet/blog/internal/logging"
)

func RespondJSON(w http.ResponseWriter, r *http.Request, statusCode int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body in RespondJSON(): %w", err)
	}
	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("failed to write body in RespondJSON(): %w", err)
	}
	return nil
}

func ResponsdBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	resp := ErrorResponse{Message: ErrMessageBadRequest}
	if err := RespondJSON(w, r, http.StatusBadRequest, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json error: %v", err))
	}
}

func ResponsdNotFound(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	resp := ErrorResponse{Message: ErrMessageNotFound}
	if err := RespondJSON(w, r, http.StatusNotFound, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json error: %v", err))
	}
}

func ResponsdInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	resp := ErrorResponse{Message: ErrMessageInternalServerError}
	if err := RespondJSON(w, r, http.StatusInternalServerError, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json error: %v", err))
	}
}

func RespondUnauthorized(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	resp := ErrorResponse{Message: ErrMessageUnauthorized}
	if err := RespondJSON(w, r, http.StatusUnauthorized, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json error: %v", err))
	}
}

func JsonToStruct(r *http.Request, v any) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode json in JsonToStruct(): %w", err)
	}
	return nil
}

type ErrorResponse struct {
	Message string `json:"message"`
}

const (
	ErrMessageBadRequest          = "BadRequest"
	ErrMessageNotFound            = "NotFound"
	ErrMessageInternalServerError = "InternalServerError"
	ErrMessageUnauthorized        = "Unauthorized"
)
