package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DownloadDocument(w http.ResponseWriter, r *http.Request) {
	// Путь к шаблону PDF
	templatePath := "../server/handlers/согласие на обработку ПД совершеннолетнего студента.pdf"

	// Проверяем, существует ли файл по указанному пути
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовки для скачивания файла
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=согласие на обработку данных.pdf")

	// Отправляем файл клиенту
	http.ServeFile(w, r, templatePath)
}
func UploadDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(vars["userId"])
	if err != nil {
		log.Printf("Invalid userId: %v", err)
		http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
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

	// Проверка существования пользователя
	collection := client.Database("diplome").Collection("users")
	filter := bson.M{"_id": userId}
	var user bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
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
func ViewConsent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := primitive.ObjectIDFromHex(vars["userId"])
	if err != nil {
		log.Printf("Invalid userId: %v", err)
		http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
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

	collection := client.Database("diplome").Collection("users")

	// Получаем путь к файлу согласия
	var user bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	filePath, ok := user["contractFilePath"].(string)
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
