package db

import (
	"context"
	"errors"
	"fmt"
	"hackathon/models"
	"hackathon/util"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var classMappings = map[string]string{
	"ECO":   "class_economy",
	"BUS":   "class_business",
	"FIRST": "class_first",
}

func InsertBatchRoutes(routes []models.FlightRoute) error {
	fmt.Print(routes)
	collection := getCollection("routes")

	batches := util.SplitIntoBatches(routes, 1000)

	for _, batch := range batches {

		_, err := collection.InsertMany(context.TODO(), batch)
		if err != nil {
			return err
		}
	}

	return err
}

func GetRecommendedRoute(request models.RecommendationRequest) ([]models.FlightRoute, error) {
	collection := getCollection("routes")
	dateString := request.Date

	// Parse the date string into a time.Time object
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return nil, err
	}

	dayOfWeekNumber := int(date.Weekday())
	fmt.Printf("%d", dayOfWeekNumber)
	updatedDay := dayOfWeekNumber
	if dayOfWeekNumber == 0 {
		updatedDay = 7
	}

	day := fmt.Sprintf("day%d", updatedDay)
	fmt.Printf("%s", day)
	class := strings.ToUpper(request.Class)
	switch class {
	case "ECONOMY":
		// Handle economy class logic
		class = fmt.Sprintf("class_economy")
	case "BUSINESS":
		// Handle business class logic
		class = fmt.Sprintf("class_business")
	case "FIRST":
		// Handle first class logic
		class = fmt.Sprintf("class_first")
	default:
		// Handle other class values or errors
		return nil, errors.New("Not a valid class")
	}
	filter := bson.M{
		"airportFrom": request.Departure,
		"airportTo":   request.Arrival,
		day:           "yes",
		class:         1,
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	uniqueRoutes := make(map[int]models.FlightRoute)

	// Iterate through the cursor and add unique routes to the map
	for cursor.Next(context.TODO()) {
		var route models.FlightRoute
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		uniqueRoutes[route.RouteId] = route
	}

	var routes []models.FlightRoute
	for _, route := range uniqueRoutes {
		routes = append(routes, route)
	}

	return routes, nil

}

func GetRecommendedDetailRoute3(request models.MultiCityRequest) (models.MultiCityResponse, error) {

	collection := getCollection("routes")
	classMappings := map[string]string{
		"ECO":   "class_economy",
		"BUS":   "class_business",
		"FIRST": "class_first",
	}
	var classFilters []bson.M

	for _, class := range request.Class {
		fieldName, exists := classMappings[class]
		if exists {
			classFilters = append(classFilters, bson.M{fieldName: true})
		}
	}

	var routeResponse []models.RouteResponse

	for _, route := range request.Routes {
		filter := bson.M{
			"airportFrom": route.Departure,
			"airportTo":   route.Arrival,
		}
		if len(request.Class) > 0 {
			filter["$or"] = classFilters
		}
		if len(request.Airlines) > 0 {
			filter["airline"] = bson.M{"$in": request.Airlines}
		}
		cursor, err := collection.Find(context.TODO(), filter)
		if err != nil {
			return models.MultiCityResponse{}, err
		}
		defer cursor.Close(context.TODO())
		var routeResults []models.FlightRoute
		uniqueRoutes := make(map[int]models.FlightRoute)

		for cursor.Next(context.TODO()) {
			var routeResult models.FlightRoute
			if err := cursor.Decode(&routeResult); err != nil {
				return models.MultiCityResponse{}, err
			}
			uniqueRoutes[routeResult.RouteId] = routeResult

		}

		for _, route := range uniqueRoutes {
			routeResults = append(routeResults, route)
		}

		responseRoute := models.RouteResponse{
			Date:      route.Date,
			Departure: route.Departure,
			Arrival:   route.Arrival,
			Route:     routeResults,
		}
		routeResponse = append(routeResponse, responseRoute)

	}

	return models.MultiCityResponse{
		MultiCity: routeResponse,
	}, nil

}
func GetRecommendedDetailRoute2(request models.MultiCityRequest) (models.MultiCityResponse2, error) {

	collection := getCollection("routes")

	classMappings := map[string]string{
		"ECO":   "class_economy",
		"BUS":   "class_business",
		"FIRST": "class_first",
	}
	var classFilters []bson.M

	for _, class := range request.Class {
		fieldName, exists := classMappings[class]
		if exists {
			classFilters = append(classFilters, bson.M{fieldName: true})
		}
	}

	var responseRoutes []models.RouteResponse2

	for _, route := range request.Routes {
		var pipeline []bson.M
		date, err := util.GetWeekDay(route.Date)
		if err != nil {
			return models.MultiCityResponse2{}, errors.New("invalid date format")
		}
		fmt.Print(date)
		matchFilters := []bson.M{
			{
				date:          "yes",
				"airportFrom": route.Departure,
				"airportTo":   route.Arrival,
			},
		}

		if len(request.Class) > 0 {
			classFilter := bson.M{
				"$or": classFilters,
			}
			matchFilters = append(matchFilters, classFilter)
		}

		if len(request.Airlines) > 0 {
			airlineFilter := bson.M{
				"airline": bson.M{
					"$in": request.Airlines,
				},
			}
			matchFilters = append(matchFilters, airlineFilter)
		}

		matchStage := bson.M{
			"$match": bson.M{
				"$and": matchFilters,
			},
		}

		pipeline = append(pipeline, matchStage)

		lookupStages := []bson.M{
			{
				"$lookup": bson.M{
					"from":         "airports",
					"localField":   "airportFrom",
					"foreignField": "iata",
					"as":           "resAirportFrom",
				},
			},
			{
				"$lookup": bson.M{
					"from":         "airports",
					"localField":   "airportTo",
					"foreignField": "iata",
					"as":           "resAirportTo",
				},
			},
			{
				"$lookup": bson.M{
					"from":         "airlines",
					"localField":   "airline",
					"foreignField": "iata",
					"as":           "resAirline",
				},
			},
		}

		pipeline = append(pipeline, lookupStages...)

		unwindStages := []bson.M{
			{
				"$unwind": bson.M{
					"path":                       "$resAirportFrom",
					"preserveNullAndEmptyArrays": true,
				},
			},
			{
				"$unwind": bson.M{
					"path":                       "$resAirportTo",
					"preserveNullAndEmptyArrays": true,
				},
			},
			{
				"$unwind": bson.M{
					"path":                       "$resAirline",
					"preserveNullAndEmptyArrays": true,
				},
			},
		}

		pipeline = append(pipeline, unwindStages...)

		cursor, err := collection.Aggregate(context.TODO(), pipeline)
		if err != nil {
			fmt.Printf("cursor error %s", err)
			return models.MultiCityResponse2{}, err
		}
		defer cursor.Close(context.TODO())

		var routeResponse []models.RouteInfo
		fmt.Printf("from api:%s", route.Departure)
		for cursor.Next(context.TODO()) {
			var routeResult models.RouteInfo
			if err := cursor.Decode(&routeResult); err != nil {
				fmt.Print("decode error")
				return models.MultiCityResponse2{}, err
			}

			// Adapt the routeResult data to match your desired response structure

			routeResponse = append(routeResponse, routeResult)
		}

		responseRoute := models.RouteResponse2{
			Date:      route.Date,
			Departure: route.Departure,
			Arrival:   route.Arrival,
			Route:     routeResponse,
		}
		responseRoutes = append(responseRoutes, responseRoute)

	}

	return models.MultiCityResponse2{
		MultiCity: responseRoutes,
	}, nil
}

func GetRecommendedDetailRoute(request models.MultiCityRequest) (models.MultiCityResponse2, error) {
	// ...
	collection := getCollection("routes")
	var responseRoutes []models.RouteResponse2
	var wg sync.WaitGroup
	routeChannel := make(chan models.RouteResponse2, len(request.Routes))

	classMappings := map[string]string{
		"ECO":   "class_economy",
		"BUS":   "class_business",
		"FIRST": "class_first",
	}
	var classFilters []bson.M

	for _, class := range request.Class {
		fieldName, exists := classMappings[class]
		if exists {
			classFilters = append(classFilters, bson.M{fieldName: true})
		}
	}
	if len(request.Airlines) > 0 {
		airlineFilter := bson.M{
			"airline": bson.M{
				"$in": request.Airlines,
			},
		}
		classFilters = append(classFilters, airlineFilter)
	}
	if len(request.Class) > 0 {
		classFilter := bson.M{
			"$or": classFilters,
		}
		classFilters = append(classFilters, classFilter)
	}
	for _, route := range request.Routes {
		wg.Add(1)
		go func(route models.DetailRecommendationRequest) {
			defer wg.Done()
			// Process the route and send the result to the channel
			response, err := processRoute(route, collection, request, classFilters)
			if err != nil {
				return
			}
			routeChannel <- response
		}(route)
	}

	// Close the channel when all Goroutines have finished
	go func() {
		wg.Wait()
		close(routeChannel)
	}()

	for response := range routeChannel {
		responseRoutes = append(responseRoutes, response)
	}

	return models.MultiCityResponse2{
		MultiCity: responseRoutes,
	}, nil
}

func processRoute(route models.DetailRecommendationRequest, collection *mongo.Collection, request models.MultiCityRequest, classFilters []bson.M) (models.RouteResponse2, error) {
	getCollection("routes")
	var pipeline []bson.M
	date, err := util.GetWeekDay(route.Date)
	if err != nil {
		return models.RouteResponse2{}, errors.New("invalid date format")
	}
	fmt.Print(date)

	matchFilters := []bson.M{
		{
			date:          "yes",
			"airportFrom": route.Departure,
			"airportTo":   route.Arrival,
		},
	}

	matchFilters = append(matchFilters, classFilters...)

	matchStage := bson.M{
		"$match": bson.M{
			"$and": matchFilters,
		},
	}

	pipeline = append(pipeline, matchStage)

	lookupStages := []bson.M{
		{
			"$lookup": bson.M{
				"from":         "airports",
				"localField":   "airportFrom",
				"foreignField": "iata",
				"as":           "resAirportFrom",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "airports",
				"localField":   "airportTo",
				"foreignField": "iata",
				"as":           "resAirportTo",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "airlines",
				"localField":   "airline",
				"foreignField": "iata",
				"as":           "resAirline",
			},
		},
	}

	pipeline = append(pipeline, lookupStages...)

	unwindStages := []bson.M{
		{
			"$unwind": bson.M{
				"path":                       "$resAirportFrom",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$resAirportTo",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$resAirline",
				"preserveNullAndEmptyArrays": true,
			},
		},
	}

	pipeline = append(pipeline, unwindStages...)

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		fmt.Printf("cursor error %s", err)
		return models.RouteResponse2{}, err
	}
	defer cursor.Close(context.TODO())

	var routeResponse []models.RouteInfo
	fmt.Printf("from api:%s", route.Departure)
	for cursor.Next(context.TODO()) {
		var routeResult models.RouteInfo
		if err := cursor.Decode(&routeResult); err != nil {
			fmt.Print("decode error")
			return models.RouteResponse2{}, err
		}

		// Adapt the routeResult data to match your desired response structure

		routeResponse = append(routeResponse, routeResult)
	}

	responseRoute := models.RouteResponse2{
		Date:      route.Date,
		Departure: route.Departure,
		Arrival:   route.Arrival,
		Route:     routeResponse,
	}
	return responseRoute, nil
}

func FindRoutesWithId(routeId []int) ([]models.ResponseFlightRoute2, error) {
	collection := getCollection("routes")
	var routes []models.ResponseFlightRoute2
	uniqueRoutes := make(map[int]models.ResponseFlightRoute2)
	for _, id := range routeId {
		filter := bson.A{
			bson.D{{Key: "$match", Value: bson.D{{Key: "routeId", Value: id}}}},
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "airports"},
						{Key: "localField", Value: "airportFrom"},
						{Key: "foreignField", Value: "iata"},
						{Key: "as", Value: "resAirportFrom"},
					},
				},
			},
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "airports"},
						{Key: "localField", Value: "airportTo"},
						{Key: "foreignField", Value: "iata"},
						{Key: "as", Value: "resAirportTo"},
					},
				},
			},
			bson.D{
				{Key: "$lookup",
					Value: bson.D{
						{Key: "from", Value: "airlines"},
						{Key: "localField", Value: "airline"},
						{Key: "foreignField", Value: "iata"},
						{Key: "as", Value: "resAirline"},
					},
				},
			},
			bson.D{
				{Key: "$unwind",
					Value: bson.D{
						{Key: "path", Value: "$resAirportFrom"},
						{Key: "includeArrayIndex", Value: "string"},
						{Key: "preserveNullAndEmptyArrays", Value: true},
					},
				},
			},
			bson.D{
				{Key: "$unwind",
					Value: bson.D{
						{Key: "path", Value: "$resAirportTo"},
						{Key: "includeArrayIndex", Value: "string"},
						{Key: "preserveNullAndEmptyArrays", Value: true},
					},
				},
			},
			bson.D{
				{Key: "$unwind",
					Value: bson.D{
						{Key: "path", Value: "$resAirline"},
						{Key: "includeArrayIndex", Value: "string"},
						{Key: "preserveNullAndEmptyArrays", Value: true},
					},
				},
			},
		}

		cursor, err := collection.Aggregate(context.TODO(), filter)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var route models.ResponseFlightRoute2
			if err := cursor.Decode(&route); err != nil {
				return nil, err
			}
			uniqueRoutes[route.RouteId] = route

		}

		if err := cursor.Err(); err != nil {
			return nil, err
		}

	}

	for _, route := range uniqueRoutes {

		routes = append(routes, route)
	}

	return routes, nil
}

func FindRoutesWithIds(routeIds []int) ([]models.ResponseFlightRoute2, error) {
	collection := getCollection("routes")

	// Define a match stage to filter routes by routeIds
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "routeId", Value: bson.D{
				{Key: "$in", Value: routeIds},
			}},
		}},
	}

	// Define lookup stages to join with airports and airlines
	airportFromLookup := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "airports"},
			{Key: "localField", Value: "airportFrom"},
			{Key: "foreignField", Value: "iata"},
			{Key: "as", Value: "resAirportFrom"},
		}},
	}

	airportToLookup := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "airports"},
			{Key: "localField", Value: "airportTo"},
			{Key: "foreignField", Value: "iata"},
			{Key: "as", Value: "resAirportTo"},
		}},
	}

	airlineLookup := bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "airlines"},
			{Key: "localField", Value: "airline"},
			{Key: "foreignField", Value: "iata"},
			{Key: "as", Value: "resAirline"},
		}},
	}

	// Define unwind stages to flatten the arrays
	airportFromUnwind := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$resAirportFrom"},
			{Key: "includeArrayIndex", Value: "string"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}

	airportToUnwind := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$resAirportTo"},
			{Key: "includeArrayIndex", Value: "string"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}

	airlineUnwind := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$resAirline"},
			{Key: "includeArrayIndex", Value: "string"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}

	// Define the pipeline stages
	pipeline := []bson.D{
		matchStage,
		airportFromLookup,
		airportToLookup,
		airlineLookup,
		airportFromUnwind,
		airportToUnwind,
		airlineUnwind,
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var routes []models.ResponseFlightRoute2

	for cursor.Next(context.TODO()) {
		var route models.ResponseFlightRoute2
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}
