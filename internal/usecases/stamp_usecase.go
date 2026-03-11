package usecases

import (
	"context"
	"errors"
	"sort"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
)

const (
	MinStampsToRedeem = 5
)

var (
	ErrAlreadyRedeemed = errors.New("stamps were redeemed")
	ErrNotEnoughStamps = errors.New("not enough stamps to redeem")
)

type StampUsecase interface {
	GetUserStamps(ctx context.Context, userID int64) (*models.UserStamps, error)
	GetMyStampPosters(ctx context.Context, userID int64) (*models.StampRedemptionStatus, error)
	RedeemStamps(ctx context.Context, userID int64, category models.StampType) error
}

type stampUsecaseImpl struct {
	stampRepo   repositories.StampRepo
	bookingRepo repositories.BookingRepo
	boothRepo   repositories.BoothRepo
}

func NewStampUsecase(
	stampRepo repositories.StampRepo,
	bookingRepo repositories.BookingRepo,
	boothRepo repositories.BoothRepo,
) StampUsecase {
	return &stampUsecaseImpl{
		stampRepo:   stampRepo,
		bookingRepo: bookingRepo,
		boothRepo:   boothRepo,
	}
}

func (u *stampUsecaseImpl) GetUserStamps(ctx context.Context, userID int64) (*models.UserStamps, error) {

	// Get booth stamps (booth check-ins)
	stamps, err := u.boothRepo.GetBoothCheckInsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	sort.Slice(stamps, func(i, j int) bool {
		return stamps[i].CheckedInAt.After(stamps[j].CheckedInAt)
	})

	departmentStamps := make([]models.StampItem, 0, len(stamps))
	clubStamps := make([]models.StampItem, 0, len(stamps))
	exhibitionStamps := make([]models.StampItem, 0, len(stamps))

	for _, s := range stamps {
		item := models.StampItem{
			ID:          s.ID,
			Type:        s.Type,
			Name:        s.Name,
			CheckedInAt: s.CheckedInAt,
		}
		switch s.Type {
		case models.StampTypeDepartment:
			departmentStamps = append(departmentStamps, item)
		case models.StampTypeClub:
			clubStamps = append(clubStamps, item)
		case models.StampTypeExhibition:
			exhibitionStamps = append(exhibitionStamps, item)
		}
	}

	return &models.UserStamps{
		TotalCount:           int64(len(stamps)),
		DepartmentStampCount: int64(len(departmentStamps)),
		ClubStampCount:       int64(len(clubStamps)),
		ExhibitionStampCount: int64(len(exhibitionStamps)),
		DepartmentStamps:     departmentStamps,
		ClubStamps:           clubStamps,
		ExhibitionStamps:     exhibitionStamps,
	}, nil
}

func (u *stampUsecaseImpl) GetMyStampPosters(ctx context.Context, userID int64) (*models.StampRedemptionStatus, error) {
	posters, err := u.stampRepo.GetUserStampPosters(ctx, userID)
	if err != nil {
		return nil, err
	}

	stamps, err := u.GetUserStamps(ctx, userID)
	if err != nil {
		return nil, err
	}

	// check sufficient stamps count
	result := &models.StampRedemptionStatus{
		DepartmentRedeemable: stamps.DepartmentStampCount >= MinStampsToRedeem,
		ClubRedeemable:       stamps.ClubStampCount >= MinStampsToRedeem,
		ExhibitionRedeemable: stamps.ExhibitionStampCount >= MinStampsToRedeem,
	}

	for _, p := range posters {
		switch p.Type {
		case models.StampTypeDepartment:
			result.DepartmentIsRedeemed = p.IsRedeemed
		case models.StampTypeClub:
			result.ClubIsRedeemed = p.IsRedeemed
		case models.StampTypeExhibition:
			result.ExhibitionIsRedeemed = p.IsRedeemed
		}
	}

	// check if stamps had been redeemed
	result.DepartmentRedeemable = result.DepartmentRedeemable && !result.DepartmentIsRedeemed
	result.ClubRedeemable = result.ClubRedeemable && !result.ClubIsRedeemed
	result.ExhibitionRedeemable = result.ExhibitionRedeemable && !result.ExhibitionIsRedeemed

	return result, nil
}

func (u *stampUsecaseImpl) RedeemStamps(ctx context.Context, userID int64, category models.StampType) error {

	status, err := u.GetMyStampPosters(ctx, userID)
	if err != nil {
		return err
	}

	redeemable := false
	alreadyRedeemed := false

	switch category {
	case models.StampTypeDepartment:
		redeemable = status.DepartmentRedeemable
		alreadyRedeemed = status.DepartmentIsRedeemed
	case models.StampTypeClub:
		redeemable = status.ClubRedeemable
		alreadyRedeemed = status.ClubIsRedeemed
	case models.StampTypeExhibition:
		redeemable = status.ExhibitionRedeemable
		alreadyRedeemed = status.ExhibitionIsRedeemed
	default:
		return nil
	}

	if alreadyRedeemed {
		return ErrAlreadyRedeemed
	}

	if !redeemable {
		return ErrNotEnoughStamps
	}

	if err := u.stampRepo.RedeemStamps(ctx, userID, category); err != nil {
		return err
	}
	return nil

}
