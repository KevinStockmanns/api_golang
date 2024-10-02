package middlewares

import (
	"net/http"
	"strings"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/encryptor"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		// log.Println(authHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Message: "el token de seguridad es requerido"})
		}

		token := authHeader[7:]
		if _, err := encryptor.VerifyJWT(token); err != nil {
			return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Message: err.Error()})
		}
		return next(c)
	}
}
