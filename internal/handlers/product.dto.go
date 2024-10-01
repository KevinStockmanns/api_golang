package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ProductPost(c echo.Context) error {
	return c.JSON(http.StatusOK, "post product")
}
