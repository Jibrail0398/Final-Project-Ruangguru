package main

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/handler"
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"a21hc3NpZ25tZW50/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
	
)



func main() {

	//make connection database
	database := db.NewDatabase()

	credential := model.Credential{
		Host		:"localhost",
		Username	:"postgres",
		Password	:"jibrailadji02",
		DatabaseName:"FP-Ruangguru",
		Port		: 5432,
		
	}

	conn,err := database.Connect(&credential)


	if err!=nil{
		panic(err)
	}

	//Database Migration
	err = database.Migrate()
	if err!=nil{
		fmt.Println(err)
	}
	
	// Inisialisasi repository
	userRepo := repository.NewUserRepo(conn)

	// Inisialisasi service
	userService := service.NewUserService(userRepo)

	// Inisialisasi handler
	apiHandler := handler.NewAPIHandler(userService)

	// Set up the router
	router := mux.NewRouter()

	//register endpoint
	router.HandleFunc("/register",apiHandler.RegisterHandler).Methods("POST")

	//login endpoint
	router.HandleFunc("/login",apiHandler.LoginHandler).Methods("POST")

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
