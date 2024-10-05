package models

import (
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
)

type User struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"type:varchar(50)"`
	LastName string    `json:"lastname" gorm:"type:varchar(50)"`
	Email    string    `json:"email" gorm:"type:varchar(255);unique"`
	Password string    `json:"password" gorm:"type:varchar(255)"`
	Birthday time.Time `json:"birthday" gorm:"type:date"`
	Status   bool      `json:"status"`
	Phone    string    `json:"phone" gorm:"type:varchar(30)"`
	RolId    uint      `json:"rolId"`
	Rol      Rol       `json:"rol" gorm:"foreignKey:RolId"`
}

func (u *User) Normalize() {
	u.Name = strings.ToTitle(strings.TrimSpace(u.Name))
	u.LastName = strings.ToTitle(strings.TrimSpace(u.LastName))
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)
	u.Phone = strings.TrimSpace(u.Phone)
}

func (u *User) IsAdmin() bool {
	return u.Rol.Name == string(constants.Admin) || u.Rol.Name == string(constants.SuperAdmin)
}

func (u *User) Update(data dtos.UserUpdateDTO) {
	if data.Birthday != nil {
		bDay, _ := time.Parse("2006-01-02", *data.Birthday)
		u.Birthday = bDay
	}
	if data.Email != nil {
		u.Email = *data.Email
	}
	if data.LastName != nil {
		u.LastName = strings.ToTitle(strings.TrimSpace(*data.LastName))
	}
	if data.Name != nil {
		u.Name = strings.ToTitle(strings.TrimSpace(*data.Name))
	}
	if data.Phone != nil {
		u.Phone = *data.Phone
	}
	if data.Status != nil {
		u.Status = *data.Status
	}
}

func (u User) GetID() uint            { return u.ID }
func (u User) GetName() string        { return u.Name }
func (u User) GetLastName() string    { return u.LastName }
func (u User) GetEmail() string       { return u.Email }
func (u User) GetBirthday() time.Time { return u.Birthday }
func (u User) GetPhone() string       { return u.Phone }
func (u User) GetRol() string         { return u.Rol.Name }
func (u User) IsStatus() bool         { return u.Status }
