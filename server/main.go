package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Course представляет структуру данных курса
type Course struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"` // Продолжительность в часах
	Price       int    `json:"price"`    // Стоимость курса
	Type        string `json:"type"`     // Тип курса
}

// corsMiddleware sets the necessary headers to allow CORS requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Подключение к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	// Создание маршрутов
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/add-course", addCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/courses", getCourses).Methods("GET")
	r.HandleFunc("/update-course/{id}", updateCourse).Methods("PUT", "OPTIONS")

	// Запуск сервера
	fmt.Println("Server is running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}

// addCourse обрабатывает запросы на добавление нового курса
func addCourse(w http.ResponseWriter, r *http.Request) {
	var course Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Подключение к коллекции courses
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		http.Error(w, "Ошибка подключения к базе данных", http.StatusInternalServerError)
		return
	}
	collection := client.Database("diplome").Collection("courses")

	// Вставка данных курса в коллекцию
	_, err = collection.InsertOne(context.Background(), course)
	if err != nil {
		http.Error(w, "Ошибка при добавлении курса в базу данных", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Курс '%s' успешно добавлен!", course.Title)
}

// getCourses обрабатывает запросы на получение списка курсов
func getCourses(w http.ResponseWriter, r *http.Request) {
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

	var courses []Course
	if err = cursor.All(context.Background(), &courses); err != nil {
		http.Error(w, "Ошибка при обработке данных курсов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

// updateCourse обрабатывает запросы на обновление курса по его идентификатору
func updateCourse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var course Course
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
