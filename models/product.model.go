package models

import (
	"github.com/KevinStockmanns/api_golang/models/wrapper"
	"github.com/KevinStockmanns/api_golang/utils"
	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"type:varchar(50);unique"`
	Status   bool      `json:"status"`
	Versions []Version `json:"versions"`
}

type PostProduct struct {
	Name     string        `json:"name" validate:"required,min=3,max=50,regexp=^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+( [a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+)*$"`
	Status   bool          `json:"status"`
	Versions []VersionPost `json:"versions" validate:"required,min=1,max=6,dive"`
}

type PutProduct struct {
	Name     *string       `json:"name" validate:"omitempty,min=3,max=50,regexp=^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+( [a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+)*$"`
	Status   *bool         `json:"status"`
	Versions *[]PutVersion `json:"versions" validate:"omitempty,min=1,dive"`
}

func (p *PutProduct) NormalizeAndValidate() []wrapper.ErrorWrapper {
	if p.Versions != nil {
		for i := range *p.Versions {
			(*p.Versions)[i].Normalize()
		}
	}

	if err := utils.Validate.Struct(p); err != nil {
		errors := utils.ValidateErrors(err.(validator.ValidationErrors))
		return errors
	}

	return make([]wrapper.ErrorWrapper, 0)
}
