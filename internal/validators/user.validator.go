package validators

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
)

func UserSignUp(user models.User, data dtos.UserSignUpDTO) (int, dtos.ErrorsDTO) {
	validations := []ValidationFunc{
		UniqueValueInDB(user, "email", data.Email, "el correo ya se encuentra en uso"),
		requiredAge(18, user.Birthday, "birthday"),
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
	}
	for _, v := range validators {
		if status, errs := v(); status != http.StatusOK {
			return status, errs
		}
	}
	return http.StatusOK, dtos.ErrorsDTO{}
}

func requiredAge(requiredAge int8, date time.Time, field string) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {

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
