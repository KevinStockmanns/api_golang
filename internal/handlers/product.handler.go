package handlers

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
	"github.com/KevinStockmanns/api_golang/internal/validators"
	"github.com/labstack/echo/v4"
)

func ProductPostHandler(c echo.Context) error {
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

	return c.JSON(http.StatusOK, product)
}
