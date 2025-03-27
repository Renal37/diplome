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

// AddEducation добавляет новый уровень образования
func AddEducation(w http.ResponseWriter, r *http.Request) {
    var education models.Education
    if err := json.NewDecoder(r.Body).Decode(&education); err != nil {
        http.Error(w, "Неверный формат данных", http.StatusBadRequest)
        return
    }

    collection := db.GetCollection(db.EducationsCollection)
    result, err := collection.InsertOne(context.Background(), education)
    if err != nil {
        http.Error(w, "Ошибка при добавлении уровня образования", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "id":      result.InsertedID,
    })
}

// GetEducations возвращает список всех уровней образования
func GetEducations(w http.ResponseWriter, r *http.Request) {
    collection := db.GetCollection(db.EducationsCollection)
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        http.Error(w, "Ошибка при получении уровней образования", http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.Background())

    var educations []models.Education
    if err := cursor.All(context.Background(), &educations); err != nil {
        http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(educations)
}

// UpdateEducation обновляет уровень образования
func UpdateEducation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(vars["id"])
    if err != nil {
        http.Error(w, "Неверный формат ID", http.StatusBadRequest)
        return
    }

    var education models.Education
    if err := json.NewDecoder(r.Body).Decode(&education); err != nil {
        http.Error(w, "Неверный формат данных", http.StatusBadRequest)
        return
    }

    collection := db.GetCollection(db.EducationsCollection)
    _, err = collection.UpdateOne(
        context.Background(),
        bson.M{"_id": id},
        bson.M{"$set": bson.M{"name": education.Name}},
    )
    if err != nil {
        http.Error(w, "Ошибка при обновлении уровня образования", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// DeleteEducation удаляет уровень образования
func DeleteEducation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := primitive.ObjectIDFromHex(vars["id"])
    if err != nil {
        http.Error(w, "Неверный формат ID", http.StatusBadRequest)
        return
    }

    // Проверяем, используется ли уровень образования
    usersCollection := db.GetCollection(db.UsersCollection)
    count, err := usersCollection.CountDocuments(
        context.Background(),
        bson.M{"educationId": id},
    )
    if err != nil {
        http.Error(w, "Ошибка при проверке использования уровня образования", http.StatusInternalServerError)
        return
    }
    if count > 0 {
        http.Error(w, "Невозможно удалить уровень образования, так как он используется", http.StatusBadRequest)
        return
    }

    collection := db.GetCollection(db.EducationsCollection)
    _, err = collection.DeleteOne(context.Background(), bson.M{"_id": id})
    if err != nil {
        http.Error(w, "Ошибка при удалении уровня образования", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}