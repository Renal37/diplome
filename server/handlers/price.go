package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Price структура для хранения информации о стоимости
type Price struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Amount      int                `bson:"amount" json:"amount"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	Description string             `bson:"description" json:"description"`
}

// AddPrice добавляет новую стоимость в базу данных
func AddPrice(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Amount      int    `json:"amount"`
		Description string `json:"description"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	price := Price{
		Amount:      request.Amount,
		Description: request.Description,
		CreatedAt:   time.Now(),
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("prices")
	result, err := collection.InsertOne(context.Background(), price)
	if err != nil {
		http.Error(w, "Ошибка при добавлении стоимости", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      result.InsertedID,
	})
}

// GetPrices возвращает список всех стоимостей, отсортированных по дате (новые сначала)
func GetPrices(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("prices")
	
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

// GetActivePrice возвращает последнюю добавленную стоимость
func GetActivePrice(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("prices")
	
	// Получаем последнюю добавленную стоимость
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"createdAt", -1}})
	
	var price Price
	err = collection.FindOne(context.Background(), bson.M{}, findOptions).Decode(&price)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Нет доступных стоимостей", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при получении стоимости", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(price)
}