package usecases

import (
	"context"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
)

type WorkshopUsecase interface {
	GetWorkshop(ctx context.Context, id int64, fields []string) (*models.Workshop, error)
	ListWorkshop(ctx context.Context, filter models.WorkshopFilter) ([]*models.Workshop, error)
}

type workshopUsecaseImpl struct {
	repo repositories.WorkshopRepo
}

func NewWorkshopUsecase(repo repositories.WorkshopRepo) WorkshopUsecase {
	return &workshopUsecaseImpl{
		repo: repo,
	}
}
func (u *workshopUsecaseImpl) GetWorkshop(ctx context.Context, id int64, fields []string) (*models.Workshop, error) {
	return u.repo.GetWorkshopById(ctx, id, fields)
}

func (u *workshopUsecaseImpl) ListWorkshop(ctx context.Context, filter models.WorkshopFilter) ([]*models.Workshop, error) {
	return u.repo.ListWorkshop(ctx, filter)
}
