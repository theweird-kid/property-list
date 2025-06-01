package user_service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/theweird-kid/property-list/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (us *UserService) GetFavourites(ctx *gin.Context) ([]models.Property, error) {
	email := ctx.Keys["email"]
	userCollection := us.DB.Collection("users")

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	favouriteCollection := us.DB.Collection("favourites")
	cursor, err := favouriteCollection.Find(ctx, bson.M{"userId": user.ID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var properties []models.Property
	propertyCollection := us.DB.Collection("properties")

	for cursor.Next(ctx) {

		var fav models.Favorite
		if err := cursor.Decode(&fav); err != nil {
			return nil, err
		}

		var prop models.Property
		if err := propertyCollection.FindOne(ctx, bson.M{"_id": fav.PropertyID}).Decode(&prop); err != nil {
			return nil, err
		}

		properties = append(properties, prop)
	}

	log.Println("here")

	if cursor.Err() != nil {
		return nil, err
	}

	return properties, nil
}

func (us *UserService) FavouriteProperty(ctx *gin.Context) error {

	propID := ctx.Query("property")
	if propID == "" {
		return fmt.Errorf("invalid request")
	}

	req := ctx.Query("req")
	if req == "" {
		return fmt.Errorf("invalid request")
	}
	requestType, _ := strconv.ParseBool(req)

	userCollection := us.DB.Collection("users")
	propertyCollection := us.DB.Collection("properties")

	// Get User
	email := ctx.Keys["email"]
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return err
	}

	// Validate property
	var prop models.Property
	err := propertyCollection.FindOne(ctx, bson.M{"_id": propID}).Decode(&prop)
	if err == mongo.ErrNoDocuments {
		return err
	}

	favouriteCollection := us.DB.Collection("favourites")
	if requestType == false {
		_, err := favouriteCollection.DeleteOne(ctx, bson.M{"propertyId": propID, "userId": user.ID})
		if err != nil {
			return err
		}
	} else {
		fav := models.Favorite{
			UserID:     user.ID,
			PropertyID: propID,
		}
		_, err := favouriteCollection.InsertOne(ctx, fav)
		if err != nil {
			return err
		}
	}

	return nil
}
