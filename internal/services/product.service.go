package services

import (
	"fmt"
	"time"

	"github.com/KevinStockmanns/api_golang/internal/constants"
	"github.com/KevinStockmanns/api_golang/internal/dtos"
	"github.com/KevinStockmanns/api_golang/internal/models"
	"gorm.io/gorm"
)

func ProductUpdate(product *models.Product, productDto dtos.ProductUpdateDTO, db *gorm.DB) {
	if productDto.Name != nil {
		product.Name = *productDto.Name
	}
	if productDto.Status != nil {
		product.Status = *productDto.Status
	}
	if productDto.Versions != nil {
		for _, vDto := range *productDto.Versions {
			if vDto.Action == string(constants.Create) {
				version := models.Version{
					ProductId:   product.ID,
					Name:        *vDto.Name,
					Price:       *vDto.Price,
					ResalePrice: vDto.ResalePrice,
					Status:      *vDto.Status,
					Date:        time.Now().UTC(),
					Stock:       *vDto.Stock,
					Views:       0,
				}
				version.Normalize()
				if err := db.Create(&version).Error; err != nil {
					db.Rollback()
				}
				product.Versions = append(product.Versions, version)
			} else {
				for i := range product.Versions {
					if *vDto.ID == product.Versions[i].ID {
						if vDto.Action == string(constants.Update) {
							if vDto.Price != nil || vDto.ResalePrice != nil {
								RegisterPrice(product.Versions[i], db)
							}
							if vDto.Name != nil {
								product.Versions[i].Name = *vDto.Name
							}
							if vDto.Price != nil {
								product.Versions[i].Price = *vDto.Price
							}
							if vDto.ResalePrice != nil {
								product.Versions[i].ResalePrice = vDto.ResalePrice
							}
							if vDto.Status != nil {
								product.Versions[i].Status = *vDto.Status
							}
							if vDto.Stock != nil {
								product.Versions[i].Stock = *vDto.Stock
							}
							product.Versions[i].Normalize()
						}
						if vDto.Action == string(constants.Delete) {
							product.Versions[i].Status = false
						}

						db.Save(product.Versions[i])
					}
				}
			}
		}
	}
	product.Normalize()
	if len(product.Versions) > 6 {
		db.Rollback()
	} else {
		if err := db.Save(&product).Error; err != nil {
			db.Rollback()
		}
	}
}

func ChangeProductsPrice(data dtos.ProductChangePrice, db *gorm.DB) error {
	var ids []uint

	for _, vDto := range data.Versions {
		ids = append(ids, vDto.ID)
	}

	var versions []models.Version
	if db.Model(versions).Where("id IN ?", ids).Find(&versions).Error != nil {
		return fmt.Errorf("ocurrio un error al obtener las versiones")
	}

	for _, vDto := range data.Versions {
		for _, v := range versions {
			if vDto.ID == v.ID {
				RegisterPrice(v, db)
				if vDto.Price != nil {
					v.Price = *vDto.Price
				}
				if vDto.ResalePrice != nil {
					v.ResalePrice = vDto.ResalePrice
				}
				db.Model(v).Save(&v)
			}
		}
	}

	return nil
}
