package dtos

type UserPostDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=50,propername"`
	LastName string `json:"lastname" validate:"required,min=3,max=50,propername"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Birthday string `json:"birthday" validate:"required,date"`
	Phone    string `json:"phone" validate:"required,phone"`
}
