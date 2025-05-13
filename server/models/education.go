package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Education struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string             `json:"name"`
}
