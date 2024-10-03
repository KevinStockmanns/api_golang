package v1

import (
	"github.com/KevinStockmanns/api_golang/internal/handlers"
	"github.com/KevinStockmanns/api_golang/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")

	// User´s Endpoint
	v1.POST("/user/signup", handlers.UserSignUp)
	v1.POST("/user/login", handlers.UserLogin)
	v1.PUT("/user/:id", handlers.UserUpdate, middlewares.JwtMiddleware)
	v1.GET("/user/:id", handlers.GetUser, middlewares.JwtMiddleware)

	v1.POST("/product", handlers.ProductPost, middlewares.JwtMiddleware)
}
