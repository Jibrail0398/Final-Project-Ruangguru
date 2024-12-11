package model

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type AIRequest struct {
	Inputs Inputs `json:"inputs"`
}

type TapasResponse struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

type ChatResponse struct {
	GeneratedText []string `json:"generated_text"`
}

type ChatAPIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type MessagesAIRequest struct{
	Role string `json:"role"`
	Content string `json:"content"`
}

type ChatAIRequest struct{
	Messages []MessagesAIRequest `json:"messages"`
	MaxTokens int `json:"max_tokens"`
}

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	
}