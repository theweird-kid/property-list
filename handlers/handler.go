package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/theweird-kid/property-list/services/prop_service"
	"go.mongodb.org/mongo-driver/mongo"
)

type API struct {
	MongoDB     *mongo.Database
	RedisClient *redis.Client
}

func (api *API) Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}

func (api *API) GetProperties(ctx *gin.Context) {
	properties, err := prop_service.GetAllProperties(ctx, api.MongoDB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, properties)
}
