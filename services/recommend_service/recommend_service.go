package recommend_service

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

func (rs *RecommendationService) GetRecommendations(ctx *gin.Context) {

}

func (rs *RecommendationService) RecommendProperty(ctx *gin.Context) {

}
