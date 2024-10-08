package models

import (
	"strings"

	"github.com/KevinStockmanns/api_golang/internal/dtos"
)

type Product struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(50);unique"`
	Status   bool
	Versions []Version
}

func (p *Product) Create(data dtos.ProductCreateDTO) {
	p.Name = data.Name
	p.Status = data.Status
	oneVersionActive := false
	for _, vDto := range data.Versions {
		var version Version
		version.Create(vDto)
		p.Versions = append(p.Versions, version)
		if data.Status {
			oneVersionActive = true
		}
	}
	if p.Status && !oneVersionActive {
		p.Status = false
	}

	p.Normalize()
}

func (p *Product) Normalize() {
	p.Name = strings.ToTitle(strings.TrimSpace(p.Name))
}

func (p Product) GetID() uint     { return p.ID }
func (p Product) GetName() string { return p.Name }
func (p Product) GetStatus() bool { return p.Status }
func (p Product) GetVersions() []dtos.VersionModel {
	versions := make([]dtos.VersionModel, len(p.Versions))
	for i, version := range p.Versions {
		versions[i] = version
	}
	return versions
}
