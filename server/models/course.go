package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Course struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title             string             `bson:"title" json:"title"`
	Description       string             `bson:"description" json:"description"`
	Duration          int                `bson:"duration" json:"duration"`
	PriceId           primitive.ObjectID `bson:"priceId" json:"priceId"`
	Price             int                `bson:"price" json:"price"`
	TypeId            primitive.ObjectID `bson:"typeId" json:"typeId"`
	Type              string             `bson:"type" json:"type"`
	CreatedAt         time.Time          `bson:"createdAt" json:"createdAt"`
	RegistrationStart time.Time          `bson:"registrationStart" json:"registrationStart"`
	RegistrationEnd   time.Time          `bson:"registrationEnd" json:"registrationEnd"`
	StudentsCount     int                `bson:"studentsCount" json:"studentsCount"`
	MaxStudents       int                `bson:"maxStudents" json:"maxStudents"` // Новое поле
}
