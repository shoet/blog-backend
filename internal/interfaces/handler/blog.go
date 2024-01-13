package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"

	"github.com/shoet/blog/internal/infrastracture/services"
	"github.com/shoet/blog/internal/options"
)

type BlogListHandler struct {
	Service BlogManager
}

func NewBlogListHandler(blogService BlogManager) *BlogListHandler {
	return &BlogListHandler{
		Service: blogService,
	}
}

func (l *BlogListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	option := options.ListBlogOptions{
		IsPublic: true,
	}
	resp, err := l.Service.ListBlog(ctx, option)
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

type BlogGetHandler struct {
	Service BlogManager
	jwter   services.JWTer
}

func NewBlogGetHandler(blogService BlogManager, jwter services.JWTer) *BlogGetHandler {
	return &BlogGetHandler{
		Service: blogService,
		jwter:   jwter,
	}
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
	blog, err := l.Service.GetBlog(ctx, models.BlogId(idInt))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	if blog == nil {
		response.ResponsdNotFound(w, r, err)
		return
	}
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
	if err := response.RespondJSON(w, r, http.StatusOK, blog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogAddHandler struct {
	Service   BlogManager
	Validator *validator.Validate
}

func NewBlogAddHandler(
	blogService BlogManager, validator *validator.Validate,
) *BlogAddHandler {
	return &BlogAddHandler{
		Service:   blogService,
		Validator: validator,
	}
}

func (a *BlogAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		Title                  string        `json:"title" validate:"required"`
		Content                string        `json:"content" validate:"required"`
		Description            string        `json:"description" validate:"required"`
		AuthorId               models.UserId `json:"authorId" validate:"required"`
		ThumbnailImageFileName string        `json:"thumbnailImageFileName"`
		IsPublic               bool          `json:"isPublic" default:"false"`
		Tags                   []string      `json:"tags" default:"[]"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	blog := &models.Blog{
		Title:                  reqBody.Title,
		Content:                reqBody.Content,
		Description:            reqBody.Description,
		AuthorId:               reqBody.AuthorId,
		ThumbnailImageFileName: reqBody.ThumbnailImageFileName,
		IsPublic:               reqBody.IsPublic,
		Tags:                   reqBody.Tags,
	}

	newBlog, err := a.Service.AddBlog(ctx, blog)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogDeleteHandler struct {
	Service   BlogManager
	Validator *validator.Validate
}

func NewBlogDeleteHandler(
	blogService BlogManager, validator *validator.Validate,
) *BlogDeleteHandler {
	return &BlogDeleteHandler{
		Service:   blogService,
		Validator: validator,
	}
}

func (d *BlogDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		Id models.BlogId `json:"id" validate:"required"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	if err := d.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	err := d.Service.DeleteBlog(ctx, reqBody.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	resp := struct {
		Id int `json:"id"`
	}{
		Id: int(reqBody.Id),
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogPutHandler struct {
	Service   BlogManager
	Validator *validator.Validate
}

func NewBlogPutHandler(
	blogService BlogManager, validator *validator.Validate,
) *BlogPutHandler {
	return &BlogPutHandler{
		Service:   blogService,
		Validator: validator,
	}
}

func (p *BlogPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		Id                     models.BlogId `json:"id" validate:"required"`
		AuthorId               models.UserId `json:"authorId" validate:"required"`
		Title                  string        `json:"title"`
		Content                string        `json:"content"`
		Description            string        `json:"description"`
		ThumbnailImageFileName string        `json:"thumbnailImageFileName"`
		IsPublic               bool          `json:"isPublic"`
		Tags                   []string      `json:"tags"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	if err := p.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	blog := &models.Blog{
		Id:                     reqBody.Id,
		AuthorId:               reqBody.AuthorId,
		Title:                  reqBody.Title,
		Content:                reqBody.Content,
		Description:            reqBody.Description,
		ThumbnailImageFileName: reqBody.ThumbnailImageFileName,
		IsPublic:               reqBody.IsPublic,
		Tags:                   reqBody.Tags,
	}

	newBlog, err := p.Service.PutBlog(ctx, blog)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to put blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type BlogListAdminHandler struct {
	Service BlogManager
}

func NewBlogListAdminHandler(blogService BlogManager) *BlogListAdminHandler {
	return &BlogListAdminHandler{
		Service: blogService,
	}
}

func (l *BlogListAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	option := options.ListBlogOptions{}
	resp, err := l.Service.ListBlog(ctx, option)
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
