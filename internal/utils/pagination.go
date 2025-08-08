package utils

import "math"

type Pagination struct {
	Page      int
	Limit     int
	Total     int64
	TotalPage int
	Offset    int
	BaseUrl   string
}

func NewPagination(page, limit int, baseUrl string) *Pagination {
	if page <= 0 || limit <= 0 {
		return nil
	}
	return &Pagination{
		Page:    page,
		Limit:   limit,
		Offset:  (page - 1) * limit,
		BaseUrl: baseUrl,
	}
}

func (p *Pagination) CalculateTotalPage() {
	p.TotalPage = int(math.Ceil(float64(p.Total) / float64(p.Limit)))
}
