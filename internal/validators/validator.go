package validators

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
)

type ValidationFunc func() (int, dtos.ErrorsDTO)

func UniqueValueInDB(model interface{}, tableName string, value string, textError string) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		var count int64
		db.DB.Model(model).Where(tableName+" = ?", value).Count(&count)

		if count > 0 {
			errors := dtos.ErrorsDTO{
				Errors: []dtos.ErrorDTO{
					{Field: tableName, Error: textError},
				},
			}

			return http.StatusBadRequest, errors
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}
