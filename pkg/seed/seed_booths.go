package seed

import (
	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
)

func GetBoothSeedData() []models.Booth {
	return []models.Booth{
		// Department Category
		{Name: "Computer Engineering", Category: models.BoothCategoryDepartment, CheckInCode: "00000000-0000-0000-0000-000000000001"},
		{Name: "Electrical Engineering", Category: models.BoothCategoryDepartment, CheckInCode: "00000000-0000-0000-0000-000000000002"},
		{Name: "Mechanical Engineering", Category: models.BoothCategoryDepartment, CheckInCode: "00000000-0000-0000-0000-000000000003"},
		{Name: "Civil Engineering", Category: models.BoothCategoryDepartment, CheckInCode: "00000000-0000-0000-0000-000000000004"},
		{Name: "Chemical Engineering", Category: models.BoothCategoryDepartment, CheckInCode: "00000000-0000-0000-0000-000000000005"},
		{Name: "Industrial Engineering", Category: models.BoothCategoryDepartment, CheckInCode: "00000000-0000-0000-0000-000000000006"},

		// Club Category
		{Name: "Robotics Club", Category: models.BoothCategoryClub, CheckInCode: "10000000-0000-0000-0000-000000000001"},
		{Name: "AI Research Group", Category: models.BoothCategoryClub, CheckInCode: "10000000-0000-0000-0000-000000000002"},
		{Name: "Sustainable Energy Club", Category: models.BoothCategoryClub, CheckInCode: "10000000-0000-0000-0000-000000000003"},
		{Name: "Engineering Music Club", Category: models.BoothCategoryClub, CheckInCode: "10000000-0000-0000-0000-000000000004"},
		{Name: "Sports Engineering Club", Category: models.BoothCategoryClub, CheckInCode: "10000000-0000-0000-0000-000000000005"},

		// Exhibition Category
		{Name: "Future Transport Exhibition", Category: models.BoothCategoryExhibition, CheckInCode: "20000000-0000-0000-0000-000000000001"},
		{Name: "Smart City Showcase", Category: models.BoothCategoryExhibition, CheckInCode: "20000000-0000-0000-0000-000000000002"},
		{Name: "Space Tech Exhibit", Category: models.BoothCategoryExhibition, CheckInCode: "20000000-0000-0000-0000-000000000003"},
		{Name: "Medical Robotics Display", Category: models.BoothCategoryExhibition, CheckInCode: "20000000-0000-0000-0000-000000000004"},
		{Name: "Green Hydrogen Project", Category: models.BoothCategoryExhibition, CheckInCode: "20000000-0000-0000-0000-000000000005"},
	}
}
