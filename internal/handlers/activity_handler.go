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
	ErrActivityNotFound = huma.Error404NotFound("activity not found")
)

type activityHandler struct {
	api     huma.API
	usecase usecases.ActivityUsecase
	mid     middlewares.Middleware
}

func InitActivityHandler(api huma.API, usecase usecases.ActivityUsecase, mid middlewares.Middleware) {
	handler := &activityHandler{
		api:     api,
		usecase: usecase,
		mid:     mid,
	}

	huma.Get(api, "", handler.ListActivities, func(o *huma.Operation) {
		o.Summary = "List activities"
		o.Description = "Retrieve a list of activities with optional search, filtering, and sorting."
		o.DefaultStatus = 200
	})

	huma.Get(api, "/{id}", handler.GetActivity, func(o *huma.Operation) {
		o.Summary = "Get activity details"
		o.Description = "Retrieve detailed information about a specific activity."
		o.DefaultStatus = 200
	})
}

type ListActivitiesRequest struct {
	Search       string `query:"search" doc:"Search by title, description, or location"`
	HidePast     bool   `query:"hide_past" doc:"Exclude activities that have already ended" default:"false"`
	HappeningNow bool   `query:"happening_now" doc:"Include only activities currently in progress" default:"false"`
	SortBy       string `query:"sort_by" enum:"start_time,title,location" default:"start_time" doc:"Sort results by field"`
	Order        string `query:"order" enum:"asc,desc" default:"asc" doc:"Sort order"`
}

type ListActivitiesResponse struct {
	Body ListActivitiesResponseBody `json:"body"`
}

type ListActivitiesResponseBody struct {
	Activities []ActivityItem `json:"activities"`
}

type ActivityItem struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	BuildingName string    `json:"building_name,omitempty"`
	Floor        string    `json:"floor,omitempty"`
	RoomName     string    `json:"room_name,omitempty"`
	Description  string    `json:"description"`
	Image        string    `json:"image,omitempty"`
	IsHappening  bool      `json:"is_happening"`
}

func (h *activityHandler) ListActivities(ctx context.Context, input *ListActivitiesRequest) (*ListActivitiesResponse, error) {
	filter := models.ActivityFilter{
		Search:       input.Search,
		HidePast:     input.HidePast,
		HappeningNow: input.HappeningNow,
		SortBy:       input.SortBy,
		Order:        input.Order,
	}

	activities, err := h.usecase.ListActivities(ctx, filter)
	if err != nil {
		return nil, ErrInternalServerError
	}

	now := time.Now()
	items := make([]ActivityItem, 0, len(activities))
	for _, a := range activities {
		items = append(items, ActivityItem{
			ID:           a.ID,
			Title:        a.Title,
			StartTime:    a.StartTime,
			EndTime:      a.EndTime,
			BuildingName: a.BuildingName,
			Floor:        a.Floor,
			RoomName:     a.RoomName,
			Description:  a.Description,
			Image:        a.Image,
			IsHappening:  now.After(a.StartTime) && now.Before(a.EndTime),
		})
	}

	return &ListActivitiesResponse{
		Body: ListActivitiesResponseBody{
			Activities: items,
		},
	}, nil
}

type GetActivityRequest struct {
	ID int64 `path:"id" doc:"Activity ID"`
}

type GetActivityResponse struct {
	Body ActivityItem `json:"body"`
}

func (h *activityHandler) GetActivity(ctx context.Context, input *GetActivityRequest) (*GetActivityResponse, error) {
	activity, err := h.usecase.GetActivity(ctx, input.ID)
	if err != nil {
		if err == repositories.ErrActivityNotFound {
			return nil, ErrActivityNotFound
		}
		return nil, ErrInternalServerError
	}

	now := time.Now()
	return &GetActivityResponse{
		Body: ActivityItem{
			ID:           activity.ID,
			Title:        activity.Title,
			StartTime:    activity.StartTime,
			EndTime:      activity.EndTime,
			BuildingName: activity.BuildingName,
			Floor:        activity.Floor,
			RoomName:     activity.RoomName,
			Description:  activity.Description,
			Image:        activity.Image,
			IsHappening:  now.After(activity.StartTime) && now.Before(activity.EndTime),
		},
	}, nil
}
