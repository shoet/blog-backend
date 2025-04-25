package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shoet/blog/internal/infrastracture/models"
	"github.com/shoet/blog/internal/infrastracture/repository"
	"github.com/shoet/blog/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func Test_UserProfileRepository_Create(t *testing.T) {
	ptrStr := func(s string) *string {
		return &s
	}

	type args struct {
		user                models.User
		userId              models.UserId
		nickname            string
		avatarImageFileName *string
		bioGraphy           *string
	}
	type wants struct {
		err         error
		userProfile *models.UserProfile
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "normal case",
			args: args{
				user: models.User{
					Id:       models.UserId(1),
					Name:     "test",
					Email:    "test@example.com",
					Password: "password",
				},
				userId:              models.UserId(1),
				nickname:            "nickname",
				avatarImageFileName: nil,
				bioGraphy:           nil,
			},
			wants: wants{
				err: nil,
				userProfile: &models.UserProfile{
					UserId:              models.UserId(1),
					Nickname:            "nickname",
					AvatarImageFileName: nil,
					Biography:           nil,
				},
			},
		},
		{
			name: "オプショナルなフィールドが入っている",
			args: args{
				user: models.User{
					Id:       models.UserId(1),
					Name:     "test",
					Email:    "test@example.com",
					Password: "password",
				},
				userId:              models.UserId(1),
				nickname:            "nickname",
				avatarImageFileName: ptrStr("avatar.png"),
				bioGraphy:           ptrStr("bio"),
			},
			wants: wants{
				err: nil,
				userProfile: &models.UserProfile{
					UserId:              models.UserId(1),
					Nickname:            "nickname",
					AvatarImageFileName: ptrStr("avatar.png"),
					Biography:           ptrStr("bio"),
				},
			},
		},
	}

	for _, tt := range tests {

		db, err := testutil.NewDBPostgreSQLForTest(t, context.Background())
		if err != nil {
			t.Fatalf("failed to create test db: %v", err)
		}

		sut := repository.NewUserProfileRepository()

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			tx, err := db.BeginTxx(ctx, nil)
			if err != nil {
				t.Fatalf("failed to begin: %v", err)
			}
			defer tx.Rollback()

			createUserBuilder := goqu.Insert("users").Rows(
				goqu.Record{
					"id":       tt.args.user.Id,
					"name":     tt.args.user.Name,
					"email":    tt.args.user.Email,
					"password": tt.args.user.Password,
				},
			)
			queryUser, params, err := createUserBuilder.ToSQL()
			if err != nil {
				t.Fatalf("failed to build query: %v", err)
			}
			if _, err := tx.ExecContext(ctx, queryUser, params...); err != nil {
				t.Fatalf("failed to insert user: %v", err)
			}

			userProfile, err := sut.Create(
				ctx, tx,
				tt.args.userId,
				tt.args.nickname,
				tt.args.avatarImageFileName,
				tt.args.bioGraphy,
			)

			if diff := cmp.Diff(
				tt.wants.err,
				err,
			); diff != "" {
				t.Errorf("error mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(
				tt.wants.userProfile,
				userProfile,
				cmpopts.IgnoreFields(models.UserProfile{}, "UserProfileId", "Created", "Modified"),
			); diff != "" {
				t.Errorf("userProfile mismatch (-want +got):\n%s", diff)
			}
		})
	}

}

func Test_UserProfileRepository_Update(t *testing.T) {
	ptrStr := func(s string) *string { return &s }
	type args struct {
		userProfile         models.UserProfile
		userId              models.UserId
		nickname            *string
		avatarImageFileName *string
		bioGraphy           *string
	}
	type wants struct {
		err         error
		userProfile *models.UserProfile
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "noupdate",
			args: args{
				userProfile: models.UserProfile{
					UserId:              models.UserId(1),
					Nickname:            "shoet",
					AvatarImageFileName: nil,
					Biography:           nil,
				},
				userId:              models.UserId(1),
				nickname:            nil,
				bioGraphy:           nil,
				avatarImageFileName: nil,
			},
			wants: wants{
				err:         errors.New("no update columns"),
				userProfile: nil,
			},
		},
		{
			name: "update",
			args: args{
				userProfile: models.UserProfile{
					UserId:              models.UserId(1),
					Nickname:            "shoet",
					AvatarImageFileName: nil,
					Biography:           nil,
				},
				userId:              models.UserId(1),
				nickname:            nil,
				bioGraphy:           ptrStr("bio"),
				avatarImageFileName: nil,
			},
			wants: wants{
				err: nil,
				userProfile: &models.UserProfile{
					UserId:              models.UserId(1),
					Nickname:            "shoet",
					AvatarImageFileName: nil,
					Biography:           ptrStr("bio"),
				},
			},
		},
	}

	db, err := testutil.NewDBPostgreSQLForTest(t, context.Background())
	if err != nil {
		t.Fatalf("failed to create test db: %v", err)
	}

	sut := repository.NewUserProfileRepository()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			tx, err := db.BeginTxx(ctx, nil)
			if err != nil {
				t.Fatalf("failed to create transaction: %v", err)
			}
			defer tx.Rollback()

			createQuery, params, err := goqu.
				Insert("users").
				Rows(goqu.Record{
					"id":       tt.args.userProfile.UserId,
					"name":     "test",
					"email":    "test@example.com",
					"password": "test",
				}).
				ToSQL()
			if err != nil {
				t.Fatalf("failed to create query: %v", err)
			}
			if _, err := tx.ExecContext(ctx, createQuery, params...); err != nil {
				t.Fatalf("failed create test user: %v", err)
			}

			createQuery, params, err = goqu.
				Insert("user_profile").
				Rows(goqu.Record{
					"user_id":                tt.args.userProfile.UserId,
					"nickname":               tt.args.userProfile.Nickname,
					"avatar_image_file_name": tt.args.userProfile.AvatarImageFileName,
					"bio":                    tt.args.userProfile.Biography,
				}).
				ToSQL()
			if err != nil {
				t.Fatalf("failed to create query: %v", err)
			}
			if _, err := tx.ExecContext(ctx, createQuery, params...); err != nil {
				t.Fatalf("failed create test user profile: %v", err)
			}

			userProfile, gotErr := sut.Update(
				ctx, tx,
				tt.args.userId, tt.args.nickname, tt.args.avatarImageFileName, tt.args.bioGraphy,
			)

			if err != nil {
				assert.ErrorContains(t, gotErr, tt.wants.err.Error())
			}

			if diff := cmp.Diff(
				tt.wants.userProfile,
				userProfile,
				cmpopts.IgnoreFields(models.UserProfile{}, "UserProfileId", "Created", "Modified"),
			); diff != "" {
				t.Errorf("xxx mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
