package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

const (
	activityInputFilename = "OPH Data - activity.csv"
	boothInputFilename    = "OPH Data - booth.csv"
	workshopInputFilename = "OPH Data - workshop.csv"

	activityOutputFilename = "postgres_activity.csv"
	boothOutputFilename    = "postgres_booth.csv"
	workshopOutputFilename = "postgres_workshop.csv"
)

const (
	TruncateCommand = `
TRUNCATE TABLE workshops RESTART IDENTITY CASCADE;
TRUNCATE TABLE activities RESTART IDENTITY CASCADE;
TRUNCATE TABLE booths RESTART IDENTITY CASCADE;
`
	CopyCommand = `
COPY activities (title, description, event_date, start_time, end_time, building_name, floor, room_name, image, link) FROM '/home/postgres_activity.csv' DELIMITER ',' CSV HEADER;
COPY booths (name, category, check_in_code) FROM '/home/postgres_booth.csv' DELIMITER ',' CSV HEADER;
COPY workshops (name, description, category, affiliation, event_date, start_time, end_time, location, total_seats, image) FROM '/home/postgres_workshop.csv' DELIMITER ',' CSV HEADER;
`
)

type ActivityInputRow struct {
	Title        string `csv:"title"`
	EventDate    string `csv:"event_date"`
	StartTime    string `csv:"start_time"`
	EndTime      string `csv:"end_time"`
	BuildingName string `csv:"building_name"`
	Floor        string `csv:"floor"`
	RoomName     string `csv:"room_name"`
	Description  string `csv:"description"`
	ImageUrl     string `csv:"image_url  (4:3)"`
	Link         string `csv:"link"`
}

type BoothInputRow struct {
	Name                  string `csv:"name"`
	Type                  string `csv:"type"`
	CheckInCodeWithPrefix string `csv:"check_in_code"`
}

type WorkshopInputRow struct {
	Name        string `csv:"name"`
	Description string `csv:"description"`
	Category    string `csv:"category"`
	Affiliation string `csv:"affiliation"`
	EventDate   string `csv:"event_date"`
	StartTime   string `csv:"start_time"`
	EndTime     string `csv:"end_time"`
	Location    string `csv:"location"`
	TotalSeats  string `csv:"total_seats"`
	ImageUrl    string `csv:"image_url (4:3)"`
}

func main() {
	fmt.Println("Starting data transformation...")

	transformActivityCsv()
	transformBoothCsv()
	transformWorkshopCsv()

	fmt.Println("All transformations complete.")
}

func transformActivityCsv() {
	inputFile, err := os.Open(activityInputFilename)
	if err != nil {
		fmt.Printf("Error opening activity input: %v\n", err)
		return
	}
	defer inputFile.Close()

	var input []ActivityInputRow
	if err := gocsv.UnmarshalFile(inputFile, &input); err != nil {
		fmt.Printf("Error unmarshaling activity input: %v\n", err)
		return
	}

	outputFile, err := os.Create(activityOutputFilename)
	if err != nil {
		fmt.Printf("Error creating activity output: %v\n", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	newHeader := []string{"title", "description", "event_date", "start_time", "end_time", "building_name", "floor", "room_name", "image", "link"}
	writer.Write(newHeader)

	for i, row := range input {
		fullStart := fmt.Sprintf("%s %s+07", row.EventDate, row.StartTime)
		fullEnd := fmt.Sprintf("%s %s+07", row.EventDate, row.EndTime)

		imagePath := ""
		if row.ImageUrl != "" {
			fileId, found := strings.CutPrefix(row.ImageUrl, "https://drive.google.com/file/d/")
			if !found {
				fmt.Printf("Invalid activity image url %q on line %d\n", row.ImageUrl, i+2)
				continue
			}
			fileId = strings.Split(fileId, "/")[0]
			imagePath = fmt.Sprintf("/activity/%s", fileId)
		}

		newRow := []string{
			row.Title,
			row.Description,
			row.EventDate,
			fullStart,
			fullEnd,
			row.BuildingName,
			row.Floor,
			row.RoomName,
			imagePath,
			row.Link,
		}
		writer.Write(newRow)
	}
	fmt.Println(" - Activity CSV transformed")
}

func transformBoothCsv() {
	inputFile, err := os.Open(boothInputFilename)
	if err != nil {
		fmt.Printf("Error opening booth input: %v\n", err)
		return
	}
	defer inputFile.Close()

	var input []BoothInputRow
	if err := gocsv.UnmarshalFile(inputFile, &input); err != nil {
		fmt.Printf("Error unmarshaling booth input: %v\n", err)
		return
	}

	outputFile, err := os.Create(boothOutputFilename)
	if err != nil {
		fmt.Printf("Error creating booth output: %v\n", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	newHeader := []string{"name", "category", "check_in_code"}
	writer.Write(newHeader)

	for i, row := range input {
		var category string
		switch row.Type {
		case "Club":
			category = string(models.BoothCategoryClub)
		case "Department":
			category = string(models.BoothCategoryDepartment)
		case "Innovation exhibition":
			category = string(models.BoothCategoryExhibition)
		default:
			fmt.Printf("Invalid booth type %s on line %d, row data is %v\n", row.Type, i+2, row)
			continue
		}

		if len(row.CheckInCodeWithPrefix) <= 2 || row.CheckInCodeWithPrefix[0:2] != "B-" {
			fmt.Printf("Invalid check_in_code %s on line %d, row data is %v\n", row.CheckInCodeWithPrefix, i+2, row)
			continue
		}
		checkInCode := row.CheckInCodeWithPrefix[2:]
		if err := uuid.Validate(checkInCode); err != nil {
			fmt.Printf("Invalid check_in_code uuid format %s on line %d, row data is %v (%v)\n", checkInCode, i+2, row, err)
			continue
		}

		newRow := []string{
			row.Name,
			category,
			checkInCode,
		}
		writer.Write(newRow)
	}
	fmt.Println(" - Booth CSV transformed")
}

func transformWorkshopCsv() {
	inputFile, err := os.Open(workshopInputFilename)
	if err != nil {
		fmt.Printf("Error opening workshop input: %v\n", err)
		return
	}
	defer inputFile.Close()

	var input []WorkshopInputRow
	if err := gocsv.UnmarshalFile(inputFile, &input); err != nil {
		fmt.Printf("Error unmarshaling workshop input: %v\n", err)
		return
	}

	outputFile, err := os.Create(workshopOutputFilename)
	if err != nil {
		fmt.Printf("Error creating workshop output: %v\n", err)
		return
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	newHeader := []string{"name", "description", "category", "affiliation", "event_date", "start_time", "end_time", "location", "total_seats", "image"}
	writer.Write(newHeader)

	for i, row := range input {
		fullStart := fmt.Sprintf("%s %s+07", row.EventDate, row.StartTime)
		fullEnd := fmt.Sprintf("%s %s+07", row.EventDate, row.EndTime)

		var category string
		switch row.Category {
		case "Club":
			category = string(models.WorkShopCategoryClub)
		case "Department":
			category = string(models.WorkShopCategoryDepartment)
		default:
			fmt.Printf("Invalid workshop category %s on line %d, row data is %v\n", row.Category, i+2, row)
			continue
		}

		imagePath := ""
		if row.ImageUrl != "" {
			fileId, found := strings.CutPrefix(row.ImageUrl, "https://drive.google.com/file/d/")
			if !found {
				fmt.Printf("Invalid workshop image url %q on line %d\n", row.ImageUrl, i+2)
				continue
			}
			fileId = strings.Split(fileId, "/")[0]
			imagePath = fmt.Sprintf("/workshop/%s", fileId)
		}

		newRow := []string{
			row.Name,
			row.Description,
			category,
			row.Affiliation,
			row.EventDate,
			fullStart,
			fullEnd,
			row.Location,
			row.TotalSeats,
			imagePath,
		}
		writer.Write(newRow)
	}
	fmt.Println(" - Workshop CSV transformed")
}
