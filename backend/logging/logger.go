package logging

import (
	"context"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

const LoggerKey = "log"

type Logger = zerolog.Logger

func NewLogger() *zerolog.Logger {
	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger()
	return &logger
}

func WithLoggerMiddleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), LoggerKey, logger)
			logger.Info().Msgf("request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogger(ctx context.Context) zerolog.Logger {
	return ctx.Value(LoggerKey).(zerolog.Logger)
}
