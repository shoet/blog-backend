package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/put_blog"
	"net/http"
)

type BlogPutHandler struct {
	Usecase   *put_blog.Usecase
	Validator *validator.Validate
}

func NewBlogPutHandler(
	usecase *put_blog.Usecase,
	validator *validator.Validate,
) *BlogPutHandler {
	return &BlogPutHandler{
		Usecase:   usecase,
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
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := p.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
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

	newBlog, err := p.Usecase.Run(ctx, blog)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to put blog: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, newBlog); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
