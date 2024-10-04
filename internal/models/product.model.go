package models

type Product struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(50);unique"`
	Status   bool
	Versions []Version
}
