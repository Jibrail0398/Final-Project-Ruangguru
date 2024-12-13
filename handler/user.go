package handler

import(
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

type APIHandler struct {
	UserService service.UserService
}

func NewAPIHandler(userService service.UserService) *APIHandler {
	return &APIHandler{
		UserService: userService,
	}
}

/// Berkaitan dengan user

func (h *APIHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var dataRegister model.Register

	//decode json
	err := json.NewDecoder(r.Body).Decode(&dataRegister)

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

func (h *APIHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

	var dataLogin model.Login
	//decode json
	err := json.NewDecoder(r.Body).Decode(&dataLogin)

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

	// Panggil fungsi Login dari UserService
	user, err := h.UserService.Login(dataLogin.Email, dataLogin.Password)
	if err != nil {
		w.WriteHeader(400)
		errorResponse := model.Error{Error: err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	tokenjwt, err := GenerateTokenJWT(user.Username, user.Email)
	if err != nil {
		log.Println("Error occured when generate token", err.Error())
		return
	}

	response := model.LoginResponse{
		Id: user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    "Bearer " + tokenjwt,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}

func GenerateTokenJWT(username string, email string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute) // 5 menit

	// Buat claims berisi data username dan role yang akan kita embed ke JWT
	claims := &model.Claims{
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			// expiry time menggunakan time millisecond
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Buat token menggunakan encoded claim dengan salah satu algoritma yang dipakai
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Buat JWT string dari token yang sudah dibuat menggunakan JWT key yang telah dideklarasikan (proses encoding JWT)
	jwtkey := []byte(getSecretKey())
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		// return internal error ketika ada kesalahan saat pembuatan JWT string
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(401)
			errorResponse := model.Error{Error: "Token is empty"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// Periksa format token (Bearer)
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.WriteHeader(401)
			errorResponse := model.Error{Error: "Invalid token format"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		tokenString := headerParts[1]

		// Parse dan validasi token
		claims := &model.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Gunakan fungsi yang sama dengan generate token
			return []byte(getSecretKey()), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(401)
				errorResponse := model.Error{Error: err.Error()}
				json.NewEncoder(w).Encode(errorResponse)
				return
			}
			w.WriteHeader(401)
			errorResponse := model.Error{Error: err.Error()}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// Periksa apakah token valid
		if !token.Valid {
			w.WriteHeader(401)
			errorResponse := model.Error{Error: "Invalid token"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		// Periksa apakah token sudah expire
		if claims.ExpiresAt < time.Now().Unix() {
			w.WriteHeader(401)
			errorResponse := model.Error{Error: "Token has expired"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		

		// Lanjutkan ke handler selanjutnya jika token valid
		next.ServeHTTP(w, r)
	}
}