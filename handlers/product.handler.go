package handlers

import (
	"net/http"
	"strconv"

	"github.com/KevinStockmanns/api_golang/db"
	"github.com/KevinStockmanns/api_golang/models"
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
