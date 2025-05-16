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

// writeJSONError отправляет JSON-ошибку с указанным сообщением и кодом состояния
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func AddCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		writeJSONError(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Отладка входных данных
	fmt.Printf("Полученные данные: %+v\n", course)

	// Проверяем, что priceId и typeId являются валидными ObjectID
	if course.PriceId.IsZero() || course.TypeId.IsZero() {
		writeJSONError(w, "Неверный формат ID стоимости или типа курса", http.StatusBadRequest)
		return
	}

	// Проверяем, что даты регистрации валидны
	if course.RegistrationStart.IsZero() || course.RegistrationEnd.IsZero() {
		writeJSONError(w, "Даты начала и окончания регистрации обязательны", http.StatusBadRequest)
		return
	}

	// Проверяем, что дата окончания регистрации не раньше даты начала
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

	collection := db.GetCollection(db.CoursesCollection)

	fullCourse := bson.M{
		"title":             course.Title,
		"description":       course.Description,
		"duration":          course.Duration,
		"priceId":           course.PriceId,
		"price":             course.Price,
		"typeId":            course.TypeId,
		"type":              course.Type,
		"createdAt":         time.Now(),
		"registrationStart": course.RegistrationStart,
		"registrationEnd":   course.RegistrationEnd,
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
			"_id":               result.InsertedID,
			"title":             course.Title,
			"description":       course.Description,
			"duration":          course.Duration,
			"priceId":           course.PriceId,
			"price":             course.Price,
			"typeId":            course.TypeId,
			"type":              course.Type,
			"registrationStart": course.RegistrationStart,
			"registrationEnd":   course.RegistrationEnd,
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

	// Отладка входных данных
	fmt.Printf("Полученные данные: %+v\n", course)

	// Проверяем даты регистрации
	if course.RegistrationStart.IsZero() || course.RegistrationEnd.IsZero() {
		writeJSONError(w, "Даты начала и окончания регистрации обязательны", http.StatusBadRequest)
		return
	}

	// Проверяем, что дата окончания регистрации не раньше даты начала
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

	collection := db.GetCollection(db.CoursesCollection)

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":             course.Title,
			"description":       course.Description,
			"duration":          course.Duration,
			"price":             course.Price,
			"type":              course.Type,
			"priceId":           course.PriceId,
			"typeId":            course.TypeId,
			"registrationStart": course.RegistrationStart,
			"registrationEnd":   course.RegistrationEnd,
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

// ... остальные функции остаются без изменений ...
func GetCourseByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CoursesCollection)

	var course models.Course
	filter := bson.M{"_id": courseID}
	err = collection.FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Курс не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при получении курса", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}
func RegisterForCourse(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CourseID primitive.ObjectID `json:"courseId"`
		UserID   primitive.ObjectID `json:"userId"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Проверяем, что курс существует и даты регистрации актуальны
	courseCollection := db.GetCollection(db.CoursesCollection)
	var course models.Course
	err = courseCollection.FindOne(context.Background(), bson.M{"_id": request.CourseID}).Decode(&course)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Курс не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при получении курса", http.StatusInternalServerError)
		}
		return
	}

	currentTime := time.Now()
	if currentTime.Before(course.RegistrationStart) {
		http.Error(w, "Регистрация на курс еще не началась", http.StatusBadRequest)
		return
	}
	if currentTime.After(course.RegistrationEnd) {
		http.Error(w, "Регистрация на курс уже закончилась", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)

	registration := bson.M{
		"courseId":          request.CourseID,
		"userId":            request.UserID,
		"status":            "Ожидание",
		"registerDate":      time.Now().Format("2006-01-02"),
		"registrationStart": course.RegistrationStart,
		"registrationEnd":   course.RegistrationEnd,
	}

	_, err = collection.InsertOne(context.Background(), registration)
	if err != nil {
		http.Error(w, "Ошибка при записи на курс", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func ApproveCourseRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
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
		http.Error(w, "Ошибка при одобрении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func GetCourseRegistrations(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	// Агрегация для получения данных о курсе и пользователе
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
						"Unknown Course", // Значение по умолчанию, если курс не найден
					},
				},
				"userName": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$user.username", 0}},
						"Unknown User", // Значение по умолчанию, если пользователь не найден
					},
				},
				"status":           1,
				"contractFilePath": 1,
				"groupId":          1,
				"userId":           1,
				"rejectReason":     1,
				"expelOrderId":     1,
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
		http.Error(w, "Ошибка при получении заявок", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var registrations []bson.M
	if err = cursor.All(context.Background(), &registrations); err != nil {
		http.Error(w, "Ошибка при обработке данных заявок", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registrations)
}

func ApproveRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
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
		http.Error(w, "Ошибка при одобрении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func RejectRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Reason string `json:"reason"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.Reason == "" {
		http.Error(w, "Причина отклонения обязательна", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status":           "Отклоненный",
			"rejectReason":     requestBody.Reason, // Сохраняем причину отклонения
			"contractFilePath": bson.TypeNull,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Ошибка при отклонении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func GetCoursesByStatus(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie с именем "token"
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	// Парсим токен
	tokenStr := tokenCookie.Value
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		fmt.Println("Неверный токен:", err)
		return
	}

	// Извлекаем userId из claims
	userId := claims.UserID
	if userId == "" {
		http.Error(w, "Идентификатор пользователя не найден", http.StatusUnauthorized)
		return
	}

	// Преобразуем userId в ObjectID
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}

	// Получаем параметр status из запроса
	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "Не указан статус", http.StatusBadRequest)
		return
	}

	// Подключаемся к базе данных
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	// Агрегация для получения курсов конкретного пользователя с указанным статусом
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
		http.Error(w, "Ошибка при получении курсов", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []bson.M
	if err = cursor.All(context.Background(), &courses); err != nil {
		http.Error(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func GetCoursesForUser(w http.ResponseWriter, r *http.Request) {
	// Получаем cookie с именем "token"
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	// Парсим токен
	tokenStr := tokenCookie.Value
	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	// Извлекаем userId из claims
	userId := claims.UserID
	if userId == "" {
		http.Error(w, "Идентификатор пользователя не найден", http.StatusUnauthorized)
		return
	}

	// Преобразуем userId в ObjectID
	userIdObj, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}

	// Подключаемся к базе данных
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	// Агрегация для получения курсов конкретного пользователя
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
		http.Error(w, "Ошибка при получении курсов", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []bson.M
	if err = cursor.All(context.Background(), &courses); err != nil {
		http.Error(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// Отчисление пользователя
func ExpelRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		OrderID primitive.ObjectID `json:"orderId"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.CourseRegistrationsCollection)

	// Получаем тип приказа
	orderCollection := db.GetCollection(db.OrderCollection)
	var order models.Order
	err = orderCollection.FindOne(context.Background(), bson.M{"_id": requestBody.OrderID}).Decode(&order)
	if err != nil {
		http.Error(w, "Приказ не найден", http.StatusBadRequest)
		return
	}

	var newStatus string
	if order.OrderType == "О выпуске обучающихся" {
		newStatus = "Завершил"

	} else if order.OrderType == "О зачислении обучающихся" {
		newStatus = "Проходит курс"
	} else if order.OrderType == "Об отчислении обучающихся" {
		newStatus = "Отчисленный"
	} else {
		newStatus = "Отчисленный" // дефолтный статус
	}

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status":       newStatus,
			"expelOrderId": requestBody.OrderID,
			"expelDate":    time.Now(),
		},
	}

	if newStatus == "Завершил" {
		update["$set"].(bson.M)["documentType"] = "Диплом"
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Ошибка при отчислении", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// Выдача документа (сертификат/диплом)
func IssueDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		DocumentType string `json:"documentType"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.DocumentType == "" {
		http.Error(w, "Тип документа обязателен", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"documentType": requestBody.DocumentType, // Сохраняем тип документа
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Ошибка при выдаче документа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func DeleteRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	filter := bson.M{"_id": registrationID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Ошибка при удалении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func WithdrawRegistration(w http.ResponseWriter, r *http.Request) {
	// Получаем ID заявки из URL
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	// Получаем ID пользователя из токена
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Токен отсутствует", http.StatusUnauthorized)
		return
	}

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}

	// Подключаемся к базе данных
	collection := db.GetCollection(db.CourseRegistrationsCollection)

	// Проверяем, что заявка принадлежит пользователю
	filter := bson.M{"_id": registrationID, "userId": userID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Ошибка при удалении заявки", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Заявка не найдена или не принадлежит пользователю", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func PayCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registrationID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
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
		http.Error(w, "Ошибка при обновлении статуса оплаты", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
