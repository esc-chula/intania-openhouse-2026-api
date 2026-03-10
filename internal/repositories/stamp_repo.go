package repositories

import (
	"context"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
)

type StampRepo interface {
	CreateStampPosters(ctx context.Context, stampPosters []models.StampPoster) error
}

type stampRepoImpl struct {
	exec baserepo.Executor
}

func NewStampRepo(db *bun.DB) StampRepo {
	return &stampRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}

func (r *stampRepoImpl) CreateStampPosters(ctx context.Context, stampPosters []models.StampPoster) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		_, err := idb.NewInsert().Model(&stampPosters).Exec(ctx)
		return err
	})
}
