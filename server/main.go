package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Renal37/handlers"
	"github.com/Renal37/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	r.Use(middleware.CorsMiddleware)
	r.HandleFunc("/add-course", handlers.AddCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/courses", handlers.GetCourses).Methods("GET")
	r.HandleFunc("/users", handlers.GetUser).Methods("GET")
	r.HandleFunc("/update-course/{id}", handlers.UpdateCourse).Methods("PUT", "OPTIONS")
	r.HandleFunc("/delete-course/{id}", handlers.DeleteCourse).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/check-token", handlers.CheckToken).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/profile", handlers.GetProfile).Methods("GET")
	r.HandleFunc("/logout", handlers.LogoutUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/update-profile", handlers.UpdateProfile).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/course-registrations", handlers.GetCourseRegistrations).Methods("GET")
	r.HandleFunc("/admin/approve-registration/{id}", handlers.ApproveRegistration).Methods("POST")
	r.HandleFunc("/admin/reject-registration/{id}", handlers.RejectRegistration).Methods("POST")
	r.HandleFunc("/courses/{id}", handlers.GetCourseByID).Methods("GET")
	r.HandleFunc("/courses/register", handlers.RegisterForCourse).Methods("POST", "OPTIONS")
	fmt.Println("Server is running on port 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
