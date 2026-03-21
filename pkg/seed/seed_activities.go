package seed

import (
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/esc-chula/intania-openhouse-2026-api/pkg/utils"
)

func GetActivitySeedData() []models.Activity {
	return []models.Activity{
		{
			Title:        "Opening Ceremony",
			Description:  "Join us for the grand opening of Intania Openhouse 2026.",
			StartTime:    time.Date(1, 1, 1, 9, 0, 0, 0, time.UTC),
			EndTime:      time.Date(1, 1, 1, 10, 0, 0, 0, time.UTC),
			BuildingName: utils.Ptr("Engineering Building 3"),
			Floor:        utils.Ptr("1"),
			RoomName:     utils.Ptr("Hall of Intania"),
			Image:        utils.Ptr("https://example.com/opening.jpg"),
			EventDate:    "2026-03-28",
			Link:         utils.Ptr("https://example.com/events/opening-ceremony"),
		},
		{
			Title:        "Robotics Workshop",
			Description:  "Learn how to build and program your first robot.",
			StartTime:    time.Date(1, 1, 1, 13, 0, 0, 0, time.UTC),
			EndTime:      time.Date(1, 1, 1, 15, 0, 0, 0, time.UTC),
			BuildingName: utils.Ptr("Engineering Building 4"),
			Floor:        utils.Ptr("G"),
			RoomName:     utils.Ptr("Robotics Lab"),
			Image:        utils.Ptr("https://example.com/robotics.jpg"),
			EventDate:    "2026-03-28",
			Link:         utils.Ptr("https://example.com/events/robotics-workshop"),
		},
		{
			Title:        "Engineering Fair",
			Description:  "Explore projects and innovations from various departments.",
			StartTime:    time.Date(1, 1, 1, 10, 30, 0, 0, time.UTC),
			EndTime:      time.Date(1, 1, 1, 16, 0, 0, 0, time.UTC),
			BuildingName: utils.Ptr("Engineering Library"),
			Floor:        utils.Ptr("1-2"),
			RoomName:     utils.Ptr("Main Hall"),
			Image:        utils.Ptr("https://example.com/fair.jpg"),
			EventDate:    "2026-03-28",
			Link:         utils.Ptr("https://example.com/events/engineering-fair"),
		},
		// Past Event
		{
			Title:        "Orientation for Volunteers",
			Description:  "Preparation for the volunteers of Openhouse 2026.",
			StartTime:    time.Date(1, 1, 1, 9, 0, 0, 0, time.UTC),
			EndTime:      time.Date(1, 1, 1, 12, 0, 0, 0, time.UTC),
			BuildingName: utils.Ptr("Engineering Building 3"),
			Floor:        utils.Ptr("1"),
			RoomName:     utils.Ptr("Hall of Intania"),
			Image:        utils.Ptr("https://example.com/orientation.jpg"),
			EventDate:    "2026-02-01",
			Link:         utils.Ptr("https://example.com/events/volunteer-orientation"),
		},
		// Happening Now
		{
			Title:        "Midnight Hackathon Setup",
			Description:  "Setting up the equipment for the overnight hackathon.",
			StartTime:    time.Date(1, 1, 1, 22, 0, 0, 0, time.UTC),
			EndTime:      time.Date(1, 1, 1, 23, 59, 0, 0, time.UTC),
			BuildingName: utils.Ptr("Engineering Building 100"),
			Floor:        utils.Ptr("3"),
			RoomName:     utils.Ptr("Tech Hub"),
			Image:        utils.Ptr("https://example.com/setup.jpg"),
			EventDate:    "2026-03-04",
			Link:         utils.Ptr("https://example.com/events/hackathon-setup"),
		},
		// Future Event
		{
			Title:        "Final Props Inspection",
			Description:  "Checking all physical assets before the big day.",
			StartTime:    time.Date(1, 1, 1, 14, 0, 0, 0, time.UTC),
			EndTime:      time.Date(1, 1, 1, 16, 0, 0, 0, time.UTC),
			BuildingName: utils.Ptr("Engineering Building 4"),
			Floor:        utils.Ptr("G"),
			RoomName:     utils.Ptr("Storage Site X"),
			Image:        utils.Ptr("https://example.com/inspection.jpg"),
			EventDate:    "2026-04-15",
			Link:         utils.Ptr("https://example.com/events/props-inspection"),
		},
	}
}
