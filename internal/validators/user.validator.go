package validators

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
)

func UserValidations(user models.User, data dtos.UserPostDTO) (int, dtos.ErrorsDTO) {
	validations := []ValidationFunc{
		UniqueValueInDB(user, "email", data.Email, "el correo ya se encuentra en uso"),
	}

	for _, v := range validations {
		if status, errs := v(); status != http.StatusOK {
			return status, errs
		}
	}

	return http.StatusOK, dtos.ErrorsDTO{}
}
