package models

import "time"

type Property struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title"`
	Type      string `json:"type" bson:"type"`
	Price     int64  `json:"price" bson:"price"`
	State     string `json:"state" bson:"state"`
	City      string `json:"city" bson:"city"`
	AreaSqFt  int    `json:"areaSqFt" bson:"areaSqFt"`
	Bedrooms  int    `json:"bedrooms" bson:"bedrooms"`
	Bathrooms int    `json:"bathrooms" bson:"bathrooms"`

	Amenities     []string  `json:"amenities" bson:"amenities"`
	Furnished     string    `json:"furnished" bson:"furnished"` // Can be "Furnished", "Unfurnished", "Semi"
	AvailableFrom time.Time `json:"availableFrom" bson:"availableFrom"`
	ListedBy      string    `json:"listedBy" bson:"listedBy"`
	Tags          []string  `json:"tags" bson:"tags"`
	ColorTheme    string    `json:"colorTheme" bson:"colorTheme"`
	Rating        float32   `json:"rating" bson:"rating"` // 0 - 5
	IsVerified    bool      `json:"isVerified" bson:"isVerified"`
	ListingType   string    `json:"listingType" bson:"listingType"` // "rent" or "sale"

	// CreatedBy will store the User ID who created this property.
	CreatedBy string `json:"createdBy" bson:"createdBy"`

	// CreatedAt and UpdatedAt for auditing purposes.
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
