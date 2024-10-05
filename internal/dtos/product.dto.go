package dtos

type ProductModel interface {
	GetID() uint
	GetName() string
	GetStatus() bool
	// GetVersions() []VersionModel
}

type ProductCreateDTO struct {
	Name     string             `json:"name" validate:"required,objectname,min=3,max=50"`
	Status   bool               `json:"status"`
	Versions []VersionCreateDTO `json:"versions" validate:"required,min=1,max=6,dive"`
}

type ProductResponseDTO struct {
	ID       uint                 `json:"id"`
	Name     string               `json:"name"`
	Status   bool                 `json:"status"`
	Versions []VersionResponseDTO `json:"versions"`
}

func (p *ProductResponseDTO) Init(product ProductModel) {
	p.ID = product.GetID()
	p.Name = product.GetName()
	p.Status = product.GetStatus()
	// p.Versions = product.GetVersions()
}
