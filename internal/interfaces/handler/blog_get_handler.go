package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_blog_detail"
)

type BlogGetHandler struct {
	Usecase *get_blog_detail.Usecase
	jwter   JWTService
}

func NewBlogGetHandler(usecase *get_blog_detail.Usecase, jwter JWTService) *BlogGetHandler {
	return &BlogGetHandler{
		Usecase: usecase,
		jwter:   jwter,
	}
}

type BlogGetResponse struct {
	*models.Blog
	Comments []*models.Comment `json:"comments"`
}

func (l *BlogGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error("failed to get id from url")
		response.ResponsdBadRequest(w, r, nil)
		return
	}
	idInt, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to convert id to int: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}
	blog, comments, err := l.Usecase.Run(ctx, models.BlogId(idInt))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	if blog == nil {
		response.ResponsdNotFound(w, r, err)
		return
	}
	// 非公開のBlogは認証が必要
	if !blog.IsPublic {
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error("failed to get authorization header")
			response.ResponsdNotFound(w, r, err)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			logger.Error("failed to get authorization token")
			response.ResponsdNotFound(w, r, err)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		_, err := l.jwter.VerifyToken(ctx, token)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to verify token: %v", err))
			response.ResponsdNotFound(w, r, err)
			return
		}
	}
	res := &BlogGetResponse{
		Blog: blog,
	}
	if comments != nil {
		res.Comments = comments
	}
	if err := response.RespondJSON(w, r, http.StatusOK, res); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
