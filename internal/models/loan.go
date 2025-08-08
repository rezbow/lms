package models

import "time"

var (
	StatusReturned = "returned"
	StatusBorrowed = "borrowed"
)

var LoanSafeSortList = []string{
	"id",
	"borrow_date",
	"return_date",
	"due_date",
	"status",
}

type Loan struct {
	ID         uint
	BookId     uint
	MemberId   uint
	BorrowDate time.Time
	DueDate    time.Time
	ReturnDate *time.Time
	Status     string `gorm:"default:borrowed"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	version    uint
	Book       *Book
	Member     *Member
}

type LoanFilter struct {
	BookId   uint
	MemberId uint
	Status   string
}
