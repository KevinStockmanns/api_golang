package main

import (
	"github.com/KevinStockmanns/api_golang/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	routes.Index(e)

	e.Logger.Fatal(e.Start(":8080"))

}
