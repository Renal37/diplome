package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/unidoc/unipdf/v3/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
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
	user := userArray[0].(primitive.M)

	// Обработка данных курса
	courseArray, ok := result[0]["course"].(primitive.A)
	if !ok || len(courseArray) == 0 {
		log.Println("Course data not found or invalid format")
		http.Error(w, "Данные курса не найдены", http.StatusNotFound)
		return
	}
	course := courseArray[0].(primitive.M)

	// Путь к шаблону PDF
	templatePath := "../server/handlers/ДОГОВОР_fix.pdf"
	outputPath := "../server/contract_output.pdf"

	// Проверка наличия интерактивных форм
	hasForms, err := checkInteractiveForms(templatePath)
	if err != nil {
		log.Printf("Error checking interactive forms: %v", err)
		http.Error(w, "Ошибка при проверке интерактивных форм", http.StatusInternalServerError)
		return
	}

	if hasForms {
		log.Println("Интерактивные формы найдены. Заполняем их.")
		err = fillInteractivePDF(templatePath, outputPath, map[string]string{
			"FullName":       fmt.Sprintf("%s %s %s", user["lastname"], user["firstname"], user["middlename"]),
			"Address":        fmt.Sprintf("%s", user["homeAddress"]),
			"Passport":       fmt.Sprintf("%s", user["passportData"]),
			"Phone":          fmt.Sprintf("%s", user["phone"]),
			"Email":          fmt.Sprintf("%s", user["email"]),
			"CourseName":     fmt.Sprintf("%s", course["title"]),
			"CourseDuration": fmt.Sprintf("%d часов", course["duration"]),
			"CoursePrice":    fmt.Sprintf("%d рублей", course["price"]),
		})
	} else {
		log.Println("Интерактивные формы отсутствуют. Добавляем текст поверх PDF.")
		// Здесь можно добавить альтернативную логику для добавления текста поверх PDF.
	}

	if err != nil {
		log.Printf("Error filling template: %v", err)
		http.Error(w, "Ошибка при заполнении шаблона", http.StatusInternalServerError)
		return
	}

	// Отправка PDF-файла клиенту
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract.pdf")
	http.ServeFile(w, r, outputPath)
}

// checkInteractiveForms проверяет наличие интерактивных форм в PDF
func checkInteractiveForms(filePath string) (bool, error) {
	// Загружаем PDF
	pdfReader, pdfErr, err := model.NewPdfReaderFromFile(filePath, nil)
	if err != nil {
		return false, fmt.Errorf("failed to load PDF: %w", err)
	}
	if pdfErr != nil {
		return false, fmt.Errorf("PDF read error: %v", pdfErr)
	}

	// Получаем AcroForm (интерактивные формы)
	form := pdfReader.AcroForm
	if form == nil {
		return false, nil
	}

	return true, nil
}

// fillInteractivePDF заполняет интерактивные формы в PDF
func fillInteractivePDF(inputPath, outputPath string, userData map[string]string) error {
	// Загружаем PDF
	pdfReader, pdfErr, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return fmt.Errorf("failed to load PDF: %w", err)
	}
	if pdfErr != nil {
		return fmt.Errorf("PDF read error: %v", pdfErr)
	}

	// Получаем AcroForm (интерактивные формы)
	form := pdfReader.AcroForm
	if form == nil {
		return fmt.Errorf("interactive forms not found in the PDF")
	}

	// Заполняем значения полей данными из userData
	for _, field := range *form.Fields {
		fieldName := field.T.String() // Получаем имя поля
		if fieldValue, exists := userData[fieldName]; exists {
			// Устанавливаем значение поля
			if field, ok := field.GetContext().(*model.PdfFieldText); ok {
				// Используем SetValue для установки значения
				field.V = model.MakeString(fieldValue)
			} else {
				log.Printf("Field '%s' is not a text field and cannot be set", fieldName)
			}
		}
	}

	// Создаем новый PDF с заполненными данными
	pdfWriter := model.NewPdfWriter()
	pdfWriter.SetForms(form)

	// Сохраняем результат
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	err = pdfWriter.Write(outputFile)
	if err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}