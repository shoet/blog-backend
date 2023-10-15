package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/options"
)

type BlogListHandler struct {
	Service BlogService
}

func (l *BlogListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authorId := models.UserId(1) // TODO
	option := options.ListBlogOptions{
		AuthorId: authorId,
	}
	resp, err := l.Service.ListBlog(ctx, option)
	if err != nil {
		log.Printf("failed to list blog: %v", err)
		resp := ErrorResponse{Message: ErrMessageInternalServerError}
		if err := RespondJSON(w, http.StatusInternalServerError, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}
	if resp == nil {
		if err := RespondJSON(w, http.StatusOK, []interface{}{}); err != nil {
			log.Printf("failed to respond json response: %v", err)
		}
		return
	}
	if err := RespondJSON(w, http.StatusOK, resp); err != nil {
		log.Printf("failed to respond json response: %v", err)
	}
	return
}

type BlogGetHandler struct {
	Service BlogService
}

func (l *BlogGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Printf("failed to get id from url param")
		resp := ErrorResponse{Message: ErrMessageBadRequest}
		if err := RespondJSON(w, http.StatusBadRequest, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}
	idInt, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		log.Printf("failed to convert id to int: %v", err)
		resp := ErrorResponse{Message: ErrMessageBadRequest}
		if err := RespondJSON(w, http.StatusBadRequest, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}
	blog, err := l.Service.GetBlog(ctx, models.BlogId(idInt))
	if blog == nil {
		resp := ErrorResponse{Message: ErrMessageNotFound}
		if err := RespondJSON(w, http.StatusNotFound, resp); err != nil {
			log.Printf("failed to respond json response: %v", err)
		}
		return
	}
	if err := RespondJSON(w, http.StatusOK, blog); err != nil {
		log.Printf("failed to respond json response: %v", err)
	}
	return
}

type BlogAddHandler struct {
	Service   BlogService
	Validator *validator.Validate
}

func (a *BlogAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqBody struct {
		Title                  string        `json:"title" validate:"required"`
		Content                string        `json:"content" validate:"required"`
		Description            string        `json:"description" validate:"required"`
		AuthorId               models.UserId `json:"authorId" validate:"required"`
		ThumbnailImageFileName string        `json:"thumbnailImage_file_name"`
		IsPublic               bool          `json:"isPublic" default:"false"`
		Tags                   []string      `json:"tags" default:"[]"`
	}
	defer r.Body.Close()
	if err := JsonToStruct(r, &reqBody); err != nil {
		log.Printf("failed to parse request body: %v", err)
		resp := ErrorResponse{Message: ErrMessageInternalServerError}
		if err := RespondJSON(w, http.StatusInternalServerError, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		log.Printf("failed to validate request body: %v", err)
		resp := ErrorResponse{Message: ErrMessageBadRequest}
		if err := RespondJSON(w, http.StatusBadRequest, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
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

	err := a.Service.AddBlog(ctx, blog)
	if err != nil {
		resp := ErrorResponse{Message: ErrMessageInternalServerError}
		if err := RespondJSON(w, http.StatusInternalServerError, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}
	resp := struct {
		Id int `json:"id"`
	}{
		Id: int(blog.Id),
	}
	if err := RespondJSON(w, http.StatusOK, resp); err != nil {
		log.Printf("failed to respond json response: %v", err)
	}
	return
}

type BlogDeleteHandler struct {
	Service BlogService
}

func (d *BlogDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqBody struct {
		Id models.BlogId `json:"id"`
	}
	defer r.Body.Close()
	if err := JsonToStruct(r, reqBody); err != nil {
		resp := ErrorResponse{Message: ErrMessageBadRequest}
		if err := RespondJSON(w, http.StatusBadRequest, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}

	err := d.Service.DeleteBlog(ctx, reqBody.Id)
	if err != nil {
		resp := ErrorResponse{Message: ErrMessageInternalServerError}
		if err := RespondJSON(w, http.StatusInternalServerError, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}
	resp := struct {
		Id int `json:"id"`
	}{
		Id: int(reqBody.Id),
	}
	if err := RespondJSON(w, http.StatusOK, resp); err != nil {
		log.Printf("failed to respond json response: %v", err)
	}
	return
}

type BlogPutHandler struct {
	Service BlogService
}

func (p *BlogPutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqBody struct {
		Title                  string   `json:"title"`
		Content                string   `json:"content"`
		ThumbnailImageFileName string   `json:"thumbnailImage_file_name"`
		IsPublic               bool     `json:"isPublic"`
		Tags                   []string `json:"tags"`
	}
	defer r.Body.Close()
	if err := JsonToStruct(r, reqBody); err != nil {
		resp := ErrorResponse{Message: ErrMessageBadRequest}
		if err := RespondJSON(w, http.StatusBadRequest, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}

	blog := &models.Blog{
		Title:                  reqBody.Title,
		Content:                reqBody.Content,
		ThumbnailImageFileName: reqBody.ThumbnailImageFileName,
		IsPublic:               reqBody.IsPublic,
		Tags:                   reqBody.Tags,
	}

	err := p.Service.PutBlog(ctx, blog)
	if err != nil {
		resp := ErrorResponse{Message: ErrMessageInternalServerError}
		if err := RespondJSON(w, http.StatusInternalServerError, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}
	resp := struct {
		Id int `json:"id"`
	}{
		Id: int(blog.Id),
	}
	if err := RespondJSON(w, http.StatusOK, resp); err != nil {
		log.Printf("failed to respond json response: %v", err)
	}
	return
}
