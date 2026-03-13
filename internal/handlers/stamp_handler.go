package handlers

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
)

var (
	ErrStampPosterAlreadyRedeemed = huma.Error400BadRequest("stamps were redeemd")
	ErrNotEnoughStamps            = huma.Error400BadRequest("not enough stamps to redeem")
	ErrStampPosterNotFound        = huma.Error404NotFound("stamp poster not found")
	ErrInvalidStampCategory       = huma.Error400BadRequest("invalid stamp category")
)

type stampHandler struct {
	stampUsecase usecases.StampUsecase
	userUsecase  usecases.UserUsecase
	mid          middlewares.Middleware
}

func InitStampHandler(
	stampGroup huma.API,
	userGroup huma.API,
	stampUsecase usecases.StampUsecase,
	userUsecase usecases.UserUsecase,
	mid middlewares.Middleware,
) {
	handler := &stampHandler{
		stampUsecase: stampUsecase,
		userUsecase:  userUsecase,
		mid:          mid,
	}

	stampTag := "stamp"

	huma.Get(userGroup, "/me/stamps", handler.GetUserStamps, func(o *huma.Operation) {
		errDoc, errCodes := buildErrorsDocumentation(getUserStampsErrorList)
		o.Summary = "Get user stamps"
		o.Description = "Retrieve user stamps and checked in details for booth and workshop. (to count stamps in each category)" + errDoc
		o.DefaultStatus = 200
		o.Tags = []string{stampTag}
		o.Errors = errCodes
	})

	huma.Get(userGroup, "/me/redemption-status", handler.GetRedemptionStatus, func(o *huma.Operation) {
		errDoc, errCodes := buildErrorsDocumentation(getRedemptionStatusErrorList)
		o.Summary = "Get stamp redemption status"
		o.Description = "Retrieve redemption status for each category. (to check whether the redemption is possible)" + errDoc
		o.DefaultStatus = 200
		o.Tags = []string{stampTag}
		o.Errors = errCodes
	})

	huma.Post(stampGroup, "/redemptions", handler.RedeemStamps, func(o *huma.Operation) {
		errDoc, errCodes := buildErrorsDocumentation(redeemStampsErrorList)
		o.Summary = "Redeem stamps"
		o.Description = "Redeem stamps for a specific category." + errDoc
		o.DefaultStatus = 200
		o.Tags = []string{stampTag}
		o.Errors = errCodes
	})
}

var (
	getUserStampsErrorList       = []huma.StatusError{ErrEmailNotFound, ErrUserNotFound, ErrInternalServerError}
	getRedemptionStatusErrorList = []huma.StatusError{ErrEmailNotFound, ErrUserNotFound, ErrInternalServerError}
	redeemStampsErrorList        = []huma.StatusError{ErrEmailNotFound, ErrUserNotFound, ErrStampPosterAlreadyRedeemed, ErrNotEnoughStamps, ErrStampPosterNotFound, ErrInternalServerError}
)

type GetUserStampsRequest struct{}

type GetUserStampsResponse struct {
	Body GetUserStampsResponseBody `json:"body"`
}

type GetUserStampsResponseBody struct {
	TotalCount           int64           `json:"total_count" doc:"Overall count of all stamps collected"`
	DepartmentStampCount int64           `json:"department_stamp_count" doc:"Number of department stamps collected"`
	ClubStampCount       int64           `json:"club_stamp_count" doc:"Number of club stamps collected"`
	ExhibitionStampCount int64           `json:"exhibition_stamp_count" doc:"Number of exhibition stamps collected"`
	DepartmentStamps     []StampItemBody `json:"department_stamps" doc:"List of specific department stamps"`
	ClubStamps           []StampItemBody `json:"club_stamps" doc:"List of specific club stamps"`
	ExhibitionStamps     []StampItemBody `json:"exhibition_stamps" doc:"List of specific exhibition stamps"`
}

type StampItemBody struct {
	ID          int64            `json:"id"`
	Type        models.StampType `json:"type"`
	Name        string           `json:"name"`
	CheckedInAt time.Time        `json:"checked_in_at"`
}

func (h *stampHandler) GetUserStamps(ctx context.Context, input *GetUserStampsRequest) (*GetUserStampsResponse, error) {
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	user, err := h.userUsecase.GetUser(ctx, email, []string{"id"})
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServerError
	}

	stamps, err := h.stampUsecase.GetUserStamps(ctx, user.ID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	var departmentStamps, clubStamps, exhibitionStamps []StampItemBody

	for _, s := range stamps.DepartmentStamps {
		departmentStamps = append(departmentStamps, StampItemBody{
			ID:          s.ID,
			Type:        s.Type,
			Name:        s.Name,
			CheckedInAt: s.CheckedInAt,
		})
	}
	for _, s := range stamps.ClubStamps {
		clubStamps = append(clubStamps, StampItemBody{
			ID:          s.ID,
			Type:        s.Type,
			Name:        s.Name,
			CheckedInAt: s.CheckedInAt,
		})
	}
	for _, s := range stamps.ExhibitionStamps {
		exhibitionStamps = append(exhibitionStamps, StampItemBody{
			ID:          s.ID,
			Type:        s.Type,
			Name:        s.Name,
			CheckedInAt: s.CheckedInAt,
		})
	}

	return &GetUserStampsResponse{
		Body: GetUserStampsResponseBody{
			TotalCount:           stamps.TotalCount,
			DepartmentStampCount: stamps.DepartmentStampCount,
			ClubStampCount:       stamps.ClubStampCount,
			ExhibitionStampCount: stamps.ExhibitionStampCount,
			DepartmentStamps:     departmentStamps,
			ClubStamps:           clubStamps,
			ExhibitionStamps:     exhibitionStamps,
		},
	}, nil
}

type GetRedemptionStatusRequest struct{}

type GetRedemptionStatusResponse struct {
	Body GetRedemptionStatusResponseBody
}

type GetRedemptionStatusResponseBody struct {
	Department RedemptionStatusItem `json:"department" doc:"Redemption status for department stamps"`
	Club       RedemptionStatusItem `json:"club" doc:"Redemption status for club stamps"`
	Exhibition RedemptionStatusItem `json:"exhibition" doc:"Redemption status for exhibition stamps"`
}

type RedemptionStatusItem struct {
	Redeemable bool `json:"redeemable" doc:"True if user has enough stamps and hasn't redeemed yet"`
	IsRedeemed bool `json:"is_redeemed" doc:"True if user has already redeemed reward"`
}

func (h *stampHandler) GetRedemptionStatus(ctx context.Context, input *GetRedemptionStatusRequest) (*GetRedemptionStatusResponse, error) {
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	user, err := h.userUsecase.GetUser(ctx, email, []string{"id"})
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServerError
	}

	status, err := h.stampUsecase.GetMyStampPosters(ctx, user.ID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return &GetRedemptionStatusResponse{
		Body: GetRedemptionStatusResponseBody{
			Department: RedemptionStatusItem{
				Redeemable: status.DepartmentRedeemable,
				IsRedeemed: status.DepartmentIsRedeemed,
			},
			Club: RedemptionStatusItem{
				Redeemable: status.ClubRedeemable,
				IsRedeemed: status.ClubIsRedeemed,
			},
			Exhibition: RedemptionStatusItem{
				Redeemable: status.ExhibitionRedeemable,
				IsRedeemed: status.ExhibitionIsRedeemed,
			},
		},
	}, nil
}

type RedeemStampsRequest struct {
	Category models.StampType `query:"category" enum:"department,club,exhibition" doc:"Stamp category to redeem (department, club, exhibition)"`
}

type RedeemStampsResponse struct{}

func (h *stampHandler) RedeemStamps(ctx context.Context, input *RedeemStampsRequest) (*RedeemStampsResponse, error) {
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	user, err := h.userUsecase.GetUser(ctx, email, []string{"id"})
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServerError
	}

	err = h.stampUsecase.RedeemStamps(ctx, user.ID, input.Category)
	if err != nil {
		switch err {
		case usecases.ErrStampPosterAlreadyRedeemed:
			return nil, ErrStampPosterAlreadyRedeemed
		case usecases.ErrInvalidStampCategory:
			return nil, ErrInvalidStampCategory
		case usecases.ErrNotEnoughStamps:
			return nil, ErrNotEnoughStamps
		case repositories.ErrStampPosterNotFound:
			return nil, ErrStampPosterNotFound
		default:
			return nil, ErrInternalServerError
		}
	}

	return &RedeemStampsResponse{}, nil
}
