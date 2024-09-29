package handlers

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/labstack/echo/v4"
)

func UserPost(c echo.Context) error {
	var userDto dtos.UserPostDTO

	if err := c.Bind(&userDto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "no se pudo leer el cuerpo enviado",
		})
	}

	c.Response().Header().Set("Location", "/user/1")
	return c.JSON(http.StatusCreated, userDto)
}
