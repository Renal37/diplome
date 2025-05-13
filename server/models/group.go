package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Group struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	GroupName string             `bson:"groupName" json:"groupName"`
	CourseID  primitive.ObjectID `bson:"courseId" json:"courseId"` // Исправлено на ObjectID
}