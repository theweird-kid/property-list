package prop_service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/theweird-kid/property-list/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllProperties(ctx *gin.Context, db *mongo.Database) ([]models.Property, error) {
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

func GetPropertiesByUser(ctx *gin.Context, userEmail string, db *mongo.Database) ([]models.Property, error) {
	userCollection := db.Collection("users")
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		return nil, err
	}

	propertyCollection := db.Collection("properties")
	cursor, err := propertyCollection.Find(ctx, bson.M{"createdBy": user.ID})
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

	fmt.Println(len(properties))

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return properties, nil
}
