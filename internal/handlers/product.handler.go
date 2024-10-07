package handlers

import (
	"fmt"
	"net/http"

	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
	"github.com/KevinStockmanns/api_golang/internal/services"
	"github.com/KevinStockmanns/api_golang/internal/utils"
	"github.com/KevinStockmanns/api_golang/internal/validators"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ProductCreate(c echo.Context) error {
	var productDto dtos.ProductCreateDTO
	if err := c.Bind(&productDto); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}

	if errs, ok := validators.ValidateDTOs(productDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}
	var product models.Product
	if status, errs := validators.ProductCreate(product, productDto); status != http.StatusOK {
		return c.JSON(http.StatusOK, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}
	product.Create(productDto)

	tx := db.DB.Begin()

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al guardar el producto"})
	}

	tx.Commit()

	var productResponse dtos.ProductResponseDTO
	productResponse.Init(product)

	c.Response().Header().Set("Location", fmt.Sprintf("/product/%d", product.ID))
	return c.JSON(http.StatusOK, productResponse)
}

func ProductGet(c echo.Context) error {
	idParam := c.Param("id")
	if !utils.IsInteger(idParam) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "el id debe ser un número entero"})
	}

	var product models.Product
	if err := db.DB.Preload("Versions").First(&product, idParam).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Message: "no se encontro el producto"})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al obtener el producto"})
		}
	}

	var productResponse dtos.ProductResponseDTO
	productResponse.Init(product)

	return c.JSON(http.StatusOK, productResponse)
}

func ProductUpViews(c echo.Context) error {
	var idVersions dtos.VersionUpViewDTO
	if err := c.Bind(&idVersions); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un errro al leer el cuerpo de la petición"})
	}
	if errs, ok := validators.ValidateDTOs(idVersions); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}
	var versions []models.Version
	tx := db.DB.Begin()
	if err := tx.Where("id IN ?", idVersions.IDVersions).Find(&versions); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "ocurrio un error al obtener las versiones",
		})
	}
	for i := range versions {
		versions[i].Views++
	}

	if err := tx.Save(&versions).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Message: "Ocurrió un error al guardar las versiones",
		})
	}

	tx.Commit()

	return c.NoContent(http.StatusNoContent)
}

func ProductUpdate(c echo.Context) error {
	idParam := c.Param("id")
	if !utils.IsInteger(idParam) {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "el id ingresado debe ser un número entero"})
	}
	var productDto dtos.ProductUpdateDTO
	if c.Bind(&productDto) != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}
	if errs, ok := validators.ValidateDTOs(productDto); !ok {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}

	var product models.Product
	tx := db.DB.Begin()

	if err := tx.Preload("Versions").First(&product, idParam).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Message: "no se encontro el producto en la base de datos"})
		} else {
			return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Message: "ocurrio un error al buscar el producto"})
		}
	}

	if status, errs := validators.ProductUpdate(product, productDto); status != http.StatusOK {
		return c.JSON(status, dtos.ErrorResponse{
			Message: "error de validación",
			Errors:  errs.Errors,
		})
	}

	services.ProductUpdate(&product, productDto, tx)

	tx.Commit()
	var productResponse dtos.ProductResponseDTO
	productResponse.Init(product)

	return c.JSON(http.StatusOK, productResponse)
}
