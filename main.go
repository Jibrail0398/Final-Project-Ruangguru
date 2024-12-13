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
	reportRepo := repository.NewReportRepository(conn)
	chatRepo := repository.NewChatRepository(conn)
	// Inisialisasi service
	userService := service.NewUserService(userRepo)
	reportService := service.NewReportService(reportRepo)
	chatService := service.NewChatService(chatRepo)
	// Inisialisasi handler
	apiHandler := handler.NewAPIHandler(userService)
	reportHandler := handler.NewReportHandler(reportService)
	chatCRUDHandler := handler.NewChatCRUDHandler(chatService)

	// Set up the router
	router := mux.NewRouter()

	//register endpoint
	router.HandleFunc("/register",apiHandler.RegisterHandler).Methods("POST")

	//login endpoint
	router.HandleFunc("/login",apiHandler.LoginHandler).Methods("POST")

	//upload report endpoint
	router.HandleFunc("/upload",handler.AuthMiddleware(reportHandler.Upload)).Methods("POST")

	// getreportByuser endpoint
	router.HandleFunc("/get/report/{id}",handler.AuthMiddleware(reportHandler.GetReportByUser)).Methods("GET")

	//delete report endpoint
	router.HandleFunc("/delete/report/{id}",handler.AuthMiddleware(reportHandler.Delete)).Methods("DELETE")

	// Chat endpoint
	router.HandleFunc("/chat", handler.AuthMiddleware(handler.ChatHandler)).Methods("POST")

	//save chat endpoint
	router.HandleFunc("/chat/save",handler.AuthMiddleware(chatCRUDHandler.SaveChat)).Methods("POST")

	//Get chat by report endpoint
	router.HandleFunc("/chat/get/{id}",handler.AuthMiddleware(chatCRUDHandler.GetChatByReport)).Methods("GET")

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
