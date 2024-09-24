package models

type Product struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"type:varchar(50);unique"`
	Status   bool      `json:"status"`
	Versions []Version `json:"versions"`
}
