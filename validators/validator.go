package validators

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/db"
	"github.com/KevinStockmanns/api_golang/models/wrapper"
)

type ValidationFunc func() (int, wrapper.ErrorWrapper)

func UniqueField(model interface{}, column string, value *string, errorText string) ValidationFunc {
	return func() (int, wrapper.ErrorWrapper) {
		if value == nil {
			return http.StatusOK, wrapper.ErrorWrapper{}
		}
		var cant int64 = 0

		db.DB.Model(model).Where(column+" = ?", value).Count(&cant)

		if cant > 0 {
			return http.StatusBadRequest, wrapper.ErrorWrapper{Field: column, Error: errorText}
		}

		return http.StatusOK, wrapper.ErrorWrapper{}
	}
}
