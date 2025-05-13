package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/Renal37/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Prompt struct {
	ID      string `json:"id" bson:"_id"`
	Content string `json:"content" bson:"content"`
}

func GetPrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	collection := db.GetCollection(db.Prompts)
	ctx := context.Background()

	var prompt Prompt
	err := collection.FindOne(ctx, bson.M{}).Decode(&prompt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Создаем дефолтный промпт если нет в базе
			defaultPrompt := Prompt{
				Content: `Ты - AI-ассистент информационной системы для записи на курсы дополнительного профессионального образования...`, // ваш текущий промпт
			}
			json.NewEncoder(w).Encode(defaultPrompt)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(prompt)
}

func UpdatePrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var prompt Prompt
	if err := json.NewDecoder(r.Body).Decode(&prompt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := db.GetCollection(db.Prompts)
	ctx := context.Background()

	// Обновляем или создаем промпт
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{"$set": bson.M{"content": prompt.Content}},
		opts,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func DeletePrompt(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	collection := db.GetCollection(db.Prompts)
	ctx := context.Background()

	_, err := collection.DeleteOne(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}