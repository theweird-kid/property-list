package propservice

import (
	"context"

	"github.com/theweird-kid/property-list/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllProperties(ctx context.Context, db *mongo.Database) ([]models.Property, error) {
	collection := db.Collection("properties")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []models.Property
	for cursor.Next(ctx) {
		var prop models.Property
		if err := cursor.Decode(&prop); err != nil {
			return nil, err
		}
		properties = append(properties, prop)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return properties, nil
}
