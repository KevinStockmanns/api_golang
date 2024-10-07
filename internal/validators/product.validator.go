package validators

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
)

func ProductCreate(product models.Product, data dtos.ProductCreateDTO) (int, dtos.ErrorsDTO) {
	statusDto := data.Status
	validatros := []ValidationFunc{
		UniqueValueInDB(product, "name", &data.Name, "el nombre del producto se encuentra en uso"),
		oneVersionActiveCreate(data, &statusDto),
		validatePricesCreate(data.Versions),
		noRepeatVersionCreate(data.Versions),
	}

	for _, v := range validatros {
		if status, errs := v(); status != http.StatusOK {
			return status, errs
		}
	}

	return http.StatusOK, dtos.ErrorsDTO{}
}

func ProductUpdate(product models.Product, data dtos.ProductUpdateDTO) (int, dtos.ErrorsDTO) {
	validators := []ValidationFunc{
		validateActions(data.Versions),
		requiredDataUpdate(data),
		idIntegrityProductUpdate(product, data),
		UniqueValueInDB(product, "name", data.Name, "el nombre del producto ya esta en uso"),
		oneVersionActiveUpdate(product, data),
		noRepeatVersionUpdate(product.Versions, data.Versions),
		validatePricesUpdate(product.Versions, data.Versions),
	}
	for _, v := range validators {
		if status, err := v(); status != http.StatusOK {
			return status, err
		}
	}

	return http.StatusOK, dtos.ErrorsDTO{}
}

func oneVersionActiveCreate(data dtos.ProductCreateDTO, statusDto *bool) ValidationFunc {
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

func oneVersionActiveUpdate(product models.Product, productDto dtos.ProductUpdateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if productDto.Versions != nil {
			versionsActive := 0
			for _, v := range product.Versions {
				if v.Status {
					versionsActive++
					for _, vDto := range *productDto.Versions {
						if vDto.ID != nil && v.ID == *vDto.ID {
							if vDto.Status != nil && !*vDto.Status {
								versionsActive--
							}
							if vDto.Action == string(constants.Delete) {
								versionsActive--
							}
						}
					}
				}
			}
			for _, vDto := range *productDto.Versions {
				if vDto.Action == string(constants.Create) && vDto.Status != nil && *vDto.Status {
					versionsActive++
				}
			}

			if versionsActive <= 0 && product.Status && (productDto.Status != nil && !*productDto.Status) {
				return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Error: "es requerido al menos una versión activa"}}}
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func validatePricesCreate(versions []dtos.VersionCreateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		for i, vDto := range versions {
			if vDto.ResalePrice != nil && *vDto.ResalePrice >= vDto.Price {
				return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: fmt.Sprintf("versions[%d].resalePrice", i), Error: "el precio debe ser mayor al original"}}}
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func validatePricesUpdate(version []models.Version, versionDto *[]dtos.VersionUpdateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if versionDto != nil {
			for i, vDto := range *versionDto {
				if vDto.ResalePrice != nil {
					if vDto.Price != nil && *vDto.Price < *vDto.ResalePrice {
						return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: fmt.Sprintf("versions[%d].resalePrice", i), Error: "el precio de reventa debe ser mayor al precio"}}}
					} else {
						if vDto.Action == string(constants.Update) {
							for _, v := range version {
								if v.ID == *vDto.ID && *vDto.ResalePrice > v.Price {
									return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: fmt.Sprintf("versions[%d].resalePrice", i), Error: "el precio de reventa debe ser mayor al precio"}}}
								}
							}
						}
					}
				}
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func noRepeatVersionCreate(versions []dtos.VersionCreateDTO) ValidationFunc {
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

func noRepeatVersionUpdate(versions []models.Version, versionsDto *[]dtos.VersionUpdateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {

		if versionsDto != nil {
			versionsName := make(map[uint]string)

			for _, v := range versions {
				versionsName[v.ID] = v.Name
			}

			for i, vDto := range *versionsDto {
				if vDto.Action == string(constants.Create) || (vDto.Action == string(constants.Update) && vDto.Name != nil) {
					for key, value := range versionsName {
						if value == *vDto.Name {
							willChange := false
							for _, vDto2 := range *versionsDto {
								if vDto2.Action == string(constants.Update) && *vDto2.ID == key {
									if vDto2.Name != nil && *vDto2.Name != value {
										willChange = true
									}
								}
							}
							if !willChange {
								return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: fmt.Sprintf("versions[%d].name", i), Error: "el nombre de la versión no se puede repetir"}}}
							}
						}
					}
				}
			}
		}

		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func validateActions(versions *[]dtos.VersionUpdateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if versions != nil {
			var errors dtos.ErrorsDTO
			var validActions []string
			for action := range constants.Actions {
				validActions = append(validActions, string(action))
			}

			for i, vDto := range *versions {
				if _, isValid := constants.Actions[constants.Action(vDto.Action)]; !isValid {
					errors.Errors = append(errors.Errors, dtos.ErrorDTO{
						Field: fmt.Sprintf("versions[%d].action", i),
						Error: "la acción ingresada no es válida. Las opciones validas son: " + strings.Join(validActions, ", "),
					})
				}
			}
			if len(errors.Errors) > 0 {
				return http.StatusBadRequest, errors
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func requiredDataUpdate(data dtos.ProductUpdateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if data.Name == nil && data.Status == nil && data.Versions == nil {
			return http.StatusBadRequest, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Error: "es requerido al menos un dato para actualizar"}}}
		}

		if data.Versions != nil {
			var errors dtos.ErrorsDTO
			for i, vDto := range *data.Versions {
				field := fmt.Sprintf("versions[%d].", i)
				if vDto.Action == string(constants.Create) {
					if vDto.Name == nil {
						errors.Errors = append(errors.Errors, dtos.ErrorDTO{Field: field + "name", Error: "el nombre es requerido para crear una versión"})
					}
					if vDto.Price == nil {
						errors.Errors = append(errors.Errors, dtos.ErrorDTO{Field: field + "price", Error: "el precio es requerido"})
					}
				}
				if vDto.Action == string(constants.Update) {
					if vDto.ID == nil {
						errors.Errors = append(errors.Errors, dtos.ErrorDTO{Field: field + "id", Error: "es requerido el id para actualizar la versión"})
					}
				}
				if vDto.Action == string(constants.Delete) {
					if vDto.ID == nil {
						errors.Errors = append(errors.Errors, dtos.ErrorDTO{Field: field + "id", Error: "es requerido el id para desactivar la versión"})
					}
				}
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}

func idIntegrityProductUpdate(product models.Product, data dtos.ProductUpdateDTO) ValidationFunc {
	return func() (int, dtos.ErrorsDTO) {
		if data.Versions != nil {
			for i, vDto := range *data.Versions {
				isIn := false
				for _, v := range product.Versions {
					if vDto.ID != nil && *vDto.ID == v.ID {
						isIn = true
					}
				}
				if !isIn {
					return http.StatusNotFound, dtos.ErrorsDTO{Errors: []dtos.ErrorDTO{{Field: fmt.Sprintf("versiones[%d].id", i), Error: "la versión del producto no se encontro en la base de datos"}}}
				}
			}
		}
		return http.StatusOK, dtos.ErrorsDTO{}
	}
}
