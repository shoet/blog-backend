package store

import (
	"context"
	"testing"

	"github.com/shoet/blog/testutil"
)

func Test_NewDBMySQL(t *testing.T) {
	ctx := context.Background()
	db, err := testutil.NewDBMySQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed new db: %v", err)
	}
	_ = db
}
