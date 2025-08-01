package utils

type Pagination struct {
	Page   int
	Limit  int
	Offset int
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
