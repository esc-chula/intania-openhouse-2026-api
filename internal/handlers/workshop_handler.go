package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/myValidator"
	"log"
	"time"
)

var (
	ErrWorkshopNotFound = huma.Error404NotFound("workshop not found")
	ErrInvalidCategory  = huma.Error400BadRequest("invalid workshop category")
	ErrInvalidEventDate = huma.Error400BadRequest("invalid event date format, expected YYYY-MM-DD")
)

type workshopHandler struct {
	api     huma.API
	usecase usecases.WorkshopUsecase
	mid     middlewares.Middleware
}

func InitWorkshopHandler(api huma.API, usecase usecases.WorkshopUsecase, mid middlewares.Middleware) {
	handler := &workshopHandler{
		api:     api,
		usecase: usecase,
		mid:     mid,
	}

	// api.UseMiddleware(mid.WithAuthContext)

	huma.Get(api, "/{id}", handler.GetWorkshop, func(o *huma.Operation) {
		o.Summary = "Get workshop details"
		o.Description = "Retrieve workshop details by ID path parameter"
		o.DefaultStatus = 200
	})

	huma.Get(api, "", handler.ListWorkshop, func(o *huma.Operation) {
		o.Summary = "Get a list of workshops"
		o.Description = "Retrieve a list of workshops with optional filters"
		o.DefaultStatus = 200
	})
}

type GetWorkshopRequest struct {
	ID     int64    `path:"id"`
	Fields []string `query:"fields" explode:"true" enum:"id,name,description,category,affiliation,event_date,start_time,end_time,location,total_seats,registered_count"`
}
type GetWorkshopResponse struct {
	Body GetWorkshopResponseBody `json:"body"`
}
type GetWorkshopResponseBody struct {
	ID              int64           `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	Category        models.Category `json:"category,omitempty"`
	Affiliation     string          `json:"affiliation,omitempty"`
	EventDate       string          `json:"event_date,omitempty"`
	StartTime       time.Time       `json:"start_time,omitempty"`
	EndTime         time.Time       `json:"end_time,omitempty"`
	Location        string          `json:"location,omitempty"`
	TotalSeats      int             `json:"total_seats,omitempty"`
	RegisteredCount int             `json:"registered_count,omitempty"`
}

func (h *workshopHandler) GetWorkshop(ctx context.Context, input *GetWorkshopRequest) (*GetWorkshopResponse, error) {
	id := input.ID
	fields := input.Fields
	log.Println(fields)
	// default
	if len(fields) == 0 {
		fields = []string{"name"}
	}

	w, err := h.usecase.GetWorkshop(ctx, id, fields)
	if err != nil {
		if err == repositories.ErrWorkshopNotFound {
			return nil, ErrWorkshopNotFound
		}
		return nil, ErrInternalServerError
	}

	return &GetWorkshopResponse{
		Body: GetWorkshopResponseBody{
			ID:              w.ID,
			Name:            w.Name,
			Description:     w.Description,
			Category:        w.Category,
			Affiliation:     w.Affiliation,
			EventDate:       w.EventDate,
			StartTime:       w.StartTime,
			EndTime:         w.EndTime,
			Location:        w.Location,
			TotalSeats:      w.TotalSeats,
			RegisteredCount: w.RegisteredCount,
		},
	}, nil
}

type ListWorkshopRequest struct {
	Search    string `query:"search"`
	Category  string `query:"category"`
	EventDate string `query:"event_date"`
	HideFull  bool   `query:"hide_full" default:"false"`
	SortBy    string `query:"sort_by" enum:"start_time,name" default:"start_time"`
	Order     string `query:"order" enum:"asc,desc" default:"asc"`
}
type ListWorkshopResponse struct {
	Body ListWorkshopResponseBody `json:"body"`
}
type ListWorkshopResponseBody struct {
	Workshops []WorkshopItem `json:"workshops"`
}
type WorkshopItem struct {
	ID              int64           `json:"id"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Category        models.Category `json:"category"`
	Affiliation     string          `json:"affiliation"`
	EventDate       string          `json:"event_date"`
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	Location        string          `json:"location"`
	TotalSeats      int             `json:"total_seats"`
	RegisteredCount int             `json:"registered_count"`
}

func (h *workshopHandler) ListWorkshop(ctx context.Context, input *ListWorkshopRequest) (*ListWorkshopResponse, error) {
	// Validate category if provided
	if input.Category != "" {
		if err := myValidator.ValidateWorkshopCategory(input.Category); err != nil {
			return nil, ErrInvalidCategory
		}
	}

	// Validate event_date format if provided
	if input.EventDate != "" {
		if err := myValidator.ValidateEventDate(input.EventDate); err != nil {
			return nil, ErrInvalidEventDate
		}
	}

	filter := models.WorkshopFilter{
		Search:    input.Search,
		Category:  input.Category,
		EventDate: input.EventDate,
		HideFull:  input.HideFull,
		SortBy:    input.SortBy,
		Order:     input.Order,
	}

	workshops, err := h.usecase.ListWorkshop(ctx, filter)
	if err != nil {
		return nil, ErrInternalServerError
	}

	items := make([]WorkshopItem, 0, len(workshops))

	for _, w := range workshops {
		items = append(items, WorkshopItem{
			ID:              w.ID,
			Name:            w.Name,
			Description:     w.Description,
			Category:        w.Category,
			Affiliation:     w.Affiliation,
			EventDate:       w.EventDate,
			StartTime:       w.StartTime,
			EndTime:         w.EndTime,
			Location:        w.Location,
			TotalSeats:      w.TotalSeats,
			RegisteredCount: w.RegisteredCount,
		})
	}

	return &ListWorkshopResponse{
		Body: ListWorkshopResponseBody{
			Workshops: items,
		},
	}, nil
}
