package repository

import (
	"a21hc3NpZ25tZW50/model"
	"database/sql"
	"fmt"
	
	
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Register(username string, email string, password string)error
	GetUserByEmail(email string)(model.UserByEmail,error)
	Login(email string,password string)(model.UserByEmail,error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Register(username string, email string, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }

	query := `
		INSERT INTO users(username,email,password) VALUES ($1,$2,$3)
	`
	_, err = u.db.Exec(query,username,email,hashedPassword)

	if err != nil {
		return fmt.Errorf("registration failed")
	}

	return nil
}


func (u *userRepository) GetUserByEmail(email string)(model.UserByEmail,error){
	query :=`
		SELECT username,email,password FROM users WHERE email = $1
	`
	row := u.db.QueryRow(query,email)

	var user model.UserByEmail

	err := row.Scan(&user.Username,&user.Email,&user.Password)

	if err!=nil{
		//Error jika data tidak ditemukan
		if err == sql.ErrNoRows {
            return model.UserByEmail{}, fmt.Errorf("user with email %s not found", email)
        }
        //Error jika terjadi kesalahan query
        return model.UserByEmail{}, fmt.Errorf("failed to query user by email: %w", err)
	}
	return user,nil

}

func(u*userRepository)Login(email string,password string)(model.UserByEmail,error){

	var user model.UserByEmail
	
	
	user,err:= u.GetUserByEmail(email)

	if err!=nil{
		return model.UserByEmail{},err
	}

	// Bandingkan password yang dimasukkan dengan password tersimpan
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.UserByEmail{}, fmt.Errorf("invalid password")
	}

	return user, nil
}