package models

type Product struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name" gorm:"type:varchar(50);unique"`
	Status   bool      `json:"status"`
	Versions []Version `json:"versions"`
}

type PostProduct struct {
	Name     string        `json:"name" validate:"required,min=3,max=50,regexp=^[a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+( [a-zA-Z0-9ñÑáéíóúÁÉÍÓÚ]+)*$"`
	Status   bool          `json:"status"`
	Versions []VersionPost `json:"versions" validate:"required,min=1,max=6,dive"`
}
