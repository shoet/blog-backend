package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) error {
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

type ErrorResponse struct {
	Message string `json:"message"`
}
