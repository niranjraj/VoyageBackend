package ingestion

import (
	"encoding/json"
	"hackathon/models"
	"net/http"
	"time"
)

func GetAirlines(responseFlightRoute []models.ResponseFlightRoute) ([]models.Airline, error) {
	// Specify the URL you want to send the GET request to

	airlines, err := fetchAirlinesFromURL("https://cdn.jsdelivr.net/gh/besrourms/airlines@latest/airlines.json")
	if err != nil {
		return nil, err
	}

	routeDataMap := make(map[string]models.ResponseFlightRoute)
	for _, route := range responseFlightRoute {
		routeDataMap[route.Airline.IATA] = route

	}

	for i, airline := range airlines {
		route, found := routeDataMap[airline.IATA]

		airlines[i].ModifiedAt = time.Now()

		if found {
			airlines[i].IsScheduledPassenger = route.Airline.IsScheduledPassenger
			airlines[i].IsCargo = route.Airline.IsCargo

		}
	}

	return airlines, nil

}

func fetchAirlinesFromURL(url string) ([]models.Airline, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var airlines []models.Airline
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&airlines); err != nil {
		return nil, err
	}

	return airlines, nil
}
