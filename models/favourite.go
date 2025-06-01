package models

import "time"

type Favorite struct {
	UserID      string    `json:"userId" bson:"userId"`
	PropertyID  string    `json:"propertyId" bson:"propertyId"`
	FavoritedAt time.Time `json:"favoritedAt" bson:"favoritedAt"`
}
