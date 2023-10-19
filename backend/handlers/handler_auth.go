package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/shoet/blog/config"
)

type AuthLoginHandler struct {
	Service   AuthManager
	Validator *validator.Validate
}

func (a *AuthLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
	var reqBody struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	defer r.Body.Close()
	if err := JsonToStruct(r, &reqBody); err != nil {
		logger.Error().Msgf("failed to parse request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		logger.Error().Msgf("failed to validate request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	token, err := a.Service.Login(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		logger.Error().Msgf("failed login: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}
	resp := struct {
		AuthToken string `json:"auth_token"`
	}{
		AuthToken: token,
	}
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}

type AuthAdminLoginHandler struct {
	Service   AuthManager
	Validator *validator.Validate
	config    *config.Config
}

func (a *AuthAdminLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := GetLogger(ctx)
	var reqBody struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	defer r.Body.Close()
	if err := JsonToStruct(r, &reqBody); err != nil {
		logger.Error().Msgf("failed to parse request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	if err := a.Validator.Struct(reqBody); err != nil {
		logger.Error().Msgf("failed to validate request body: %v", err)
		ResponsdBadRequest(w, r, err)
		return
	}

	token, err := a.Service.LoginAdmin(ctx, a.config, reqBody.Email, reqBody.Password)
	if err != nil {
		logger.Error().Msgf("failed login admin: %v", err)
		ResponsdInternalServerError(w, r, err)
		return
	}
	resp := struct {
		AuthToken string `json:"auth_token"`
	}{
		AuthToken: token,
	}
	if err := RespondJSON(w, r, http.StatusOK, resp); err != nil {
		logger.Error().Msgf("failed to respond json response: %v", err)
	}
	return
}

type AuthSessionLoginHandler struct {
	Service   AuthManager
	Validator *validator.Validate
}

func (a *AuthSessionLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get authorization header
	return
}
