package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FullName     string             `json:"fullName"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	Password     string             `json:"password"`
	BirthDate    string             `json:"birthDate"`
	Residence    string             `json:"residence,omitempty"`
	Education    string             `json:"education,omitempty"`
	BirthPlace   string             `json:"birthPlace,omitempty"`
	HomeAddress  string             `json:"homeAddress,omitempty"`
	PassportData string             `json:"passportData,omitempty"`
	SNILS        string             `json:"snils,omitempty"`
}