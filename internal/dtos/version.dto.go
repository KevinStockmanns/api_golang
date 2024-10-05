package dtos

import "time"

type VersionModel interface {
	GetID() uint
	GetName() string
	GetPrice() float64
	GetResalePrice() float64
	GetStatus() bool
	GetDate() time.Time
	GetStock() uint
	GetViews() uint
}

type VersionCreateDTO struct {
	Name        string   `json:"name" validate:"required,objectname,min=3,max=50"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	ResalePrice *float64 `json:"resalePrice" validate:"omitempty,gte=0"`
	Status      bool     `json:"status"`
	Stock       uint     `json:"stock"`
}

type VersionResponseDTO struct {
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	ResalePrice *float64  `json:"resalePrice,omitempty"`
	Status      bool      `json:"status"`
	Date        time.Time `json:"date"`
	Stock       uint      `json:"stock"`
	Views       uint      `json:"views,omitempty"`
}
