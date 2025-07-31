package models

import "time"

var (
	StatusReturned = "returned"
	StatusBorrowed = "borrowed"
)

type Loan struct {
	ID         int
	BookId     int
	MemberId   int
	BorrowDate time.Time
	DueDate    time.Time
	ReturnDate *time.Time
	Status     string `gorm:"default:borrowed"`
	Book       Book   `gorm:"foreignKey:BookId"`
	Member     Member `gorm:"foreignKey:MemberId"`
}
