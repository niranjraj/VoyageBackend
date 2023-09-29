package ingestion

import (
	"encoding/json"
	"time"

	"hackathon/models"
	"net/http"
)

func GetAirports(responseFlightRoute []models.ResponseFlightRoute) ([]models.Airport, error) {
	// Fetch airport data
	airports, err := fetchAirportsFromURL("https://gist.githubusercontent.com/tdreyno/4278655/raw/7b0762c09b519f40397e4c3e100b097d861f5588/airports.json")
	if err != nil {
		return nil, err
	}

	// Create a map to store flight route data for quick access
	routeDataMap := make(map[string]models.ResponseFlightRoute)
	for _, route := range responseFlightRoute {
		routeDataMap[route.AirportFrom.IATA] = route
		routeDataMap[route.AirportTo.IATA] = route
	}

	// Update airport data based on flight route information
	for i, airport := range airports {
		route, found := routeDataMap[airport.IATA]
		airports[i].ModifiedAt = time.Now()
		if found {
			airports[i].NormalizedScore = route.AirportTo.NormalizedScore
			airports[i].Canceled = route.AirportTo.Canceled
			airports[i].Delayed15 = route.AirportTo.Delayed15
			airports[i].Delayed30 = route.AirportTo.Delayed30
			airports[i].Delayed45 = route.AirportTo.Delayed45
			airports[i].OnTime = route.AirportTo.OnTime
		}
	}

	return airports, nil
}

func fetchAirportsFromURL(url string) ([]models.Airport, error) {
	// Fetch and decode airport data from the URL
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var airports []models.Airport
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&airports); err != nil {
		return nil, err
	}

	return airports, nil
}
