package models

import (
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
)

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

func (v *Version) Create(data dtos.VersionCreateDTO) {
	v.Date = time.Now().UTC()
	v.Name = data.Name
	v.Price = data.Price
	v.ResalePrice = data.ResalePrice
	v.Status = data.Status
	v.Stock = v.Stock
	v.Views = 0

	v.Normalize()
}

func (v *Version) Normalize() {
	v.Name = strings.ToTitle(strings.TrimSpace(v.Name))
}
