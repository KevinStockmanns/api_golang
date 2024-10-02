package dtos

import "time"

type UserModel interface {
	GetID() uint
	GetName() string
	GetLastName() string
	GetEmail() string
	GetBirthday() time.Time
	GetPhone() string
	GetRol() string
	IsStatus() bool
}

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

func (dto *UserResponseDTO) Init(user UserModel) {
	dto.ID = user.GetID()
	dto.Birthday = user.GetBirthday()
	dto.Email = user.GetEmail()
	dto.LastName = user.GetLastName()
	dto.Name = user.GetName()
	dto.Phone = user.GetPhone()
	dto.Rol = user.GetRol()
	dto.Status = user.IsStatus()
}

type UserWithTokenResponseDTO struct {
	Token string `json:"token"`
	UserResponseDTO
}

type UserLoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password,min=8,max=20"`
}

type UserUpdateDTO struct {
	Name     *string `json:"name" validate:"omitempty,required,min=3,max=50,propername"`
	LastName *string `json:"lastname" validate:"omitempty,required,min=3,max=50,propername"`
	Email    *string `json:"email" validate:"omitempty,required,email"`
	Birthday *string `json:"birthday" validate:"omitempty,required,date"`
	Status   *bool   `json:"status"`
	Phone    *string `json:"phone" validate:"omitempty,required,phone"`
	Rol      *string `json:"rol" validate:"omitempty,required"`
}
