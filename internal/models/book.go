package models

import (
	"time"
)

type Book struct {
	ID              int
	TitleFa         string
	TitleEn         string
	ISBN            string
	AuthorId        int
	TotalCopies     int
	AvailableCopies int
	CreatedAt       time.Time `gorm:"default:current_timestamp"`
	Author          Author    `gorm:"foreignKey:AuthorId"`
}
