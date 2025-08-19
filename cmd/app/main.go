package main

import (
	"encoding/gob"
	"fmt"
	"lms/internal"
	"lms/internal/database"
	"lms/internal/handlers"
	"lms/internal/models"
	"lms/internal/repositories"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func getDsn() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		log.Fatal(".env: missing db_user")
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		log.Fatal(".env: missing db_password")
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		log.Fatal(".env: missing db_name")
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, name, port,
	)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	validate := validator.New()
	db := database.SetupDataBase(getDsn())

	gob.Register(models.Staff{})

	logRepo := repositories.ActivityRepo{DB: db}

	categoryRepo := repositories.CategoryRepo{DB: db}
	categoryHandler := handlers.CategoryHandler{Repo: &categoryRepo, Validator: validate}

	bookRepo := repositories.BookRepo{DB: db}
	bookHandler := handlers.BookHandler{
		BookRepo:     &bookRepo,
		LogRepo:      &logRepo,
		CategoryRepo: &categoryRepo,
		Validator:    validate,
	}

	authorRepo := repositories.AuthorRepo{DB: db}
	authorHandler := handlers.AuthorHandler{Repo: authorRepo, Validator: validate}

	memberRepo := repositories.MemberRepo{DB: db}
	memberHandler := handlers.MemberHandler{Repo: &memberRepo, LogRepo: &logRepo, Validator: validate}

	loanRepo := repositories.LoanRepo{DB: db}
	loanHandler := handlers.LoanHandler{Repo: &loanRepo, LogRepo: &logRepo, Validator: validate}

	staffRepo := repositories.StaffRepo{DB: db}
	staffHandler := handlers.StaffHandler{Repo: &staffRepo, Validator: validate}

	dashboardRepo := repositories.DashboardRepo{DB: db}
	dashboardHandler := handlers.DashboardHanlder{Repo: &dashboardRepo, LogRepo: &logRepo}

	r := internal.SetupRouter(&bookHandler, &authorHandler, &memberHandler, &loanHandler, &categoryHandler, &staffHandler, &dashboardHandler)
	r.Run(":8080")
}
