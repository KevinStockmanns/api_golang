package models

type Rol struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"varchar(30),unique,not null"`
}
