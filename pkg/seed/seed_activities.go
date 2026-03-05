package seed

import (
	"time"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
)

func GetActivitySeedData() []models.Activity {
	return []models.Activity{
		{
			Title:        "Opening Ceremony",
			Description:  "Join us for the grand opening of Intania Openhouse 2026.",
			StartTime:    time.Date(2026, 3, 28, 9, 0, 0, 0, time.Local),
			EndTime:      time.Date(2026, 3, 28, 10, 0, 0, 0, time.Local),
			BuildingName: "Engineering Building 3",
			Floor:        "1",
			RoomName:     "Hall of Intania",
			Image:        "https://example.com/opening.jpg",
		},
		{
			Title:        "Robotics Workshop",
			Description:  "Learn how to build and program your first robot.",
			StartTime:    time.Date(2026, 3, 28, 13, 0, 0, 0, time.Local),
			EndTime:      time.Date(2026, 3, 28, 15, 0, 0, 0, time.Local),
			BuildingName: "Engineering Building 4",
			Floor:        "G",
			RoomName:     "Robotics Lab",
			Image:        "https://example.com/robotics.jpg",
		},
		{
			Title:        "Engineering Fair",
			Description:  "Explore projects and innovations from various departments.",
			StartTime:    time.Date(2026, 3, 28, 10, 30, 0, 0, time.Local),
			EndTime:      time.Date(2026, 3, 28, 16, 0, 0, 0, time.Local),
			BuildingName: "Engineering Library",
			Floor:        "1-2",
			RoomName:     "Main Hall",
			Image:        "https://example.com/fair.jpg",
		},
		// Past Event
		{
			Title:        "Orientation for Volunteers",
			Description:  "Preparation for the volunteers of Openhouse 2026.",
			StartTime:    time.Date(2026, 2, 1, 9, 0, 0, 0, time.Local),
			EndTime:      time.Date(2026, 2, 1, 12, 0, 0, 0, time.Local),
			BuildingName: "Engineering Building 3",
			Floor:        "1",
			RoomName:     "Hall of Intania",
			Image:        "https://example.com/orientation.jpg",
		},
		// Happening Now
		{
			Title:        "Midnight Hackathon Setup",
			Description:  "Setting up the equipment for the overnight hackathon.",
			StartTime:    time.Date(2026, 3, 4, 22, 0, 0, 0, time.Local),
			EndTime:      time.Date(2026, 4, 5, 2, 0, 0, 0, time.Local),
			BuildingName: "Engineering Building 100",
			Floor:        "3",
			RoomName:     "Tech Hub",
			Image:        "https://example.com/setup.jpg",
		},
		// Future Event
		{
			Title:        "Final Props Inspection",
			Description:  "Checking all physical assets before the big day.",
			StartTime:    time.Date(2026, 4, 15, 14, 0, 0, 0, time.Local),
			EndTime:      time.Date(2026, 4, 15, 16, 0, 0, 0, time.Local),
			BuildingName: "Engineering Building 4",
			Floor:        "G",
			RoomName:     "Storage Site X",
			Image:        "https://example.com/inspection.jpg",
		},
	}
}
