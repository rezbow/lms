package models

import (
	"time"
)

var (
	StatusActive    = "active"
	StatusSuspended = "suspended"
)

var MemberSafeSortList = []string{
	"id",
	"full_name",
	"email",
	"phone_number",
	"national_id",
	"joined_at",
	"status",
}

type Member struct {
	ID          uint
	FullName    string
	Email       string
	PhoneNumber string
	NationalId  string
	JoinedAt    time.Time `gorm:"default:current_timestamp"`
	Status      string    `gorm:"default:active"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	version     uint
	Loans       []*Loan
}

type MemberFilter struct {
	FullName    string
	Email       string
	PhoneNumber string
	NationalId  string
	JoinedAt    time.Time
}
