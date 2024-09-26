package validators

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/models"
	"github.com/KevinStockmanns/api_golang/models/wrapper"
)

func PutProductValidation(product models.Product, productDto models.PutProduct) (int, wrapper.ErrorWrapper) {
	validators := []ValidationFunc{
		UniqueField(product, "name", productDto.Name, "el nombre del producto se encuentra en uso"),
	}
	for _, v := range validators {
		if status, err := v(); status != http.StatusOK {
			return status, err
		}
	}

	return http.StatusOK, wrapper.ErrorWrapper{}
}
