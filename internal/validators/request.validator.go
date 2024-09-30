package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidations() {
	Validate = validator.New()

	Validate.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}", fl.Field().String())
		return matched
	})

	Validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		regexPattern := `^\+?(\d{1,3})?[-. ]?(\(?\d{1,4}\)?)?[-. ]?\d{1,4}[-. ]?\d{1,4}[-. ]?\d{1,9}$`
		matched, _ := regexp.MatchString(regexPattern, fl.Field().String())
		return matched
	})

	Validate.RegisterValidation("propername", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.MatchString("^[A-Za-zÁÉÍÓÚáéíóúÑñ]+([ '-][A-Za-zÁÉÍÓÚáéíóúÑñ]+)*$", fl.Field().String())
		return matched
	})

	Validate.RegisterValidation("objectname", func(fl validator.FieldLevel) bool {
		matched, _ := regexp.MatchString("^[A-Za-zÁÉÍÓÚáéíóúÑñ0-9]+([ '-][A-Za-zÁÉÍÓÚáéíóúÑñ0-9]+)*$", fl.Field().String())
		return matched
	})
}

func ValidateDTOs(model interface{}) (dtos.ErrorsDTO, bool) {
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
			case "min":
				errorDto.Error = fmt.Sprintf("el dato debe tener una longitud mínima de %s", e.Param())
			case "max":
				errorDto.Error = fmt.Sprintf("el dato puede tener una logintud máxima de %s", e.Param())
			case "email":
				errorDto.Error = "el formato es inválido"
			case "date":
				errorDto.Error = "la fecha es inválida. debe ser YYYY-MM-DD"
			case "propername":
				errorDto.Error = "el dato acepta solo letras y espacios en blanco"
			case "objectname":
				errorDto.Error = "el dato acepta solo letras, números y espacios en blanco"
			default:
				errorDto.Error = "Dato inválido"

			}
			errors.Errors = append(errors.Errors, errorDto)
		}
		return errors, false
	}

	return dtos.ErrorsDTO{}, true
}
