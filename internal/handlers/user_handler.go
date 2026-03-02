package handlers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/myValidator"
)

var (
	ErrExtraAttributesRequired = huma.Error400BadRequest("extra attributes required")
	ErrInvalidParticipantType  = huma.Error400BadRequest("invalid participant type")
	ErrExtraAttributesInvalid  = huma.Error400BadRequest("extra attributes is invalid")
	ErrAttendanceDateInvalid   = huma.Error400BadRequest("attendance date is invalid")
	ErrInvalidGender           = huma.Error400BadRequest("invalid gender")
	ErrInvalidTransportMode    = huma.Error400BadRequest("invalid transport mode")
	ErrInvalidOriginLocation   = huma.Error400BadRequest("invalid origin location")
	ErrEmailNotFound           = huma.Error401Unauthorized("email not found in context")
	ErrUserNotFound            = huma.Error404NotFound("user not found")
	ErrUserAlreadyExists       = huma.Error400BadRequest("user already exists")
	ErrInternalServerError     = huma.Error500InternalServerError("internal server error")
)

type userHandler struct {
	api     huma.API
	usecase usecases.UserUsecase
	mid     middlewares.Middleware
}

func InitUserHandler(api huma.API, usecase usecases.UserUsecase, mid middlewares.Middleware) {
	handler := &userHandler{
		api:     api,
		usecase: usecase,
		mid:     mid,
	}

	api.UseMiddleware(mid.WithAuthContext)

	huma.Post(api, "", handler.CreateUser, func(o *huma.Operation) {
		o.Summary = "Register new user"
		o.Description = "Register a new user with the provided details."
		o.DefaultStatus = 200
	})

	huma.Get(api, "/me", handler.GetUser, func(o *huma.Operation) {
		o.Summary = "Get user details"
		o.Description = "Retrieve the user details for the current user, based on the Authorization header."
	})
}

// Request and Response structs
type CreateUserRequest struct {
	Body struct {
		FirstName       string                 `json:"first_name"`
		LastName        string                 `json:"last_name"`
		Gender          models.Gender          `json:"gender"`
		PhoneNumber     string                 `json:"phone_number"`
		ParticipantType models.ParticipantType `json:"participant_type"`
		TransportMode   models.TransportMode   `json:"transport_mode"`
		IsFromBangkok   bool                   `json:"is_from_bangkok"`
		OriginLocation  models.OriginLocation  `json:"origin_location"`

		AttendanceDates      []string `json:"attendance_dates" doc:"date in format 2024-12-31"`
		InterestedActivities []string `json:"interested_activities"`
		DiscoveryChannel     []string `json:"discovery_channel"`

		models.ExtraAttributesFields
	}
}

func (h *userHandler) CreateUser(ctx context.Context, input *CreateUserRequest) (*struct{}, error) {
	// Retrieve email from context
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	extraAttributes, err := myValidator.ValidateExtraAttributes(input.Body.ParticipantType, &input.Body.ExtraAttributesFields)
	if err != nil {
		switch err {
		case myValidator.ErrExtraAttributesRequired:
			return nil, ErrExtraAttributesRequired
		case myValidator.ErrInvalidParticipantType:
			return nil, ErrInvalidParticipantType
		case myValidator.ErrExtraAttributesInvalid:
			return nil, ErrExtraAttributesInvalid
		default:
			return nil, ErrInternalServerError
		}
	}

	user := &models.User{
		FirstName:   input.Body.FirstName,
		LastName:    input.Body.LastName,
		Gender:      input.Body.Gender,
		PhoneNumber: input.Body.PhoneNumber,
		Email:       email,

		ParticipantType:      input.Body.ParticipantType,
		TransportMode:        input.Body.TransportMode,
		IsFromBangkok:        input.Body.IsFromBangkok,
		OriginLocation:       input.Body.OriginLocation,
		AttendanceDates:      input.Body.AttendanceDates,
		InterestedActivities: input.Body.InterestedActivities,
		DiscoveryChannel:     input.Body.DiscoveryChannel,
		ExtraAttributes:      extraAttributes,
	}

	if err := myValidator.ValidateAttendanceDate(user); err != nil {
		return nil, ErrAttendanceDateInvalid
	}

	if err := myValidator.ValidateUserEnums(user); err != nil {
		switch err {
		case myValidator.ErrInvalidGender:
			return nil, ErrInvalidGender
		case myValidator.ErrInvalidParticipantType:
			return nil, ErrInvalidParticipantType
		case myValidator.ErrInvalidTransportMode:
			return nil, ErrInvalidTransportMode
		case myValidator.ErrInvalidOriginLocation:
			return nil, ErrInvalidOriginLocation
		default:
			return nil, huma.Error400BadRequest(err.Error())
		}
	}

	err = h.usecase.CreateUser(ctx, user)
	if err != nil {
		if err == repositories.ErrUserAlreadyExists {
			return nil, ErrUserAlreadyExists
		}
		return nil, ErrInternalServerError
	}

	return nil, nil
}

type GetUserRequest struct {
	Fields []string `query:"fields" explode:"true" enum:"id,email,first_name,last_name,gender,phone_number,participant_type,transport_mode,is_from_bangkok,origin_location,attendance_dates,interested_activities,discovery_channel,extra_attributes"`
}

type GetUserResponse struct {
	Body GetUserResponseBody `json:"body"`
}

type GetUserResponseBody struct {
	ID                   int64                  `json:"id,omitempty"`
	Email                string                 `json:"email,omitempty"`
	FirstName            string                 `json:"first_name,omitempty"`
	LastName             string                 `json:"last_name,omitempty"`
	Gender               models.Gender          `json:"gender,omitempty"`
	PhoneNumber          string                 `json:"phone_number,omitempty"`
	ParticipantType      models.ParticipantType `json:"participant_type,omitempty"`
	TransportMode        models.TransportMode   `json:"transport_mode,omitempty"`
	IsFromBangkok        bool                   `json:"is_from_bangkok,omitempty"`
	OriginLocation       models.OriginLocation  `json:"origin_location,omitempty"`
	AttendanceDates      []string               `json:"attendance_dates,omitempty"`
	InterestedActivities []string               `json:"interested_activities,omitempty"`
	DiscoveryChannel     []string               `json:"discovery_channel,omitempty"`
	ExtraAttributes      json.RawMessage        `json:"extra_attributes,omitempty"`
}

func (h *userHandler) GetUser(ctx context.Context, input *GetUserRequest) (*GetUserResponse, error) {
	// Retrieve email from context
	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return nil, ErrEmailNotFound
	}

	fields := input.Fields
	log.Println(fields)
	// default
	if len(fields) == 0 {
		fields = []string{"email"}
	}

	user, err := h.usecase.GetUser(ctx, email, fields)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternalServerError
	}

	return &GetUserResponse{
		Body: GetUserResponseBody{
			ID:                   user.ID,
			Email:                user.Email,
			FirstName:            user.FirstName,
			LastName:             user.LastName,
			Gender:               user.Gender,
			PhoneNumber:          user.PhoneNumber,
			ParticipantType:      user.ParticipantType,
			TransportMode:        user.TransportMode,
			IsFromBangkok:        user.IsFromBangkok,
			OriginLocation:       user.OriginLocation,
			AttendanceDates:      user.AttendanceDates,
			InterestedActivities: user.InterestedActivities,
			DiscoveryChannel:     user.DiscoveryChannel,
			ExtraAttributes:      user.ExtraAttributes,
		},
	}, nil
}
