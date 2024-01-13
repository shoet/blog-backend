package handler

import (
	"net/http"

	"github.com/shoet/blog/internal/interfaces/response"
)

type HealthCheckHandler struct{}

type ResponseHealthCheck struct {
	Message string `json:"message"`
}

func (hh *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseHealthCheck{Message: "OK"}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		errResp := &response.ErrorResponse{Message: "NG"}
		response.RespondJSON(w, r, http.StatusInternalServerError, errResp)
	}
}
