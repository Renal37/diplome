package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/phpdave11/gofpdi"
	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
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

	userArray, ok := result[0]["user"].(primitive.A)
	if !ok || len(userArray) == 0 {
		log.Println("User data not found or invalid format")
		http.Error(w, "Данные пользователя не найдены", http.StatusNotFound)
		return
	}
	user := userArray[0].(primitive.M)

	courseArray, ok := result[0]["course"].(primitive.A)
	if !ok || len(courseArray) == 0 {
		log.Println("Course data not found or invalid format")
		http.Error(w, "Данные курса не найдены", http.StatusNotFound)
		return
	}
	course := courseArray[0].(primitive.M)

	// Логирование данных для отладки
	log.Printf("User Data: %+v", user)
	log.Printf("Course Data: %+v", course)

	// Создание PDF-документа
	pdf := gofpdf.New("P", "mm", "A4", "")
	importer := gofpdi.NewImporter()

	// Загрузка шаблона PDF
	templatePath := "server/handlers/ДОГОВОР.pdf"
	pageNumber := 1 // Номер страницы шаблона для импорта

	// Импорт страницы
	tpl := importer.ImportPageFromFile(templatePath, pageNumber, "/MediaBox")

	// Добавление страницы
	pdf.AddPage()

	// Применение импортированного шаблона
	importer.UseImportedTemplate(pdf, tpl, 0, 0, 210, 297)

	// Установка шрифта
	pdf.SetFont("Arial", "", 12)

	// Заполнение полей PDF
	pdf.SetXY(50, 50) // Координаты для ФИО
	pdf.CellFormat(80, 10, fmt.Sprintf("%v %v %v", user["lastname"], user["firstname"], user["middlename"]), "", 0, "L", false, 0, "")

	pdf.SetXY(50, 60) // Координаты для Email
	pdf.CellFormat(80, 10, fmt.Sprintf("%v", user["email"]), "", 0, "L", false, 0, "")

	pdf.SetXY(50, 70) // Координаты для Телефона
	pdf.CellFormat(80, 10, fmt.Sprintf("%v", user["phone"]), "", 0, "L", false, 0, "")

	pdf.SetXY(50, 80) // Координаты для Курса
	pdf.CellFormat(80, 10, fmt.Sprintf("%v", course["title"]), "", 0, "L", false, 0, "")

	pdf.SetXY(50, 90) // Координаты для Продолжительности
	pdf.CellFormat(80, 10, fmt.Sprintf("%v часов", course["duration"]), "", 0, "L", false, 0, "")

	pdf.SetXY(50, 100) // Координаты для Цены
	pdf.CellFormat(80, 10, fmt.Sprintf("%v рублей", course["price"]), "", 0, "L", false, 0, "")

	// Отправка PDF-файла клиенту
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract.pdf")
	if err := pdf.Output(w); err != nil {
		log.Printf("PDF generation error: %v", err)
		http.Error(w, "Ошибка при генерации PDF", http.StatusInternalServerError)
	}
}
