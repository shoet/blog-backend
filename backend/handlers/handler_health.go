package handlers

import "net/http"

type HealthCheckHandler struct{}

type ResponseHealthCheck struct {
	Message string `json:"message"`
}

func (hh *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseHealthCheck{Message: "OK"}
	if err := RespondJSON(w, http.StatusOK, resp); err != nil {
		errResp := &ErrorResponse{Message: "NG"}
		RespondJSON(w, http.StatusInternalServerError, errResp)
	}
}
