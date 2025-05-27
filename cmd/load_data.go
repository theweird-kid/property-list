package main

import (
	"context"
	"time"

	"github.com/theweird-kid/property-list/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedUsers(db *mongo.Database) error {
	users := []models.User{
		{Name: "Alice", Email: "alice@example.com", Password: "hashedpassword1"},
		{Name: "Bob", Email: "bob@example.com", Password: "hashedpassword2"},
		{Name: "Charlie", Email: "charlie@example.com", Password: "hashedpassword3"},
	}
	collection := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, user := range users {
		count, err := collection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			return err
		}
		if count == 0 {
			_, err := collection.InsertOne(ctx, user)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
