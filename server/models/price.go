package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Price struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Amount      int                `bson:"amount" json:"amount"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	Description string             `bson:"description" json:"description"`
}
