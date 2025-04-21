package handler

import (
	"fmt"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_blogs_offset_paging"
	"net/http"
	"strconv"
)

type BlogGetOffsetPagingHandler struct {
	Usecase *get_blogs_offset_paging.Usecase
}

func NewBlogGetOffsetPagingHandler(usecase *get_blogs_offset_paging.Usecase) *BlogGetOffsetPagingHandler {
	return &BlogGetOffsetPagingHandler{
		Usecase: usecase,
	}
}

func (l *BlogGetOffsetPagingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)

	isPublicOnly := func() *bool { var v = true; return &v }()
	input := &get_blogs_offset_paging.Input{IsPublicOnly: isPublicOnly}
	v := r.URL.Query()
	tag := v.Get("tag")
	if tag != "" {
		input.Tag = &tag
	}
	keyword := v.Get("keyword")
	if keyword != "" {
		input.KeyWord = &keyword
	}
	limit := v.Get("limit")
	if limit != "" {
		v, err := strconv.Atoi(limit)
		if err != nil {
			err := fmt.Errorf("limit is invalid")
			logger.Error(err.Error())
			response.RespondBadRequest(w, r, err)
			return
		}
		l := int64(v)
		input.Limit = &l
	}
	page := v.Get("page")
	if page != "" {
		v, err := strconv.Atoi(page)
		if err != nil {
			err := fmt.Errorf("page is invalid")
			logger.Error(err.Error())
			response.RespondBadRequest(w, r, err)
			return
		}
		p := int64(v)
		input.Page = &p
	}

	blogs, blogsCount, err := l.Usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to list blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if blogs == nil {
		if err := response.RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
		}
		return
	}

	type ResponseBody struct {
		Blog       []*models.Blog `json:"blogs"`
		TotalCount int64          `json:"totalCount"`
	}

	body := &ResponseBody{
		Blog:       blogs,
		TotalCount: blogsCount,
	}

	if err := response.RespondJSON(w, r, http.StatusOK, body); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
