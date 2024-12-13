package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}


func (s *AIService) ChatWithAI( query, token string) (model.ChatResponse, error) {

	payload := model.ChatAIRequest{
		Messages: []model.MessagesAIRequest{},
		MaxTokens: 500,
	}

	payload.Messages = append(payload.Messages,model.MessagesAIRequest{
		Role: "user",
		Content: query,
	})


	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return model.ChatResponse{}, err
	}

	url := "https://api-inference.huggingface.co/models/Qwen/Qwen2.5-Coder-32B-Instruct/v1/chat/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return model.ChatResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.ChatResponse{}, errors.New("failed to chat with AI: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ChatResponse{}, err
	}

	var chatapiResponse model.ChatAPIResponse

	err = json.Unmarshal([]byte(body), &chatapiResponse)
	if err != nil {
		return model.ChatResponse{},errors.New("failed Parsing Response")
	}

	var chatResponses model.ChatResponse
	for _,choice := range chatapiResponse.Choices{
		chatResponses.GeneratedText = append(chatResponses.GeneratedText, choice.Message.Content)
	}

	
	return chatResponses, nil
}
