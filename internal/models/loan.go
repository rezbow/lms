package models

import "time"

var (
	StatusReturned = "returned"
	StatusBorrowed = "borrowed"
)

type Loan struct {
	ID         int
	BookId     uint
	MemberId   uint
	BorrowDate time.Time
	DueDate    time.Time
	ReturnDate *time.Time
	Status     string `gorm:"default:borrowed"`
}
