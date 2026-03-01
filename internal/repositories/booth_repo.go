package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	ErrBoothNotFound         = errors.New("booth not found")
	ErrAlreadyCheckedInBooth = errors.New("already check-in this booth")
)

type BoothRepo interface {
	GetBoothIDFromCheckInCode(ctx context.Context, checkInCode string) (int64, error)
	CreateBoothCheckIn(ctx context.Context, userID int64, boothID int64) error
}

type boothRepoImpl struct {
	exec baserepo.Executor
}

func NewBoothRepo(db *bun.DB) BoothRepo {
	return &boothRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}

func (r *boothRepoImpl) GetBoothIDFromCheckInCode(ctx context.Context, checkInCode string) (int64, error) {
	var boothID int64
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.
			NewSelect().
			Model((*models.Booth)(nil)).
			Column("id").
			Where("check_in_code = ?", checkInCode).
			Scan(ctx, &boothID)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrBoothNotFound
		}
		return 0, err
	}

	return boothID, nil
}

func (r *boothRepoImpl) CreateBoothCheckIn(ctx context.Context, userID int64, boothID int64) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		_, err := idb.NewInsert().Model(&models.BoothCheckIn{
			UserID:      userID,
			BoothID:     boothID,
			CheckedInAt: time.Now(),
		}).Exec(ctx)
		if err != nil {
			if pgErr, ok := err.(pgdriver.Error); ok && pgErr.Field('C') == "23505" {
				return ErrAlreadyCheckedInBooth
			}

			return err
		}

		return nil
	})
}
