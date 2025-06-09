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
	r.HandleFunc("/admin/enroll-registration/{id}", handlers.EnrollRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/issue-document/{id}", handlers.IssueDocument).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/approve-contract/{id}", handlers.ApproveContract).Methods("POST")
	r.HandleFunc("/admin/expel-registration/{id}", handlers.ExpelRegistration).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/expel-registrations", handlers.ExpelMultipleRegistrations).Methods("POST", "OPTIONS")

	// Добавление уровня образования
	r.HandleFunc("/admin/educations", handlers.GetEducations).Methods("GET")
	r.HandleFunc("/admin/educations/add", handlers.AddEducation).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/educations/update/{id}", handlers.UpdateEducation).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/educations/delete/{id}", handlers.DeleteEducation).Methods("DELETE", "OPTIONS")

	// Маршруты для курсов
	r.HandleFunc("/update-course/{id}", handlers.UpdateCourse).Methods("PUT", "OPTIONS")
	r.HandleFunc("/delete-course/{id}", handlers.DeleteCourse).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/add-course", handlers.AddCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/courses", handlers.GetCourses).Methods("GET")
	r.HandleFunc("/courses/{id}", handlers.GetCourseById).Methods("GET")
	r.HandleFunc("/courses/register", handlers.RegisterForCourse).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/delete-registration/{id}", handlers.DeleteRegistration).Methods("POST")
	r.HandleFunc("/user/view-contract/{id}", handlers.ViewContract).Methods("GET")
	r.HandleFunc("/user/withdraw-registration/{id}", handlers.WithdrawRegistration).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/user/pay-course/{id}", handlers.PayCourse).Methods("POST", "OPTIONS")

	// Маршруты для пользовательских курсов
	r.HandleFunc("/user/courses", handlers.GetCoursesForUser).Methods("GET")
	r.HandleFunc("/user/courses/status", handlers.GetCoursesByStatus).Methods("GET")

	// Маршруты для документов
	r.HandleFunc("/user/download-contract/{courseId}", handlers.DownloadContract).Methods("GET")
	r.HandleFunc("/user/upload-contract/{courseId}", handlers.UploadContract).Methods("POST")
	r.HandleFunc("/user/download-document", handlers.DownloadDocument).Methods("GET")
	r.HandleFunc("/user/download-document/{courseId}", handlers.DownloadDocs).Methods("GET")
	r.HandleFunc("/user/document-template/{courseId}", handlers.GetDocsTemplate).Methods("GET")
	r.HandleFunc("/user/upload-document/{userId}", handlers.UploadDocument).Methods("POST")
	r.HandleFunc("/user/view-consent/{userId}", handlers.ViewConsent).Methods("GET")
	r.HandleFunc("/user/fill-consent", handlers.FillConsent).Methods("GET")

	// Маршруты для групп
	r.HandleFunc("/groups", handlers.GetGroups).Methods("GET")
	r.HandleFunc("/admin/create-group", handlers.CreateGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/update-group/{id}", handlers.UpdateGroup).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/delete-group/{id}", handlers.DeleteGroup).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/admin/assign-group/{id}", handlers.AssignGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/group-members/{id}", handlers.GetGroupMembers).Methods("GET")
	r.HandleFunc("/admin/enroll-group/{groupId}", handlers.EnrollGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/expel-group/{groupId}", handlers.ExpelGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/reject-group-registrations/{id}", handlers.RejectGroupRegistrations).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/complete-group/{groupId}", handlers.CompleteGroup).Methods("POST", "OPTIONS")

	// Маршруты для стоимости
	r.HandleFunc("/prices", handlers.GetPrices).Methods("GET", "OPTIONS")
	r.HandleFunc("/add-price", handlers.AddPrice).Methods("POST", "OPTIONS")
	r.HandleFunc("/update-price/{id}", handlers.UpdatePrice).Methods("PUT", "OPTIONS")
	r.HandleFunc("/delete-price/{id}", handlers.DeletePrice).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/bulk-update-prices", handlers.BulkUpdatePrices).Methods("POST", "OPTIONS")

	// Маршруты для типов курсов
	r.HandleFunc("/admin/course-types", handlers.GetCourseTypes).Methods("GET")
	r.HandleFunc("/admin/course-types/add", handlers.AddCourseType).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/course-types/update/{id}", handlers.UpdateCourseType).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/course-types/delete/{id}", handlers.DeleteCourseType).Methods("DELETE", "OPTIONS")

	// Маршруты для типов приказов
	r.HandleFunc("/admin/order-types", handlers.GetOrderTypes).Methods("GET")
	r.HandleFunc("/admin/order-types/add", handlers.AddOrderType).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/order-types/update", handlers.UpdateOrderType).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/order-types/delete", handlers.DeleteOrderType).Methods("DELETE", "OPTIONS")

	// Маршруты для приказов
	r.HandleFunc("/admin/orders", handlers.GetOrders).Methods("GET")
	r.HandleFunc("/admin/orders/add", handlers.AddOrder).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/orders/update", handlers.UpdateOrder).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/orders/delete", handlers.DeleteOrder).Methods("DELETE", "OPTIONS")

	// Маршруты для проверки токена
	r.HandleFunc("/check-token", handlers.CheckToken).Methods("POST", "OPTIONS")

	// Маршруты для ИИ-ассистента
	r.HandleFunc("/ai-assistant", handlers.AIAssistantHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/admin/prompt", handlers.GetPrompt).Methods("GET")
	r.HandleFunc("/admin/prompt/update", handlers.UpdatePrompt).Methods("PUT", "OPTIONS")
	r.HandleFunc("/admin/prompt/delete", handlers.DeletePrompt).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/user/contract-data/{courseId}", handlers.GetContractData).Methods("GET")
	r.HandleFunc("/contract-template", handlers.GetContractTemplate).Methods("GET")
}
