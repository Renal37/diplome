package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	yandexGPTAPIURL = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
)

type YandexGPTClient struct {
	iamToken  string
	folderID  string
	modelName string
}

func NewYandexGPTClient(iamToken, folderID string) *YandexGPTClient {
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
			Message Message `json:"message"`
			Status  string  `json:"status"`
		} `json:"alternatives"`
		Usage struct {
			InputTextTokens  int `json:"inputTextTokens"`
			CompletionTokens int `json:"completionTokens"`
			TotalTokens      int `json:"totalTokens"`
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
				Text: "Ты - помощник на сайте курсов. Отвечай кратко и по делу. Основные темы: запись на курсы, стоимость, расписание, документы.",
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
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if len(completionResp.Result.Alternatives) == 0 {
		return "", fmt.Errorf("no alternatives in response")
	}

	return completionResp.Result.Alternatives[0].Message.Text, nil
}
