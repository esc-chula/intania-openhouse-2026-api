package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
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
		o.Summary = "Get user stamps"
		o.Description = "Retrieve user stamps and checked in details for booth and workshop."
		o.Tags = []string{stampTag}
	})
}

type GetUserStampsRequest struct{}

type GetUserStampsResponse struct {
	Body GetUserStampsResponseBody `json:"body"`
}

type GetUserStampsResponseBody struct {
	TotalCount           int64           `json:"total_count"`
	DepartmentStampCount int64           `json:"department_stamp_count"`
	ClubStampCount       int64           `json:"club_stamp_count"`
	ExhibitionStampCount int64           `json:"exhibition_stamp_count"`
	DepartmentStamps     []StampItemBody `json:"department_stamps"`
	ClubStamps           []StampItemBody `json:"club_stamps"`
	ExhibitionStamps     []StampItemBody `json:"exhibition_stamps"`
}

type StampItemBody struct {
	ID          int64  `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	CheckedInAt string `json:"checked_in_at"`
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

	departmentStamps := make([]StampItemBody, 0, len(stamps.Stamps))
	clubStamps := make([]StampItemBody, 0, len(stamps.Stamps))
	exhibitionStamps := make([]StampItemBody, 0, len(stamps.Stamps))
	for _, s := range stamps.Stamps {
		item := StampItemBody{
			ID:          s.ID,
			Type:        string(s.Type),
			Name:        s.Name,
			CheckedInAt: s.CheckedInAt.Format("2006-01-02T15:04:05Z07:00"),
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

	return &GetUserStampsResponse{
		Body: GetUserStampsResponseBody{
			TotalCount:           stamps.TotalCount,
			DepartmentStampCount: int64(len(departmentStamps)),
			ClubStampCount:       int64(len(clubStamps)),
			ExhibitionStampCount: int64(len(exhibitionStamps)),
			DepartmentStamps:     departmentStamps,
			ClubStamps:           clubStamps,
			ExhibitionStamps:     exhibitionStamps,
		},
	}, nil
}
