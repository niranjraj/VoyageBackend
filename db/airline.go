package db

import (
	"context"
	"errors"
	"hackathon/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertManyAirlines(airlines []models.Airline) error {

	collection := getCollection("airlines")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "iata", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return err
	}

	var documents []interface{}

	for _, airline := range airlines {
		documents = append(documents, airline)
	}

	_, err := collection.InsertMany(context.TODO(), documents)
	return err
}

func GetAirlineByIATACode(iataCode string) (models.Airline, error) {
	collection := getCollection("airlines")
	filter := bson.M{"iata": iataCode}
	var airline models.Airline
	err = collection.FindOne(context.TODO(), filter).Decode(&airline)
	if err != nil {

		return models.Airline{}, errors.New("Airline not found")
	}
	return airline, nil
}
