package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
)

var (
	ErrInvalidCode      = huma.Error400BadRequest("invalid code")
	ErrAlreadyCheckedIn = huma.Error400BadRequest("already checked in")
)

type checkInHandler struct {
	checkInUsecase usecases.CheckInUsecase
	mid            middlewares.Middleware
}

func InitCheckInHandler(
	checkInGroup huma.API,
	checkInUsecase usecases.CheckInUsecase,
	mid middlewares.Middleware,
) {
	handler := &checkInHandler{
		checkInUsecase: checkInUsecase,
		mid:            mid,
	}
	checkInTag := "check-in"

	checkInGroup.UseMiddleware(mid.WithAuthContext)

	huma.Post(checkInGroup, "", handler.CheckIn, func(o *huma.Operation) {
		o.Summary = "Check-in with code"
		o.Description = "The code should be formatted in `<type>-<uuid>` where <type> is either `W` for workshop or `B` for booth, and <uuid> is the identifier for workshop and booth"
		o.DefaultStatus = 201
		o.Tags = []string{checkInTag}
	})
}

type CheckInRequest struct {
	Body struct {
		Code string `json:"code"`
	}
}

type CheckInResponse struct{}

func (h *checkInHandler) CheckIn(ctx context.Context, input *CheckInRequest) (*CheckInResponse, error) {
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	err := h.checkInUsecase.CheckIn(ctx, email, input.Body.Code)
	if err != nil {
		switch err {
		case usecases.ErrInvalidCodeFormat:
			return nil, ErrInvalidCode
		case usecases.ErrAlreadyAttended:
			return nil, ErrAlreadyCheckedIn
		case repositories.ErrAlreadyCheckedInBooth:
			return nil, ErrAlreadyCheckedIn
		case repositories.ErrInvalidBookingStatus:
			return nil, ErrInvalidCode
		case repositories.ErrBoothNotFound:
			return nil, ErrInvalidCode
		case repositories.ErrInvalidCheckInCode:
			return nil, ErrInvalidCode
		default:
			return nil, ErrInternalServerError
		}
	}

	return &CheckInResponse{}, nil
}
