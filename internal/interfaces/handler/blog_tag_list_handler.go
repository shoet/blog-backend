package handler

import (
	"fmt"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/options"
	"github.com/shoet/blog/internal/usecase/get_tags"
	"net/http"
)

type TagListHandler struct {
	Usecase get_tags.Usecase
}

func NewTagListHandler(
	usecase get_tags.Usecase,
) *TagListHandler {
	return &TagListHandler{
		Usecase: usecase,
	}
}

func (t *TagListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	option := options.ListTagsOptions{
		Limit: 100,
	}
	resp, err := t.Usecase.Run(ctx, option)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list tags: %v", err))
		response.RespondInternalServerError(w, r, err)
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
