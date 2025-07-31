package models

import "time"

type Member struct {
	ID       int
	Name     string
	Email    string
	Phone    string
	JoinDate time.Time `gorm:"default:current_timestamp"`
	Status   string    `gorm:"default:active"`
}
