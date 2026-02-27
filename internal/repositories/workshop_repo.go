package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
)

var (
	ErrWorkshopNotFound = errors.New("workshop not found")
	ErrWorkshopFull     = errors.New("workshop is full")
)

type WorkshopRepo interface {
	GetWorkshopById(ctx context.Context, id int64, fields []string) (*models.Workshop, error)
	ListWorkshop(ctx context.Context, filter models.WorkshopFilter) ([]*models.Workshop, error)
	IncrementRegisteredCount(ctx context.Context, workshopID int64) error
	DecrementRegisteredCount(ctx context.Context, workshopID int64) error
}

type workshopRepoImpl struct {
	exec baserepo.Executor
}

func NewWorkshopRepo(db *bun.DB) WorkshopRepo {
	return &workshopRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}

func (r *workshopRepoImpl) GetWorkshopById(ctx context.Context, id int64, fields []string) (*models.Workshop, error) {
	workshop := new(models.Workshop)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		query := idb.NewSelect().Model(workshop).Where("id = ?", id)
		if len(fields) > 0 {
			query.Column(fields...)
		}
		return query.Scan(ctx)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkshopNotFound
		}
		return nil, err
	}
	return workshop, nil
}
func (r *workshopRepoImpl) ListWorkshop(ctx context.Context, filter models.WorkshopFilter) ([]*models.Workshop, error) {
	workshops := make([]*models.Workshop, 0)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		query := idb.NewSelect().Model(&workshops)
		if filter.Search != "" {
			query.Where(
				"(name ILIKE ? OR description ILIKE ?)",
				"%"+filter.Search+"%",
				"%"+filter.Search+"%",
			)
		}
		if filter.Category != "" {
			query.Where("category = ?", filter.Category)
		}
		if filter.EventDate != "" {
			query.Where("event_date = ?", filter.EventDate)
		}
		if filter.HideFull {
			query.Where("registered_count < total_seats")
		}
		if filter.SortBy != "" {
			query.Order(filter.SortBy + " " + filter.Order)
		}
		return query.Scan(ctx)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorkshopNotFound
		}
		return nil, err
	}
	return workshops, nil
}

func (r *workshopRepoImpl) IncrementRegisteredCount(ctx context.Context, workshopID int64) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		result, err := idb.NewUpdate().
			Table("workshops").
			Set("registered_count = registered_count + 1").
			Where("id = ?", workshopID).
			Where("registered_count < total_seats"). // race safe
			Exec(ctx)
		if err != nil {
			return err
		}
		if n, err := result.RowsAffected(); err == nil && n == 0 {
			return ErrWorkshopFull
		}
		return nil
	})
}

func (r *workshopRepoImpl) DecrementRegisteredCount(ctx context.Context, workshopID int64) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		result, err := idb.NewUpdate().
			Table("workshops").
			Set("registered_count = registered_count - 1").
			Where("id = ?", workshopID).
			Where("registered_count > 0").
			Exec(ctx)
		if err != nil {
			return err
		}
		if n, err := result.RowsAffected(); err == nil && n == 0 {
			return ErrWorkshopNotFound
		}
		return nil
	})
}
