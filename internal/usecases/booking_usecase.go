package usecases

import (
	"context"
	"errors"
	"sort"
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
	GetMyBookings(ctx context.Context, userID int64) ([]models.BookingWithWorkshop, error)
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
	workshop, err := u.workshopRepo.GetWorkshopById(ctx, workshopID, []string{"id", "event_date", "start_time", "end_time", "total_seats", "registered_count", "category"})
	if err != nil {
		return err
	}
	if *workshop.RegisteredCount >= *workshop.TotalSeats {
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
	// check for club's workshop
	if *workshop.Category == models.WorkShopCategoryClub && user.ParticipantType != models.ParticipantTypeStudent {
		return ErrParticipantTypeNotAllowed
	}

	// Get user's existing confirmed bookings for the same date (with time info)
	existingBookings, err := u.bookingRepo.GetUserBookings(ctx, userID)
	if err != nil {
		return err
	}
	// Check for time overlap
	targetStart := *workshop.StartTime
	targetEnd := *workshop.EndTime
	for _, b := range existingBookings {
		if targetStart.Before(b.EndTime) && targetEnd.After(b.StartTime) && *workshop.EventDate == b.EventDate && b.Status == models.StatusConfirmed {
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

func (u *bookingUsecaseImpl) GetMyBookings(ctx context.Context, userID int64) ([]models.BookingWithWorkshop, error) {
	bookings, err := u.bookingRepo.GetUserBookings(ctx, userID)
	if err != nil {
		return nil, err
	}

	statusPriority := map[models.Status]int{
		models.StatusConfirmed: 0,
		models.StatusAttended:  1,
		models.StatusAbsent:    2,
	}

	// Sort by status priority, then by event_date and start_time
	sort.Slice(bookings, func(i, j int) bool {
		// First sort by status priority
		pi := statusPriority[bookings[i].Status]
		pj := statusPriority[bookings[j].Status]
		if pi != pj {
			return pi < pj
		}
		// Then sort by event_date
		if bookings[i].EventDate != bookings[j].EventDate {
			return bookings[i].EventDate < bookings[j].EventDate
		}
		// Then sort by start_time
		return bookings[i].StartTime.Before(bookings[j].StartTime)
	})

	return bookings, nil
}

func (u *bookingUsecaseImpl) UpdateBookingStatus(ctx context.Context, bookingID int64, status models.Status) error {
	return u.bookingRepo.UpdateBookingStatus(ctx, bookingID, status)
}
