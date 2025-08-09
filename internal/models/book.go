package models

import (
	"time"
)

var BookSafeSortList = []string{
	"id",
	"title",
	"isbn",
	"publisher",
	"language",
	"total_copies",
	"available_copies",
	"author_id",
}

type Book struct {
	ID              uint
	Title           string
	ISBN            string
	Publisher       string
	Language        string
	Summary         string
	Translator      string
	CoverImageUrl   string
	AuthorId        uint
	TotalCopies     int
	AvailableCopies int
	CreatedAt       time.Time `gorm:"default:current_timestamp"`
	UpdatedAt       time.Time `gorm:"default:current_timestamp"`
	version         uint
	Loans           []*Loan
	Categories      []*Category `gorm:"many2many:category_books"`
	Author          *Author
}

type BookDashboard struct {
	TotalBooks           int64
	TotalCopies          int64
	TotalAvailableCopies int64
	RecentBooks          []Book
	PopularBooks         []Book
	LowStockBooks        []Book
}

func BookValidateSearchData(data *SearchData) bool {
	if data.SortBy == "" {
		return true
	}
	sortByMap := map[string]bool{}
	return sortByMap[data.SortBy]
}
