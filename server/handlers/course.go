package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Renal37/models"
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
			"status": "Отклоненный",
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
	// Получаем userId из контекста или токена
	userIdHex := r.Context().Value("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdHex)
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
				"userId": userId, // Фильтр по userId
				"status": status, // Фильтр по статусу
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
				"status": 1,
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
	// Получаем userId из контекста или токена
	userIdHex := r.Context().Value("userId").(string)
	userId, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		http.Error(w, "Неверный формат идентификатора пользователя", http.StatusBadRequest)
		return
	}

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
			"$match": bson.M{"userId": userId}, // Фильтр по userId
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
				"status": 1,
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
