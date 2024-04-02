package handler

import (
	"fmt"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_blogs"
	"net/http"
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

	isPublic := func() *bool { var v = true; return &v }()
	input := get_blogs.NewGetBlogsInput(isPublic, nil, nil)
	v := r.URL.Query()
	tag := v.Get("tag")
	if tag != "" {
		input.Tag = &tag
	}
	keyword := v.Get("keyword")
	if keyword != "" {
		input.KeyWord = &keyword
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
