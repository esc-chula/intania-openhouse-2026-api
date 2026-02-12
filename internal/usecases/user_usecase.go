package usecases

import (
	"context"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
)

// TODO:
type UserUsecase interface {
	CreateUser(ctx context.Context, user *models.User) error
}

type userUsecaseImpl struct {
	repo repositories.UserRepo
}

func NewUserUsecase(repo repositories.UserRepo) UserUsecase {
	return &userUsecaseImpl{
		repo: repo,
	}
}
func (u *userUsecaseImpl) CreateUser(ctx context.Context, user *models.User) error {
	return u.repo.CreateUser(ctx, user)
}
