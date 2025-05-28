package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


func AddCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Отладка входных данных
	fmt.Printf("Полученные данные: %+v\n", course)

	// Проверяем обязательные поля
	if course.Title == "" || course.Description == "" {
		writeJSONError(w, "Заголовок и описание обязательны", http.StatusBadRequest)
		return
	}
	if course.Duration < 1 {
		writeJSONError(w, "Продолжительность должна быть положительной", http.StatusBadRequest)
		return
	}
	if course.MaxStudents < 1 {
		writeJSONError(w, "Максимальное количество студентов должно быть положительным", http.StatusBadRequest)
		return
	}
	if course.PriceId.IsZero() || course.TypeId.IsZero() {
		writeJSONError(w, "Неверный формат ID стоимости или типа курса", http.StatusBadRequest)
		return
	}
	if course.RegistrationStart.IsZero() || course.RegistrationEnd.IsZero() {
		writeJSONError(w, "Даты начала и окончания регистрации обязательны", http.StatusBadRequest)
		return
	}
	if course.RegistrationEnd.Before(course.RegistrationStart) {
		writeJSONError(w, "Дата окончания регистрации не может быть раньше даты начала", http.StatusBadRequest)
		return
	}

	// Проверяем, что дата начала регистрации не раньше текущей даты
	today := time.Now().Truncate(24 * time.Hour)
	fmt.Printf("course.RegistrationStart: %v, course.RegistrationEnd: %v, today: %v\n",
		course.RegistrationStart, course.RegistrationEnd, today)
	if course.RegistrationStart.Before(today) {
		writeJSONError(w, "Дата начала регистрации не может быть раньше сегодняшней даты", http.StatusBadRequest)
		return
	}

	// Проверка существования priceId
	priceCollection := db.GetCollection(db.PricesCollection)
	var price bson.M
	err = priceCollection.FindOne(context.Background(), bson.M{"_id": course.PriceId}).Decode(&price)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Цена не найдена", http.StatusBadRequest)
		} else {
			writeJSONError(w, "Ошибка при проверке цены", http.StatusInternalServerError)
		}
		return
	}

	// Проверка существования typeId
	typeCollection := db.GetCollection(db.CourseTypesCollection)
	var courseType bson.M
	err = typeCollection.FindOne(context.Background(), bson.M{"_id": course.TypeId}).Decode(&courseType)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Тип курса не найден", http.StatusBadRequest)
		} else {
			writeJSONError(w, "Ошибка при проверке типа курса", http.StatusInternalServerError)
		}
		return
	}

	collection := db.GetCollection(db.CoursesCollection)

	fullCourse := bson.M{
		"title":             course.Title,
		"description":       course.Description,
		"duration":          course.Duration,
		"priceId":           course.PriceId,
		"price":             price["amount"],
		"typeId":            course.TypeId,
		"type":              courseType["name"],
		"createdAt":         time.Now(),
		"registrationStart": course.RegistrationStart,
		"registrationEnd":   course.RegistrationEnd,
		"studentsCount":     0,
		"maxStudents":       course.MaxStudents,
	}

	result, err := collection.InsertOne(context.Background(), fullCourse)
	if err != nil {
		writeJSONError(w, "Ошибка при добавлении курса в базу данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"course": map[string]interface{}{
			"_id":               result.InsertedID.(primitive.ObjectID).Hex(),
			"title":             course.Title,
			"description":       course.Description,
			"duration":          course.Duration,
			"priceId":           course.PriceId.Hex(),
			"price":             price["amount"],
			"typeId":            course.TypeId.Hex(),
			"type":              courseType["name"],
			"registrationStart": course.RegistrationStart,
			"registrationEnd":   course.RegistrationEnd,
			"studentsCount":     0,
			"maxStudents":       course.MaxStudents,
		},
	})
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.CoursesCollection)

	pipeline := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         db.PricesCollection,
				"localField":   "priceId",
				"foreignField": "_id",
				"as":           "priceInfo",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         db.CourseTypesCollection,
				"localField":   "typeId",
				"foreignField": "_id",
				"as":           "typeInfo",
			},
		},
		bson.M{
			"$project": bson.M{
				"title":             1,
				"description":       1,
				"duration":          1,
				"price":             bson.M{"$arrayElemAt": bson.A{"$priceInfo.amount", 0}},
				"priceId":           1,
				"typeId":            1,
				"type":              bson.M{"$arrayElemAt": bson.A{"$typeInfo.name", 0}},
				"createdAt":         1,
				"registrationStart": 1,
				"registrationEnd":   1,
				"studentsCount":     1,
				"maxStudents":       1, // Добавляем maxStudents
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		writeJSONError(w, "Ошибка при получении курсов из базы данных", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []bson.M
	if err = cursor.All(context.Background(), &courses); err != nil {
		writeJSONError(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	for i := range courses {
		courses[i]["_id"] = courses[i]["_id"].(primitive.ObjectID).Hex()
		if courses[i]["priceId"] != nil {
			courses[i]["priceId"] = courses[i]["priceId"].(primitive.ObjectID).Hex()
		}
		if courses[i]["typeId"] != nil {
			courses[i]["typeId"] = courses[i]["typeId"].(primitive.ObjectID).Hex()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var course models.Course
	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if course.Title == "" || course.Description == "" {
		writeJSONError(w, "Заголовок и описание обязательны", http.StatusBadRequest)
		return
	}
	if course.Duration < 1 {
		writeJSONError(w, "Продолжительность должна быть положительной", http.StatusBadRequest)
		return
	}
	if course.MaxStudents < 1 {
		writeJSONError(w, "Максимальное количество студентов должно быть положительным", http.StatusBadRequest)
		return
	}
	if course.PriceId.IsZero() || course.TypeId.IsZero() {
		writeJSONError(w, "Неверный формат ID стоимости или типа курса", http.StatusBadRequest)
		return
	}
	if course.RegistrationStart.IsZero() || course.RegistrationEnd.IsZero() {
		writeJSONError(w, "Даты начала и окончания регистрации обязательны", http.StatusBadRequest)
		return
	}
	if course.RegistrationEnd.Before(course.RegistrationStart) {
		writeJSONError(w, "Дата окончания регистрации не может быть раньше даты начала", http.StatusBadRequest)
		return
	}

	// Проверяем, что дата начала регистрации не раньше текущей даты
	today := time.Now().Truncate(24 * time.Hour)
	if course.RegistrationStart.Before(today) {
		writeJSONError(w, "Дата начала регистрации не может быть раньше сегодняшней даты", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CoursesCollection)

	// Получаем текущий курс для проверки studentsCount
	var currentCourse models.Course
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&currentCourse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Курс не найден", http.StatusNotFound)
		} else {
			writeJSONError(w, "Ошибка при получении курса", http.StatusInternalServerError)
		}
		return
	}

	// Проверяем, что maxStudents не меньше текущего studentsCount
	if course.MaxStudents < currentCourse.StudentsCount {
		writeJSONError(w, "Максимальное количество студентов не может быть меньше текущего количества", http.StatusBadRequest)
		return
	}

	// Проверка существования priceId
	priceCollection := db.GetCollection(db.PricesCollection)
	var price bson.M
	err = priceCollection.FindOne(context.Background(), bson.M{"_id": course.PriceId}).Decode(&price)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Цена не найдена", http.StatusBadRequest)
		} else {
			writeJSONError(w, "Ошибка при проверке цены", http.StatusInternalServerError)
		}
		return
	}

	// Проверка существования typeId
	typeCollection := db.GetCollection(db.CourseTypesCollection)
	var courseType bson.M
	err = typeCollection.FindOne(context.Background(), bson.M{"_id": course.TypeId}).Decode(&courseType)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Тип курса не найден", http.StatusBadRequest)
		} else {
			writeJSONError(w, "Ошибка при проверке типа курса", http.StatusInternalServerError)
		}
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":             course.Title,
			"description":       course.Description,
			"duration":          course.Duration,
			"priceId":           course.PriceId,
			"price":             price["amount"],
			"typeId":            course.TypeId,
			"type":              courseType["name"],
			"registrationStart": course.RegistrationStart,
			"registrationEnd":   course.RegistrationEnd,
			"studentsCount":     currentCourse.StudentsCount,
			"maxStudents":       course.MaxStudents,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении курса в базе данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Курс '%s' успешно обновлен!", course.Title)})
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	courseCollection := db.GetCollection(db.CoursesCollection)
	filter := bson.M{"_id": id}
	_, err = courseCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		writeJSONError(w, "Ошибка при удалении курса из базы данных", http.StatusInternalServerError)
		return
	}

	// Удаляем все заявки, связанные с этим курсом
	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	_, err = registrationCollection.DeleteMany(context.Background(), bson.M{"courseId": id})
	if err != nil {
		writeJSONError(w, "Ошибка при удалении заявок на курс", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Курс и все связанные заявки успешно удалены!"})
}

func RegisterForCourse(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CourseID primitive.ObjectID `json:"courseId"`
		UserID   primitive.ObjectID `json:"userId"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	fmt.Printf("Полученные данные: CourseID=%s, UserID=%s\n", request.CourseID.Hex(), request.UserID.Hex())

	// Проверяем, что пользователь существует и профиль заполнен
	userCollection := db.GetCollection(db.UsersCollection)
	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": request.UserID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			writeJSONError(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
		}
		return
	}

	// Проверка обязательных полей пользователя (без изменений)
	requiredFields := map[string]string{
		"lastname":          "Фамилия не заполнена",
		"firstname":         "Имя не заполнено",
		"middlename":        "Отчество не заполнено",
		"birthdate":         "Дата рождения не указана",
		"birthplace":        "Место рождения не указано",
		"educationid":       "Образование не указано",
		"email":             "Email не указан",
		"homeaddress":       "Домашний адрес не указан",
		"jobtitle":          "Должность не указана",
		"passportdata":      "Паспортные данные не указаны",
		"phone":             "Телефон не указан",
		"snils":             "СНИЛС не указан",
		"workplace":         "Место работы не указано",
		"passportissuedby":  "Кем выдан паспорт не указано",
		"passportissuedate": "Дата выдачи паспорта не указана",
		"agreetoprocessing": "Согласие на обработку данных не получено",
		"contractuploaded":  "Соглашение не загружено",
	}

	userMap := bson.M{
		"lastname":          user.LastName,
		"firstname":         user.FirstName,
		"middlename":        user.MiddleName,
		"birthdate":         user.BirthDate,
		"birthplace":        user.BirthPlace,
		"educationid":       user.EducationID,
		"email":             user.Email,
		"homeaddress":       user.HomeAddress,
		"jobtitle":          user.JobTitle,
		"passportdata":      user.PassportData,
		"phone":             user.Phone,
		"snils":             user.SNILS,
		"workplace":         user.WorkPlace,
		"passportissuedby":  user.PassportIssuedBy,
		"passportissuedate": user.PassportIssueDate,
		"agreetoprocessing": user.AgreeToProcessing,
		"contractuploaded":  user.ContractUploaded,
	}

	for field, message := range requiredFields {
		if field == "agreetoprocessing" || field == "contractuploaded" {
			if !userMap[field].(bool) {
				writeJSONError(w, message, http.StatusBadRequest)
				return
			}
		} else if field == "educationid" {
			if userMap[field].(primitive.ObjectID).IsZero() {
				writeJSONError(w, message, http.StatusBadRequest)
				return
			}
		} else if userMap[field] == "" {
			writeJSONError(w, message, http.StatusBadRequest)
			return
		}
	}

	// Проверяем, что курс существует и даты регистрации актуальны
	courseCollection := db.GetCollection(db.CoursesCollection)
	var course models.Course
	err = courseCollection.FindOne(context.Background(), bson.M{"_id": request.CourseID}).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Курс не найден", http.StatusNotFound)
		} else {
			writeJSONError(w, "Ошибка при получении курса", http.StatusInternalServerError)
		}
		return
	}

	// Проверка maxStudents
	if course.StudentsCount >= course.MaxStudents {
		writeJSONError(w, "Курс достиг максимального количества студентов", http.StatusBadRequest)
		return
	}

	currentTime := time.Now()
	fmt.Printf("currentTime: %v, registrationStart: %v, registrationEnd: %v\n",
		currentTime, course.RegistrationStart, course.RegistrationEnd)
	if currentTime.Before(course.RegistrationStart) {
		writeJSONError(w, "Регистрация на курс ещё не началась", http.StatusBadRequest)
		return
	}
	if currentTime.After(course.RegistrationEnd) {
		writeJSONError(w, "Регистрация на курс уже закончилась", http.StatusBadRequest)
		return
	}

	// Проверка уровня образования для курсов типа "Профессиональная переподготовка" (без изменений)
	if course.Type == "Профессиональная переподготовка" {
		educationCollection := db.GetCollection(db.EducationsCollection)
		var education bson.M
		err = educationCollection.FindOne(context.Background(), bson.M{"_id": user.EducationID}).Decode(&education)
		if err != nil {
			writeJSONError(w, "Ошибка при получении данных об образовании", http.StatusInternalServerError)
			return
		}
		educationName := education["name"].(string)
		allowedEducations := []string{
			"Среднее профессиональное",
			"Высшее",
			"Высшее образование",
		}
		isAllowed := false
		for _, allowed := range allowedEducations {
			if educationName == allowed {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			writeJSONError(w,
				"Для записи на этот курс требуется среднее профессиональное или высшее образование",
				http.StatusBadRequest)
			return
		}
	}

	// Проверка, не зарегистрирован ли пользователь уже на этот курс
	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	var existingRegistration bson.M
	err = registrationCollection.FindOne(context.Background(),
		bson.M{"courseId": request.CourseID, "userId": request.UserID}).Decode(&existingRegistration)
	if err == nil {
		writeJSONError(w, "Вы уже зарегистрированы на этот курс", http.StatusBadRequest)
		return
	}

	// Начинаем транзакцию
	session, err := db.GetMongoClient().StartSession()
	if err != nil {
		writeJSONError(w, "Ошибка при создании сессии", http.StatusInternalServerError)
		return
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		writeJSONError(w, "Ошибка при старте транзакции", http.StatusInternalServerError)
		return
	}

	// Регистрация пользователя на курс с добавлением цены
	registration := bson.M{
		"courseId":          request.CourseID,
		"userId":            request.UserID,
		"status":            "Ожидание",
		"registerDate":      time.Now().Format("2006-01-02"),
		"registrationStart": course.RegistrationStart,
		"registrationEnd":   course.RegistrationEnd,
		"price":             course.Price, // Добавляем цену из курса
	}

	_, err = registrationCollection.InsertOne(context.Background(), registration)
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при записи на курс", http.StatusInternalServerError)
		return
	}

	// Увеличиваем studentsCount
	_, err = courseCollection.UpdateOne(context.Background(),
		bson.M{"_id": request.CourseID},
		bson.M{"$inc": bson.M{"studentsCount": 1}},
	)
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при обновлении количества студентов", http.StatusInternalServerError)
		return
	}

	err = session.CommitTransaction(context.Background())
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при фиксации транзакции", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func DeleteRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	// Получаем информацию о регистрации
	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	var registration bson.M
	err = registrationCollection.FindOne(context.Background(), bson.M{"_id": registrationID}).Decode(&registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Заявка не найдена", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Ошибка при получении заявки", http.StatusInternalServerError)
		return
	}

	courseID := registration["courseId"].(primitive.ObjectID)

	// Начинаем транзакцию
	session, err := db.GetMongoClient().StartSession()
	if err != nil {
		writeJSONError(w, "Ошибка при создании сессии", http.StatusInternalServerError)
		return
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		writeJSONError(w, "Ошибка при старте транзакции", http.StatusInternalServerError)
		return
	}

	// Удаляем заявку
	_, err = registrationCollection.DeleteOne(context.Background(), bson.M{"_id": registrationID})
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при удалении заявки", http.StatusInternalServerError)
		return
	}

	// Уменьшаем studentsCount
	courseCollection := db.GetCollection(db.CoursesCollection)
	_, err = courseCollection.UpdateOne(context.Background(),
		bson.M{"_id": courseID},
		bson.M{"$inc": bson.M{"studentsCount": -1}},
	)
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при обновлении количества студентов", http.StatusInternalServerError)
		return
	}

	err = session.CommitTransaction(context.Background())
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при фиксации транзакции", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func WithdrawRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из токена
	cookie, err := r.Cookie("token")
	if err != nil {
		writeJSONError(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		writeJSONError(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}

	// Получаем информацию о регистрации
	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	var registration bson.M
	err = registrationCollection.FindOne(context.Background(), bson.M{"_id": registrationID, "userId": userID}).Decode(&registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Заявка не найдена или не принадлежит пользователю", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Ошибка при получении заявки", http.StatusInternalServerError)
		return
	}

	courseID := registration["courseId"].(primitive.ObjectID)

	// Начинаем транзакцию
	session, err := db.GetMongoClient().StartSession()
	if err != nil {
		writeJSONError(w, "Ошибка при создании сессии", http.StatusInternalServerError)
		return
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		writeJSONError(w, "Ошибка при старте транзакции", http.StatusInternalServerError)
		return
	}

	// Удаляем заявку
	_, err = registrationCollection.DeleteOne(context.Background(), bson.M{"_id": registrationID, "userId": userID})
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при удалении заявки", http.StatusInternalServerError)
		return
	}

	// Уменьшаем studentsCount
	courseCollection := db.GetCollection(db.CoursesCollection)
	_, err = courseCollection.UpdateOne(context.Background(),
		bson.M{"_id": courseID},
		bson.M{"$inc": bson.M{"studentsCount": -1}},
	)
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при обновлении количества студентов", http.StatusInternalServerError)
		return
	}

	err = session.CommitTransaction(context.Background())
	if err != nil {
		session.AbortTransaction(context.Background())
		writeJSONError(w, "Ошибка при фиксации транзакции", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func GetCourseRegistrations(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	pipeline := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "courseId",
				"foreignField": "_id",
				"as":           "course",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "userId",
				"foreignField": "_id",
				"as":           "user",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "groups",
				"localField":   "groupId",
				"foreignField": "_id",
				"as":           "group",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "orders",
				"localField":   "expelOrderId",
				"foreignField": "_id",
				"as":           "order",
			},
		},
		bson.M{
			"$project": bson.M{
				"courseTitle": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$course.title", 0}},
						"Unknown Course",
					},
				},
				"userName": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.username", 0}},
						"Unknown User",
					},
				},
				"status":           1,
				"contractFilePath": 1,
				"groupId":          1,
				"userId":           1,
				"rejectReason":     1,
				"expelOrderId":     1,
				"price":            1, // Убедимся, что price включен
				"orderType": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$order.orderType", 0}},
						"Unknown orderType",
					},
				},
				"groupName": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$group.groupName", 0}},
						"Unknown group",
					},
				},
				"userEmail": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.email", 0}},
						"Unknown Email",
					},
				},
				"userFullname": bson.M{
					"$concat": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.lastname", 0}},
						" ",
						bson.M{"$arrayElemAt": bson.A{"$user.firstname", 0}},
						" ",
						bson.M{"$arrayElemAt": bson.A{"$user.middlename", 0}},
					},
				},
				"userBirthdate": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.birthdate", 0}},
						"Unknown Birthdate",
					},
				},
				"userBirthplace": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.birthplace", 0}},
						"Unknown Birthplace",
					},
				},
				"userEducation": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.education", 0}},
						"Unknown Education",
					},
				},
				"userWorkplace": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.workplace", 0}},
						"Unknown Workplace",
					},
				},
				"userJobtitle": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.jobtitle", 0}},
						"Unknown Jobtitle",
					},
				},
				"userHomeaddress": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.homeaddress", 0}},
						"Unknown Homeaddress",
					},
				},
				"userPhone": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.phone", 0}},
						"Unknown Phone",
					},
				},
				"userPassportdata": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.passportdata", 0}},
						"Unknown Passportdata",
					},
				},
				"userSnils": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.snils", 0}},
						"Unknown Snils",
					},
				},
				"usercontractFilePath": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.contractFilePath", 0}},
						"Unknown contractFilePath",
					},
				},
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		writeJSONError(w, "Ошибка при получении заявок", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var registrations []bson.M
	if err = cursor.All(context.Background(), &registrations); err != nil {
		writeJSONError(w, "Ошибка при обработке данных заявок", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registrations)
}

func GetCourseByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CoursesCollection)

	var course models.Course
	filter := bson.M{"_id": courseID}
	err = collection.FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSONError(w, "Курс не найден", http.StatusNotFound)
		} else {
			writeJSONError(w, "Ошибка при получении курса", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func ApproveCourseRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status": "Пройдено",
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при одобрении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func ApproveRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status": "Одобренный",
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при одобрении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func RejectRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Reason string `json:"reason"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.Reason == "" {
		writeJSONError(w, "Причина отклонения обязательна", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status":           "Отклоненный",
			"rejectReason":     requestBody.Reason,
			"contractFilePath": bson.TypeNull,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при отклонении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func GetCoursesByStatus(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		writeJSONError(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	tokenStr := tokenCookie.Value
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		writeJSONError(w, "Неверный токен", http.StatusUnauthorized)
		fmt.Println("Неверный токен:", err)
		return
	}

	userId := claims.UserID
	if userId == "" {
		writeJSONError(w, "Идентификатор пользователя не найден", http.StatusUnauthorized)
		return
	}

	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		writeJSONError(w, "Не указан статус", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{
				"userId": userIdObj,
				"status": status,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "courseId",
				"foreignField": "_id",
				"as":           "course",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "groups",
				"localField":   "groupId",
				"foreignField": "_id",
				"as":           "group",
			},
		},
		bson.M{
			"$project": bson.M{
				"courseTitle": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$course.title", 0}},
						"Unknown Course",
					},
				},
				"groupId": 1,
				"groupName": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$group.groupName", 0}},
						"Unknown group",
					},
				},
				"status":           1,
				"rejectReason":     1,
				"contractFilePath": 1,
				"contractUploaded": 1,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		writeJSONError(w, "Ошибка при получении курсов", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []bson.M
	if err = cursor.All(context.Background(), &courses); err != nil {
		writeJSONError(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func GetCoursesForUser(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		writeJSONError(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	tokenStr := tokenCookie.Value
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		writeJSONError(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	userId := claims.UserID
	if userId == "" {
		writeJSONError(w, "Идентификатор пользователя не найден", http.StatusUnauthorized)
		return
	}

	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)

	pipeline := bson.A{
		bson.M{
			"$match": bson.M{"userId": userIdObj},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "courses",
				"localField":   "courseId",
				"foreignField": "_id",
				"as":           "course",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "groups",
				"localField":   "groupId",
				"foreignField": "_id",
				"as":           "group",
			},
		},
		bson.M{
			"$project": bson.M{
				"courseTitle": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$course.title", 0}},
						"Unknown Course",
					},
				},
				"groupId": 1,
				"groupName": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$group.groupName", 0}},
						"Unknown group",
					},
				},
				"status":           1,
				"rejectReason":     1,
				"contractFilePath": 1,
				"contractUploaded": 1,
			},
		},
	}

	cursor, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		writeJSONError(w, "Ошибка при получении курсов", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []bson.M
	if err = cursor.All(context.Background(), &courses); err != nil {
		writeJSONError(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func EnrollRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		OrderID primitive.ObjectID `json:"orderId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверяем, что приказ существует и имеет тип "О зачислении обучающихся"
	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": requestBody.OrderID}).Decode(&order)
	if err != nil {
		writeJSONError(w, "Приказ не найден", http.StatusBadRequest)
		return
	}
	if order.OrderType != "О зачислении обучающихся" {
		writeJSONError(w, "Неверный тип приказа для зачисления", http.StatusBadRequest)
		return
	}

	// Проверяем текущий статус заявки
	registrationCollection := db.GetCollection(db.CourseRegistrationsCollection)
	var registration bson.M
	err = registrationCollection.FindOne(context.Background(), bson.M{"_id": registrationID}).Decode(&registration)
	if err != nil {
		writeJSONError(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	if registration["status"] != "Оплаченный" {
		writeJSONError(w, "Зачисление возможно только для статуса 'Оплаченный'", http.StatusBadRequest)
		return
	}

	// Обновляем статус заявки
	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status":        "Проходит курс",
			"enrollOrderId": requestBody.OrderID,
			"enrollDate":    time.Now(),
		},
	}

	_, err = registrationCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при зачислении", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func IssueDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		DocumentType string `json:"documentType"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.DocumentType == "" {
		writeJSONError(w, "Тип документа обязателен", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"documentType": requestBody.DocumentType,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при выдаче документа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func PayCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		writeJSONError(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status": "Оплаченный",
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		writeJSONError(w, "Ошибка при обновлении статуса оплаты", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
