package models

import (
	"time"
)

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
}

type BookFilter struct {
	Title      string
	ISBN       string
	Publisher  string
	Language   string
	Translator string
	AuthorId   uint
}
