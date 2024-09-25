package main

import (
	"github.com/KevinStockmanns/api_golang/db"
	"github.com/KevinStockmanns/api_golang/models"
	"github.com/KevinStockmanns/api_golang/routes"
	"github.com/KevinStockmanns/api_golang/utils"
	"github.com/labstack/echo/v4"
)

func main() {
	utils.InitValidations()
	db.Connection()
	db.DB.AutoMigrate(models.Product{}, models.Version{})
	e := echo.New()

	routes.Index(e)
	routes.ProductRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))

}
