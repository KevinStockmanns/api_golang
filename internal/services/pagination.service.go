package services

import "gorm.io/gorm"

// Pagination es una estructura genérica que permite la paginación de cualquier tipo T
type Pagination[T any] struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
	Content    []T   `json:"content"`
}

func NewPagination[T any](page, size, limit int) *Pagination[T] {
	if size > limit || size <= 0 {
		size = limit
	}
	return &Pagination[T]{
		Page: page,
		Size: size,
	}
}

// RunQuery es un método que ejecuta la consulta y llena los datos en la estructura Pagination
func (p *Pagination[T]) RunQuery(db *gorm.DB, condition string, values []interface{}, order string) error {
	var total int64
	offset := (p.Page - 1) * p.Size

	query := db.Model(&p.Content)

	// Aplicamos la condición si existe
	if condition != "" {
		query = query.Where(condition, values...)
	}

	// Contamos el total de registros
	if err := query.Count(&total).Error; err != nil {
		return err
	}

	// Aplicamos el orden si existe
	if order != "" {
		query = query.Order(order)
	}

	// Aplicamos la paginación
	if err := query.Offset(offset).Limit(p.Size).Find(&p.Content).Error; err != nil {
		return err
	}

	// Asignamos los resultados a la estructura Pagination
	p.Total = total

	return nil
}
