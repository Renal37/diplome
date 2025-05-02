package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Renal37/db"
	"github.com/gorilla/mux"
	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DownloadContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Invalid courseId: %v", err)
		http.Error(w, "Invalid course ID format", http.StatusBadRequest)
		return
	}

	registration, user, course, err := getRegistrationData(courseId)
	if err != nil {
		log.Printf("Error getting registration data: %v", err)
		http.Error(w, "Error getting course data", http.StatusInternalServerError)
		return
	}

	pdfBytes, err := generateFilledContract(registration, user, course)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, "Error generating contract", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract_"+courseId.Hex()+".pdf")
	w.Write(pdfBytes)
}

func getRegistrationData(courseId primitive.ObjectID) (bson.M, bson.M, bson.M, error) {
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
		return nil, nil, nil, fmt.Errorf("aggregation error: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil || len(results) == 0 {
		return nil, nil, nil, fmt.Errorf("no data found")
	}

	result := results[0]
	user := result["user"].(primitive.A)[0].(primitive.M)
	course := result["course"].(primitive.A)[0].(primitive.M)

	return result, user, course, nil
}

func generateFilledContract(registration, user, course bson.M) ([]byte, error) {
	// Открываем файл шаблона
	file, err := os.Open("../server/document_donwload/ДОГОВОР_fix.pdf")
	if err != nil {
		return nil, fmt.Errorf("could not open template file: %v", err)
	}
	defer file.Close()

	// Создаем PDF reader
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return nil, fmt.Errorf("could not create PDF reader: %v", err)
	}

	// Получаем форму из PDF
	acroForm := pdfReader.AcroForm
	if acroForm == nil {
		return nil, fmt.Errorf("PDF template has no form fields")
	}

	// Заполняем поля формы
	fields := acroForm.AllFields()
	for _, field := range fields {
		fieldName, err := field.FullName()
		if err != nil {
			continue
		}

		var value string
		switch fieldName {
		case "FullName":
			value = fmt.Sprintf("%s %s %s",
				user["lastName"], user["firstName"], user["middleName"])
		case "CourseName":
			value = course["title"].(string)
		case "CourseDuration":
			value = fmt.Sprintf("%d часов", course["duration"])
		case "CoursePrice":
			value = fmt.Sprintf("%.2f руб.", course["price"])
		case "SignDate":
			value = time.Now().Format("02.01.2006")
		case "Address":
			if address, ok := user["address"].(string); ok {
				value = address
			}
		case "Passport":
			if passport, ok := user["passport"].(string); ok {
				value = passport
			}
		}

		// Устанавливаем значение поля
		if value != "" {
			field.V = core.MakeString(value)
		}
	}

	// Создаем новый PDF writer
	var buf bytes.Buffer
	writer := model.NewPdfWriter()

	// Копируем все страницы из исходного PDF
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("could not get page count: %v", err)
	}

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return nil, fmt.Errorf("could not get page %d: %v", i, err)
		}

		if err := writer.AddPage(page); err != nil {
			return nil, fmt.Errorf("could not add page %d: %v", i, err)
		}
	}

	// Устанавливаем заполненную форму
	if err := writer.SetForms(acroForm); err != nil {
		return nil, fmt.Errorf("could not set form: %v", err)
	}

	// Записываем в буфер
	if err := writer.Write(&buf); err != nil {
		return nil, fmt.Errorf("could not write PDF: %v", err)
	}

	return buf.Bytes(), nil
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
	collection := db.GetCollection(db.CourseRegistrationsCollection)
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

	collection := db.GetCollection(db.CourseRegistrationsCollection)
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
	collection := db.GetCollection(db.CourseRegistrationsCollection)

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
