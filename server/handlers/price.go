package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddPrice(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if request.Amount <= 0 {
		sendError(w, "Сумма должна быть положительной", http.StatusBadRequest)
		return
	}

	now := time.Now()
	price := models.Price{
		Amount:      request.Amount,
		Description: request.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
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
			"updatedAt":   price.UpdatedAt,
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

	if request.Amount <= 0 {
		sendError(w, "Сумма должна быть положительной", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.PricesCollection)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"amount":      request.Amount,
			"description": request.Description,
			"updatedAt":   time.Now(),
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
		PriceIds []string `json:"priceIds"`
		Percent  float64  `json:"percent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if request.Percent <= 0 {
		sendError(w, "Процент должен быть положительным", http.StatusBadRequest)
		return
	}

	if len(request.PriceIds) == 0 {
		sendError(w, "Не выбраны стоимости для обновления", http.StatusBadRequest)
		return
	}

	// Convert priceIds to ObjectIDs
	var objectIds []primitive.ObjectID
	for _, id := range request.PriceIds {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			sendError(w, fmt.Sprintf("Неверный формат ID: %s", id), http.StatusBadRequest)
			return
		}
		objectIds = append(objectIds, objId)
	}

	collection := db.GetCollection(db.PricesCollection)

	// Fetch current prices to calculate new amounts
	cursor, err := collection.Find(context.Background(), bson.M{"_id": bson.M{"$in": objectIds}})
	if err != nil {
		sendError(w, "Ошибка при получении стоимостей", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var prices []struct {
		ID     primitive.ObjectID `bson:"_id"`
		Amount int                `bson:"amount"`
	}
	if err = cursor.All(context.Background(), &prices); err != nil {
		sendError(w, "Ошибка при обработке стоимостей", http.StatusInternalServerError)
		return
	}

	if len(prices) == 0 {
		sendError(w, "Ни одна стоимость не найдена", http.StatusNotFound)
		return
	}

	// Update each price
	now := time.Now()
	updatedCount := 0
	for _, price := range prices {
		newAmount := int(float64(price.Amount) * (1 + request.Percent/100))
		if newAmount <= 0 {
			sendError(w, "Новая сумма должна быть положительной", http.StatusBadRequest)
			return
		}

		filter := bson.M{"_id": price.ID}
		update := bson.M{
			"$set": bson.M{
				"amount":    newAmount,
				"updatedAt": now,
			},
		}

		result, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			sendError(w, "Ошибка при обновлении стоимости", http.StatusInternalServerError)
			return
		}
		if result.ModifiedCount > 0 {
			updatedCount++
		}
	}

	sendJSON(w, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("%d стоимостей успешно обновлено", updatedCount),
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
