package handlers

import (
	"lms/internal/models"
	"lms/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthorHandler struct {
	Repo      repositories.AuthorRepo
	Validator *validator.Validate
}

func (ah *AuthorHandler) Get(ctx *gin.Context) {
	authorId := ctx.Param("id")
	book, err := ah.Repo.GetById(authorId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "author not found"})
		return
	}
	ctx.JSON(http.StatusOK, book)
}

func (ah *AuthorHandler) Delete(ctx *gin.Context) {
	authorId := ctx.Param("id")
	err := ah.Repo.DeleteById(authorId)
	if err != nil {
		if err == repositories.ErrNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "author not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (ah *AuthorHandler) Add(ctx *gin.Context) {
	var author models.Author
	if err := ctx.ShouldBindBodyWithJSON(&author); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if err := ah.Validator.Struct(author); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fields", "validation_errors": extractValidationErrs(err)})
		return
	}

	if err := ah.Repo.Insert(&author); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusCreated, author)

}
