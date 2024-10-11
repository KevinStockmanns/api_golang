package services

import (
	"time"

	"github.com/KevinStockmanns/api_golang/internal/models"
	"gorm.io/gorm"
)

func RegisterPrice(version models.Version, db *gorm.DB) *models.PriceHistory {
	historyItem := models.PriceHistory{
		VersionID:   version.ID,
		Price:       version.Price,
		ResalePrice: version.ResalePrice,
		Date:        time.Now().UTC(),
	}

	if db.Create(&historyItem).Error != nil {
		db.Rollback()
		return nil
	}

	return &historyItem
}

func GetHistory(db *gorm.DB, history *[]models.PriceHistory, initTime time.Time, endTime time.Time, id string) error {
	endTime = endTime.AddDate(0, 0, 1)
	return db.Model(history).Where("date >= ? && date < ? && version_id = ?", initTime, endTime, id).Find(&history).Error
}
