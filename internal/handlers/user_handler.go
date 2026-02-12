package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
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

	huma.Register(api, huma.Operation{
		Method:      "POST",
		Path:        "/",
		Summary:     "Register new user",
		Description: "Register a new user with the provided details.",
	}, handler.CreateUser)

	huma.Get(api, "/me", handler.GetUser, func(o *huma.Operation) {
		o.Summary = "Get user"
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
		// Add other fields as necessary based on models.User
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
		return nil, huma.Error401Unauthorized("Unauthorized: email not found in context")
	}

	user := &models.User{
		FirstName:       input.Body.FirstName,
		LastName:        input.Body.LastName,
		Gender:          input.Body.Gender,
		PhoneNumber:     input.Body.PhoneNumber,
		ParticipantType: input.Body.ParticipantType,
		Email:           email,
	}

	err := h.usecase.CreateUser(ctx, user)
	if err != nil {
		return nil, huma.Error500InternalServerError(err.Error())
	}

	return &CreateUserResponse{
		Body: struct{ User *models.User `json:"user"` }{User: user},
	}, nil
}

func (h *userHandler) GetUser(ctx context.Context, input *struct{}) (*struct{}, error) {
	return nil, nil // Placeholder
}
