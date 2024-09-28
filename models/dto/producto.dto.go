package dto

import "github.com/KevinStockmanns/api_golang/models"

type ProductResponse struct {
	ID       uint              `json:"id"`
	Name     string            `json:"name"`
	Status   bool              `json:"status"`
	Versions []VersionResponse `json:"versions"`
}

func InitProduct(product models.Product) *ProductResponse {
	var p ProductResponse
	p.ID = product.ID
	p.Name = product.Name
	p.Status = product.Status
	for i := range product.Versions {
		var v VersionResponse
		v.Init(product.Versions[i])
		p.Versions = append(p.Versions, v)
	}
	return &p
}
