package db

import (
	"context"
	"errors"
	"fmt"
	"hackathon/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getAirportSearchPipeline(searchString string) bson.A {
	return bson.A{
		bson.D{
			{Key: "$search",
				Value: bson.D{
					{Key: "index", Value: "default"},
					{Key: "compound",
						Value: bson.D{
							{Key: "should",
								Value: bson.A{
									bson.D{
										{Key: "wildcard",
											Value: bson.D{
												{Key: "path", Value: "iata"},
												{Key: "query", Value: searchString},
												{Key: "allowAnalyzedField", Value: true},
											},
										},
									},
									bson.D{
										{Key: "autocomplete",
											Value: bson.D{
												{Key: "path", Value: "name"},
												{Key: "query", Value: searchString},
											},
										},
									},
									bson.D{
										{Key: "autocomplete",
											Value: bson.D{
												{Key: "path", Value: "city"},
												{Key: "query", Value: searchString},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$limit", Value: 3},
		},
	}
}

func InsertManyAirPorts(airports []models.Airport) error {

	collection := getCollection("airports")

	var documents []interface{}
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "iata", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	for _, airport := range airports {
		documents = append(documents, airport)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	return err
}

func GetAirportByIATACode(iataCode string) (models.Airport, error) {
	collection := getCollection("airports")
	filter := bson.M{"iata": iataCode}
	var airport models.Airport
	err = collection.FindOne(context.TODO(), filter).Decode(&airport)
	if err != nil {

		return models.Airport{}, errors.New("Airport not found")
	}
	return airport, nil
}

func FindAirportsWithRegex(searchString string) ([]models.Airport, error) {
	var results []models.Airport

	// Create a MongoDB cursor for querying the airport collection
	// i option for case insensitive
	filter := bson.M{
		"$or": []bson.M{
			{"code": bson.M{"$regex": primitive.Regex{Pattern: searchString, Options: "i"}}},
			{"name": bson.M{"$regex": primitive.Regex{Pattern: searchString, Options: "i"}}},
		},
	}

	collection := getCollection("airports")
	cursor, err := collection.Find(context.TODO(), filter) // Limit the results to 3
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var airport models.Airport
		if err := cursor.Decode(&airport); err != nil {
			return nil, err
		}
		results = append(results, airport)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func FindAirportsWithSearch(searchString string) ([]models.AutoCompleteResponse, error) {
	var results []models.AutoCompleteResponse
	fmt.Print("in airline")
	collection := getCollection("airports")
	cursor, err := collection.Aggregate(context.TODO(), getAirportSearchPipeline(searchString))
	if err != nil {
		fmt.Print("pipeline")
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var airport models.AutoCompleteResponse
		if err := cursor.Decode(&airport); err != nil {
			return nil, err
		}
		results = append(results, airport)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
