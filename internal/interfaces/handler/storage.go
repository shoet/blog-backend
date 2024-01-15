package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
)

type GenerateThumbnailImageSignedURLHandler struct {
	ContentsService ContentsService
	Validator       *validator.Validate
}

func NewGenerateThumbnailImageSignedURLHandler(
	ContentsService ContentsService,
	Validator *validator.Validate,
) *GenerateThumbnailImageSignedURLHandler {
	return &GenerateThumbnailImageSignedURLHandler{
		ContentsService: ContentsService,
		Validator:       Validator,
	}
}

func (g *GenerateThumbnailImageSignedURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		FileName string `json:"fileName" validate:"required"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	if err := g.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	signedUrl, destinationUrl, err := g.ContentsService.GenerateThumbnailPutURL(reqBody.FileName)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}

	resp := struct {
		SignedUrl string `json:"signedUrl"`
		PutedUrl  string `json:"putUrl"`
	}{
		SignedUrl: signedUrl,
		PutedUrl:  destinationUrl,
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
	}
}

type GenerateContentsImageSignedURLHandler struct {
	ContentsService ContentsService
	Validator       *validator.Validate
}

func NewGenerateContentsImageSignedURLHandler(
	ContentsService ContentsService,
	Validator *validator.Validate,
) *GenerateContentsImageSignedURLHandler {
	return &GenerateContentsImageSignedURLHandler{
		ContentsService: ContentsService,
		Validator:       Validator,
	}
}

func (g *GenerateContentsImageSignedURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		FileName string `json:"fileName" validate:"required"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	if err := g.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdBadRequest(w, r, err)
		return
	}

	signedUrl, destinationUrl, err := g.ContentsService.GenerateContentImagePutURL(reqBody.FileName)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.ResponsdInternalServerError(w, r, err)
		return
	}

	resp := struct {
		SignedUrl string `json:"signedUrl"`
		PutedUrl  string `json:"putUrl"`
	}{
		SignedUrl: signedUrl,
		PutedUrl:  destinationUrl,
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
	}
}
