package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() {
	dsn := "root:root@tcp(127.0.0.1:3307)/api_go?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	log.Println("DB Conected")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error al conectar a la base de datos")
		log.Println(err)
	}
}
