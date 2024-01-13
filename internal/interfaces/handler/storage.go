package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
)

type GenerateThumbnailImageSignedURLHandler struct {
	StorageService Storager
	Validator      *validator.Validate
}

func NewGenerateThumbnailImageSignedURLHandler(
	StorageService Storager,
	Validator *validator.Validate,
) *GenerateThumbnailImageSignedURLHandler {
	return &GenerateThumbnailImageSignedURLHandler{
		StorageService: StorageService,
		Validator:      Validator,
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

	signedUrl, destinationUrl, err := g.StorageService.GenerateThumbnailPutURL(reqBody.FileName)
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
	StorageService Storager
	Validator      *validator.Validate
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

	signedUrl, destinationUrl, err := g.StorageService.GenerateContentImagePutURL(reqBody.FileName)
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
