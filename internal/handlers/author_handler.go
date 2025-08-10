package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"
	authorViews "lms/internal/views/authors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthorHandler struct {
	Repo      repositories.AuthorRepo
	Validator *validator.Validate
}

func (ah *AuthorHandler) Index(ctx *gin.Context) {
	totalAuthors, err := ah.Repo.Total()
	recentAuthors, err := ah.Repo.RecentAuthors(5)
	popularAuthors, err := ah.Repo.PopularAuthors(5)
	if err != nil {
		serverError(ctx)
		return
	}
	data := models.AuthorDashboard{
		TotalAuthors:   totalAuthors,
		RecentAuthors:  recentAuthors,
		PopularAuthors: popularAuthors,
	}
	render(ctx, authorViews.AuthorsDashboard(&data), "authors")
}

func (ah *AuthorHandler) Get(ctx *gin.Context) {
	authorId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	author, err := ah.Repo.GetById(authorId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	render(ctx, authorViews.AuthorInfo(author), author.FullName)
}

func (ah *AuthorHandler) Delete(ctx *gin.Context) {
	authorId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	err = ah.Repo.DeleteById(authorId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	redirect(ctx, "/authors")
}

func (ah *AuthorHandler) AddPage(ctx *gin.Context) {
	render(ctx, authorViews.AuthorsForm(nil, views.Errors{}, "/authors/add"), "add a author")
}

func (ah *AuthorHandler) Add(ctx *gin.Context) {
	var userInput models.AuthorFormData
	action := "/authors/add"

	if err := ctx.ShouldBind(&userInput); err != nil {
		errors := parseValidationErrors(err)
		render(ctx, authorViews.AuthorsForm(&userInput, errors, action), "add new author")
		return
	}

	if err := ah.Validator.Struct(&userInput); err != nil {
		errors := parseValidationErrors(err)
		render(ctx, authorViews.AuthorsForm(&userInput, errors, action), "add new author")
		return
	}

	author := models.Author{
		FullName:    userInput.FullName,
		Nationality: userInput.Nationality,
		Bio:         userInput.Bio,
	}

	if err := ah.Repo.Insert(&author); err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
			return
		}
		render(ctx, authorViews.AuthorsForm(&userInput, ah.Repo.ErrorToFromError(err), action), "add new author")
		return
	}

	redirect(ctx, fmt.Sprintf("/authors/%d", author.ID))
}

func (ah *AuthorHandler) EditPage(ctx *gin.Context) {
	id, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	author, err := ah.Repo.GetById(id)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	render(ctx, authorViews.AuthorsForm(&models.AuthorFormData{
		Nationality: author.Nationality,
		Bio:         author.Bio,
		FullName:    author.FullName,
	}, views.Errors{}, fmt.Sprintf("/authors/%d/edit", author.ID)), author.FullName)
}

func (ah *AuthorHandler) Edit(ctx *gin.Context) {
	id, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	author, err := ah.Repo.GetById(id)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	formData := models.AuthorFormData{
		Bio:         author.Bio,
		Nationality: author.Nationality,
		FullName:    author.FullName,
	}

	action := fmt.Sprintf("/authors/%d/edit", author.ID)
	var userInput models.AuthorEditFormData

	if err := ctx.ShouldBind(&userInput); err != nil {
		errors := parseValidationErrors(err)
		render(ctx, authorViews.AuthorsForm(&formData, errors, action), "add new author")
		return
	}

	if err := ah.Validator.Struct(&userInput); err != nil {
		errors := parseValidationErrors(err)
		render(ctx, authorViews.AuthorsForm(&formData, errors, action), "add new author")
		return
	}

	if userInput.FullName != nil {
		author.FullName = *userInput.FullName
	}
	if userInput.Bio != nil {
		author.Bio = *userInput.Bio
	}
	if userInput.Nationality != nil {
		author.Nationality = *userInput.Nationality

	}

	if err := ah.Repo.Update(author); err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
			return
		}
		render(ctx, authorViews.AuthorsForm(&formData, ah.Repo.ErrorToFromError(err), action), "add new author")
		return
	}

	redirect(ctx, fmt.Sprintf("/authors/%d", author.ID))
}
