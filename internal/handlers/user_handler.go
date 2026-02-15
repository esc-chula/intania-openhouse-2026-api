package handlers

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/myValidator"
)

var (
	ErrExtraAttributesInvalid = huma.Error400BadRequest("extra attributes is invalid")
	ErrAttendanceDateInvalid  = huma.Error400BadRequest("attendance date is invalid")
	ErrEmailNotFound          = huma.Error401Unauthorized("email not found in context")
	ErrInternalServerError    = huma.Error500InternalServerError("internal server error")
)

type userHandler struct {
	api     huma.API
	usecase usecases.UserUsecase
}

func InitUserHandler(api huma.API, usecase usecases.UserUsecase) {
	handler := &userHandler{
		api:     api,
		usecase: usecase,
	}

	huma.Post(api, "/", handler.CreateUser, func(o *huma.Operation) {
		o.Summary = "Register new user"
		o.Description = "Register a new user with the provided details."
	})

	huma.Get(api, "/me", handler.GetUser, func(o *huma.Operation) {
		o.Summary = "Get user details"
		o.Description = "Retrieve the user details for the current user, based on the Authorization header."
	})
}

// Request and Response structs
type CreateUserRequest struct {
	Body struct {
		FirstName       string                 `json:"first_name" validate:"required"`
		LastName        string                 `json:"last_name" validate:"required"`
		Gender          models.Gender          `json:"gender" validate:"required,oneof=male female prefer_not-to-say other"`
		PhoneNumber     string                 `json:"phone_number" validate:"required"`
		ParticipantType models.ParticipantType `json:"participant_type" validate:"required"`

		AttendanceDates      []string        `json:"attendance_dates" validate:"dive,datetime=2006-01-02"`
		InterestedActivities []string        `json:"interested_activities"`
		DiscoveryChannel     []string        `json:"discovery_channel"`
		ExtraAttributes      json.RawMessage `json:"extra_attributes"`
	}
}

type CreateUserResponse struct {
	Body struct {
		User *models.User `json:"user"`
	}
}

func (h *userHandler) CreateUser(ctx context.Context, input *CreateUserRequest) (*CreateUserResponse, error) {
	// Retrieve email from context
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	user := &models.User{
		FirstName:   input.Body.FirstName,
		LastName:    input.Body.LastName,
		Gender:      input.Body.Gender,
		PhoneNumber: input.Body.PhoneNumber,
		Email:       email,

		ParticipantType:      input.Body.ParticipantType,
		AttendanceDates:      input.Body.AttendanceDates,
		InterestedActivities: input.Body.InterestedActivities,
		DiscoveryChannel:     input.Body.DiscoveryChannel,
		ExtraAttributes:      input.Body.ExtraAttributes,
	}

	if err := myValidator.ValidateAttendanceDate(user); err != nil {
		return nil, ErrAttendanceDateInvalid
	}

	if err := myValidator.ValidateExtraAttributes(user); err != nil {
		return nil, ErrExtraAttributesInvalid
	}

	err := h.usecase.CreateUser(ctx, user)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return &CreateUserResponse{
		Body: struct {
			User *models.User `json:"user"`
		}{User: user},
	}, nil
}

type GetUserRequest struct {
	Fields string `query:"fields" explode:"true"`
}

type GetUserResponse struct {
	Body struct {
		User *models.User `json:"user"`
	}
}

func (h *userHandler) GetUser(ctx context.Context, input *GetUserRequest) (*GetUserResponse, error) {
	// Retrieve email from context
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	fields := []string{}

	if input.Fields != "" {
		fields = strings.Split(input.Fields, ",")
	}

	log.Println(fields)
	// default
	if len(fields) == 0 {
		fields = []string{"email"}
	}

	user, err := h.usecase.GetUser(ctx, email, fields)
	if err != nil {
		return nil, ErrInternalServerError
	}

	return &GetUserResponse{
		Body: struct {
			User *models.User `json:"user"`
		}{User: user},
	}, nil
}
