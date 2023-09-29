package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	err    error
)

func getCollection(colName string) *mongo.Collection {
	dbName := os.Getenv("DB_Name")
	return client.Database(dbName).Collection(colName)

}

func InitDB(dbUrl string) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(dbUrl)

	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
