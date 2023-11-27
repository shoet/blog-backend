package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type BlogListHandler struct {
	Service BlogManager
}

func (l *BlogListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
	option := options.ListBlogOptions{
		IsPublic: true,
	}
	resp, err := l.Service.ListBlog(ctx, option)
	if err != nil {
		logger.Error().Msgf("failed to list blog: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}
	if resp == nil {
		if err := RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error().Msgf("failed to respond json response: %v", err)
		}
		return
	}
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}

type BlogGetHandler struct {
	Service BlogManager
}

func (l *BlogGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
	id := chi.URLParam(r, "id")
	if id == "" {
		logger.Error().Msgf("failed to get id from url param")
		ResponsdBadRequest(w, r, nil)
		return
	}
	idInt, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		logger.Error().Msgf("failed to convert id to int: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}
	blog, err := l.Service.GetBlog(ctx, models.BlogId(idInt))
	if blog == nil {
		ResponsdNotFound(w, r, err)
		return
	}
	if err := RespondJSON(w, r, http.StatusOK, blog); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}

type BlogAddHandler struct {
	Service   BlogManager
	Validator *validator.Validate
}

func (a *BlogAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
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
	if err := JsonToStruct(r, &reqBody); err != nil {
		logger.Error().Msgf("failed to parse request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		logger.Error().Msgf("failed to validate request body: %v", err)
		ResponsdBadRequest(w, r, err)
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
		logger.Error().Msgf("failed to add blog: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}
	if err := RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}

type BlogDeleteHandler struct {
	Service   BlogManager
	Validator *validator.Validate
}

func (d *BlogDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
	var reqBody struct {
		Id models.BlogId `json:"id" validate:"required"`
	}
	defer r.Body.Close()
	if err := JsonToStruct(r, &reqBody); err != nil {
		logger.Error().Msgf("failed to validate request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	if err := d.Validator.Struct(reqBody); err != nil {
		logger.Error().Msgf("failed to validate request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	err := d.Service.DeleteBlog(ctx, reqBody.Id)
	if err != nil {
		logger.Error().Msgf("failed to delete blog: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}
	resp := struct {
		Id int `json:"id"`
	}{
		Id: int(reqBody.Id),
	}
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}

type BlogPutHandler struct {
	Service   BlogManager
	Validator *validator.Validate
}

func (p *BlogPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
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
	if err := JsonToStruct(r, &reqBody); err != nil {
		logger.Error().Msgf("failed to parse request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	if err := p.Validator.Struct(reqBody); err != nil {
		logger.Error().Msgf("failed to validate request body: %v", err)
		ResponsdBadRequest(w, r, err)
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
		logger.Error().Msgf("failed to put blog: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}
	if err := RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}

type BlogListAdminHandler struct {
	Service BlogManager
}

func (l *BlogListAdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
	option := options.ListBlogOptions{}
	resp, err := l.Service.ListBlog(ctx, option)
	if err != nil {
		logger.Error().Msgf("failed to list blog: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}
	if resp == nil {
		if err := RespondJSON(w, r, http.StatusOK, []interface{}{}); err != nil {
			logger.Error().Msgf("failed to respond json response: %v", err)
		}
		return
	}
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}
