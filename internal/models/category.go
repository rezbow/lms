package models

import "time"

type Category struct {
	ID        uint
	Slug      string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	version   uint
	Books     []*Book `gorm:"many2many:category_books"`
}
