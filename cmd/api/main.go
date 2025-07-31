package main

import (
	"lms/internal"
	"lms/internal/database"
	"lms/internal/handlers"
	"lms/internal/repositories"

	"github.com/go-playground/validator/v10"
)

func main() {
	dsn := "host=localhost user=admin password=cfaa7e52 dbname=lms_db port=5432 sslmode=disable TimeZone=UTC"
	validate := validator.New()
	db := database.SetupDataBase(dsn)

	bookRepo := repositories.BookRepo{DB: db}
	bookHandler := handlers.BookHandler{
		BookRepo:  &bookRepo,
		Validator: validate,
	}

	authorRepo := repositories.AuthorRepo{DB: db}
	authorHandler := handlers.AuthorHandler{Repo: authorRepo, Validator: validate}

	memberRepo := repositories.MemberRepo{DB: db}
	memberHandler := handlers.MemberHandler{Repo: &memberRepo, Validator: validate}

	loanRepo := repositories.LoanRepo{DB: db}
	loanHandler := handlers.LoanHandler{Repo: &loanRepo, Validator: validate}

	r := internal.SetupRouter(&bookHandler, &authorHandler, &memberHandler, &loanHandler)
	r.Run(":8080")
}
