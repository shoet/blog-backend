package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/shoet/blog/config"
	"github.com/shoet/blog/logging"
	"github.com/shoet/blog/services"
)

func NewCORSMiddleWare(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		whiteList := getCORSWhiteList(cfg)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if originAllowed(origin, whiteList) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,UPDATE,OPTIONS")
			w.Header().Set("Content-Type", "application/json")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getCORSWhiteList(cfg *config.Config) []string {
	return strings.Split(cfg.CORSWhiteList, ",")
}

func originAllowed(origin string, whiteList []string) bool {
	for _, allowedOrigin := range whiteList {
		if allowedOrigin == origin {
			return true
		}
	}
	return false
}

type AuthorizationMiddleware struct {
	jwter services.JWTer
}

func NewAuthorizationMiddleware(jwter services.JWTer) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		jwter: jwter,
	}
}

func (a *AuthorizationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.GetLogger(ctx)
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error(fmt.Sprintf("failed to get authorization header"))
			RespondUnauthorized(w, r, fmt.Errorf("failed to get authorization header"))
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			logger.Error(fmt.Sprintf("failed invalid authorization header"))
			RespondUnauthorized(w, r, fmt.Errorf("failed to get authorization header"))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		userId, err := a.jwter.VerifyToken(ctx, token)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to verify token: %v", err))
			RespondUnauthorized(w, r, fmt.Errorf("failed to verify token"))
			return
		}

		// Todo: set user info with context
		_ = userId

		next.ServeHTTP(w, r)
	})
}
