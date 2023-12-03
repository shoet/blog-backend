package testutil

import (
	"context"
	"net/http"
	"testing"

	"github.com/shoet/blog/logging"
)

var loggerKey = logging.LoggerKey

func SetLoggerContextToRequest(t *testing.T, request *http.Request) *http.Request {
	t.Helper()
	logger := logging.NewLogger()
	return request.WithContext(context.WithValue(request.Context(), loggerKey, *logger))
}
