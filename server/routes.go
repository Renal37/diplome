package main

import (
	"github.com/Renal37/handlers"
	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router) {
	// Пользовательские маршруты
	r.HandleFunc("/users", handlers.GetUser).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/profile", handlers.GetProfile).Methods("GET")
	r.HandleFunc("/logout", handlers.LogoutUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/update-profile", handlers.UpdateProfile).Methods("PUT", "OPTIONS")

	// Административные маршруты
	r.HandleFunc("/admin/course-registrations", handlers.GetCourseRegistrations).Methods("GET")
	r.HandleFunc("/admin/approve-registration/{id}", handlers.ApproveRegistration).Methods("POST")
	r.HandleFunc("/admin/reject-registration/{id}", handlers.RejectRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/expel-registration/{id}", handlers.ExpelRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/issue-document/{id}", handlers.IssueDocument).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/approve-contract/{id}", handlers.ApproveContract).Methods("POST")

	// Маршруты для курсов
	r.HandleFunc("/update-course/{id}", handlers.UpdateCourse).Methods("PUT", "OPTIONS")
	r.HandleFunc("/delete-course/{id}", handlers.DeleteCourse).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/add-course", handlers.AddCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/courses", handlers.GetCourses).Methods("GET")
	r.HandleFunc("/courses/{id}", handlers.GetCourseByID).Methods("GET")
	r.HandleFunc("/courses/register", handlers.RegisterForCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/delete-registration/{id}", handlers.DeleteRegistration).Methods("POST")
	r.HandleFunc("/user/view-contract/{id}", handlers.ViewContract).Methods("GET")

	// Маршруты для пользовательских курсов
	r.HandleFunc("/user/courses", handlers.GetCoursesForUser).Methods("GET")
	r.HandleFunc("/user/courses/status", handlers.GetCoursesByStatus).Methods("GET")

	// Маршруты для документов
	r.HandleFunc("/user/download-contract/{courseId}", handlers.DownloadContract).Methods("GET")
	r.HandleFunc("/user/upload-contract/{courseId}", handlers.UploadContract).Methods("POST")
	r.HandleFunc("/user/download-document", handlers.DownloadDocument).Methods("GET")
	r.HandleFunc("/user/upload-document/{userId}", handlers.UploadDocument).Methods("POST")
	r.HandleFunc("/user/view-consent/{userId}", handlers.ViewConsent).Methods("GET")

	// Маршруты для групп
	r.HandleFunc("/groups", handlers.GetGroups).Methods("GET")
	r.HandleFunc("/admin/create-group", handlers.CreateGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/update-group/{id}", handlers.UpdateGroup).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/delete-group/{id}", handlers.DeleteGroup).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/admin/assign-group/{id}", handlers.AssignGroup).Methods("POST", "OPTIONS")

	// Маршруты для проверки токена
	r.HandleFunc("/check-token", handlers.CheckToken).Methods("POST", "OPTIONS")
}