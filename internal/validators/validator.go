package validators

import (
	"net/http"
	"reflect"

	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
)

type ValidationFunc func() (int, dtos.ErrorsDTO)

func UniqueValueInDB(model interface{}, tableName string, value *string, textError string) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if value != nil {
			var count int64
			db.DB.Model(model).Where(tableName+" = ?", *value).Count(&count)

			if count > 0 {
				errors := dtos.ErrorsDTO{
					Errors: []dtos.ErrorDTO{
						{Field: tableName, Error: textError},
					},
				}

				return http.StatusBadRequest, errors
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func OneDataRequired(model interface{}, errorText string) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		value := reflect.ValueOf(model)

		if value.Kind() == reflect.Struct {
			fieldCount := value.NumField()

			for i := 0; i < fieldCount; i++ {
				field := value.Field(i)

				if field.Kind() == reflect.Ptr && !field.IsNil() {
					return http.StatusOK, dtos.ErrorsDTO{}
				}
			}
			return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Error: errorText}}}
		}

		return http.StatusOK, dtos.ErrorsDTO{}
	}
}
