package main

import (
	
	"a21hc3NpZ25tZW50/handler"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)



func main() {
	
	
	

	// Set up the router
	router := mux.NewRouter()

	// File analyze endpoint
	router.HandleFunc("/analyze", handler.AnalyzeHandler).Methods("POST")

	// Chat endpoint
	router.HandleFunc("/chat", handler.ChatHandler).Methods("POST")



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
