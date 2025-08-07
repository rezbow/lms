package models

import "time"

type PopularBook struct {
	ID        uint
	Title     string
	LoanCount uint
}

type ActiveMember struct {
	ID        uint
	FullName  string
	LoanCount uint
}

type PopularCategory struct {
	Slug      string
	Name      string
	LoanCount uint
}

type UpcomingLoan struct {
	ID         uint
	BookTitle  string
	MemberName string
	DueDate    time.Time
}

type Dashboard struct {
	TotalBooks        int64
	TotalMembers      int64
	TotalActiveLoans  int64
	TotalOverdueLoans int64
	PopularBooks      []PopularBook
	ActiveMembers     []ActiveMember
	PopularCategory   []PopularCategory
	UpcomingLoans     []UpcomingLoan
	RecentActivities  []ActivityLog
}
