package db

import (
	"log"

	"github.com/KevinStockmanns/api_golang/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:root@tcp(127.0.0.1:3307)/api_go?charset=utf8mb4&parseTime=True&loc=UTC"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error al conectar a la base de datos")
		panic(err)
	}
}
func InitMigrations() {
	DB.AutoMigrate(models.User{}, models.Rol{})
}
