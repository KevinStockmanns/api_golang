package dtos

import "time"

type UserSignUpDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=50,propername"`
	LastName string `json:"lastname" validate:"required,min=3,max=50,propername"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password,min=8,max=20"`
	Birthday string `json:"birthday" validate:"required,date"`
	Phone    string `json:"phone" validate:"required,phone"`
}

type UserResponseDTO struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	LastName string    `json:"lastname"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
	Status   bool      `json:"status"`
	Phone    string    `json:"phone"`
	Rol      string    `json:"rol"`
}

type UserWithTokenResponseDTO struct {
	Token string `json:"token"`
	UserResponseDTO
}

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password,min=8,max=20"`
}
