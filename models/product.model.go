package models

import (
	"strings"
	"time"

	"github.com/KevinStockmanns/api_golang/models/wrapper"
	"github.com/KevinStockmanns/api_golang/utils"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Product struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"type:varchar(50);unique"`
	Status   bool      `json:"status"`
	Versions []Version `json:"versions"`
}

func (p *Product) Init(productPost PostProduct) {
	p.Name = strings.Trim(productPost.Name, " ")
	p.Status = productPost.Status
	for _, v := range productPost.Versions {
		p.Versions = append(p.Versions, Version{
			Name:        strings.Trim(v.Name, " "),
			Price:       v.Price,
			ResalePrice: v.ResalePrice,
			Status:      v.Status,
			Date:        time.Now().UTC(),
			Stock:       v.Stock,
			Vistas:      0,
		})
	}
}

func (p *Product) Update(productDto PutProduct, tx *gorm.DB) {
	if productDto.Name != nil {
		p.Name = strings.Trim(*productDto.Name, " ")
	}
	if productDto.Status != nil {
		p.Status = *productDto.Status
	}
	if productDto.Versions != nil {
		var historyPrices []PriceHistory
		for _, vDto := range *productDto.Versions {
			switch strings.ToLower(vDto.Action) {
			case "update":
				for i := range p.Versions {
					if p.Versions[i].ID == *vDto.ID {
						if vDto.Price != nil || vDto.ResalePrice != nil {
							history := PriceHistory{}
							history.Init(p.Versions[i])
							historyPrices = append(historyPrices, history)
						}
						if vDto.Name != nil {
							p.Versions[i].Name = strings.Trim(*vDto.Name, " ")
						}
						if vDto.Price != nil {
							p.Versions[i].Price = *vDto.Price
						}
						if vDto.ResalePrice != nil {
							p.Versions[i].ResalePrice = vDto.ResalePrice
						}
						if vDto.Status != nil {
							p.Versions[i].Status = *vDto.Status
						}
						if vDto.Stock != nil {
							p.Versions[i].Stock = *vDto.Stock
						}
					}
				}
			case "create":
				version := Version{
					Name:        strings.Trim(*vDto.Name, " "),
					Price:       *vDto.Price,
					ResalePrice: vDto.ResalePrice,
					Status:      true,
					ProductID:   p.ID,
					Date:        time.Now().UTC(),
					Stock:       0,
					Vistas:      0,
				}
				if vDto.Stock != nil {
					version.Stock = *vDto.Stock
				}
				if vDto.Status != nil {
					version.Status = *vDto.Status
				}
				tx.Create(&version)
				p.Versions = append(p.Versions, version)
			case "delete":
				for i := range p.Versions {
					p.Versions[i].Status = false
				}
			}
		}

		if len(historyPrices) > 0 {
			tx.Create(&historyPrices)
		}
		tx.Save(&p.Versions)
	}
}

type PostProduct struct {
	Name     string        `json:"name" validate:"required,min=3,max=50,regexp=^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+( [a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+)*$"`
	Status   bool          `json:"status"`
	Versions []VersionPost `json:"versions" validate:"required,min=1,max=6,dive"`
}

func (p *PostProduct) NormalizeAndValidate() []wrapper.ErrorWrapper {
	p.Name = strings.Trim(p.Name, " ")

	for _, v := range p.Versions {
		v.Normalize()
	}

	if err := utils.Validate.Struct(p); err != nil {
		errors := utils.ValidateErrors(err.(validator.ValidationErrors))
		return errors
	}

	return make([]wrapper.ErrorWrapper, 0)
}

type PutProduct struct {
	Name     *string       `json:"name" validate:"omitempty,min=3,max=50,regexp=^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+( [a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+)*$"`
	Status   *bool         `json:"status"`
	Versions *[]PutVersion `json:"versions" validate:"omitempty,min=1,dive"`
}

func (p *PutProduct) NormalizeAndValidate() []wrapper.ErrorWrapper {
	if p.Name != nil {
		*p.Name = strings.Trim(*p.Name, " ")
	}
	if p.Versions != nil {
		for i := range *p.Versions {
			(*p.Versions)[i].Normalize()
		}
	}

	if err := utils.Validate.Struct(p); err != nil {
		errors := utils.ValidateErrors(err.(validator.ValidationErrors))
		return errors
	}

	return make([]wrapper.ErrorWrapper, 0)
}
