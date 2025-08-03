package models

import (
	"time"
)

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
	Loans       []Loan
}
