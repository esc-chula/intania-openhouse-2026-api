package models

import (
	"github.com/uptrace/bun"
	"time"
)

type Category string

const (
	CategoryDepartment Category = "Department"
	CategoryClub       Category = "Club"
)

type Workshop struct {
	bun.BaseModel   `bun:"table:workshops,alias:ws"`
	ID              int64     `bun:"id,pk,autoincrement" json:"id"`
	Name            string    `bun:"name" json:"name"`
	Description     string    `bun:"description" json:"description"`
	Category        Category  `bun:"category" json:"category"`
	Affiliation     string    `bun:"affiliation" json:"affiliation"`
	EventDate       string    `bun:"event_date" json:"event_date"` // Date in format `2024-12-31`
	StartTime       time.Time `bun:"start_time" json:"start_time"`
	EndTime         time.Time `bun:"end_time" json:"end_time"`
	Location        string    `bun:"location" json:"location"`
	TotalSeats      int       `bun:"total_seats" json:"total_seats"`
	RegisteredCount int       `bun:"registered_count" json:"registered_count"`
}

type WorkshopFilter struct {
	Search    string
	Category  string
	EventDate string
	HideFull  bool
	SortBy    string // "start_time" | "name"
	Order     string // "ASC" | "DESC"
}

type Status string

const (
	StatusConfirmed Status = "Confirmed"
	StatusCancelled Status = "Cancelled"
	StatusAttended  Status = "attended"
	StatusAbsent    Status = "absent"
)

type Booking struct {
	bun.BaseModel `bun:"table:bookings,alias:bk"`
	ID            int64      `bun:"id,pk,autoincrement" json:"id"`
	UserID        int64      `bun:"user_id" json:"user_id"`
	WorkshopID    int64      `bun:"workshop_id" json:"workshop_id"`
	Status        Status     `bun:"status" json:"status"`
	CreatedAt     time.Time  `bun:"created_at,nullzero" json:"created_at"`
	CheckedInAt   *time.Time `bun:"checked_in_at,nullzero" json:"checked_in_at"`
}

// BookingWithTime is used for time-overlap checking.
// Returned by joining bookings with workshops.
type BookingWithTime struct {
	bun.BaseModel `bun:"table:bookings,alias:bk"`
	BookingID     int64     `bun:"id" json:"booking_id"`
	WorkshopID    int64     `bun:"workshop_id" json:"workshop_id"`
	StartTime     time.Time `bun:"start_time" json:"start_time"`
	EndTime       time.Time `bun:"end_time" json:"end_time"`
}
