package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Activity struct {
	bun.BaseModel `bun:"table:activities,alias:act"`
	ID            int64     `bun:"id,pk,autoincrement" json:"id"`
	Title         string    `bun:"title,notnull" json:"title"`
	Description   string    `bun:"description,notnull" json:"description"`
	StartTime     time.Time `bun:"start_time,notnull" json:"start_time"`
	EndTime       time.Time `bun:"end_time,notnull" json:"end_time"`
	EventDate     time.Time `bun:"event_date,notnull" json:"event_date"`
	BuildingName  string    `bun:"building_name" json:"building_name,omitempty"`
	Floor         string    `bun:"floor" json:"floor,omitempty"`
	RoomName      string    `bun:"room_name" json:"room_name,omitempty"`
	Image         string    `bun:"image" json:"image,omitempty"`
	Link          string    `bun:"link" json:"link,omitempty"`
}

type ActivityFilter struct {
	Search       string
	HidePast     bool
	HappeningNow bool
	SortBy       string
	Order        string
}
