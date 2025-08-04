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
