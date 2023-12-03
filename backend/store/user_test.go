package store

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/clocker"
	"github.com/shoet/blog/models"
	"github.com/shoet/blog/testutil"
)

func Test_UserRepository_Get(t *testing.T) {
	type args struct {
		user *models.User
	}

	clocker := &clocker.FiexedClocker{}

	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "success",
			args: args{
				user: &models.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "test",
					Created:  clocker.Now(),
					Modified: clocker.Now(),
				}},
			want: &models.User{
				Name:     "test",
				Created:  clocker.Now(),
				Modified: clocker.Now(),
			},
			wantErr: nil,
		},
	}

	ctx := context.Background()

	sut, err := NewUserRepository(clocker)
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}

	db, err := testutil.NewDBMySQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	for _, tt := range tests {
		tx, err := db.Beginx()

		if err != nil {
			t.Fatalf("failed to create tx: %v", err)
		}
		defer tx.Rollback()

		sql := `
		INSERT INTO users
			(name, email, password, created, modified)
		VALUES
			(?, ?, ?, ?, ?)
		;
		`
		res, err := tx.ExecContext(
			ctx,
			sql,
			tt.args.user.Name,
			tt.args.user.Email,
			tt.args.user.Password,
			tt.args.user.Created,
			tt.args.user.Modified)
		gotId, err := res.LastInsertId()
		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		got, err := sut.Get(ctx, tx, models.UserId(gotId))
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		opt := cmpopts.IgnoreFields(models.User{}, "Id")
		if diff := cmp.Diff(got, tt.want, opt); diff != "" {
			t.Errorf("(-got +want)\n%s", diff)
		}
	}
}
