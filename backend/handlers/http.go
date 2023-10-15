package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, statusCode int, body any) error {
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

var (
	ErrMessageBadRequest          = "BadRequest"
	ErrMessageNotFound            = "NotFound"
	ErrMessageInternalServerError = "InternalServerError"
	ErrMessageUnauthorized        = "Unauthorized"
)
