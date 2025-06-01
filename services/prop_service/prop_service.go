package prop_service

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/theweird-kid/property-list/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PropertyService struct {
	DB          *mongo.Database
	RedisClient *redis.Client
}

func NewPropertyService(db *mongo.Database, redis *redis.Client) *PropertyService {
	return &PropertyService{
		DB:          db,
		RedisClient: redis,
	}
}

func (ps *PropertyService) GetAllProperties(ctx *gin.Context) ([]models.Property, error) {
	collection := ps.DB.Collection("properties")
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

func (ps *PropertyService) GetPropertiesByUser(ctx *gin.Context, userEmail string) ([]models.Property, error) {
	userCollection := ps.DB.Collection("users")
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		return nil, err
	}

	propertyCollection := ps.DB.Collection("properties")
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

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return properties, nil
}

func (ps *PropertyService) GetPropertiesByAttributes(ctx *gin.Context) ([]models.Property, error) {
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
	*/

	propertyCollection := ps.DB.Collection("properties")
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

func (ps *PropertyService) NewProperty(ctx *gin.Context) error {
	var prop models.Property
	if err := ctx.ShouldBindJSON(&prop); err != nil {
		return err
	}

	userCollection := ps.DB.Collection("users")

	email := ctx.Keys["email"]
	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return err
	}
	prop.CreatedBy = user.ID
	propId, _ := ps.getNextPropertyID(ctx)
	prop.ID = propId

	propertyCollection := ps.DB.Collection("properties")
	_, err := propertyCollection.InsertOne(ctx, prop)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PropertyService) UpdateProperty(ctx *gin.Context) (models.Property, error) {
	email := ctx.Keys["email"]
	userCollection := ps.DB.Collection("users")
	propertyCollection := ps.DB.Collection("properties")

	var user models.User
	if err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return models.Property{}, err
	}

	var propertyReq models.Property
	if err := ctx.ShouldBindJSON(&propertyReq); err != nil {
		return models.Property{}, err
	}

	var existingProperty models.Property
	if err := propertyCollection.FindOne(ctx, bson.M{"_id": propertyReq.ID}).Decode(&existingProperty); err != nil {
		return models.Property{}, err
	}

	if existingProperty.CreatedBy != user.ID {
		return models.Property{}, fmt.Errorf("not authroized")
	}

	update := buildPropertyUpdateMap(&propertyReq)
	if len(update) == 0 {
		return existingProperty, nil // Nothing to update
	}
	_, err := propertyCollection.UpdateByID(ctx, existingProperty.ID, bson.M{"$set": update})
	if err != nil {
		log.Println("here", err)
		return models.Property{}, err
	}

	if err := propertyCollection.FindOne(ctx, bson.M{"_id": propertyReq.ID}).Decode(&propertyReq); err != nil {
		return models.Property{}, err
	}
	return propertyReq, nil
}

func buildPropertyUpdateMap(req *models.Property) bson.M {
	update := bson.M{}

	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Type != "" {
		update["type"] = req.Type
	}
	if req.Price != 0 {
		update["price"] = req.Price
	}
	if req.State != "" {
		update["state"] = req.State
	}
	if req.City != "" {
		update["city"] = req.City
	}
	if req.AreaSqFt != 0 {
		update["areaSqFt"] = req.AreaSqFt
	}
	if req.Bedrooms != 0 {
		update["bedrooms"] = req.Bedrooms
	}
	if req.Bathrooms != 0 {
		update["bathrooms"] = req.Bathrooms
	}
	if req.Amenities != nil {
		update["amenities"] = req.Amenities
	}
	if req.Furnished != "" {
		update["furnished"] = req.Furnished
	}
	if !req.AvailableFrom.IsZero() {
		update["availableFrom"] = req.AvailableFrom
	}
	if req.ListedBy != "" {
		update["listedBy"] = req.ListedBy
	}
	if req.Tags != nil {
		update["tags"] = req.Tags
	}
	if req.ColorTheme != "" {
		update["colorTheme"] = req.ColorTheme
	}
	if req.Rating != 0 {
		update["rating"] = req.Rating
	}
	// Only update IsVerified if the request explicitly sets it to true
	if req.IsVerified {
		update["isVerified"] = req.IsVerified
	}
	if req.ListingType != "" {
		update["listingType"] = req.ListingType
	}
	// Add more fields as needed

	update["updatedAt"] = time.Now()
	return update
}

func (ps *PropertyService) getNextPropertyID(ctx *gin.Context) (string, error) {
	propertyCollection := ps.DB.Collection("properties")

	// Only fetch the PropertyID field
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	cursor, err := propertyCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	maxID := 0
	re := regexp.MustCompile(`PROP(\d{4})`)

	for cursor.Next(ctx) {
		var result struct {
			PropertyID string `bson:"_id"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		matches := re.FindStringSubmatch(result.PropertyID)
		if len(matches) == 2 {
			num, err := strconv.Atoi(matches[1])
			if err == nil && num > maxID {
				maxID = num
			}
		}
	}

	nextID := fmt.Sprintf("PROP%04d", maxID+1)
	return nextID, nil
}
