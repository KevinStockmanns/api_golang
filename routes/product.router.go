package routes

import (
	"github.com/KevinStockmanns/api_golang/handlers"
	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Echo) {
	e.GET("/product/:id", handlers.GetProduct)
	e.POST("/product", handlers.PostProduct)
}
