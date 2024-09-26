package utils

import (
	"log"
	"regexp"
	"strings"

	"github.com/KevinStockmanns/api_golang/models/wrapper"
	"github.com/go-playground/validator/v10"
)

func ValidateErrors(err validator.ValidationErrors) []wrapper.ErrorWrapper {
	var errors []wrapper.ErrorWrapper
	for _, e := range err {
		var error string
		param := e.Param()
		switch e.Tag() {
		case "required":
			error = "el campo es requerido"
		case "min":
			error = "el campo debe ser de al menos " + param + " de valor/longitud"
		case "max":
			error = "el campo debe ser hasta " + param + " de valor/longitud"
		case "gt":
			error = "el campo debe ser mayor a " + param
		case "regexp":
			error = "el cambo debe respetar este formato " + e.Param()
		case "oneof":
			error = "el campo solo acepta estas opciones: " + param
		default:
			error = "campo inválido"
		}
		field := e.Namespace()
		firstDotIndex := strings.Index(field, ".")
		if firstDotIndex != -1 {
			field = field[firstDotIndex+1:]
		}
		errors = append(errors, wrapper.ErrorWrapper{Field: strings.ToLower(field), Error: error})
	}
	return errors
}

var Validate *validator.Validate

func InitValidations() {
	Validate = validator.New()
	// Validador personalizado que permite múltiples expresiones regulares
	Validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		pattern := fl.Param() // Obtener el patrón desde la etiqueta

		// Validar con la expresión regular proporcionada
		matched, err := regexp.MatchString(pattern, value)
		if err != nil {
			log.Println("Error en la expresión regular:", err)
			return false
		}
		return matched
	})
}

func IsInt(s string) bool {
	re := regexp.MustCompile("^[0-9]+$")
	return re.MatchString(s)
}
