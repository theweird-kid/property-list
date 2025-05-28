package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/theweird-kid/property-list/handlers"
	"github.com/theweird-kid/property-list/services/cache"
	"github.com/theweird-kid/property-list/services/database"
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

	api := &handlers.API{
		MongoDB:     database.MongoDB,
		RedisClient: cache.RedisClient,
	}

	// err = LoadData(api.MongoDB)
	// if err != nil {
	// 	log.Fatal(err)
	// } else if err == nil {
	// 	log.Println("Data Loaded Successfully")
	// }

	r := gin.New()

	r.GET("/", api.Hello)
	r.GET("/properties", api.GetProperties)

	// Example: pass api to handlers as needed
	// r.GET("/properties", api.GetProperties)

	r.Run("0.0.0.0:8080")
}
