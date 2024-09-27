package routes

import (
	"github.com/KevinStockmanns/api_golang/handlers"
	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Group) {
	e.GET("/product/:id", handlers.GetProduct)
	e.GET("/product", handlers.GetProducts)
	e.POST("/product", handlers.PostProduct)
	e.PUT("product/:id", handlers.PutProduct)
	e.DELETE("product/:id", handlers.DeleteProduct)
}
