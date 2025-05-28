package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if group.GroupName == "" || group.CourseID.IsZero() {
		http.Error(w, "Название группы и ID курса обязательны", http.StatusBadRequest)
		return
	}

	courseCollection := db.GetCollection(db.CoursesCollection)
	var course bson.M
	err = courseCollection.FindOne(context.Background(), bson.M{"_id": group.CourseID}).Decode(&course)
	if err != nil {
		log.Printf("Error finding course: %v", err)
		http.Error(w, "Курс не найден", http.StatusNotFound)
		return
	}

	collection := db.GetCollection(db.GroupsCollection)
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
	collection := db.GetCollection(db.GroupsCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		writeJSONError(w, "Ошибка при получении групп", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var groups []bson.M
	if err = cursor.All(context.Background(), &groups); err != nil {
		writeJSONError(w, "Ошибка при обработке данных групп", http.StatusInternalServerError)
		return
	}

	regCollection := db.GetCollection(db.CourseRegistrationsCollection)
	for i := range groups {
		groupId, ok := groups[i]["_id"].(primitive.ObjectID)
		if !ok {
			writeJSONError(w, "Неверный формат ID группы", http.StatusInternalServerError)
			return
		}
		groups[i]["_id"] = groupId.Hex()
		if courseId, ok := groups[i]["courseId"].(primitive.ObjectID); ok {
			groups[i]["courseId"] = courseId.Hex()
		}
		count, err := regCollection.CountDocuments(
			context.Background(),
			bson.M{"groupId": groupId, "status": bson.M{"$ne": "Отчислен"}},
		)
		if err != nil {
			writeJSONError(w, "Ошибка при подсчете участников", http.StatusInternalServerError)
			return
		}
		groups[i]["currentStudents"] = count

		hasCompleted, err := regCollection.CountDocuments(
			context.Background(),
			bson.M{"groupId": groupId, "status": "Завершил"},
		)
		if err != nil {
			writeJSONError(w, "Ошибка при проверке статуса группы", http.StatusInternalServerError)
			return
		}
		if hasCompleted > 0 {
			groups[i]["status"] = "Завершена"
		} else {
			groups[i]["status"] = "Активна"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"groups": groups})
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var updatedGroup models.Group
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

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	count, err := registrationCollection.CountDocuments(context.Background(), bson.M{"groupId": id})
	if err != nil {
		log.Printf("Error checking group members: %v", err)
		http.Error(w, "Ошибка при проверке участников группы", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Невозможно изменить курс, так как в группе есть участники", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.GroupsCollection)
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

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	count, err := registrationCollection.CountDocuments(context.Background(), bson.M{
		"groupId": objectId,
		"$or": []bson.M{
			{"status": "Проходит курс"},
			{"status": "Завершил"},
		},
	})
	if err != nil {
		log.Printf("Error checking group members: %v", err)
		http.Error(w, "Ошибка при проверке участников группы", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Нельзя удалить группу, так как в ней есть участники со статусом 'Проходит курс' или 'Завершил'", http.StatusBadRequest)
		return
	}

	groupCollection := db.GetCollection(db.GroupsCollection)
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

func RejectGroupRegistrations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "ID группы не указан", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Reason string `json:"reason"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.Reason == "" {
		http.Error(w, "Причина отклонения обязательна", http.StatusBadRequest)
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	count, err := registrationCollection.CountDocuments(context.Background(), bson.M{
		"groupId": objectId,
		"status":  "Проходит курс",
	})
	if err != nil {
		log.Printf("Error checking group members: %v", err)
		http.Error(w, "Ошибка при проверке участников группы", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Нельзя отклонить записи, так как есть участники со статусом 'Проходит курс'", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"status":       "Отклоненный",
			"rejectReason": requestBody.Reason,
		},
	}
	_, err = registrationCollection.UpdateMany(context.Background(), bson.M{"groupId": objectId}, update)
	if err != nil {
		log.Printf("Error updating registrations: %v", err)
		http.Error(w, "Ошибка при отклонении записей", http.StatusInternalServerError)
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
		writeJSONError(w, "Неверный формат идентификатора заявки", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		GroupID string `json:"groupId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.GroupID == "" {
		writeJSONError(w, "ID группы обязателен", http.StatusBadRequest)
		return
	}

	groupId, err := primitive.ObjectIDFromHex(requestBody.GroupID)
	if err != nil {
		log.Printf("Invalid group ID: %v", err)
		writeJSONError(w, "Неверный формат ID группы", http.StatusBadRequest)
		return
	}

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	var registration bson.M
	err = registrationCollection.FindOne(context.Background(), bson.M{"_id": registrationId}).Decode(&registration)
	if err != nil {
		log.Printf("Error finding registration: %v", err)
		writeJSONError(w, "Заявка не найдена", http.StatusNotFound)
		return
	}

	courseId := registration["courseId"].(primitive.ObjectID)
	groupCollection := db.GetCollection(db.GroupsCollection)
	var group bson.M
	err = groupCollection.FindOne(context.Background(), bson.M{"_id": groupId}).Decode(&group)
	if err != nil {
		log.Printf("Error finding group: %v", err)
		writeJSONError(w, "Группа не найдена", http.StatusNotFound)
		return
	}

	if group["courseId"].(primitive.ObjectID) != courseId {
		writeJSONError(w, "Группа не принадлежит этому курсу", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": bson.M{"groupId": groupId}}
	result, err := registrationCollection.UpdateOne(context.Background(), bson.M{"_id": registrationId}, update)
	if err != nil {
		log.Printf("Error updating registration: %v", err)
		writeJSONError(w, "Ошибка при обновлении заявки", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		writeJSONError(w, "Заявка не найдена", http.StatusNotFound)
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

	objectId, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)
	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"groupId": objectId},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userId",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":      1,
				"username": bson.M{"$arrayElemAt": bson.A{"$user.username", 0}},
				"email":    bson.M{"$arrayElemAt": bson.A{"$user.email", 0}},
				"status":   1,
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

	for i := range members {
		if id, ok := members[i]["_id"].(primitive.ObjectID); ok {
			members[i]["_id"] = id.Hex()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"members": members})
}

func EnrollGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId, err := primitive.ObjectIDFromHex(vars["groupId"])
	if err != nil {
		writeJSONError(w, "Неверный формат ID группы", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		OrderID primitive.ObjectID `json:"orderId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": requestBody.OrderID}).Decode(&order)
	if err != nil {
		writeJSONError(w, "Приказ не найден", http.StatusBadRequest)
		return
	}
	if order.OrderType != "О зачислении обучающихся" {
		writeJSONError(w, "Приказ должен быть типа 'О зачислении обучающихся'", http.StatusBadRequest)
		return
	}

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	filter := bson.M{
		"groupId": groupId,
		"status":  "Оплаченный",
	}
	update := bson.M{
		"$set": bson.M{
			"status":        "Проходит курс",
			"enrollOrderId": requestBody.OrderID,
			"enrollDate":    time.Now(),
		},
	}

	result, err := registrationCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении заявок группы", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		writeJSONError(w, "Нет заявок со статусом 'Оплаченный' для зачисления", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func ExpelGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId, err := primitive.ObjectIDFromHex(vars["groupId"])
	if err != nil {
		writeJSONError(w, "Неверный формат ID группы", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		OrderID primitive.ObjectID `json:"orderId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": requestBody.OrderID}).Decode(&order)
	if err != nil {
		writeJSONError(w, "Приказ не найден", http.StatusBadRequest)
		return
	}

	var newStatus string
	var documentType string
	if order.OrderType == "О выпуске обучающихся" {
		newStatus = "Завершил"
		documentType = "Диплом"
	} else if order.OrderType == "Об отчислении обучающихся" {
		newStatus = "Отчисленный"
		documentType = "Сертификат"
	} else {
		writeJSONError(w, "Приказ должен быть типа 'О выпуске обучающихся' или 'Об отчислении обучающихся'", http.StatusBadRequest)
		return
	}

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	filter := bson.M{
		"groupId": groupId,
		"status":  "Проходит курс",
	}
	update := bson.M{
		"$set": bson.M{
			"status":       newStatus,
			"expelOrderId": requestBody.OrderID,
			"expelDate":    time.Now(),
			"documentType": documentType,
		},
	}

	result, err := registrationCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении заявок группы", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		writeJSONError(w, "Нет заявок со статусом 'Проходит курс' для обработки", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"status":       newStatus,
		"documentType": documentType,
	})
}

func ExpelRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат ID заявки", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		OrderID primitive.ObjectID `json:"orderId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": requestBody.OrderID}).Decode(&order)
	if err != nil {
		writeJSONError(w, "Приказ не найден", http.StatusBadRequest)
		return
	}
	if order.OrderType != "Об отчислении обучающихся" {
		writeJSONError(w, "Приказ должен быть типа 'Об отчислении обучающихся'", http.StatusBadRequest)
		return
	}

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	var registration bson.M
	err = registrationCollection.FindOne(context.Background(), bson.M{"_id": registrationId}).Decode(&registration)
	if err != nil {
		writeJSONError(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	if registration["status"] != "Проходит курс" {
		writeJSONError(w, "Отчисление возможно только для статуса 'Проходит курс'", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": registrationId}
	update := bson.M{
		"$set": bson.M{
			"status":       "Отчисленный",
			"expelOrderId": requestBody.OrderID,
			"expelDate":    time.Now(),
			"documentType": "Сертификат",
		},
	}

	_, err = registrationCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при отчислении участника", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func ExpelMultipleRegistrations(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		RegistrationIDs []string   `json:"registrationIds"`
		OrderID         string `json:"orderId"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if len(requestBody.RegistrationIDs) == 0 {
		writeJSONError(w, "Не указаны ID заявок для отчисления", http.StatusBadRequest)
		return
	}

	orderId, err := primitive.ObjectIDFromHex(requestBody.OrderID)
	if err != nil {
		writeJSONError(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": orderId}).Decode(&order)
	if err != nil {
		writeJSONError(w, "Приказ не найден", http.StatusBadRequest)
		return
	}
	if order.OrderType != "Об отчислении обучающихся" {
		writeJSONError(w, "Приказ должен быть типа 'Об отчислении обучающихся'", http.StatusBadRequest)
		return
	}

	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	var registrationIds []primitive.ObjectID
	for _, id := range requestBody.RegistrationIDs {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			writeJSONError(w, fmt.Sprintf("Неверный формат ID заявки: %s", id), http.StatusBadRequest)
			return
		}
		var registration bson.M
		err = registrationCollection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&registration)
		if err != nil {
			writeJSONError(w, fmt.Sprintf("Заявка %s не найдена", id), http.StatusNotFound)
			return
		}
		if registration["status"] != "Проходит курс" {
			writeJSONError(w, fmt.Sprintf("Отчисление невозможно для заявки %s: статус не 'Проходит курс'", id), http.StatusBadRequest)
			return
		}
		registrationIds = append(registrationIds, objId)
	}

	filter := bson.M{
		"_id": bson.M{"$in": registrationIds},
	}
	update := bson.M{
		"$set": bson.M{
			"status":       "Отчисленный",
			"expelOrderId": orderId,
			"expelDate":    time.Now(),
			"documentType": "Сертификат",
		},
	}

	result, err := registrationCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при отчислении участников", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   result.ModifiedCount,
	})
}

func CompleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId, err := primitive.ObjectIDFromHex(vars["groupId"])
	if err != nil {
		writeJSONError(w, "Неверный ID группы", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.GroupsCollection)
	var group bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": groupId}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Группа не найдена", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Ошибка при поиске группы", http.StatusInternalServerError)
		return
	}

	if group["status"] == "Завершена" {
		writeJSONError(w, "Группа уже завершена", http.StatusBadRequest)
		return
	}

	regCollection := db.GetCollection(db.CourseRegistrationsCollection)
	cursor, err := regCollection.Find(context.Background(), bson.M{
		"groupId": groupId,
		"status":  "Проходит курс",
	})
	if err != nil {
		writeJSONError(w, "Ошибка при поиске участников группы", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var registrations []bson.M
	if err = cursor.All(context.Background(), &registrations); err != nil {
		writeJSONError(w, "Ошибка при обработке участников", http.StatusInternalServerError)
		return
	}

	if len(registrations) == 0 {
		writeJSONError(w, "В группе нет активных участников со статусом 'Проходит курс'", http.StatusBadRequest)
		return
	}

	_, err = regCollection.UpdateMany(
		context.Background(),
		bson.M{
			"groupId": groupId,
			"status":  "Проходит курс",
		},
		bson.M{
			"$set": bson.M{
				"status":       "Завершил",
				"documentType": "Диплом",
			},
		},
	)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении статуса участников", http.StatusInternalServerError)
		return
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": groupId},
		bson.M{"$set": bson.M{"status": "Завершена"}},
	)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении статуса группы", http.StatusInternalServerError)
		return
	}

	orderCollection := db.GetCollection(db.OrderCollection)
	order := bson.M{
		"_id":       primitive.NewObjectID(),
		"groupId":   groupId,
		"orderType": "О выпуске обучающихся",
		"createdAt": time.Now(),
		"status":    "Выполнен",
	}
	_, err = orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		writeJSONError(w, "Ошибка при создании приказа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Группа успешно завершена",
	})
}