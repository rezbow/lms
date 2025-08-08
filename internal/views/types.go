package views

import (
	"fmt"
	"lms/internal/utils"
)

type Data map[string]any
type Errors map[string]string

type SearchData struct {
	Term       string
	BaseUrl    string
	Sort       string
	Direction  string
	Pagination *utils.Pagination
}

func (s SearchData) ToUrl() string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page,
		s.Pagination.Limit,
		s.Sort,
		s.Direction,
	)
}

func (s SearchData) Prev() string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page-1,
		s.Pagination.Limit,
		s.Sort,
		s.Direction,
	)
}

func (s SearchData) SortBy(field, dir string) string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page,
		s.Pagination.Limit,
		field,
		dir,
	)
}

func (s SearchData) Next() string {
	return fmt.Sprintf(
		"%s?q=%s&page=%d&limit=%d&sortBy=%s&dir=%s",
		s.BaseUrl,
		s.Term,
		s.Pagination.Page+1,
		s.Pagination.Limit,
		s.Sort,
		s.Direction,
	)

}

func (d Data) GetInt(key string) (int, bool) {
	i, ok := d[key].(int)
	if !ok {
		return 0, ok
	}
	return i, ok
}

func (d Data) GetErrors() (Errors, bool) {
	errors, ok := d["errors"].(Errors)
	if !ok {
		return nil, ok
	}
	return errors, ok
}
