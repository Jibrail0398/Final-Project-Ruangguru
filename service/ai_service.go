package service

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	
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
			Table: table,
			Query: query,
		},
	}

	payload, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://api.example.com/analyze", bytes.NewBuffer(payload))
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

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	payload := map[string]string{
		"context": context,
		"query":   query,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return model.ChatResponse{}, err
	}

	req, err := http.NewRequest("POST", "http://api.example.com/chat", bytes.NewBuffer(payloadBytes))
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

	var chatResponses []model.ChatResponse
	err = json.Unmarshal(body, &chatResponses)
	if err != nil || len(chatResponses) == 0 {
		return model.ChatResponse{}, errors.New("invalid response from AI")
	}

	return chatResponses[0], nil
}
