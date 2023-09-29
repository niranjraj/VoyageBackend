package router

import (
	"hackathon/db"
	"hackathon/models"

	"github.com/gofiber/fiber/v2"
)

func CreateRouterGroup(app *fiber.App) {
	apiGroup := app.Group("/api")

	apiGroup.Get("/airline/:iataCode", getAirline)
	apiGroup.Get("/airport/:iataCode", getAirport)
	apiGroup.Post("/recommendations", getRecommendedRoute)
	apiGroup.Post("/detailrecommendations", getDetailRecommendations)
	apiGroup.Get("/airportautocomplete", airportAutoComplete)
	apiGroup.Post("/routes", getAllRoutesById)

}

func getAirline(c *fiber.Ctx) error {

	iataCode := c.Params("iataCode")
	airline, err := db.GetAirlineByIATACode(iataCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Airline not found",
		})
	}

	// Return the airline information as JSON response
	return c.JSON(airline)
}

func getAirport(c *fiber.Ctx) error {
	iataCode := c.Params("iataCode")
	airport, err := db.GetAirportByIATACode(iataCode)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Airline not found",
		})
	}

	// Return the airline information as JSON response
	return c.JSON(airport)
}

func getRecommendedRoute(c *fiber.Ctx) error {
	var req models.RecommendationRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	recommendedRoutes, err := db.GetRecommendedRoute(req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No route found",
		})
	}
	return c.JSON(recommendedRoutes)

}

func getDetailRecommendations(c *fiber.Ctx) error {
	var req models.MultiCityRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	recommendedRoutes, err := db.GetRecommendedDetailRoute(req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No route found",
		})
	}
	return c.JSON(recommendedRoutes)
}
func airportAutoComplete(c *fiber.Ctx) error {
	// Read the search_string query parameter from the URL
	req := c.Query("search_string")

	// Check if the search_string is empty
	if req == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Search string is empty",
		})
	}

	// Call your database function to find airports based on the search string
	airports, err := db.FindAirportsWithSearch(req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No route found",
		})
	}

	// Return the list of airports as JSON response
	return c.JSON(airports)
}

func getAllRoutesById(c *fiber.Ctx) error {
	var req models.RouteIdRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	routes, err := db.FindRoutesWithIds(req.RouteId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No route found",
		})
	}
	return c.JSON(routes)
}
