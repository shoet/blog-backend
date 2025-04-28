package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/infrastructure/repository"
	"github.com/shoet/blog/internal/testutil"
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
				}},
			want: &models.User{
				Name: "test",
			},
			wantErr: nil,
		},
	}

	ctx := context.Background()

	sut, err := repository.NewUserRepository(clocker)
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}

	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := db.Beginx()

			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			defer tx.Rollback()

			sql := `
			INSERT INTO users
				(name, email, password)
			VALUES
				($1, $2, $3)
			RETURNING id
			;
			`
			row := tx.QueryRowxContext(
				ctx,
				sql,
				tt.args.user.Name,
				tt.args.user.Email,
				tt.args.user.Password,
			)
			if row.Err() != nil {
				t.Fatalf("failed to insert user: %v", row.Err())
			}
			var userId models.UserId
			if err := row.Scan(&userId); err != nil {
				t.Fatalf("failed to scan user: %v", err)
			}

			got, err := sut.Get(ctx, tx, userId)
			if err != nil {
				t.Fatalf("failed to get user: %v", err)
			}

			opt := cmpopts.IgnoreFields(models.User{}, "Id", "Created", "Modified")
			if diff := cmp.Diff(got, tt.want, opt); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}

func Test_UserRepository_GetByEmail(t *testing.T) {
	type args struct {
		prepareUser []*models.User
		email       string
	}

	type want struct {
		user  *models.User
		error error
	}

	clocker := &clocker.FiexedClocker{}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr error
	}{
		{
			name: "success",
			args: args{
				prepareUser: []*models.User{
					{
						Name:     "test",
						Email:    "test@test.com",
						Password: "test",
					},
				},
				email: "test@test.com",
			},
			want: want{
				user: &models.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "test",
				},
				error: nil,
			},
		},
		{
			name: "failed not found",
			args: args{
				prepareUser: []*models.User{},
				email:       "test@test.com",
			},
			want: want{
				user:  nil,
				error: repository.ErrUserNotFound,
			},
		},
	}

	ctx := context.Background()

	sut, err := repository.NewUserRepository(clocker)
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}

	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := db.Beginx()

			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			defer tx.Rollback()

			for _, u := range tt.args.prepareUser {
				sql := `
				INSERT INTO users
					(name, email, password)
				VALUES
					($1, $2, $3)
				;
				`
				_, err = tx.ExecContext(
					ctx,
					sql,
					u.Name,
					u.Email,
					u.Password,
				)
				if err != nil {
					t.Fatalf("failed to insert user: %v", err)
				}
			}

			got, err := sut.GetByEmail(ctx, tx, tt.args.email)
			if err != nil {
				if !errors.Is(err, tt.want.error) {
					t.Fatalf("failed GetByEmail")
				}
			}

			opt := cmpopts.IgnoreFields(models.User{}, "Id", "Created", "Modified")
			if diff := cmp.Diff(got, tt.want.user, opt); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}

func Test_UserRepository_Add(t *testing.T) {
	type args struct {
		user *models.User
	}

	type want struct {
		user *models.User
	}

	clocker := &clocker.FiexedClocker{}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success",
			args: args{
				user: &models.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "test",
				},
			},
			want: want{
				user: &models.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "test",
				},
			},
		},
	}

	ctx := context.Background()

	sut, err := repository.NewUserRepository(clocker)
	if err != nil {
		t.Fatalf("failed to create user repository: %v", err)
	}

	db, err := testutil.NewDBPostgreSQLForTest(t, ctx)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	testutil.RepositoryTestPrepare(t, ctx, db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := db.Beginx()

			if err != nil {
				t.Fatalf("failed to create tx: %v", err)
			}
			defer tx.Rollback()

			user, err := sut.Add(ctx, tx, tt.args.user)
			if err != nil {
				t.Fatalf("failed to add user: %v", err)
			}

			selectSql := "SELECT * FROM users WHERE id = $1"
			var gotUser models.User
			if err := tx.QueryRowxContext(ctx, selectSql, user.Id).StructScan(&gotUser); err != nil {
				t.Fatalf("failed to get user: %v", err)
			}

			opt := cmpopts.IgnoreFields(models.User{}, "Id", "Created", "Modified")
			if diff := cmp.Diff(&gotUser, tt.want.user, opt); diff != "" {
				t.Errorf("(-got +want)\n%s", diff)
			}
		})
	}
}
