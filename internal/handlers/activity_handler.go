package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/middlewares"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/repositories"
	"github.com/esc-chula/intania-openhouse-2026-api/internal/usecases"
)

var ErrActivityNotFound = huma.Error404NotFound("activity not found")
var loc, _ = time.LoadLocation("Asia/Bangkok")
var now = time.Now().In(loc)

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

	activityTag := "activity"

	huma.Get(api, "", handler.ListActivities, func(o *huma.Operation) {
		o.Summary = "List activities"
		o.Description = "Retrieve a list of activities with optional search, filtering, and sorting."
		o.DefaultStatus = 200
		o.Tags = []string{activityTag}
	})

	huma.Get(api, "/{id}", handler.GetActivity, func(o *huma.Operation) {
		o.Summary = "Get activity details"
		o.Description = "Retrieve detailed information about a specific activity."
		o.DefaultStatus = 200
		o.Tags = []string{activityTag}
	})
}

type ListActivitiesRequest struct {
	Search       string `query:"search"        doc:"Search by title, description, or location"`
	HidePast     bool   `query:"hide_past"     doc:"Exclude activities that have already ended"    default:"false"`
	HappeningNow bool   `query:"happening_now" doc:"Include only activities currently in progress" default:"false"`
	SortBy       string `query:"sort_by"       doc:"Sort results by field"                         default:"start_time" enum:"start_time,title,location"`
	Order        string `query:"order"         doc:"Sort order"                                    default:"asc"        enum:"asc,desc"`
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
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	EventDate    string    `json:"event_date"`
	BuildingName *string   `json:"building_name"`
	Floor        *string   `json:"floor"`
	RoomName     *string   `json:"room_name"`
	Image        *string   `json:"image"`
	Link         *string   `json:"link"`
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
		return nil, ErrInternalServerError(err)
	}

	items := make([]ActivityItem, 0, len(activities))
	for _, a := range activities {
		isHappening, err := getIsHappening(now, a.EventDate, a.StartTime, a.EndTime)
		if err != nil {
			return nil, ErrInternalServerError()
		}

		items = append(items, ActivityItem{
			ID:           a.ID,
			Title:        a.Title,
			StartTime:    a.StartTime,
			EndTime:      a.EndTime,
			EventDate:    a.EventDate,
			BuildingName: a.BuildingName,
			Floor:        a.Floor,
			RoomName:     a.RoomName,
			Description:  a.Description,
			Image:        a.Image,
			IsHappening:  isHappening,
			Link:         a.Link,
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
		return nil, ErrInternalServerError(err)
	}

	isHappening, err := getIsHappening(now, activity.EventDate, activity.StartTime, activity.EndTime)
	if err != nil {
		return nil, ErrInternalServerError()
	}

	return &GetActivityResponse{
		Body: ActivityItem{
			ID:           activity.ID,
			Title:        activity.Title,
			StartTime:    activity.StartTime,
			EndTime:      activity.EndTime,
			EventDate:    activity.EventDate,
			BuildingName: activity.BuildingName,
			Floor:        activity.Floor,
			RoomName:     activity.RoomName,
			Description:  activity.Description,
			Image:        activity.Image,
			IsHappening:  isHappening,
			Link:         activity.Link,
		},
	}, nil
}

func getIsHappening(now time.Time, eventDate string, startTime time.Time, endTime time.Time) (bool, error) {
	datePart, err := time.Parse(time.RFC3339, eventDate)
	if err != nil {
		return false, fmt.Errorf("invalid event date: %w", err)
	}

	y, m, d := datePart.Date()

	fullStartTime := time.Date(y, m, d, startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Nanosecond(), now.Location())
	fullEndTime := time.Date(y, m, d, endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Nanosecond(), now.Location())
	return now.After(fullStartTime) && now.Before(fullEndTime), nil
}
