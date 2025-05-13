package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Renal37/db"
	"github.com/Renal37/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddCourseType adds a new course type
func AddCourseType(w http.ResponseWriter, r *http.Request) {
	var courseType models.CourseType
	err := json.NewDecoder(r.Body).Decode(&courseType)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("course_types")

	// Check if type already exists
	var existingType models.CourseType
	err = collection.FindOne(context.Background(), bson.M{"name": courseType.Name}).Decode(&existingType)
	if err == nil {
		http.Error(w, "Тип курса с таким названием уже существует", http.StatusConflict)
		return
	}

	result, err := collection.InsertOne(context.Background(), courseType)
	if err != nil {
		http.Error(w, "Ошибка при добавлении типа курса", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      result.InsertedID,
	})
}

// GetCourseTypes returns all course types
func GetCourseTypes(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("course_types")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Ошибка при получении типов курсов", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var courseTypes []models.CourseType
	if err = cursor.All(context.Background(), &courseTypes); err != nil {
		http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courseTypes)
}

// UpdateCourseType updates a course type
func UpdateCourseType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	var courseType models.CourseType
	err = json.NewDecoder(r.Body).Decode(&courseType)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("course_types")

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"name": courseType.Name}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Ошибка при обновлении типа курса", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// DeleteCourseType deletes a course type
func DeleteCourseType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection("course_types")

	// Check if any courses are using this type
	coursesCollection := db.GetCollection(db.CoursesCollection)
	count, err := coursesCollection.CountDocuments(context.Background(), bson.M{"typeId": id})
	if err != nil {
		http.Error(w, "Ошибка при проверке использования типа", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Невозможно удалить тип курса, так как он используется в одном или нескольких курсах", http.StatusConflict)
		return
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Ошибка при удалении типа курса", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}