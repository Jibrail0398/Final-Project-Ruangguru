package handler

import (

	"a21hc3NpZ25tZW50/service"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"strconv"
	"fmt"
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



func(h *ChatCRUDHandler) ChatHandler(w http.ResponseWriter, r *http.Request) {

	token := getTokenHuggingFace()
	
	aiService := service.AIService{Client: &http.Client{}}

	
	question := r.FormValue("query")
	document := r.FormValue("document")
	id_report := r.FormValue("id_report")
	id_user := r.FormValue("id_user")
	date := r.FormValue("date")
	
	id_report_Int,_ := strconv.Atoi(id_report)
	id_user_Int,_ := strconv.Atoi(id_user)

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
	
	responseAIString := fmt.Sprintf("%v", responseToUser["responseAI"])
	errSave := h.ChatService.SaveChat(date,question,responseAIString,id_user_Int,id_report_Int)
	if errSave!=nil{
		w.WriteHeader(400)
		response := model.Error{Error: errSave.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseToUser)
}




