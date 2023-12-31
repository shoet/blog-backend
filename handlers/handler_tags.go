package handlers

import (
	"fmt"
	"net/http"

	"github.com/shoet/blog/logging"
	"github.com/shoet/blog/options"
)

type TagListHandler struct {
	Service BlogManager
}

func NewTagListHandler(service BlogManager) *TagListHandler {
	return &TagListHandler{
		Service: service,
	}
}

func (t *TagListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	option := options.ListTagsOptions{
		Limit: 100,
	}
	resp, err := t.Service.ListTags(ctx, option)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list tags: %v", err))
		ResponsdInternalServerError(w, r, err)
		return
	}
	if resp == nil {
		if err := RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
	return
}
