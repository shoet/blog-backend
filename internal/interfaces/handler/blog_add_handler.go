package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/create_blog"
	"net/http"
)

type BlogAddHandler struct {
	Usecase   *create_blog.Usecase
	Validator *validator.Validate
}

func NewBlogAddHandler(
	usecase *create_blog.Usecase,
	validator *validator.Validate,
) *BlogAddHandler {
	return &BlogAddHandler{
		Usecase:   usecase,
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

	newBlog, err := a.Usecase.Run(ctx, blog)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to add blog: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
