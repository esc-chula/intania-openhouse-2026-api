package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/uptrace/bun"
)

var (
	ErrAlreadyBooked        = errors.New("user already booked this workshop")
	ErrBookingNotFound      = errors.New("booking not found")
	ErrInvalidBookingStatus = errors.New("invalid booking status")
	ErrInvalidCheckInCode   = errors.New("invalid check-in code")
)

type BookingRepo interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	GetUserBookings(ctx context.Context, userID int64) ([]models.BookingWithWorkshop, error)
	CancelBooking(ctx context.Context, userID int64, workshopID int64) error
	UpdateBookingStatus(ctx context.Context, bookingID int64, status models.Status) error
	GetBookingIDAndStatus(ctx context.Context, email string, checkInCode string) (int64, models.Status, error)
	AttendBooking(ctx context.Context, bookingID int64) error
	GetAttendedWorkshopsForUser(ctx context.Context, userID int64) ([]models.StampItem, error)
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
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
				return ErrAlreadyBooked
			}
			return err
		}
		return nil
	})
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

func (r *bookingRepoImpl) GetUserBookings(ctx context.Context, userID int64) ([]models.BookingWithWorkshop, error) {
	bookings := make([]models.BookingWithWorkshop, 0)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().
			TableExpr("bookings AS bk").
			ColumnExpr("bk.id").
			ColumnExpr("bk.workshop_id").
			ColumnExpr("bk.status").
			ColumnExpr("bk.created_at").
			ColumnExpr("bk.checked_in_at").
			ColumnExpr("ws.name AS workshop_name").
			ColumnExpr("ws.event_date").
			ColumnExpr("ws.start_time").
			ColumnExpr("ws.end_time").
			ColumnExpr("ws.location").
			Join("JOIN workshops AS ws ON ws.id = bk.workshop_id").
			Where("bk.user_id = ?", userID).
			Where("bk.status != ?", models.StatusCancelled).
			Scan(ctx, &bookings)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bookings, nil
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

func (r *bookingRepoImpl) GetAttendedWorkshopsForUser(ctx context.Context, userID int64) ([]models.StampItem, error) {
	stamps := make([]models.StampItem, 0)
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().
			TableExpr("bookings AS bk").
			ColumnExpr("ws.id AS id").
			ColumnExpr("ws.name AS name").
			ColumnExpr("bk.checked_in_at AS checked_in_at").
			Join("JOIN workshops AS ws ON ws.id = bk.workshop_id").
			Where("bk.user_id = ?", userID).
			Where("bk.status = ?", models.StatusAttended).
			Scan(ctx, &stamps)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return stamps, nil
		}
		return nil, err
	}
	// Set type for all items
	for i := range stamps {
		stamps[i].Type = models.StampTypeWorkshop
	}
	return stamps, nil
}
