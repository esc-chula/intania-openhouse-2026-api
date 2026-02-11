package baserepo

import (
	"context"

	"github.com/uptrace/bun"
)

type dbContextKey struct{}

type Executor interface {
	Run(ctx context.Context, fn func(idb bun.IDB) error) error
}

type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type (
	executorImpl      struct{ db *bun.DB }
	transactionerImpl struct{ db *bun.DB }
)

func NewExecutor(db *bun.DB) Executor {
	return &executorImpl{db}
}

func (e *executorImpl) Run(ctx context.Context, fn func(idb bun.IDB) error) error {
	idb, ok := ctx.Value(dbContextKey{}).(bun.IDB)
	if !ok {
		idb = e.db
	}

	return fn(idb)
}

func NewTransactioner(db *bun.DB) Transactioner {
	return &transactionerImpl{db}
}

func (e *transactionerImpl) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return e.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(context.WithValue(ctx, dbContextKey{}, tx))
	})
}
