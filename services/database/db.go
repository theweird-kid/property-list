package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func ConnectDatabase() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	databaseName := os.Getenv("DB_NAME")

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	MongoDB = client.Database(databaseName)
	return client, nil
}
