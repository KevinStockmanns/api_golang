package middlewares

import (
	"net/http"
	"strings"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/encryptor"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Message: "el token de seguridad es requerido"})
			}

			token := authHeader[7:]
			claims, err := encryptor.VerifyJWT(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dtos.ErrorResponse{Message: err.Error()})
			}

			c.Set("tokenClaims", claims)

			if len(allowedRoles) > 0 {
				rolToken := claims.Rol

				for _, userRole := range allowedRoles {
					if rolToken == userRole {
						return next(c)
					}
				}

				return c.JSON(http.StatusForbidden, dtos.ErrorResponse{Message: "No tienes los permisos necesarios"})
			}

			return next(c)
		}
	}
}
