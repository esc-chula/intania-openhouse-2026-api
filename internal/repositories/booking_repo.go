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
	ErrAlreadyBooked        = errors.New("user already booked this workshop")
	ErrBookingNotFound      = errors.New("booking not found")
	ErrInvalidBookingStatus = errors.New("invalid booking status")
	ErrInvalidCheckInCode   = errors.New("invalid check-in code")
)

type BookingRepo interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetConfirmedBookingsWithWorkshop(ctx context.Context, userID int64, eventDate string) ([]models.BookingWithTime, error)
	CancelBooking(ctx context.Context, userID int64, workshopID int64) error
	GetUserBookings(ctx context.Context, userID int64) ([]*models.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID int64, status models.Status) error
	GetBookingIDAndStatus(ctx context.Context, email string, checkInCode string) (int64, models.Status, error)
	AttendBooking(ctx context.Context, bookingID int64) error
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

func (r *bookingRepoImpl) UpdateBookingStatus(ctx context.Context, bookingID int64, status models.Status) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		booking := new(models.Booking)
		result, err := idb.NewUpdate().
			Model(booking).
			Set("status = ?", status).
			Where("id = ?", bookingID).
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

func (r *bookingRepoImpl) AttendBooking(ctx context.Context, bookingID int64) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		result, err := idb.NewUpdate().
			Model((*models.Booking)(nil)).
			Set("status = ?", models.StatusAttended).
			Set("checked_in_at = ?", time.Now()).
			Where("id = ?", bookingID).
			Where("status = ?", models.StatusConfirmed). // race safe
			Exec(ctx)
		if err != nil {
			return err
		}

		if n, err := result.RowsAffected(); err == nil && n == 0 {
			return ErrInvalidBookingStatus
		}

		return err
	})
}

func (r *bookingRepoImpl) GetBookingIDAndStatus(ctx context.Context, email string, checkInCode string) (int64, models.Status, error) {
	var booking models.Booking
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		err := idb.NewSelect().
			Model((*models.Booking)(nil)).
			ColumnExpr("bk.id, bk.status").
			Join("JOIN users AS u ON u.id = bk.user_id").
			Join("JOIN workshops AS ws ON ws.id = bk.workshop_id").
			Where("u.email = ?", email).
			Where("ws.check_in_code = ?", checkInCode).
			Scan(ctx, &booking)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrInvalidCheckInCode
			}
			return err
		}

		return nil
	})

	return booking.ID, booking.Status, err
}
