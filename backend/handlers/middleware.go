package handlers

import (
	"net/http"

	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

var whiteList = []string{
	"http://localhost:5173",
	"http://localhost:3000",
	"http://localhost:6006",
}

func CORSMiddleWare(next http.Handler) http.Handler {
	// TODO: ブラッシュアップ
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if originAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
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

func originAllowed(origin string) bool {
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
