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
	ErrAlreadyBooked   = errors.New("user already booked this workshop")
	ErrBookingNotFound = errors.New("booking not found")
)

type BookingRepo interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetConfirmedBookingsWithWorkshop(ctx context.Context, userID int64, eventDate string) ([]models.BookingWithTime, error)
	CancelBooking(ctx context.Context, userID int64, workshopID int64) error
	GetUserBookings(ctx context.Context, userID int64) ([]*models.Booking, error)
}

type bookingRepoImpl struct {
	exec baserepo.Executor
}

func NewBookingRepo(db *bun.DB) BookingRepo {
	return &bookingRepoImpl{
		exec: baserepo.NewExecutor(db),
	}
}

func (r *bookingRepoImpl) CreateBooking(ctx context.Context, booking *models.Booking) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		_, err := idb.NewInsert().Model(booking).Exec(ctx)
		if err != nil {
			if pgErr, ok := err.(pgdriver.Error); ok && pgErr.IntegrityViolation() && pgErr.Field('C') == "23505" {
				return ErrAlreadyBooked
			}
			return err
		}
		return nil
	})
}

func (r *bookingRepoImpl) GetConfirmedBookingsWithWorkshop(ctx context.Context, userID int64, eventDate string) ([]models.BookingWithTime, error) {
	bookings := make([]models.BookingWithTime, 0)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().
			TableExpr("bookings AS bk").
			Column("bk.id", "bk.workshop_id").
			ColumnExpr("ws.start_time").
			ColumnExpr("ws.end_time").
			Join("JOIN workshops AS ws ON ws.id = bk.workshop_id").
			Where("bk.user_id = ?", userID).
			Where("bk.status = ?", models.StatusConfirmed).
			Where("ws.event_date = ?", eventDate).
			Scan(ctx, &bookings)
	})
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepoImpl) CancelBooking(ctx context.Context, userID int64, workshopID int64) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		result, err := idb.NewUpdate().
			Table("bookings").
			Set("status = ?", models.StatusCancelled).
			Where("user_id = ?", userID).
			Where("workshop_id = ?", workshopID).
			Where("status = ?", models.StatusConfirmed).
			Exec(ctx)
		if err != nil {
			return err
		}
		if n, err := result.RowsAffected(); err == nil && n == 0 {
			return ErrBookingNotFound
		}
		return nil
	})
}

func (r *bookingRepoImpl) GetUserBookings(ctx context.Context, userID int64) ([]*models.Booking, error) {
	bookings := make([]*models.Booking, 0)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().
			Model(&bookings).
			Where("user_id = ?", userID).
			Where("status = ?", models.StatusConfirmed).
			Scan(ctx)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return bookings, nil
}
