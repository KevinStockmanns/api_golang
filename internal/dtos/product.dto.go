package dtos

type ProductModel interface {
	GetID() uint
	GetName() string
	GetStatus() bool
	GetVersions() []VersionModel
}

type ProductCreateDTO struct {
	Name     string             `json:"name" validate:"required,objectname,min=3,max=50"`
	Status   bool               `json:"status"`
	Versions []VersionCreateDTO `json:"versions" validate:"required,min=1,max=6,dive"`
}

type ProductUpdateDTO struct {
	Name     *string             `json:"name" validate:"omitempty,required,objectname,min=3,max=50"`
	Status   *bool               `json:"status"`
	Versions *[]VersionUpdateDTO `json:"versions" validate:"omitempty,required,min=1,dive"`
}

type ProductPriceHistoryDTO struct {
	InitTime string `json:"initTime" validate:"required,date"`
	EndTime  string `json:"endTime" validate:"required,date"`
}

type ProductChangePrice struct {
	Versions []VersionsChangePrice `json:"versions" validate:"required,min=1,dive"`
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
	p.Versions = make([]VersionResponseDTO, len(product.GetVersions()))
	for i, v := range product.GetVersions() {
		p.Versions[i] = VersionResponseDTO{
			Name:        v.GetName(),
			Price:       v.GetPrice(),
			ResalePrice: v.GetResalePrice(),
			Status:      v.GetStatus(),
			Date:        v.GetDate(),
			Stock:       v.GetStock(),
			Views:       v.GetViews(),
		}
	}
}
