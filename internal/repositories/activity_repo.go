package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
)

var ErrActivityNotFound = errors.New("activity not found")
var loc, _ = time.LoadLocation("Asia/Bangkok")
var now = time.Now().In(loc)

type ActivityRepo interface {
	GetActivityByID(ctx context.Context, id int64) (*models.Activity, error)
	ListActivities(ctx context.Context, filter models.ActivityFilter) ([]*models.Activity, error)
}

type activityRepoImpl struct {
	exec baserepo.Executor
}

func NewActivityRepo(db *bun.DB) ActivityRepo {
	return &activityRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}

func (r *activityRepoImpl) GetActivityByID(ctx context.Context, id int64) (*models.Activity, error) {
	activity := new(models.Activity)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().Model(activity).Where("id = ?", id).Scan(ctx)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrActivityNotFound
		}
		return nil, err
	}
	return activity, nil
}

func (r *activityRepoImpl) ListActivities(ctx context.Context, filter models.ActivityFilter) ([]*models.Activity, error) {
	activities := make([]*models.Activity, 0)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		query := idb.NewSelect().Model(&activities)

		if filter.Search != "" {
			query.Where(
				"(title ILIKE ? OR description ILIKE ? OR building_name ILIKE ? OR room_name ILIKE ?)",
				"%"+filter.Search+"%",
				"%"+filter.Search+"%",
				"%"+filter.Search+"%",
				"%"+filter.Search+"%",
			)
		}

		if filter.HidePast {
			query.Where("end_time >= ?", now)
		}

		if filter.HappeningNow {
			query.Where("start_time <= ? AND end_time >= ?", now, now)
		}

		if filter.SortBy != "" {
			if filter.SortBy == "location" {
				query.Order("building_name " + filter.Order).Order("room_name " + filter.Order)
			} else {
				query.Order(filter.SortBy + " " + filter.Order)
			}
		}

		return query.Scan(ctx)
	})

	if err != nil {
		return nil, err
	}
	return activities, nil
}
