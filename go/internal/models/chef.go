package models

import (
	"github.com/globalsign/mgo/bson"
)

// Chef represents a chef on the platform
type Chef struct {
	ID            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FirstName     string        `json:"first_name" bson:"first_name"`
	LastName      string        `json:"last_name" bson:"last_name"`
	Password      string        `json:"password" bson:"password"`
	Location      *Location     `json:"location" bson:"location"`
	Rating        float32       `json:"rating" bson:"rating"`
	KnownRecipies []string      `json:"known_recipies" bson:"known_recipies"`
}
