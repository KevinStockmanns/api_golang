package dtos

import "time"

type UserPostDTO struct {
	Name     string    `json:"name" validate:"required,min=3,max=50"`
	LastName string    `json:"lastname" validate:"required,min=3,max=50"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required"`
	Birthday time.Time `json:"birthday" validate:"required"`
}
