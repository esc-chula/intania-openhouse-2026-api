package usecases

import (
	"context"
	"sort"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
)

type StampUsecase interface {
	GetUserStamps(ctx context.Context, userID int64) (*models.UserStamps, error)
}

type stampUsecaseImpl struct {
	bookingRepo repositories.BookingRepo
	boothRepo   repositories.BoothRepo
}

func NewStampUsecase(
	bookingRepo repositories.BookingRepo,
	boothRepo repositories.BoothRepo,
) StampUsecase {
	return &stampUsecaseImpl{
		bookingRepo: bookingRepo,
		boothRepo:   boothRepo,
	}
}

func (u *stampUsecaseImpl) GetUserStamps(ctx context.Context, userID int64) (*models.UserStamps, error) {

	// The requirements for workshop stamp are still uncertain for now

	// Get workshop stamps (attended workshops)
	// workshopStamps, err := u.bookingRepo.GetAttendedWorkshopsForUser(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }

	// Get booth stamps (booth check-ins)
	boothStamps, err := u.boothRepo.GetBoothCheckInsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// allStamps := append(workshopStamps, boothStamps...)
	allStamps := boothStamps
	sort.Slice(allStamps, func(i, j int) bool {
		return allStamps[i].CheckedInAt.After(allStamps[j].CheckedInAt)
	})

	return &models.UserStamps{
		TotalCount: int64(len(allStamps)),
		Stamps:     allStamps,
	}, nil
}
