package repositories

import (
	"context"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
)

// TODO:
type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepoImpl struct {
	exec baserepo.Executor
}

func NewUserRepo(db *bun.DB) UserRepo {
	return &userRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}

func (r *userRepoImpl) CreateUser(ctx context.Context, user *models.User) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		_, err := idb.NewInsert().Model(user).Exec(ctx)
		return err
	})
}

func (r *userRepoImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := new(models.User)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().Model(user).Where("email = ?", email).Scan(ctx)
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
