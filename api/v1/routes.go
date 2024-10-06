package v1

import (
	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/handlers"
	"github.com/KevinStockmanns/api_golang/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")

	// User´s Endpoints
	v1.POST("/user/signup", handlers.UserSignUp)
	v1.POST("/user/login", handlers.UserLogin)
	v1.PUT("/user/:id", handlers.UserUpdate, middlewares.JwtMiddleware())
	v1.GET("/user/:id", handlers.GetUser, middlewares.JwtMiddleware())
	v1.GET("/user", handlers.UserList, middlewares.JwtMiddleware(string(constants.SuperAdmin), string(constants.Admin)))
	v1.PUT("/user/password", handlers.UserChangePassword, middlewares.JwtMiddleware())

	//Product´s Endpoints
	v1.POST("/product", handlers.ProductPostHandler, middlewares.JwtMiddleware(string(constants.SuperAdmin), string(constants.Admin)))
	v1.GET("/product/:id", handlers.GetUser)
	v1.PATCH("product/views", handlers.ProductViews, middlewares.JwtMiddleware(string(constants.Admin), string(constants.SuperAdmin)))
}
