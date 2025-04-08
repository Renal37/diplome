package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	"github.com/Renal37/db"
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

	collection := db.GetCollection(db.Prompts)
	var dbPrompt struct {
		Content string `bson:"content"`
	}
	err := collection.FindOne(ctx, bson.M{}).Decode(&dbPrompt)

	systemPrompt := `Ты - AI-ассистент...` // дефолтный промпт
	if err == nil && dbPrompt.Content != "" {
		systemPrompt = dbPrompt.Content
	}

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
				Text: systemPrompt,
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
