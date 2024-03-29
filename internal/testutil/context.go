package testutil

import (
	"context"
	// "io"
	"net/http"
	"os"
	"testing"

	"github.com/shoet/blog/internal/logging"
)

func SetLoggerContextToRequest(t *testing.T, request *http.Request) *http.Request {
	t.Helper()
	// logger := logging.NewLogger(io.Discard) // テスト時はロガーによるログ出力を捨てる
	logger := logging.NewLogger(os.Stdout, "debug")
	return request.WithContext(context.WithValue(request.Context(), logging.LoggerKey{}, logger))
}
