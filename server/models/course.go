package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Course struct {
	ID                primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title             string             `json:"title" bson:"title"`
	Description       string             `json:"description" bson:"description"`
	Duration          int                `json:"duration" bson:"duration"`
	Price             int                `json:"price" bson:"price"`
	PriceId           primitive.ObjectID `json:"priceId" bson:"priceId"`
	TypeId            primitive.ObjectID `json:"typeId" bson:"typeId"`
	Type              string             `json:"type" bson:"type"`
	CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
	RegistrationStart time.Time          `json:"registrationStart" bson:"registrationStart"`
	RegistrationEnd   time.Time          `json:"registrationEnd" bson:"registrationEnd"`
}
