package routes

import (
	"net/http"

	"github.com/KevinStockmanns/api_golang/handlers"
	"github.com/KevinStockmanns/api_golang/middlewares"
	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Echo) {
	e.GET("/product/:id", handlers.GetProduct)
	e.GET("/product", handlers.GetProducts)
	e.POST("/product", handlers.PostProduct)
	middlewares.ProtectEndPoint("/product", http.MethodPost)
	e.PUT("product/:id", handlers.PutProduct)
	e.DELETE("product/:id", handlers.DeleteProduct)
}
