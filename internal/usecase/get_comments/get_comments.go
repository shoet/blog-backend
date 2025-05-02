package get_comments

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type CommmentRepository interface {
	GetByBlogId(
		ctx context.Context,
		tx infrastructure.TX,
		blogId models.BlogId,
		excludeDeleted bool,
	) ([]*models.Comment, error)
}

type UserProfileRepository interface {
	Get(ctx context.Context, tx infrastructure.TX, userId models.UserId) (*models.UserProfile, error)
}

type Usecase struct {
	DB                    infrastructure.DB
	commentRepository     CommmentRepository
	userProfileRepository UserProfileRepository
}

func NewUsecase(
	db infrastructure.DB,
	commentRepository CommmentRepository,
	userProfileRepository UserProfileRepository,
) *Usecase {
	return &Usecase{
		DB:                    db,
		commentRepository:     commentRepository,
		userProfileRepository: userProfileRepository,
	}
}

func (u *Usecase) Run(ctx context.Context, blogId models.BlogId) ([]*models.Comment, error) {
	comments, err := u.commentRepository.GetByBlogId(ctx, u.DB, blogId, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	for _, comment := range comments {
		if comment.UserId != nil {
			profile, err := u.userProfileRepository.Get(ctx, u.DB, *comment.UserId)
			if err != nil {
				return nil, fmt.Errorf("failed to get user profile: %w", err)
			}
			comment.AvatarImageFileURL = profile.AvatarImageFileURL
			comment.Nickname = &profile.Nickname
		}
	}
	return comments, nil
}
