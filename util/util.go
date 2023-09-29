package util

import (
	"fmt"
	"hackathon/models"
	"time"
)

func SplitIntoBatches(data []models.FlightRoute, batchSize int) [][]interface{} {
	var batches [][]interface{}

	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := make([]interface{}, end-i)
		for j, item := range data[i:end] {
			batch[j] = item
		}
		batches = append(batches, batch)
	}
	return batches
}

func GetWeekDay(dateString string) (string, error) {
	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "", err
	}

	// Convert the weekday to a string
	weekdayStr := date.Weekday().String()

	return weekdayStr, nil
}

// func FormattedRoute(route  models.FlightRoute){
// 	return models.UpdatedRouteResponse{

// 	}

// }
