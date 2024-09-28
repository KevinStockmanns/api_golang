package dto

import (
	"time"

	"github.com/KevinStockmanns/api_golang/models"
)

type VersionResponse struct {
	ID     uint      `json:"id"`
	Name   string    `json:"name"`
	Price  float64   `json:"price"`
	Status bool      `json:"status"`
	Date   time.Time `json:"date"`
}

func (v *VersionResponse) Init(version models.Version) {
	v.ID = version.ID
	v.Name = version.Name
	v.Date = version.Date
	v.Price = version.Price
	v.Status = version.Status
}
