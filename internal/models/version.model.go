package models

import "time"

type Version struct {
	ID          uint `gorm:"primaryKey"`
	ProductId   uint
	Name        string `gorm:"varchar(50)"`
	Price       float64
	ResalePrice *float64
	Status      bool
	Date        time.Time
	Stock       uint
	Views       uint
}
