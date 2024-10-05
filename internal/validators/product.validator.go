package validators

import (
	"fmt"
	"net/http"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
)

func ProductCreate(product models.Product, data dtos.ProductCreateDTO) (int, dtos.ErrorsDTO) {
	statusDto := data.Status
	validatros := []ValidationFunc{
		UniqueValueInDB(product, "name", &data.Name, "el nombre del producto se encuentra en uso"),
		oneVersionActive(data, &statusDto),
		validatePrices(data.Versions),
		noRepeatVersion(data.Versions),
	}

	for _, v := range validatros {
		if status, errs := v(); status != http.StatusOK {
			return status, errs
		}
	}

	return http.StatusOK, dtos.ErrorsDTO{}
}

func oneVersionActive(data dtos.ProductCreateDTO, statusDto *bool) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if statusDto != nil && *statusDto {
			oneVersionActive := false
			for _, vDto := range data.Versions {
				if vDto.Status {
					oneVersionActive = true
					break
				}
			}
			if !oneVersionActive {
				return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: "versions", Error: "para activar el producto es necesario tener al menos una versión activa"}}}
			}
		}

		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func validatePrices(versions []dtos.VersionCreateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		for i, vDto := range versions {
			if vDto.ResalePrice != nil && *vDto.ResalePrice >= vDto.Price {
				return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: fmt.Sprintf("versions[%d].resalePrice", i), Error: "el precio debe ser mayor al original"}}}
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func noRepeatVersion(versions []dtos.VersionCreateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		for i, vDto := range versions {
			for j, vDto2 := range versions {
				if i != j && vDto2.Name == vDto.Name {
					return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{
						Field: fmt.Sprintf("versiones[%d].name", j),
						Error: "el nombre de la versión no se debe repetir por cada producto",
					}}}
				}
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}
