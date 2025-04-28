package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/interfaces/response"
	"github.com/shoet/blog/internal/logging"
	"github.com/shoet/blog/internal/session"
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

type JWTService interface {
	VerifyToken(ctx context.Context, token string) (models.UserId, error)
}

type AuthorizationMiddleware struct {
	jwter JWTService
}

func NewAuthorizationMiddleware(jwter JWTService) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		jwter: jwter,
	}
}

func (a *AuthorizationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logging.GetLogger(ctx)

		token, err := a.ChallengeAuthorizationHeader(r.Header)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get authorization header: %v", err))
			response.RespondUnauthorized(w, r, fmt.Errorf("failed to get authorization header"))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		userId, err := a.jwter.VerifyToken(ctx, token)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to verify token: %v", err))
			response.RespondUnauthorized(w, r, fmt.Errorf("failed to verify token"))
			return
		}

		// set UserId to context
		ctx = session.SetUserId(ctx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ChallengeAuthorizationHeader は、Authorizationヘッダが大文字・小文字のどちらであっても認証トークンを受け取れるようにする
func (a *AuthorizationMiddleware) ChallengeAuthorizationHeader(h http.Header) (string, error) {
	authorizationHeader := []string{"Authorization", "authorization"}
	for _, a := range authorizationHeader {
		token := h.Get(a)
		if token == "" {
			return "", fmt.Errorf("failed to get authorization header")
		}

		if !strings.HasPrefix(token, "Bearer ") {
			return "", fmt.Errorf("failed invalid authorization header")
		}
		return token, nil
	}
	return "", fmt.Errorf("failed to get authorization header is not found")
}
