package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Airport struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	IATA            string             `bson:"iata" json:"code"`
	IATACode        string             `bson:"-" json:"IATA"`
	Name            string             `bson:"name" json:"name"`
	City            string             `bson:"city" json:"city"`
	State           string             `bson:"state" json:"state"`
	Country         string             `bson:"country" json:"country"`
	URL             string             `bson:"url" json:"url"`
	ICAO            string             `bson:"icao" json:"icao"`
	DirectFlights   string             `bson:"directflights" json:"direct_flights"`
	Carriers        string             `bson:"carriers" json:"carriers"`
	NormalizedScore string             `bson:"normalizedScore"`
	Canceled        int                `bson:"canceled" `
	Delayed15       int                `bson:"delayed15"`
	Delayed30       int                `bson:"delayed30" `
	Delayed45       int                `bson:"delayed45"`
	OnTime          int                `bson:"onTime" `
	ModifiedAt      time.Time          `bson:"modifiedAt" json:"modifiedAt"`
}

type Airline struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name                 string             `bson:"name" json:"name"`
	IATA                 string             `bson:"iata" json:"code"`
	IsLowCost            bool               `bson:"is_lowcost" json:"is_lowcost"`
	Logo                 string             `bson:"logo" json:"logo"`
	IsScheduledPassenger int                `bson:"is_scheduled_passenger" json:"is_scheduled_passenger"`
	IsCargo              int                `bson:"is_cargo" json:"is_cargo"`
	ModifiedAt           time.Time          `bson:"modifiedat" json:"modifiedat"`
}

type Route struct {
}
type AirlineResponse struct {
	Airlines []Airline `json:"airlines"`
}

type AirportResponse struct {
	Airports []Airport `json:"airports"`
}

type RoutesResponse struct {
	Routes []ResponseFlightRoute `json:"data"`
}

type ResponseAirline struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name                 string             `bson:"name" json:"name"`
	IATA                 string             `bson:"iata" json:"code"`
	IATACode             string             `bson:"iata" json:"IATA"`
	IsLowCost            int                `bson:"is_lowcost" json:"is_lowcost"`
	Logo                 string             `bson:"logo" json:"logo"`
	IsScheduledPassenger int                `bson:"is_scheduled_passenger" json:"is_scheduled_passenger"`
	IsCargo              int                `bson:"is_cargo" json:"is_cargo"`
}
type ResponseFlightRoute struct {
	ID             primitive.ObjectID `bson:"-" json:"-"`
	RouteID        int                `bson:"route_id" json:"id"`
	Airline        ResponseAirline    `bson:"airline" json:"airline"`
	AirportFrom    Airport            `bson:"airportFrom" json:"airportFrom"`
	AirportTo      Airport            `bson:"airportTo" json:"airportTo"`
	AirportVia     Airport            `bson:"airportVia" json:"airportVia"`
	ClassBusiness  int                `bson:"class_business" json:"class_business"`
	ClassEconomy   int                `bson:"class_economy" json:"class_economy"`
	ClassFirst     int                `bson:"class_first" json:"class_first"`
	CommonDuration int                `bson:"common_duration" json:"common_duration"`
	Day1           string             `bson:"day1" json:"day1"`
	Day2           string             `bson:"day2" json:"day2"`
	Day3           string             `bson:"day3" json:"day3"`
	Day4           string             `bson:"day4" json:"day4"`
	Day5           string             `bson:"day5" json:"day5"`
	Day6           string             `bson:"day6" json:"day6"`
	Day7           string             `bson:"day7" json:"day7"`
	FlightsPerDay  string             `bson:"flights_per_day" json:"flights_per_day"`
	FlightsPerWeek int                `bson:"flights_per_week" json:"flights_per_week"`
	IsActive       int                `bson:"is_active" json:"is_active"`
	MaxDuration    int                `bson:"max_duration" json:"max_duration"`
	MinDuration    int                `bson:"min_duration" json:"min_duration"`
}

type FlightRoute struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	RouteId        int                `bson:"routeId,omitempty" json:"id"`
	Airline        string             `bson:"airline" json:"airline"`
	AirportFrom    string             `bson:"airportFrom" json:"airportFrom"`
	AirportTo      string             `bson:"airportTo" json:"airportTo"`
	AirportVia     string             `bson:"airportVia" json:"airportVia"`
	ClassBusiness  bool               `bson:"class_business" json:"class_business"`
	ClassEconomy   bool               `bson:"class_economy" json:"class_economy"`
	ClassFirst     bool               `bson:"class_first" json:"class_first"`
	CommonDuration int                `bson:"common_duration" json:"common_duration"`
	Monday         string             `bson:"Monday" json:"day1"`
	Tuesday        string             `bson:"Tuesday" json:"day2"`
	Wednesday      string             `bson:"Wednesday" json:"day3"`
	Thursday       string             `bson:"Thursday" json:"day4"`
	Friday         string             `bson:"Friday" json:"day5"`
	Saturday       string             `bson:"Saturday" json:"day6"`
	Sunday         string             `bson:"Sunday" json:"day7"`
	FlightsPerDay  string             `bson:"flights_per_day" json:"flights_per_day"`
	FlightsPerWeek int                `bson:"flights_per_week" json:"flights_per_week"`
	IsActive       bool               `bson:"is_active" json:"is_active"`
	MaxDuration    int                `bson:"max_duration" json:"max_duration"`
	MinDuration    int                `bson:"min_duration" json:"min_duration"`
	ModifiedAt     time.Time          `bson:"modifiedat" json:"modifiedat"`
}

type UpdatedRouteResponse struct {
	Airline     string `bson:"airline" json:"airline"`
	AirportFrom string `bson:"airportFrom" json:"airportFrom"`
	AirportTo   string `bson:"airportTo" json:"airportTo"`
	AirportVia  string `bson:"airportVia" json:"airportVia"`
}

type ResponseFlightRoute2 struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	RouteId        int                `bson:"routeId,omitempty" json:"id"`
	Airline        string             `bson:"airline" json:"airline"`
	AirportFrom    string             `bson:"airportFrom" json:"airportFrom"`
	AirportTo      string             `bson:"airportTo" json:"airportTo"`
	AirportVia     string             `bson:"airportVia" json:"airportVia"`
	ClassBusiness  bool               `bson:"class_business" json:"class_business"`
	ClassEconomy   bool               `bson:"class_economy" json:"class_economy"`
	ClassFirst     bool               `bson:"class_first" json:"class_first"`
	CommonDuration int                `bson:"common_duration" json:"common_duration"`
	Monday         string             `bson:"Monday" json:"day1"`
	Tuesday        string             `bson:"Tuesday" json:"day2"`
	Wednesday      string             `bson:"Wednesday" json:"day3"`
	Thursday       string             `bson:"Thursday" json:"day4"`
	Friday         string             `bson:"Friday" json:"day5"`
	Saturday       string             `bson:"Saturday" json:"day6"`
	Sunday         string             `bson:"Sunday" json:"day7"`
	FlightsPerDay  string             `bson:"flights_per_day" json:"flights_per_day"`
	FlightsPerWeek int                `bson:"flights_per_week" json:"flights_per_week"`
	IsActive       bool               `bson:"is_active" json:"is_active"`
	MaxDuration    int                `bson:"max_duration" json:"max_duration"`
	MinDuration    int                `bson:"min_duration" json:"min_duration"`
	ModifiedAt     time.Time          `bson:"modifiedat" json:"modifiedat"`
	ResAirportFrom Airport            `bson:"resAirportFrom" json:"resAirportFrom"`
	ResAirportTo   Airport            `bson:"resAirportTo" json:"resAirportTo"`
	ResAirline     Airline            `bson:"resAirline" json:"resAirline"`
}

type RecommendationRequest struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Date      string `json:"date"`
	Class     string `json:"class"`
}

type RouteIdRequest struct {
	RouteId []int `json:"routeId"`
}
type DetailRecommendationRequest struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Date      string `json:"date"`
}

type MultiCityRequest struct {
	Class    []string                      `json:"class"`
	Airlines []string                      `json:"airline"`
	Routes   []DetailRecommendationRequest `json:"routes"`
}

type MultiCityResponse struct {
	MultiCity []RouteResponse `json:"multicity"`
}

type RouteResponse struct {
	Date      string        `json:"date"`
	Departure string        `json:"departure"`
	Arrival   string        `json:"arrival"`
	Route     []FlightRoute `json:"route"`
}
type RouteResponse2 struct {
	Date      string      `json:"date"`
	Departure string      `json:"departure"`
	Arrival   string      `json:"arrival"`
	Route     []RouteInfo `json:"route"`
}

type MultiCityResponse2 struct {
	MultiCity []RouteResponse2 `json:"multicity"`
}
type AutoCompleteResponse struct {
	IATA     string `bson:"iata" json:"code"`
	IATACode string `bson:"-" json:"IATA"`
	Name     string `bson:"name" json:"name"`
	City     string `bson:"city" json:"city"`
	State    string `bson:"state" json:"state"`
}

type RouteInfo struct {
	RouteId        int    `bson:"routeId,omitempty" json:"id"`
	AirportFrom    string `json:"airportFrom"`
	AirportTo      string `json:"airportTo"`
	Airline        string `json:"ariline"`
	CommonDuration int    `bson:"common_duration"`
	ClassBusiness  bool   `bson:"class_bussiness"`
	ClassEconomy   bool   `bson:"class_economy"`
	ClassFirst     bool   `bson:"class_first"`
	Monday         string `bson:"Monday"`
	Tuesday        string `bson:"Tuesday"`
	Wednesday      string `bson:"Wednesday"`
	Thursday       string `bson:"Thursday"`
	Friday         string `bson:"Friday"`
	Saturday       string `bson:"Saturday"`
	Sunday         string `bson:"Sunday"`
	// Corrected field name
	ResAirportTo   Airport `bson:"resAirportTo"`
	ResAirportFrom Airport `bson:"resAirportFrom"`
	AirlineInfo    Airline `bson:"resAirline"`
}

type RouteFilterResponse struct {
	Data  []ResponseFlightRoute `json:"data"`
	Total int                   `json:"total"`
	Last  int                   `json:"last"`
}
