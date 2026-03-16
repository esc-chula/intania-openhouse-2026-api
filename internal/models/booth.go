package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Booth struct {
	bun.BaseModel `       bun:"table:booths,alias:bt"`
	ID            int64  `bun:"id,pk,autoincrement"    json:"id"`
	Name          string `bun:"name"                   json:"name"`
	CheckInCode   string `bun:"check_in_code,nullzero" json:"-"`
}

type BoothCheckIn struct {
	bun.BaseModel `bun:"table:booth_checkins,alias:btck"`
	ID            int64     `bun:"id,pk,autoincrement"             json:"id"`
	UserID        int64     `bun:"user_id"                         json:"user_id"`
	BoothID       int64     `bun:"booth_id"                        json:"booth_id"`
	CheckedInAt   time.Time `bun:"checked_in_at"                   json:"checked_in_at"`
}
