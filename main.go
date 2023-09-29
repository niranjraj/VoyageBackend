package main

import (
	"context"
	"flag"
	"fmt"
	"hackathon/db"
	"hackathon/ingestion"
	"hackathon/router"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var inject bool

func init() {
	flag.BoolVar(&inject, "inject", false, "Enable injection ")

}

func main() {
	flag.Parse()
	start := time.Now()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	dbUrl := os.Getenv("DB_HOST")
	fmt.Printf("Initializing Database\n")
	client, err := db.InitDB(dbUrl)
	if err != nil {
		fmt.Println("Error initializing the database:", err)
		return
	}
	defer client.Disconnect(context.TODO())

	if inject {
		fmt.Printf("Fetching Routes\n")
		flightRoutes, err := ingestion.FetchFlightRoutesFromLocalFile("./data/routes.json")
		if err != nil {
			fmt.Println("Error fetching routes:", err)
			return
		}
		elapsed := time.Since(start)
		fmt.Printf("Time to fetch routes: %s\n", elapsed)

		fmt.Printf("Inserting Airlines\n")
		airlines, err := ingestion.GetAirlines(flightRoutes)
		if err != nil {
			fmt.Println("Error fetching airlines:", err)
			return
		}
		err = db.InsertManyAirlines(airlines)
		if err != nil {
			fmt.Println("Error airlines injection:", err)
			return
		}

		fmt.Printf("Inserting Airports\n")
		airports, err := ingestion.GetAirports(flightRoutes)
		if err != nil {
			fmt.Println("Error fetching airport:", err)
			return
		}
		db.InsertManyAirPorts(airports)

		fmt.Printf("Inserting Routes\n")
		// flight := injection.FetchRoutesFromUrl(airports)

		updatedroutes := ingestion.ConvertResponseToFlightRoute2(flightRoutes)
		fmt.Print(updatedroutes)
		err = db.InsertBatchRoutes(updatedroutes)
		if err != nil {
			fmt.Println("Error routes injection:", err)
			return
		}
		elapsed = time.Since(start)
		fmt.Printf("Execution time: %s\n", elapsed)
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return os.Getenv("ENVIRONMENT") == "development"
		},
	}))

	router.CreateRouterGroup(app)

	err = app.Listen(":8080")

	if err != nil {
		panic(err)
	}

}
