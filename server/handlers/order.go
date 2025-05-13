package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	// Ограничиваем размер файла (10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "File too large, max 10MB"})
		return
	}

	// Получаем данные формы
	number := r.FormValue("number")
	dateStr := r.FormValue("date")
	orderTypeID := r.FormValue("orderTypeId")
	file, handler, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error retrieving the file"})
		return
	}
	defer file.Close()

	// Проверяем тип файла
	if !strings.HasSuffix(handler.Filename, ".pdf") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Only PDF files are allowed"})
		return
	}

	// Парсим дату
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	// Преобразуем orderTypeId в ObjectID
	orderTypeObjID, err := primitive.ObjectIDFromHex(orderTypeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid orderTypeId format"})
		return
	}

	// Получаем тип приказа
	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)
	var orderType models.OrderType
	err = orderTypeCollection.FindOne(context.Background(), bson.M{"_id": orderTypeObjID}).Decode(&orderType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order type not found"})
		return
	}

	// Создаем директорию для файлов, если ее нет
	uploadDir := "./uploads/orders"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Генерируем уникальное имя файла
	fileExt := filepath.Ext(handler.Filename)
	newFileName := primitive.NewObjectID().Hex() + fileExt
	filePath := filepath.Join(uploadDir, newFileName)

	// Сохраняем файл
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error saving file"})
		return
	}
	defer f.Close()

	io.Copy(f, file)

	// Создаем полную модель Order
	fullOrder := models.Order{
		Number:      number,
		Date:        date,
		OrderTypeID: orderTypeObjID,
		OrderType:   orderType.Name,
		FileURL:     "/uploads/orders/" + newFileName,
		FileName:    handler.Filename,
	}

	collection := db.GetCollection(db.OrderCollection)
	fullOrder.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.Background(), fullOrder)
	if err != nil {
		// Удаляем сохраненный файл в случае ошибки
		os.Remove(filePath)
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

// Добавляем маршрут для просмотра файла приказа
func ServeOrderFile(w http.ResponseWriter, r *http.Request) {
	filePath := "." + r.URL.Path
	http.ServeFile(w, r, filePath)
}
func DeleteOrderType(w http.ResponseWriter, r *http.Request) {
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

	// Получаем коллекции
	orderCollection := db.GetCollection(db.OrderCollection)
	orderTypeCollection := db.GetCollection(db.OrderTypesCollection)

	// Удаляем все приказы этого типа
	_, err = orderCollection.DeleteMany(context.Background(), bson.M{"orderTypeId": orderTypeID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error deleting related orders"})
		return
	}

	// Удаляем сам тип приказа
	result, err := orderTypeCollection.DeleteOne(context.Background(), bson.M{"_id": orderTypeID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error deleting order type"})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order type not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order type and related orders deleted successfully"})
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
