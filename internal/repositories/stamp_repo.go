package repositories

import (
	"context"
	"errors"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
)

var (
	ErrStampPosterNotFound = errors.New("stamp poster not found")
)

type StampRepo interface {
	CreateStampPosters(ctx context.Context, stampPosters []models.StampPoster) error
	GetUserStampPosters(ctx context.Context, userID int64) ([]models.StampPoster, error)
	RedeemStamps(ctx context.Context, userID int64, category models.StampType) error
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

func (r *stampRepoImpl) GetUserStampPosters(ctx context.Context, userID int64) ([]models.StampPoster, error) {

	var stampPosters []models.StampPoster

	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.
			NewSelect().
			TableExpr("stamp_posters").
			ColumnExpr("id").
			ColumnExpr("user_id").
			ColumnExpr("type").
			ColumnExpr("is_redeemed").
			Where("user_id = ?", userID).
			Scan(ctx, &stampPosters)
	})

	if err != nil {
		return nil, err
	}

	return stampPosters, nil
}

func (r *stampRepoImpl) RedeemStamps(ctx context.Context, userID int64, category models.StampType) error {

	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		res, err := idb.
			NewUpdate().
			Table("stamp_posters").
			Set("is_redeemed = ?", true).
			Where("user_id = ? AND type = ?", userID, category).
			Exec(ctx)
		if err != nil {
			return err
		}

		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrStampPosterNotFound
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
