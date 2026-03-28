package handlers

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
)

var (
	ErrInvalidCode   = huma.Error400BadRequest("invalid code")
	ErrInvalidCodeFn = func(code string) huma.StatusError {
		return huma.Error400BadRequest(fmt.Sprintf("invalid code %s", code))
	}
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

	huma.Post(checkInGroup, "", handler.CheckIn, func(o *huma.Operation) {
		errDoc, errCodes := buildErrorsDocumentation(checkInErrorList)

		o.Summary = "Check-in with code"
		o.Description = "The code should be formatted in `<type>-<uuid>` where `<type>` is either `W` for workshop or `B` for booth, and `<uuid>` is the identifier for workshop and booth"
		o.Description += errDoc
		o.DefaultStatus = 201
		o.Tags = []string{checkInTag}
		o.Errors = errCodes
	})
}

type CheckInRequest struct {
	Body struct {
		Code string `json:"code"`
	}
}

type CheckInResponse struct {
	Body CheckInResponseBody
}

type CheckInResponseBody struct {
	Type     string               `json:"type"     enum:"workshop,booth"`
	ID       int64                `json:"id"`
	Name     string               `json:"name"`
	Category models.BoothCategory `json:"category" enum:"department,club,exhibition"`
}

var checkInErrorList = []huma.StatusError{ErrEmailNotFound, ErrInvalidCode, ErrAlreadyCheckedIn, ErrInternalServerError()}

func (h *checkInHandler) CheckIn(ctx context.Context, input *CheckInRequest) (*CheckInResponse, error) {
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	result, err := h.checkInUsecase.CheckIn(ctx, email, input.Body.Code)
	if err != nil {
		switch err {
		case usecases.ErrInvalidCodeFormat:
			return nil, ErrInvalidCodeFn(input.Body.Code + " (invalid code format)")

		// case workshop check-in
		case repositories.ErrInvalidCheckInCode:
			return nil, ErrInvalidCodeFn(input.Body.Code + " (invalid check in code)")
		case usecases.ErrAlreadyAttended:
			return nil, ErrAlreadyCheckedIn
		case repositories.ErrInvalidBookingStatus:
			return nil, ErrInvalidCodeFn(input.Body.Code + " (invalid booking status)")

		// case booth check-in
		case repositories.ErrUserNotFound:
			return nil, ErrUserNotFound
		case repositories.ErrBoothNotFound:
			return nil, ErrInvalidCodeFn(input.Body.Code)
		case repositories.ErrAlreadyCheckedInBooth:
			return nil, ErrAlreadyCheckedIn

		default:
			return nil, ErrInternalServerError(err)
		}
	}

	return &CheckInResponse{
		Body: CheckInResponseBody{
			Type:     result.Type,
			ID:       result.ID,
			Name:     result.Name,
			Category: result.Category,
		},
	}, nil
}
