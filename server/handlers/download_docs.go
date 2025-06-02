package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

// DownloadDocs - обработчик для скачивания диплома или сертификата
func DownloadDocs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Неверный courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	log.Printf("Получение данных регистрации для courseId: %s", courseId.Hex())
	user, course, registration, err := getRegistrationDataWithStatus(courseId)
	if err != nil {
		log.Printf("Ошибка получения данных регистрации: %v", err)
		http.Error(w, "Ошибка получения данных курса", http.StatusInternalServerError)
		return
	}

	// Проверка статуса
	status := registration["status"].(string)
	if status != "Завершил" && status != "Отчисленный" {
		log.Printf("Недопустимый статус для скачивания документа: %s", status)
		http.Error(w, "Документ доступен только для статусов 'Завершил' или 'Отчисленный'", http.StatusBadRequest)
		return
	}

	data := prepareDocumentData(user, course, registration)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	log.Printf("Данные документа успешно отправлены для courseId: %s", courseId.Hex())
}

// GetDocsTemplate - эндпоинт для возврата шаблона PDF (диплом или сертификат)
func GetDocsTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Неверный courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	// Получаем статус регистрации
	_, _, registration, err := getRegistrationDataWithStatus(courseId)
	if err != nil {
		log.Printf("Ошибка получения данных регистрации: %v", err)
		http.Error(w, "Ошибка получения данных курса", http.StatusInternalServerError)
		return
	}

	// Определяем тип документа в зависимости от статуса
	status := registration["status"].(string)
	var pdfFileName string
	if status == "Завершил" {
		pdfFileName = "Диплом.pdf"
	} else if status == "Отчисленный" {
		pdfFileName = "Сертефикат.pdf"
	} else {
		log.Printf("Недопустимый статус для шаблона: %s", status)
		http.Error(w, "Документ доступен только для статусов 'Завершил' или 'Отчисленный'", http.StatusBadRequest)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Ошибка получения текущей директории: %v", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	pdfPath := filepath.Join(cwd, "document_download", pdfFileName)
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		log.Printf("Файл шаблона PDF не существует: %s", pdfPath)
		http.Error(w, "Шаблон PDF не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", pdfFileName))
	http.ServeFile(w, r, pdfPath)
	log.Printf("Шаблон PDF (%s) успешно отправлен", pdfFileName)
}

// getRegistrationDataWithStatus - получение данных пользователя, курса и регистрации
func getRegistrationDataWithStatus(courseId primitive.ObjectID) (bson.M, bson.M, bson.M, error) {
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
		return nil, nil, nil, fmt.Errorf("ошибка агрегации: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil || len(results) == 0 {
		return nil, nil, nil, fmt.Errorf("данные не найдены")
	}

	result := results[0]
	user := result["user"].(primitive.A)[0].(primitive.M)
	course := result["course"].(primitive.A)[0].(primitive.M)
	registration := result

	return user, course, registration, nil
}

// prepareDocumentData - подготовка данных для диплома или сертификата
func prepareDocumentData(user, course, registration bson.M) map[string]string {
	data := make(map[string]string)

	// Полное имя
	if lastName, ok := user["lastname"].(string); ok {
		firstName, _ := user["firstname"].(string)
		middleName, _ := user["middlename"].(string)
		data["Full_Name"] = fmt.Sprintf("%s %s %s", lastName, firstName, middleName)
	}

	// Данные курса
	if title, ok := course["title"].(string); ok {
		data["Course_Title"] = title
	}

	// Даты начала и окончания курса
	if registrationStart, ok := course["registrationStart"].(primitive.DateTime); ok {
		data["Start_course"] = registrationStart.Time().Format("02.01.2006")
	}
	if registrationEnd, ok := course["registrationEnd"].(primitive.DateTime); ok {
		data["End_course"] = registrationEnd.Time().Format("02.01.2006")
	}

	// Дата отчисления (для сертификата) или завершения (для диплома)
	if expelDate, ok := registration["expelDate"].(primitive.DateTime); ok && expelDate.Time().Unix() != 0 {
		data["Date"] = expelDate.Time().Format("02.01.2006")
	} else {
		data["Date"] = time.Now().Format("02.01.2006") // Если expelDate отсутствует, используем текущую дату
	}

	// Случайный номер документа
	data["Random_number"] = fmt.Sprintf("%d", time.Now().Nanosecond()%1000000) // Пример случайного числа

	// Логирование данных
	for field, value := range data {
		log.Printf("Поле %s: %s", field, value)
	}

	return data
}