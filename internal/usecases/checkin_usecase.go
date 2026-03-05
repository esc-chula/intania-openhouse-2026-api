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

const (
	PrefixWorkshop = "W-"
	PrefixBooth    = "B-"
	PrefixLength   = 2
)

func (u *checkInUsecaseImpl) CheckIn(ctx context.Context, email string, code string) error {
	if len(code) <= PrefixLength {
		return ErrInvalidCodeFormat
	}

	prefixCode := code[0:PrefixLength]
	checkInCode := code[PrefixLength:]

	if _, err := uuid.Parse(checkInCode); err != nil {
		return ErrInvalidCodeFormat
	}

	switch prefixCode {
	case PrefixWorkshop:
		return u.handleWorkshopCheckIn(ctx, email, checkInCode)
	case PrefixBooth:
		return u.handleBoothCheckIn(ctx, email, checkInCode)
	default:
		return ErrInvalidCodeFormat
	}
}

func (u *checkInUsecaseImpl) handleWorkshopCheckIn(ctx context.Context, email string, checkInCode string) error {
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

	return u.bookingRepo.AttendBooking(ctx, bookingID)
}

func (u *checkInUsecaseImpl) handleBoothCheckIn(ctx context.Context, email string, checkInCode string) error {
	user, err := u.userRepo.GetUserByEmail(ctx, email, []string{"id"})
	if err != nil {
		return err
	}
	userID := user.ID

	boothID, err := u.boothRepo.GetBoothIDFromCheckInCode(ctx, checkInCode)
	if err != nil {
		return err
	}

	return u.boothRepo.CreateBoothCheckIn(ctx, userID, boothID)
}
