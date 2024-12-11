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


func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	

	if len(table) == 0 {
		return "", errors.New("table is empty")
	}

	inputs := model.AIRequest{
		Inputs: model.Inputs{
			Query: query,
			Table: table,
		},
	}

	payload, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-large-finetuned-wtq", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to analyze data: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response model.TapasResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if len(response.Cells) == 0 {
		return "", errors.New("no results found")
	}

	return response.Cells[0], nil
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
