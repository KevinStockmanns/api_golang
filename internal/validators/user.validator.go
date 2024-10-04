package validators

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/encryptor"
	"github.com/KevinStockmanns/api_golang/internal/models"
)

func UserSignUp(user models.User, data dtos.UserSignUpDTO) (int, dtos.ErrorsDTO) {
	validations := []ValidationFunc{
		UniqueValueInDB(user, "email", &data.Email, "el correo ya se encuentra en uso"),
		requiredAge(18, &data.Birthday, "birthday"),
	}

	for _, v := range validations {
		if status, errs := v(); status != http.StatusOK {
			return status, errs
		}
	}

	return http.StatusOK, dtos.ErrorsDTO{}
}

func UserLogin(user models.User, data dtos.UserLoginDTO) (int, dtos.ErrorsDTO) {
	validators := []ValidationFunc{
		userActive(user),
		validatePassword(data.Password, user.Password, "password"),
	}
	for _, v := range validators {
		if status, errs := v(); status != http.StatusOK {
			return status, errs
		}
	}
	return http.StatusOK, dtos.ErrorsDTO{}
}

func UserUpdate(user models.User, data dtos.UserUpdateDTO, idSent uint) (int, dtos.ErrorsDTO) {
	validators := []ValidationFunc{
		idCorrespondient(user, idSent),
		validateRol(user, data.Rol),
		UniqueValueInDB(user, "email", data.Email, "el correo ya se encuentra en uso"),
		OneDataRequired(data, "el requerido al menos un dato para actualizar el usuario"),
		requiredAge(18, data.Birthday, "birthday"),
	}
	for _, v := range validators {
		if status, errors := v(); status != http.StatusOK {
			return status, errors
		}
	}
	return http.StatusOK, dtos.ErrorsDTO{}
}

func validatePassword(password string, hash string, field string) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if err := encryptor.VerifyPassword(password, hash); err != nil {
			return http.StatusUnauthorized, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: field, Error: "credenciales inválidas"}}}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func requiredAge(requiredAge int8, initTime *string, field string) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if initTime == nil {
			return http.StatusOK, dtos.ErrorsDTO{}
		}

		date, err := time.Parse("2006-01-02", *initTime)
		if err != nil {
			return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{
				{Field: field, Error: fmt.Sprintf("fecha inválida: %s. Formato esperado: YYYY-MM-DD", *initTime)},
			}}
		}

		now := time.Now()
		age := int8(now.Year() - date.Year())

		if now.YearDay() < date.YearDay() {
			age--
		}

		if age < requiredAge {
			return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{
				{Field: field, Error: fmt.Sprintf("la edad mínima requerida es de %d años", requiredAge)},
			}}
		}

		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func userActive(user models.User) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if !user.Status && user.Rol.Name == string(constants.User) {
			return http.StatusUnauthorized, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{
				{Error: "no puedes ingresar porque la cuenta esta desactivada"},
			}}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func validateRol(user models.User, newRol *string) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if newRol != nil {
			if !user.IsAdmin() {
				return http.StatusUnauthorized, dtos.ErrorsDTO{
					Errors: []dtos.ErrorDTO{{Field: "rol", Error: "no tienes los permisos requeridos para cambiar el rol"}},
				}
			}
			newRolString := strings.ToUpper(*newRol)

			if _, exists := constants.UserRoles[constants.UserRole(newRolString)]; !exists {
				validRoles := []string{}
				for role := range constants.UserRoles {
					validRoles = append(validRoles, string(role))
				}
				return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: "rol", Error: fmt.Sprintf("el rol ingresado no existe las opciones validas son: %v", strings.Join(validRoles, ", "))}}}
			}
		}

		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func idCorrespondient(user models.User, idSent uint) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if !user.IsAdmin() && user.ID != idSent {
			return http.StatusUnauthorized, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Error: "no tienes los permisos necesarios para actualizar los datos de otro usuario"}}}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}
