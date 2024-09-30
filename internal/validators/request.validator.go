package validators

import (
	"strings"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func ValidateDTOs(model interface{}) (dtos.ErrorsDTO, bool) {
	Validate = validator.New()

	if err := Validate.Struct(model); err != nil {
		var errors dtos.ErrorsDTO
		for _, e := range err.(validator.ValidationErrors) {
			errorDto := dtos.ErrorDTO{
				// Field: e.Namespace()[strings.Index(e.Namespace(), ".")+1:],
				Field: strings.ToLower(e.Namespace()[strings.Index(e.Namespace(), ".")+1:][:1]) + e.Namespace()[strings.Index(e.Namespace(), ".")+1:][1:],
			}
			switch e.Tag() {
			case "required":
				errorDto.Error = "el dato es requerido"
			default:
				errorDto.Error = "Dato inválido"

			}
			errors.Errors = append(errors.Errors, errorDto)
		}
		return errors, false
	}

	return dtos.ErrorsDTO{}, true
}
