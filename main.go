package main

import (
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

// Initialize the services
var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}

func main() {
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

	// Set up the router
	router := mux.NewRouter()

	// File analyze endpoint
	router.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseMultipartForm(1024) 
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file content", http.StatusInternalServerError)
			return
		}

		resultFile, err := fileService.ProcessFile(string(content))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		question  := r.FormValue("query")


		tapasRequest,errTapasRequest := aiService.AnalyzeData(resultFile,question,token)
		if errTapasRequest!=nil{
			http.Error(w,errTapasRequest.Error(),500)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tapasRequest)
	}).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		
		var requestBody struct {
			Query   string `json:"query"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody)

		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	
		aiService := service.AIService{Client: &http.Client{}}

		response, err := aiService.ChatWithAI( requestBody.Query, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		
	}).Methods("POST")



	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
