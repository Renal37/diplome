package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type AIAssistantRequest struct {
	Question string `json:"question"`
	Context  string `json:"context"`
}

type AIAssistantResponse struct {
	Answer string `json:"answer"`
}

var (
	yandexGPTClient *YandexGPTClient
	yandexGPTOnce   sync.Once
)

func initYandexGPT() error {
	var initErr error
	yandexGPTOnce.Do(func() {
		yandexGPTClient = NewYandexGPTClient() 
	})
	return initErr
}

func AIAssistantHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AIAssistantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := initYandexGPT(); err != nil {
		log.Printf("Yandex GPT init error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	prompt := fmt.Sprintf("Контекст: %s\nВопрос: %s\nОтветь кратко и по делу.", req.Context, req.Question)

	answer, err := yandexGPTClient.GenerateResponse(r.Context(), prompt)
	if err != nil {
		log.Printf("Yandex GPT error: %v", err)
		answer = "Извините, не могу ответить сейчас. Пожалуйста, попробуйте позже."
	}

	response := AIAssistantResponse{
		Answer: answer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
