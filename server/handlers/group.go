package handlers

import (
	"context"
	"encoding/json"
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

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
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

	// Проверяем, что курс существует
	courseCollection := db.GetCollection(db.CoursesCollection)

	var course bson.M
	err = courseCollection.FindOne(context.Background(), bson.M{"_id": group.CourseID}).Decode(&course)
	if err != nil {
		log.Printf("Error finding course: %v", err)
		http.Error(w, "Курс не найден", http.StatusNotFound)
		return
	}

	// Создаем группу
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
		// Подсчет активных участников
		count, err := regCollection.CountDocuments(
			context.Background(),
			bson.M{"groupId": groupId, "status": bson.M{"$ne": "Отчислен"}},
		)
		if err != nil {
			writeJSONError(w, "Ошибка при подсчете участников", http.StatusInternalServerError)
			return
		}
		groups[i]["currentStudents"] = count

		// Установка статуса группы
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

	// Проверяем, есть ли участники в группе
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

	// Проверяем, есть ли участники со статусом "Проходит курс"
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
		http.Error(w, "Нельзя удалить группу, так как есть участники со статусом 'Проходит курс'", http.StatusBadRequest)
		return
	}

	// Если нет участников, удаляем группу и связанные заявки
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

	// Удаляем все заявки, связанные с этой группой
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

// Новый endpoint для отклонения записей группы
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

	// Проверяем, что нет участников со статусом "Проходит курс"
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

	// Обновляем статус заявок на "Отклоненный" с указанием причины
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

	// Получаем информацию о заявке
	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
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
	groupCollection := db.GetCollection(db.GroupsCollection)
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

	// Преобразуем строковый ID в ObjectID
	objectId, err := primitive.ObjectIDFromHex(groupId)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	// Получаем участников группы из коллекции course_registrations
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

	// Проверяем существование приказа
	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": requestBody.OrderID}).Decode(&order)
	if err != nil {
		writeJSONError(w, "Приказ не найден", http.StatusBadRequest)
		return
	}
	if order.OrderType != "О зачислении обучающихся" {
		writeJSONError(w, "Неверный тип приказа для зачисления", http.StatusBadRequest)
		return
	}

	// Обновляем статус всех заявок в группе
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
		writeJSONError(w, "Ошибка при зачислении группы", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		writeJSONError(w, "Нет подходящих для зачисления заявок", http.StatusBadRequest)
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

	// Проверяем существование приказа
	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": requestBody.OrderID}).Decode(&order)
	if err != nil {
		writeJSONError(w, "Приказ не найден", http.StatusBadRequest)
		return
	}

	// Определяем статус и тип документа
	var newStatus string
	var documentType string
	if order.OrderType == "О выпуске обучающихся" {
		newStatus = "Завершил"
		documentType = "Диплом"
	} else if order.OrderType == "Об отчислении обучающихся" {
		newStatus = "Отчисленный"
		documentType = "Сертификат"
	} else {
		writeJSONError(w, "Неверный тип приказа", http.StatusBadRequest)
		return
	}

	// Обновляем статус всех заявок в группе
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
		writeJSONError(w, "Ошибка при обработке", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		writeJSONError(w, "Нет подходящих для обработки заявок", http.StatusBadRequest)
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

	// Проверка статуса группы
	if group["status"] == "Завершена" {
		writeJSONError(w, "Группа уже завершена", http.StatusBadRequest)
		return
	}

	// Получение участников группы
	regCollection := db.GetCollection(db.CourseRegistrationsCollection)
	cursor, err := regCollection.Find(context.Background(), bson.M{"groupId": groupId, "status": bson.M{"$ne": "Отчислен"}})
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
		writeJSONError(w, "В группе нет активных участников", http.StatusBadRequest)
		return
	}

	// Обновление статуса участников
	_, err = regCollection.UpdateMany(
		context.Background(),
		bson.M{"groupId": groupId, "status": bson.M{"$ne": "Отчислен"}},
		bson.M{"$set": bson.M{"status": "Завершил"}},
	)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении статуса участников", http.StatusInternalServerError)
		return
	}

	// Обновление статуса группы
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": groupId},
		bson.M{"$set": bson.M{"status": "Завершена"}},
	)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении статуса группы", http.StatusInternalServerError)
		return
	}

	// Создание приказа
	orderCollection := db.GetCollection(db.OrderCollection)
	order := bson.M{
		"_id":       primitive.NewObjectID(),
		"groupId":   groupId,
		"orderType": "Завершение",
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
