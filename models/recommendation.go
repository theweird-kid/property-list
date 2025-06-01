package models

import "time"

// Recommendation represents a property recommended by one user to another.
type Recommendation struct {
	RecommendingUserID  string    `json:"recommendingUserId" bson:"recommendingUserId"`
	RecommendedToUserID string    `json:"recommendedToUserId" bson:"recommendedToUserId"`
	PropertyID          string    `json:"propertyId" bson:"propertyId"`
	RecommendedAt       time.Time `json:"recommendedAt" bson:"recommendedAt"`
}

type RecommendationResponse struct {
	FromUserEmail string   `json:"from_user"`
	PropertyData  Property `json:"property"`
}
