package handlers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type Group struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	GroupName string             `bson:"groupName" json:"groupName"`
	CourseID  string             `bson:"courseId" json:"courseId"`
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if group.GroupName == "" || group.CourseID == "" {
		http.Error(w, "Название группы и ID курса обязательны", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("groups")
	_, err = collection.InsertOne(context.Background(), group)
	if err != nil {
		log.Printf("Error inserting group: %v", err)
		http.Error(w, "Ошибка при создании группы", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func GetGroups(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("groups")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error fetching groups: %v", err)
		http.Error(w, "Ошибка при получении групп", http.StatusInternalServerError)
		return
	}

	var groups []Group
	if err := cursor.All(context.Background(), &groups); err != nil {
		log.Printf("Error decoding groups: %v", err)
		http.Error(w, "Ошибка при декодировании групп", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"groups": groups})
}
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var updatedGroup Group
	err := json.NewDecoder(r.Body).Decode(&updatedGroup)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if updatedGroup.GroupName == "" || updatedGroup.CourseID == "" {
		http.Error(w, "Название группы и ID курса обязательны", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID группы не указан", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("groups")

	// Преобразуем строковый ID в ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{
		"groupName": updatedGroup.GroupName,
		"courseId":  updatedGroup.CourseID,
	}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating group: %v", err)
		http.Error(w, "Ошибка при обновлении группы", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Группа не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID группы не указан", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("groups")

	// Преобразуем строковый ID в ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectId}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting group: %v", err)
		http.Error(w, "Ошибка при удалении группы", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Группа не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func AssignGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		log.Printf("Invalid registration ID: %v", err)
		http.Error(w, "Неверный формат идентификатора заявки", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		GroupID string `json:"groupId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.GroupID == "" {
		http.Error(w, "ID группы обязателен", http.StatusBadRequest)
		return
	}

	groupId, err := primitive.ObjectIDFromHex(requestBody.GroupID)
	if err != nil {
		log.Printf("Invalid group ID: %v", err)
		http.Error(w, "Неверный формат ID группы", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("course_registrations")
	filter := bson.M{"_id": registrationId}
	update := bson.M{"$set": bson.M{"groupId": groupId}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating registration: %v", err)
		http.Error(w, "Ошибка при обновлении заявки", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
