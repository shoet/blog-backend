package handlers

import (
	"log"
	"net/http"

	"github.com/shoet/blog/models"
)

// TODO: validate

type BlogListHandler struct {
	Service BlogService
}

func (l *BlogListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := l.Service.ListBlog(nil)
	if err != nil {
		resp := ErrorResponse{Message: ErrMessageInternalServerError}
		if err := RespondJSON(w, http.StatusInternalServerError, resp); err != nil {
			log.Printf("failed to respond json error: %v", err)
		}
		return
	}
	if err := RespondJSON(w, http.StatusOK, resp); err != nil {
		log.Printf("failed to respond json response: %v", err)
	}
	return
}

type BlogAddHandler struct {
	Service BlogService
}

func (a *BlogAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Title                  string        `json:"title"`
		Content                string        `json:"content"`
		AuthorId               models.UserId `json:"authorId"`
		ThumbnailImageFileName string        `json:"thumbnailImage_file_name"`
		IsPublic               bool          `json:"isPublic"`
		Tags                   []models.Tag  `json:"tags"`
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
		AuthorId:               reqBody.AuthorId,
		ThumbnailImageFileName: reqBody.ThumbnailImageFileName,
		IsPublic:               reqBody.IsPublic,
		Tags:                   reqBody.Tags,
	}

	err := a.Service.AddBlog(blog)
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

	err := d.Service.DeleteBlog(reqBody.Id)
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
	var reqBody struct {
		Title                  string       `json:"title"`
		Content                string       `json:"content"`
		ThumbnailImageFileName string       `json:"thumbnailImage_file_name"`
		IsPublic               bool         `json:"isPublic"`
		Tags                   []models.Tag `json:"tags"`
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

	err := p.Service.PutBlog(blog)
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
