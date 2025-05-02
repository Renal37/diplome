package handlers

import (
	"bytes"
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
	"github.com/pdfcpu/pdfcpu/pkg/api"
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

	log.Printf("Fetching registration data for courseId: %s", courseId.Hex())
	registration, user, course, err := getRegistrationData(courseId)
	if err != nil {
		log.Printf("Error getting registration data: %v", err)
		http.Error(w, "Error getting course data", http.StatusInternalServerError)
		return
	}

	log.Printf("Generating PDF for courseId: %s", courseId.Hex())
	pdfBytes, err := generateFilledContract(registration, user, course)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		http.Error(w, "Error generating contract", http.StatusInternalServerError)
		return
	}

	log.Printf("PDF generated successfully, sending response for courseId: %s", courseId.Hex())
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract_"+courseId.Hex()+".pdf")
	_, err = w.Write(pdfBytes)
	if err != nil {
		log.Printf("Error writing PDF response: %v", err)
		return
	}
	log.Printf("PDF response sent successfully for courseId: %s", courseId.Hex())
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
	inputPath := "../server/document_donwload/ДОГОВОРFullName.pdf"
	log.Printf("Opening PDF template at: %s", inputPath)

	// Проверяем, существует ли файл шаблона
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("PDF template file does not exist: %s", inputPath)
	}

	// Открываем входной файл
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("could not open PDF template: %v", err)
	}
	defer inputFile.Close()

	// Формируем данные для заполнения формы с проверкой типов
	formData := map[string]string{}

	// Основные поля из вывода pdfcpu form list
	if lastName, ok := user["lastname"].(string); ok {
		firstName, _ := user["firstname"].(string)
		middleName, _ := user["middlename"].(string)
		formData["FullName"] = fmt.Sprintf("%s %s %s", lastName, firstName, middleName)
	} else {
		return nil, fmt.Errorf("invalid or missing user lastName")
	}

	// if title, ok := course["title"].(string); ok {
	// 	formData["CourseTitle"] = title
	// } else {
	// 	return nil, fmt.Errorf("invalid or missing course title")
	// }

	// var durationVal int
	// switch v := course["duration"].(type) {
	// case int:
	// 	durationVal = v
	// case float64:
	// 	durationVal = int(v)
	// case int32:
	// 	durationVal = int(v)
	// case int64:
	// 	durationVal = int(v)
	// default:
	// 	return nil, fmt.Errorf("invalid or missing course duration: got type %T", v)
	// }
	// formData["CoruseDuration"] = fmt.Sprintf("%d часов", durationVal)

	// var priceVal float64
	// switch v := course["price"].(type) {
	// case int:
	// 	priceVal = float64(v)
	// case float64:
	// 	priceVal = v
	// case int32:
	// 	priceVal = float64(v)
	// case int64:
	// 	priceVal = float64(v)
	// default:
	// 	return nil, fmt.Errorf("invalid or missing course price: got type %T", v)
	// }
	// formData["CoursePrice"] = fmt.Sprintf("%.0f руб.", priceVal)

	// formData["document_day"] = time.Now().Format("02")
	// formData["DocumentId"] = time.Now().Format("02")

	// if address, ok := user["homeaddress"].(string); ok {
	// 	formData["Adress"] = address
	// } else {
	// 	formData["Adress"] = ""
	// }

	// if passport, ok := user["passportdata"].(string); ok {
	// 	formData["PasportDate"] = passport
	// } else {
	// 	formData["PasportDate"] = ""
	// }

	// if snils, ok := user["snils"].(string); ok {
	// 	formData["SNILS"] = snils
	// } else {
	// 	formData["SNILS"] = ""
	// }

	// if phone, ok := user["phone"].(string); ok {
	// 	formData["Phone"] = phone
	// } else {
	// 	formData["Phone"] = ""
	// }

	// if email, ok := user["email"].(string); ok {
	// 	formData["Email"] = email
	// } else {
	// 	formData["Email"] = ""
	// }

	// formData["HowGive"] = ""
	// formData["CourseEnd"] = ""
	// formData["Time"] = ""
	// formData["TimeDayStart"] = ""
	// formData["TimeMonthStart"] = ""
	// formData["TimeDayEnd"] = ""
	// formData["TimeMonthEnd"] = ""

	// Создаем JSON-данные для формы
	form := map[string]interface{}{
		"Fields": formData,
	}
	formGroup := map[string]interface{}{
		"Forms": []interface{}{form},
	}

	jsonData, err := json.Marshal(formGroup)
	if err != nil {
		return nil, fmt.Errorf("could not marshal form data to JSON: %v", err)
	}

	log.Printf("Filling PDF form with JSON data: %s", string(jsonData))

	// Создаем io.Reader для JSON-данных
	jsonReader := bytes.NewReader(jsonData)

	// Подготавливаем выходной буфер
	var buf bytes.Buffer

	// Заполняем форму
	err = api.FillForm(inputFile, jsonReader, &buf, nil)
	if err != nil {
		return nil, fmt.Errorf("could not fill PDF form: %v", err)
	}

	log.Printf("PDF form filled successfully, size: %d bytes", buf.Len())
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

	collection := db.GetCollection(db.CourseRegistrationsCollection)
	filter := bson.M{"_id": courseId}
	var courseRegistration bson.M
	err = collection.FindOne(context.Background(), filter).Decode(&courseRegistration)
	if err != nil {
		log.Printf("Course not found: %v", err)
		http.Error(w, "Курс не найден", http.StatusNotFound)
		return
	}

	file, handler, err := r.FormFile("contract")
	if err != nil {
		log.Printf("Error retrieving file: %v", err)
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

	collection := db.GetCollection(db.CourseRegistrationsCollection)
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

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("File does not exist: %v", err)
		http.Error(w, "Файл договора не найден на сервере", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=contract.pdf")
	http.ServeFile(w, r, filePath)
}
