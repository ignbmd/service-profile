package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongodb *mongo.Database
var Mongocontext context.Context

func Connect(connection string, collection string) error {
	if Mongodb == nil {
		Mongocontext = context.Background()
		clientOptions := options.Client()
		clientOptions.ApplyURI(connection)
		client, err := mongo.NewClient(clientOptions)
		if err != nil {
			return err
		}
		err = client.Connect(Mongocontext)
		if err != nil {
			return err
		}
		Mongodb = client.Database(collection)
	}

	return nil
}
