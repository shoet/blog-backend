package repository

import (
	"context"

	"github.com/shoet/blog/internal/clocker"
	"github.com/shoet/blog/internal/infrastructure"
	"github.com/shoet/blog/internal/infrastructure/models"
)

type PrivacyPolicyRepository struct {
	Clocker clocker.Clocker
}

func NewPrivacyPolicyRepository(clocker clocker.Clocker) *PrivacyPolicyRepository {
	return &PrivacyPolicyRepository{
		Clocker: clocker,
	}
}

func (r *PrivacyPolicyRepository) Get(ctx context.Context, tx infrastructure.TX, id string) (*models.PrivacyPolicy, error) {
	return nil, nil
}

func (r *PrivacyPolicyRepository) Create(ctx context.Context, tx infrastructure.TX, id string, content string) (*models.PrivacyPolicy, error) {
	return nil, nil
}

func (r *PrivacyPolicyRepository) UpdateContent(ctx context.Context, tx infrastructure.TX, id string, content string) (*models.PrivacyPolicy, error) {
	return nil, nil
}
