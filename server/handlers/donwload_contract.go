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
	"github.com/jung-kurt/gofpdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DownloadContract остается без изменений
func DownloadContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseId, err := primitive.ObjectIDFromHex(vars["courseId"])
	if err != nil {
		log.Printf("Неверный courseId: %v", err)
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	log.Printf("Получение данных регистрации для courseId: %s", courseId.Hex())
	registration, user, course, err := getRegistrationData(courseId)
	if err != nil {
		log.Printf("Ошибка получения данных регистрации: %v", err)
		http.Error(w, "Ошибка получения данных курса", http.StatusInternalServerError)
		return
	}

	log.Printf("Генерация PDF для courseId: %s", courseId.Hex())
	pdfBytes, err := generateFilledContract(registration, user, course)
	if err != nil {
		log.Printf("Ошибка генерации PDF: %v", err)
		http.Error(w, "Ошибка генерации договора", http.StatusInternalServerError)
		return
	}

	log.Printf("PDF успешно сгенерирован, отправка ответа для courseId: %s", courseId.Hex())
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract_"+courseId.Hex()+".pdf")
	_, err = w.Write(pdfBytes)
	if err != nil {
		log.Printf("Ошибка записи ответа PDF: %v", err)
		return
	}
	log.Printf("Ответ PDF успешно отправлен для courseId: %s", courseId.Hex())
}

// getRegistrationData остается без изменений
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

	return result, user, course, nil
}

type FieldPosition struct {
	X, Y, Width float64
}

var fieldPositions = map[string]FieldPosition{
	"FullName":       {X: 100, Y: 700, Width: 400},
	"CourseTitle":    {X: 100, Y: 600, Width: 400},
	"CourseDuration": {X: 100, Y: 580, Width: 100},
	"CoursePrice":    {X: 100, Y: 560, Width: 100},
	"DocumentId":     {X: 100, Y: 720, Width: 100},
	"Adress":         {X: 100, Y: 200, Width: 400},
	"PasportDate":    {X: 100, Y: 180, Width: 200},
	"SNILS":          {X: 100, Y: 160, Width: 100},
	"Phone":          {X: 100, Y: 140, Width: 100},
	"Email":          {X: 100, Y: 120, Width: 200},
}

func generateFilledContract(registration, user, course bson.M) ([]byte, error) {
	inputPath := "./document_download/ДОГОВОР.pdf"
	log.Printf("Открытие шаблона PDF по пути: %s", inputPath)

	// Проверка существования файла шаблона
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Файл шаблона PDF не существует: %s", inputPath)
	}

	// Инициализация gofpdf
	pdf := gofpdf.New("P", "pt", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("arial", "", "./fonts/arial.ttf")
	pdf.SetFont("arial", "", 12)

	// Подготовка данных для заполнения
	data := make(map[string]string)

	if lastName, ok := user["lastname"].(string); ok {
		firstName, _ := user["firstname"].(string)
		middleName, _ := user["middlename"].(string)
		data["FullName"] = fmt.Sprintf("%s %s %s", lastName, firstName, middleName)
	} else {
		return nil, fmt.Errorf("недействительная или отсутствующая фамилия пользователя")
	}

	if title, ok := course["title"].(string); ok {
		data["CourseTitle"] = title
	} else {
		return nil, fmt.Errorf("недействительное или отсутствующее название курса")
	}

	var durationVal int
	switch v := course["duration"].(type) {
	case int:
		durationVal = v
	case float64:
		durationVal = int(v)
	case int32:
		durationVal = int(v)
	case int64:
		durationVal = int(v)
	default:
		return nil, fmt.Errorf("недействительная или отсутствующая длительность курса: получен тип %T", v)
	}
	data["CourseDuration"] = fmt.Sprintf("%d часов", durationVal)

	var priceVal float64
	switch v := course["price"].(type) {
	case int:
		priceVal = float64(v)
	case float64:
		priceVal = v
	case int32:
		priceVal = float64(v)
	case int64:
		priceVal = float64(v)
	default:
		return nil, fmt.Errorf("недействительная или отсутствующая цена курса: получен тип %T", v)
	}
	data["CoursePrice"] = fmt.Sprintf("%.0f руб.", priceVal)

	data["DocumentId"] = time.Now().Format("02")

	if address, ok := user["homeaddress"].(string); ok {
		data["Adress"] = address
	} else {
		data["Adress"] = ""
	}

	if passport, ok := user["passportdata"].(string); ok {
		data["PasportDate"] = passport
	} else {
		data["PasportDate"] = ""
	}

	if snils, ok := user["snils"].(string); ok {
		data["SNILS"] = snils
	} else {
		data["SNILS"] = ""
	}

	if phone, ok := user["phone"].(string); ok {
		data["Phone"] = phone
	} else {
		data["Phone"] = ""
	}

	if email, ok := user["email"].(string); ok {
		data["Email"] = email
	} else {
		data["Email"] = ""
	}

	// Размещение текста по координатам
	for field, value := range data {
		if pos, exists := fieldPositions[field]; exists {
			pdf.SetXY(pos.X, pos.Y)
			pdf.Write(0, value)
		}
	}

	// Запись в буфер
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("не удалось записать PDF: %v", err)
	}

	log.Printf("PDF успешно сгенерирован, размер: %d байт", buf.Len())
	return buf.Bytes(), nil
}
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

	dir := "../server/document"
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
		http.Error(w, "Ошибка при обновлении статуса договора", http.StatusInternalServerError)
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
