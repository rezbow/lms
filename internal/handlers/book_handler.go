package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	bookViews "lms/internal/views/books"
	commonViews "lms/internal/views/common"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookHandler struct {
	BookRepo  *repositories.BookRepo
	Validator *validator.Validate
}

func (bh *BookHandler) _get(ctx *gin.Context) (*models.Book, error) {
	bookId := ctx.Param("id")
	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			ctx.HTML(http.StatusNotFound, "", commonViews.NotFound())
			return nil, err
		}
		ctx.HTML(http.StatusInternalServerError, "", commonViews.ServerError(err.Error()))
		return nil, err
	}
	return book, nil

}

func (bh *BookHandler) Get(ctx *gin.Context) {
	bookId := ctx.Param("id")
	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			ctx.HTML(http.StatusNotFound, "", commonViews.NotFound())
			return
		}
		ctx.HTML(http.StatusInternalServerError, "", commonViews.ServerError(err.Error()))
		return
	}
	ctx.HTML(http.StatusOK, "", bookViews.Book(book))
}

func (bh *BookHandler) Delete(ctx *gin.Context) {
	bookId := ctx.Param("id")
	err := bh.BookRepo.DeleteById(bookId)
	if err != nil {
		switch err {
		case repositories.ErrBookHasActiveLoan:
			ctx.HTML(http.StatusConflict, "", commonViews.Flash("this book has active loans and can't be deleted", "red"))
		case repositories.ErrNotFound:
			ctx.HTML(http.StatusNotFound, "", commonViews.NotFound())
		default:
			ctx.HTML(http.StatusInternalServerError, "", commonViews.ServerError(err.Error()))
		}
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/books?flash=deleted")
}

func (bh *BookHandler) AddPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "", bookViews.BookForm())
}

func (bh *BookHandler) Add(ctx *gin.Context) {
	var bookForm struct {
		TitleFa     string `form:"title_fa" binding:"required" validate:"required,min=1,max=100"`
		TitleEn     string `form:"title_en" binding:"required" validate:"required,min=1,max=100"`
		ISBN        string `form:"isbn" binding:"required" validate:"required"`
		TotalCopies int    `form:"total_copies" binding:"required" validate:"required,min=1"`
		AuthorId    int    `form:"author_id" binding:"required" validate:"required,min=1"`
	}

	if err := ctx.ShouldBind(&bookForm); err != nil {
		ctx.HTML(http.StatusBadRequest, "", commonViews.FormErrors([]string{err.Error()}))
		return
	}

	if err := bh.Validator.Struct(bookForm); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fields", "validation_errors": extractValidationErrs(err)})
		return
	}

	book := models.Book{
		TitleFa:         bookForm.TitleFa,
		TitleEn:         bookForm.TitleEn,
		ISBN:            bookForm.ISBN,
		TotalCopies:     bookForm.TotalCopies,
		AvailableCopies: bookForm.TotalCopies,
		AuthorId:        bookForm.AuthorId,
	}

	if err := bh.BookRepo.Insert(&book); err != nil {
		// check for invalid author ID
		if err == repositories.ErrAuthorIdNotFound {
			ctx.HTML(http.StatusBadRequest, "", commonViews.FormErrors([]string{err.Error()}))
			return
		}
		ctx.HTML(http.StatusInternalServerError, "", commonViews.ServerError(err.Error()))
		return
	}

	HXRedirect(ctx, fmt.Sprintf("/books/%d?flash=added!", book.ID))
}

func (bh *BookHandler) Search(ctx *gin.Context) {
	var total int64

	title := ctx.Query("title")
	author := ctx.Query("author")
	isbn := ctx.Query("isbn")

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	limit, err := strconv.Atoi(limitStr)

	books, err := bh.BookRepo.Filter(
		title,
		author,
		isbn,
		limit,
		page,
		&total,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    CodeInternalErr,
				"message": "internal server error",
				"details": err.Error(),
			},
		})
		return
	}

	ctx.HTML(http.StatusOK, "", bookViews.BookList(books))
}

func (bh *BookHandler) Index(ctx *gin.Context) {
	flash := ctx.Query("flash")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	books, err := bh.BookRepo.All(page, pageSize)
	if err != nil {
		ctx.HTML(http.StatusOK, "", commonViews.ServerError(err.Error()))
		return
	}
	ctx.HTML(http.StatusOK, "", bookViews.BookPage(books, flash))
}

func (bh *BookHandler) EditPage(ctx *gin.Context) {
	bookId := ctx.Param("id")
	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			ctx.HTML(http.StatusNotFound, "", commonViews.NotFound())
			return
		}
		ctx.HTML(http.StatusInternalServerError, "", commonViews.ServerError(err.Error()))
		return
	}
	ctx.HTML(http.StatusOK, "", bookViews.BookEditForm(book))
}

func (bh *BookHandler) Update(ctx *gin.Context) {

	bookId := ctx.Param("id")
	id, err := strconv.Atoi(bookId)
	if err != nil {
		ctx.HTML(http.StatusBadRequest, "", commonViews.FormErrors([]string{err.Error()}))
		return
	}

	var bookUpdateForm struct {
		TitleFa     *string `form:"title_fa" validate:"omitempty,min=1,max=100"`
		TitleEn     *string `form:"title_en" validate:"omitempty,min=1,max=100"`
		ISBN        *string `form:"isbn"`
		TotalCopies *int    `form:"total_copies" validate:"omitempty,min=1"`
		AuthorId    *int    `form:"author_id" validate:"min=0"`
	}

	if err := ctx.ShouldBind(&bookUpdateForm); err != nil {
		ctx.HTML(http.StatusBadRequest, "", commonViews.FormErrors([]string{err.Error()}))
		return
	}

	if err := bh.Validator.Struct(&bookUpdateForm); err != nil {
		ctx.HTML(http.StatusBadRequest, "", commonViews.FormErrors([]string{err.Error()}))
		return
	}

	book := models.Book{
		ID: id,
	}
	if bookUpdateForm.TitleEn != nil {
		book.TitleEn = *bookUpdateForm.TitleEn
	}
	if bookUpdateForm.TitleFa != nil {
		book.TitleFa = *bookUpdateForm.TitleFa
	}
	if bookUpdateForm.ISBN != nil {
		book.ISBN = *bookUpdateForm.ISBN
	}
	if bookUpdateForm.TotalCopies != nil {
		book.TotalCopies = *bookUpdateForm.TotalCopies
	}
	if bookUpdateForm.AuthorId != nil {
		book.AuthorId = *bookUpdateForm.AuthorId
	}

	if err := bh.BookRepo.Update(&book); err != nil {
		if err == repositories.ErrAuthorIdNotFound {
			ctx.HTML(http.StatusOK, "", commonViews.FormErrors([]string{err.Error()}))
			return
		}
		ctx.HTML(http.StatusOK, "", commonViews.FormErrors([]string{err.Error()}))
		return
	}

	HXRedirect(ctx, "/books/"+bookId)

}
