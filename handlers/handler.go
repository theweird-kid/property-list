package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theweird-kid/property-list/services/prop_service"
	"github.com/theweird-kid/property-list/services/recommend_service"
	"github.com/theweird-kid/property-list/services/user_service"
)

type API struct {
	PropertyService       *prop_service.PropertyService
	UserService           *user_service.UserService
	RecommendationService *recommend_service.RecommendationService
}

func (api *API) Hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}

func (api *API) GetProperties(ctx *gin.Context) {
	properties, err := api.PropertyService.GetAllProperties(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, properties)
}
