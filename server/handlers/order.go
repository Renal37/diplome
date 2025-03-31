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
		// Возвращаем JSON вместо текста
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid data format"})
		return
	}

	collection := db.GetCollection(db.OrderTypesCollection)

	_, err := collection.InsertOne(context.Background(), orderType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error adding order type"})
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Error fetching order types"})
		return
	}
	defer cursor.Close(context.Background())

	var orderTypes []models.OrderType
	if err = cursor.All(context.Background(), &orderTypes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error processing data"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderTypes)
}

func AddOrder(w http.ResponseWriter, r *http.Request) {
	var order struct {
		Number      string `json:"number"`
		Date        string `json:"date"` // Получаем как строку
		OrderTypeID string `json:"orderTypeId"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid data format: " + err.Error()})
		return
	}

	// Преобразуем orderTypeId в ObjectID
	orderTypeID, err := primitive.ObjectIDFromHex(order.OrderTypeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid orderTypeId format"})
		return
	}

	// Парсим дату
	date, err := time.Parse("2006-01-02", order.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	// Получаем тип приказа
	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)
	var orderType models.OrderType
	err = orderTypeCollection.FindOne(context.Background(), bson.M{"_id": orderTypeID}).Decode(&orderType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order type not found"})
		return
	}

	// Создаем полную модель Order
	fullOrder := models.Order{
		Number:      order.Number,
		Date:        date,
		OrderTypeID: orderTypeID,
		OrderType:   orderType.Name,
		Description: order.Description,
	}

	collection := db.GetCollection(db.OrderCollection)
	fullOrder.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.Background(), fullOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error adding order: " + err.Error()})
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Error fetching orders"})
		return
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err = cursor.All(context.Background(), &orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error processing data"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func UpdateOrderType(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing id parameter"})
		return
	}

	orderTypeID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid order type ID format"})
		return
	}

	var orderType models.OrderType
	if err := json.NewDecoder(r.Body).Decode(&orderType); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid data format"})
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
		json.NewEncoder(w).Encode(map[string]string{"error": "Error updating order type"})
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order type not found"})
		return
	}

	// Получаем обновленный документ
	var updatedOrderType models.OrderType
	err = collection.FindOne(context.Background(), bson.M{"_id": orderTypeID}).Decode(&updatedOrderType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving updated order type"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrderType)
}
	

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing id parameter"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid order ID format"})
		return
	}

	var updateData struct {
		Number      string `json:"number"`
		Date        string `json:"date"`
		OrderTypeID string `json:"orderTypeId"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid data format"})
		return
	}

	// Парсим дату
	date, err := time.Parse("2006-01-02", updateData.Date)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	// Преобразуем orderTypeId в ObjectID
	orderTypeID, err := primitive.ObjectIDFromHex(updateData.OrderTypeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid orderTypeId format"})
		return
	}

	// Получаем имя типа приказа
	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)
	var orderType models.OrderType
	err = orderTypeCollection.FindOne(context.Background(), bson.M{"_id": orderTypeID}).Decode(&orderType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order type not found"})
		return
	}

	// Обновляем приказ
	collection := db.GetCollection(db.OrderCollection)
	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{
			"number":      updateData.Number,
			"date":        date,
			"orderTypeId": orderTypeID,
			"orderType":   orderType.Name,
			"description": updateData.Description,
		}},
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error updating order"})
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	// Получаем обновленный документ
	var updatedOrder models.Order
	err = collection.FindOne(context.Background(), bson.M{"_id": orderID}).Decode(&updatedOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving updated order"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из URL
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing id parameter"})
		return
	}

	orderID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid order ID format"})
		return
	}

	collection := db.GetCollection(db.OrderCollection)
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": orderID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error deleting order"})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order deleted successfully"})
}
