package models

type Version struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name" gorm:"tpye:varchar(50)"`
	Price       float64  `json:"price"`
	ResalePrice *float64 `json:"resalePrice"`
	Status      bool     `json:"status"`
	ProductID   uint     `json:"productId"`
}
