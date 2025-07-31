package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDataBase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: false})
	if err != nil {
		log.Fatal("failed to setupt database", err)
	}
	return db
}
