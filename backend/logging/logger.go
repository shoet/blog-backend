package logging

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

const LoggerKey = "log"

func WithLoggerMiddleware(logger *Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), LoggerKey, logger)
			logger.Info(fmt.Sprintf("request: %s %s", r.Method, r.URL.Path))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetLogger(ctx context.Context) *Logger {
	l := ctx.Value(LoggerKey)
	if l == nil {
		fmt.Println("logger is not set in context")
	}
	return ctx.Value(LoggerKey).(*Logger)
}

type Logger struct {
	logger *zerolog.Logger
}

func NewLogger(w io.Writer) *Logger {
	logger := zerolog.
		New(w).
		With().
		Timestamp().
		Logger()
	return &Logger{
		logger: &logger,
	}
}

func (l *Logger) Info(message string) {
	l.logger.Info().Msg(message)
}

func (l *Logger) Error(message string) {
	l.logger.Error().Msg(message)
}
