package validators

import (
	"net/http"
	"strconv"

	"github.com/KevinStockmanns/api_golang/models"
	"github.com/KevinStockmanns/api_golang/models/wrapper"
)

func PutProductValidation(product models.Product, productDto models.PutProduct) (int, wrapper.ErrorWrapper) {
	validators := []ValidationFunc{
		validateActions(productDto),
		UniqueField(product, "name", productDto.Name, "el nombre del producto se encuentra en uso"),
		hasOneVersionActive(product, productDto),
		idCorrespondient(product, productDto),
	}
	for _, v := range validators {
		if status, err := v(); status != http.StatusOK {
			return status, err
		}
	}

	return http.StatusOK, wrapper.ErrorWrapper{}
}

func hasOneVersionActive(product models.Product, productDto models.PutProduct) ValidationFunc {
	return func() (int, wrapper.ErrorWrapper) {
		valids := 0

		if product.Status || (productDto.Status != nil && *productDto.Status) {
			for _, v := range product.Versions {
				for _, vDto := range *productDto.Versions {
					if vDto.ID != nil && v.ID == *vDto.ID {
						if (vDto.Status != nil && !*vDto.Status) || vDto.Action == "delete" {
							continue
						}
						if v.Status || (vDto.Status != nil && *vDto.Status) || vDto.Action == "create" {
							valids++
						}
					}
				}
			}
		} else {
			valids++
		}
		if valids == 0 {
			return http.StatusBadRequest, wrapper.ErrorWrapper{Field: "versions", Error: "Debe haber al menos una versión activa"}
		}
		return http.StatusOK, wrapper.ErrorWrapper{}
	}
}

func validateActions(productoDto models.PutProduct) ValidationFunc {
	return func() (int, wrapper.ErrorWrapper) {
		for i, v := range *productoDto.Versions {
			switch v.Action {
			case "create":
				if v.Name == nil {
					return http.StatusBadRequest, wrapper.ErrorWrapper{Field: "versions[" + strconv.Itoa(i) + "].name", Error: "el nombre es requerido para agregar"}
				}
				if v.Price == nil {
					return http.StatusBadRequest, wrapper.ErrorWrapper{Field: "versions[" + strconv.Itoa(i) + "].price", Error: "el precio es requerido para agregar"}
				}
			case "update":
				if v.ID == nil {
					return http.StatusBadRequest, wrapper.ErrorWrapper{Field: "versions[" + strconv.Itoa(i) + "].id", Error: "el id es requerido para actualizar"}
				}
			case "delete":
				if v.ID == nil {
					return http.StatusBadRequest, wrapper.ErrorWrapper{Field: "versions[" + strconv.Itoa(i) + "].id", Error: "el id es requerido para eliminar"}
				}
			}

		}
		return http.StatusOK, wrapper.ErrorWrapper{}
	}
}
func idCorrespondient(product models.Product, productDto models.PutProduct) ValidationFunc {
	return func() (int, wrapper.ErrorWrapper) {

		for i, vDto := range *productDto.Versions {
			if vDto.Action == "create" {
				continue
			}
			find := false
			for _, v := range product.Versions {
				if vDto.ID != nil {
					if v.ID == *vDto.ID {
						find = true
						break
					}
				}
			}
			if !find {
				return http.StatusNotFound, wrapper.ErrorWrapper{Field: "versions[" + strconv.Itoa(i) + "].id", Error: "el id no coincide con una versión del producto"}
			}
		}

		return http.StatusOK, wrapper.ErrorWrapper{}
	}
}
