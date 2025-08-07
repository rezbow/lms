package models

import (
	"database/sql"
	"time"
)

var (
	// user
	ActivityTypeDeleteMember = "member_deleted"
	ActivityTypeUpdateMember = "member_updated"
	ActivityTypeAddMember    = "member_added"
	// loan
	ActivityTypeAddLoan    = "loan_added"
	ActivityTypeUpdateLoan = "loan_updated"
	ActivityTypeDeleteLoan = "loan_deleted"
	ActivityTypeReturnLoan = "loan_returned"
	// book
	ActivityTypeAddBook    = "book_added"
	ActivityTypeUpdateBook = "book_updated"
	ActivityTypeDeleteBook = "book_deleted"
)

var (
	ActorTypeMember = "member"
	ActorTypeStaff  = "staff"
	ActorTypeSystem = "system"
)

var (
	EntityTypeLoan     = "loan"
	EntityTypeMember   = "member"
	EntityTypeBook     = "book"
	EntityTypeAuthor   = "author"
	EntityTypeCategory = "category"
	EntityTypeStaff    = "staff"
)

type ActivityLog struct {
	ID           uint
	ActivityType string
	ActorId      sql.NullInt32
	ActorType    string
	Description  string
	EntityId     sql.NullInt32
	EntityType   sql.NullString
	CreatedAt    time.Time `gorm:"default:current_timestamp"`
}
