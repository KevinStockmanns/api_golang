package models

import "time"

type Version struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"tpye:varchar(50)"`
	Price       float64   `json:"price"`
	ResalePrice *float64  `json:"resalePrice"`
	Status      bool      `json:"status"`
	ProductID   uint      `json:"productId"`
	Date        time.Time `json:"date"`
	Stock       uint      `json:"stock"`
	Vistas      uint      `json:"vistas"`
}

type VersionPost struct {
	Name        string   `json:"name" validate:"required,min=3,max=50,regexp=^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+( [a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+)*$"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	ResalePrice *float64 `json:"resalePrice" validate:"omitempty,gt=0"`
	Status      bool     `json:"status"`
	Stock       uint     `json:"stock" validate:"min=0"`
}
