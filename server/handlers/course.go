package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddCourse(w http.ResponseWriter, r *http.Request) {
	var course models.Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

	_, err = collection.InsertOne(context.Background(), course)
	if err != nil {
		http.Error(w, "Ошибка при добавлении курса в базу данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс '%s' успешно добавлен!", course.Title)
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Ошибка при получении курсов из базы данных", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courses []bson.M
	if err = cursor.All(context.Background(), &courses); err != nil {
		http.Error(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	for i := range courses {
		courses[i]["_id"] = courses[i]["_id"].(primitive.ObjectID).Hex()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func UpdateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var course models.Course
	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"title":       course.Title,
			"description": course.Description,
			"duration":    course.Duration,
			"price":       course.Price,
			"type":        course.Type,
		},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Ошибка при обновлении курса в базе данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс '%s' успешно обновлен!", course.Title)
}

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Удаляем курс
	courseCollection := client.Database("diplome").Collection("courses")
	filter := bson.M{"_id": id}
	_, err = courseCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Ошибка при удалении курса из базы данных", http.StatusInternalServerError)
		return
	}

	// Удаляем все заявки, связанные с этим курсом
	registrationCollection := client.Database("diplome").Collection("course_registrations")
	_, err = registrationCollection.DeleteMany(context.Background(), bson.M{"courseId": id})
	if err != nil {
		http.Error(w, "Ошибка при удалении заявок на курс", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс и все связанные заявки успешно удалены!")
}
func GetCourseByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора курса", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("course_registrations")

	registration := bson.M{
		"courseId": request.CourseID,
		"userId":   request.UserID,
		"status":   "Ожидание", // Статус ожидания одобрения
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

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("course_registrations")

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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("course_registrations")

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

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("course_registrations")

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

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("course_registrations")

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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("course_registrations")

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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("diplome").Collection("course_registrations")

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
		Reason string `json:"reason"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	if requestBody.Reason == "" {
		http.Error(w, "Причина отчисления обязательна", http.StatusBadRequest)
		return
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("course_registrations")

	filter := bson.M{"_id": registrationID}
	update := bson.M{
		"$set": bson.M{
			"status":       "Отчисленный",
			"rejectReason": requestBody.Reason, // Сохраняем причину отчисления
		},
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

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("course_registrations")

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

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("course_registrations")

	filter := bson.M{"_id": registrationID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Ошибка при удалении заявки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
