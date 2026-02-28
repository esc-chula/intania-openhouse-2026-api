package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/baserepo"
)

var (
	ErrTimeConflict              = errors.New("time conflict with existing booking")
	ErrParticipantTypeNotAllowed = errors.New("participant type is not allowed")
	ErrBookingNotFound           = errors.New("booking not found")
)

type BookingUsecase interface {
	BookWorkshop(ctx context.Context, userID int64, userEmail string, workshopID int64) error
	CancelBooking(ctx context.Context, userID int64, workshopID int64) error
	GetMyBookings(ctx context.Context, userID int64) ([]*models.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID int64, status models.Status) error
}

type bookingUsecaseImpl struct {
	bookingRepo   repositories.BookingRepo
	workshopRepo  repositories.WorkshopRepo
	userRepo      repositories.UserRepo
	transactioner baserepo.Transactioner
}

func NewBookingUsecase(
	bookingRepo repositories.BookingRepo,
	workshopRepo repositories.WorkshopRepo,
	userRepo repositories.UserRepo,
	transactioner baserepo.Transactioner,
) BookingUsecase {
	return &bookingUsecaseImpl{
		bookingRepo:   bookingRepo,
		workshopRepo:  workshopRepo,
		userRepo:      userRepo,
		transactioner: transactioner,
	}
}

func (u *bookingUsecaseImpl) BookWorkshop(ctx context.Context, userID int64, userEmail string, workshopID int64) error {
	workshop, err := u.workshopRepo.GetWorkshopById(ctx, workshopID, []string{"id", "event_date", "start_time", "end_time", "total_seats", "registered_count"})
	if err != nil {
		return err
	}
	if workshop.RegisteredCount >= workshop.TotalSeats {
		return repositories.ErrWorkshopFull
	}

	// check participant type
	user, err := u.userRepo.GetUserByEmail(ctx, userEmail, []string{"participant_type"})
	if err != nil {
		return err
	}
	switch user.ParticipantType {
	case models.ParticipantTypeAlumni, models.ParticipantTypeTeacher, models.ParticipantTypeOther:
		return ErrParticipantTypeNotAllowed
	}

	// Get user's existing confirmed bookings for the same date (with time info)
	existingBookings, err := u.bookingRepo.GetConfirmedBookingsWithWorkshop(ctx, userID, workshop.EventDate)
	if err != nil {
		return err
	}
	// Check for time overlap
	targetStart := workshop.StartTime
	targetEnd := workshop.EndTime
	for _, b := range existingBookings {
		if targetStart.Before(b.EndTime) && targetEnd.After(b.StartTime) {
			return ErrTimeConflict
		}
	}

	return u.transactioner.Transaction(ctx, func(ctx context.Context) error {
		booking := &models.Booking{
			UserID:     userID,
			WorkshopID: workshopID,
			Status:     models.StatusConfirmed,
			CreatedAt:  time.Now(),
		}
		if err := u.bookingRepo.CreateBooking(ctx, booking); err != nil {
			return err
		}
		return u.workshopRepo.IncrementRegisteredCount(ctx, workshopID)
	})
}

func (u *bookingUsecaseImpl) CancelBooking(ctx context.Context, userID int64, workshopID int64) error {
	return u.transactioner.Transaction(ctx, func(ctx context.Context) error {
		err := u.bookingRepo.CancelBooking(ctx, userID, workshopID)
		if err != nil {
			return err
		}
		return u.workshopRepo.DecrementRegisteredCount(ctx, workshopID)
	})
}

func (u *bookingUsecaseImpl) GetMyBookings(ctx context.Context, userID int64) ([]*models.Booking, error) {
	return u.bookingRepo.GetUserBookings(ctx, userID)
}

func (u *bookingUsecaseImpl) UpdateBookingStatus(ctx context.Context, bookingID int64, status models.Status) error {
	return u.bookingRepo.UpdateBookingStatus(ctx, bookingID, status)
}
