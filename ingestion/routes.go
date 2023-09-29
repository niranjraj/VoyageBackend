package ingestion

import (
	"encoding/json"
	"fmt"
	"hackathon/models"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

func ConvertResponseToFlightRoute(responseFlightRoute []models.ResponseFlightRoute) []models.FlightRoute {
	var flightRoutes []models.FlightRoute
	flightRouteCh := make(chan models.FlightRoute)
	var wg sync.WaitGroup

	for _, responseRoute := range responseFlightRoute {

		wg.Add(1)
		go func(route models.ResponseFlightRoute) {
			defer wg.Done()
			flightRoute := models.FlightRoute{

				Airline:        route.Airline.IATACode,
				RouteId:        route.RouteID,
				AirportFrom:    route.AirportFrom.IATACode,
				AirportTo:      route.AirportTo.IATACode,
				AirportVia:     "", // Initialize AirportVia
				ClassBusiness:  route.ClassBusiness == 1,
				ClassEconomy:   route.ClassEconomy == 1,
				ClassFirst:     route.ClassFirst == 1,
				CommonDuration: route.CommonDuration,
				Monday:         route.Day1,
				Tuesday:        route.Day2,
				Wednesday:      route.Day3,
				Thursday:       route.Day4,
				Friday:         route.Day5,
				Saturday:       route.Day6,
				Sunday:         route.Day7,
				FlightsPerDay:  route.FlightsPerDay,
				FlightsPerWeek: route.FlightsPerWeek,
				IsActive:       route.IsActive == 1,
				MaxDuration:    route.MaxDuration,
				MinDuration:    route.MinDuration,
				ModifiedAt:     time.Now(),
			}

			// Check if AirportVia is not nil before accessing its IATA field
			if route.AirportVia.IATA != "" {
				flightRoute.AirportVia = route.AirportVia.IATA
			}

			flightRouteCh <- flightRoute
		}(responseRoute)
	}

	go func() {
		defer close(flightRouteCh)
		wg.Wait()

	}()

	for flightRoute := range flightRouteCh {
		flightRoutes = append(flightRoutes, flightRoute)
	}

	return flightRoutes
}

func ConvertResponseToFlightRoute2(responseFlightRoute []models.ResponseFlightRoute) []models.FlightRoute {
	var flightRoutes []models.FlightRoute
	var existingRoutes = make(map[int]struct{})

	for _, responseRoute := range responseFlightRoute {
		uniqueID := responseRoute.RouteID

		// Check if the route is already in existingRoutes
		_, exists := existingRoutes[uniqueID]
		if exists {
			continue // Skip if the route already exists
		}
		fmt.Printf("%d", responseRoute.RouteID)

		flightRoute := models.FlightRoute{

			Airline:        responseRoute.Airline.IATACode,
			RouteId:        responseRoute.RouteID,
			AirportFrom:    responseRoute.AirportFrom.IATACode,
			AirportTo:      responseRoute.AirportTo.IATACode,
			AirportVia:     "", // Initialize AirportVia
			ClassBusiness:  responseRoute.ClassBusiness == 1,
			ClassEconomy:   responseRoute.ClassEconomy == 1,
			ClassFirst:     responseRoute.ClassFirst == 1,
			CommonDuration: responseRoute.CommonDuration,
			Monday:         responseRoute.Day1,
			Tuesday:        responseRoute.Day2,
			Wednesday:      responseRoute.Day3,
			Thursday:       responseRoute.Day4,
			Friday:         responseRoute.Day5,
			Saturday:       responseRoute.Day6,
			Sunday:         responseRoute.Day7,
			FlightsPerDay:  responseRoute.FlightsPerDay,
			FlightsPerWeek: responseRoute.FlightsPerWeek,
			IsActive:       responseRoute.IsActive == 1,
			MaxDuration:    responseRoute.MaxDuration,
			MinDuration:    responseRoute.MinDuration,
			ModifiedAt:     time.Now(),
		}
		// Update the existingRoutes map
		existingRoutes[responseRoute.RouteID] = struct{}{}

		flightRoutes = append(flightRoutes, flightRoute)
	}

	return flightRoutes
}

func FetchFlightRoutesFromLocalFile(filename string) ([]models.ResponseFlightRoute, error) {
	// Fetch and decode flight route data from a local file

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var responseFlightRoute []models.ResponseFlightRoute
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&responseFlightRoute); err != nil {
		return nil, err
	}

	return responseFlightRoute, nil
}

func FetchRoutesFromUrl(airports []models.Airport) []models.ResponseFlightRoute {
	var responseFlightRoute []models.ResponseFlightRoute
	baseURL := "https://www.flightroutes.com/api/loadMore/newroutes/"

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the HTTP requests
	}

	for _, airport := range airports {
		airportURL := baseURL + airport.IATA

		for page := 1; ; page++ {
			u, err := url.Parse(airportURL)
			if err != nil {
				fmt.Print("here")
				// Log the error and continue
				fmt.Println("Error parsing URL:", err)
				continue
			}

			q := u.Query()
			q.Add("page", strconv.Itoa(page))
			u.RawQuery = q.Encode()

			response, err := client.Get(u.String())
			if err != nil {
				fmt.Print("here")
				// Log the error and continue
				fmt.Println("HTTP request error:", err)
				continue
			}
			defer response.Body.Close()

			// Handle 403 Forbidden error
			if response.StatusCode == http.StatusForbidden {
				fmt.Println("HTTP request is Forbidden. Skipping...")
				break // Break out of the loop for this airport
			}

			// Handle EOF error
			if err != nil && err.Error() == "EOF" {
				fmt.Println("EOF error occurred. Skipping...")
				break // Break out of the loop for this airport
			}

			if response.StatusCode != http.StatusOK {
				fmt.Println("HTTP request failed with status code:", response.StatusCode)
				break // Break out of the loop for this airport
			}

			var routeResponse models.RouteFilterResponse
			decoder := json.NewDecoder(response.Body)
			if err := decoder.Decode(&routeResponse); err != nil {
				fmt.Println("JSON decoding error:", err)
				break // Break out of the loop for this airport
			}

			responseFlightRoute = append(responseFlightRoute, routeResponse.Data...)
		}
	}

	return responseFlightRoute
}
