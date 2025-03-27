package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Renal37/db"
	// "github.com/Renal37/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

// Price структура для хранения информации о стоимости
type Price struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Amount      int                `bson:"amount" json:"amount"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	Description string             `bson:"description" json:"description"`
}

// AddPrice добавляет новую стоимость
func AddPrice(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	price := Price{
		Amount:      request.Amount,
		Description: request.Description,
		CreatedAt:   time.Now(),
	}

	collection := db.GetCollection(db.PricesCollection)
	result, err := collection.InsertOne(context.Background(), price)
	if err != nil {
		sendError(w, "Ошибка при добавлении стоимости", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"price": map[string]interface{}{
			"id":          result.InsertedID,
			"amount":      price.Amount,
			"description": price.Description,
			"createdAt":   price.CreatedAt,
		},
	})
}

// GetPrices возвращает список всех стоимостей, отсортированных по дате (новые сначала)
func GetPrices(w http.ResponseWriter, r *http.Request) {
	
	collection := db.GetCollection(db.PricesCollection)

	// Сортировка по дате создания (новые сначала)
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		http.Error(w, "Ошибка при получении стоимостей", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var prices []bson.M
	if err = cursor.All(context.Background(), &prices); err != nil {
		http.Error(w, "Ошибка при обработке данных стоимостей", http.StatusInternalServerError)
		return
	}

	// Преобразуем ObjectID в строку для корректного отображения в JSON
	for i := range prices {
		prices[i]["_id"] = prices[i]["_id"].(primitive.ObjectID).Hex()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

// UpdatePrice обновляет существующую стоимость
func UpdatePrice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		sendError(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.PricesCollection)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"amount":      request.Amount,
			"description": request.Description,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		sendError(w, "Ошибка при обновлении стоимости", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		sendError(w, "Стоимость не найдена", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"message": "Стоимость успешно обновлена",
	})
}

// DeletePrice удаляет стоимость
func DeletePrice(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		sendError(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.PricesCollection)
	filter := bson.M{"_id": id}

	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		sendError(w, "Ошибка при удалении стоимости", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		sendError(w, "Стоимость не найдена", http.StatusNotFound)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"message": "Стоимость успешно удалена",
	})
}

// BulkUpdatePrices массово обновляет стоимости
func BulkUpdatePrices(w http.ResponseWriter, r *http.Request) {

	var request struct {
		PriceIDs []string `json:"priceIds"`
		Percent  float64  `json:"percent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	var objectIDs []primitive.ObjectID
	for _, id := range request.PriceIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			sendError(w, fmt.Sprintf("Неверный формат ID: %s", id), http.StatusBadRequest)
			return
		}
		objectIDs = append(objectIDs, objID)
	}

	collection := db.GetCollection(db.PricesCollection)
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}
	update := bson.M{
		"$mul": bson.M{
			"amount": 1 + request.Percent/100,
		},
	}

	result, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		sendError(w, "Ошибка при массовом обновлении стоимостей", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Обновлено %d стоимостей", result.ModifiedCount),
	})
}


func sendJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
