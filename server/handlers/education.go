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

func AddEducation(w http.ResponseWriter, r *http.Request) {
	var education models.Education
	if err := json.NewDecoder(r.Body).Decode(&education); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.EducationsCollection)
	result, err := collection.InsertOne(context.Background(), education)
	if err != nil {
		http.Error(w, "Error adding education", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"education": map[string]interface{}{
			"id":   result.InsertedID,
			"name": education.Name,
		},
	})
}

func GetEducations(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection(db.EducationsCollection)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching educations", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var educations []models.Education
	if err := cursor.All(context.Background(), &educations); err != nil {
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	if educations == nil {
		educations = []models.Education{} // Ensure empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(educations)
}

func UpdateEducation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var education models.Education
	if err := json.NewDecoder(r.Body).Decode(&education); err != nil {
		http.Error(w, "Invalid data format", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.EducationsCollection)
	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"name": education.Name}},
	)
	if err != nil {
		http.Error(w, "Error updating education", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func DeleteEducation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Проверка использования уровня образования
	usersCollection := db.GetCollection(db.UsersCollection)
	count, err := usersCollection.CountDocuments(
		context.Background(),
		bson.M{"educationId": id},
	)
	if err != nil {
		http.Error(w, "Error checking education usage", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Cannot delete - education is in use", http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.EducationsCollection)
	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Error deleting education", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
