package handler

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func getTokenHuggingFace()string{
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

func getSecretKey()string{
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("secret_key")
	log.Println(key)
	return key
}


func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {

	token := getTokenHuggingFace()

	var fileService = &service.FileService{}
	var aiService = &service.AIService{Client: &http.Client{}}
	
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

	question := r.FormValue("query")

	tapasRequest, errTapasRequest := aiService.AnalyzeData(resultFile, question, token)
	if errTapasRequest != nil {
		http.Error(w, errTapasRequest.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tapasRequest)
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {

	token := getTokenHuggingFace()
	
	var requestBody struct {
		Query string `json:"query"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	aiService := service.AIService{Client: &http.Client{}}

	response, err := aiService.ChatWithAI(requestBody.Query, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/// Berkaitan dengan user

type APIHandler struct {
	UserService service.UserService
}

func NewAPIHandler(userService service.UserService) *APIHandler {
	return &APIHandler{
		UserService: userService,
	}
}

func(h *APIHandler) RegisterHandler(w http.ResponseWriter, r *http.Request){
	var dataRegister model.Register
	
	//decode json
	err:=json.NewDecoder(r.Body).Decode(&dataRegister)
	
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	
	// Validasi input, pastikan semua field diisi
	if dataRegister.Username == "" || dataRegister.Email == "" || dataRegister.Password == "" {
		w.WriteHeader(400)
		errorResponse := model.Error{Error: "field should not be empty"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	// Panggil fungsi Register dari UserService
	err = h.UserService.Register(dataRegister.Username, dataRegister.Email, dataRegister.Password)
	if err != nil {
		// Jika error karena email sudah terdaftar
		if err.Error() == "email already registered" {
			w.WriteHeader(409)
			errorResponse := model.Error{Error: err.Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		
		// Error lainnya
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	

	w.Header().Set("Content-Type", "application/json")

	
	response := model.Success{Message: "Successfully registered the user"}
	json.NewEncoder(w).Encode(response)

}

func(h *APIHandler) LoginHandler(w http.ResponseWriter, r *http.Request){

	var dataLogin model.Login
	//decode json
	err:=json.NewDecoder(r.Body).Decode(&dataLogin)
	
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validasi input, pastikan semua field diisi
	if dataLogin.Email == "" || dataLogin.Password == "" {
		w.WriteHeader(400)
		errorResponse := model.Error{Error: "password and email should not be empty"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	// Panggil fungsi Register dari UserService
	user,err := h.UserService.Login(dataLogin.Email,dataLogin.Password)
	if err != nil {
		w.WriteHeader(400)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	tokenjwt,err:= GenerateTokenJWT(user.Username,user.Email)
	if err!=nil{
		log.Println("Error occured when generate token",err.Error())
		return
	}

	response:= model.LoginResponse{
		Username: user.Username,
		Email: user.Email,
		Token: "Bearer "+tokenjwt,
	}

	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(response)

}

func GenerateTokenJWT(username string, email string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute) // 5 menit

	// Buat claims berisi data username dan role yang akan kita embed ke JWT
	claims := &model.Claims{
		Username: username,
		Email:     email,
		StandardClaims: jwt.StandardClaims{
			// expiry time menggunakan time millisecond
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Buat token menggunakan encoded claim dengan salah satu algoritma yang dipakai
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Buat JWT string dari token yang sudah dibuat menggunakan JWT key yang telah dideklarasikan (proses encoding JWT)
	jwtkey:=[]byte(getSecretKey())
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		// return internal error ketika ada kesalahan saat pembuatan JWT string
		return "", err
	}

	return tokenString, nil
}

