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
	"time"

	"github.com/Renal37/db"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DownloadContract - обработчик для скачивания договора (возвращает JSON с данными)
func DownloadContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Неверный courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	log.Printf("Получение данных регистрации для courseId: %s", courseId.Hex())
	user, course, err := getRegistrationData(courseId)
	if err != nil {
		log.Printf("Ошибка получения данных регистрации: %v", err)
		http.Error(w, "Ошибка получения данных курса", http.StatusInternalServerError)
		return
	}

	data := prepareContractData(user, course)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	log.Printf("Данные договора успешно отправлены для courseId: %s", courseId.Hex())
}

// GetContractData - эндпоинт для возврата данных договора
func GetContractData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Неверный courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	log.Printf("Получение данных договора для courseId: %s", courseId.Hex())
	user, course, err := getRegistrationData(courseId)
	if err != nil {
		log.Printf("Ошибка получения данных регистрации: %v", err)
		http.Error(w, "Ошибка получения данных курса", http.StatusInternalServerError)
		return
	}

	data := prepareContractData(user, course)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	log.Printf("Данные договора успешно отправлены для courseId: %s", courseId.Hex())
}

// GetContractTemplate - эндпоинт для возврата шаблона PDF
func GetContractTemplate(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Ошибка получения текущей директории: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	pdfPath := filepath.Join(cwd, "document_download", "ДОГОВОР123.pdf")
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		log.Printf("Файл шаблона PDF не существует: %s", pdfPath)
		http.Error(w, "Шаблон PDF не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=ДОГОВОР123.pdf")
	http.ServeFile(w, r, pdfPath)
	log.Printf("Шаблон PDF успешно отправлен")
}

// getRegistrationData - получение данных пользователя и курса
func getRegistrationData(courseId primitive.ObjectID) (bson.M, bson.M, error) {
	collection := db.GetCollection(db.CourseRegistrationsCollection)
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
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка агрегации: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil || len(results) == 0 {
		return nil, nil, fmt.Errorf("данные не найдены")
	}

	result := results[0]
	user := result["user"].(primitive.A)[0].(primitive.M)
	course := result["course"].(primitive.A)[0].(primitive.M)

	return user, course, nil
}

// prepareContractData - подготовка данных для договора
func prepareContractData(user, course bson.M) map[string]string {
	data := make(map[string]string)

	// ФИО
	if lastName, ok := user["lastname"].(string); ok {
		firstName, _ := user["firstname"].(string)
		middleName, _ := user["middlename"].(string)
		data["FullName"] = fmt.Sprintf("%s %s %s", lastName, firstName, middleName)
	}

	// Данные курса
	if title, ok := course["title"].(string); ok {
		data["CourseTitle"] = title
	}

	switch v := course["duration"].(type) {
	case int, int32, int64:
		data["CourseDuration"] = fmt.Sprintf("%d часов", v)
	case float64:
		data["CourseDuration"] = fmt.Sprintf("%.0f часов", v)
	}

	switch v := course["price"].(type) {
	case int, int32, int64:
		data["CoursePrice"] = fmt.Sprintf("%d руб.", v)
	case float64:
		data["CoursePrice"] = fmt.Sprintf("%.0f руб.", v)
	}

	// Номер и дата договора
	data["DocumentId"] = time.Now().Format("02")
	data["DocumentDay"] = time.Now().Format("02.01.2006") // Например, 14.05.2025

	// Личные данные
	if address, ok := user["homeaddress"].(string); ok {
		data["Adress"] = address
	}
	if passport, ok := user["passportdata"].(string); ok {
		data["PasportDate"] = passport
	}
	if snils, ok := user["snils"].(string); ok {
		data["SNILS"] = snils
	}
	if phone, ok := user["phone"].(string); ok {
		data["Phone"] = phone
	}
	if email, ok := user["email"].(string); ok {
		data["Email"] = email
	}

	// Новые поля (заглушки, замени на реальные данные из базы)
	data["CourseEnd"] = "31.12.2025" // Дата окончания курса
	data["Time"] = "10:00-14:00"     // Время занятий
	data["TimeDayStart"] = "15"      // День начала
	data["TimeMonthStart"] = "05"    // Месяц начала
	data["TimeDayEnd"] = "31"        // День окончания
	data["TimeMonthEnd"] = "12"      // Месяц окончания
	data["HowGive"] = "Очно"         // Формат обучения

	// Логирование данных
	for field, value := range data {
		log.Printf("Поле %s: %s", field, value)
	}

	return data
}

// UploadContract - загрузка договора
func UploadContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Неверный courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)
	filter := bson.M{"_id": courseId}
	var courseRegistration bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&courseRegistration)
	if err != nil {
		log.Printf("Курс не найден: %v", err)
		http.Error(w, "Курс не найден", http.StatusNotFound)
		return
	}

	file, handler, err := r.FormFile("contract")
	if err != nil {
		log.Printf("Ошибка получения файла: %v", err)
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dir := filepath.Join("document")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
	filePath := filepath.Join(dir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Ошибка создания файла: %v", err)
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Ошибка сохранения файла: %v", err)
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"contractUploaded": true,
			"contractFilePath": filePath,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Ошибка обновления базы данных: %v", err)
		http.Error(w, "Ошибка при обновлении данных", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "Договор успешно загружен!"}`)
}

// ApproveContract - утверждение договора
func ApproveContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	contractId, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		log.Printf("Неверный contractId: %v", err)
		http.Error(w, "Неверный формат идентификатора договора", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)
	filter := bson.M{"_id": contractId}
	update := bson.M{"$set": bson.M{"status": "Принят"}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Ошибка обновления статуса договора: %v", err)
		http.Error(w, "Ошибка при обновления статуса договора", http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Договор не найден или статус уже обновлен", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"success": true}`)
}

// ViewContract - просмотр договора
func ViewContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		log.Printf("Неверный registration ID: %v", err)
		http.Error(w, "Неверный формат идентификатора заявки", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)
	var registration bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": registrationID}).Decode(&registration)
	if err != nil {
		log.Printf("Запись о регистрации не найдена: %v", err)
		http.Error(w, "Запись о регистрации не найдена", http.StatusNotFound)
		return
	}

	filePath, ok := registration["contractFilePath"].(string)
	if !ok || filePath == "" {
		log.Println("Путь к файлу договора не найден")
		http.Error(w, "Файл договора не найден", http.StatusNotFound)
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("Файл не существует: %v", err)
		http.Error(w, "Файл договора не найден на сервере", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=contract.pdf")
	http.ServeFile(w, r, filePath)
}
