package models

import (
	"fmt"
	"math"
)

type SearchData struct {
	Term       string
	Pagination *Pagination
	SortBy     string
	Dir        string
	BaseUrl    string
}

func (s *SearchData) dir() string {
	switch s.Dir {
	case "asc":
		return "desc"
	default:
		return "asc"
	}
}

func (s *SearchData) SortWith(sortBy string) string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page,
		s.Pagination.Limit,
		sortBy,
		s.dir(),
	)
}

func (s *SearchData) URL() string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page,
		s.Pagination.Limit,
		s.SortBy,
		s.Dir,
	)
}

func (s *SearchData) PrevPageUrl() string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page-1,
		s.Pagination.Limit,
		s.SortBy,
		s.Dir,
	)
}

func (s *SearchData) NextPageUrl() string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page+1,
		s.Pagination.Limit,
		s.SortBy,
		s.Dir,
	)
}

func (s *SearchData) Valid(safeSort []string) bool {
	if s.SortBy == "" {
		return true
	}
	for _, safe := range safeSort {
		if s.SortBy == safe {
			return true
		}
	}
	return false
}

type Pagination struct {
	Page      int
	Limit     int
	Total     int64
	TotalPage int
	Offset    int
	BaseUrl   string
}

func NewPagination(page, limit int) *Pagination {
	if page <= 0 || limit <= 0 {
		return nil
	}
	return &Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}

func (p *Pagination) CalculateTotalPage() {
	p.TotalPage = int(math.Ceil(float64(p.Total) / float64(p.Limit)))
}
