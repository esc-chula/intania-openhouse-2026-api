package usecases

import (
	"context"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
)

type ActivityUsecase interface {
	GetActivity(ctx context.Context, id int64) (*models.Activity, error)
	ListActivities(ctx context.Context, filter models.ActivityFilter) ([]*models.Activity, error)
}

type activityUsecaseImpl struct {
	repo repositories.ActivityRepo
}

func NewActivityUsecase(repo repositories.ActivityRepo) ActivityUsecase {
	return &activityUsecaseImpl{
		repo: repo,
	}
}

func (u *activityUsecaseImpl) GetActivity(ctx context.Context, id int64) (*models.Activity, error) {
	return u.repo.GetActivityByID(ctx, id)
}

func (u *activityUsecaseImpl) ListActivities(ctx context.Context, filter models.ActivityFilter) ([]*models.Activity, error) {
	return u.repo.ListActivities(ctx, filter)
}
