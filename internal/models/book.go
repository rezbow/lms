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

type BookFilter struct {
	Title      string
	ISBN       string
	Publisher  string
	Language   string
	Translator string
	AuthorId   uint
}

func BookValidateSearchData(data *SearchData) bool {
	if data.SortBy == "" {
		return true
	}
	sortByMap := map[string]bool{}
	return sortByMap[data.SortBy]
}
