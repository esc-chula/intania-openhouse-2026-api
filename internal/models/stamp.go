package models

import "time"

type StampType string

const (
	StampTypeDepartment StampType = "Department"
	StampTypeClub       StampType = "Club"
	StampTypeExhibition StampType = "Exhibition"
)

type StampPoster struct {
	ID         int64     `bun:"id,pk,autoincrement"`
	UserID     int64     `bun:"user_id"`
	Type       StampType `bun:"stamp_type"`
	IsRedeemed bool      `bun:"is_redeemed"`
}

type StampItem struct {
	ID          int64     `bun:"id"`
	Type        StampType `bun:"-"`
	Name        string    `bun:"name"`
	CheckedInAt time.Time `bun:"checked_in_at"`
}

type UserStamps struct {
	TotalCount int64
	Stamps     []StampItem
}
