package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Duration    int                `json:"duration"` // Продолжительность в часах
	Price       int                `json:"price"`    // Стоимость курса
	PriceId     primitive.ObjectID `json:"priceId"`
	TypeId      primitive.ObjectID `json:"typeId"`
	Type        string             `json:"type"` // Тип курса
}
