package login_user

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
}

type Usecase struct {
	authService AuthService
}

func NewUsecase(authService AuthService) *Usecase {
	return &Usecase{
		authService: authService,
	}
}

func (a *Usecase) Run(ctx context.Context, email string, password string) (string, error) {
	token, err := a.authService.Login(ctx, email, password)
	if err != nil {
		return "", err
	}
	return token, nil
}
