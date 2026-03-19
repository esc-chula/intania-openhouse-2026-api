package models

import "time"

type StampType string

const (
	StampTypeDepartment StampType = "department"
	StampTypeClub       StampType = "club"
	StampTypeExhibition StampType = "exhibition"
)

type StampPoster struct {
	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID     int64     `bun:"user_id"`
	Type       StampType `bun:"type"`
	IsRedeemed bool      `bun:"is_redeemed"`
}

type StampItem struct {
	ID          int64     `bun:"id"`
	Type        StampType `bun:"type"`
	Name        string    `bun:"name"`
	CheckedInAt time.Time `bun:"checked_in_at"`
}
