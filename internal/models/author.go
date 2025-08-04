package models

import "time"

type Author struct {
	ID          int
	FullName    string
	Nationality string
	Bio         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	version     uint
	Books       []*Book
}
