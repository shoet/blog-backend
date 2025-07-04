package repository

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
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
	query, params, err := goqu.
		Select("id", "content", "created", "modified").
		From("privacy_policy").
		Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	var privacy_policies []*models.PrivacyPolicy
	if err := tx.SelectContext(ctx, &privacy_policies, query, params...); err != nil {
		return nil, fmt.Errorf("failed to select privacy_policies: %w", err)
	}
	if len(privacy_policies) == 0 {
		return nil, ErrResourceNotFound
	}
	return privacy_policies[0], nil
}

func (r *PrivacyPolicyRepository) Create(ctx context.Context, tx infrastructure.TX, id string, content string) error {
	builder := goqu.
		Insert("privacy_policy").
		Cols("id", "content").
		Rows(
			goqu.Record{"id": id, "content": content},
		)
	query, params, err := builder.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	if _, err := tx.ExecContext(ctx, query, params...); err != nil {
		return fmt.Errorf("failed to insert privacy policy: %w", err)
	}
	return nil
}

func (r *PrivacyPolicyRepository) UpdateContent(ctx context.Context, tx infrastructure.TX, id string, content string) error {
	builder := goqu.
		Update("privacy_policy").
		Set(goqu.Record{"content": content}).
		Where(goqu.Ex{"id": id})
	query, params, err := builder.ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	if _, err := tx.ExecContext(ctx, query, params...); err != nil {
		return fmt.Errorf("failed to update privacy policy: %w", err)
	}
	return nil
}

func (r *PrivacyPolicyRepository) Delete(ctx context.Context, tx infrastructure.TX, id string) error {
	query, params, err := goqu.
		Delete("privacy_policy").
		Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	if _, err := tx.ExecContext(ctx, query, params...); err != nil {
		return fmt.Errorf("failed to delete privacy_policies: %w", err)
	}
	return nil
}
