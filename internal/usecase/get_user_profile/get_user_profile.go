package get_user_profile

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

var ErrNotFound = fmt.Errorf("user profile not found")

type UserProfileRepository interface {
	Get(
		ctx context.Context,
		tx infrastracture.TX,
		userId models.UserId,
	) (*models.UserProfile, error)
}

type Usecase struct {
	config                *config.Config
	DB                    infrastracture.DB
	UserProfileRepository UserProfileRepository
}

func NewUsecase(
	config *config.Config,
	db infrastracture.DB,
	userProfileRepository UserProfileRepository,
) *Usecase {
	return &Usecase{
		config:                config,
		DB:                    db,
		UserProfileRepository: userProfileRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, userId models.UserId) (*models.UserProfile, error) {

	userProfile, err := u.UserProfileRepository.Get(ctx, u.DB, userId)
	if err != nil {
		return nil, err
	}

	if userProfile == nil {
		return nil, ErrNotFound
	}

	if userProfile.AvatarImageFileName != nil {
		file, err := models.NewFile(models.FileType("avatar"), *userProfile.AvatarImageFileName)
		if err != nil {
			return nil, fmt.Errorf("invalid avatar image url: %w", err)
		}
		avatarImageURL, err := file.GetFileURL(u.config)
		if err != nil {
			return nil, fmt.Errorf("failed to get avatar image url: %w", err)
		}
		userProfile.AvatarImageFileURL = &avatarImageURL
	}
	return userProfile, nil
}
