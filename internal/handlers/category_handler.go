package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"
	categoryViews "lms/internal/views/categories"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryHandler struct {
	Repo      *repositories.CategoryRepo
	Validator *validator.Validate
}

func (ch *CategoryHandler) AddPage(ctx *gin.Context) {
	render(ctx, categoryViews.CategoryForm(nil, nil, "/categories/add"), "add category")
}

func (ch *CategoryHandler) Add(ctx *gin.Context) {
	var userInput struct {
		Name string `form:"name" binding:"required"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, categoryViews.CategoryForm(&models.Category{Name: userInput.Name}, views.Errors{"name": err.Error()}, "/categories/add"), "add category")
	}

	slug := slugify(userInput.Name)

	category := &models.Category{
		Name: userInput.Name,
		Slug: slug,
	}

	err := ch.Repo.Insert(category)

	if err != nil {
		switch err {
		case repositories.ErrInternal:
			serverError(ctx)
		default:
			render(ctx, categoryViews.CategoryForm(category, views.Errors{"name": err.Error()}, "/categories/add"), "add category")
		}
		return
	}

	redirect(ctx, fmt.Sprintf("/categories/%s", category.Slug))

}

func (ch *CategoryHandler) Get(ctx *gin.Context) {
	slug := ctx.Param("slug")
	category, err := ch.Repo.GetBySlug(slug)
	if err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
		} else if err == repositories.ErrNotFound {
			notfound(ctx)
		}
		return
	}
	render(ctx, categoryViews.Category(category), category.Name)
}

func (ch *CategoryHandler) EditPage(ctx *gin.Context) {
	slug := ctx.Param("slug")
	category, err := ch.Repo.GetBySlug(slug)
	if err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
		} else if err == repositories.ErrNotFound {
			notfound(ctx)
		}
		return
	}
	render(ctx, categoryViews.CategoryForm(category, nil, fmt.Sprintf("/categories/%s/edit", category.Slug)), category.Name)
}

func (ch *CategoryHandler) Edit(ctx *gin.Context) {
	oldSlug := ctx.Param("slug")

	oldCategory, err := ch.Repo.GetBySlug(oldSlug)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	endpoint := fmt.Sprintf("/categories/%s/edit", oldSlug)

	var userInput struct {
		Name string `form:"name" binding:"required"`
	}

	if err := ctx.ShouldBind(&userInput); err != nil {
		render(ctx, categoryViews.CategoryForm(&models.Category{Name: userInput.Name}, views.Errors{"name": err.Error()}, endpoint), "")
		return
	}

	slug := slugify(userInput.Name)

	category := &models.Category{
		ID:   oldCategory.ID,
		Name: userInput.Name,
		Slug: slug,
	}

	err = ch.Repo.Update(category)

	if err != nil {
		switch err {
		case repositories.ErrInternal:
			serverError(ctx)
		case repositories.ErrNotFound:
			notfound(ctx)
		default:
			if err := ctx.ShouldBind(&userInput); err != nil {
				render(ctx, categoryViews.CategoryForm(&models.Category{Name: userInput.Name}, views.Errors{"name": err.Error()}, endpoint), "add category")
				return
			}
		}
		return
	}

	redirect(ctx, fmt.Sprintf("/categories/%s", category.Slug))

}

func (ch *CategoryHandler) Delete(ctx *gin.Context) {
	slug := ctx.Param("slug")
	err := ch.Repo.DeleteBySlug(slug)
	if err != nil {
		if err == repositories.ErrInternal {
			serverError(ctx)
		} else if err == repositories.ErrNotFound {
			notfound(ctx)
		}
		return
	}
	redirect(ctx, "/categories")
}
