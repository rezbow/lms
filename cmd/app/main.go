package main

import (
	"encoding/gob"
	"lms/internal"
	"lms/internal/database"
	"lms/internal/handlers"
	"lms/internal/models"
	"lms/internal/repositories"

	"github.com/go-playground/validator/v10"
)

func main() {

	dsn := "host=localhost user=admin password=cfaa7e52 dbname=lms_db port=5432 sslmode=disable TimeZone=UTC"
	validate := validator.New()
	db := database.SetupDataBase(dsn)

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
