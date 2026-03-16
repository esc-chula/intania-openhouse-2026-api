package models

import "time"

type StampType string

const (
	StampTypeWorkshop StampType = "workshop"
	StampTypeBooth    StampType = "booth"
)

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
