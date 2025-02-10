package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Renal37/models"
	"github.com/Renal37/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
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
	collection := client.Database("diplome").Collection("courses")

	filter := bson.M{"_id": id}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Ошибка при удалении курса из базы данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс успешно удален!")
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
				"status": 1,
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
			"status":       "Отклоненный",
			"rejectReason": requestBody.Reason, // Сохраняем причину отклонения
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
			"$project": bson.M{
				"courseTitle": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$course.title", 0}},
						"Unknown Course",
					},
				},
				"status":       1,
				"rejectReason": 1,
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
			"$project": bson.M{
				"courseTitle": bson.M{
					"$ifNull": bson.A{
						bson.M{"$arrayElemAt": bson.A{"$course.title", 0}},
						"Unknown Course",
					},
				},
				"status":       1,
				"rejectReason": 1,
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

	// Создание PDF-документа
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Добавление шрифта Times New Roman
	fontPathRegular := "../server/font/timesnewromanpsmt.ttf"
	fontPathBold := "../server/font/timesnewromanbold.ttf"

	_, err = os.Stat(fontPathRegular)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Regular font file not found: %v", err)
			http.Error(w, "Шрифт (Regular) не найден", http.StatusInternalServerError)
			return
		}
		log.Printf("Error checking regular font file: %v", err)
		http.Error(w, "Ошибка при проверке шрифта (Regular)", http.StatusInternalServerError)
		return
	}

	_, err = os.Stat(fontPathBold)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Bold font file not found: %v", err)
			http.Error(w, "Шрифт (Bold) не найден", http.StatusInternalServerError)
			return
		}
		log.Printf("Error checking bold font file: %v", err)
		http.Error(w, "Ошибка при проверке шрифта (Bold)", http.StatusInternalServerError)
		return
	}

	pdf.AddUTF8Font("TimesNewRoman", "", fontPathRegular)
	pdf.AddUTF8Font("TimesNewRoman", "B", fontPathBold)

	// Настройка отступов и шрифтов
	topMargin := 20.0
	leftMargin := 20.0
	rightMargin := 20.0
	lineHeight := 7.0

	pdf.SetMargins(leftMargin, topMargin, rightMargin)
	pdf.SetAutoPageBreak(true, 20)

	// Заголовок
	pdf.SetFont("TimesNewRoman", "B", 14)
	pdf.CellFormat(0, 10, "ДОГОВОР № _______", "", 0, "C", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.CellFormat(0, 10, "об образовании на обучение по дополнительным образовательным программам", "", 0, "C", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(0, 10, "г. Альметьевск                                  «___» ____________ 20__ г.", "", 0, "L", false, 0, "")
	pdf.Ln(20)

	// Раздел I. Предмет Договора
	pdf.SetFont("TimesNewRoman", "B", 12)
	pdf.CellFormat(0, 10, "I. Предмет Договора", "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.MultiCell(0, lineHeight, "1.1. Исполнитель обязуется предоставить образовательную услугу по обучению по программе, указанной в п.1.2. настоящего договора, в пределах федерального государственного образовательного стандарта и (или) профессиональных стандартов в соответствии с учебным планом, в том числе индивидуальным, и образовательной программой Исполнителя, а Обучающийся обязуется оплатить указанную образовательную услугу.", "", "L", false)
	pdf.Ln(5)
	pdf.CellFormat(0, 10, fmt.Sprintf("1.2. Наименование дополнительной образовательной программы: %s", course["title"]), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, fmt.Sprintf("1.3. Срок обучения составляет %d часов.", course["duration"]), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, "1.4. Форма обучения – очно-заочная.", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.MultiCell(0, lineHeight, "1.5. После освоения Обучающимся образовательной программы, успешного прохождения итоговой аттестации и полностью оплатившему за обучение, в соответствии с условиями договора выдается документ установленного образца - диплом о профессиональной переподготовке.", "", "L", false)
	pdf.Ln(20)

	// Раздел IX. Адреса и реквизиты сторон
	pdf.SetFont("TimesNewRoman", "B", 12)
	pdf.CellFormat(0, 10, "IX. Адреса и реквизиты сторон", "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("TimesNewRoman", "", 12)
	pdf.CellFormat(0, 10, "Исполнитель:", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, "ГАПОУ «Альметьевский политехнический техникум»", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, "423457, РТ, г. Альметьевск, ул. Мира д.10.", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, "ИНН/КПП 1644005722/164401001", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, "ОГРН 1021601625352", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, "Тел.8(8553)33-48-60", "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.CellFormat(0, 10, "Обучающийся:", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, fmt.Sprintf("ФИО: %s %s %s", user["lastname"], user["firstname"], user["middlename"]), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, fmt.Sprintf("Адрес: %s", user["homeAddress"]), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, fmt.Sprintf("Паспорт: %s", user["passportData"]), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, fmt.Sprintf("Телефон: %s", user["phone"]), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 10, fmt.Sprintf("E-mail: %s", user["email"]), "", 0, "L", false, 0, "")

	// Отправка PDF-файла клиенту
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=contract.pdf")
	if err := pdf.Output(w); err != nil {
		log.Printf("PDF generation error: %v", err)
		http.Error(w, "Ошибка при генерации PDF", http.StatusInternalServerError)
	}
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
