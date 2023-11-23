package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/shoet/blog/config"
	"github.com/shoet/blog/services"
	"golang.org/x/net/context"
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

const loggerKey = "log"

func WithLoggerMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loggerKey, logger)
			logger.Info().Msgf("request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogger(ctx context.Context) zerolog.Logger {
	return ctx.Value(loggerKey).(zerolog.Logger)
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
		logger := GetLogger(ctx)
		token := r.Header.Get("Authorization")
		if token == "" {
			logger.Error().Msgf("failed to get authorization header")
			RespondUnauthorized(w, r, fmt.Errorf("failed to get authorization header"))
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			logger.Error().Msgf("failed invalid authorization header")
			RespondUnauthorized(w, r, fmt.Errorf("failed to get authorization header"))
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		userId, err := a.jwter.VerifyToken(ctx, token)
		if err != nil {
			logger.Error().Msgf("failed to verify token: %v", err)
			RespondUnauthorized(w, r, fmt.Errorf("failed to verify token"))
			return
		}

		// Todo: set user info with context
		_ = userId

		next.ServeHTTP(w, r)
	})
}
