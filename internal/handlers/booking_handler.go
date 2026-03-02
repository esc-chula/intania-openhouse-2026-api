package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
)

var (
	ErrWorkshopFull              = huma.Error400BadRequest("workshop is full")
	ErrTimeConflict              = huma.Error400BadRequest("time conflict with existing booking")
	ErrParticipantTypeNotAllowed = huma.Error403Forbidden("participant type is not allowed")
	ErrAlreadyBooked             = huma.Error400BadRequest("already booked this workshop")
	ErrBookingNotFound           = huma.Error404NotFound("booking not found")
)

type bookingHandler struct {
	bookingUsecase usecases.BookingUsecase
	userUsecase    usecases.UserUsecase
	mid            middlewares.Middleware
}

func InitBookingHandler(
	workshopGroup huma.API,
	userGroup huma.API,
	bookingUsecase usecases.BookingUsecase,
	userUsecase usecases.UserUsecase,
	mid middlewares.Middleware,
) {
	handler := &bookingHandler{
		bookingUsecase: bookingUsecase,
		userUsecase:    userUsecase,
		mid:            mid,
	}

	workshopGroup.UseMiddleware(mid.WithAuthContext)

	huma.Post(workshopGroup, "/{workshop_id}/book", handler.BookWorkshop, func(o *huma.Operation) {
		o.Summary = "Book a workshop"
		o.Description = "Create a booking for a workshop. Prevents double-booking and checks seat availability."
		o.DefaultStatus = 201
	})

	huma.Delete(workshopGroup, "/{workshop_id}/book", handler.CancelBooking, func(o *huma.Operation) {
		o.Summary = "Cancel a workshop booking"
		o.Description = "Cancel an existing workshop booking"
		o.DefaultStatus = 204
	})

	userGroup.UseMiddleware(mid.WithAuthContext)

	huma.Get(userGroup, "/me/bookings", handler.GetMyBookings, func(o *huma.Operation) {
		o.Summary = "Get my bookings"
		o.Description = "Retrieve all confirmed bookings for the current user"
	})
}

type BookWorkshopRequest struct {
	WorkshopID int64 `path:"workshop_id"`
}

type BookWorkshopResponse struct {
	Body *struct{}
}

func (h *bookingHandler) BookWorkshop(ctx context.Context, input *BookWorkshopRequest) (*BookWorkshopResponse, error) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	userEmail, ok := ctx.Value("email").(string)
	if !ok || userEmail == "" {
		return nil, ErrEmailNotFound
	}
	err = h.bookingUsecase.BookWorkshop(ctx, userID, userEmail, input.WorkshopID)
	if err != nil {
		switch err {
		case repositories.ErrWorkshopNotFound:
			return nil, ErrWorkshopNotFound
		case repositories.ErrWorkshopFull:
			return nil, ErrWorkshopFull
		case repositories.ErrAlreadyBooked:
			return nil, ErrAlreadyBooked
		case usecases.ErrTimeConflict:
			return nil, ErrTimeConflict
		case usecases.ErrParticipantTypeNotAllowed:
			return nil, ErrParticipantTypeNotAllowed
		default:
			return nil, ErrInternalServerError
		}
	}

	return &BookWorkshopResponse{}, nil
}

type CancelBookingRequest struct {
	WorkshopID int64 `path:"workshop_id"`
}

type CancelBookingResponse struct {
	Body *struct{}
}

func (h *bookingHandler) CancelBooking(ctx context.Context, input *CancelBookingRequest) (*CancelBookingResponse, error) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = h.bookingUsecase.CancelBooking(ctx, userID, input.WorkshopID)
	if err != nil {
		switch err {
		case repositories.ErrBookingNotFound:
			return nil, ErrBookingNotFound
		default:
			return nil, ErrInternalServerError
		}
	}

	return &CancelBookingResponse{}, nil
}

type GetMyBookingsRequest struct{}

type GetMyBookingsResponse struct {
	Body GetMyBookingsResponseBody `json:"body"`
}

type GetMyBookingsResponseBody struct {
	Bookings []BookingItem `json:"bookings"`
}

type BookingItem struct {
	ID         int64         `json:"id"`
	WorkshopID int64         `json:"workshop_id"`
	Status     models.Status `json:"status"`
	CreatedAt  string        `json:"created_at"`
}

func (h *bookingHandler) GetMyBookings(ctx context.Context, input *GetMyBookingsRequest) (*GetMyBookingsResponse, error) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	bookings, err := h.bookingUsecase.GetMyBookings(ctx, userID)
	if err != nil {
		return nil, ErrInternalServerError
	}

	items := make([]BookingItem, 0, len(bookings))
	for _, b := range bookings {
		items = append(items, BookingItem{
			ID:         b.ID,
			WorkshopID: b.WorkshopID,
			Status:     b.Status,
			CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &GetMyBookingsResponse{
		Body: GetMyBookingsResponseBody{
			Bookings: items,
		},
	}, nil
}

func (h *bookingHandler) getUserIDFromContext(ctx context.Context) (int64, error) {
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return 0, ErrEmailNotFound
	}

	user, err := h.userUsecase.GetUser(ctx, email, []string{"id"})
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return 0, ErrUserNotFound
		}
		return 0, ErrInternalServerError
	}

	return user.ID, nil
}
