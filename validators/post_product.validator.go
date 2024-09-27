package validators

import (
	"fmt"
	"net/http"

	"github.com/KevinStockmanns/api_golang/models"
	"github.com/KevinStockmanns/api_golang/models/wrapper"
)

func PostValidations(productDto models.PostProduct) (int, wrapper.ErrorWrapper) {
	validators := []ValidationFunc{
		UniqueValieInDB(models.Product{}, "name", &productDto.Name, "el nombre del producto ya se encuentra en uso"),
		isActive(productDto),
		nonRepeatVersion(productDto.Versions),
		priceResale(productDto),
	}
	for _, validator := range validators {
		if status, err := validator(); status != http.StatusOK {
			return status, err
		}
	}
	return http.StatusOK, wrapper.ErrorWrapper{}
}

func isActive(data models.PostProduct) ValidationFunc {
	return func() (int, wrapper.ErrorWrapper) {
		if data.Status {
			isOneVersionEnabled := false
			for _, v := range data.Versions {
				if v.Status {
					isOneVersionEnabled = true
				}
			}
			if !isOneVersionEnabled {
				return http.StatusBadRequest, wrapper.ErrorWrapper{Error: "es necesario una versión activa"}
			}
		}
		return http.StatusOK, wrapper.ErrorWrapper{}
	}
}
func nonRepeatVersion(versions []models.VersionPost) ValidationFunc {
	return func() (int, wrapper.ErrorWrapper) {
		for i, v := range versions {
			for j, v2 := range versions {
				if j != i {
					if v2.Name == v.Name {
						return http.StatusBadRequest, wrapper.ErrorWrapper{
							Field: fmt.Sprintf("versions[%d].name", j),
							Error: "el nombre no debe ser repetido",
						}
					}
				}
			}
		}
		return http.StatusOK, wrapper.ErrorWrapper{}
	}
}
func priceResale(product models.PostProduct) ValidationFunc {
	return func() (int, wrapper.ErrorWrapper) {
		for i, v := range product.Versions {
			if v.ResalePrice != nil && *v.ResalePrice >= v.Price {
				return http.StatusBadRequest, wrapper.ErrorWrapper{
					Field: fmt.Sprintf("versions[%d]", i),
					Error: "el precio de reventa no debe ser mayor al precio",
				}
			}
		}
		return http.StatusOK, wrapper.ErrorWrapper{}
	}
}
