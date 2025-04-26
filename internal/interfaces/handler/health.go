package handler

import (
	"net/http"

	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
)

type HealthCheckHandler struct{}

type ResponseHealthCheck struct {
	Message string `json:"message"`
}

func (hh *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	logger.Info("health check")
	resp := &ResponseHealthCheck{Message: "OK"}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		errResp := &response.Response{Message: "NG"}
		if err := response.RespondJSON(w, r, http.StatusInternalServerError, errResp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
