package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) GetRecommendations(ctx *gin.Context) {
	res, err := api.RecommendationService.GetRecommendations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal problem!",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) RecommendProperty(ctx *gin.Context) {
	err := api.RecommendationService.RecommendProperty(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal problem",
		})
		return
	}

	ctx.JSON(http.StatusCreated, "recommended successfully!")
}
