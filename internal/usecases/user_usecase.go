package usecases

import (
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
)

// TODO:
type UserUsecase interface{}

type userUsecaseImpl struct {
	repo repositories.UserRepo
}

func NewUserUsecase(repo repositories.UserRepo) UserUsecase {
	return &userUsecaseImpl{
		repo: repo,
	}
}
