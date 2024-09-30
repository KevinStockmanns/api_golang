package main

import (
	"log"

	v1 "github.com/KevinStockmanns/api_golang/api/v1"
	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/validators"
	"github.com/labstack/echo/v4"
)

func main() {
	db.ConnectDB()
	log.Println("DB Connected")
	db.InitMigrations()
	log.Println("Listening...")

	e := echo.New()
	v1.RegisterRoutes(e)
	validators.InitValidations()
	if err := e.Start(":8080"); err != nil {
		log.Println("ocurrio un error al levantar el servidor")
		log.Println(err)
	}
}
