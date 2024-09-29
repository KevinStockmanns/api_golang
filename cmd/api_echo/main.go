package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	log.Println("Listening...")
	e := echo.New()

	if err := e.Start(":8080"); err != nil {
		log.Println("ocurrio un error al levantar el servidor")
		log.Println(err)
	}
}
