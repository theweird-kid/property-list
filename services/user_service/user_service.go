package user_service

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/theweird-kid/property-list/models"
	"github.com/theweird-kid/property-list/services/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	DB          *mongo.Database
	RedisClient *redis.Client
}

func NewUserService(db *mongo.Database, redis *redis.Client) *UserService {
	return &UserService{
		DB:          db,
		RedisClient: redis,
	}
}

func (us *UserService) GetUsers(ctx *gin.Context) ([]models.User, error) {
	usersCollection := us.DB.Collection("users")
	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) RegisterUser(ctx *gin.Context, user models.User) error {
	hashedPass, _ := auth.HashPassword(user.Password)
	user.Password = string(hashedPass)

	userCollection := us.DB.Collection("users")
	// check if already exists
	var existringUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existringUser)
	if err == mongo.ErrNoDocuments {
		_, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			return fmt.Errorf("Failed to register user: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("Failed to register user: %w", err)
	}

	return fmt.Errorf("user already exists with email %s", existringUser.Email)
}

func (us *UserService) LoginUser(ctx *gin.Context, user *models.User) (string, error) {
	userCollection := us.DB.Collection("users")
	// check if user exists
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		log.Println(err)
		return "", fmt.Errorf("invalid")
	}

	err = auth.CheckPassword(existingUser.Password, user.Password)
	if err != nil {
		log.Println("meow", err)
		return "", fmt.Errorf("invalid")
	}

	return auth.CreateToken(user.Email)
}
