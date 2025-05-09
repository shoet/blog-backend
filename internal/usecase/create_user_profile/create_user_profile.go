package create_user_profile

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/config"
	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type FileRepository interface {
	ExistsFile(ctx context.Context, file *models.File) (bool, error)
}

type UserProfileRepository interface {
	Create(
		ctx context.Context,
		tx infrastructure.TX,
		userId models.UserId, nickname string, avatarImageFileName *string, bioGraphy *string,
	) (*models.UserProfile, error)
}

type Usecase struct {
	Config                *config.Config
	DB                    infrastructure.DB
	FileRepository        FileRepository
	UserProfileRepository UserProfileRepository
}

func NewUsecase(
	config *config.Config,
	db infrastructure.DB,
	fileRepository FileRepository,
	userProfileRepository UserProfileRepository,
) *Usecase {
	return &Usecase{
		Config:                config,
		DB:                    db,
		FileRepository:        fileRepository,
		UserProfileRepository: userProfileRepository,
	}
}

const MAX_NICKNAME_LENGTH = 30

type CreateUserProfileInput struct {
	UserId         models.UserId
	Nickname       string
	AvatarImageURL *string
	BioGraphy      *string
}

func (u *Usecase) Run(ctx context.Context, input CreateUserProfileInput) (*models.UserProfile, error) {
	if len(input.Nickname) >= MAX_NICKNAME_LENGTH {
		return nil, fmt.Errorf("nickname is too long: %s", input.Nickname)
	}
	var avatarImageFileName *string
	if input.AvatarImageURL != nil {
		file, err := models.NewFileFromURL(u.Config, *input.AvatarImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create file: %w", err)
		}
		if exists, err := u.FileRepository.ExistsFile(ctx, file); err != nil {
			return nil, fmt.Errorf("failed to check exist file: %w", err)
		} else if !exists {
			return nil, fmt.Errorf("file not found: %v", file)
		}
		avatarImageFileName = &file.FileName
	}

	userProfile, err := u.UserProfileRepository.Create(
		ctx, u.DB, input.UserId, input.Nickname, avatarImageFileName, input.BioGraphy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user profile: %w", err)
	}

	return userProfile, nil
}
