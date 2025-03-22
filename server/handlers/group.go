package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Group struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	GroupName string             `bson:"groupName" json:"groupName"`
	CourseID  primitive.ObjectID `bson:"courseId" json:"courseId"` // Исправлено на ObjectID
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if group.GroupName == "" || group.CourseID.IsZero() { // Проверка на пустой ObjectID
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

	// Проверяем, что курс существует
	courseCollection := client.Database("diplome").Collection("courses")
	var course bson.M
	err = courseCollection.FindOne(context.Background(), bson.M{"_id": group.CourseID}).Decode(&course)
	if err != nil {
		log.Printf("Error finding course: %v", err)
		http.Error(w, "Курс не найден", http.StatusNotFound)
		return
	}

	// Создаем группу
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

	// Преобразуем строковый ID в ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// Удаляем группу
	groupCollection := client.Database("diplome").Collection("groups")
	filter := bson.M{"_id": objectId}
	result, err := groupCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Error deleting group: %v", err)
		http.Error(w, "Ошибка при удалении группы", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Группа не найдена", http.StatusNotFound)
		return
	}

	// Удаляем все заявки, связанные с этой группой
	registrationCollection := client.Database("diplome").Collection("course_registrations")
	_, err = registrationCollection.DeleteMany(context.Background(), bson.M{"groupId": objectId})
	if err != nil {
		log.Printf("Error deleting registrations: %v", err)
		http.Error(w, "Ошибка при удалении заявок на курс", http.StatusInternalServerError)
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат идентификатора заявки"})
		return
	}

	var requestBody struct {
		GroupID string `json:"groupId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат данных"})
		return
	}

	if requestBody.GroupID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID группы обязателен"})
		return
	}

	groupId, err := primitive.ObjectIDFromHex(requestBody.GroupID)
	if err != nil {
		log.Printf("Invalid group ID: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат ID группы"})
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer client.Disconnect(context.Background())

	// Получаем информацию о заявке
	registrationCollection := client.Database("diplome").Collection("course_registrations")
	var registration bson.M
	err = registrationCollection.FindOne(context.Background(), bson.M{"_id": registrationId}).Decode(&registration)
	if err != nil {
		log.Printf("Error finding registration: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Заявка не найдена"})
		return
	}

	// Получаем ID курса из заявки
	courseId := registration["courseId"].(primitive.ObjectID)

	// Получаем информацию о группе
	groupCollection := client.Database("diplome").Collection("groups")
	var group bson.M
	err = groupCollection.FindOne(context.Background(), bson.M{"_id": groupId}).Decode(&group)
	if err != nil {
		log.Printf("Error finding group: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Группа не найдена"})
		return
	}

	// Проверяем, что группа принадлежит тому же курсу
	if group["courseId"].(primitive.ObjectID) != courseId {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Группа не принадлежит этому курсу"})
		return
	}

	// Обновляем группу в заявке
	update := bson.M{"$set": bson.M{"groupId": groupId}}
	result, err := registrationCollection.UpdateOne(context.Background(), bson.M{"_id": registrationId}, update)
	if err != nil {
		log.Printf("Error updating registration: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка при обновлении заявки"})
		return
	}

	if result.MatchedCount == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Заявка не найдена"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["id"]
	if groupId == "" {
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

	// Преобразуем строковый ID в ObjectID
	objectId, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
			http.Error(w, "Неверный формат ID", http.StatusBadRequest)
			return
	}

	// Получаем участников группы из коллекции course_registrations
	collection := client.Database("diplome").Collection("course_registrations")
	pipeline := bson.A{
			bson.M{
					"$match": bson.M{"groupId": objectId},
			},
			bson.M{
					"$lookup": bson.M{
							"from":         "users",
							"localField":   "userId",
							"foreignField":  "_id",
							"as":           "user",
					},
			},
			bson.M{
					"$project": bson.M{
							"username": bson.M{
									"$arrayElemAt": bson.A{"$user.username", 0},
							},
							"email": bson.M{
									"$arrayElemAt": bson.A{"$user.email", 0},
							},
					},
			},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
			log.Printf("Error fetching group members: %v", err)
			http.Error(w, "Ошибка при получении участников группы", http.StatusInternalServerError)
			return
	}
	defer cursor.Close(context.Background())

	var members []bson.M
	if err := cursor.All(context.Background(), &members); err != nil {
			log.Printf("Error decoding group members: %v", err)
			http.Error(w, "Ошибка при декодировании участников группы", http.StatusInternalServerError)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"members": members})
}