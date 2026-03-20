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
	CheckInBooth(ctx context.Context, email string, code string) (*models.Booth, error)
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
	PrefixBooth  = "B-"
	PrefixLength = 2
)

func (u *checkInUsecaseImpl) CheckInBooth(ctx context.Context, email string, code string) (*models.Booth, error) {
	if len(code) <= PrefixLength {
		return nil, ErrInvalidCodeFormat
	}

	prefixCode := code[0:PrefixLength]
	checkInCode := code[PrefixLength:]

	if err := uuid.Validate(checkInCode); err != nil {
		return nil, ErrInvalidCodeFormat
	}

	if prefixCode != PrefixBooth {
		return nil, ErrInvalidCodeFormat
	}

	user, err := u.userRepo.GetUserByEmail(ctx, email, []string{"id"})
	if err != nil {
		return nil, err
	}
	userID := user.ID

	booth, err := u.boothRepo.GetBoothFromCheckInCode(ctx, checkInCode)
	if err != nil {
		return nil, err
	}

	if err := u.boothRepo.CreateBoothCheckIn(ctx, userID, booth.ID); err != nil {
		return nil, err
	}

	return booth, nil
}
