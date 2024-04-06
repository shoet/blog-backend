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
	cursor_id := v.Get("cursor_id") // ページネーションのカーソルID
	if cursor_id != "" {
		v, err := strconv.Atoi(cursor_id)
		if err != nil {
			err := fmt.Errorf("cursor_id is invalid")
			logger.Error(err.Error())
			response.ResponsdBadRequest(w, r, err)
			return
		}
		blogId := models.BlogId(v)
		input.CursorId = &blogId
	}
	direction := v.Get("direction") // ページネーションの方向
	if direction != "" {
		if direction != "prev" && direction != "next" {
			err := fmt.Errorf("direction is invalid")
			logger.Error(err.Error())
			response.ResponsdBadRequest(w, r, err)
			return
		}
		input.PageDirection = &direction
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
		l := int64(v)
		input.Limit = &l
	}

	blogs, prevEOF, nextEOF, err := l.Usecase.Run(ctx, input)
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

	type ResponseBody struct {
		Blog    []*models.Blog `json:"blogs"`
		PrevEOF bool           `json:"prevEOF"`
		NextEOF bool           `json:"nextEOF"`
	}

	body := &ResponseBody{
		Blog:    blogs,
		PrevEOF: prevEOF,
		NextEOF: nextEOF,
	}

	if err := response.RespondJSON(w, r, http.StatusOK, body); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
