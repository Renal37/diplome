package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FillConsent(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf("Токен отсутствует: %v", err)
		http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		log.Printf("Неверный токен: %v", err)
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		log.Printf("Неверный формат ID пользователя: %v", err)
		http.Error(w, "Неверный формат ID пользователя", http.StatusBadRequest)
		return
	}

	userCollection := db.GetCollection(db.UsersCollection)
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		log.Printf("Пользователь не найден: %v", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	consentData := map[string]interface{}{
		"FullName":          fmt.Sprintf("%s %s %s", user.LastName, user.FirstName, user.MiddleName),
		"PassportData":      user.PassportData,
		"PassportIssuedBy":  user.PassportIssuedBy,
		"PassportIssueDate": user.PassportIssueDate,
		"HomeAddress":       user.HomeAddress,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(consentData); err != nil {
		log.Printf("Ошибка при формировании ответа: %v", err)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}
func DownloadDocument(w http.ResponseWriter, r *http.Request) {
	templatePath := "./document_download/согласие1.pdf"
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("Файл не найден: %v", err)
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	fileBytes, err := os.ReadFile(templatePath)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		http.Error(w, "Ошибка при чтении файла", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=согласие на обработку данных.pdf")
	w.Write(fileBytes)
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
	collection := db.GetCollection(db.UsersCollection)

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
	collection := db.GetCollection(db.UsersCollection)

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
