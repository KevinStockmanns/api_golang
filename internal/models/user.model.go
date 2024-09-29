package models

import "time"

type User struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"type:varchar(50)"`
	LastName string    `json:"lastname" gorm:"type:varchar(50)"`
	Email    string    `json:"email" gorm:"type:varchar(255);unique"`
	Password string    `json:"password" gorm:"type:varchar(255)"`
	Birthday time.Time `json:"birthday" gorm:"type:date"`
	Status   bool      `json:"status"`
}
