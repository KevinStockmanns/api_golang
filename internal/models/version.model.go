package models

import (
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/utils"
)

type Version struct {
	ID          uint `gorm:"primaryKey"`
	ProductId   uint
	Name        string   `gorm:"varchar(50)"`
	Price       float64  `gorm:"type:decimal(10,2)"`
	ResalePrice *float64 `gorm:"type:decimal(10,2)"`
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
	v.Stock = data.Stock
	v.Views = 0

	v.Normalize()
}

func (v *Version) Normalize() {
	v.Name = strings.ToTitle(strings.TrimSpace(v.Name))
	v.Price = utils.RoundDecimal(v.Price)
	if v.ResalePrice != nil {
		*v.ResalePrice = utils.RoundDecimal(*v.ResalePrice)
	}
}
