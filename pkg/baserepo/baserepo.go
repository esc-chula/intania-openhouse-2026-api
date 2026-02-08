package baserepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/uptrace/bun"
)

func ErrNotFound(name string) error {
	return huma.Error400BadRequest(fmt.Sprintf("%s not found", name))
}

type BaseRepo[C, U, Q any] interface {
	CreateOne(ctx context.Context, model *C) (id int64, err error)
	FindOneByID(ctx context.Context, id int64) (*Q, error)
	ExistOneByID(ctx context.Context, id int64) (bool, error)
	UpdateOneByID(ctx context.Context, id int64, model *U) error
	DeleteOneByID(ctx context.Context, id int64) error
}

type baseRepo[C, U, Q any] struct {
	executor Executor
	name     string
}

func NewBaseRepo[C, U, Q any](db *bun.DB, name string) BaseRepo[C, U, Q] {
	return &baseRepo[C, U, Q]{
		executor: NewExecutor(db),
		name:     name,
	}
}

func (b *baseRepo[C, U, Q]) CreateOne(ctx context.Context, model *C) (id int64, err error) {
	err = b.executor.Run(ctx, func(idb bun.IDB) error {
		err := idb.NewInsert().Model(model).Returning("id").Scan(ctx, &id)
		return err
	})

	return id, err
}

func (b *baseRepo[C, U, Q]) FindOneByID(ctx context.Context, id int64) (*Q, error) {
	model := new(Q)
	err := b.executor.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().Model(model).Where("id = ?", id).Scan(ctx, model)
	})
	err = TransformErrNotFound(err, b.name)

	return model, err
}

func (b *baseRepo[C, U, Q]) ExistOneByID(ctx context.Context, id int64) (exist bool, err error) {
	err = b.executor.Run(ctx, func(idb bun.IDB) error {
		model := new(Q)
		exist, err = idb.NewSelect().Model(model).Where("id = ?", id).Exists(ctx)
		return err
	})
	return exist, err
}

func (b *baseRepo[C, U, Q]) UpdateOneByID(ctx context.Context, id int64, model *U) error {
	err := b.executor.Run(ctx, func(idb bun.IDB) error {
		q := idb.NewUpdate().Model(model).Where("id = ?", id)

		result, err := q.Exec(ctx)
		if err != nil {
			return err
		}

		if n, err := result.RowsAffected(); err == nil && n == 0 {
			return ErrNotFound(b.name)
		}
		return nil
	})

	return err
}

func (b *baseRepo[C, U, Q]) DeleteOneByID(ctx context.Context, id int64) error {
	err := b.executor.Run(ctx, func(idb bun.IDB) error {
		model := new(Q)
		result, err := idb.NewDelete().Model(model).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return err
		}

		if n, err := result.RowsAffected(); err == nil && n == 0 {
			return ErrNotFound(b.name)
		}
		return nil
	})

	return err
}

func TransformErrNotFound(err error, name string) error {
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound(name)
	}
	return err
}
