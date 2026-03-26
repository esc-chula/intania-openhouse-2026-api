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
	CheckIn(ctx context.Context, email string, code string) (CheckInOutput, error)
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

type CheckInOutput struct {
	Type     string
	ID       int64
	Name     string
	Category models.BoothCategory
}

func (u *checkInUsecaseImpl) CheckIn(ctx context.Context, email string, code string) (CheckInOutput, error) {
	if len(code) <= PrefixLength {
		return CheckInOutput{}, ErrInvalidCodeFormat
	}

	prefixCode := code[0:PrefixLength]
	checkInCode := code[PrefixLength:]

	if err := uuid.Validate(checkInCode); err != nil {
		return CheckInOutput{}, ErrInvalidCodeFormat
	}

	switch prefixCode {
	case PrefixWorkshop:
		return u.handleWorkshopCheckIn(ctx, email, checkInCode)
	case PrefixBooth:
		return u.handleBoothCheckIn(ctx, email, checkInCode)
	default:
		return CheckInOutput{}, ErrInvalidCodeFormat
	}
}

func (u *checkInUsecaseImpl) handleWorkshopCheckIn(ctx context.Context, email string, checkInCode string) (CheckInOutput, error) {
	bookingData, err := u.bookingRepo.GetBookingData(ctx, email, checkInCode)
	if err != nil {
		return CheckInOutput{}, err
	}

	if bookingData.Status != models.StatusConfirmed {
		if bookingData.Status == models.StatusAttended {
			return CheckInOutput{}, ErrAlreadyAttended
		}
		return CheckInOutput{}, repositories.ErrInvalidBookingStatus
	}

	if err := u.bookingRepo.AttendBooking(ctx, bookingData.ID); err != nil {
		return CheckInOutput{}, err
	}

	return CheckInOutput{
		Type:     "workshop",
		ID:       bookingData.WorkshopID,
		Name:     bookingData.WorkshopName,
		Category: models.BoothCategory(bookingData.WorkshopCategory),
	}, nil
}

func (u *checkInUsecaseImpl) handleBoothCheckIn(ctx context.Context, email string, checkInCode string) (CheckInOutput, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email, []string{"id"})
	if err != nil {
		return CheckInOutput{}, err
	}
	userID := user.ID

	booth, err := u.boothRepo.GetBoothFromCheckInCode(ctx, checkInCode)
	if err != nil {
		return CheckInOutput{}, err
	}

	if err := u.boothRepo.CreateBoothCheckIn(ctx, userID, booth.ID); err != nil {
		return CheckInOutput{}, err
	}

	return CheckInOutput{
		Type:     "booth",
		ID:       booth.ID,
		Name:     booth.Name,
		Category: booth.Category,
	}, nil
}
