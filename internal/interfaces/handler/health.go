package handler

import (
	"net/http"

	"github.com/shoet/blog/internal/interfaces"
)

type HealthCheckHandler struct{}

type ResponseHealthCheck struct {
	Message string `json:"message"`
}

func (hh *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseHealthCheck{Message: "OK"}
	if err := interfaces.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		errResp := &interfaces.ErrorResponse{Message: "NG"}
		interfaces.RespondJSON(w, r, http.StatusInternalServerError, errResp)
	}
}
