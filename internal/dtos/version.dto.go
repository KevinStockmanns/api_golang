package dtos

type VersionCreateDTO struct {
	Name        string   `json:"name" validate:"required,objectname,min=3,max=50"`
	Price       float64  `json:"price" validate:"required,gt=0"`
	ResalePrice *float64 `json:"resalePrice" validate:"omitempty,gte=0"`
	Status      bool     `json:"status"`
	Stock       uint     `json:"stock"`
}
