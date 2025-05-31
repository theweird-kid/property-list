package prop_service

import (
	"fmt"
	"strconv"
	"time"

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

func GetPropertiesByAttributes(ctx *gin.Context, db *mongo.Database) ([]models.Property, error) {
	filters := bson.M{}

	if title := ctx.Query("title"); title != "" {
		filters["title"] = bson.M{"$regex": title}
	}
	if prop_type := ctx.Query("type"); prop_type != "" {
		filters["type"] = bson.M{"$regex": prop_type}
	}
	if price := ctx.Query("price"); price != "" {
		priceInt, _ := strconv.ParseInt(price, 10, 64)
		filters["price"] = bson.M{
			"$gte": max(0, priceInt-50000),
			"$lte": priceInt + 50000,
		}
	}
	if city := ctx.Query("city"); city != "" {
		filters["city"] = bson.M{"$regex": city}
	}
	if state := ctx.Query("state"); state != "" {
		filters["state"] = bson.M{"$regex": state}
	}
	if furnished := ctx.Query("furnished"); furnished != "" {
		filters["furnished"] = bson.M{"$regex": furnished}
	}
	if bedrooms := ctx.Query("bedrooms"); bedrooms != "" {
		bed_cnt, _ := strconv.ParseInt(bedrooms, 10, 32)
		filters["bedrooms"] = bson.M{"$gte": bed_cnt}
	}
	if bathrooms := ctx.Query("bathrooms"); bathrooms != "" {
		bath_cnt, _ := strconv.ParseInt(bathrooms, 10, 32)
		filters["bathrooms"] = bson.M{"$gte": bath_cnt}
	}
	if listedBy := ctx.Query("listedBy"); listedBy != "" {
		filters["listedBy"] = bson.M{"$regex": listedBy}
	}
	if colorTheme := ctx.Query("colorTheme"); colorTheme != "" {
		filters["colorTheme"] = bson.M{"$regex": colorTheme}
	}
	if listingType := ctx.Query("listingType"); listingType != "" {
		filters["listingType"] = bson.M{"$regex": listingType}
	}
	if availableFrom := ctx.Query("availableFrom"); availableFrom != "" {
		date, _ := time.Parse("2006-01-02", availableFrom)
		till_date := date.AddDate(0, 6, 0)
		filters["availableFrom"] = bson.M{
			"$gte": date,
			"$lte": till_date,
		}
	}
	if rating := ctx.Query("rating"); rating != "" {
		ratingFloat, _ := strconv.ParseFloat(rating, 32)
		filters["rating"] = bson.M{"$gte": ratingFloat}
	}
	if isVerified := ctx.Query("isVerified"); isVerified != "" {
		verifiedBool, _ := strconv.ParseBool(isVerified)
		filters["isVerified"] = verifiedBool
	}

	/*

		json:"areaSqFt" bson:"areaSqFt"`

		// Amenities and Tags are pipe-separated in the CSV, so they are represented as slices of strings.
		Amenities     []string  `json:"amenities" bson:"amenities"`


		Tags          []string  `json:"tags" bson:"
		IsVerified    bool      `json:"isVerified" bson:"isVerified"`

	*/

	propertyCollection := db.Collection("properties")
	cursor, err := propertyCollection.Find(ctx, filters)
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
