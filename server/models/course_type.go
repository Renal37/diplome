package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourseType struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string             `json:"name"`
}
