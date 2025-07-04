package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/session"
	"github.com/shoet/blog/internal/usecase/create_user_profile"
	"github.com/shoet/blog/internal/usecase/get_user_profile"
	"github.com/shoet/blog/internal/usecase/update_user_profile"
)

type GetUserProfileHandler struct {
	jwtService JWTService
	usecase    *get_user_profile.Usecase
}

func NewGetUserProfileHandler(
	jwtService JWTService,
	usecase *get_user_profile.Usecase,
) *GetUserProfileHandler {
	return &GetUserProfileHandler{
		jwtService: jwtService,
		usecase:    usecase,
	}
}

/*
RequuestBody:

	QueryParameter:
		user_id: int

Response:

	userProfileId: int
	userId: int
	nickname: string
	avatarImageFileURL: string | null
	biography: string | null
	created: datetime
	updated: datetime
*/
func (h *GetUserProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	q := r.URL.Query()
	userId := q.Get("user_id")
	if userId == "" {
		logger.Error("user_id is required")
		response.RespondBadRequest(w, r, nil)
		return
	}
	userIdInt, err := strconv.Atoi(strings.TrimSpace(userId))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to convert id to int: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	userProfile, err := h.usecase.Run(ctx, models.UserId(userIdInt))
	if err != nil {
		if errors.Is(err, get_user_profile.ErrNotFound) {
			logger.Error(fmt.Sprintf("user profile not found: %v", err))
			response.RespondNotFound(w, r, err)
			return
		}
		logger.Error(fmt.Sprintf("failed to get user profile: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, userProfile); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type CreateUserProfileHandler struct {
	validator  *validator.Validate
	jwtService JWTService
	usecase    *create_user_profile.Usecase
}

func NewCreateUserProfileHandler(
	validator *validator.Validate,
	jwtService JWTService,
	usecase *create_user_profile.Usecase,
) *CreateUserProfileHandler {
	return &CreateUserProfileHandler{
		validator:  validator,
		jwtService: jwtService,
		usecase:    usecase,
	}
}

/*
RequuestBody:

	userId: int
	nickname: string
	avatarImageUrl: string | null
	biography: string | null

Response:

	user_profile: UserProfile
*/
func (h *CreateUserProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		UserID         models.UserId `json:"userId" validate:"required"`
		Nickname       string        `json:"nickname" validate:"required"`
		AvatarImageURL *string       `json:"avatarImageUrl"`
		BioGraphy      *string       `json:"biography"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := h.validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	requestUserId, err := session.GetUserId(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user id from context: %v", err))
		response.RespondUnauthorized(w, r, err)
		return
	}
	if reqBody.UserID != requestUserId {
		logger.Error(fmt.Sprintf("user id in request body is not equal to user id in context: %v", err))
		response.RespondUnauthorized(w, r, err)
		return
	}

	input := create_user_profile.CreateUserProfileInput{
		UserId:         reqBody.UserID,
		Nickname:       reqBody.Nickname,
		AvatarImageURL: reqBody.AvatarImageURL,
		BioGraphy:      reqBody.BioGraphy,
	}
	userProfile, err := h.usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create user profile: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, userProfile); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type UpdateUserProfileHandler struct {
	validator  *validator.Validate
	jwtService JWTService
	usecase    *update_user_profile.Usecase
}

func NewUpdateUserProfileHandler(
	validator *validator.Validate,
	jwtService JWTService,
	usecase *update_user_profile.Usecase,
) *UpdateUserProfileHandler {
	return &UpdateUserProfileHandler{
		validator:  validator,
		jwtService: jwtService,
		usecase:    usecase,
	}
}

/*
RequuestBody:

	user_id: int
	nickname: string | null
	avatar_image_url: string | null
	biography: string | null

Response:

	user_profile: UserProfile
*/
func (h *UpdateUserProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		UserID         models.UserId `json:"userId" validate:"required"`
		Nickname       string        `json:"nickname" validate:"required"`
		AvatarImageURL *string       `json:"avatarImageUrl"`
		BioGraphy      *string       `json:"biography"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := h.validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	requestUserId, err := session.GetUserId(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get user id from context: %v", err))
		response.RespondUnauthorized(w, r, err)
		return
	}
	if reqBody.UserID != requestUserId {
		logger.Error(fmt.Sprintf("user id in request body is not equal to user id in context: %v", err))
		response.RespondUnauthorized(w, r, err)
		return
	}

	input := update_user_profile.UpdateUserProfileInput{
		UserId:         reqBody.UserID,
		Nickname:       reqBody.Nickname,
		AvatarImageURL: reqBody.AvatarImageURL,
		BioGraphy:      reqBody.BioGraphy,
	}

	userProfile, err := h.usecase.Run(ctx, input)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to update user profile: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}

	if userProfile == nil {
		logger.Error("update not occurred")
		response.RespondNoContent(w, r)
		return
	}

	if err := response.RespondJSON(w, r, http.StatusOK, userProfile); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
