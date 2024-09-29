package middlewares

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type protectedEndpoint struct {
	Endpoint string
	Method   string
	Roles    []string
}

var protectedEndpoints []protectedEndpoint = make([]protectedEndpoint, 0)

func ProtectEndPoint(endpoint string, method string, roles ...string) {
	protectedEndpoints = append(protectedEndpoints, protectedEndpoint{
		Endpoint: endpoint,
		Method:   method,
		Roles:    roles,
	})
}

func JwtMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Path())
		log.Println(c.Request().URL)

		for _, e := range protectedEndpoints {
			if e.Endpoint == c.Path() && e.Method == c.Request().Method {

				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "no estas autorizado",
				})
			}
		}

		return next(c)
	}
}
