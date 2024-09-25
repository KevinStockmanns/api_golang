package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/KevinStockmanns/api_golang/db"
	"github.com/KevinStockmanns/api_golang/models"
	"github.com/KevinStockmanns/api_golang/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetProduct(c echo.Context) error {
	var product models.Product
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "El id debe ser un número")
	}
	if err := db.DB.Model(&models.Product{}).First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, "Producto no encontrado")
		} else {
			return c.JSON(http.StatusInternalServerError, "Error al buscar el producto")
		}
	}
	return c.JSON(200, &product)
}

func PostProduct(c echo.Context) error {
	var product models.Product
	var productPost models.PostProduct

	if err := c.Bind(&productPost); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(productPost); err != nil {
		errors := utils.ValidateErrors(err.(validator.ValidationErrors))
		if len(errors) > 0 {
			return c.JSON(http.StatusBadRequest, map[string][]models.ErrorWrapper{
				"errors": errors,
			})
		}
	}
	result := db.DB.Create(&product)
	if result.Error != nil {
		log.Println(result.Error)
		if strings.Contains(result.Error.Error(), "Error 1062") {
			return c.JSON(http.StatusBadRequest, "El nombre del produccto ya esta en uso")
		} else {
			return c.JSON(http.StatusInternalServerError, "error al agregar")
		}
	}

	return c.JSON(http.StatusOK, product)
}
