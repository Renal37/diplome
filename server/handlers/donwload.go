package handlers

import (
	"context"
	"log"
	"net/http"

	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DownloadContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Invalid courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
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
	pipeline := bson.A{
		bson.M{"$match": bson.M{"_id": courseId}},
		bson.M{"$lookup": bson.M{
			"from":         "users",
			"localField":   "userId",
			"foreignField": "_id",
			"as":           "user",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "courses",
			"localField":   "courseId",
			"foreignField": "_id",
			"as":           "course",
		}},
		bson.M{"$project": bson.M{
			"user":   1,
			"course": 1,
			"status": 1,
		}},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Printf("Aggregation error: %v", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var result []bson.M
	if err = cursor.All(context.Background(), &result); err != nil || len(result) == 0 {
		log.Println("No data found")
		http.Error(w, "Данные не найдены", http.StatusNotFound)
		return
	}

	// Обработка данных пользователя
	userArray, ok := result[0]["user"].(primitive.A)
	if !ok || len(userArray) == 0 {
		log.Println("User data not found or invalid format")
		http.Error(w, "Данные пользователя не найдены", http.StatusNotFound)
		return
	}
	// user := userArray[0].(primitive.M)

	// Обработка данных курса
	courseArray, ok := result[0]["course"].(primitive.A)
	if !ok || len(courseArray) == 0 {
		log.Println("Course data not found or invalid format")
		http.Error(w, "Данные курса не найдены", http.StatusNotFound)
		return
	}
	// course := courseArray[0].(primitive.M)

	var usercourse []bson.M
	if err = cursor.All(context.Background(), &usercourse); err != nil {
		http.Error(w, "Ошибка при обработке данных заявок", http.StatusInternalServerError)
		return
	}

	// Путь к шаблону PDF
	templatePath := "../server/handlers/ДОГОВОР_fix.pdf"

	// Отправка PDF-шаблона клиенту
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract_template.pdf")
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, templatePath)
}
func UploadContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Invalid courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Проверка существования курса
	collection := client.Database("diplome").Collection("course_registrations")
	filter := bson.M{"_id": courseId}
	var courseRegistration bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&courseRegistration)
	if err != nil {
		log.Printf("Course not found: %v", err)
		http.Error(w, "Курс не найден", http.StatusNotFound)
		return
	}

	// Получение файла из запроса
	file, handler, err := r.FormFile("contract")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Сохранение файла на диск
	dir := "../server/document"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
	filePath := filepath.Join(dir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	// Обновление записи в базе данных
	update := bson.M{
		"$set": bson.M{
			"contractUploaded": true,
			"contractFilePath": filePath,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating database: %v", err)
		http.Error(w, "Ошибка при обновлении данных", http.StatusInternalServerError)
		return
	}

	// Ответ клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Договор успешно загружен!"}`)
}

func ApproveContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		log.Printf("Invalid contractId: %v", err)
		http.Error(w, "Неверный формат идентификатора договора", http.StatusBadRequest)
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
	filter := bson.M{"_id": contractId}
	update := bson.M{"$set": bson.M{"status": "Принят"}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating contract status: %v", err)
		http.Error(w, "Ошибка при обновлении статуса договора", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Договор не найден или статус уже обновлен", http.StatusNotFound)
		return
	}

	// Возвращаем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true}`)
}

func ViewContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		log.Printf("Invalid registration ID: %v", err)
		http.Error(w, "Неверный формат идентификатора заявки", http.StatusBadRequest)
		return
	}

	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("course_registrations")

	// Получаем путь к файлу договора
	var registration bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": registrationID}).Decode(&registration)
	if err != nil {
		log.Printf("Registration not found: %v", err)
		http.Error(w, "Запись о регистрации не найдена", http.StatusNotFound)
		return
	}

	filePath, ok := registration["contractFilePath"].(string)
	if !ok || filePath == "" {
		log.Println("Contract file path not found")
		http.Error(w, "Файл договора не найден", http.StatusNotFound)
		return
	}

	// Проверяем, существует ли файл на диске
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File does not exist: %v", err)
		http.Error(w, "Файл договора не найден на сервере", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовки для просмотра PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=contract.pdf")

	// Отправляем файл клиенту
	http.ServeFile(w, r, filePath)
}
