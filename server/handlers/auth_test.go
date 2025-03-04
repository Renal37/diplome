package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginUser(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(LoginUser))
	defer ts.Close()

	// Создаем тестовые данные для входа
	credentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testuser",
		Password: "testpassword",
	}

	// Преобразуем данные в JSON
	jsonData, err := json.Marshal(credentials)
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

	// Проверяем наличие токена в ответе
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	if result["token"] == "" {
		t.Error("Токен не был возвращен")
	}
}
