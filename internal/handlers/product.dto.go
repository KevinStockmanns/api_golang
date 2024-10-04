package handlers

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/labstack/echo/v4"
)

func ProductPost(c echo.Context) error {
	var productDto dtos.ProductCreateDTO
	if err := c.Bind(&productDto); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.ErrorResponse{Message: "ocurrio un error al leer el cuerpo de la petición"})
	}
	return c.JSON(http.StatusOK, productDto)
}
