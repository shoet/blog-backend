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
	"github.com/shoet/blog/internal/usecase/delete_blog"
)

type BlogDeleteHandler struct {
	Usecase   *delete_blog.Usecase
	Validator *validator.Validate
}

func NewBlogDeleteHandler(
	usecase *delete_blog.Usecase,
	validator *validator.Validate,
) *BlogDeleteHandler {
	return &BlogDeleteHandler{
		Usecase:   usecase,
		Validator: validator,
	}
}

func (d *BlogDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	blogId, err := d.Usecase.Run(ctx, models.BlogId(idInt))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to delete blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	resp := struct {
		Id int `json:"id"`
	}{
		Id: int(blogId),
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
