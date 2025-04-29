package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/usecase/login_user"
	"github.com/shoet/blog/internal/usecase/login_user_session"
)

type AuthLoginHandler struct {
	Usecase   *login_user.Usecase
	Validator *validator.Validate
	Cookie    Cookier
}

func NewAuthLoginHandler(
	usecase *login_user.Usecase,
	validator *validator.Validate,
	cookie Cookier,
) *AuthLoginHandler {
	return &AuthLoginHandler{
		Usecase:   usecase,
		Validator: validator,
		Cookie:    cookie,
	}
}

func (a *AuthLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	var reqBody struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	defer r.Body.Close()
	if err := response.JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		response.RespondBadRequest(w, r, err)
		return
	}

	token, err := a.Usecase.Run(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		logger.Error(fmt.Sprintf("failed login: %v", err))
		response.RespondUnauthorized(w, r, err)
		return
	}
	resp := struct {
		AuthToken string `json:"authToken"`
	}{
		AuthToken: token,
	}
	if err := a.Cookie.SetCookie(w, "authToken", resp.AuthToken); err != nil {
		logger.Error(fmt.Sprintf("failed to set cookie: %v", err))
		response.RespondInternalServerError(w, r, err)
		return
	}

	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type AuthSessionLoginHandler struct {
	Usecase *login_user_session.Usecase
}

func NewAuthSessionLoginHandler(
	usecase *login_user_session.Usecase,
) *AuthSessionLoginHandler {
	return &AuthSessionLoginHandler{
		Usecase: usecase,
	}
}

/*
RequuestBody:

	RequestHeader:
		Authorization: string

Response:

	id int
	name string
	email string
	password string
	profile: UserProfile | null
	created: timestamp
	updated: timestamp
*/
func (a *AuthSessionLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	token := r.Header.Get("Authorization")
	if token == "" {
		msg := "failed to get authorization header"
		logger.Error(msg)
		response.RespondUnauthorized(w, r, errors.New(msg))
		return
	}

	if !strings.HasPrefix(token, "Bearer ") {
		msg := "passing invalid authorization header format"
		logger.Error(msg)
		response.RespondUnauthorized(w, r, errors.New(msg))
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	u, err := a.Usecase.Run(ctx, token)
	if err != nil {
		logger.Error(fmt.Sprintf("failed login: %v", err))
		response.RespondUnauthorized(w, r, err)
		return
	}
	if err := response.RespondJSON(w, r, http.StatusOK, u); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}

type AuthLogoutHandler struct {
	Cookie Cookier
}

func NewAuthLogoutHandler(
	cookie Cookier,
) *AuthLogoutHandler {
	return &AuthLogoutHandler{
		Cookie: cookie,
	}
}

func (a *AuthLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	a.Cookie.ClearCookie(w, "authToken")
	resp := struct {
		Message string `json:"message"`
	}{
		Message: "success",
	}
	if err := response.RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
}
