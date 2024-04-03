package handler

import (
	"fmt"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_blogs"
	"net/http"
)

type BlogListAdminHandler struct {
	Usecase *get_blogs.Usecase
}

func NewBlogListAdminHandler(usecase *get_blogs.Usecase) *BlogListAdminHandler {
	return &BlogListAdminHandler{
		Usecase: usecase,
	}
}

func (l *BlogListAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)

	input := &get_blogs.GetBlogsInput{}
	resp, err := l.Usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	if resp == nil {
		if err := response.RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
