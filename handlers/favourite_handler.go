package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
 * protectedRoutes.GET("/fav", api.GetFavourites)
		protectedRoutes.POST("/fav-prop", api.FavouriteStatus)
*/

func (api *API) GetFavourites(ctx *gin.Context) {
	res, err := api.UserService.GetFavourites(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannnot retrieve favourites",
		})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (api *API) FavouriteProperty(ctx *gin.Context) {
	err := api.UserService.FavouriteProperty(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to process request",
		})
		return
	}

	ctx.JSON(http.StatusOK, "done")
}
