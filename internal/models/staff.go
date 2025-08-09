package models

import (
	"time"
)

var (
	RoleAdmin     = "admin"
	RoleLibrarian = "librarian"
)

var StaffSafeSortList = []string{
	"id", "full_name", "username",
	"phone_number", "email", "role",
	"last_login", "status",
}

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

type StaffDashboard struct {
	TotalStaff int64
	Admin      Staff
	Librarians []Staff
}

func (Staff) TableName() string {
	return "staff"
}
