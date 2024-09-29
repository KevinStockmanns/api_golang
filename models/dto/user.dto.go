package dto

import (
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/models/wrapper"
	"github.com/KevinStockmanns/api_golang/utils"
	"github.com/go-playground/validator/v10"
)

type UserCreate struct {
	Name        string    `json:"name" validate:"required,min=3,max=100,regexp=^[a-zA-ZñÑáéíóúÁÉÍÓÚ]+( [a-zA-ZñÑáéíóúÁÉÍÓÚ]+)*$"`
	LastName    string    `json:"lastName" validate:"regexp=^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+( [a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+)*$,required,min=3,max=50"`
	Birthday    time.Time `json:"birthday" validate:"required"`
	Password    string    `json:"password" validate:"required,regexp=^[ñÑA-Za-z0-9-_]$,min=8,max=18"`
	Email       string    `json:"email" validate:"required,email"`
	NumberPhone string    `json:"numberPhone" validate:"required,regexp=^\\+?[1-9]\\d{0,2}[-.\\s]?(\\(?\\d{1,4}?\\)?[-.\\s]?)?\\d{1,4}[-.\\s]?\\d{1,4}[-.\\s]?\\d{1,9}$"`
}

func (u *UserCreate) NormalizeAndValidate() []wrapper.ErrorWrapper {
	u.Name = strings.Trim(u.Name, " ")
	u.LastName = strings.Trim(u.LastName, " ")
	u.Email = strings.Trim(u.Email, " ")
	u.NumberPhone = strings.Trim(u.NumberPhone, " ")

	if err := utils.Validate.Struct(u); err != nil {
		errors := utils.ValidateErrors(err.(validator.ValidationErrors))
		return errors
	}

	return make([]wrapper.ErrorWrapper, 0)
}
