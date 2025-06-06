package main

import (
	"github.com/gin-gonic/gin"
	"github.com/theweird-kid/property-list/handlers"
	"github.com/theweird-kid/property-list/services/auth"
)

func setupRoutes(r *gin.Engine, api *handlers.API) {
	r.GET("/", api.Hello)

	r.GET("/properties", api.GetProperties)
	r.GET("/prop-search", api.GetPropertiesByAttributes)

	// Query filter for email
	r.GET("/users", api.GetUsers)

	r.POST("/register", api.RegisterUser)
	r.POST("/login", api.LoginUser)

	protectedRoutes := r.Group("/auth")
	protectedRoutes.Use(auth.AuthMiddleware())
	{
		protectedRoutes.GET("/my-props", api.GetUserProperties)
		protectedRoutes.POST("/add-prop", api.NewProperty)
		protectedRoutes.PUT("/update-prop", api.UpdateProperty)

		protectedRoutes.GET("/fav", api.GetFavourites)
		protectedRoutes.POST("/fav-prop", api.FavouriteProperty)

		protectedRoutes.GET("/my-rec", api.GetRecommendations)
		protectedRoutes.POST("/rec-prop", api.RecommendProperty)
	}
}
