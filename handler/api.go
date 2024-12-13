package handler

import (

	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	
)

func getTokenHuggingFace() string {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("myAI_Token")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}
	return token
}

func getSecretKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("secret_key")
	return key
}



func ChatHandler(w http.ResponseWriter, r *http.Request) {

	token := getTokenHuggingFace()
	
	aiService := service.AIService{Client: &http.Client{}}
	

	
	question := r.FormValue("query")
	document := r.FormValue("document")
	
	resultstring := document + "Bacalah format file tersebut! dan jawablah pertanyaan ini berdasarkan file tersebut, jika pertanyaan diluar konteks, maka berikan respon 'pertanyaan diluar konteks dokumen'. Berikut ini pertanyaannya:" + question

	
	
	responseAI, err := aiService.ChatWithAI(resultstring, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseToUser := map[string]interface{}{
		"responseAI" : responseAI.GeneratedText,
		"Question":question,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseToUser)
}




