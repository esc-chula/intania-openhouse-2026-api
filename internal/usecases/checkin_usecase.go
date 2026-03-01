package usecases

import (
	"context"
	"errors"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/google/uuid"
)

var (
	ErrInvalidCodeFormat = errors.New("invalid code format")
	ErrAlreadyAttended   = errors.New("already attended")
)

type CheckInUsecase interface {
	CheckIn(ctx context.Context, email string, code string) error
}

type checkInUsecaseImpl struct {
	bookingRepo repositories.BookingRepo
	boothRepo   repositories.BoothRepo
	userRepo    repositories.UserRepo
}

func NewCheckInUsecase(
	bookingRepo repositories.BookingRepo,
	boothRepo repositories.BoothRepo,
	userRepo repositories.UserRepo,
) CheckInUsecase {
	return &checkInUsecaseImpl{
		bookingRepo: bookingRepo,
		boothRepo:   boothRepo,
		userRepo:    userRepo,
	}
}

func (u *checkInUsecaseImpl) CheckIn(ctx context.Context, email string, code string) error {
	prefixCode := code[0:2]
	checkInCode := code[2:]

	if _, err := uuid.Parse(checkInCode); err != nil {
		return ErrInvalidCodeFormat
	}

	switch prefixCode {
	case "W-":
		bookingID, bookingStatus, err := u.bookingRepo.GetBookingIDAndStatus(ctx, email, checkInCode)
		if err != nil {
			return err
		}

		if bookingStatus != models.StatusConfirmed {
			if bookingStatus == models.StatusAttended {
				return ErrAlreadyAttended
			}
			return repositories.ErrInvalidBookingStatus
		}

		err = u.bookingRepo.AttendBooking(ctx, bookingID)
		if err != nil {
			return err
		}
	case "B-":
		user, err := u.userRepo.GetUserByEmail(ctx, email, []string{"id"})
		if err != nil {
			return err
		}
		userID := user.ID

		boothID, err := u.boothRepo.GetBoothIDFromCheckInCode(ctx, checkInCode)
		if err != nil {
			return err
		}

		err = u.boothRepo.CreateBoothCheckIn(ctx, userID, boothID)
		if err != nil {
			return err
		}

	default:
		return ErrInvalidCodeFormat
	}
	return nil
}
