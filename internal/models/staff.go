package models

import (
	"time"
)

var (
	RoleAdmin     = "admin"
	RoleLibrarian = "librarian"
)

type Staff struct {
	ID           uint
	FullName     string
	Username     string
	PhoneNumber  string
	Email        string
	Role         string
	PasswordHash string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLogin    *time.Time
	version      uint
}

type StaffFilter struct {
	FullName string
	Username string
	Role     string
	Status   string
}

func (Staff) TableName() string {
	return "staff"
}
