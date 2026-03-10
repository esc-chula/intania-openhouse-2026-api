package usecases

import (
	"context"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
)

// TODO:
type UserUsecase interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, email string, fields []string) (*models.User, error)
}

type userUsecaseImpl struct {
	repo          repositories.UserRepo
	stampRepo     repositories.StampRepo
	transactioner baserepo.Transactioner
}

func NewUserUsecase(repo repositories.UserRepo, stampRepo repositories.StampRepo, transactioner baserepo.Transactioner) UserUsecase {
	return &userUsecaseImpl{
		repo:          repo,
		stampRepo:     stampRepo,
		transactioner: transactioner,
	}
}

func (u *userUsecaseImpl) CreateUser(ctx context.Context, user *models.User) error {
	return u.transactioner.Transaction(ctx, func(ctx context.Context) error {
		if err := u.repo.CreateUser(ctx, user); err != nil {
			return err
		}

		stampPosters := []models.StampPoster{
			{UserID: user.ID, Type: models.StampTypeDepartment},
			{UserID: user.ID, Type: models.StampTypeClub},
			{UserID: user.ID, Type: models.StampTypeExhibition},
		}

		return u.stampRepo.CreateStampPosters(ctx, stampPosters)
	})
}

func (u *userUsecaseImpl) GetUser(ctx context.Context, email string, fields []string) (*models.User, error) {
	return u.repo.GetUserByEmail(ctx, email, fields)
}
