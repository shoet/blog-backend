package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertResponse(t *testing.T, got *http.Response, status int, want []byte) error {
	t.Helper()
	t.Cleanup(func() {
		got.Body.Close()
	})

	gb, err := io.ReadAll(got.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %v", err)
	}

	if got.StatusCode != status {
		return fmt.Errorf("got status %d, want %d", got.StatusCode, status)
	}

	if len(gb) == 0 && len(want) == 0 {
		return fmt.Errorf("got empty response body")
	}

	return AssertJSON(t, want, gb)
}

func AssertJSON(t *testing.T, want, got []byte) error {
	t.Helper()

	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		return fmt.Errorf("cannot unmarshal want %q: %v", want, err)
	}

	if err := json.Unmarshal(got, &jg); err != nil {
		return fmt.Errorf("cannot unmarshal got %q: %v", got, err)
	}

	if diff := cmp.Diff(jg, jw); diff != "" {
		return fmt.Errorf("(-got +want)\n%s", diff)
	}
	return nil
}
