package testutil

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/shoet/blog/logging"
)

func SetLoggerContextToRequest(t *testing.T, request *http.Request) *http.Request {
	t.Helper()
	logger := logging.NewLogger(io.Discard) // テスト時はロガーによるログ出力を捨てる
	return request.WithContext(context.WithValue(request.Context(), logging.LoggerKey, logger))
}
