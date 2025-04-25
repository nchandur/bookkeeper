package db

import (
	"bookkeeper/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() error {
	var err error

	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		return err
	}

	return nil

}

func InsertBooks(document model.Document) error {
	collection := Client.Database("books").Collection("works")

	_, err := collection.InsertOne(context.TODO(), document)

	if err != nil {
		return err
	}

	return nil
}
