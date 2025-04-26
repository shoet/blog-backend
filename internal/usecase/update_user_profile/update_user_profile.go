package update_user_profile

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastracture"
	"github.com/shoet/blog/internal/infrastracture/models"
)

type FileRepository interface {
	ExistsFile(ctx context.Context, file *models.File) (bool, error)
}

type UserProfileRepository interface {
	Update(
		ctx context.Context,
		tx infrastracture.TX,
		userId models.UserId, nickname string, avatarImageFileName *string, bioGraphy *string,
	) (*models.UserProfile, error)
}

type Usecase struct {
	config                *config.Config
	DB                    infrastracture.DB
	FileRepository        FileRepository
	UserProfileRepository UserProfileRepository
}

func NewUsecase(
	config *config.Config,
	db infrastracture.DB,
	fileRepository FileRepository,
	userProfileRepository UserProfileRepository,
) *Usecase {
	return &Usecase{
		config:                config,
		DB:                    db,
		FileRepository:        fileRepository,
		UserProfileRepository: userProfileRepository,
	}
}

type UpdateUserProfileInput struct {
	UserId         models.UserId
	Nickname       string
	AvatarImageURL *string
	BioGraphy      *string
}

const MAX_NICKNAME_LENGTH = 30

func (u *Usecase) Run(ctx context.Context, input UpdateUserProfileInput) (*models.UserProfile, error) {
	if len(input.Nickname) >= MAX_NICKNAME_LENGTH {
		return nil, fmt.Errorf("nickname is too long: %s", input.Nickname)
	}
	var avatarImageFileName *string
	if input.AvatarImageURL != nil {
		file, err := models.NewFileFromURL(u.config, *input.AvatarImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create file from url: %w", err)
		}
		if exists, err := u.FileRepository.ExistsFile(ctx, file); err != nil {
			return nil, fmt.Errorf("failed to check exist file: %w", err)
		} else if !exists {
			return nil, fmt.Errorf("file not found: %v", file)
		}
		avatarImageFileName = &file.FileName
	}

	userProfile, err := u.UserProfileRepository.Update(
		ctx, u.DB, input.UserId, input.Nickname, avatarImageFileName, input.BioGraphy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return userProfile, nil
}
