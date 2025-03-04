package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Renal37/models" // Импортируем пакет models
)

func TestAddCourse(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(AddCourse))
	defer ts.Close()

	// Создаем тестовый курс
	course := models.Course{
		Title:       "Test Course",
		Description: "This is a test course",
		Duration:    10,
		Price:       100,
		Type:        "Online",
	}

	// Преобразуем курс в JSON
	jsonData, err := json.Marshal(course)
	if err != nil {
		t.Fatalf("Ошибка при маршалинге данных: %v", err)
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус код 200, получен %d", resp.StatusCode)
	}
}

func TestGetCourses(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(GetCourses))
	defer ts.Close()

	// Создаем запрос
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус код 200, получен %d", resp.StatusCode)
	}

	// Проверяем, что ответ содержит список курсов
	var courses []models.Course
	err = json.NewDecoder(resp.Body).Decode(&courses)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	if len(courses) == 0 {
		t.Error("Ожидался непустой список курсов")
	}
}