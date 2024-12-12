package service

import (
	"a21hc3NpZ25tZW50/model"
	"a21hc3NpZ25tZW50/repository"
	"fmt"
)

type UserService interface {
	Register(username string, email string, password string) error
	Login(email string,password string)(model.UserByEmail,error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) Register(username string, email string, password string) error {

	_, err := s.userRepository.GetUserByEmail(email)

	//jika user dengan email tertentu ditemukan, maka registrasi error
	if err == nil {

		return fmt.Errorf("email already registered")
	}

	errRegist := s.userRepository.Register(username, email, password)

	if errRegist != nil {
		return errRegist
	}

	return nil
}

func(s *userService) Login(email string,password string)(model.UserByEmail,error){
	
	user,err := s.userRepository.Login(email,password)
	//jika user dengan kombinasi email dan password tidak ditemukan
	if err!=nil{
		return model.UserByEmail{},err
	}
	return user,nil
}
//
