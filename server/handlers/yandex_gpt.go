package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	yandexGPTAPIURL = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
)

type YandexGPTClient struct {
	iamToken  string
	folderID  string
	modelName string
}

func NewYandexGPTClient() *YandexGPTClient {
	iamToken := os.Getenv("YANDEX_IAM_TOKEN")
	folderID := os.Getenv("YANDEX_FOLDER_ID")

	if iamToken == "" || folderID == "" {
		log.Fatal("Yandex Cloud credentials not found in .env file")
	}

	return &YandexGPTClient{
		iamToken:  iamToken,
		folderID:  folderID,
		modelName: "yandexgpt-lite",
	}
}

type CompletionRequest struct {
	ModelURI          string            `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []Message         `json:"messages"`
}

type CompletionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type CompletionResponse struct {
	Result struct {
		Alternatives []struct {
			Message struct {
				Role string `json:"role"`
				Text string `json:"text"`
			} `json:"message"`
			Status string `json:"status"`
		} `json:"alternatives"`
		Usage struct {
			InputTextTokens  interface{} `json:"inputTextTokens"`  // Принимаем как interface{}
			CompletionTokens interface{} `json:"completionTokens"` // Принимаем как interface{}
			TotalTokens      interface{} `json:"totalTokens"`      // Принимаем как interface{}
		} `json:"usage"`
		ModelVersion string `json:"modelVersion"`
	} `json:"result"`
}

func (y *YandexGPTClient) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	modelURI := fmt.Sprintf("gpt://%s/yandexgpt-lite", y.folderID)

	reqBody := CompletionRequest{
		ModelURI: modelURI,
		CompletionOptions: CompletionOptions{
			Stream:      false,
			Temperature: 0.5,
			MaxTokens:   1000,
		},
		Messages: []Message{
			{
				Role: "system",
				Text: `Ты - AI-ассистент информационной системы для записи на курсы дополнительного профессионального образования. Система разработана на стеке Golang (бэкенд), React (фронтенд) и MongoDB (база данных). 

		Основные функции системы:
1. Для пользователей:
- Регистрация/авторизация с валидацией данных
- Просмотр каталога курсов с фильтрацией
- Запись на курсы и отслеживание статуса заявки
- Личный кабинет с историей обучения
- Загрузка/скачивание договоров

2. Для администраторов:
- Управление курсами (добавление/редактирование)
- Модерация заявок (одобрение/отклонение)
- Формирование приказов и сертификатов
- Просмотр аналитики

Технические особенности:
- Безопасность: JWT-аутентификация, хеширование паролей (bcrypt)
- API: RESTful на Golang с CORS middleware
- Данные: Хранение в MongoDB (пользователи, курсы, заявки, документы)

Отвечай кратко и по делу, используя только информацию о функционале системы. Если вопрос не связан с курсами, вежливо сообщи, что не можешь помочь.

Примеры ответов:
- 'Вы можете записаться на курс через кнопку "Записаться" на странице курса'
- 'Статус заявки можно проверить в личном кабинете'
- 'Администратор рассматривает заявки в течение 3 рабочих дней'
- 'Для восстановления пароля обратитесь в поддержку
`,
			},
			{
				Role: "user",
				Text: prompt,
			},
		},
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", yandexGPTAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", y.iamToken))
	req.Header.Set("x-folder-id", y.folderID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s, body: %s", resp.Status, string(body))
	}

	var completionResp CompletionResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Добавляем отладочный вывод
	fmt.Printf("Raw API response: %s\n", string(body))

	if err := json.Unmarshal(body, &completionResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v (body: %s)", err, string(body))
	}

	if len(completionResp.Result.Alternatives) == 0 {
		return "", fmt.Errorf("no alternatives in response")
	}

	return completionResp.Result.Alternatives[0].Message.Text, nil
}
