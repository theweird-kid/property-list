package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

	properties, err := api.PropertyService.GetPropertiesByUser(ctx, userEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal problem",
		})
		return
	}

	ctx.JSON(http.StatusOK, properties)
}

func (api *API) GetPropertiesByAttributes(ctx *gin.Context) {
	results, err := api.PropertyService.GetPropertiesByAttributes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to find properties",
		})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (api *API) NewProperty(ctx *gin.Context) {
	if err := api.PropertyService.NewProperty(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to create property",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "property created",
	})
}

func (api *API) UpdateProperty(ctx *gin.Context) {
	res, err := api.PropertyService.UpdateProperty(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to update property",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
