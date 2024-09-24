package routes

import (
	"github.com/KevinStockmanns/api_golang/handlers"
	"github.com/labstack/echo/v4"
)

func Index(e *echo.Echo) {
	e.GET("/", handlers.Index)
}
