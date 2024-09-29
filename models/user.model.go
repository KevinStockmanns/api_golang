package models

import "time"

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(100)"`
	LastName    string    `json:"lastName" gorm:"type:varchar(50)"`
	Birthday    time.Time `json:"birthday" gorm:"type:date"`
	Password    string    `json:"password"`
	Email       string    `json:"email" gorm:"type:varchar(255)"`
	NumberPhone string    `json:"numberPhone" gorm:"type:varchar(30)"`
}
