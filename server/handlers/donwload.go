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

	// Находим запись о регистрации курса
	var registration bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": courseId}).Decode(&registration)
	if err != nil {
		log.Printf("Registration not found: %v", err)
		http.Error(w, "Запись о регистрации не найдена", http.StatusNotFound)
		return
	}

	// Получаем путь к файлу договора
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

	// Определяем режим отображения файла (просмотр или скачивание)
	displayMode := r.URL.Query().Get("mode") // mode=inline для просмотра, mode=attachment для скачивания

	// Устанавливаем заголовки ответа
	w.Header().Set("Content-Type", "application/pdf")
	if displayMode == "inline" {
		w.Header().Set("Content-Disposition", "inline; filename=contract.pdf")
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename=contract.pdf")
	}

	// Отправляем файл клиенту
	http.ServeFile(w, r, filePath)
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
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating contract status: %v", err)
		http.Error(w, "Ошибка при обновлении статуса договора", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Договор успешно принят!"}`)
}
