package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string, fields []string) (*models.User, error)
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
		if err != nil {
			if pgErr, ok := err.(pgdriver.Error); ok && pgErr.IntegrityViolation() && pgErr.Field('C') == "23505" {
				return ErrUserAlreadyExists
			}
			return err
		}
		return nil
	})
}

func (r *userRepoImpl) GetUserByEmail(ctx context.Context, email string, fields []string) (*models.User, error) {
	user := new(models.User)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {

		query := idb.NewSelect().
			Model(user).
			Where("email = ?", email)

		if len(fields) > 0 {
			query.Column(fields...)
		}

		return query.Scan(ctx)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
