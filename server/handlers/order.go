package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddOrderType(w http.ResponseWriter, r *http.Request) {
	var orderType models.OrderType
	if err := json.NewDecoder(r.Body).Decode(&orderType); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат данных"})
		return
	}

	collection := db.GetCollection(db.OrderTypesCollection)

	_, err := collection.InsertOne(context.Background(), orderType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка добавления типа приказа"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderType)
}

func GetOrderTypes(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.OrderTypesCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка получения типов приказов"})
		return
	}
	defer cursor.Close(context.Background())

	var orderTypes []models.OrderType
	if err = cursor.All(context.Background(), &orderTypes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка обработки данных"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderTypes)
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	var orderData struct {
		Number      string `json:"number"`
		Date        string `json:"date"`
		OrderTypeID string `json:"orderTypeId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&orderData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат данных"})
		return
	}

	// Парсим дату
	date, err := time.Parse("2006-01-02", orderData.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат даты, используйте ГГГГ-ММ-ДД"})
		return
	}

	// Преобразуем orderTypeId в ObjectID
	orderTypeObjID, err := primitive.ObjectIDFromHex(orderData.OrderTypeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат orderTypeId"})
		return
	}

	// Получаем тип приказа
	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)
	var orderType models.OrderType
	err = orderTypeCollection.FindOne(context.Background(), bson.M{"_id": orderTypeObjID}).Decode(&orderType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Тип приказа не найден"})
		return
	}

	// Создаем полную модель Order
	fullOrder := models.Order{
		Number:      orderData.Number,
		Date:        date,
		OrderTypeID: orderTypeObjID,
		OrderType:   orderType.Name,
	}

	collection := db.GetCollection(db.OrderCollection)
	fullOrder.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.Background(), fullOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка добавления приказа: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fullOrder)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.OrderCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка получения приказов"})
		return
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err = cursor.All(context.Background(), &orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка обработки данных"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func DeleteOrderType(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Отсутствует параметр id"})
		return
	}

	orderTypeID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат ID типа приказа"})
		return
	}

	orderCollection := db.GetCollection(db.OrderCollection)
	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)

	_, err = orderCollection.DeleteMany(context.Background(), bson.M{"orderTypeId": orderTypeID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка удаления связанных приказов"})
		return
	}

	result, err := orderTypeCollection.DeleteOne(context.Background(), bson.M{"_id": orderTypeID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка удаления типа приказа"})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Тип приказа не найден"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Тип приказа и связанные приказы успешно удалены"})
}

func UpdateOrderType(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Отсутствует параметр id"})
		return
	}

	orderTypeID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат ID типа приказа"})
		return
	}

	var orderType models.OrderType
	if err := json.NewDecoder(r.Body).Decode(&orderType); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат данных"})
		return
	}

	collection := db.GetCollection(db.OrderTypesCollection)
	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": orderTypeID},
		bson.M{"$set": bson.M{"name": orderType.Name}},
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка обновления типа приказа"})
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Тип приказа не найден"})
		return
	}

	var updatedOrderType models.OrderType
	err = collection.FindOne(context.Background(), bson.M{"_id": orderTypeID}).Decode(&updatedOrderType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка получения обновленного типа приказа"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrderType)
}

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Отсутствует параметр id"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат ID приказа"})
		return
	}

	var updateData struct {
		Number      string `json:"number"`
		Date        string `json:"date"`
		OrderTypeID string `json:"orderTypeId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат данных"})
		return
	}

	date, err := time.Parse("2006-01-02", updateData.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат даты, используйте ГГГГ-ММ-ДД"})
		return
	}

	orderTypeID, err := primitive.ObjectIDFromHex(updateData.OrderTypeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат orderTypeId"})
		return
	}

	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)
	var orderType models.OrderType
	err = orderTypeCollection.FindOne(context.Background(), bson.M{"_id": orderTypeID}).Decode(&orderType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Тип приказа не найден"})
		return
	}

	collection := db.GetCollection(db.OrderCollection)
	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{
			"number":      updateData.Number,
			"date":        date,
			"orderTypeId": orderTypeID,
			"orderType":   orderType.Name,
		}},
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка обновления приказа"})
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Приказ не найден"})
		return
	}

	var updatedOrder models.Order
	err = collection.FindOne(context.Background(), bson.M{"_id": orderID}).Decode(&updatedOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка получения обновленного приказа"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Отсутствует параметр id"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Неверный формат ID приказа"})
		return
	}

	collection := db.GetCollection(db.OrderCollection)
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": orderID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка удаления приказа"})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Приказ не найден"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Приказ успешно удален"})
}
