package models

import "time"

type Property struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Title     string `json:"title" bson:"title"`
	Type      string `json:"type" bson:"type"`
	Price     int64  `json:"price" bson:"price"` // Using int64 for price to avoid floating point issues and handle large values
	State     string `json:"state" bson:"state"`
	City      string `json:"city" bson:"city"`
	AreaSqFt  int    `json:"areaSqFt" bson:"areaSqFt"`
	Bedrooms  int    `json:"bedrooms" bson:"bedrooms"`
	Bathrooms int    `json:"bathrooms" bson:"bathrooms"`
	// Amenities and Tags are pipe-separated in the CSV, so they are represented as slices of strings.
	Amenities     []string  `json:"amenities" bson:"amenities"`
	Furnished     string    `json:"furnished" bson:"furnished"`         // Can be "Furnished", "Unfurnished", "Semi"
	AvailableFrom time.Time `json:"availableFrom" bson:"availableFrom"` // Use time.Time for date parsing
	ListedBy      string    `json:"listedBy" bson:"listedBy"`           // Corresponds to 'createdBy' in requirements
	Tags          []string  `json:"tags" bson:"tags"`
	ColorTheme    string    `json:"colorTheme" bson:"colorTheme"`
	Rating        float32   `json:"rating" bson:"rating"` // Using float32 for decimal ratings
	IsVerified    bool      `json:"isVerified" bson:"isVerified"`
	ListingType   string    `json:"listingType" bson:"listingType"` // "rent" or "sale"

	// CreatedBy will store the User ID who created this property.
	// This is crucial for the authorization logic (only creator can update/delete).
	CreatedBy string `json:"createdBy" bson:"createdBy"`

	// CreatedAt and UpdatedAt for auditing purposes.
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
