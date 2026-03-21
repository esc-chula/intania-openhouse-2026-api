package usecases

import (
	"context"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
)

type WorkshopUsecase interface {
	GetWorkshop(ctx context.Context, userEmail string, workshopId int64, fields []string) (*models.WorkshopDetail, error)
	ListWorkshop(ctx context.Context, filter models.WorkshopFilter) ([]*models.Workshop, error)
}

type workshopUsecaseImpl struct {
	workshopRepo repositories.WorkshopRepo
	userRepo     repositories.UserRepo
}

func NewWorkshopUsecase(workshopRepo repositories.WorkshopRepo, userRepo repositories.UserRepo) WorkshopUsecase {
	return &workshopUsecaseImpl{
		workshopRepo: workshopRepo,
		userRepo:     userRepo,
	}
}

func (u *workshopUsecaseImpl) GetWorkshop(ctx context.Context, userEmail string, workshopId int64, fields []string) (*models.WorkshopDetail, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, userEmail, []string{"id"})
	if err != nil {
		return nil, err
	}
	userId := user.ID

	return u.workshopRepo.GetWorkshopDetail(ctx, userId, workshopId, fields)
}

func (u *workshopUsecaseImpl) ListWorkshop(ctx context.Context, filter models.WorkshopFilter) ([]*models.Workshop, error) {
	return u.workshopRepo.ListWorkshop(ctx, filter)
}
