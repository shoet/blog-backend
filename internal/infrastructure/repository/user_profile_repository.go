package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type UserProfileRepository struct {
	config *config.Config
}

func NewUserProfileRepository(config *config.Config) *UserProfileRepository {
	return &UserProfileRepository{
		config: config,
	}
}

/*
Get は、userIdに一致するユーザープロフィールを取得する。

レコードが存在しない場合は、nilを返す。
*/
func (r *UserProfileRepository) Get(
	ctx context.Context,
	tx infrastructure.TX,
	userId models.UserId,
) (*models.UserProfile, error) {

	builder := goqu.
		Select("id", "user_id", "nickname", "avatar_image_file_name", "bio", "created", "modified").
		From("user_profile").
		Where(goqu.Ex{"user_id": userId})

	query, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := tx.QueryRowxContext(ctx, query, params...)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to execute query: %w", row.Err())
	}

	var userProfile models.UserProfile
	if err := row.StructScan(&userProfile); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan struct: %w", err)
	}

	if userProfile.AvatarImageFileName != nil {
		file, err := models.NewFile("avatar_image", *userProfile.AvatarImageFileName)
		if err != nil {
			return nil, fmt.Errorf("failed to get file: %w", err)
		}
		avatarImageFileURL, err := file.GetFileURL(r.config)
		if err != nil {
			return nil, fmt.Errorf("failed to get file url: %w", err)
		}
		userProfile.AvatarImageFileURL = &avatarImageFileURL
	}
	return &userProfile, nil
}

func (r *UserProfileRepository) Create(
	ctx context.Context,
	tx infrastructure.TX,
	userId models.UserId, nickname string, avatarImageFileName *string, bioGraphy *string,
) (*models.UserProfile, error) {
	builder := goqu.
		Insert("user_profile").
		Rows(
			goqu.Record{
				"user_id":                userId,
				"nickname":               nickname,
				"avatar_image_file_name": avatarImageFileName,
				"bio":                    bioGraphy,
			}).
		Returning("*")
	query, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var userProfile models.UserProfile
	row := tx.QueryRowxContext(ctx, query, params...)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to insert: %w", row.Err())
	}
	if err := row.StructScan(&userProfile); err != nil {
		return nil, fmt.Errorf("failed to scan: %v", err)
	}
	return &userProfile, nil
}

func (r *UserProfileRepository) Update(
	ctx context.Context,
	tx infrastructure.TX,
	userId models.UserId, nickname string, avatarImageFileName *string, bioGraphy *string,
) (*models.UserProfile, error) {

	builder := goqu.
		Update("user_profile").
		Set(goqu.Record{
			"nickname":               nickname,
			"avatar_image_file_name": avatarImageFileName,
			"bio":                    bioGraphy,
		}).
		Where(goqu.Ex{"user_id": userId}).
		Returning("*")

	query, params, err := builder.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := tx.QueryRowxContext(ctx, query, params...)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	var userProfile models.UserProfile
	if err := row.StructScan(&userProfile); err != nil {
		return nil, fmt.Errorf("failed to scan struct: %w", err)
	}
	return &userProfile, err
}
