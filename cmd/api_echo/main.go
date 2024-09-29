package main

import (
	"log"

	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/labstack/echo/v4"
)

func main() {
	log.Println("Listening...")
	e := echo.New()

	db.ConnectDB()
	db.InitMigrations()

	if err := e.Start(":8080"); err != nil {
		log.Println("ocurrio un error al levantar el servidor")
		log.Println(err)
	}
}
