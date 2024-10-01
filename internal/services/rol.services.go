package services

import (
	"errors"

	"github.com/KevinStockmanns/api_golang/internal/db"
	"github.com/KevinStockmanns/api_golang/internal/models"
	"gorm.io/gorm"
)

func GetOrCreateRol(rol *models.Rol, rolName string) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("name = ?", rolName).First(rol).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			rol.Name = rolName
			if err := tx.Create(rol).Error; err != nil {
				tx.Rollback()
				return errors.New("ocurrió un error al insertar el rol")
			}
		} else {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
