package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_handlename"
)

type GetHandlenameHandler struct {
	Usecase *get_handlename.Usecase
}

func NewGetHandlenameHandler(usecase *get_handlename.Usecase) *GetHandlenameHandler {
	return &GetHandlenameHandler{
		Usecase: usecase,
	}
}

type GetHandlenameResponse struct {
	Handlename string `json:"handlename"`
}

func (h *GetHandlenameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	blogId := r.URL.Query().Get("blogId")
	if blogId == "" {
		logger.Error(fmt.Sprintf("blogId is required"))
		response.RespondBadRequest(w, r, nil)
		return
	}
	blogIdNum, err := strconv.Atoi(blogId)
	if err != nil {
		logger.Error(fmt.Sprintf("invalid blogId: %d", blogIdNum))
		response.RespondBadRequest(w, r, nil)
		return
	}
	ip := r.Header.Get("x-forwarded-for")
	if ip == "" {
		logger.Error(fmt.Sprintf("IPAddr not found"))
		response.RespondBadRequest(w, r, nil)
		return
	}
	handlename, err := h.Usecase.Run(ctx, models.BlogId(blogIdNum), ip)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get handlename: %v", err))
		response.RespondInternalServerError(w, r, nil)
		return
	}
	res := &GetHandlenameResponse{
		Handlename: handlename,
	}
	if err := response.RespondJSON(w, r, http.StatusOK, res); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
