package handlers

import (
	"fmt"
	"lms/internal/models"
	"lms/internal/repositories"
	"lms/internal/views"
	bookViews "lms/internal/views/books"
	loanViews "lms/internal/views/loans"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BookHandler struct {
	BookRepo  *repositories.BookRepo
	Validator *validator.Validate
}

func (bh *BookHandler) Get(ctx *gin.Context) {

	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	render(ctx, bookViews.Book(book), book.TitleFa)
}

func (bh *BookHandler) Delete(ctx *gin.Context) {

	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	err = bh.BookRepo.DeleteById(bookId)
	if err != nil {
		switch err {
		case repositories.ErrBookHasActiveLoan:
			redirect(ctx, fmt.Sprintf("/books/%d", bookId))
			return
		case repositories.ErrNotFound:
			notfound(ctx)
			return
		default:
			serverError(ctx)
			return

		}
	}
	redirect(ctx, "/books")
}

func (bh *BookHandler) AddLoanPage(ctx *gin.Context) {
	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	render(ctx, loanViews.LoanAddForm(views.Data{"bookId": book.ID}), "add loan")

}

func (bh *BookHandler) AddPage(ctx *gin.Context) {
	render(ctx, bookViews.BookAddForm(views.Errors{}), "add book")
}

func (bh *BookHandler) Add(ctx *gin.Context) {
	var bookForm struct {
		TitleFa     string `form:"titleFa" binding:"required" validate:"required,min=1,max=100"`
		TitleEn     string `form:"titleEn" binding:"required" validate:"required,min=1,max=100"`
		ISBN        string `form:"isbn" binding:"required" validate:"required"`
		TotalCopies int    `form:"totalCopies" binding:"required" validate:"required,min=1"`
		AuthorId    uint   `form:"authorId" binding:"required" validate:"required,min=1"`
	}

	if err := ctx.ShouldBind(&bookForm); err != nil {
		render(ctx, bookViews.BookAddForm(parseValidationErrors(err)), "add book")
		return
	}

	if err := bh.Validator.Struct(bookForm); err != nil {
		render(ctx, bookViews.BookAddForm(parseValidationErrors(err)), "add book")
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
			render(ctx, bookViews.BookAddForm(views.Errors{
				"authorId": "doesn't exist",
			}), "add book")
			return
		}
		serverError(ctx)
		return
	}

	redirect(ctx, fmt.Sprintf("/books/%d", book.ID))
}

func (bh *BookHandler) Search(ctx *gin.Context) {
	var total int64

	title := ctx.Query("title")
	author := ctx.Query("author")
	isbn := ctx.Query("isbn")

	pagination, err := readPagination(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	books, err := bh.BookRepo.Filter(
		title,
		author,
		isbn,
		pagination,
		&total,
	)

	if err != nil {
		serverError(ctx)
		return
	}
	render(ctx, bookViews.BookList(books), "search results")
}

func (bh *BookHandler) Index(ctx *gin.Context) {
	pagination, err := readPagination(ctx)
	var total int64

	books, err := bh.BookRepo.Filter("", "", "", pagination, &total)
	if err != nil {
		serverError(ctx)
		return
	}
	render(ctx, bookViews.BookPage(books), "books")
}

func (bh *BookHandler) EditPage(ctx *gin.Context) {
	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}
	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}
	render(ctx, bookViews.BookEditForm(book, views.Errors{}), book.TitleFa)
}

func (bh *BookHandler) Update(ctx *gin.Context) {

	bookId, err := readID(ctx)
	if err != nil {
		notfound(ctx)
		return
	}

	var bookUpdateForm struct {
		TitleFa     *string `form:"titleFa" validate:"omitempty,min=1,max=100"`
		TitleEn     *string `form:"titleEn" validate:"omitempty,min=1,max=100"`
		ISBN        *string `form:"isbn"`
		TotalCopies *int    `form:"totalCopies" validate:"omitempty,min=1"`
		AuthorId    *uint   `form:"authorId" validate:"min=0"`
	}

	book, err := bh.BookRepo.GetById(bookId)
	if err != nil {
		if err == repositories.ErrNotFound {
			notfound(ctx)
			return
		}
		serverError(ctx)
		return
	}

	if err := ctx.ShouldBind(&bookUpdateForm); err != nil {
		render(ctx, bookViews.BookEditForm(book, parseValidationErrors(err)), book.TitleFa)
		return
	}

	if err := bh.Validator.Struct(&bookUpdateForm); err != nil {
		render(ctx, bookViews.BookEditForm(book, parseValidationErrors(err)), book.TitleFa)
		return
	}

	newBook := models.Book{
		ID: bookId,
	}
	if bookUpdateForm.TitleEn != nil {
		newBook.TitleEn = *bookUpdateForm.TitleEn
	}
	if bookUpdateForm.TitleFa != nil {
		newBook.TitleFa = *bookUpdateForm.TitleFa
	}
	if bookUpdateForm.ISBN != nil {
		newBook.ISBN = *bookUpdateForm.ISBN
	}
	if bookUpdateForm.TotalCopies != nil {
		newBook.TotalCopies = *bookUpdateForm.TotalCopies
	}
	if bookUpdateForm.AuthorId != nil {
		newBook.AuthorId = *bookUpdateForm.AuthorId
	}

	if err := bh.BookRepo.Update(&newBook); err != nil {
		if err == repositories.ErrAuthorIdNotFound {
			render(ctx, bookViews.BookEditForm(book, views.Errors{"authorId": "author doesn't exist"}), book.TitleFa)
			return
		}
		serverError(ctx)
		return
	}

	redirect(ctx, fmt.Sprintf("/books/%d", bookId))
}
