package delete_privacy_policy

import (
	"context"
	"fmt"

	"github.com/shoet/blog/internal/infrastructure"
)

type PrivacyPolicyRepository interface {
	Delete(ctx context.Context, tx infrastructure.TX, id string) error
}

type Usecase struct {
	DB   infrastructure.DB
	repo PrivacyPolicyRepository
}

func NewUsecase(DB infrastructure.DB, repository PrivacyPolicyRepository) *Usecase {
	return &Usecase{
		DB:   DB,
		repo: repository,
	}
}

func (u *Usecase) Run(ctx context.Context, id string) error {
	if err := u.repo.Delete(ctx, u.DB, id); err != nil {
		return fmt.Errorf("failed to delete privacy policy")
	}
	return nil
}
