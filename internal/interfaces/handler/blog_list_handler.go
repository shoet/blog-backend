package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_blogs"
)

type BlogListHandler struct {
	Usecase *get_blogs.Usecase
}

func NewBlogListHandler(usecase *get_blogs.Usecase) *BlogListHandler {
	return &BlogListHandler{
		Usecase: usecase,
	}
}

func (l *BlogListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)

	isPublicOnly := func() *bool { var v = true; return &v }()
	input := &get_blogs.GetBlogsInput{IsPublicOnly: isPublicOnly}
	v := r.URL.Query()
	tag := v.Get("tag")
	if tag != "" {
		input.Tag = &tag
	}
	keyword := v.Get("keyword")
	if keyword != "" {
		input.KeyWord = &keyword
	}
	offsetBlogId := v.Get("offset_id")
	if offsetBlogId != "" {
		v, err := strconv.Atoi(offsetBlogId)
		if err != nil {
			err := fmt.Errorf("offset_id is invalid")
			logger.Error(err.Error())
			response.ResponsdBadRequest(w, r, err)
			return
		}
		blogId := models.BlogId(v)
		input.OffsetBlogId = &blogId
	}
	limit := v.Get("limit")
	if limit != "" {
		v, err := strconv.Atoi(limit)
		if err != nil {
			err := fmt.Errorf("limit is invalid")
			logger.Error(err.Error())
			response.ResponsdBadRequest(w, r, err)
			return
		}
		l := uint(v)
		input.Limit = &l
	}

	blogs, err := l.Usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	if blogs == nil {
		if err := response.RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}

	if err := response.RespondJSON(w, r, http.StatusOK, blogs); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
