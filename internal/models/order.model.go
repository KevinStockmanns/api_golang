package models

import "time"

type Order struct {
	ID          uint   `gorm:"primaryKey"`
	Status      string `gorm:"type:varchar(30)"`
	Date        time.Time
	UpdatedDate *time.Time
	UserID      *uint
	User        *User  `gorm:"foreignKey:UserID"`
	FullName    string `gorm:"type:varchar(70)"`
	Email       string `gorm:"type:varchar(255)"`
	Phone       string `gorm:"type:varchar(30)"`
	OrderItems  []OrderItem
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint
	Price     float64 `gorm:"type:decimal(10,2)"`
	Discount  float32 `gorm:"type:decimal(5,2)"`
	VersionID uint
	Version   Version `gorm:"foreignKey:VersionID"`
}
