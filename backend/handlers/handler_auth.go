package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/logging"
)

type AuthLoginHandler struct {
	Service   AuthManager
	Validator *validator.Validate
	Cookie    *CookieManager
}

func NewAuthLoginHandler(
	service AuthManager,
	validator *validator.Validate,
	cookie *CookieManager,
) *AuthLoginHandler {
	return &AuthLoginHandler{
		Service:   service,
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
	if err := JsonToStruct(r, &reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to parse request body: %v", err))
		ResponsdBadRequest(w, r, err)
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		logger.Error(fmt.Sprintf("failed to validate request body: %v", err))
		ResponsdBadRequest(w, r, err)
		return
	}

	token, err := a.Service.Login(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		logger.Error(fmt.Sprintf("failed login: %v", err))
		RespondUnauthorized(w, r, err)
		return
	}
	resp := struct {
		AuthToken string `json:"authToken"`
	}{
		AuthToken: token,
	}
	if err := a.Cookie.SetCookie(w, "authToken", resp.AuthToken); err != nil {
		logger.Error(fmt.Sprintf("failed to set cookie: %v", err))
		ResponsdInternalServerError(w, r, err)
		return
	}

	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
	return
}

type AuthSessionLoginHandler struct {
	Service AuthManager
}

func NewAuthSessionLoginHandler(
	service AuthManager,
) *AuthSessionLoginHandler {
	return &AuthSessionLoginHandler{
		Service: service,
	}
}

func (a *AuthSessionLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.GetLogger(ctx)
	token := r.Header.Get("Authorization")
	if token == "" {
		logger.Error(fmt.Sprintf("failed to get authorization header"))
		RespondUnauthorized(w, r, fmt.Errorf("failed to get authorization header"))
		return
	}

	if !strings.HasPrefix(token, "Bearer ") {
		logger.Error(fmt.Sprintf("failed to get authorization header"))
		RespondUnauthorized(w, r, fmt.Errorf("failed to get authorization header"))
		return
	}

	token = strings.TrimPrefix(token, "Bearer ")

	u, err := a.Service.LoginSession(ctx, token)
	if err != nil {
		logger.Error(fmt.Sprintf("failed login: %v", err))
		RespondUnauthorized(w, r, err)
		return
	}
	if err := RespondJSON(w, r, http.StatusOK, u); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
	return
}

type AuthLogoutHandler struct {
	Cookie *CookieManager
}

func NewAuthLogoutHandler(
	cookie *CookieManager,
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
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error(fmt.Sprintf("failed to respond json response: %v", err))
	}
	return
}
