package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Number      string             `json:"number" bson:"number"`
	Date        time.Time          `json:"date" bson:"date"`
	OrderTypeID primitive.ObjectID `json:"orderTypeId" bson:"orderTypeId"`
	OrderType   string             `json:"orderType" bson:"orderType"`
	FileURL     string             `json:"fileUrl" bson:"fileUrl"`   // URL к файлу приказа
	FileName    string             `json:"fileName" bson:"fileName"` // Имя файла
}

type OrderType struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}
