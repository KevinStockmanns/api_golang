package models

import "time"

type PriceHistory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	VersionID   uint      `json:"versionId"`
	Date        time.Time `json:"date"`
	Price       float64   `json:"price"`
	ResalePrcie *float64  `json:"resalePrice"`
}

func (h *PriceHistory) Init(version Version) {
	h.VersionID = version.ID
	h.Date = time.Now().UTC()
	h.Price = version.Price
	h.ResalePrcie = version.ResalePrice
}
