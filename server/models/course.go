package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Duration    int                `json:"duration" bson:"duration"` // Продолжительность в часах
	Price       int                `json:"price" bson:"price"`       // Стоимость курса
	PriceId     primitive.ObjectID `json:"priceId" bson:"priceId"`   // ID стоимости
	TypeId      primitive.ObjectID `json:"typeId" bson:"typeId"`     // ID типа курса
	Type        string             `json:"type" bson:"type"`         // Название типа курса
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}
