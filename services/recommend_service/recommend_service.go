package recommend_service

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/theweird-kid/property-list/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecommendationService struct {
	DB          *mongo.Database
	RedisClient *redis.Client
}

func NewRecommendationService(db *mongo.Database, redis *redis.Client) *RecommendationService {
	return &RecommendationService{
		DB:          db,
		RedisClient: redis,
	}
}

func (rs *RecommendationService) GetRecommendations(ctx *gin.Context) ([]models.RecommendationResponse, error) {
	// get users
	userEmail := ctx.Keys["email"]
	userCollection := rs.DB.Collection("users")
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&user)
	if err != nil {
		return nil, err
	}

	propertyCollection := rs.DB.Collection("properties")

	// Get Recommendations
	var recommendations []models.RecommendationResponse
	recommendationCollection := rs.DB.Collection("recommendations")
	cursor, err := recommendationCollection.Find(ctx, bson.M{"recommendedToUserId": user.ID})
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var rec models.Recommendation
		if err := cursor.Decode(&rec); err != nil {
			return nil, err
		}
		var property models.Property
		if err := propertyCollection.FindOne(ctx, bson.M{"_id": rec.PropertyID}).Decode(&property); err != nil {
			return nil, err
		}

		var recommendingUser models.User
		objID, err := primitive.ObjectIDFromHex(rec.RecommendingUserID)
		if err != nil {
			return nil, err
		}
		if err := userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&recommendingUser); err != nil {
			log.Println("user prob", err)
			return nil, err
		}

		recommendation := models.RecommendationResponse{
			FromUserEmail: recommendingUser.Email,
			PropertyData:  property,
		}

		recommendations = append(recommendations, recommendation)
	}

	if cursor.Err() != nil {
		return nil, err
	}

	return recommendations, nil
}

type RecommendationRequest struct {
	ToEmail    string `json:"to_email"`
	PropertyID string `json:"prop_id"`
}

func (rs *RecommendationService) RecommendProperty(ctx *gin.Context) error {
	// get recommenders email
	recommendersEmail := ctx.Keys["email"]
	var req RecommendationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return err
	}

	// Validate user
	userCollection := rs.DB.Collection("users")
	var toUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": req.ToEmail}).Decode(&toUser)
	if err != nil {
		return err
	}

	var fromUser models.User
	err = userCollection.FindOne(ctx, bson.M{"email": recommendersEmail}).Decode(&fromUser)
	if err != nil {
		return err
	}

	// Validate property
	propertyCollection := rs.DB.Collection("properties")
	var property models.Property
	err = propertyCollection.FindOne(ctx, bson.M{"_id": req.PropertyID}).Decode(&property)
	if err != nil {
		return err
	}

	recommendationCollection := rs.DB.Collection("recommendations")
	recommendation := models.Recommendation{
		RecommendingUserID:  fromUser.ID,
		RecommendedToUserID: toUser.ID,
		PropertyID:          property.ID,
	}

	// Check if already available
	var existingRec models.Recommendation
	err = recommendationCollection.FindOne(ctx, bson.M{
		"recommendingUserID":  recommendation.RecommendingUserID,
		"recommendedToUserID": recommendation.RecommendedToUserID,
		"propertyID":          recommendation.PropertyID,
	}).Decode(&existingRec)

	if err == mongo.ErrNoDocuments {
		_, err = recommendationCollection.InsertOne(ctx, recommendation)
		if err != nil {
			return nil
		}
	}

	return err
}
