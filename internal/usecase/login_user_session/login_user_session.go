package login_user_session

import (
	"context"

	"github.com/shoet/blog/internal/infrastructure/models"
)

type AuthService interface {
	LoginSession(ctx context.Context, token string) (*models.User, error)
}

type Usecase struct {
	authService AuthService
}

func NewUsecase(authService AuthService) *Usecase {
	return &Usecase{
		authService: authService,
	}
}

func (u *Usecase) Run(ctx context.Context, token string) (*models.User, error) {
	return u.authService.LoginSession(ctx, token)
}
