package v1

import (
	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/handlers"
	"github.com/KevinStockmanns/api_golang/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	v1 := e.Group("/v1")

	admin := string(constants.Admin)
	superAdmin := string(constants.SuperAdmin)
	admins := []string{admin, superAdmin}

	// User´s Endpoints
	v1.POST("/user/signup", handlers.UserSignUp)
	v1.POST("/user/login", handlers.UserLogin)
	v1.PUT("/user/:id", handlers.UserUpdate, middlewares.JwtMiddleware(true))
	v1.GET("/user/:id", handlers.GetUser, middlewares.JwtMiddleware(true))
	v1.GET("/user", handlers.UserList, middlewares.JwtMiddleware(true, admins...))
	v1.PUT("/user/password", handlers.UserChangePassword, middlewares.JwtMiddleware(true))

	//Product´s Endpoints
	v1.POST("/product", handlers.ProductCreate, middlewares.JwtMiddleware(true, admins...))
	v1.GET("/product/:id", handlers.ProductGet)
	v1.GET("/product", handlers.ProductList, middlewares.JwtMiddleware(false))
	v1.PATCH("product/views", handlers.ProductUpViews, middlewares.JwtMiddleware(true, admins...))
	v1.PUT("/product/:id", handlers.ProductUpdate, middlewares.JwtMiddleware(true, admins...))
	v1.DELETE("/product/:id", handlers.ProductDelete, middlewares.JwtMiddleware(true, admins...))
	v1.GET("/product/version/:id/price", handlers.ProductPriceHistoyList, middlewares.JwtMiddleware(true, admins...))
	v1.POST("/product/price", handlers.ProductChangePrice, middlewares.JwtMiddleware(true, admins...))
}
