package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
)

var (
	ErrBoothNotFound         = errors.New("booth not found")
	ErrAlreadyCheckedInBooth = errors.New("already check-in this booth")
)

type BoothRepo interface {
	GetBoothFromCheckInCode(ctx context.Context, checkInCode string) (*models.Booth, error)
	CreateBoothCheckIn(ctx context.Context, userID int64, boothID int64) error
	GetBoothCheckInsForUser(ctx context.Context, userID int64) ([]models.StampItem, error)
}

type boothRepoImpl struct {
	exec baserepo.Executor
}

func NewBoothRepo(db *bun.DB) BoothRepo {
	return &boothRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}

func (r *boothRepoImpl) GetBoothFromCheckInCode(ctx context.Context, checkInCode string) (*models.Booth, error) {
	var booth models.Booth
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.
			NewSelect().
			Model((*models.Booth)(nil)).
			Where("check_in_code = ?", checkInCode).
			Scan(ctx, &booth)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBoothNotFound
		}
		return nil, err
	}

	return &booth, nil
}

func (r *boothRepoImpl) CreateBoothCheckIn(ctx context.Context, userID int64, boothID int64) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		_, err := idb.NewInsert().Model(&models.BoothCheckIn{
			UserID:      userID,
			BoothID:     boothID,
			CheckedInAt: time.Now(),
		}).Exec(ctx)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
				return ErrAlreadyCheckedInBooth
			}

			return err
		}

		return nil
	})
}

func (r *boothRepoImpl) GetBoothCheckInsForUser(ctx context.Context, userID int64) ([]models.StampItem, error) {
	stamps := make([]models.StampItem, 0)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().
			TableExpr("booth_checkins AS btck").
			ColumnExpr("bt.id AS id").
			ColumnExpr("bt.name AS name").
			ColumnExpr("bt.category AS type").
			ColumnExpr("btck.checked_in_at AS checked_in_at").
			Join("JOIN booths AS bt ON bt.id = btck.booth_id").
			Where("btck.user_id = ?", userID).
			Scan(ctx, &stamps)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return stamps, nil
		}
		return nil, err
	}
	return stamps, nil
}
