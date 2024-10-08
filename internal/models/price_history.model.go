package models

import "time"

type PriceHistory struct {
	ID          uint `gorm:"primaryKey"`
	VersionID   uint
	Price       float64  `gorm:"type:decimal(10,2)"`
	ResalePrice *float64 `gorm:"type:decimal(10,2)"`
	Date        time.Time
}
