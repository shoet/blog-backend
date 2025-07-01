package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/get_privacy_policy"
	"github.com/shoet/blog/internal/usecase/put_privacy_policy"
)

type GetPrivacyPolicyHandler struct {
	Usecase *get_privacy_policy.Usecase
}

func NewGetPrivacyPolicyHandler(
	usecase *get_privacy_policy.Usecase,
) *GetPrivacyPolicyHandler {
	return &GetPrivacyPolicyHandler{
		Usecase: usecase,
	}
}

func (h *GetPrivacyPolicyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	id := chi.URLParam(r, "id")

	privacyPolicy, err := h.Usecase.Run(ctx, id)
	if err != nil {
		if errors.Is(err, get_privacy_policy.ErrResourceNotFound) {
			logger.Error(fmt.Sprintf("privacy policy not found: %v", err))
			response.RespondNotFound(w, r, err)
			return
		}
		logger.Error(fmt.Sprintf("failed to run usecase: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}

	res := struct {
		PrivacyPolicy *models.PrivacyPolicy `json:"privacy_policy"`
	} {
		PrivacyPolicy: privacyPolicy,
	}

	if err := response.RespondJSON(w, r, http.StatusOK, res); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type PutPrivacyPolicyHandler struct {
	Usecase *put_privacy_policy.Usecase
	Validator *validator.Validate
}

func NewPutPrivacyPolicyHandler(
	usecase *put_privacy_policy.Usecase,
	validator *validator.Validate,
) *PutPrivacyPolicyHandler {
	return &PutPrivacyPolicyHandler{
		Usecase: usecase,
		Validator: validator,
	}
}

func (h *PutPrivacyPolicyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	id := chi.URLParam(r, "id")

	var reqBody struct {
		Content                string        `json:"content" validate:"required"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := h.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := h.Usecase.Run(ctx, id, reqBody.Content); err != nil {
		logger.Error(fmt.Sprintf("failed to run usecase: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}

	res := struct {
		Message string `json:"message"`
	} {
		Message: "ok",
	}

	if err := response.RespondJSON(w, r, http.StatusOK, res); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
