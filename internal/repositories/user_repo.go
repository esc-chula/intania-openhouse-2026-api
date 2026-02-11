package repositories

import (
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
)

// TODO:
type UserRepo interface{}

type userRepoImpl struct {
	exec baserepo.Executor
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &userRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}
