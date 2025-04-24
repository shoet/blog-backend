package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/storage_presigned_content"
	"github.com/shoet/blog/internal/usecase/storage_presigned_thumbnail"
	"github.com/shoet/blog/internal/usecase/upload_file"
)

type UploadFileHandler struct {
	Usecase   *upload_file.Usecase
	Validator *validator.Validate
}

func NewUploadFileHandler(
	Usecase *upload_file.Usecase,
	Validator *validator.Validate,
) *UploadFileHandler {
	return &UploadFileHandler{
		Usecase:   Usecase,
		Validator: Validator,
	}
}

func (g *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		FileName string `json:"fileName" validate:"required"`
		FileType string `json:"fileType" validate:"required"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := g.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	uploadURL, destinationUrl, err := g.Usecase.Run(ctx, upload_file.UploadFileInput{
		FileType: reqBody.FileType,
		FileName: reqBody.FileName,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}

	resp := struct {
		UploadURL      string `json:"uploadUrl"`
		DestinationURL string `json:"destinationUrl"`
	}{
		UploadURL:      uploadURL,
		DestinationURL: destinationUrl,
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
	}
}

type GenerateThumbnailImageSignedURLHandler struct {
	Usecase   *storage_presigned_thumbnail.Usecase
	Validator *validator.Validate
}

func NewGenerateThumbnailImageSignedURLHandler(
	Usecase *storage_presigned_thumbnail.Usecase,
	Validator *validator.Validate,
) *GenerateThumbnailImageSignedURLHandler {
	return &GenerateThumbnailImageSignedURLHandler{
		Usecase:   Usecase,
		Validator: Validator,
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
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := g.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	signedUrl, destinationUrl, err := g.Usecase.Run(ctx, reqBody.FileName)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondInternalServerError(w, r, err)
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
	Usecase   *storage_presigned_content.Usecase
	Validator *validator.Validate
}

func NewGenerateContentsImageSignedURLHandler(
	Usecase *storage_presigned_content.Usecase,
	Validator *validator.Validate,
) *GenerateContentsImageSignedURLHandler {
	return &GenerateContentsImageSignedURLHandler{
		Usecase:   Usecase,
		Validator: Validator,
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
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := g.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	signedUrl, destinationUrl, err := g.Usecase.Run(ctx, reqBody.FileName)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondInternalServerError(w, r, err)
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
