package handlers

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/models/dto"
	"github.com/labstack/echo/v4"
)

func RegisterUser(c echo.Context) error {
	var userDto dto.UserCreate

	if err := c.Bind(&userDto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "ocurrio un error al leer el cuerpo",
		})
	}
	return nil
}
