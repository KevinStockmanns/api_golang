package utils

import (
	"math"
	"strconv"

	"github.com/KevinStockmanns/api_golang/db"
)

type Pagination struct {
	Page          int
	Size          int
	TotalPages    int64
	TotalElements int64
}

func InitPagination(page string, size string, limit int) *Pagination {
	p := Pagination{}
	p.Page = 1
	p.Size = limit
	if size != "" {
		if siz, err := strconv.Atoi(size); err == nil && siz > 0 && siz <= limit {
			p.Size = siz
		}
	}
	if page != "" {
		if pag, err := strconv.Atoi(page); err == nil {
			if pag > 0 {
				p.Page = pag
			}
		}
	}
	return &p
}

func (p *Pagination) Query(model interface{}, preloads []string, condition string, values ...interface{}) error {

	var total int64
	db.DB.Debug().Model(model).Where(condition, values...).Count(&total)

	p.TotalPages = int64(math.Ceil(float64(total) / float64(p.Size)))
	p.TotalElements = total

	query := db.DB.Model(model).Where(condition, values...)
	for _, relation := range preloads {
		query = query.Preload(relation)
	}

	return query.Limit(p.Size).Offset((p.Page - 1) * p.Size).Find(model).Error
}
