package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SeedData struct {
	Activities []models.Activity
	Workshops  []models.Workshop
	Booths     []models.Booth
}

var modelsToDelete = []any{
	(*models.User)(nil),
	(*models.Activity)(nil),
	(*models.Workshop)(nil),
	(*models.Booth)(nil),
}

const (
	MAX_DATA_ITEM = 100
	TOTAL_SEAT    = 100
)

func seedData(ctx context.Context, db *bun.DB) *SeedData {
	activities := make([]models.Activity, MAX_DATA_ITEM)
	for i := range activities {
		activities[i] = models.Activity{
			ID:           int64(i + 1),
			Title:        fmt.Sprintf("Activity %d", i),
			Description:  "Desc",
			StartTime:    time.Now(),
			EndTime:      time.Now().Add(time.Hour),
			EventDate:    "2026-03-28",
			BuildingName: utils.Ptr("ENG3"),
			Floor:        utils.Ptr("4"),
			RoomName:     utils.Ptr("409"),
			Image:        nil,
			Link:         nil,
		}
	}

	workshops := make([]models.Workshop, MAX_DATA_ITEM)
	for i := range workshops {
		category := models.WorkShopCategoryDepartment
		if rand.Intn(2) == 1 {
			category = models.WorkShopCategoryClub
		}

		workshops[i] = models.Workshop{
			ID:          int64(i + 1),
			Name:        fmt.Sprintf("Workshop %d", i),
			Description: "Desc",
			Category:    category,
			Affiliation: "Something",
			EventDate:   "2026-03-28",
			StartTime:   time.Now(),
			EndTime:     time.Now().Add(time.Hour),
			Location:    "ENG3 409",
			TotalSeats:  TOTAL_SEAT,
			CheckInCode: uuid.NewString(),
		}
	}

	booths := make([]models.Booth, MAX_DATA_ITEM)
	for i := range booths {
		categoryIndex := rand.Intn(3)
		var category models.BoothCategory
		switch categoryIndex {
		case 0:
			category = models.BoothCategoryDepartment
		case 1:
			category = models.BoothCategoryClub
		case 2:
			category = models.BoothCategoryExhibition
		}

		booths[i] = models.Booth{
			ID:          int64(i + 1),
			Name:        fmt.Sprintf("Booth %d", i),
			Category:    category,
			CheckInCode: uuid.NewString(),
		}
	}

	dataToInsert := []any{&activities, &workshops, &booths}

	err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for _, model := range modelsToDelete {
			if _, err := tx.NewTruncateTable().Model(model).Cascade().Exec(ctx); err != nil {
				return err
			}
		}

		for _, data := range dataToInsert {
			if _, err := tx.NewInsert().Model(data).Exec(ctx); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("failed to seed data: %v", err)
	}

	return &SeedData{
		Activities: activities,
		Workshops:  workshops,
		Booths:     booths,
	}
}

func ensureUserCreated(baseURL string, auth string) error {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/users", strings.NewReader(`{
	"first_name": "John",
	"last_name": "Doe",
	"gender": "male",
	"phone_number": "0123456789",
	"participant_type": "student",
	"attendance_dates": ["2026-03-28"],
	"interested_activities": [],
	"discovery_channel": [],
	"is_from_bangkok": true,
	"origin_location": "phra_nakhon",
	"transport_mode": "personal_car",
	"student_extra_attributes": {
	  "education_level": "high school",
	  "emergency_contact": "0123456789",
	  "interested_major": "",
	  "province": "",
	  "school_name": "",
	  "study_plan": "",
	  "tcas_rank": 0
	}
}`))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
