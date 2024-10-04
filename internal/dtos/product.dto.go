package dtos

type ProductCreateDTO struct {
	Name     string             `json:"name" validate:"required,objectname,min=3,max=50"`
	Status   bool               `json:"status"`
	Versions []VersionCreateDTO `json:"versions" validate:"required,min=1,max=6"`
}
