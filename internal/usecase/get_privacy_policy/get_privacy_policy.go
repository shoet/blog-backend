package get_privacy_policy

import (
	"context"
	"errors"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/infrastructure/repository"
)

var ErrResourceNotFound = fmt.Errorf("privacy policy not found")

type PrivacyPolicyRepository interface {
	Get(ctx context.Context, tx infrastructure.TX, id string) (*models.PrivacyPolicy, error)
}

type Usecase struct {
	DB             infrastructure.DB
	PrivacyPolicyRepository PrivacyPolicyRepository
}

func NewUsecase(
	db infrastructure.DB,
	repo PrivacyPolicyRepository) *Usecase {
	return &Usecase{
		DB: db,
		PrivacyPolicyRepository: repo,
	}
}

func (u *Usecase) Run(ctx context.Context, id string) (*models.PrivacyPolicy, error) {
	privacyPolicy, err := u.PrivacyPolicyRepository.Get(ctx, u.DB, id)
	if err != nil {
		if errors.Is(err, repository.ErrResourceNotFound) {
			return nil, ErrResourceNotFound
		}
		return nil, fmt.Errorf("failed to get privacy policy: %w", err)
	}
	return privacyPolicy, nil
}
