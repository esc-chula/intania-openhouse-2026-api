package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
)

var ErrActivityNotFound = errors.New("activity not found")

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

		currentDate := "(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Bangkok')::date"
		currentTime := "(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Bangkok')::time"

		if filter.HidePast {
			query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
				return q.Where(fmt.Sprintf("event_date > %s", currentDate)).
					WhereOr(fmt.Sprintf("event_date = %s AND end_time >= %s", currentDate, currentTime))
			})
		}

		if filter.HappeningNow {
			query.Where(fmt.Sprintf("event_date = %s", currentDate)).
				Where(fmt.Sprintf("start_time <= %s", currentTime)).
				Where(fmt.Sprintf("end_time >= %s", currentTime))
		}

		if filter.SortBy != "" {
			if filter.SortBy == "location" {
				query.OrderExpr("building_name ?, room_name ?", bun.Safe(filter.Order), bun.Safe(filter.Order))
			} else {
				query.OrderExpr("? ?", bun.Ident(filter.SortBy), bun.Safe(filter.Order))
			}
		}

		return query.Scan(ctx)
	})
	if err != nil {
		return nil, err
	}
	return activities, nil
}
