package models

import (
	"github.com/globalsign/mgo/bson"
)

// User represents a user on the platform
type User struct {
	ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FirstName      string        `json:"first_name" bson:"first_name"`
	LastName       string        `json:"last_name" bson:"last_name"`
	Password       string        `json:"password" bson:"password"`
	Location       *Location     `json:"location" bson:"location"`
	Allergens      []string      `json:"allergens" bson:"allergens"`
	FoodPreference string        `json:"food_preference" bson:"food_preference"`
}

// Location represents a users location
type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}
