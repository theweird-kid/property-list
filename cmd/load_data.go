package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/theweird-kid/property-list/models"
	"github.com/theweird-kid/property-list/services/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const FILE_PATH string = "./static/properties.csv"

func seedUsers(ctx context.Context, db *mongo.Database) ([]string, error) {
	pass1, _ := auth.HashPassword("hashedpassword1")
	pass2, _ := auth.HashPassword("hashedpassword2")
	pass3, _ := auth.HashPassword("hashedpassword3")
	users := []models.User{
		{Name: "Alice", Email: "alice@example.com", Password: string(pass1)},
		{Name: "Bob", Email: "bob@example.com", Password: string(pass2)},
		{Name: "Charlie", Email: "charlie@example.com", Password: string(pass3)},
	}
	userCollection := db.Collection("users")
	var userIDs []string

	for _, user := range users {
		res, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}
		// Convert ObjectID to string
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			userIDs = append(userIDs, oid.Hex())
		} else if idStr, ok := res.InsertedID.(string); ok {
			userIDs = append(userIDs, idStr)
		}
	}

	return userIDs, nil
}

func LoadData(db *mongo.Database) error {
	var ctx context.Context

	// Create Dummy Users
	userIDs, err := seedUsers(ctx, db)
	if err != nil {
		return fmt.Errorf("Failed to seed users: %w", err)
	}

	// Load Properties
	properties, err := importPropertiesFromCSV(FILE_PATH)
	if err != nil {
		return fmt.Errorf("failed to import properties: %w", err)
	}

	// Adjust CreatedBy and time
	now := time.Now()
	for i := range properties {
		properties[i].CreatedBy = userIDs[i%len(userIDs)]
		properties[i].CreatedAt = now
		properties[i].UpdatedAt = now
	}

	// Store in DB
	propertyCollection := db.Collection("properties")
	var docs []interface{}
	for _, p := range properties {
		docs = append(docs, p)
	}
	if len(docs) > 0 {
		_, err = propertyCollection.InsertMany(ctx, docs)
		if err != nil {
			return fmt.Errorf("failed to insert properties: %w", err)
		}
	}

	return nil
}

func importPropertiesFromCSV(filePath string) ([]models.Property, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var properties []models.Property

	for i, row := range records {
		if i == 0 {
			continue
		}

		price, _ := strconv.ParseInt(row[3], 10, 64)
		areaSqFt, _ := strconv.Atoi(row[6])
		bedrooms, _ := strconv.Atoi(row[7])
		bathrooms, _ := strconv.Atoi(row[8])
		amenities := strings.Split(row[9], "|")
		availableFrom, _ := time.Parse("2006-01-02", row[11])
		tags := strings.Split(row[13], "|")
		rating, _ := strconv.ParseFloat(row[15], 32)
		isVerified := ParseBool(row[16])

		property := models.Property{
			ID:            row[0],
			Title:         row[1],
			Type:          row[2],
			Price:         price,
			State:         row[4],
			City:          row[5],
			AreaSqFt:      areaSqFt,
			Bedrooms:      bedrooms,
			Bathrooms:     bathrooms,
			Amenities:     amenities,
			Furnished:     row[10],
			AvailableFrom: availableFrom,
			ListedBy:      row[12],
			Tags:          tags,
			ColorTheme:    row[14],
			Rating:        float32(rating),
			IsVerified:    isVerified,
			ListingType:   row[17],
		}

		properties = append(properties, property)
	}

	return properties, nil
}

func ParseBool(val string) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	return val == "true" || val == "yes" || val == "1"
}
