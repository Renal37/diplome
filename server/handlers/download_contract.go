package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Renal37/db"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DownloadContract - обработчик для скачивания договора
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

// getRegistrationData - получение данных регистрации
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

// FieldPosition - структура для координат полей
type FieldPosition struct {
	X, Y  float64
	Page  int
	Width float64
}

// fieldPositions - карта позиций полей в PDF
var fieldPositions = map[string]FieldPosition{
	"FullName":       {X: 120, Y: 680, Page: 4, Width: 400},
	"CourseTitle":    {X: 150, Y: 500, Page: 1, Width: 300},
	"CourseDuration": {X: 350, Y: 500, Page: 1, Width: 100},
	"CoursePrice":    {X: 450, Y: 500, Page: 1, Width: 100},
	"DocumentId":     {X: 300, Y: 750, Page: 1, Width: 100},
	"Adress":         {X: 120, Y: 220, Page: 4, Width: 400},
	"PasportDate":    {X: 120, Y: 200, Page: 4, Width: 200},
	"SNILS":          {X: 120, Y: 180, Page: 4, Width: 100},
	"Phone":          {X: 120, Y: 160, Page: 4, Width: 100},
	"Email":          {X: 120, Y: 140, Page: 4, Width: 200},
}

// generateFilledContract - генерация заполненного PDF
func generateFilledContract(registration, user, course bson.M) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения текущей директории: %v", err)
	}
	log.Printf("Текущая рабочая директория: %s", cwd)

	inputPath := filepath.Join(cwd, "document_download", "ДОГОВОР.pdf")
	log.Printf("Открытие шаблона PDF по пути: %s", inputPath)

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("файл шаблона PDF не существует: %s", inputPath)
	}

	// Подготовка данных
	data := prepareContractData(user, course)

	// Создание временного файла
	outputPath := filepath.Join(os.TempDir(), fmt.Sprintf("contract_%d.pdf", time.Now().UnixNano()))
	defer os.Remove(outputPath)

	// Копирование исходного PDF
	if err := copyFile(inputPath, outputPath); err != nil {
		return nil, fmt.Errorf("не удалось скопировать PDF: %v", err)
	}

	// Добавление текста в PDF с помощью pdftk
	if err := addTextWithPdftk(data, outputPath); err != nil {
		return nil, fmt.Errorf("ошибка добавления текста в PDF: %v", err)
	}

	// Чтение результата
	return os.ReadFile(outputPath)
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

	// Номер договора (текущий день)
	data["DocumentId"] = time.Now().Format("02")

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

	// Логирование данных
	for field, value := range data {
		log.Printf("Поле %s: %s", field, value)
	}

	return data
}

// copyFile - копирование файла
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// addTextWithPdftk - добавление текста с помощью pdftk
func addTextWithPdftk(data map[string]string, pdfPath string) error {
	// Проверяем наличие pdftk
	if _, err := exec.LookPath("pdftk"); err != nil {
		return fmt.Errorf("pdftk не найден, установите его: sudo apt install pdftk")
	}

	// Создаем временный файл для данных
	dataPath := filepath.Join(os.TempDir(), fmt.Sprintf("data_%d.fdf", time.Now().UnixNano()))
	defer os.Remove(dataPath)

	// Генерируем FDF файл с данными
	if err := generateFdfFile(data, dataPath); err != nil {
		return fmt.Errorf("ошибка генерации FDF файла: %v", err)
	}

	// Выполняем команду pdftk для заполнения формы
	outputPath := pdfPath + ".filled.pdf"
	defer os.Remove(outputPath)

	cmd := exec.Command("pdftk", pdfPath, "fill_form", dataPath, "output", outputPath, "flatten")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка выполнения pdftk: %v, stderr: %s", err, stderr.String())
	}

	// Заменяем оригинальный файл заполненным
	if err := os.Rename(outputPath, pdfPath); err != nil {
		return fmt.Errorf("ошибка замены файла: %v", err)
	}

	return nil
}

// generateFdfFile - генерация FDF файла для pdftk
func generateFdfFile(data map[string]string, outputPath string) error {
	var builder strings.Builder

	builder.WriteString("%FDF-1.2\n")
	builder.WriteString("1 0 obj\n")
	builder.WriteString("<<\n")
	builder.WriteString("/FDF << /Fields [\n")

	for field, value := range data {
		builder.WriteString(fmt.Sprintf("<< /T (%s) /V (%s) >>\n", field, value))
	}

	builder.WriteString("] >>\n")
	builder.WriteString(">>\n")
	builder.WriteString("endobj\n")
	builder.WriteString("trailer\n")
	builder.WriteString("<<\n")
	builder.WriteString("/Root 1 0 R\n")
	builder.WriteString(">>\n")
	builder.WriteString("%%EOF\n")

	return os.WriteFile(outputPath, []byte(builder.String()), 0644)
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