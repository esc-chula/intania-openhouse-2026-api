package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/esc-chula/intania-openhouse-2026-api/internal/models"
	"github.com/gocarina/gocsv"
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
	CopyActivityCommand = `COPY activities (title, description, start_time, end_time, building_name, floor, room_name, image) FROM '/home/postgres_activity.csv' DELIMITER ',' CSV HEADER;`
	CopyBoothCommand    = `COPY booths (name, category) FROM '/home/postgres_booth.csv' DELIMITER ',' CSV HEADER;`
	CopyWorkshopCommand = `COPY workshops (name, description, category, affiliation, event_date, start_time, end_time, location, total_seats, image) FROM '/home/postgres_workshop.csv' DELIMITER ',' CSV HEADER;`
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
	Name string `csv:"name"`
	Type string `csv:"type"`
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

	newHeader := []string{"title", "description", "start_time", "end_time", "building_name", "floor", "room_name", "image"}
	writer.Write(newHeader)

	for _, row := range input {
		fullStart := fmt.Sprintf("%s %s+07", row.EventDate, row.StartTime)
		fullEnd := fmt.Sprintf("%s %s+07", row.EventDate, row.EndTime)

		newRow := []string{
			row.Title,
			row.Description,
			fullStart,
			fullEnd,
			row.BuildingName,
			row.Floor,
			row.RoomName,
			row.ImageUrl,
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

	newHeader := []string{"name", "category"}
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
			fmt.Printf("Invalid booth type %s on line %d, row data is %v", row.Type, i+2, row)
			continue
		}

		newRow := []string{
			row.Name,
			category,
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
			fmt.Printf("Invalid workshop category %s on line %d, row data is %v", row.Category, i+2, row)
			continue
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
			row.ImageUrl,
		}
		writer.Write(newRow)
	}
	fmt.Println(" - Workshop CSV transformed")
}
