package dtos

type UserPostDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	LastName string `json:"lastname" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Birthday string `json:"birthday" validate:"required"`
}
