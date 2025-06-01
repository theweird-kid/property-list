package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/theweird-kid/property-list/handlers"
	"github.com/theweird-kid/property-list/services/cache"
	"github.com/theweird-kid/property-list/services/database"
	"github.com/theweird-kid/property-list/services/prop_service"
	"github.com/theweird-kid/property-list/services/recommend_service"
	"github.com/theweird-kid/property-list/services/user_service"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found or error loading .env")
	}

	fmt.Println("Property Listing system!!")

	// Connect to MongoDB
	_, err = database.ConnectDatabase()
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	// Connect to Redis
	err = cache.ConnectRedis()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	propertyService := prop_service.NewPropertyService(database.MongoDB, cache.RedisClient)
	userService := user_service.NewUserService(database.MongoDB, cache.RedisClient)
	recommendationService := recommend_service.NewRecommendationService(database.MongoDB, cache.RedisClient)

	api := &handlers.API{
		PropertyService:       propertyService,
		UserService:           userService,
		RecommendationService: recommendationService,
	}

	// err = LoadData(api.MongoDB)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Data Loaded Successfully")
	// }

	r := gin.Default()
	setupRoutes(r, api)

	r.Run("0.0.0.0:8080")
}
