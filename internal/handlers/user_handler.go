package handlers

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
)

type userHandler struct {
	api     huma.API
	usecase usecases.UserUsecase
}

func InitUserHandler(api huma.API, usecase usecases.UserUsecase) {
	// TODO:
	// handler := &userHandler{
	// 	api:     api,
	// 	usecase: usecase,
	// }
	//
	// api.UseMiddleware(mid.WithAuthContext)
	//
	// huma.Get(api, "/me", handler.GetUser, func(o *huma.Operation) {
	// 	o.Summary = "Get user"
	// 	o.Description = "Retrieve the user details for the current user, based on the Authorization header."
	// })
}
