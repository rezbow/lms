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

type AuthorFormData struct {
	FullName    string `form:"fullName" binding:"required" validator:"required,min=1"`
	Nationality string `form:"nationality" binding:"required" validator:"required,min=1,max=50"`
	Bio         string `form:"bio" binding:"omitempty" validator:"omitempty,min=1,max=500"`
}

type AuthorEditFormData struct {
	FullName    *string `form:"fullName" binding:"omitempty" validator:"omitempty,min=1"`
	Nationality *string `form:"nationality" binding:"omitempty" validator:"omitempty,min=1,max=50"`
	Bio         *string `form:"bio" binding:"omitempty" validator:"omitempty,min=1,max=500"`
}

type AuthorDashboard struct {
	TotalAuthors   int64
	RecentAuthors  []Author
	PopularAuthors []Author
}
