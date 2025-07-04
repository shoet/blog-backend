package put_privacy_policy

import (
	"context"
	"errors"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
	"github.com/shoet/blog/internal/infrastructure/repository"
)

type PrivacyPolicyRepository interface {
	Get(ctx context.Context, tx infrastructure.TX, id string) (*models.PrivacyPolicy, error)
	Create(ctx context.Context, tx infrastructure.TX, id string, content string) error
	UpdateContent(ctx context.Context, tx infrastructure.TX, id string, content string) error
}

type Usecase struct {
	DB                      infrastructure.DB
	PrivacyPolicyRepository PrivacyPolicyRepository
}

func NewUsecase(
	db infrastructure.DB,
	repo PrivacyPolicyRepository) *Usecase {
	return &Usecase{
		DB:                      db,
		PrivacyPolicyRepository: repo,
	}
}

func (u *Usecase) Run(ctx context.Context, id string, content string) error {
	p, err := u.PrivacyPolicyRepository.Get(ctx, u.DB, id)
	if err != nil {
		if !errors.Is(err, repository.ErrResourceNotFound) {
			return fmt.Errorf("failed to get privacy policy: %w", err)
		}
	}
	if p == nil {
		if err := u.PrivacyPolicyRepository.Create(ctx, u.DB, id, content); err != nil {
			return fmt.Errorf("failed to create privacy policy: %w", err)
		}
	} else {
		if err := u.PrivacyPolicyRepository.UpdateContent(ctx, u.DB, id, string(content)); err != nil {
			return fmt.Errorf("failed to update privacy policy content: %w", err)
		}
	}
	return nil
}
