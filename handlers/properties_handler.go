package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theweird-kid/property-list/services/prop_service"
)

func (api *API) GetUserProperties(ctx *gin.Context) {
	email, exists := ctx.Get("email")
	userEmail := email.(string)
	if exists == false {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "not logged in",
		})
		return
	}

	properties, err := prop_service.GetPropertiesByUser(ctx, userEmail, api.MongoDB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal problem",
		})
		return
	}

	ctx.JSON(http.StatusOK, properties)
}

func (api *API) GetPropertiesByAttributes(ctx *gin.Context) {
	results, err := prop_service.GetPropertiesByAttributes(ctx, api.MongoDB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to find properties",
		})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (api *API) NewProperty(ctx *gin.Context) {

}

func (api *API) UpdateProperty(ctx *gin.Context) {

}
