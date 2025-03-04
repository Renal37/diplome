package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/Renal37/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Подключение к MongoDB
	client := connectToMongoDB()
	defer client.Disconnect(context.Background())

	// Создание маршрутов
	r := setupRouter()

	fmt.Println("Server is running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func connectToMongoDB() *mongo.Client {
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
	return client
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(middleware.CorsMiddleware)

	// Регистрация маршрутов
	registerRoutes(r)

	return r
}